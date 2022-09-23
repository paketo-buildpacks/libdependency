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
	SemverVersion *semver.Version `json:"-"`
	Target        string          `json:"target,omitempty"`
}

func NewDependency(configMetadataDependency cargo.ConfigMetadataDependency, target string) (Dependency, error) {
	if semverVersion, err := semver.NewVersion(configMetadataDependency.Version); err != nil {
		return Dependency{}, err
	} else {
		return Dependency{
			ConfigMetadataDependency: configMetadataDependency,
			SemverVersion:            semverVersion,
			Target:                   target,
		}, nil
	}
}

func NewDependencyArray(configMetadataDependency cargo.ConfigMetadataDependency, target string) ([]Dependency, error) {
	dependency, err := NewDependency(configMetadataDependency, target)
	if err != nil {
		return []Dependency{}, err
	}

	return []Dependency{dependency}, nil
}

func (d Dependency) Version() *semver.Version {
	return d.SemverVersion
}

// Versions will return an array of strings that represents each version of the input
// Primarily intended as a test helper.
func Versions(dependencies []Dependency) []string {
	return collections.TransformFunc(dependencies, func(dep Dependency) string {
		return dep.SemverVersion.String()
	})
}
