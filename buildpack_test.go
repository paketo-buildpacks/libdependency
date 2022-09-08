package libdependency_test

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/paketo-buildpacks/libdependency"
	"github.com/paketo-buildpacks/libdependency/versionology"
	"github.com/paketo-buildpacks/occam"
	"github.com/paketo-buildpacks/packit/v2/cargo"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
)

func testBuildpackToml(t *testing.T, context spec.G, it spec.S) {
	Expect := NewWithT(t).Expect

	context("ParseBuildpackToml", func() {
		it("parses bundler's buildpack.toml", func() {
			config, err := libdependency.ParseBuildpackToml(filepath.Join("testdata", "bundler", "buildpack.toml"))
			Expect(err).NotTo(HaveOccurred())
			Expect(config.Metadata.Dependencies).To(HaveLen(3))
		})

		context("failure cases", func() {
			it("returns an error when path not found", func() {
				_, err := libdependency.ParseBuildpackToml("/bad/path")
				Expect(err).To(MatchError(os.ErrNotExist))
				Expect(err).To(MatchError("unable to open buildpack.toml: open /bad/path: no such file or directory"))
			})

			it("returns an error when buildpack cannot be parsed", func() {
				_, err := libdependency.ParseBuildpackToml(filepath.Join("testdata", "invalid", "buildpack.toml"))
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
			dependencies, err := libdependency.GetDependenciesById("id2", config)
			Expect(err).NotTo(HaveOccurred())

			Expect(dependencies.Versions()).To(ConsistOf("2.2.2", "3.3.3", "4.4.4"))
		})

		context("failure cases", func() {
			it.Before(func() {
				config.Metadata.Dependencies = append(config.Metadata.Dependencies, cargo.ConfigMetadataDependency{
					ID:      "id-invalid",
					Version: "not valid",
				})
			})

			it("will return error when version is not valid semver", func() {
				_, err := libdependency.GetDependenciesById("id-invalid", config)
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
			dependencies, err := libdependency.GetDependenciesByIdAndStack("id2", "stack1", config)
			Expect(err).NotTo(HaveOccurred())

			Expect(dependencies.Versions()).To(ConsistOf("2.2.2", "4.4.4"))
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
				_, err := libdependency.GetDependenciesByIdAndStack("id-invalid", "stack", config)
				Expect(err).To(HaveOccurred())
			})
		})
	})

	context("GetConstraintsById", func() {
		it("will filter by id", func() {
			constraints, err := libdependency.GetConstraintsById("id2", cargo.Config{
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
				_, err := libdependency.GetConstraintsById("id1", cargo.Config{
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

	context("GetNewVersionsForId", func() {
		it("will get new versions for id and stack", func() {
			config, err := libdependency.ParseBuildpackToml(filepath.Join("testdata", "bundler", "buildpack.toml"))
			Expect(err).NotTo(HaveOccurred())

			newVersions, err := libdependency.GetNewVersionsForId(
				"bundler",
				config,
				func() ([]versionology.HasVersion, error) {
					return []versionology.HasVersion{
						versionology.NewSimpleHasVersion(semver.MustParse("0.1.1")),
						versionology.NewSimpleHasVersion(semver.MustParse("1.17.3")),
						versionology.NewSimpleHasVersion(semver.MustParse("1.17.4")),
						versionology.NewSimpleHasVersion(semver.MustParse("2.3.16")),
						versionology.NewSimpleHasVersion(semver.MustParse("2.3.17")),
						versionology.NewSimpleHasVersion(semver.MustParse("2.4.0")),
						versionology.NewSimpleHasVersion(semver.MustParse("3.0.0")),
					}, nil
				},
			)
			Expect(err).NotTo(HaveOccurred())

			Expect(versionology.HasVersionToString(newVersions)).To(ConsistOf("1.17.4", "2.3.17", "2.4.0"))
		})

		context("when there are more new versions than allowed patches", func() {
			it("will return all new versions", func() {
				config, err := libdependency.ParseBuildpackToml(filepath.Join("testdata", "no-deps", "buildpack.toml"))
				Expect(err).NotTo(HaveOccurred())

				newVersions, err := libdependency.GetNewVersionsForId(
					"dep",
					config,
					func() ([]versionology.HasVersion, error) {
						return []versionology.HasVersion{
							versionology.NewSimpleHasVersion(semver.MustParse("0.1.1")),
							versionology.NewSimpleHasVersion(semver.MustParse("1.0.0")),
							versionology.NewSimpleHasVersion(semver.MustParse("1.1.0")),
							versionology.NewSimpleHasVersion(semver.MustParse("1.1.1")),
							versionology.NewSimpleHasVersion(semver.MustParse("1.1.2")),
							versionology.NewSimpleHasVersion(semver.MustParse("1.1.3")),
							versionology.NewSimpleHasVersion(semver.MustParse("2.0.0")),
						}, nil
					},
				)
				Expect(err).NotTo(HaveOccurred())
				Expect(versionology.HasVersionToString(newVersions)).To(ConsistOf("1.1.2", "1.1.3"))
			})
		})

		context("when there are no constraints", func() {
			it("will return only versions not found in buildpack.toml", func() {
				config, err := libdependency.ParseBuildpackToml(filepath.Join("testdata", "no-constraints", "buildpack.toml"))
				Expect(err).NotTo(HaveOccurred())

				newVersions, err := libdependency.GetNewVersionsForId(
					"dep1",
					config,
					func() ([]versionology.HasVersion, error) {
						return []versionology.HasVersion{
							versionology.NewSimpleHasVersion(semver.MustParse("1.0.0")),
							versionology.NewSimpleHasVersion(semver.MustParse("2.0.0")),
							versionology.NewSimpleHasVersion(semver.MustParse("3.0.0")),
						}, nil
					},
				)
				Expect(err).NotTo(HaveOccurred())
				Expect(versionology.HasVersionToString(newVersions)).To(ConsistOf("3.0.0"))
			})
		})

		context("when there are no existing constraints or dependencies for that id or stack", func() {
			it("will return all new versions", func() {
				config, err := libdependency.ParseBuildpackToml(filepath.Join("testdata", "empty", "buildpack.toml"))
				Expect(err).NotTo(HaveOccurred())

				newVersions, err := libdependency.GetNewVersionsForId(
					"id",
					config,
					func() ([]versionology.HasVersion, error) {
						return []versionology.HasVersion{
							versionology.NewSimpleHasVersion(semver.MustParse("1.0.0")),
							versionology.NewSimpleHasVersion(semver.MustParse("2.0.0")),
							versionology.NewSimpleHasVersion(semver.MustParse("3.0.0")),
						}, nil
					},
				)
				Expect(err).NotTo(HaveOccurred())
				Expect(versionology.HasVersionToString(newVersions)).To(ConsistOf("1.0.0", "2.0.0", "3.0.0"))
			})
		})

		context("failure cases", func() {
			context("when getNewVersions returns an error", func() {
				it("will return the error", func() {
					config, err := libdependency.ParseBuildpackToml(filepath.Join("testdata", "deps-only", "buildpack.toml"))
					Expect(err).NotTo(HaveOccurred())

					_, err = libdependency.GetNewVersionsForId(
						"id",
						config,
						func() ([]versionology.HasVersion, error) {
							return []versionology.HasVersion{}, errors.New("hi")
						},
					)
					Expect(err).To(MatchError("hi"))
				})
			})
		})
	})

	context("PruneBuildpackToml", func() {
		it("will prune", func() {
			source, err := occam.Source(filepath.Join("testdata", "deps-only"))
			Expect(err).NotTo(HaveOccurred())

			err = libdependency.PruneBuildpackToml(filepath.Join(source, "buildpack.toml"))
			Expect(err).NotTo(HaveOccurred())

			bytes, err := os.ReadFile(filepath.Join(source, "buildpack.toml"))
			Expect(err).NotTo(HaveOccurred())

			Expect(string(bytes)).To(Equal(`[buildpack]

[metadata]

  [[metadata.dependencies]]
    id = "dep1"
    stacks = ["stack1"]
    version = "1.2.0"

  [[metadata.dependencies]]
    id = "dep1"
    stacks = ["stack2"]
    version = "1.2.0"

  [[metadata.dependency-constraints]]
    constraint = "1.*.*"
    id = "dep1"
    patches = 1
`))
		})

		context("failure cases", func() {
			context("when buildpack.toml does not exist", func() {
				it("will return the error", func() {
					err := libdependency.PruneBuildpackToml("/bad/path")
					Expect(err).To(MatchError(os.ErrNotExist))
					Expect(err).To(MatchError("unable to open buildpack.toml: open /bad/path: no such file or directory"))
				})
			})

			context("when buildpack.toml is not writeable", func() {
				it("will return the error", func() {
					source, err := occam.Source(filepath.Join("testdata", "deps-only"))
					Expect(err).NotTo(HaveOccurred())

					Expect(os.Chmod(filepath.Join(source, "buildpack.toml"), 0444)).To(Succeed())

					err = libdependency.PruneBuildpackToml(filepath.Join(source, "buildpack.toml"))
					Expect(err).To(MatchError(os.ErrPermission))
					Expect(err).To(MatchError(fmt.Sprintf("unable to open buildpack.toml for writing: open %s: permission denied", filepath.Join(source, "buildpack.toml"))))
				})
			})

			context("when configuration cannot be encoded", func() {
				var (
					tempDir string
					err     error
				)

				it.Before(func() {
					tempDir, err = os.MkdirTemp("", "toml")
					Expect(err).NotTo(HaveOccurred())

					Expect(os.WriteFile(filepath.Join(tempDir, "buildpack.toml"), []byte(`[metadata]
  [[metadata.dependency-constraints]]
    constraint = "1.*.*"
    id = "dep1"
    patches = 0
`), os.ModePerm)).To(Succeed())
				})

				it.After(func() {
					Expect(os.RemoveAll(tempDir)).To(Succeed())
				})

				it("will return the error", func() {
					err = libdependency.PruneBuildpackToml(filepath.Join(tempDir, "buildpack.toml"))
					Expect(err).To(HaveOccurred())
					Expect(err).To(MatchError(ContainSubstring("failure to assert type: unexpected data in constraint patches")))
				})
			})
		})
	})

	context("PruneConfig", func() {
		it("will remove dependencies that are no longer part of the constraints", func() {
			config := libdependency.PruneConfig(cargo.Config{
				Metadata: cargo.ConfigMetadata{
					Dependencies: []cargo.ConfigMetadataDependency{
						{ID: "dep1", Stacks: []string{"stack1"}, Version: "0.1.0"},
						{ID: "dep1", Stacks: []string{"stack1"}, Version: "1.0.0"},
						{ID: "dep1", Stacks: []string{"stack2"}, Version: "1.0.0"},
						{ID: "dep1", Stacks: []string{"stack1"}, Version: "1.0.1"},
						{ID: "dep1", Stacks: []string{"stack2"}, Version: "1.0.1"},
						{ID: "dep1", Stacks: []string{"stack1"}, Version: "1.0.2"},
						{ID: "dep1", Stacks: []string{"stack2"}, Version: "1.0.2"},
						{ID: "dep1", Stacks: []string{"stack1"}, Version: "1.0.0"},
					},
					DependencyConstraints: []cargo.ConfigMetadataDependencyConstraint{
						{ID: "dep1", Constraint: "1.*.*", Patches: 2},
					},
				},
			})

			Expect(config.Metadata.Dependencies).To(ConsistOf(
				And(HaveField("Version", "1.0.1"), HaveField("Stacks", []string{"stack1"})),
				And(HaveField("Version", "1.0.1"), HaveField("Stacks", []string{"stack2"})),
				And(HaveField("Version", "1.0.2"), HaveField("Stacks", []string{"stack1"})),
				And(HaveField("Version", "1.0.2"), HaveField("Stacks", []string{"stack2"})),
			))
		})

		it("will respect having multiple stacks", func() {
			config := libdependency.PruneConfig(cargo.Config{
				Metadata: cargo.ConfigMetadata{
					Dependencies: []cargo.ConfigMetadataDependency{
						{ID: "dep1", Stacks: []string{"stack1", "stack2"}, Version: "2.0.0"},
						{ID: "dep1", Stacks: []string{"stack1", "stack2"}, Version: "2.1.0"},
						{ID: "dep1", Stacks: []string{"stack1", "stack2"}, Version: "2.2.0"},
					},
					DependencyConstraints: []cargo.ConfigMetadataDependencyConstraint{
						{ID: "dep1", Constraint: "1.*.*", Patches: 2},
						{ID: "dep1", Constraint: "2.*.*", Patches: 2},
						{ID: "dep1", Constraint: "3.*.*", Patches: 2},
					},
				},
			})

			Expect(config.Metadata.Dependencies).To(ConsistOf(
				And(HaveField("Version", "2.1.0"), HaveField("Stacks", []string{"stack1", "stack2"})),
				And(HaveField("Version", "2.2.0"), HaveField("Stacks", []string{"stack1", "stack2"})),
			))
		})
	})
}
