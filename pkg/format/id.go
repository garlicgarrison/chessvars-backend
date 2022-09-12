package format

import "fmt"

func invalidIDError(id, reason string) error {
	return fmt.Errorf("invalid id: %s -- %s", id, reason)
}

type IDMethod int

const (
	IDMETHOD_RANDOM = iota
	IDMETHOD_HASH
)

// IDType is the interface for creating a new type of id in fwaygo
//
// In order to use this interface,
type IDType interface {
	IDMethod() IDMethod
	// Prefix must return the prefix of the this id type
	//
	// Example: "iusr"
	//
	// Note: Ensure there are no collisions with other id types
	Prefix() string
	Size() uint
}

type ID[T IDType] struct {
	idType     T
	identifier string
}

func NewID[T IDType](idType T, args ...interface{}) ID[T] {
	switch idType.IDMethod() {
	case IDMETHOD_HASH:
		return ID[T]{
			idType:     idType,
			identifier: Hash(args),
		}
	case IDMETHOD_RANDOM:
		fallthrough
	default:
		return ID[T]{
			idType:     idType,
			identifier: Random(idType.Size()),
		}
	}
}

func ParseID[T IDType](idType T, id string) (ID[T], error) {
	prefix := idType.Prefix()

	if len(id) <= len(prefix) {
		return ID[T]{}, invalidIDError(id, "insufficient characters")
	}

	if id[:len(prefix)] != prefix {
		return ID[T]{}, invalidIDError(id,
			fmt.Sprintf("does not match prefix: %s", prefix))
	}

	return ID[T]{
		idType:     idType,
		identifier: id[len(prefix):],
	}, nil
}

func (id ID[T]) String() string {
	return id.idType.Prefix() + id.identifier
}

func (id ID[T]) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%q", id.String())), nil
}

func (id *ID[T]) UnmarshalJSON(b []byte) error {
	if len(b) < 4 {
		return fmt.Errorf("invalid ID: %s", b)
	}

	// Remove quote marks
	b = b[1 : len(b)-1]

	parsed, err := ParseID(id.idType, string(b))
	*id = parsed
	return err
}
