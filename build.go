package pythonpackagers

import (
	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/chronos"
	"github.com/paketo-buildpacks/packit/v2/pexec"
	"github.com/paketo-buildpacks/packit/v2/sbom"
	"github.com/paketo-buildpacks/packit/v2/scribe"

	pipinstall "github.com/paketo-buildpacks/python-packagers/pkg/pip"
)

type Generator struct{}

func (f Generator) Generate(dir string) (sbom.SBOM, error) {
	return sbom.Generate(dir)
}

func Build(logger scribe.Emitter) packit.BuildFunc {
	return func(context packit.BuildContext) (packit.BuildResult, error) {

		pipResult, err := pipinstall.Build(
			pipinstall.NewPipInstallProcess(pexec.NewExecutable("pip"), logger),
			pipinstall.NewSiteProcess(pexec.NewExecutable("python")),
			Generator{},
			chronos.DefaultClock,
			logger,
		)(context)


		return pipResult, err
	}
}
