package integration_test

import (
	"os"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/paketo-buildpacks/libdependency/github"
	"github.com/paketo-buildpacks/libdependency/retrieve"
	"github.com/sclevine/spec"
)

func testGithubReleases(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect          = NewWithT(t).Expect
		allVersionsFunc retrieve.GetAllVersionsFunc
	)

	context("GetReleasesFromGithub", func() {
		context("nodejs/node", func() {
			it.Before(func() {
				allVersionsFunc = github.GetAllVersions(os.Getenv("GIT_TOKEN"), "nodejs", "node")
			})

			it("will return a list of github releases", func() {
				// https://github.com/nodejs/node/releases
				fromGithub, err := allVersionsFunc()
				Expect(err).NotTo(HaveOccurred())
				Expect(fromGithub).NotTo(BeNil())

				Expect(fromGithub.GetVersionStrings()).To(ContainElements("18.7.0", "6.11.3"))
				Expect(len(fromGithub) > 300).To(BeTrue())
			})
		})

		context("curl/curl", func() {
			it.Before(func() {
				allVersionsFunc = github.GetAllVersions(os.Getenv("GIT_TOKEN"), "curl", "curl")
			})

			it("will return a list of github releases", func() {
				// https://github.com/curl/curl/releases
				fromGithub, err := allVersionsFunc()
				Expect(err).NotTo(HaveOccurred())
				Expect(fromGithub).NotTo(BeNil())

				Expect(fromGithub.GetVersionStrings()).To(ContainElements("7.78.0", "7.64.0"))
				Expect(len(fromGithub) > 40).To(BeTrue())
			})
		})

		context("failure cases", func() {
			context("non-existing org/space", func() {
				it.Before(func() {
					allVersionsFunc = github.GetAllVersions(os.Getenv("GIT_TOKEN"), "a612403a", "99b59f037720")
				})

				it("will return error", func() {
					_, err := allVersionsFunc()
					Expect(err).To(HaveOccurred())
				})
			})
		})
	})
}
