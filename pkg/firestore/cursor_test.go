package firestore

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCursorString(t *testing.T) {
	tests := []struct {
		name  string
		curr  func() *string
		after func() *time.Time
	}{
		{
			name: "Current Only",
			curr: func() *string {
				str := "current"
				return &str
			},
			after: func() *time.Time {
				return nil
			},
		},
		{
			name: "After Only",
			curr: func() *string {
				return nil
			},
			after: func() *time.Time {
				after := time.Now().UTC()
				return &after
			},
		},
		{
			name: "All",
			curr: func() *string {
				curr := "current"
				return &curr
			},
			after: func() *time.Time {
				after := time.Now().UTC()
				return &after
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cursor := Cursor{
				Current: tt.curr(),
				After:   tt.after(),
			}

			decoded, err := DecodeCursor(cursor.Encode())
			if err != nil {
				t.Fatal("Error decoding", err)
			}

			assert.Equal(t, &cursor, decoded)
		})
	}

}

func TestEmptyCursor(t *testing.T) {
	parsed, err := DecodeCursor("")
	if err != nil {
		t.Fatal("Unexpected error", err)
	}

	if parsed != nil {
		t.Fatal("Unexpected cursor", parsed)
	}
}

func TestCursorError(t *testing.T) {
	_, err := DecodeCursor("hello")
	if err == nil {
		t.Fatal("unexpected nil error")
	}

	if !IsInvalidCursorError(err) {
		t.Fatal("unexpected error", err)
	}
}
