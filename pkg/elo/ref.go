package elo

import (
	"github.com/garlicgarrison/chessvars-backend/pkg/firestore"
	"github.com/garlicgarrison/chessvars-backend/pkg/format"
	"github.com/garlicgarrison/chessvars-backend/pkg/game"
	"github.com/garlicgarrison/chessvars-backend/pkg/users"
)

const (
	FS_ELO_COLL = "elo"
)

func (s *service) getElosRef(userID format.UserID) *firestore.CollectionRef {
	return s.fs.Collection(users.FS_USERS_COLL).Doc(userID.String()).Collection(FS_ELO_COLL)
}

func (s *service) getEloRef(userID format.UserID, game game.GameType) *firestore.DocumentRef {
	return s.getElosRef(userID).Doc(game.String())
}
