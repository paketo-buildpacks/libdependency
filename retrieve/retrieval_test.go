package retrieve_test

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/Masterminds/semver/v3"
	. "github.com/onsi/gomega"
	"github.com/paketo-buildpacks/libdependency/buildpack_config"
	"github.com/paketo-buildpacks/libdependency/retrieve"
	"github.com/paketo-buildpacks/libdependency/versionology"
	"github.com/paketo-buildpacks/occam/matchers"
	"github.com/paketo-buildpacks/packit/v2/cargo"
	"github.com/sclevine/spec"
)

func testNewMetadata(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect
		output string

		getAllVersions   retrieve.GetAllVersionsFunc
		generateMetadata retrieve.GenerateMetadataFunc
	)

	it.Before(func() {
		output = filepath.Join(t.TempDir(), "metadata.json")
	})

	context("given fake versions and fake metadata", func() {
		it.Before(func() {
			retrieve.FetchArgs = func() (string, string) {
				buildpackTomlPath := filepath.Join("testdata", "happy_path", "buildpack.toml")
				return buildpackTomlPath, output
			}

			getAllVersions = func() (versionology.VersionFetcherArray, error) {
				return versionology.VersionFetcherArray{
					versionology.NewSimpleVersionFetcher(semver.MustParse("1.0.0")),
					versionology.NewSimpleVersionFetcher(semver.MustParse("1.1.0")),
					versionology.NewSimpleVersionFetcher(semver.MustParse("1.2.0")),
					versionology.NewSimpleVersionFetcher(semver.MustParse("1.3.0")),
					versionology.NewSimpleVersionFetcher(semver.MustParse("1.4.0")),
					versionology.NewSimpleVersionFetcher(semver.MustParse("1.5.0")),
				}, nil
			}

			generateMetadata = func(versionFetcher versionology.VersionFetcher) ([]versionology.Dependency, error) {
				dependency := cargo.ConfigMetadataDependency{
					ID:      "fake-dependency-id",
					Stacks:  []string{"jammy-stack", "bionic-stack"},
					Version: versionFetcher.Version().String(),
				}

				return versionology.NewDependencyArray(dependency, "linux-64")
			}
		})

		it("should generate metadata.json in the output dir", func() {
			retrieve.NewMetadata("fake-dependency-id", getAllVersions, generateMetadata)

			Expect(output).To(matchers.BeAFileMatching(MatchJSON(`
[
	{"id":"fake-dependency-id","stacks":["jammy-stack","bionic-stack"],"version":"1.5.0","target":"linux-64"},
	{"id":"fake-dependency-id","stacks":["jammy-stack","bionic-stack"],"version":"1.4.0","target":"linux-64"}
]`)))
		})
	})

	context("cpython", func() {
		it.Before(func() {
			retrieve.FetchArgs = func() (string, string) {
				buildpackTomlPath := filepath.Join("testdata", "cpython-de13b843", "buildpack.toml")
				return buildpackTomlPath, output
			}

			getAllVersions = func() (versionology.VersionFetcherArray, error) {
				return versionology.NewSimpleVersionFetcherArray("3.10.5", "3.10.6", "3.10.7", "3.10.8", "3.10.9", "3.10.10")
			}

			generateMetadata = func(versionFetcher versionology.VersionFetcher) ([]versionology.Dependency, error) {
				dependency := cargo.ConfigMetadataDependency{
					ID:      "python",
					Stacks:  []string{"io.buildpacks.stacks.jammy"},
					Version: versionFetcher.Version().String(),
				}

				return versionology.NewDependencyArray(dependency, "jammy")
			}
		})

		it("should only generate metadata for the two newest versions", func() {
			retrieve.NewMetadata("python", getAllVersions, generateMetadata)

			Expect(output).To(matchers.BeAFileMatching(MatchJSON(`
[
	{"id":"python","stacks":["io.buildpacks.stacks.jammy"],"version":"3.10.10","target":"jammy"},
	{"id":"python","stacks":["io.buildpacks.stacks.jammy"],"version":"3.10.9","target":"jammy"}
]`)))
		})
	})

	context("when the dependency id is not found in buildpack.toml", func() {
		it.Before(func() {
			retrieve.FetchArgs = func() (string, string) {
				buildpackTomlPath := filepath.Join("testdata", "happy_path", "buildpack.toml")
				return buildpackTomlPath, output
			}

			getAllVersions = func() (versionology.VersionFetcherArray, error) {
				return versionology.VersionFetcherArray{
					versionology.NewSimpleVersionFetcher(semver.MustParse("999.888.777")),
					versionology.NewSimpleVersionFetcher(semver.MustParse("666.555.444")),
					versionology.NewSimpleVersionFetcher(semver.MustParse("333.222.111")),
				}, nil
			}

			generateMetadata = func(versionFetcher versionology.VersionFetcher) ([]versionology.Dependency, error) {
				dependency := cargo.ConfigMetadataDependency{
					ID:      "not a real dependency id",
					Version: versionFetcher.Version().String(),
				}

				return versionology.NewDependencyArray(dependency, "jammy")
			}
		})

		it("all versions are found in the metadata.json", func() {
			retrieve.NewMetadata("not a real dependency id", getAllVersions, generateMetadata)

			Expect(output).To(matchers.BeAFileMatching(MatchJSON(`
[
	{"id":"not a real dependency id","version":"999.888.777","target":"jammy"},
	{"id":"not a real dependency id","version":"666.555.444","target":"jammy"},
	{"id":"not a real dependency id","version":"333.222.111","target":"jammy"}
]`)))
		})
	})

	context("GetNewVersionsForId", func() {
		it("will get new versions for id and stack", func() {
			config, err := buildpack_config.ParseBuildpackToml(filepath.Join("testdata", "bundler", "buildpack.toml"))
			Expect(err).NotTo(HaveOccurred())

			newVersions, err := retrieve.GetNewVersionsForId(
				"bundler",
				config,
				func() (versionology.VersionFetcherArray, error) {
					return versionology.VersionFetcherArray{
						versionology.NewSimpleVersionFetcher(semver.MustParse("0.1.1")),
						versionology.NewSimpleVersionFetcher(semver.MustParse("1.17.3")),
						versionology.NewSimpleVersionFetcher(semver.MustParse("1.17.4")),
						versionology.NewSimpleVersionFetcher(semver.MustParse("2.3.16")),
						versionology.NewSimpleVersionFetcher(semver.MustParse("2.3.17")),
						versionology.NewSimpleVersionFetcher(semver.MustParse("2.4.0")),
						versionology.NewSimpleVersionFetcher(semver.MustParse("3.0.0")),
					}, nil
				},
			)
			Expect(err).NotTo(HaveOccurred())

			Expect(versionology.VersionFetcherToString(newVersions)).To(ConsistOf("1.17.4", "2.3.17", "2.4.0"))
		})

		context("when there are more new versions than allowed patches", func() {
			it("will return all new versions", func() {
				config, err := buildpack_config.ParseBuildpackToml(filepath.Join("testdata", "no-deps", "buildpack.toml"))
				Expect(err).NotTo(HaveOccurred())

				newVersions, err := retrieve.GetNewVersionsForId(
					"dep",
					config,
					func() (versionology.VersionFetcherArray, error) {
						return versionology.VersionFetcherArray{
							versionology.NewSimpleVersionFetcher(semver.MustParse("0.1.1")),
							versionology.NewSimpleVersionFetcher(semver.MustParse("1.0.0")),
							versionology.NewSimpleVersionFetcher(semver.MustParse("1.1.0")),
							versionology.NewSimpleVersionFetcher(semver.MustParse("1.1.1")),
							versionology.NewSimpleVersionFetcher(semver.MustParse("1.1.2")),
							versionology.NewSimpleVersionFetcher(semver.MustParse("1.1.3")),
							versionology.NewSimpleVersionFetcher(semver.MustParse("2.0.0")),
						}, nil
					},
				)
				Expect(err).NotTo(HaveOccurred())
				Expect(versionology.VersionFetcherToString(newVersions)).To(ConsistOf("1.1.2", "1.1.3"))
			})
		})

		context("when there are no constraints", func() {
			it("will return only versions not found in buildpack.toml", func() {
				config, err := buildpack_config.ParseBuildpackToml(filepath.Join("testdata", "no-constraints", "buildpack.toml"))
				Expect(err).NotTo(HaveOccurred())

				newVersions, err := retrieve.GetNewVersionsForId(
					"dep1",
					config,
					func() (versionology.VersionFetcherArray, error) {
						return versionology.VersionFetcherArray{
							versionology.NewSimpleVersionFetcher(semver.MustParse("1.0.0")),
							versionology.NewSimpleVersionFetcher(semver.MustParse("2.0.0")),
							versionology.NewSimpleVersionFetcher(semver.MustParse("3.0.0")),
						}, nil
					},
				)
				Expect(err).NotTo(HaveOccurred())
				Expect(versionology.VersionFetcherToString(newVersions)).To(ConsistOf("3.0.0"))
			})
		})

		context("when there are no existing constraints or dependencies for that id or stack", func() {
			it("will return all new versions", func() {
				config, err := buildpack_config.ParseBuildpackToml(filepath.Join("testdata", "empty", "buildpack.toml"))
				Expect(err).NotTo(HaveOccurred())

				newVersions, err := retrieve.GetNewVersionsForId(
					"id",
					config,
					func() (versionology.VersionFetcherArray, error) {
						return versionology.VersionFetcherArray{
							versionology.NewSimpleVersionFetcher(semver.MustParse("1.0.0")),
							versionology.NewSimpleVersionFetcher(semver.MustParse("2.0.0")),
							versionology.NewSimpleVersionFetcher(semver.MustParse("3.0.0")),
						}, nil
					},
				)
				Expect(err).NotTo(HaveOccurred())
				Expect(versionology.VersionFetcherToString(newVersions)).To(ConsistOf("1.0.0", "2.0.0", "3.0.0"))
			})
		})

		context("failure cases", func() {
			context("when getNewVersions returns an error", func() {
				it("will return the error", func() {
					config, err := buildpack_config.ParseBuildpackToml(filepath.Join("testdata", "deps-only", "buildpack.toml"))
					Expect(err).NotTo(HaveOccurred())

					_, err = retrieve.GetNewVersionsForId(
						"id",
						config,
						func() (versionology.VersionFetcherArray, error) {
							return versionology.VersionFetcherArray{}, errors.New("hi")
						},
					)
					Expect(err).To(MatchError("hi"))
				})
			})
		})
	})

}
