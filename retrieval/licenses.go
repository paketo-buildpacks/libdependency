package retrieval

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"

	"github.com/go-enry/go-license-detector/v4/licensedb"
	"github.com/go-enry/go-license-detector/v4/licensedb/filer"
)

type DecompressArtifactFunc func(artifact io.Reader, destination string) error

// LookupLicenses uses licensedb to detect licenses contained within a compressed directory
func LookupLicenses(sourceURL string, f DecompressArtifactFunc) []interface{} {
	// getting the dependency artifact from sourceURL
	resp, err := http.Get(sourceURL)
	if err != nil {
		panic(fmt.Errorf("failed to query url: %w", err))
	}
	if resp.StatusCode != http.StatusOK {
		panic(fmt.Errorf("failed to query url %s with: status code %d", sourceURL, resp.StatusCode))
	}

	// decompressing the dependency artifact
	tempDir, err := os.MkdirTemp("", "destination")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(tempDir)

	err = f(resp.Body, tempDir)
	if err != nil {
		panic(err)
	}

	// scanning artifact for license file
	filer, err := filer.FromDirectory(tempDir)
	if err != nil {
		panic(fmt.Errorf("failed to setup a licensedb filer: %w", err))
	}

	licenses, err := licensedb.Detect(filer)
	// if no licenses are found, just return an empty slice.
	if err != nil {
		if err.Error() != "no license file was found" {
			panic(fmt.Errorf("failed to detect licenses: %w", err))
		}
		return []interface{}{}
	}

	// Only return the license IDs, in alphabetical order
	var licenseIDs []string
	for key := range licenses {
		licenseIDs = append(licenseIDs, key)
	}
	sort.Strings(licenseIDs)

	var licenseIDsAsInterface []interface{}
	for _, licenseID := range licenseIDs {
		licenseIDsAsInterface = append(licenseIDsAsInterface, licenseID)
	}

	return licenseIDsAsInterface
}
