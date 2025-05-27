package ccvalidation

import (
	"errors"
	"fmt"
	"strings"
)

// Result holds validation results and errors
type Result struct {
	errors []error
}

// Error implements the error interface. Returns all error messages joined by commas
func (r Result) Error() string {
	if r.IsSuccess() {
		return ""
	}
	var messages []string
	for _, err := range r.errors {
		messages = append(messages, err.Error())
	}
	return strings.Join(messages, ", ")
}

// Is checks if an error is of type Result
func (r Result) Is(err error) bool {
	_, ok := err.(Result)
	return ok
}

// AddErrorMessage adds a validation error message to the Result
// If the message is empty or whitespace, nothing is added
func (r *Result) AddErrorMessage(msg string) {
	if strings.TrimSpace(msg) == "" {
		return
	}
	err := errors.New(msg)
	r.AddError(err)
}

// AddError adds a validation error to the Result
// If the error is nil, nothing is added
func (r *Result) AddError(err error) {
	if err == nil {
		return
	}
	r.errors = append(r.errors, err)
}

// AddInvalidParameterError adds an invalid parameter value message to the result
// name is the name of the parameter that is invalid
// value is the invalid value for the stated parameter
func (r *Result) AddInvalidParameterError(name string, value any) {
	if strings.TrimSpace(name) == "" {
		return
	}
	var valueStr string
	if s, ok := value.(fmt.Stringer); ok {
		valueStr = s.String()
	} else {
		valueStr = fmt.Sprintf("%v", value)
	}
	err := fmt.Errorf("%s is not valid: %s", name, valueStr)
	r.AddError(err)
}

// IsSuccess returns true when no error has been added to the result
func (r Result) IsSuccess() bool {
	return len(r.errors) == 0
}

// IsFailure returns true when any error has been added to the result
func (r Result) IsFailure() bool {
	return !r.IsSuccess()
}

// GetErrors returns a list of all errors in the result
// If no errors are found, returns an empty slice
func (r Result) GetErrors() []error {
	return r.errors
}

// GetErrorMessages returns a list of all error messages in the result
// If no errors are found, returns an empty slice
func (r Result) GetErrorMessages() []string {
	messages := make([]string, len(r.errors))
	for i, err := range r.errors {
		messages[i] = err.Error()
	}
	return messages
}
