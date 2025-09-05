package retrieve_test

import (
	"testing"

	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestUnitRetrieve(t *testing.T) {
	suite := spec.New("retrieve", spec.Report(report.Terminal{}))
	suite("NewMetadata", testNewMetadata, spec.Sequential())
	suite("NewMetadataWithPlatforms", testNewMetadataWithPlatforms, spec.Sequential())
	suite("GetNewVersionsForId", testGetNewVersionsForId, spec.Sequential())
	suite("purl", testPurl)
	suite.Run(t)
}
