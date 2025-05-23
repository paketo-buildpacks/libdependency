package retrieve

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/paketo-buildpacks/libdependency/buildpack_config"
	"github.com/paketo-buildpacks/libdependency/versionology"
	"github.com/paketo-buildpacks/packit/v2/fs"
)

// GetAllVersionsFunc is a function type that buildpack authors will implement and pass in to NewMetadata.
// The implementation should return all known upstream versions of a dependency.
// Buildpack authors can choose the source of these versions. Some examples include:
//
// - `nginx` versions from https://github.com/nginx/nginx/tags
// - `bundler` versions from https://rubygems.org/api/v1/versions/bundler.json
type GetAllVersionsFunc func() (versionology.VersionFetcherArray, error)

// GenerateMetadataFunc is a function type that buildpack authors will implement and pass in to NewMetadata.
// Given a versionology.VersionFetcher, the implementation must return the associated metadata for that version.
// If there are multiple targets for the same version, return multiple versionology.Dependency.
type GenerateMetadataFunc func(version versionology.VersionFetcher) ([]versionology.Dependency, error)

// NewMetadata is the entrypoint for a buildpack to retrieve new versions and the metadata thereof.
// Given a way to retrieve all versions (getNewVersions) and a way to generate metadata for a version (generateMetadata),
// this function will take in the dependency workflow inputs and the dependency workflow outputs
func NewMetadata(id string, getAllVersions GetAllVersionsFunc, generateMetadata GenerateMetadataFunc) {
	buildpackTomlPath, output := FetchArgs()
	validate(buildpackTomlPath, output)

	config, err := buildpack_config.ParseBuildpackToml(buildpackTomlPath)
	if err != nil {
		panic(err)
	}

	newVersions, err := GetNewVersionsForId(id, config, getAllVersions)
	if err != nil {
		panic(err)
	}

	dependencies := GenerateAllMetadata(newVersions, generateMetadata)

	metadataJson, err := toWorkflowJson(dependencies)
	if err != nil {
		panic(fmt.Errorf("unable to marshall metadata json, with error=%w", err))
	}

	if err = os.WriteFile(output, []byte(metadataJson), os.ModePerm); err != nil {
		panic(fmt.Errorf("cannot write to %s: %w", output, err))
	} else {
		fmt.Printf("Wrote metadata to %s\n", output)
	}
}

// toWorkflowJson will return a string containing JSON formatted as a GitHub workflow expects, with
// no whitespace outside of strings.
//
// Use this when printing output or writing a file intended for use by a workflow.
// https://github.com/orgs/community/discussions/26288
func toWorkflowJson(item any) (string, error) {
	if bytes, err := json.Marshal(item); err != nil {
		return "", err
	} else {
		return string(bytes), nil
	}
}

// GenerateAllMetadata is public for testing purposes only
func GenerateAllMetadata(newVersions versionology.VersionFetcherArray, generateMetadata GenerateMetadataFunc) []versionology.Dependency {
	var dependencies []versionology.Dependency
	for _, version := range newVersions {
		metadata, err := generateMetadata(version)
		if err != nil {
			panic(err)
		}

		var targets []string
		for _, metadatum := range metadata {
			targets = append(targets, metadatum.Target)
		}
		fmt.Printf("Generating metadata for %s, with targets [%s]\n", version.Version().String(), strings.Join(targets, ", "))
		dependencies = append(dependencies, metadata...)
	}
	return dependencies
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
