package versionology_test

import (
	"github.com/Masterminds/semver/v3"
	"github.com/joshuatcasey/libdependency/versionology"
)

//nolint:govet
func ExampleUnitLogAllVersions() {
	versionology.LogAllVersions("dep-id", "description", []versionology.HasVersion{
		versionology.NewSimpleHasVersion(semver.MustParse("888.777.666")),
		versionology.NewSimpleHasVersion(semver.MustParse("1.0.0")),
		versionology.NewSimpleHasVersion(semver.MustParse("1.0.1")),
		versionology.NewSimpleHasVersion(semver.MustParse("1.0.2")),
		versionology.NewSimpleHasVersion(semver.MustParse("1.0.3")),
		versionology.NewSimpleHasVersion(semver.MustParse("1.0.4")),
		versionology.NewSimpleHasVersion(semver.MustParse("1.0.5")),
		versionology.NewSimpleHasVersion(semver.MustParse("1.0.6")),
		versionology.NewSimpleHasVersion(semver.MustParse("1.0.7")),
		versionology.NewSimpleHasVersion(semver.MustParse("1.0.8")),
		versionology.NewSimpleHasVersion(semver.MustParse("1.0.9")),
		versionology.NewSimpleHasVersion(semver.MustParse("1.10.0")),
		versionology.NewSimpleHasVersion(semver.MustParse("2.0.0")),
		versionology.NewSimpleHasVersion(semver.MustParse("2.0.1")),
		versionology.NewSimpleHasVersion(semver.MustParse("2.0.2")),
		versionology.NewSimpleHasVersion(semver.MustParse("2.0.3")),
		versionology.NewSimpleHasVersion(semver.MustParse("2.0.4")),
		versionology.NewSimpleHasVersion(semver.MustParse("2.0.5")),
		versionology.NewSimpleHasVersion(semver.MustParse("2.0.6")),
		versionology.NewSimpleHasVersion(semver.MustParse("2.0.7")),
		versionology.NewSimpleHasVersion(semver.MustParse("2.0.8")),
		versionology.NewSimpleHasVersion(semver.MustParse("2.0.9")),
		versionology.NewSimpleHasVersion(semver.MustParse("2.10.0")),
	})

	// Output: Found 23 versions of dep-id description
	//[
	//   "888.777.666", "2.10.0", "2.0.9",  "2.0.8", "2.0.7",
	//   "2.0.6",       "2.0.5",  "2.0.4",  "2.0.3", "2.0.2",
	//   "2.0.1",       "2.0.0",  "1.10.0", "1.0.9", "1.0.8",
	//   "1.0.7",       "1.0.6",  "1.0.5",  "1.0.4", "1.0.3",
	//   "1.0.2",       "1.0.1",  "1.0.0"
	//]

}
