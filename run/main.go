package main

import (
	"os"

	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/chronos"
	"github.com/paketo-buildpacks/packit/v2/scribe"
	"github.com/paketo-buildpacks/python-packagers"
	pkgcommon "github.com/paketo-buildpacks/python-packagers/pkg/common"
)

func main() {
	logger := scribe.NewEmitter(os.Stdout).WithLevel(os.Getenv("BP_LOG_LEVEL"))

	buildParameters := pkgcommon.CommonBuildParameters{
		pkgcommon.Generator{},
		chronos.DefaultClock,
		logger,
	}

	packit.Run(
		pythonpackagers.Detect(logger),
		pythonpackagers.Build(logger, buildParameters),
	)
}
