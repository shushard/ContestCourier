package web

import "fmt"

type HTTPError struct {
	StatusCode int
	Body       []byte
	BodyErr    error
}

func newHTTPError(statusCode int, body []byte, bodyErr error) HTTPError {
	return HTTPError{
		StatusCode: statusCode,
		Body:       body,
		BodyErr:    bodyErr,
	}
}

func (e HTTPError) Error() string {
	if e.BodyErr != nil {
		return fmt.Sprintf("get status %d and can't read body: %v", e.StatusCode, e.BodyErr)
	}
	return fmt.Sprintf("get status %d with body: %s", e.StatusCode, string(e.Body))
}
