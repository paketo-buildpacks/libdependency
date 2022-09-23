package integration_test

import (
	"os"
	"testing"

	"github.com/joshuatcasey/libdependency/github"
	"github.com/joshuatcasey/libdependency/retrieve"
	. "github.com/onsi/gomega"
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
