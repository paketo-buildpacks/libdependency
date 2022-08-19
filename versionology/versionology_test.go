package versionology_test

import (
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/joshuatcasey/libdependency/versionology"
	"github.com/paketo-buildpacks/packit/v2/cargo"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
)

func testVersionology(t *testing.T, context spec.G, it spec.S) {
	Expect := NewWithT(t).Expect

	context("SemverToString", func() {
		it("translates semver to strings", func() {
			stringVersions := versionology.SemverToString([]*semver.Version{
				semver.MustParse("v1.2.3"),
				semver.MustParse("2.3.4"),
				semver.MustParse("v9.8.7"),
			})

			Expect(stringVersions).To(HaveLen(3))
			Expect(stringVersions).To(ContainElements("1.2.3", "2.3.4", "9.8.7"))
		})
	})

	context("ConstraintsToString", func() {
		it("translates constraints to strings", func() {
			c1, _ := versionology.NewConstraint(cargo.ConfigMetadataDependencyConstraint{
				Constraint: ">=1.2.3",
			})
			c2, _ := versionology.NewConstraint(cargo.ConfigMetadataDependencyConstraint{
				Constraint: "4.*.*",
			})
			c3, _ := versionology.NewConstraint(cargo.ConfigMetadataDependencyConstraint{
				Constraint: "5.6.*",
			})

			stringVersions := versionology.ConstraintsToString([]versionology.Constraint{
				c1,
				c2,
				c3,
			})

			Expect(stringVersions).To(HaveLen(3))
			Expect(stringVersions).To(ContainElements(">=1.2.3", "4.*.*", "5.6.*"))
		})
	})

	context("FilterVersionsByConstraints", func() {
		mustParseConstraint := func(c string) *semver.Constraints {
			constraint, err := semver.NewConstraint(c)
			Expect(err).NotTo(HaveOccurred())
			return constraint
		}

		it("will filter versions by constraints", func() {
			results := versionology.FilterVersionsByConstraints([]*semver.Version{
				semver.MustParse("0.1.0"),
				semver.MustParse("1.1.0"),
				semver.MustParse("1.2.2"),
				semver.MustParse("1.2.3"),
				semver.MustParse("2.5.9"),
				semver.MustParse("3.2.0"),
				semver.MustParse("3.4.3"),
				semver.MustParse("3.4.5"),
				semver.MustParse("3.4.10"),
				semver.MustParse("3.7.0"),
				semver.MustParse("36.36.36"),
				semver.MustParse("99.99.99"),
			}, []*semver.Constraints{
				mustParseConstraint("1.2.*"),
				mustParseConstraint("3.4.*"),
				mustParseConstraint("36.*.*"),
			})

			Expect(results).To(HaveLen(6))
			Expect(versionology.SemverToString(results)).To(ContainElements("1.2.2", "1.2.3", "3.4.3", "3.4.5", "3.4.10", "36.36.36"))
		})
	})

	context("FilterUpstreamVersionsByConstraints", func() {
		it("will return only those upstream versions that match constraints and are newer than existing versions", func() {
			upstreamVersions := []versionology.HasVersion{
				versionology.NewSimpleHasVersion(semver.MustParse("6.0.0")),
				versionology.NewSimpleHasVersion(semver.MustParse("6.0.1")),
				versionology.NewSimpleHasVersion(semver.MustParse("6.0.2")),
				versionology.NewSimpleHasVersion(semver.MustParse("6.0.3")),
				versionology.NewSimpleHasVersion(semver.MustParse("6.0.4")),
				versionology.NewSimpleHasVersion(semver.MustParse("6.0.5")),
				versionology.NewSimpleHasVersion(semver.MustParse("6.0.6")),
				versionology.NewSimpleHasVersion(semver.MustParse("6.1.0")),
				versionology.NewSimpleHasVersion(semver.MustParse("6.1.1")),
				versionology.NewSimpleHasVersion(semver.MustParse("6.1.2")),
				versionology.NewSimpleHasVersion(semver.MustParse("6.1.3")),
				versionology.NewSimpleHasVersion(semver.MustParse("6.1.4")),
				versionology.NewSimpleHasVersion(semver.MustParse("6.1.5")),
				versionology.NewSimpleHasVersion(semver.MustParse("6.1.6")),
				versionology.NewSimpleHasVersion(semver.MustParse("7.0.0")),
				versionology.NewSimpleHasVersion(semver.MustParse("7.0.1")),
				versionology.NewSimpleHasVersion(semver.MustParse("7.0.2")),
				versionology.NewSimpleHasVersion(semver.MustParse("7.0.3")),
				versionology.NewSimpleHasVersion(semver.MustParse("7.0.4")),
				versionology.NewSimpleHasVersion(semver.MustParse("7.0.5")),
				versionology.NewSimpleHasVersion(semver.MustParse("7.0.6")),
			}

			c61, err := semver.NewConstraint("6.1.*")
			Expect(err).NotTo(HaveOccurred())
			c7, err := semver.NewConstraint("7.*.*")
			Expect(err).NotTo(HaveOccurred())
			constraints := []versionology.Constraint{
				{
					Constraint: c61,
					Patches:    4,
				},
				{
					Constraint: c7,
					Patches:    2,
				},
			}
			dependencies := []versionology.HasVersion{
				versionology.NewSimpleHasVersion(semver.MustParse("6.1.0")),
				versionology.NewSimpleHasVersion(semver.MustParse("7.0.4")),
			}

			filteredVersions := versionology.FilterUpstreamVersionsByConstraints("dep", upstreamVersions, constraints, dependencies)

			Expect(filteredVersions.GetVersionStrings()).To(ConsistOf("6.1.3", "6.1.4", "6.1.5", "6.1.6", "7.0.5", "7.0.6"))
		})
	})
}
