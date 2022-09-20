package retrieve

import (
	"flag"
	"fmt"
	"os"

	"github.com/joshuatcasey/libdependency/buildpack_config"
	"github.com/joshuatcasey/libdependency/versionology"
	"github.com/joshuatcasey/libdependency/workflows"
	"github.com/paketo-buildpacks/packit/v2/fs"
)

// GenerateMetadataFunc is a function type that buildpack authors will implement and pass in to NewMetadata.
// Given a versionology.VersionFetcher, the implementation must return the associated metadata for that version.
// If there are multiple targets for the same version, return multiple versionology.Dependency.
type GenerateMetadataFunc func(version versionology.VersionFetcher) ([]versionology.Dependency, error)

// NewMetadata is the entrypoint for a buildpack to retrieve new versions and the metadata thereof.
// Given a way to retrieve all versions (getNewVersions) and a way to generate metadata for a version (generateMetadata),
// this function will take in the dependency workflow inputs and the dependency workflow outputs
func NewMetadata(id string, getNewVersions buildpack_config.VersionFetcherFunc, generateMetadata GenerateMetadataFunc) {
	buildpackTomlPath, output := FetchArgs()
	validate(buildpackTomlPath, output)

	config, err := buildpack_config.ParseBuildpackToml(buildpackTomlPath)
	if err != nil {
		panic(err)
	}

	newVersions, err := buildpack_config.GetNewVersionsForId(id, config, getNewVersions)
	if err != nil {
		panic(err)
	}

	var dependencies []versionology.Dependency
	for _, version := range newVersions {
		fmt.Printf("Generating metadata for %s\n", version.Version().String())
		metadata, err := generateMetadata(version)
		if err != nil {
			panic(err)
		}
		dependencies = append(dependencies, metadata...)
	}

	metadataJson, err := workflows.ToWorkflowJson(dependencies)
	if err != nil {
		panic(fmt.Errorf("unable to marshall metadata json, with error=%w", err))
	}

	if err = os.WriteFile(output, []byte(metadataJson), os.ModePerm); err != nil {
		panic(fmt.Errorf("cannot write to %s: %w", output, err))
	} else {
		fmt.Printf("Wrote metadata to %s\n", output)
	}
}

func validate(buildpackTomlPath, metadataFile string) {
	if exists, err := fs.Exists(buildpackTomlPath); err != nil {
		panic(err)
	} else if !exists {
		panic(fmt.Errorf("could not locate buildpack.toml at '%s'", buildpackTomlPath))
	}

	if metadataFile == "" {
		panic("metadataFile is required")
	}
}

type FetchArgsFunc func() (string, string)

// FetchArgs is public for testing purposes
var FetchArgs = func() (buildpackTomlPath, output string) {
	buildpackTomlPathUsage := "full path to the buildpack.toml file, using only one of camelCase, snake_case, or dash_case"

	flag.StringVar(&buildpackTomlPath, "buildpackTomlPath", "", buildpackTomlPathUsage)
	flag.StringVar(&buildpackTomlPath, "buildpack_toml_path", buildpackTomlPath, buildpackTomlPathUsage)
	flag.StringVar(&buildpackTomlPath, "buildpack-toml-path", buildpackTomlPath, buildpackTomlPathUsage)
	flag.StringVar(&output, "output", "", "filename for the output JSON metadata")
	flag.Parse()
	return
}
