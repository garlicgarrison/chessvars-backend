package users

import (
	"github.com/garlicgarrison/chessvars-backend/pkg/firestore"
	"github.com/garlicgarrison/chessvars-backend/pkg/format"
)

const (
	FS_USERS_COLL = "users"
)

func (s *service) getUsersRef() *firestore.CollectionRef {
	return s.fs.Collection(FS_USERS_COLL)
}

func (s *service) getUserRef(userID format.UserID) *firestore.DocumentRef {
	return s.getUsersRef().Doc(userID.String())
}
