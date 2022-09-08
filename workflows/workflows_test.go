package workflows_test

import (
	"testing"

	"github.com/paketo-buildpacks/libdependency/workflows"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
)

func testWorkflows(t *testing.T, context spec.G, it spec.S) {
	Expect := NewWithT(t).Expect

	context("ToWorkflowJson", func() {
		it("translates string array to json", func() {
			item := []string{
				"v1.2.3",
				"2.3.4",
				"v9.8.7",
				"string with spaces",
			}

			json, err := workflows.ToWorkflowJson(item)
			Expect(err).NotTo(HaveOccurred())
			Expect(json).To(Equal(`["v1.2.3","2.3.4","v9.8.7","string with spaces"]`))
		})

		it("translates object into json", func() {
			type inner struct {
				Str string
				I   int
			}

			type outer struct {
				Str     string
				I       int
				Inner   inner
				private string
			}

			item := outer{
				Str: "string with spaces",
				I:   -999,
				Inner: inner{
					Str: "inner string with \ttabs",
					I:   999,
				},
				private: "will not appear in json",
			}

			json, err := workflows.ToWorkflowJson(item)
			Expect(err).NotTo(HaveOccurred())
			Expect(json).To(Equal(`{"Str":"string with spaces","I":-999,"Inner":{"Str":"inner string with \ttabs","I":999}}`))
		})
	})
}
