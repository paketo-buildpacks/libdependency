package versionology

import "github.com/Masterminds/semver/v3"

// VersionFetcher exists to allow buildpack authors to use a type of their own choosing when passing structs
// in and out of libdependency APIs.
type VersionFetcher interface {
	Version() *semver.Version
}

// SimpleVersionFetcher only contains a semverVersion and implements VersionFetcher
type SimpleVersionFetcher struct {
	version *semver.Version
}

type VersionFetcherArray []VersionFetcher

func NewVersionFetcherArray() []VersionFetcher {
	return make([]VersionFetcher, 0)
}

func (array VersionFetcherArray) GetVersionStrings() []string {
	return VersionFetcherToString(array)
}

func NewSimpleVersionFetcher(version *semver.Version) SimpleVersionFetcher {
	return SimpleVersionFetcher{
		version: version,
	}
}

func (s SimpleVersionFetcher) Version() *semver.Version {
	return s.version
}
