package github

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/joshuatcasey/collections"
	"github.com/joshuatcasey/libdependency/retrieve"
	"github.com/joshuatcasey/libdependency/versionology"
)

type GithubReleaseNamesDTO struct {
	Name    string `json:"name"`
	TagName string `json:"tag_name"`
}

// SanitizeGithubReleaseName will determine whether to use the name or the tag as the semver version
func SanitizeGithubReleaseName(release GithubReleaseNamesDTO) (*semver.Version, error) {
	if version, err := semver.NewVersion(strings.TrimSpace(release.Name)); err != nil {
		return semver.NewVersion(strings.TrimSpace(release.TagName))
	} else {
		return version, nil
	}
}

// GetAllVersions will return a libdependency.VersionFetcherFunc that can retrieve all versions for a given
// GitHub org/repo.
func GetAllVersions(githubToken, org, repo string) retrieve.GetAllVersionsFunc {
	return func() (versionology.VersionFetcherArray, error) {
		return getReleasesFromGithub(githubToken, org, repo)
	}
}

// getReleasesFromGithub will return all semver-compatible versions from the releases of the given repo
// as documented by https://docs.github.com/en/rest/releases/releases#list-releases
func getReleasesFromGithub(githubToken, org, repo string) (versionology.VersionFetcherArray, error) {
	client := &http.Client{}

	perPage := 100

	allVersions := make([]*semver.Version, 0)

	for page := 1; ; page++ {
		urlString := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases?per_page=%d&page=%d", org, repo, perPage, page)
		req, err := http.NewRequest("GET", urlString, nil)
		if err != nil {
			return versionology.NewVersionFetcherArray(), err
		}

		req.Header.Set("Accept", "application/vnd.github.v3+json")
		if githubToken != "" {
			req.Header.Set("Authorization", fmt.Sprintf("token %s", githubToken))
		}

		res, err := client.Do(req)
		if err != nil {
			return versionology.NewVersionFetcherArray(), err
		}

		if res.StatusCode != http.StatusOK {
			return versionology.NewVersionFetcherArray(),
				fmt.Errorf("failed to query url %s with: status code %d", urlString, res.StatusCode)
		}

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return versionology.NewVersionFetcherArray(), err
		}

		err = res.Body.Close()
		if err != nil {
			return versionology.NewVersionFetcherArray(), err
		}

		var githubReleaseNames []GithubReleaseNamesDTO
		err = json.Unmarshal(body, &githubReleaseNames)
		if err != nil {
			return versionology.NewVersionFetcherArray(), err
		}

		for _, release := range githubReleaseNames {
			if version, err := SanitizeGithubReleaseName(release); err == nil {
				allVersions = append(allVersions, version)
			}
		}

		if len(githubReleaseNames) < perPage {
			break
		}
	}

	sort.Slice(allVersions, func(i, j int) bool {
		return allVersions[i].GreaterThan(allVersions[j])
	})

	return collections.TransformFunc(allVersions, func(version *semver.Version) versionology.VersionFetcher {
		return versionology.NewSimpleVersionFetcher(version)
	}), nil
}
