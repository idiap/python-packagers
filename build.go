package pythonpackagers

import (
	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/chronos"
	"github.com/paketo-buildpacks/packit/v2/fs"
	"github.com/paketo-buildpacks/packit/v2/pexec"
	"github.com/paketo-buildpacks/packit/v2/sbom"
	"github.com/paketo-buildpacks/packit/v2/scribe"

	conda "github.com/paketo-buildpacks/python-packagers/pkg/conda"
	pipinstall "github.com/paketo-buildpacks/python-packagers/pkg/pip"
	pipenvinstall "github.com/paketo-buildpacks/python-packagers/pkg/pipenv"
)

type Generator struct{}

func (f Generator) Generate(dir string) (sbom.SBOM, error) {
	return sbom.Generate(dir)
}

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

func Build(logger scribe.Emitter) packit.BuildFunc {
	return func(context packit.BuildContext) (packit.BuildResult, error) {
		planEntries := filtered(context.Plan.Entries, pipinstall.SitePackages)

		for _, entry := range planEntries {
			logger.Title("Handling %s", entry.Name)

			switch entry.Name {
			case pipinstall.Manager:
				pipResult, err := pipinstall.Build(
					pipinstall.NewPipInstallProcess(pexec.NewExecutable("pip"), logger),
					pipinstall.NewSiteProcess(pexec.NewExecutable("python")),
					Generator{},
					chronos.DefaultClock,
					logger,
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
					Generator{},
					chronos.DefaultClock,
					logger,
				)(context)

				if err != nil {
					return packit.BuildResult{}, err
				}

				return pipEnvResult, err
			case conda.CondaEnvPlanEntry:
				condaResult, err := conda.Build(
					conda.NewCondaRunner(pexec.NewExecutable("conda"), fs.NewChecksumCalculator(), logger),
					Generator{},
					logger,
					chronos.DefaultClock,
				)(context)

				if err != nil {
					return packit.BuildResult{}, err
				}

				return condaResult, err
			default:
				return packit.BuildResult{}, packit.Fail.WithMessage("unknown plan: %s", entry.Name)
			}
		}

		return packit.BuildResult{}, packit.Fail.WithMessage("empty plan should not happen")
	}
}
