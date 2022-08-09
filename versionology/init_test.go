package versionology_test

import (
	"testing"

	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestUnitBundler(t *testing.T) {
	suite := spec.New("versionology", spec.Report(report.Terminal{}))
	suite("Versionology", testVersionology)
	suite.Run(t)
}
