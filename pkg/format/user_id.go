package format

const (
	USER_ID_PREFIX = "iusr"
)

type UserIDType int

func (id UserIDType) IDMethod() IDMethod {
	return IDMETHOD_RANDOM
}

func (id UserIDType) Prefix() string {
	return USER_ID_PREFIX
}

func (id UserIDType) Size() uint {
	return 32
}

type UserID string

func NewUserID() UserID {
	return UserID(NewID(UserIDType(0)).String())
}

func NewUserIDFromIdentifer(id string) UserID {
	return UserID(USER_ID_PREFIX + id)
}

func ParseUserID(id string) (UserID, error) {
	parsed, err := ParseID(UserIDType(0), id)
	if err != nil {
		return "", err
	}

	return UserID(parsed.String()), nil
}

func (u UserID) String() string {
	return string(u)
}

func (u UserID) Identifier() string {
	return string(u[4:])
}
