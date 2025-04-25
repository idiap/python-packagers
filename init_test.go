package pythonpackagers_test

import (
	"testing"

	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestUnitPythonPackagers(t *testing.T) {
	suite := spec.New("python-packagers", spec.Report(report.Terminal{}), spec.Sequential())
	suite("Detect", testDetect)
	//	suite("Build", testBuild)
	suite.Run(t)
}
