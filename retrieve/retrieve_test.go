package retrieve_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/joshuatcasey/libdependency/matchers"
	"github.com/joshuatcasey/libdependency/retrieve"
	"github.com/joshuatcasey/libdependency/versionology"
	"github.com/paketo-buildpacks/packit/v2/cargo"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
)

func testRetrieve(t *testing.T, context spec.G, it spec.S) {
	Expect := NewWithT(t).Expect

	var (
		savedArgs []string
		output    string
	)

	it.Before(func() {
		output = filepath.Join(t.TempDir(), "metadata.json")

		savedArgs = os.Args
		t.Cleanup(func() {
			os.Args = savedArgs
		})

		os.Args = []string{"/path/to-binary",
			"--buildpack-toml-path", filepath.Join("..", "testdata", "empty", "buildpack.toml"),
			"--output", output}
	})

	context("RetrieveNewMetadata", func() {
		it("will write the output to the given location", func() {
			getNewVersions := func() (versionology.VersionFetcherArray, error) {
				return versionology.VersionFetcherArray{
					versionology.NewSimpleVersionFetcher(semver.MustParse("1.1.1")),
					versionology.NewSimpleVersionFetcher(semver.MustParse("2.2.2")),
				}, nil
			}

			generateMetadata := func(version versionology.VersionFetcher) (versionology.Dependency, error) {
				return versionology.NewDependency(cargo.ConfigMetadataDependency{
					Version: version.Version().String(),
				}, "target0")
			}

			retrieve.NewMetadata("id", getNewVersions, generateMetadata)
			Expect(output).To(matchers.BeAFileWithContents(
				`[{"version":"2.2.2","target":"target0"},{"version":"1.1.1","target":"target0"}]`))
		})
	})
}
