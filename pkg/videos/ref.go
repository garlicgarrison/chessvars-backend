package videos

import (
	"github.com/pafkiuq/backend/pkg/firestore"
	"github.com/pafkiuq/backend/pkg/format"
)

const (
	FS_VIDEO_COLL = "videos"
)

func (s *Service) getVideoRef(videoID format.VideoID) *firestore.DocumentRef {
	return s.fs.Collection(FS_VIDEO_COLL).
		Doc(videoID.String())
}
