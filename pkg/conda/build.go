// SPDX-FileCopyrightText: Copyright (c) 2013-Present CloudFoundry.org Foundation, Inc. All Rights Reserved.
//
// SPDX-License-Identifier: Apache-2.0

package condaenvupdate

import (
	"os"
	"time"

	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/draft"
	"github.com/paketo-buildpacks/packit/v2/fs"
	"github.com/paketo-buildpacks/packit/v2/sbom"

	pythonpackagers "github.com/paketo-buildpacks/python-packagers/pkg/common"
)

//go:generate faux --interface Runner --output fakes/runner.go
//go:generate faux --interface SBOMGenerator --output fakes/sbom_generator.go

// Runner defines the interface for setting up the conda environment.
type Runner interface {
	Execute(condaEnvPath string, condaCachePath string, workingDir string) error
	ShouldRun(workingDir string, metadata map[string]interface{}) (bool, string, error)
}

// CondaBuildParameters encapsulates the conda specific parameters for the
// Build function
type CondaBuildParameters struct {
	Runner Runner
}

// Build will return a packit.BuildFunc that will be invoked during the build
// phase of the buildpack lifecycle.
//
// Build updates the conda environment and stores the result in a layer. It may
// reuse the environment layer from a previous build, depending on conditions
// determined by the runner.
func Build(
	buildParameters CondaBuildParameters,
	parameters pythonpackagers.CommonBuildParameters,
) packit.BuildFunc {
	return func(context packit.BuildContext) (packit.BuildResult, error) {
		runner := buildParameters.Runner
		sbomGenerator := parameters.SbomGenerator
		clock := parameters.Clock
		logger := parameters.Logger

		logger.Title("%s %s", context.BuildpackInfo.Name, context.BuildpackInfo.Version)

		condaLayer, err := context.Layers.Get(CondaEnvLayer)
		if err != nil {
			return packit.BuildResult{}, err
		}

		condaCacheLayer, err := context.Layers.Get(CondaEnvCache)
		if err != nil {
			return packit.BuildResult{}, err
		}

		run, sha, err := runner.ShouldRun(context.WorkingDir, condaLayer.Metadata)
		if err != nil {
			return packit.BuildResult{}, err
		}

		if run {
			condaLayer, err = condaLayer.Reset()
			if err != nil {
				return packit.BuildResult{}, err
			}

			logger.Process("Executing build process")

			duration, err := clock.Measure(func() error {
				return runner.Execute(condaLayer.Path, condaCacheLayer.Path, context.WorkingDir)
			})
			if err != nil {
				return packit.BuildResult{}, err
			}

			logger.Action("Completed in %s", duration.Round(time.Millisecond))
			logger.Break()

			logger.GeneratingSBOM(condaLayer.Path)

			var sbomContent sbom.SBOM
			duration, err = clock.Measure(func() error {
				sbomContent, err = sbomGenerator.Generate(context.WorkingDir)
				return err
			})
			if err != nil {
				return packit.BuildResult{}, err
			}
			logger.Action("Completed in %s", duration.Round(time.Millisecond))
			logger.Break()

			logger.FormattingSBOM(context.BuildpackInfo.SBOMFormats...)

			condaLayer.SBOM, err = sbomContent.InFormats(context.BuildpackInfo.SBOMFormats...)
			if err != nil {
				return packit.BuildResult{}, err
			}

			condaLayer.Metadata = map[string]interface{}{
				"lockfile-sha": sha,
			}
		} else {
			logger.Process("Reusing cached layer %s", condaLayer.Path)
			logger.Break()
		}

		planner := draft.NewPlanner()
		condaLayer.Launch, condaLayer.Build = planner.MergeLayerTypes(CondaEnvPlanEntry, context.Plan.Entries)
		condaLayer.Cache = condaLayer.Build
		condaCacheLayer.Cache = true

		layers := []packit.Layer{condaLayer}
		if _, err := os.Stat(condaCacheLayer.Path); err == nil {
			if !fs.IsEmptyDir(condaCacheLayer.Path) {
				layers = append(layers, condaCacheLayer)
			}
		}

		return packit.BuildResult{
			Layers: layers,
		}, nil
	}
}
