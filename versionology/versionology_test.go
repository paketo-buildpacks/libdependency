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

	context("FilterUpstreamVersionsByConstraints", func() {
		it("will return only those upstream versions that match constraints and are newer than existing versions", func() {
			upstreamVersions := []versionology.VersionFetcher{
				versionology.NewSimpleVersionFetcher(semver.MustParse("6.0.0")),
				versionology.NewSimpleVersionFetcher(semver.MustParse("6.0.1")),
				versionology.NewSimpleVersionFetcher(semver.MustParse("6.0.2")),
				versionology.NewSimpleVersionFetcher(semver.MustParse("6.0.3")),
				versionology.NewSimpleVersionFetcher(semver.MustParse("6.0.4")),
				versionology.NewSimpleVersionFetcher(semver.MustParse("6.0.5")),
				versionology.NewSimpleVersionFetcher(semver.MustParse("6.0.6")),
				versionology.NewSimpleVersionFetcher(semver.MustParse("6.1.0")),
				versionology.NewSimpleVersionFetcher(semver.MustParse("6.1.1")),
				versionology.NewSimpleVersionFetcher(semver.MustParse("6.1.2")),
				versionology.NewSimpleVersionFetcher(semver.MustParse("6.1.3")),
				versionology.NewSimpleVersionFetcher(semver.MustParse("6.1.4")),
				versionology.NewSimpleVersionFetcher(semver.MustParse("6.1.5")),
				versionology.NewSimpleVersionFetcher(semver.MustParse("6.1.6")),
				versionology.NewSimpleVersionFetcher(semver.MustParse("7.0.0")),
				versionology.NewSimpleVersionFetcher(semver.MustParse("7.0.1")),
				versionology.NewSimpleVersionFetcher(semver.MustParse("7.0.2")),
				versionology.NewSimpleVersionFetcher(semver.MustParse("7.0.3")),
				versionology.NewSimpleVersionFetcher(semver.MustParse("7.0.4")),
				versionology.NewSimpleVersionFetcher(semver.MustParse("7.0.5")),
				versionology.NewSimpleVersionFetcher(semver.MustParse("7.0.6")),
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
			dependencies := []versionology.VersionFetcher{
				versionology.NewSimpleVersionFetcher(semver.MustParse("6.1.0")),
				versionology.NewSimpleVersionFetcher(semver.MustParse("7.0.4")),
			}

			filteredVersions := versionology.FilterUpstreamVersionsByConstraints("dep", upstreamVersions, constraints, dependencies)

			Expect(filteredVersions.GetVersionStrings()).To(ConsistOf("6.1.3", "6.1.4", "6.1.5", "6.1.6", "7.0.5", "7.0.6"))
		})
	})
}
