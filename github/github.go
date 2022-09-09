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
	"github.com/joshuatcasey/libdependency"
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

func GithubAllVersions(githubToken, org, repo string) libdependency.HasVersionsFunc {
	return func() (versionology.HasVersionArray, error) {
		return getReleasesFromGithub(githubToken, org, repo)
	}
}

// GetReleasesFromGithub will return all semver-compatible versions from the releases of the given repo
// as documented by https://docs.github.com/en/rest/releases/releases#list-releases
func getReleasesFromGithub(githubToken, org, repo string) (versionology.HasVersionArray, error) {
	client := &http.Client{}

	perPage := 100
	page := 1

	var allVersions []*semver.Version

	for ; ; page++ {
		urlString := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases?per_page=%d&page=%d", org, repo, perPage, page)
		req, err := http.NewRequest("GET", urlString, nil)
		if err != nil {
			panic(err)
		}

		req.Header.Set("Accept", "application/vnd.github.v3+json")
		if githubToken != "" {
			req.Header.Set("Authorization", fmt.Sprintf("token %s", githubToken))
		}

		res, err := client.Do(req)
		if err != nil {
			panic(err)
		}

		body, err := io.ReadAll(res.Body)
		if err != nil {
			panic(err)
		}

		err = res.Body.Close()
		if err != nil {
			panic(err)
		}

		var githubReleaseNames []GithubReleaseNamesDTO
		err = json.Unmarshal(body, &githubReleaseNames)
		if err != nil {
			panic(err)
		}

		for _, release := range githubReleaseNames {
			if version, err := SanitizeGithubReleaseName(release); err != nil {
				panic(err)
			} else {
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

	return collections.TransformFunc(allVersions, func(version *semver.Version) versionology.HasVersion {
		return versionology.NewSimpleHasVersion(version)
	}), nil
}
