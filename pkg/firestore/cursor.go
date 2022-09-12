package firestore

import (
	"encoding/base64"
	"encoding/json"
	"time"
)

type Cursor struct {
	// ID of first item of the next list
	Current *string `json:"current,omitempty"`
	// Timestamp of the last item previously returned
	After *time.Time `json:"after,omitempty"`
}

func DecodeCursor(str string) (*Cursor, error) {
	b, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return nil, NewInvalidCursorError(err)
	}

	if len(b) == 0 {
		return nil, nil
	}

	var cursor Cursor
	err = json.Unmarshal(b, &cursor)
	if err != nil {
		return nil, NewInvalidCursorError(err)
	}

	return &cursor, nil
}

func (c *Cursor) String() string {
	b, err := json.Marshal(c)
	if err != nil {
		return ""
	}

	return string(b)
}

func (c *Cursor) Encode() string {
	b, err := json.Marshal(c)
	if err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(b)
}
