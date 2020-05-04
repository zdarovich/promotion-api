package errorcodes

import "fmt"

type (
	CodeError struct {
		ErrorField string
		ErrorCode  int
	}
	GenericError struct {
		Generic   error
		ErrorCode int
	}
)

func New(ErrorField string, ErrorCode int) error {
	return &CodeError{
		ErrorField: ErrorField,
		ErrorCode:  ErrorCode,
	}
}

func Wrap(error error, ErrorCode int) error {
	return &GenericError{
		Generic:   error,
		ErrorCode: ErrorCode,
	}
}

func (e *CodeError) Error() string {
	return fmt.Sprintf("errorField %s: ErrorCode: %d", e.ErrorField, e.ErrorCode)
}

func (e *GenericError) Error() string {
	return fmt.Sprintf("message %s: ErrorCode: %d", e.Generic.Error(), e.ErrorCode)
}
