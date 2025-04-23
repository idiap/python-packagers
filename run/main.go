// SPDX-FileCopyrightText: © 2025 Idiap Research Institute <contact@idiap.ch>
// SPDX-FileContributor: Samuel Gaist <samuel.gaist@idiap.ch>
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"os"

	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/chronos"
	"github.com/paketo-buildpacks/packit/v2/scribe"
	pythonpackagers "github.com/paketo-buildpacks/python-packagers"
	pkgcommon "github.com/paketo-buildpacks/python-packagers/pkg/common"
)

func main() {
	logger := scribe.NewEmitter(os.Stdout).WithLevel(os.Getenv("BP_LOG_LEVEL"))

	buildParameters := pkgcommon.CommonBuildParameters{
		SbomGenerator: pkgcommon.Generator{},
		Clock:         chronos.DefaultClock,
		Logger:        logger,
	}

	packit.Run(
		pythonpackagers.Detect(logger),
		pythonpackagers.Build(logger, buildParameters),
	)
}
