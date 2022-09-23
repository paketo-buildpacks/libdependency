package retrieve_test

import (
	"github.com/joshuatcasey/libdependency/retrieve"
	"github.com/joshuatcasey/libdependency/versionology"
	"github.com/paketo-buildpacks/packit/v2/cargo"
)

func ExampleGenerateAllMetadata() {

	versions, _ := versionology.NewSimpleVersionFetcherArray("1.2.3", "4.5.6", "7.8.9")

	generateMetadata := func(version versionology.VersionFetcher) ([]versionology.Dependency, error) {
		dep := cargo.ConfigMetadataDependency{ID: "dep-id", Version: version.Version().String()}

		switch version.Version().String() {
		case "1.2.3":
			return versionology.NewDependencyArray(dep, "target1")
		case "4.5.6":
			dep1, _ := versionology.NewDependency(dep, "target1")
			dep2, _ := versionology.NewDependency(dep, "target2")
			return []versionology.Dependency{dep1, dep2}, nil
		case "7.8.9":
			dep1, _ := versionology.NewDependency(dep, "target1")
			dep2, _ := versionology.NewDependency(dep, "target2")
			dep3, _ := versionology.NewDependency(dep, "target3")
			return []versionology.Dependency{dep1, dep2, dep3}, nil
		default:
			panic("unknown version")
		}
	}

	retrieve.GenerateAllMetadata(versions, generateMetadata)

	// Output:
	// Generating metadata for 1.2.3, with targets [target1]
	// Generating metadata for 4.5.6, with targets [target1, target2]
	// Generating metadata for 7.8.9, with targets [target1, target2, target3]

}
