package versionology_test

import (
	"testing"

	"github.com/joshuatcasey/libdependency/versionology"
	. "github.com/onsi/gomega"
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
}
