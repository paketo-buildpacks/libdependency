package retrieval_test

import (
	"testing"

	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestUnitRetrieval(t *testing.T) {
	suite := spec.New("retrieval", spec.Report(report.Terminal{}))
	suite("NewMetadata", testNewMetadata, spec.Sequential())
	suite("GetNewVersionsForId", testGetNewVersionsForId, spec.Sequential())
	suite("purl", testPurl)
	suite.Run(t)
}
