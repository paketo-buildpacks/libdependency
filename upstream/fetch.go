package upstream

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/paketo-buildpacks/packit/v2/fs"
)

func GetAndUnmarshal(url string, v any) error {
	response, err := http.DefaultClient.Get(url)
	if err != nil {
		return fmt.Errorf("could not get project metadata: %w", err)
	}

	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to query url %s with: status code %d", url, response.StatusCode)
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("could not read response: %w", err)
	}

	err = json.Unmarshal(body, v)
	if err != nil {
		return fmt.Errorf("could not unmarshal project metadata: %w", err)
	}

	return nil
}

func GetSHA256OfRemoteFile(sourceURL string) (string, error) {
	resp, err := http.Get(sourceURL)
	if err != nil {
		return "", fmt.Errorf("failed to query url: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to query url %s with: status code %d", sourceURL, resp.StatusCode)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	tempDir, err := os.MkdirTemp("", "temp")
	if err != nil {
		return "", err
	}

	defer os.RemoveAll(tempDir)

	tempFilePath := filepath.Join(tempDir, "temp-file")
	err = os.WriteFile(tempFilePath, body, os.ModePerm)
	if err != nil {
		return "", err
	}

	calculator := fs.NewChecksumCalculator()
	return calculator.Sum(tempFilePath)
}
