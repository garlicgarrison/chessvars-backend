package elo

import (
	"time"

	"github.com/garlicgarrison/chessvars-backend/pkg/firestore"
	"github.com/garlicgarrison/chessvars-backend/pkg/format"
	"github.com/garlicgarrison/chessvars-backend/pkg/users"
)

const (
	FS_ELO_COLL        = "elo"
	FS_GAME_ELOS_COLL  = "elos"
	FS_CURRENT_ELO_DOC = "current"
)

// user collection ->
// user doc ->
// elo collection (janggi/shogi or other variants) ->
// elo document ->
// elos collection ->
// elo document

func (s *service) getElosRef(userID format.UserID) *firestore.CollectionRef {
	return s.fs.Collection(users.FS_USERS_COLL).
		Doc(userID.String()).
		Collection(FS_ELO_COLL)
}

func (s *service) getCurrentEloRef(userID format.UserID, game GameType) *firestore.DocumentRef {
	return s.getElosRef(userID).
		Doc(game.String()).
		Collection(FS_GAME_ELOS_COLL).
		Doc(FS_CURRENT_ELO_DOC)
}

func (s *service) getTimestampEloRef(userID format.UserID, game GameType, timestamp time.Time) *firestore.DocumentRef {
	return s.getElosRef(userID).
		Doc(game.String()).
		Collection(FS_GAME_ELOS_COLL).
		Doc(timestamp.String())
}
