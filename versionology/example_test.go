package versionology_test

import (
	"github.com/paketo-buildpacks/libdependency/versionology"
)

func ExampleLogAllVersions() {
	array, err := versionology.NewSimpleVersionFetcherArray(
		"888.777.666", "1.0.0", "1.0.1", "1.0.2", "1.0.3",
		"1.0.4", "1.0.5", "1.0.6", "1.0.7", "1.0.8",
		"1.0.9", "1.10.0", "2.0.0", "2.0.1", "2.0.2",
		"2.0.3", "2.0.4", "2.0.5", "2.0.6", "2.0.7",
		"2.0.8", "2.0.9", "2.10.0")
	if err != nil {
		panic(err)
	}
	versionology.LogAllVersions("dep-id", "description", array)

	// Output: Found 23 versions of dep-id description
	//[
	//   "888.777.666", "2.10.0", "2.0.9",  "2.0.8", "2.0.7",
	//   "2.0.6",       "2.0.5",  "2.0.4",  "2.0.3", "2.0.2",
	//   "2.0.1",       "2.0.0",  "1.10.0", "1.0.9", "1.0.8",
	//   "1.0.7",       "1.0.6",  "1.0.5",  "1.0.4", "1.0.3",
	//   "1.0.2",       "1.0.1",  "1.0.0"
	//]
}

func ExampleLogAllVersions_oneVersion() {
	array, err := versionology.NewSimpleVersionFetcherArray("888.777.666")
	if err != nil {
		panic(err)
	}
	versionology.LogAllVersions("dep-id", "description", array)

	// Output: Found 1 version of dep-id description
	//[
	//   "888.777.666"
	//]
}
