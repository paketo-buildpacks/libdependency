package versionology

import "github.com/Masterminds/semver/v3"

// HasVersion exists to allow buildpack authors to use a type of their own choosing when passing structs
// in and out of libdependency APIs.
type HasVersion interface {
	Version() *semver.Version
}

type SimpleHasVersion struct {
	version *semver.Version
}

type HasVersionArray []HasVersion

func (array HasVersionArray) GetVersionStrings() []string {
	return HasVersionToString(array)
}

func NewSimpleHasVersion(version *semver.Version) SimpleHasVersion {
	return SimpleHasVersion{
		version: version,
	}
}

func (s SimpleHasVersion) Version() *semver.Version {
	return s.version
}
