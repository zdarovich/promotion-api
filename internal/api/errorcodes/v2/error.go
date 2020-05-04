package errorcodes

import "fmt"

// CodeError struct
type CodeError struct {
	ErrorField       string
	ErrorCode        int
	ErrorDescription string
}

// New returns new error
func New(ErrorField string, ErrorCode int) *CodeError {

	c := &CodeError{
		ErrorCode:        ErrorCode,
		ErrorDescription: GetDescription(ErrorCode),
	}

	if len(ErrorField) > 0 {
		c.ErrorField = ErrorField
	}

	return c
}

// Error converts CodeError to error string
func (e *CodeError) Error() string {

	return fmt.Sprintf(
		"errorField %s: ErrorCode: %d ErrorDescription: %s",
		e.ErrorField,
		e.ErrorCode,
		e.ErrorDescription,
	)
}
