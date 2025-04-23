// SPDX-FileCopyrightText: © 2025 Idiap Research Institute <contact@idiap.ch>
// SPDX-FileContributor: Samuel Gaist <samuel.gaist@idiap.ch>
//
// SPDX-License-Identifier: Apache-2.0

package pythonpackagers

import (
	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/scribe"

	conda "github.com/paketo-buildpacks/python-packagers/pkg/conda"
	pipinstall "github.com/paketo-buildpacks/python-packagers/pkg/pip"
	pipenvinstall "github.com/paketo-buildpacks/python-packagers/pkg/pipenv"
	poetryinstall "github.com/paketo-buildpacks/python-packagers/pkg/poetry"

	pythonpackagers "github.com/paketo-buildpacks/python-packagers/pkg/common"
)

// filtered returns the slice passed in parameter with the needle removed
func filtered(haystack []packit.BuildpackPlanEntry, needle string) []packit.BuildpackPlanEntry {
	output := []packit.BuildpackPlanEntry{}

	for _, entry := range haystack {
		if entry.Name != needle {
			output = append(output, entry)
		}
	}

	return output
}

type PackagerParameters interface {
}

func Build(
	logger scribe.Emitter,
	commonBuildParameters pythonpackagers.CommonBuildParameters,
	buildParameters map[string]PackagerParameters,
) packit.BuildFunc {
	return func(context packit.BuildContext) (packit.BuildResult, error) {
		planEntries := filtered(context.Plan.Entries, pipinstall.SitePackages)

		for _, entry := range planEntries {
			logger.Title("Handling %s", entry.Name)

			switch entry.Name {
			case pipinstall.Manager:
				pipParameters := buildParameters[pipinstall.Manager].(pipinstall.PipBuildParameters)
				pipResult, err := pipinstall.Build(
					pipParameters.InstallProcess,
					pipParameters.SitePackagesProcess,
					commonBuildParameters,
				)(context)

				if err != nil {
					return packit.BuildResult{}, err
				}

				return pipResult, err
			case pipenvinstall.Manager:
				pipenvParameters := buildParameters[pipenvinstall.Manager].(pipenvinstall.PipEnvBuildParameters)
				pipEnvResult, err := pipenvinstall.Build(
					pipenvParameters.InstallProcess,
					pipenvParameters.SiteProcess,
					pipenvParameters.VenvDirLocator,
					commonBuildParameters,
				)(context)

				if err != nil {
					return packit.BuildResult{}, err
				}

				return pipEnvResult, err
			case conda.CondaEnvPlanEntry:
				condaParameters := buildParameters[conda.CondaEnvPlanEntry].(conda.CondaBuildParameters)
				condaResult, err := conda.Build(
					condaParameters.Runner,
					commonBuildParameters,
				)(context)

				if err != nil {
					return packit.BuildResult{}, err
				}

				return condaResult, err
			case poetryinstall.PoetryVenv:
				poetryParameters := buildParameters[poetryinstall.PoetryVenv].(poetryinstall.PoetryEnvBuildParameters)
				poetryResult, err := poetryinstall.Build(
					poetryParameters.EntryResolver,
					poetryParameters.InstallProcess,
					poetryParameters.PythonPathLookupProcess,
					commonBuildParameters,
				)(context)

				if err != nil {
					return packit.BuildResult{}, err
				}

				return poetryResult, err
			default:
				return packit.BuildResult{}, packit.Fail.WithMessage("unknown plan: %s", entry.Name)
			}
		}

		return packit.BuildResult{}, packit.Fail.WithMessage("empty plan should not happen")
	}
}
