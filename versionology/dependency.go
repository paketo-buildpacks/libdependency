package versionology

import (
	"github.com/Masterminds/semver/v3"
	"github.com/joshuatcasey/collections"
	"github.com/paketo-buildpacks/packit/v2/cargo"
)

// Dependency exists as a way to "extend" cargo.ConfigMetadataDependency
// and also "implement" HasVersion
type Dependency struct {
	cargo.ConfigMetadataDependency
	version *semver.Version
}

func NewDependency(dependency cargo.ConfigMetadataDependency) (Dependency, error) {
	if semverVersion, err := semver.NewVersion(dependency.Version); err != nil {
		return Dependency{}, err
	} else {
		return Dependency{
			ConfigMetadataDependency: dependency,
			version:                  semverVersion,
		}, nil
	}
}

func (d Dependency) Version() *semver.Version {
	return d.version
}

type DependencyArray []Dependency

func (array DependencyArray) Versions() []string {
	return collections.TransformFunc(array, func(dep Dependency) string {
		return dep.Version().String()
	})
}
