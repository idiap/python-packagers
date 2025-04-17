package pythonpackagers

import (
	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/draft"
	"github.com/paketo-buildpacks/packit/v2/fs"
	"github.com/paketo-buildpacks/packit/v2/pexec"
	"github.com/paketo-buildpacks/packit/v2/scribe"

	conda "github.com/paketo-buildpacks/python-packagers/pkg/conda"
	pipinstall "github.com/paketo-buildpacks/python-packagers/pkg/pip"
	pipenvinstall "github.com/paketo-buildpacks/python-packagers/pkg/pipenv"
	poetryinstall "github.com/paketo-buildpacks/python-packagers/pkg/poetry"

	"github.com/paketo-buildpacks/python-packagers/pkg/common"
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

func Build(logger scribe.Emitter, commonBuildParameters pythonpackagers.CommonBuildParameters) packit.BuildFunc {
	return func(context packit.BuildContext) (packit.BuildResult, error) {
		planEntries := filtered(context.Plan.Entries, pipinstall.SitePackages)

		for _, entry := range planEntries {
			logger.Title("Handling %s", entry.Name)

			switch entry.Name {
			case pipinstall.Manager:
				pipResult, err := pipinstall.Build(
					pipinstall.NewPipInstallProcess(pexec.NewExecutable("pip"), logger),
					pipinstall.NewSiteProcess(pexec.NewExecutable("python")),
					commonBuildParameters,
				)(context)

				if err != nil {
					return packit.BuildResult{}, err
				}

				return pipResult, err
			case pipenvinstall.Manager:
				pipEnvResult, err := pipenvinstall.Build(
					pipenvinstall.NewPipenvInstallProcess(pexec.NewExecutable("pipenv"), logger),
					pipenvinstall.NewSiteProcess(pexec.NewExecutable("python")),
					pipenvinstall.NewVenvLocator(),
					commonBuildParameters,
				)(context)

				if err != nil {
					return packit.BuildResult{}, err
				}

				return pipEnvResult, err
			case conda.CondaEnvPlanEntry:
				condaResult, err := conda.Build(
					conda.NewCondaRunner(pexec.NewExecutable("conda"), fs.NewChecksumCalculator(), logger),
					commonBuildParameters,
				)(context)

				if err != nil {
					return packit.BuildResult{}, err
				}

				return condaResult, err
			case poetryinstall.PoetryVenv:
				poetryResult, err := poetryinstall.Build(
					draft.NewPlanner(),
					poetryinstall.NewPoetryInstallProcess(pexec.NewExecutable("poetry"), logger),
					poetryinstall.NewPythonPathProcess(),
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
