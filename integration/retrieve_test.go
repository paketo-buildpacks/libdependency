package integration_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/joshuatcasey/libdependency/matchers"
	"github.com/joshuatcasey/libdependency/retrieve"
	"github.com/joshuatcasey/libdependency/versionology"
	. "github.com/onsi/gomega"
	"github.com/paketo-buildpacks/packit/v2/cargo"
	"github.com/sclevine/spec"
)

func testRetrieve(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect            = NewWithT(t).Expect
		buildpackTomlPath string
		output            string
	)

	it.Before(func() {
		buildpackTomlPath = filepath.Join("testdata", "retrieve", "happy_path", "buildpack.toml")
		output = filepath.Join(t.TempDir(), "metadata.json")
	})

	context("given fake versions and fake metadata", func() {
		it.Before(func() {
			savedArgs := os.Args
			t.Cleanup(func() {
				os.Args = savedArgs
			})

			os.Args = []string{
				"path/to/exe",
				"--buildpackTomlPath", buildpackTomlPath,
				"--output", output,
			}
		})

		it("should generate metadata.json in the output dir", func() {
			getVersions := func() (versionology.VersionFetcherArray, error) {
				versions := versionology.VersionFetcherArray{
					versionology.NewSimpleVersionFetcher(semver.MustParse("1.0.0")),
					versionology.NewSimpleVersionFetcher(semver.MustParse("1.1.0")),
					versionology.NewSimpleVersionFetcher(semver.MustParse("1.2.0")),
					versionology.NewSimpleVersionFetcher(semver.MustParse("1.3.0")),
					versionology.NewSimpleVersionFetcher(semver.MustParse("1.4.0")),
					versionology.NewSimpleVersionFetcher(semver.MustParse("1.5.0")),
				}
				return versions, nil
			}

			getMetadata := func(versionFetcher versionology.VersionFetcher) (versionology.Dependency, error) {
				switch version := versionFetcher.Version().String(); version {
				case "1.4.0":
					fallthrough
				case "1.5.0":
					dependency := cargo.ConfigMetadataDependency{
						ID:      "fake-dependency-id",
						Stacks:  []string{"jammy-stack", "bionic-stack"},
						Version: version,
					}

					return versionology.NewDependency(dependency, "linux-64")
				default:
					return versionology.Dependency{}, errors.New("unexpected version")
				}
			}

			retrieve.NewMetadata("fake-dependency-id", getVersions, getMetadata)

			Expect(output).To(matchers.BeAFileWithContents(`[{"id":"fake-dependency-id","stacks":["jammy-stack","bionic-stack"],"version":"1.5.0","target":"linux-64"},{"id":"fake-dependency-id","stacks":["jammy-stack","bionic-stack"],"version":"1.4.0","target":"linux-64"}]`))
		})
	})
}
