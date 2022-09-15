package workflows

import "encoding/json"

// ToWorkflowJson will return a string containing JSON formatted as a GitHub workflow expects, with
// no whitespace outside of strings.
//
// Use this when printing output or writing a file intended for use by a workflow.
// https://github.com/orgs/community/discussions/26288
func ToWorkflowJson(item any) (string, error) {
	if bytes, err := json.Marshal(item); err != nil {
		return "", err
	} else {
		return string(bytes), nil
	}
}
