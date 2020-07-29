package reader

import "fmt"

type ConnectionError struct {
	Err     error
	Message string
}

func NewConnectionError(err error, message string) *ConnectionError {
	return &ConnectionError{err, message}
}

func (e *ConnectionError) Error() string {
	return fmt.Sprintf("A connection error occurred while reading - %s: %s", e.Message, e.Err.Error())
}

func (e *ConnectionError) Unwrap() error {
	return e.Err
}
