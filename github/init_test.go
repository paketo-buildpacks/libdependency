package github_test

import (
	"testing"

	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestUnitFuncs(t *testing.T) {
	suite := spec.New("github", spec.Report(report.Terminal{}))
	suite("Github", testGithub)
	suite.Run(t)
}
