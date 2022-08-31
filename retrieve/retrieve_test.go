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
		savedArgs  []string
		outputFile string
	)

	it.Before(func() {
		outputFile = filepath.Join(t.TempDir(), "temp.json")

		savedArgs = os.Args
		t.Cleanup(func() {
			os.Args = savedArgs
		})

		os.Args = []string{"/path/to-binary",
			"--buildpack-toml-path", filepath.Join("..", "testdata", "empty", "buildpack.toml"),
			"--output-file", outputFile}
	})

	context("RetrieveNewMetadata", func() {
		it("will write the output to the given location", func() {
			getNewVersions := func() ([]versionology.HasVersion, error) {
				return []versionology.HasVersion{
					versionology.NewSimpleHasVersion(semver.MustParse("1.1.1")),
					versionology.NewSimpleHasVersion(semver.MustParse("2.2.2")),
				}, nil
			}

			generateMetadata := func(version versionology.HasVersion) (cargo.ConfigMetadataDependency, error) {
				return cargo.ConfigMetadataDependency{
					Version: version.GetVersion().String(),
				}, nil
			}

			retrieve.NewMetadata("id", getNewVersions, generateMetadata)
			Expect(outputFile).To(matchers.BeAFileWithContents(`[{"version":"2.2.2"},{"version":"1.1.1"}]`))
		})
	})
}
