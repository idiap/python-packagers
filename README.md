<!--
SPDX-FileCopyrightText: © 2025 Idiap Research Institute <contact@idiap.ch>
SPDX-FileContributor: Samuel Gaist <samuel.gaist@idiap.ch>

SPDX-License-Identifier: Apache-2.0
-->

# Python Packagers Cloud Native Buildpack

The Paketo Buildpack for Python Packagers is a Cloud Native Buildpack that
installs packages using the adequate tool selected based on the contect of the
application sources and makes it available to it.

The buildpack is published for consumption at
`gcr.io/paketo-buildpacks/python-packagers` and
`paketobuildpacks/python-packageres`.

## Behavior
This buildpack participates if `requirements.txt` exists at the root the app.

The buildpack will do the following:
* At build time:
  - Installs the application packages to a layer made available to the app.
* At run time:
  - Does nothing

## Usage

To package this buildpack for consumption:
```
$ ./scripts/package.sh --version x.x.x
```
This will create a `buildpackage.cnb` file under the build directory which you
can use to build your app as follows: `pack build <app-name> -p <path-to-app>
-b <cpython buildpack> -b <pip buildpack> -b build/buildpackage.cnb -b
<other-buildpacks..>`.

To run the unit and integration tests for this buildpack:
```
$ ./scripts/unit.sh && ./scripts/integration.sh
```
