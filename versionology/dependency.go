package versionology

import (
	"github.com/Masterminds/semver/v3"
	"github.com/joshuatcasey/collections"
	"github.com/paketo-buildpacks/packit/v2/cargo"
)

// Dependency exists as a way to "extend" cargo.ConfigMetadataDependency
// and also "implements" VersionFetcher
type Dependency struct {
	cargo.ConfigMetadataDependency
	version *semver.Version
	Target  string `json:"target,omitempty"`
}

func NewDependency(dependency cargo.ConfigMetadataDependency, target string) (Dependency, error) {
	if semverVersion, err := semver.NewVersion(dependency.Version); err != nil {
		return Dependency{}, err
	} else {
		return Dependency{
			ConfigMetadataDependency: dependency,
			version:                  semverVersion,
			Target:                   target,
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
