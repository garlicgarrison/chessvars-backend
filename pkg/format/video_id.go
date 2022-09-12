package format

const (
	VIDEO_ID_PREFIX = "ivid"
)

type VideoIDType int

func (id VideoIDType) IDMethod() IDMethod {
	return IDMETHOD_RANDOM
}

func (id VideoIDType) Prefix() string {
	return USER_ID_PREFIX
}

func (id VideoIDType) Size() uint {
	return 32
}

type VideoID string

func NewVideoID() VideoID {
	return VideoID(NewID(VideoIDType(0)).String())
}

func NewVideoIDFromIdentifer(id string) VideoID {
	return VideoID(VIDEO_ID_PREFIX + id)
}

func ParseVideoID(id string) (VideoID, error) {
	parsed, err := ParseID(VideoIDType(0), id)
	if err != nil {
		return "", err
	}

	return VideoID(parsed.String()), nil
}

func (u VideoID) String() string {
	return string(u)
}

func (u VideoID) Identifier() string {
	return string(u[4:])
}
