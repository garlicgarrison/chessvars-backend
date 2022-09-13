package format

const (
	GAME_ID_PREFIX = "igam"
)

type GameIDType int

func (id GameIDType) IDMethod() IDMethod {
	return IDMETHOD_RANDOM
}

func (id GameIDType) Prefix() string {
	return GAME_ID_PREFIX
}

func (id GameIDType) Size() uint {
	return 32
}

type GameID string

func NewGameID() GameID {
	return GameID(NewID(GameIDType(0)).String())
}

func NewGameIDFromIdentifer(id string) GameID {
	return GameID(GAME_ID_PREFIX + id)
}

func ParseGameID(id string) (GameID, error) {
	parsed, err := ParseID(GameIDType(0), id)
	if err != nil {
		return "", err
	}

	return GameID(parsed.String()), nil
}

func (u GameID) String() string {
	return string(u)
}

func (u GameID) Identifier() string {
	return string(u[4:])
}
