package retrieval_test

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/joshuatcasey/libdependency/buildpack_config"
	"github.com/joshuatcasey/libdependency/retrieval"
	"github.com/joshuatcasey/libdependency/versionology"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
)

func testGetNewVersionsForId(t *testing.T, context spec.G, it spec.S) {
	Expect := NewWithT(t).Expect

	context("GetNewVersionsForId", func() {
		it("will get new versions for id and stack", func() {
			config, err := buildpack_config.ParseBuildpackToml(filepath.Join("testdata", "bundler", "buildpack.toml"))
			Expect(err).NotTo(HaveOccurred())

			newVersions, err := retrieval.GetNewVersionsForId(
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

				newVersions, err := retrieval.GetNewVersionsForId(
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

				newVersions, err := retrieval.GetNewVersionsForId(
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

				newVersions, err := retrieval.GetNewVersionsForId(
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

					_, err = retrieval.GetNewVersionsForId(
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
