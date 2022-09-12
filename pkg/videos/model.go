package videos

import "github.com/pafkiuq/backend/pkg/format"

type VideoState string

const (
	DRAFT   VideoState = "draft"
	LIVE    VideoState = "live"
	ARCHIVE VideoState = "archive"
	NONE    VideoState = "none"
)

type VideoDocument struct {
	VideoID     format.VideoID    `firestore:"video_id"`
	OwnerID     format.UserID     `firestore:"owner_id"`
	Description string            `firestore:"description"`
	URL         string            `firestore:"url"`
	Thumbnails  map[string]string `firestore:"thumbnails"`
	Tags        []string          `firestore:"tags"`
	State       VideoState        `firestore:"state"`
}
