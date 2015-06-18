package warrant

import "github.com/pivotal-cf-experimental/warrant/internal/network"

func translateError(err error) error {
	switch err.(type) {
	case network.NotFoundError:
		return NotFoundError{err}
	case network.UnauthorizedError:
		return UnauthorizedError{err}
	case network.UnexpectedStatusError:
		return UnexpectedStatusError{err}
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
