package retrieve_test

import (
	"path/filepath"
	"testing"

	"github.com/Masterminds/semver/v3"
	. "github.com/onsi/gomega"
	"github.com/paketo-buildpacks/libdependency/retrieve"
	"github.com/paketo-buildpacks/libdependency/versionology"
	"github.com/paketo-buildpacks/occam/matchers"
	"github.com/paketo-buildpacks/packit/v2/cargo"
	"github.com/sclevine/spec"
)

func testNewMetadataWithPlatforms(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect
		output string

		getAllVersions               retrieve.GetAllVersionsFunc
		generateMetadataWithPlatform retrieve.GenerateMetadataWithPlatformFunc
		transformPlatforms           retrieve.TransformsPlatformsFunc
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

			generateMetadataWithPlatform = func(versionFetcher versionology.VersionFetcher, platform retrieve.Platform) ([]versionology.Dependency, error) {
				dependency := cargo.ConfigMetadataDependency{
					ID:      "fake-dependency-id",
					Stacks:  []string{"jammy-stack", "bionic-stack"},
					Version: versionFetcher.Version().String(),
				}

				return versionology.NewDependencyArray(dependency, "linux-64")
			}

			transformPlatforms = func(platforms []retrieve.Platform) []retrieve.Platform {
				var transformed []retrieve.Platform

				if len(platforms) == 0 {
					platforms = append(platforms, retrieve.Platform{
						OS:   "linux",
						Arch: "amd64",
					})
				}

				for i := range platforms {
					if platforms[i].Arch == "amd64" {
						platforms[i].Arch = "x64"
					}
					transformed = append(transformed, platforms[i])
				}
				return transformed
			}
		})

		it("should generate metadata.json in the output dir", func() {
			retrieve.NewMetadataWithPlatforms("fake-dependency-id", getAllVersions, generateMetadataWithPlatform, transformPlatforms)

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

			generateMetadataWithPlatform = func(versionFetcher versionology.VersionFetcher, platform retrieve.Platform) ([]versionology.Dependency, error) {
				dependency := cargo.ConfigMetadataDependency{
					ID:      "python",
					Stacks:  []string{"io.buildpacks.stacks.jammy"},
					Version: versionFetcher.Version().String(),
				}

				return versionology.NewDependencyArray(dependency, "jammy")
			}

			transformPlatforms = func(platforms []retrieve.Platform) []retrieve.Platform {
				var transformed []retrieve.Platform

				if len(platforms) == 0 {
					platforms = append(platforms, retrieve.Platform{
						OS:   "linux",
						Arch: "amd64",
					})
				}

				for i := range platforms {
					if platforms[i].Arch == "amd64" {
						platforms[i].Arch = "x64"
					}
					transformed = append(transformed, platforms[i])
				}
				return transformed
			}
		})

		it("should only generate metadata for the two newest versions", func() {
			retrieve.NewMetadataWithPlatforms("python", getAllVersions, generateMetadataWithPlatform, transformPlatforms)

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

			generateMetadataWithPlatform = func(versionFetcher versionology.VersionFetcher, platform retrieve.Platform) ([]versionology.Dependency, error) {
				dependency := cargo.ConfigMetadataDependency{
					ID:      "not a real dependency id",
					Version: versionFetcher.Version().String(),
				}

				return versionology.NewDependencyArray(dependency, "jammy")
			}

			transformPlatforms = func(platforms []retrieve.Platform) []retrieve.Platform {
				var transformed []retrieve.Platform

				if len(platforms) == 0 {
					platforms = append(platforms, retrieve.Platform{
						OS:   "linux",
						Arch: "amd64",
					})
				}

				for i := range platforms {
					if platforms[i].Arch == "amd64" {
						platforms[i].Arch = "x64"
					}
					transformed = append(transformed, platforms[i])
				}
				return transformed
			}
		})

		it("all versions are found in the metadata.json", func() {
			retrieve.NewMetadataWithPlatforms("not a real dependency id", getAllVersions, generateMetadataWithPlatform, transformPlatforms)

			Expect(output).To(matchers.BeAFileMatching(MatchJSON(`
	[
		{"id":"not a real dependency id","version":"999.888.777","target":"jammy"},
		{"id":"not a real dependency id","version":"666.555.444","target":"jammy"},
		{"id":"not a real dependency id","version":"333.222.111","target":"jammy"}
	]`)))
		})
	})
}
