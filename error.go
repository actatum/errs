package errs

import (
	"errors"
	"fmt"
)

// Code represents different types of errors that can arise in applications.
type Code struct {
	slug string
}

// String returns the string representation of the error code.
func (c Code) String() string {
	return c.slug
}

var (
	Internal         = Code{"internal"}
	NotFound         = Code{"not_found"}
	Invalid          = Code{"invalid"}
	Unauthorized     = Code{"unauthorized"}
	PermissionDenied = Code{"permission_denied"}
	Conflict         = Code{"conflict"}
)

// CodeFromString parses a Code from the given string. If no code exists for the
// string 'Internal' is returned.
func CodeFromString(s string) Code {
	switch s {
	case NotFound.slug:
		return NotFound
	case Invalid.slug:
		return Invalid
	case Unauthorized.slug:
		return Unauthorized
	case PermissionDenied.slug:
		return PermissionDenied
	case Conflict.slug:
		return Conflict
	default:
		return Internal
	}
}

// Error represents an application-specific error. Application errors can be
// unwrapped by the caller to extract the code and message.
type Error struct {
	code    Code
	message string
}

// Code returns the errors underlying code.
func (e Error) Code() Code {
	return e.code
}

// Message returns the errors underlying message.
func (e Error) Message() string {
	return e.message
}

// Error implements the error interface.
func (e Error) Error() string {
	return fmt.Sprintf("error: code=%s message=%s", e.code.slug, e.message)
}

// ErrorCode unwraps an Error and returns its code.
// Non Error types always return 'Internal'.
func ErrorCode(err error) Code {
	var e Error
	if errors.As(err, &e) {
		return e.Code()
	}

	return Internal
}

// ErrorMessage unwraps an Error and returns its message.
// Non Error types always return 'internal error'.
func ErrorMessage(err error) string {
	var e Error
	if errors.As(err, &e) {
		return e.Message()
	}

	return "internal error"
}

// Errorf constructs a new error with the given code and formatted message.
func Errorf(code Code, format string, args ...interface{}) Error {
	return Error{
		code:    code,
		message: fmt.Sprintf(format, args...),
	}
}
