package errors

import (
	"fmt"
	"strings"
)

func AggregateErrors(errors []error) error {
	var errMsgs []string
	for i, e := range errors {
		errMsgs = append(errMsgs, fmt.Sprintf("%d. Error: %v", i+1, e))
	}
	return fmt.Errorf("errors occurred during deletion:\n%s", strings.Join(errMsgs, "\n"))
}
