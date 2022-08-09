package workflows

import "encoding/json"

func ToWorkflowJson(item any) (string, error) {
	if bytes, err := json.Marshal(item); err != nil {
		return "", err
	} else {
		return string(bytes), nil
	}
}
