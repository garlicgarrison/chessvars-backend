package game

import (
	"github.com/garlicgarrison/chessvars-backend/pkg/firestore"
	"github.com/garlicgarrison/chessvars-backend/pkg/format"
)

const (
	FS_GAMES_COLL = "games"
)

func (s *service) getGamesRef() *firestore.CollectionRef {
	return s.fs.Collection(FS_GAMES_COLL)
}

func (s *service) getGameRef(gameID format.GameID) *firestore.DocumentRef {
	return s.getGamesRef().Doc(gameID.String())
}
