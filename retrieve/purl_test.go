package retrieve_test

import (
	"testing"

	"github.com/joshuatcasey/libdependency/retrieve"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
)

func testPurl(t *testing.T, context spec.G, it spec.S) {
	Expect := NewWithT(t).Expect

	context("GeneratePURL", func() {
		it("will generate a purl", func() {
			purl := retrieve.GeneratePURL("NAME", "VERSION", "CHECKSUM", "source")

			Expect(purl).To(Equal("pkg:generic/NAME@VERSION?checksum=CHECKSUM&download_url=source"))
		})
	})

}
