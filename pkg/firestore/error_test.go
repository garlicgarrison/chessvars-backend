package firestore

import (
	"errors"
	"testing"
)

func TestInvalidCursorError(t *testing.T) {
	other := errors.New("error")
	invalidCursor := NewInvalidCursorError(errors.New("invalid cursor"))

	t.Run("invalid cursor", func(t *testing.T) {
		if !IsInvalidCursorError(invalidCursor) {
			t.Fatal("unexpected false")
		}
	})

	t.Run("other", func(t *testing.T) {
		if IsInvalidCursorError(other) {
			t.Fatal("unexpected true")
		}
	})

	t.Run("returned", func(t *testing.T) {
		returned := func() error {
			return invalidCursor
		}

		err := returned()
		if !IsInvalidCursorError(err) {
			t.Fatal("unexpected false")
		}
	})
}

func TestNotAllowedError(t *testing.T) {
	other := errors.New("error")
	notAllowed := NewNotAllowedError(errors.New("not allowed"))

	t.Run("not allowed", func(t *testing.T) {
		if !IsNotAllowedError(notAllowed) {
			t.Fatal("unexpected false")
		}
	})

	t.Run("other", func(t *testing.T) {
		if IsNotAllowedError(other) {
			t.Fatal("unexpected true")
		}
	})

	t.Run("returned", func(t *testing.T) {
		returned := func() error {
			return notAllowed
		}

		err := returned()
		if !IsNotAllowedError(err) {
			t.Fatal("unexpected false")
		}
	})
}
