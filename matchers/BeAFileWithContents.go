package matchers

import (
	"fmt"
	"os"

	"github.com/onsi/gomega/format"
)

type BeAFileWithContentsMatcher struct {
	expectedContents string
	actualContents   string
}

func BeAFileWithContents(expectedContents string) *BeAFileWithContentsMatcher {
	return &BeAFileWithContentsMatcher{
		expectedContents: expectedContents,
	}
}

func (matcher *BeAFileWithContentsMatcher) Match(actual interface{}) (success bool, err error) {
	actualFilename, ok := actual.(string)
	if !ok {
		return false, fmt.Errorf("BeARegularFileMatcher matcher expects a file path")
	}

	bytes, err := os.ReadFile(actualFilename)
	if err != nil {
		return false, err
	}

	matcher.actualContents = string(bytes)

	return matcher.actualContents == matcher.expectedContents, nil
}

func (matcher *BeAFileWithContentsMatcher) FailureMessage(actual interface{}) string {
	return format.Message(actual, fmt.Sprintf("to have contents '%s', but were '%s'", matcher.expectedContents, matcher.actualContents))
}

func (matcher *BeAFileWithContentsMatcher) NegatedFailureMessage(actual interface{}) string {
	return format.Message(actual, fmt.Sprintf("not to have contents: %s", matcher.expectedContents))
}
