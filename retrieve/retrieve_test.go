package retrieve_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/joshuatcasey/libdependency/matchers"
	"github.com/joshuatcasey/libdependency/retrieve"
	"github.com/paketo-buildpacks/packit/v2/cargo"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
)

func testRetrieve(t *testing.T, context spec.G, it spec.S) {
	Expect := NewWithT(t).Expect

	var (
		savedArgs  []string
		outputFile string
	)

	it.Before(func() {
		outputFile = filepath.Join(t.TempDir(), "temp.json")

		savedArgs = os.Args
		t.Cleanup(func() {
			os.Args = savedArgs
		})

		os.Args = []string{"/path/to-binary",
			"--buildpack-toml-path", filepath.Join("..", "testdata", "empty", "buildpack.toml"),
			"--id", "depId",
			"--name", "depName",
			"--output-file", outputFile}
	})

	context("RetrieveNewMetadata", func() {
		it("will write the output to the given location", func() {
			retrieveNewMetadata := func(id, name string, config cargo.Config) ([]cargo.ConfigMetadataDependency, error) {
				Expect(id).To(Equal("depId"))
				Expect(name).To(Equal("depName"))

				return []cargo.ConfigMetadataDependency{
					{
						Version: "1.1.1",
					},
					{
						Version: "2.2.2",
					},
				}, nil
			}

			retrieve.NewMetadata(retrieveNewMetadata)
			Expect(outputFile).To(matchers.BeAFileWithContents(`[{"version":"1.1.1"},{"version":"2.2.2"}]`))
		})
	})
}
