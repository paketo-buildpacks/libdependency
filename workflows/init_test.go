package workflows_test

import (
	"testing"

	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestUnitFuncs(t *testing.T) {
	suite := spec.New("workflows", spec.Report(report.Terminal{}))
	suite("Workflows", testWorkflows)
	suite.Run(t)
}
