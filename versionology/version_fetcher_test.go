package versionology_test

import (
	"testing"

	. "github.com/onsi/gomega"
	"github.com/paketo-buildpacks/libdependency/versionology"
	"github.com/sclevine/spec"
)

func testVersionFetcher(t *testing.T, context spec.G, it spec.S) {
	Expect := NewWithT(t).Expect

	context("NewSimpleVersionFetcherArray", func() {
		it("will return an array containing Semver Versions", func() {
			array, err := versionology.NewSimpleVersionFetcherArray("1.2.3", "4.5.6", "7.8.9")
			Expect(err).NotTo(HaveOccurred())
			Expect(array.GetVersionStrings()).To(ConsistOf("1.2.3", "4.5.6", "7.8.9"))
		})

		context("failure cases", func() {
			it("will return an error when a semver is invalid", func() {
				_, err := versionology.NewSimpleVersionFetcherArray("hi")
				Expect(err).To(HaveOccurred())
			})
		})
	})

	context("GetNewestVersion", func() {
		var (
			versions versionology.VersionFetcherArray
			err      error
		)

		it("will return the newest version, using semver ordering not lexical ordering", func() {
			versions, err = versionology.NewSimpleVersionFetcherArray("1.9", "1.10")
			Expect(err).NotTo(HaveOccurred())

			Expect(versions.GetNewestVersion()).To(Equal("1.10.0"))
		})

		it("will return the newest version", func() {
			versions, err = versionology.NewSimpleVersionFetcherArray("1.2.3", "4.5.6", "7.8.9")
			Expect(err).NotTo(HaveOccurred())

			Expect(versions.GetNewestVersion()).To(Equal("7.8.9"))
		})

		context("failure cases", func() {
			context("with no versions", func() {
				it.Before(func() {
					versions = versionology.NewVersionFetcherArray()
				})

				it("will return empty string", func() {
					Expect(versions.GetNewestVersion()).To(Equal(""))
				})
			})
		})
	})
}
