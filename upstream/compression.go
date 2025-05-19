package upstream

import (
	"fmt"
	"io"

	"github.com/paketo-buildpacks/packit/v2/vacation"
)

func DefaultDecompress(artifact io.Reader, destination string) error {
	archive := vacation.NewArchive(artifact)

	err := archive.StripComponents(1).Decompress(destination)
	if err != nil {
		return fmt.Errorf("failed to decompress source file: %w", err)
	}

	return nil
}
