#!/usr/bin/env bash

# Copyright (c) 2013-Present CloudFoundry.org Foundation, Inc. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0

set -eu
set -o pipefail

readonly PROGDIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
readonly BUILDPACKDIR="$(cd "${PROGDIR}/.." && pwd)"

function main() {
  while [[ "${#}" != 0 ]]; do
    case "${1}" in
      --help|-h)
        shift 1
        usage
        exit 0
        ;;

      "")
        # skip if the argument is empty
        shift 1
        ;;

      *)
        util::print::error "unknown argument \"${1}\""
    esac
  done

  mkdir -p "${BUILDPACKDIR}/bin"

  run::build
  cmd::build
}

function usage() {
  cat <<-USAGE
build.sh [OPTIONS]

Builds the buildpack executables.

OPTIONS
  --help  -h  prints the command usage
USAGE
}

function run::build() {
  if [[ -f "${BUILDPACKDIR}/run/main.go" ]]; then
    pushd "${BUILDPACKDIR}/bin" > /dev/null || return
      printf "%s" "Building run... "

      GOOS=linux \
      CGO_ENABLED=0 \
        go build \
          -ldflags="-s -w" \
          -o "run" \
            "${BUILDPACKDIR}/run"

      echo "Success!"

      names=("detect")

      if [ -f "${BUILDPACKDIR}/extension.toml" ]; then
        names+=("generate")
      else
        names+=("build")
      fi

      for name in "${names[@]}"; do
        printf "%s" "Linking ${name}... "

        ln -sf "run" "${name}"

        echo "Success!"
      done
    popd > /dev/null || return
  fi
}

function cmd::build() {
  if [[ -d "${BUILDPACKDIR}/cmd" ]]; then
    local name
    for src in "${BUILDPACKDIR}"/cmd/*; do
      name="$(basename "${src}")"

      if [[ -f "${src}/main.go" ]]; then
        printf "%s" "Building ${name}... "

        GOOS="linux" \
        CGO_ENABLED=0 \
          go build \
            -ldflags="-s -w" \
            -o "${BUILDPACKDIR}/bin/${name}" \
              "${src}/main.go"

        echo "Success!"
      else
        printf "%s" "Skipping ${name}... "
      fi
    done
  fi
}

main "${@:-}"
