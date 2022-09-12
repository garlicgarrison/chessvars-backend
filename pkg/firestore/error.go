package firestore

type InvalidCursorError struct {
	error
}

func NewInvalidCursorError(err error) *InvalidCursorError {
	return &InvalidCursorError{err}
}

func IsInvalidCursorError(err error) bool {
	_, ok := err.(*InvalidCursorError)
	return ok
}

type NotAllowedError struct {
	error
}

func NewNotAllowedError(err error) *NotAllowedError {
	return &NotAllowedError{err}
}

func IsNotAllowedError(err error) bool {
	_, ok := err.(*NotAllowedError)
	return ok
}
