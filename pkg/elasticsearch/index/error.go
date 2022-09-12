package index

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

type VersionConflictError struct {
	error
}

func NewVersionConflictError(err error) *VersionConflictError {
	return &VersionConflictError{err}
}

func IsVersionConflictError(err error) bool {
	_, ok := err.(*VersionConflictError)
	return ok
}
