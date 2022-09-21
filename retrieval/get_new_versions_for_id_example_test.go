package retrieval_test

import (
	"path/filepath"

	"github.com/Masterminds/semver/v3"
	"github.com/joshuatcasey/libdependency/buildpack_config"
	"github.com/joshuatcasey/libdependency/retrieval"
	"github.com/joshuatcasey/libdependency/versionology"
)

func ExampleGetNewVersionsForId() {
	config, _ := buildpack_config.ParseBuildpackToml(filepath.Join("..", "retrieval", "testdata", "happy_path", "buildpack.toml"))
	getAllVersions := func() (versionology.VersionFetcherArray, error) {
		return versionology.VersionFetcherArray{
			versionology.NewSimpleVersionFetcher(semver.MustParse("1.0.0")),
			versionology.NewSimpleVersionFetcher(semver.MustParse("1.1.0")),
			versionology.NewSimpleVersionFetcher(semver.MustParse("1.2.0")),
			versionology.NewSimpleVersionFetcher(semver.MustParse("1.3.0")),
			versionology.NewSimpleVersionFetcher(semver.MustParse("1.4.0")),
			versionology.NewSimpleVersionFetcher(semver.MustParse("1.5.0")),
		}, nil
	}
	_, _ = retrieval.GetNewVersionsForId("fake-dependency-id", config, getAllVersions)

	// Output:
	// Found 6 versions of fake-dependency-id from upstream
	// [
	//   "1.5.0", "1.4.0", "1.3.0", "1.2.0", "1.1.0",
	//   "1.0.0"
	// ]
	// Found 6 versions of fake-dependency-id for constraint 1.*.*
	// [
	//   "1.5.0", "1.4.0", "1.3.0", "1.2.0", "1.1.0",
	//   "1.0.0"
	// ]
	// Found 2 versions of fake-dependency-id newer than '1.1.0' for constraint 1.*.*, after limiting for 2 patches
	// [
	//   "1.5.0", "1.4.0"
	// ]
	// Found 2 versions of fake-dependency-id as new versions
	// [
	//   "1.5.0", "1.4.0"
	// ]
}
