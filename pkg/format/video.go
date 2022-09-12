package format

import (
	"fmt"
)

type VideoFormat string

const (
	MP4  VideoFormat = "mp4"
	MOV  VideoFormat = "mov"
	AVI  VideoFormat = "avi"
	FLV  VideoFormat = "flv"
	MKV  VideoFormat = "mkv"
	WEBM VideoFormat = "webm"
)

func (f VideoFormat) String() string {
	return string(f)
}

func ParseVideoFormat(format string) (VideoFormat, error) {
	switch format {
	case "mp4", "MP4":
		return MP4, nil
	case "mov", "MOV":
		return MOV, nil
	case "avi", "AVI":
		return AVI, nil
	case "flv", "FLV":
		return FLV, nil
	case "mkv", "MKV":
		return MKV, nil
	case "webm", "WEBM":
		return WEBM, nil
	default:
		return "", fmt.Errorf("unknown format: %s", format)
	}
}
