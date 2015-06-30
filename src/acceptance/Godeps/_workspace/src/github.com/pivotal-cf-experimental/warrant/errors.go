package warrant

import (
	"fmt"
	"net/http"

	"github.com/pivotal-cf-experimental/warrant/internal/network"
)

func translateError(err error) error {
	switch s := err.(type) {
	case network.NotFoundError:
		return NotFoundError{err}
	case network.UnauthorizedError:
		return UnauthorizedError{err}
	case network.UnexpectedStatusError:
		switch s.Status {
		case http.StatusBadRequest:
			return BadRequestError{err}
		default:
			return UnexpectedStatusError{err}
		}
	default:
		return UnknownError{err}
	}
}

type UnexpectedStatusError struct {
	err error
}

func (e UnexpectedStatusError) Error() string {
	return e.err.Error()
}

type UnauthorizedError struct {
	err error
}

func (e UnauthorizedError) Error() string {
	return e.err.Error()
}

type NotFoundError struct {
	err error
}

func (e NotFoundError) Error() string {
	return e.err.Error()
}

type UnknownError struct {
	err error
}

func (e UnknownError) Error() string {
	return e.err.Error()
}

type InvalidTokenError struct {
	err error
}

func (e InvalidTokenError) Error() string {
	return e.err.Error()
}

type MalformedResponseError struct {
	err error
}

func (e MalformedResponseError) Error() string {
	return fmt.Sprintf("malformed response: %s", e.err)
}

type BadRequestError struct {
	err error
}

func (e BadRequestError) Error() string {
	return fmt.Sprintf("bad request: %s", e.err.(network.UnexpectedStatusError).Body)
}
