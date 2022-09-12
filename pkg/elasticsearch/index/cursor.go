package index

import (
	"encoding/base64"
	"encoding/json"
)

type Cursor struct {
	// PID is the id of this request's Point-In-Time
	PID string `json:"pid"`
	// Offset is used to tell if all results of the query has been seen
	//
	// It is NOT sent to elasticsearch since we use "search_after" and not offset
	Offset int `json:"offset"`
	// After is the input for "search_after"
	//
	// Its value is dependent on the search request body's "sort" value and therefore should be simply
	// stored as a string and injected into the request body.
	After string `json:"after"`
	// Metadata contains application defined metadata for the cursor
	Metadata map[string]interface{} `json:"metadata"`
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
