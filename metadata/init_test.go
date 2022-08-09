package metadata_test

import (
	"testing"

	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestUnitBundler(t *testing.T) {
	suite := spec.New("metadata", spec.Report(report.Terminal{}))
	suite("metadata", testMetadata)
	suite.Run(t)
}
