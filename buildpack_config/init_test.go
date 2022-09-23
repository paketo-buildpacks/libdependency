package buildpack_config_test

import (
	"testing"

	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestUnitFuncs(t *testing.T) {
	suite := spec.New("libdependency", spec.Report(report.Terminal{}))
	suite("buildpackToml", testBuildpackToml)
	suite.Run(t)
}
