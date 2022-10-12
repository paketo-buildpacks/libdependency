package upstream_test

import (
	"testing"

	"github.com/joshuatcasey/libdependency/upstream"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
)

func testFetch(t *testing.T, context spec.G, it spec.S) {
	Expect := NewWithT(t).Expect

	context("GetSHA256OfRemoteFile", func() {
		it("works for curl", func() {
			sha256, err := upstream.GetSHA256OfRemoteFile("https://curl.se/download/curl-7.85.0.tar.gz")
			Expect(err).NotTo(HaveOccurred())
			Expect(sha256).To(Equal("78a06f918bd5fde3c4573ef4f9806f56372b32ec1829c9ec474799eeee641c27"))
		})

		it("works for curl", func() {
			sha256, err := upstream.GetSHA256OfRemoteFile("https://www.python.org/ftp/python/3.10.7/Python-3.10.7.tgz")
			Expect(err).NotTo(HaveOccurred())
			Expect(sha256).To(Equal("1b2e4e2df697c52d36731666979e648beeda5941d0f95740aafbf4163e5cc126"))
		})

		context("failure cases", func() {
			it("returns error when file not found", func() {
				_, err := upstream.GetSHA256OfRemoteFile("https://example.com/hello")
				Expect(err).To(MatchError("failed to query url https://example.com/hello with: status code 404"))
			})
		})
	})
}
