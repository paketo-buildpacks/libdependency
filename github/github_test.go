package github_test

import (
	"testing"

	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
)

func testGithub(t *testing.T, context spec.G, it spec.S) {
	Expect := NewWithT(t).Expect

	context("SanitizeGithubReleaseName", func() {
		context("with both a name and a tag", func() {
			context("when the name is a valid semver", func() {
				it("will choose the name", func() {
					releaseName, err := SanitizeGithubReleaseName(GithubReleaseNamesDTO{
						Name:    "1.2.3",
						TagName: "2.3.4",
					})

					Expect(err).NotTo(HaveOccurred())
					Expect(releaseName).NotTo(BeNil())
					Expect(releaseName.String()).To(Equal("1.2.3"))
				})
			})

			context("when the name is not a valid semver", func() {
				it("will choose the tag", func() {
					releaseName, err := SanitizeGithubReleaseName(GithubReleaseNamesDTO{
						Name:    "Release 1.0! 😊",
						TagName: "2.3.4",
					})

					Expect(err).NotTo(HaveOccurred())
					Expect(releaseName).NotTo(BeNil())
					Expect(releaseName.String()).To(Equal("2.3.4"))
				})
			})

			context("when the name present but empty", func() {
				it("will choose the tag", func() {
					releaseName, err := SanitizeGithubReleaseName(GithubReleaseNamesDTO{
						Name:    "",
						TagName: "2.3.4",
					})

					Expect(err).NotTo(HaveOccurred())
					Expect(releaseName).NotTo(BeNil())
					Expect(releaseName.String()).To(Equal("2.3.4"))
				})
			})

		})

		context("when there's only a tag", func() {
			it("will return the tag", func() {
				releaseName, err := SanitizeGithubReleaseName(GithubReleaseNamesDTO{
					TagName: "1.2.3",
				})

				Expect(err).NotTo(HaveOccurred())
				Expect(releaseName).NotTo(BeNil())
				Expect(releaseName.String()).To(Equal("1.2.3"))
			})

			it("will strip a leading 'v' from the tag", func() {
				releaseName, err := SanitizeGithubReleaseName(GithubReleaseNamesDTO{
					TagName: "v2.3.4",
				})

				Expect(err).NotTo(HaveOccurred())
				Expect(releaseName).NotTo(BeNil())
				Expect(releaseName.String()).To(Equal("2.3.4"))
			})

			it("will strip leading whitespace", func() {
				releaseName, err := SanitizeGithubReleaseName(GithubReleaseNamesDTO{
					TagName: "   3.4.5",
				})

				Expect(err).NotTo(HaveOccurred())
				Expect(releaseName).NotTo(BeNil())
				Expect(releaseName.String()).To(Equal("3.4.5"))
			})

			it("will strip trailing whitespace", func() {
				releaseName, err := SanitizeGithubReleaseName(GithubReleaseNamesDTO{
					TagName: "   4.5.6",
				})

				Expect(err).NotTo(HaveOccurred())
				Expect(releaseName).NotTo(BeNil())
				Expect(releaseName.String()).To(Equal("4.5.6"))
			})
		})

		context("when there's only a name", func() {
			it("will return the name", func() {
				releaseName, err := SanitizeGithubReleaseName(GithubReleaseNamesDTO{
					Name: "9.8.7",
				})

				Expect(err).NotTo(HaveOccurred())
				Expect(releaseName).NotTo(BeNil())
				Expect(releaseName.String()).To(Equal("9.8.7"))
			})

			it("will strip a leading 'v' from the name", func() {
				releaseName, err := SanitizeGithubReleaseName(GithubReleaseNamesDTO{
					Name: "v8.7.6",
				})

				Expect(err).NotTo(HaveOccurred())
				Expect(releaseName).NotTo(BeNil())
				Expect(releaseName.String()).To(Equal("8.7.6"))
			})

			it("will strip leading whitespace", func() {
				releaseName, err := SanitizeGithubReleaseName(GithubReleaseNamesDTO{
					Name: "   7.6.5",
				})

				Expect(err).NotTo(HaveOccurred())
				Expect(releaseName).NotTo(BeNil())
				Expect(releaseName.String()).To(Equal("7.6.5"))
			})

			it("will strip trailing whitespace", func() {
				releaseName, err := SanitizeGithubReleaseName(GithubReleaseNamesDTO{
					Name: "   6.5.4",
				})

				Expect(err).NotTo(HaveOccurred())
				Expect(releaseName).NotTo(BeNil())
				Expect(releaseName.String()).To(Equal("6.5.4"))
			})
		})

		context("failure cases", func() {
			it("will return the error", func() {
				_, err := SanitizeGithubReleaseName(GithubReleaseNamesDTO{
					TagName: "not a semver",
				})

				Expect(err).To(MatchError("Invalid Semantic Version"))
			})
		})
	})
}
