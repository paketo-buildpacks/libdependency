package buildpack_config

import (
	"errors"
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/joshuatcasey/collections"
	"github.com/paketo-buildpacks/libdependency/versionology"
	"github.com/paketo-buildpacks/packit/v2/cargo"
)

// ParseBuildpackToml takes in a path to a buildpack.toml and parses that into a cargo.Config
func ParseBuildpackToml(buildpackTomlPath string) (cargo.Config, error) {
	if config, err := cargo.NewBuildpackParser().Parse(buildpackTomlPath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return cargo.Config{}, fmt.Errorf("unable to open buildpack.toml: %w", err)
		} else if tomlError, ok := err.(toml.ParseError); ok {
			return cargo.Config{}, fmt.Errorf("unable to parse buildpack.toml: %s", tomlError.ErrorWithPosition())
		} else { //untested
			return cargo.Config{}, err
		}
	} else {
		return config, nil
	}
}

// GetDependenciesById will return an array of dependencies with the given id, of a type that "extends"
// cargo.ConfigMetadataDependency and implements versionology.VersionFetcher
func GetDependenciesById(id string, config cargo.Config) ([]versionology.Dependency, error) {
	dependencies := collections.FilterFunc(config.Metadata.Dependencies, func(d cargo.ConfigMetadataDependency) bool {
		return d.ID == id
	})

	return collections.TransformFuncWithError(dependencies, func(dependency cargo.ConfigMetadataDependency) (versionology.Dependency, error) {
		return versionology.NewDependency(dependency, "")
	})
}

// GetDependenciesByIdAndStack will return an array of dependencies with the given id and stack,
// of a type that "extends" cargo.ConfigMetadataDependency and implements versionology.VersionFetcher
func GetDependenciesByIdAndStack(id, stack string, config cargo.Config) ([]versionology.Dependency, error) {
	dependencies := collections.FilterFunc(config.Metadata.Dependencies, func(d cargo.ConfigMetadataDependency) bool {
		return d.ID == id && d.HasStack(stack)
	})

	return collections.TransformFuncWithError(dependencies, func(dependency cargo.ConfigMetadataDependency) (versionology.Dependency, error) {
		return versionology.NewDependency(dependency, "")
	})
}

// GetConstraintsById will return an array of constraints with the given id
func GetConstraintsById(id string, config cargo.Config) ([]versionology.Constraint, error) {
	constraints := collections.FilterFunc(config.Metadata.DependencyConstraints, func(c cargo.ConfigMetadataDependencyConstraint) bool {
		return c.ID == id
	})

	return collections.TransformFuncWithError(constraints, versionology.NewConstraint)
}
