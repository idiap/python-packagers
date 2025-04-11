package pythonpackagers

import (
	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/scribe"

	conda "github.com/paketo-buildpacks/python-packagers/pkg/conda"
	pipinstall "github.com/paketo-buildpacks/python-packagers/pkg/pip"
)

// Detect will return a packit.DetectFunc that will be invoked during the
// detect phase of the buildpack lifecycle.
//
// If this buildpack detects files that indicate your app is a Python project,
// it will pass detection.
func Detect(logger scribe.Emitter) packit.DetectFunc {
	return func(context packit.DetectContext) (packit.DetectResult, error) {
		plans := []packit.BuildPlan{}

		pipResult, err := pipinstall.Detect()(context)

		if err == nil {
			plans = append(plans, pipResult.Plan)
		} else {
			logger.Detail("%s", err)
		}

		condaResult, err := conda.Detect()(context)

		if err == nil {
			plans = append(plans, condaResult.Plan)
		} else {
			logger.Detail("%s", err)
		}

		if len(plans) == 0 {
			return packit.DetectResult{}, packit.Fail.WithMessage("No python packager manager related files found")
		}

		return packit.DetectResult{
			Plan: or(plans...),
		}, nil
	}
}

func or(plans ...packit.BuildPlan) packit.BuildPlan {
	if len(plans) < 1 {
		return packit.BuildPlan{}
	}
	combinedPlan := plans[0]

	for i := range plans {
		if i == 0 {
			continue
		}
		combinedPlan.Or = append(combinedPlan.Or, plans[i])
	}
	return combinedPlan
}
