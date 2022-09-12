package format

import (
	"fmt"
)

type ImageFormat string

const (
	JPEG ImageFormat = "jpg"
	PNG  ImageFormat = "png"
	HEIC ImageFormat = "heic"
)

func (f ImageFormat) String() string {
	return string(f)
}

func ParseImageFormat(format string) (ImageFormat, error) {
	switch format {
	case "jpeg", "jpg", "JPG", "JPEG":
		return JPEG, nil
	case "png", "PNG":
		return PNG, nil
	case "heic", "HEIC":
		return HEIC, nil
	default:
		return "", fmt.Errorf("unknown format: %s", format)
	}
}
