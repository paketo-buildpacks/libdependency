package upstream

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func GetAndUnmarshal(url string, v any) error {
	response, err := http.DefaultClient.Get(url)
	if err != nil {
		return fmt.Errorf("could not get project metadata: %w", err)
	}

	if err != nil {
		return err
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
