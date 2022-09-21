package retrieval

import (
	"github.com/joshuatcasey/libdependency/buildpack_config"
	"github.com/joshuatcasey/libdependency/versionology"
	"github.com/paketo-buildpacks/packit/v2/cargo"
)

// GetNewVersionsForId will return only those versions with the following properties:
// - returned by getAllVersions
// - match constraints
// - newer than all existing dependencies
func GetNewVersionsForId(id string, config cargo.Config, getAllVersions GetAllVersionsFunc) (versionology.VersionFetcherArray, error) {
	empty := versionology.NewVersionFetcherArray()

	allVersions, err := getAllVersions()
	if err != nil {
		return empty, err
	}

	versionology.LogAllVersions(id, "from upstream", allVersions)

	dependencies, err := buildpack_config.GetDependenciesById(id, config)
	if err != nil { //untested
		return empty, err
	}

	versionFetchers := make(versionology.VersionFetcherArray, len(dependencies))
	for i := range dependencies {
		versionFetchers[i] = dependencies[i]
	}

	constraints, err := buildpack_config.GetConstraintsById(id, config)
	if err != nil { //untested
		return empty, err
	}

	return versionology.FilterUpstreamVersionsByConstraints(id, allVersions, constraints, versionFetchers), nil
}
