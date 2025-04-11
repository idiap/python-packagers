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
)

type Generator struct{}

func (f Generator) Generate(dir string) (sbom.SBOM, error) {
	return sbom.Generate(dir)
}

func Build(logger scribe.Emitter) packit.BuildFunc {
	return func(context packit.BuildContext) (packit.BuildResult, error) {

		for _, entry := range context.Plan.Entries {
			logger.Title("Handling %s", entry.Name)
			switch entry.Name {
			case pipinstall.SitePackages:
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
				/* code */
				return packit.BuildResult{}, packit.Fail.WithMessage("unknown plan: %s", entry.Name)
			}
		}

		return packit.BuildResult{}, packit.Fail.WithMessage("empty plan should not happen")
	}
}
