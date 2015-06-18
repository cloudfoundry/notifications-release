package network

import "fmt"

type RequestBodyMarshalError struct {
	err error
}

func newRequestBodyMarshalError(err error) RequestBodyMarshalError {
	return RequestBodyMarshalError{err: err}
}

func (e RequestBodyMarshalError) Error() string {
	return fmt.Sprintf("Warrant RequestBodyMarshalError: %v", e.err)
}

type RequestConfigurationError struct {
	err error
}

func newRequestConfigurationError(err error) RequestConfigurationError {
	return RequestConfigurationError{err: err}
}

func (e RequestConfigurationError) Error() string {
	return fmt.Sprintf("Warrant RequestConfigurationError: %v", e.err)
}

type RequestHTTPError struct {
	err error
}

func newRequestHTTPError(err error) RequestHTTPError {
	return RequestHTTPError{err: err}
}

func (e RequestHTTPError) Error() string {
	return fmt.Sprintf("Warrant RequestHTTPError: %v", e.err)
}

type ResponseReadError struct {
	err error
}

func newResponseReadError(err error) ResponseReadError {
	return ResponseReadError{err: err}
}

func (e ResponseReadError) Error() string {
	return fmt.Sprintf("Warrant ResponseReadError: %v", e.err)
}

type UnexpectedStatusError struct {
	status int
	body   []byte
}

func newUnexpectedStatusError(status int, body []byte) UnexpectedStatusError {
	return UnexpectedStatusError{
		status: status,
		body:   body,
	}
}

func (e UnexpectedStatusError) Error() string {
	return fmt.Sprintf("Warrant UnexpectedStatusError: %d %s", e.status, e.body)
}

type NotFoundError struct {
	message []byte
}

func newNotFoundError(message []byte) NotFoundError {
	return NotFoundError{message: message}
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("Warrant NotFoundError: %s", e.message)
}

type UnauthorizedError struct {
	message []byte
}

func newUnauthorizedError(message []byte) UnauthorizedError {
	return UnauthorizedError{message: message}
}

func (e UnauthorizedError) Error() string {
	return fmt.Sprintf("Warrant UnauthorizedError: %s", e.message)
}
