package libdependency

import (
	"errors"
	"fmt"
	"os"
	"sort"

	"github.com/BurntSushi/toml"
	"github.com/Masterminds/semver/v3"
	"github.com/joshuatcasey/collections"
	"github.com/paketo-buildpacks/libdependency/versionology"
	"github.com/paketo-buildpacks/packit/v2/cargo"
	"golang.org/x/exp/slices"
)

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

func GetDependenciesById(id string, config cargo.Config) (versionology.DependencyArray, error) {
	dependencies := collections.FilterFunc(config.Metadata.Dependencies, func(d cargo.ConfigMetadataDependency) bool {
		return d.ID == id
	})

	return collections.TransformFuncWithError(dependencies, func(dependency cargo.ConfigMetadataDependency) (versionology.Dependency, error) {
		if semverVersion, err := semver.NewVersion(dependency.Version); err != nil {
			return versionology.Dependency{}, err
		} else {
			return versionology.Dependency{
				ConfigMetadataDependency: dependency,
				Version:                  semverVersion,
			}, nil
		}
	})
}

func GetDependenciesByIdAndStack(id, stack string, config cargo.Config) (versionology.DependencyArray, error) {
	dependencies := collections.FilterFunc(config.Metadata.Dependencies, func(d cargo.ConfigMetadataDependency) bool {
		return d.ID == id && d.HasStack(stack)
	})

	return collections.TransformFuncWithError(dependencies, func(dependency cargo.ConfigMetadataDependency) (versionology.Dependency, error) {
		if semverVersion, err := semver.NewVersion(dependency.Version); err != nil {
			return versionology.Dependency{}, err
		} else {
			return versionology.Dependency{
				ConfigMetadataDependency: dependency,
				Version:                  semverVersion,
			}, nil
		}
	})
}

func GetConstraintsById(id string, config cargo.Config) ([]versionology.Constraint, error) {
	constraints := collections.FilterFunc(config.Metadata.DependencyConstraints, func(c cargo.ConfigMetadataDependencyConstraint) bool {
		return c.ID == id
	})

	return collections.TransformFuncWithError(constraints, versionology.NewConstraint)
}

type AllVersionsFunc func() ([]*semver.Version, error)

type HasVersionsFunc func() ([]versionology.HasVersion, error)

func GetNewVersionsForId(id string, config cargo.Config, getAllVersions HasVersionsFunc) ([]versionology.HasVersion, error) {
	empty := make([]versionology.HasVersion, 0)

	allVersions, err := getAllVersions()
	if err != nil {
		return empty, err
	}

	versionology.LogAllVersions(id, "from upstream", allVersions)

	dependencies, err := GetDependenciesById(id, config)
	if err != nil { //untested
		return empty, err
	}

	hasVersionDependencies := make(versionology.HasVersionArray, len(dependencies))
	for i := range dependencies {
		hasVersionDependencies[i] = dependencies[i]
	}

	constraints, err := GetConstraintsById(id, config)
	if err != nil { //untested
		return empty, err
	}

	return versionology.FilterUpstreamVersionsByConstraints(id, allVersions, constraints, hasVersionDependencies), nil
}

func PruneBuildpackToml(buildpackTomlPath string) error {
	config, err := ParseBuildpackToml(buildpackTomlPath)
	if err != nil {
		return err
	}

	config = PruneConfig(config)

	file, err := os.OpenFile(buildpackTomlPath, os.O_RDWR|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("unable to open buildpack.toml for writing: %w", err)
	}
	defer file.Close()

	return cargo.EncodeConfig(file, config)
}

func PruneConfig(config cargo.Config) cargo.Config {

	// Get a map from constraints to dependencies
	constraintToDependencies := make(map[cargo.ConfigMetadataDependencyConstraint][]cargo.ConfigMetadataDependency)

	for _, dependency := range config.Metadata.Dependencies {
		dependencyVersionAsSemver := semver.MustParse(dependency.Version)
		for _, constraint := range config.Metadata.DependencyConstraints {
			if constraintAsSemver, err := semver.NewConstraint(constraint.Constraint); err != nil {
				panic(err)
			} else if dependency.ID == constraint.ID && constraintAsSemver.Check(dependencyVersionAsSemver) {
				constraintToDependencies[constraint] = append(constraintToDependencies[constraint], dependency)
			}
		}
	}

	constraintToPatches := make(map[cargo.ConfigMetadataDependencyConstraint][]string)

	// We can have more than one dependency with the same version,
	// so we have to figure out which versions are captured in the patches
	for constraint, dependencies := range constraintToDependencies {
		for _, dependency := range dependencies {
			if !slices.Contains(constraintToPatches[constraint], dependency.Version) {
				constraintToPatches[constraint] = append(constraintToPatches[constraint], dependency.Version)
			}
		}

		sort.Slice(constraintToPatches[constraint], func(i, j int) bool {
			iVersion := semver.MustParse(constraintToPatches[constraint][i])
			jVersion := semver.MustParse(constraintToPatches[constraint][j])
			return iVersion.LessThan(jVersion)
		})

		if constraint.Patches < len(constraintToPatches[constraint]) {
			constraintToPatches[constraint] = constraintToPatches[constraint][len(constraintToPatches[constraint])-constraint.Patches:]
		}
	}

	var patchesToKeep []string
	for _, versions := range constraintToPatches {
		patchesToKeep = append(patchesToKeep, versions...)
	}

	var dependenciesToKeep []cargo.ConfigMetadataDependency

	for _, dependency := range config.Metadata.Dependencies {
		if slices.Contains(patchesToKeep, dependency.Version) {
			dependenciesToKeep = append(dependenciesToKeep, dependency)
		}
	}

	// Sort the stacks within the dependency
	for _, dependency := range dependenciesToKeep {
		slices.Sort(dependency.Stacks)
	}

	// Sort the dependencies by:
	// 1. ID
	// 2. Version
	// 3. len(Stacks)
	sort.Slice(dependenciesToKeep, func(i, j int) bool {
		dep1 := dependenciesToKeep[i]
		dep2 := dependenciesToKeep[j]
		if dep1.ID == dep2.ID {
			dep1Version := semver.MustParse(dep1.Version)
			dep2Version := semver.MustParse(dep2.Version)

			if dep1Version.Equal(dep2Version) {
				return len(dep1.Stacks) < len(dep2.Stacks)
			}

			return dep1Version.LessThan(dep2Version)
		}
		return dep1.ID < dep2.ID
	})

	config.Metadata.Dependencies = dependenciesToKeep

	return config
}
