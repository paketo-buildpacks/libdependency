package buildpack_config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/paketo-buildpacks/libdependency/buildpack_config"
	"github.com/paketo-buildpacks/libdependency/versionology"
	"github.com/paketo-buildpacks/packit/v2/cargo"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
)

func testBuildpackToml(t *testing.T, context spec.G, it spec.S) {
	Expect := NewWithT(t).Expect

	context("ParseBuildpackToml", func() {
		it("parses bundler's buildpack.toml", func() {
			config, err := buildpack_config.ParseBuildpackToml(filepath.Join("testdata", "bundler", "buildpack.toml"))
			Expect(err).NotTo(HaveOccurred())
			Expect(config.Metadata.Dependencies).To(HaveLen(3))
		})

		context("failure cases", func() {
			it("returns an error when path not found", func() {
				_, err := buildpack_config.ParseBuildpackToml("/bad/path")
				Expect(err).To(MatchError(os.ErrNotExist))
				Expect(err).To(MatchError("unable to open buildpack.toml: open /bad/path: no such file or directory"))
			})

			it("returns an error when buildpack cannot be parsed", func() {
				_, err := buildpack_config.ParseBuildpackToml(filepath.Join("testdata", "invalid", "buildpack.toml"))
				Expect(err.Error()).To(ContainSubstring("unable to parse buildpack.toml: toml: "))
			})
		})
	})

	context("GetDependenciesById", func() {
		var config cargo.Config

		it.Before(func() {
			config = cargo.Config{
				Metadata: cargo.ConfigMetadata{
					Dependencies: []cargo.ConfigMetadataDependency{
						{ID: "id1", Version: "1.1.1"},
						{ID: "id2", Version: "2.2.2"},
						{ID: "id2", Version: "3.3.3"},
						{ID: "id2", Version: "4.4.4"},
						{ID: "id3", Version: "5.5.5"},
					},
				},
			}
		})

		it("will filter by id", func() {
			dependencies, err := buildpack_config.GetDependenciesById("id2", config)
			Expect(err).NotTo(HaveOccurred())

			Expect(versionology.Versions(dependencies)).To(ConsistOf("2.2.2", "3.3.3", "4.4.4"))
		})

		context("failure cases", func() {
			it.Before(func() {
				config.Metadata.Dependencies = append(config.Metadata.Dependencies, cargo.ConfigMetadataDependency{
					ID:      "id-invalid",
					Version: "not valid",
				})
			})

			it("will return error when version is not valid semver", func() {
				_, err := buildpack_config.GetDependenciesById("id-invalid", config)
				Expect(err).To(HaveOccurred())
			})
		})
	})

	context("GetDependenciesByIdAndStack", func() {
		var config cargo.Config

		it.Before(func() {
			config = cargo.Config{
				Metadata: cargo.ConfigMetadata{
					Dependencies: []cargo.ConfigMetadataDependency{
						{ID: "id1", Stacks: []string{"stack1"}, Version: "1.1.1"},
						{ID: "id2", Stacks: []string{"stack1"}, Version: "2.2.2"},
						{ID: "id2", Stacks: []string{"stack2"}, Version: "3.3.3"},
						{ID: "id2", Stacks: []string{"stack1", "stack2"}, Version: "4.4.4"},
						{ID: "id3", Stacks: []string{"stack1"}, Version: "5.5.5"},
					},
				},
			}
		})

		it("will filter by id", func() {
			dependencies, err := buildpack_config.GetDependenciesByIdAndStack("id2", "stack1", config)
			Expect(err).NotTo(HaveOccurred())

			Expect(versionology.Versions(dependencies)).To(ConsistOf("2.2.2", "4.4.4"))
		})

		context("failure cases", func() {
			it.Before(func() {
				config.Metadata.Dependencies = append(config.Metadata.Dependencies, cargo.ConfigMetadataDependency{
					ID:      "id-invalid",
					Stacks:  []string{"stack"},
					Version: "not valid",
				})
			})

			it("will return error when version is not valid semver", func() {
				_, err := buildpack_config.GetDependenciesByIdAndStack("id-invalid", "stack", config)
				Expect(err).To(HaveOccurred())
			})
		})
	})

	context("GetConstraintsById", func() {
		it("will filter by id", func() {
			constraints, err := buildpack_config.GetConstraintsById("id2", cargo.Config{
				Metadata: cargo.ConfigMetadata{
					DependencyConstraints: []cargo.ConfigMetadataDependencyConstraint{
						{ID: "id1", Constraint: ">=1.2.3"},
						{ID: "id2", Constraint: "2.*"},
						{ID: "id2", Constraint: ">=3.4.5"},
						{ID: "id3", Constraint: "4.*.*"},
					},
				},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(versionology.ConstraintsToString(constraints)).To(ConsistOf("2.*", ">=3.4.5"))
		})

		context("failure cases", func() {
			it("will return error if constraint is not valid semver", func() {
				_, err := buildpack_config.GetConstraintsById("id1", cargo.Config{
					Metadata: cargo.ConfigMetadata{
						DependencyConstraints: []cargo.ConfigMetadataDependencyConstraint{
							{ID: "id1", Constraint: "foo"},
						},
					},
				})
				Expect(err).To(HaveOccurred())
			})
		})
	})
}
