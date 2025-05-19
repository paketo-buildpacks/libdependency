package versionology

import (
	"sort"

	"github.com/Masterminds/semver/v3"
	"github.com/joshuatcasey/collections"
)

// VersionFetcher exists to allow buildpack authors to use a type of their own choosing when passing structs
// in and out of libdependency APIs.
type VersionFetcher interface {
	Version() *semver.Version
}

// SimpleVersionFetcher only contains a Semver Version and implements VersionFetcher
type SimpleVersionFetcher struct {
	version *semver.Version
}

func (s SimpleVersionFetcher) Version() *semver.Version {
	return s.version
}

type VersionFetcherArray []VersionFetcher

func NewVersionFetcherArray() VersionFetcherArray {
	return make([]VersionFetcher, 0)
}

func (versions VersionFetcherArray) GetVersionStrings() []string {
	return VersionFetcherToString(versions)
}

func (versions VersionFetcherArray) GetNewestVersion() string {
	if len(versions) < 1 {
		return ""
	}
	temp := versions
	sort.Slice(temp, func(i, j int) bool {
		return temp[i].Version().LessThan(temp[j].Version())
	})
	return temp[len(temp)-1].Version().String()
}

func NewSimpleVersionFetcher(version *semver.Version) SimpleVersionFetcher {
	return SimpleVersionFetcher{
		version: version,
	}
}

// NewSimpleVersionFetcherArray will return a VersionFetcherArray containing the semver representation of the input
func NewSimpleVersionFetcherArray(versions ...string) (VersionFetcherArray, error) {
	return collections.TransformFuncWithError(versions, func(version string) (VersionFetcher, error) {
		semverVersion, err := semver.NewVersion(version)
		if err != nil {
			return nil, err
		}

		return NewSimpleVersionFetcher(semverVersion), nil
	})
}
