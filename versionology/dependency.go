package versionology

import (
	"encoding/json"

	"github.com/Masterminds/semver/v3"
	"github.com/paketo-buildpacks/libdependency/collections"
	"github.com/paketo-buildpacks/packit/v2/cargo"
)

// Dependency exists as a way to "extend" cargo.ConfigMetadataDependency
// and also "implements" VersionFetcher
type Dependency struct {
	cargo.ConfigMetadataDependency
	SemverVersion *semver.Version `json:"-"`
	Target        string          `json:"target,omitempty"`
}

// MarshalJSON is needed because cargo.ConfigMetadataDependency has its own
// MarshalJSON (packit v2.26.0 and later). Embedding promotes it to Dependency,
// and encoding/json would then skip the Target field.
func (d Dependency) MarshalJSON() ([]byte, error) {
	content, err := json.Marshal(d.ConfigMetadataDependency)
	if err != nil {
		return nil, err
	}

	if d.Target == "" {
		return content, nil
	}

	target, err := json.Marshal(d.Target)
	if err != nil {
		return nil, err
	}

	field := append([]byte(`"target":`), target...)
	if string(content) == "{}" {
		return append(append([]byte("{"), field...), '}'), nil
	}

	// content is a JSON object, so it ends with '}'.
	content = append(content[:len(content)-1], ',')
	return append(append(content, field...), '}'), nil
}

// UnmarshalJSON is needed for the same reason as MarshalJSON.
func (d *Dependency) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &d.ConfigMetadataDependency); err != nil {
		return err
	}

	var extra struct {
		Target string `json:"target"`
	}
	if err := json.Unmarshal(data, &extra); err != nil {
		return err
	}

	d.Target = extra.Target
	return nil
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

func (d Dependency) Version() *semver.Version {
	return d.SemverVersion
}

func NewDependencyArray(configMetadataDependency cargo.ConfigMetadataDependency, targets ...string) ([]Dependency, error) {
	var dependencies []Dependency
	for _, target := range targets {
		dependency, err := NewDependency(configMetadataDependency, target)
		if err != nil {
			return nil, err
		}
		dependencies = append(dependencies, dependency)
	}

	return dependencies, nil
}

// Versions will return an array of strings that represents each version of the input
// Primarily intended as a test helper.
func Versions(dependencies []Dependency) []string {
	return collections.TransformFunc(dependencies, func(dep Dependency) string {
		return dep.SemverVersion.String()
	})
}
