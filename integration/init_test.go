package integration_test

import (
	"testing"

	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestIntegration(t *testing.T) {
	suite := spec.New("integration", spec.Report(report.Terminal{}))
	suite("GithubReleases", testGithubReleases)
	suite("Retrieve", testRetrieve)
	suite.Run(t)
}
