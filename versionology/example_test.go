package versionology_test

import (
	"github.com/Masterminds/semver/v3"
	"github.com/joshuatcasey/libdependency/versionology"
)

func ExampleLogAllVersions() {
	versionology.LogAllVersions("dep-id", "description", []versionology.VersionFetcher{
		versionology.NewSimpleVersionFetcher(semver.MustParse("888.777.666")),
		versionology.NewSimpleVersionFetcher(semver.MustParse("1.0.0")),
		versionology.NewSimpleVersionFetcher(semver.MustParse("1.0.1")),
		versionology.NewSimpleVersionFetcher(semver.MustParse("1.0.2")),
		versionology.NewSimpleVersionFetcher(semver.MustParse("1.0.3")),
		versionology.NewSimpleVersionFetcher(semver.MustParse("1.0.4")),
		versionology.NewSimpleVersionFetcher(semver.MustParse("1.0.5")),
		versionology.NewSimpleVersionFetcher(semver.MustParse("1.0.6")),
		versionology.NewSimpleVersionFetcher(semver.MustParse("1.0.7")),
		versionology.NewSimpleVersionFetcher(semver.MustParse("1.0.8")),
		versionology.NewSimpleVersionFetcher(semver.MustParse("1.0.9")),
		versionology.NewSimpleVersionFetcher(semver.MustParse("1.10.0")),
		versionology.NewSimpleVersionFetcher(semver.MustParse("2.0.0")),
		versionology.NewSimpleVersionFetcher(semver.MustParse("2.0.1")),
		versionology.NewSimpleVersionFetcher(semver.MustParse("2.0.2")),
		versionology.NewSimpleVersionFetcher(semver.MustParse("2.0.3")),
		versionology.NewSimpleVersionFetcher(semver.MustParse("2.0.4")),
		versionology.NewSimpleVersionFetcher(semver.MustParse("2.0.5")),
		versionology.NewSimpleVersionFetcher(semver.MustParse("2.0.6")),
		versionology.NewSimpleVersionFetcher(semver.MustParse("2.0.7")),
		versionology.NewSimpleVersionFetcher(semver.MustParse("2.0.8")),
		versionology.NewSimpleVersionFetcher(semver.MustParse("2.0.9")),
		versionology.NewSimpleVersionFetcher(semver.MustParse("2.10.0")),
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
