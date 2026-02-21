package internal

import (
	"fmt"
	"os"
)

const (
	ExitSuccess    = 0
	ExitError      = 1
	ExitAuth       = 2
	ExitNotFound   = 3
	ExitValidation = 4
)

type ZohoCliError struct {
	Message  string
	ExitCode int
}

func (e *ZohoCliError) Error() string {
	return e.Message
}

func NewError(msg string) *ZohoCliError {
	return &ZohoCliError{Message: msg, ExitCode: ExitError}
}

func NewAuthError(msg string) *ZohoCliError {
	return &ZohoCliError{Message: msg, ExitCode: ExitAuth}
}

func NewNotFoundError(msg string) *ZohoCliError {
	return &ZohoCliError{Message: msg, ExitCode: ExitNotFound}
}

func NewValidationError(msg string) *ZohoCliError {
	return &ZohoCliError{Message: msg, ExitCode: ExitValidation}
}

type ZohoAPIError struct {
	ZohoCliError
	StatusCode int
}

func NewAPIError(statusCode int, body string) *ZohoAPIError {
	exitCode := ExitError
	switch statusCode {
	case 401:
		exitCode = ExitAuth
	case 404:
		exitCode = ExitNotFound
	}
	return &ZohoAPIError{
		ZohoCliError: ZohoCliError{
			Message:  fmt.Sprintf("Zoho API error %d: %s", statusCode, body),
			ExitCode: exitCode,
		},
		StatusCode: statusCode,
	}
}

func Err(msg string) {
	fmt.Fprintln(os.Stderr, msg)
}

func Die(err error) {
	if e, ok := err.(*ZohoCliError); ok {
		Err(e.Message)
		os.Exit(e.ExitCode)
	}
	if e, ok := err.(*ZohoAPIError); ok {
		Err(e.Message)
		os.Exit(e.ExitCode)
	}
	Err(err.Error())
	os.Exit(ExitError)
}
