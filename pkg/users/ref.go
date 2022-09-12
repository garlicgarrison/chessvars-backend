package users

import (
	"github.com/pafkiuq/backend/pkg/firestore"
	"github.com/pafkiuq/backend/pkg/format"
)

const (
	FS_USERS_COLL = "users"
)

func (s *Service) getUsersRef() *firestore.CollectionRef {
	return s.fs.Collection(FS_USERS_COLL)
}

func (s *Service) getUserRef(userID format.UserID) *firestore.DocumentRef {
	return s.getUsersRef().Doc(userID.String())
}
