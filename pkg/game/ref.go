package game

import (
	"github.com/garlicgarrison/chessvars-backend/pkg/firestore"
	"github.com/garlicgarrison/chessvars-backend/pkg/format"
)

func (s *service) getGamesRef(game GameType) *firestore.CollectionRef {
	return s.fs.Collection(game.String())
}

func (s *service) getGameRef(game GameType, gameID format.GameID) *firestore.DocumentRef {
	return s.getGamesRef(game).Doc(gameID.String())
}
