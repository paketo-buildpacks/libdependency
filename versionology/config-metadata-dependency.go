package versionology

import (
	"github.com/Masterminds/semver/v3"
	"github.com/joshuatcasey/collections"
	"github.com/paketo-buildpacks/packit/v2/cargo"
)

type Dependency struct {
	cargo.ConfigMetadataDependency
	Version *semver.Version
}

func (d Dependency) GetVersion() *semver.Version {
	return d.Version
}

type DependencyArray []Dependency

func (array DependencyArray) Versions() []string {
	return collections.TransformFunc(array, func(dep Dependency) string {
		return dep.GetVersion().String()
	})
}
