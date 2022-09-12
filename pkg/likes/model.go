package likes

import (
	"time"

	"github.com/pafkiuq/backend/pkg/format"
)

const (
	FS_COLL_USERS  = "users"
	FS_COLL_VIDEOS = "videos"
)

type VideoLike struct {
	UserID    format.UserID  `firestore:"user_id"`
	VideoID   format.VideoID `firestore:"video_id"`
	Timestamp time.Time      `firestore:"timestamp"`
}
