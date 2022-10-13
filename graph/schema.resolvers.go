package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/garlicgarrison/chessvars-backend/graph/generated"
	"github.com/garlicgarrison/chessvars-backend/graph/model"
	"github.com/garlicgarrison/chessvars-backend/graph/resolver"
	"github.com/garlicgarrison/chessvars-backend/pkg/format"
	"github.com/garlicgarrison/chessvars-backend/pkg/game"
	"google.golang.org/grpc/codes"
)

// Type is the resolver for the type field.
func (r *gameResolver) Type(ctx context.Context, obj *resolver.Game) (*model.GameType, error) {
	panic(fmt.Errorf("not implemented"))
}

// UserEdit is the resolver for the userEdit field.
func (r *mutationResolver) UserEdit(ctx context.Context, input model.UserEditInput) (*model.UserMutationResponse, error) {
	panic(fmt.Errorf("not implemented"))
}

// UserDelete is the resolver for the userDelete field.
func (r *mutationResolver) UserDelete(ctx context.Context) (*model.BasicMutationResponse, error) {
	panic(fmt.Errorf("not implemented"))
}

// GameCreate is the resolver for the gameCreate field.
func (r *mutationResolver) GameCreate(ctx context.Context, typeArg model.GameType, limit resolver.TimeLimit) (*model.GameMutationResponse, error) {
	userID, ok := resolver.GetAuthUserID(ctx)
	if !ok {
		return nil, fmt.Errorf("could not validate user")
	}

	var timeLimit game.TimeLimit
	switch limit {
	case resolver.BULLET:
		timeLimit = game.BULLET
	case resolver.BLITZ:
		timeLimit = game.BLITZ
	case resolver.BLITZ2:
		timeLimit = game.BLITZ2
	case resolver.RAPID:
		timeLimit = game.RAPID
	case resolver.RAPID2:
		timeLimit = game.RAPID2
	case resolver.RAPID3:
		timeLimit = game.RAPID3
	case resolver.RAPID4:
		timeLimit = game.RAPID4
	default:
		return nil, fmt.Errorf("time limit not valid")
	}

	var gameType game.GameType
	switch typeArg {
	case model.GameTypeJanggi:
		gameType = game.JANGGI
	default:
		return nil, fmt.Errorf("game not implemented")
	}

	game, err := r.Services.Game.CreateGame(ctx, game.CreateGameRequest{
		UserID:    userID,
		TimeLimit: timeLimit,
		Type:      gameType,
	})
	if err != nil {
		return nil, err
	}

	return &model.GameMutationResponse{
		Code:    http.StatusOK,
		Success: true,
		Message: "game was successfully created",
		Game:    resolver.NewGameWithData(r.Services, game),
	}, nil
}

// GameJoin is the resolver for the gameJoin field.
func (r *mutationResolver) GameJoin(ctx context.Context, id string) (*model.GameMutationResponse, error) {
	userID, ok := resolver.GetAuthUserID(ctx)
	if !ok {
		return nil, fmt.Errorf("could not validate user")
	}

	gameID, err := format.ParseGameID(id)
	if err != nil {
		return nil, err
	}

	game, err := r.Services.Game.JoinGame(ctx, game.JoinGameRequest{
		GameID: gameID,
		UserID: userID,
	})
	if err != nil {
		return nil, err
	}

	return &model.GameMutationResponse{
		Code:    int(codes.OK),
		Success: true,
		Message: "game successfully joined",
		Game:    resolver.NewGameWithData(r.Services, game),
	}, err
}

// GameMove is the resolver for the gameMove field.
func (r *mutationResolver) GameMove(ctx context.Context, id string, move string, status *model.GameStatus) (*model.GameMutationResponse, error) {
	userID, ok := resolver.GetAuthUserID(ctx)
	if !ok {
		return nil, fmt.Errorf("could not validate user")
	}

	gameID, err := format.ParseGameID(id)
	if err != nil {
		return nil, err
	}

	moveN, err := game.ParseMoveNotation(move)
	if err != nil {
		return nil, err
	}

	gameStatus := game.INGAME
	if status != nil {
		switch *status {
		case model.GameStatusWin:
			gameStatus = game.WIN
		case model.GameStatusIngame:
			gameStatus = game.INGAME
		case model.GameStatusLoss:
			gameStatus = game.LOSS
		case model.GameStatusDraw:
			gameStatus = game.DRAW
		}
	}

	gameReply, err := r.Services.Game.EditGame(ctx, game.EditGameRequest{
		UserID: userID,
		GameID: gameID,
		Status: gameStatus,
		Move:   &moveN,
	})
	if err != nil {
		return nil, err
	}

	// send move to all channels with given gameID
	moveObservers := r.getObserverMap(gameID)
	log.Printf("moveObsercers %v", moveObservers)
	moveObservers.MoveObservers.Range(func(_, value interface{}) bool {
		observer := value.(*MoveObserver)
		log.Printf("[gameMove] -- move: %v, userID: %s", observer.Move, observer.UserID.String())
		if observer.UserID != userID {
			observer.Move <- resolver.NewMove(r.Services, &gameReply.Moves[len(gameReply.Moves)-1])
		}

		return true
	})

	return &model.GameMutationResponse{
		Code:    http.StatusOK,
		Success: true,
		Message: "move was successfully added",
		Game:    resolver.NewGameWithData(r.Services, gameReply),
	}, nil
}

// GameAbort is the resolver for the gameAbort field.
func (r *mutationResolver) GameAbort(ctx context.Context, id string) (*model.GameMutationResponse, error) {
	userID, ok := resolver.GetAuthUserID(ctx)
	if !ok {
		return nil, fmt.Errorf("could not validate user")
	}

	gameID, err := format.ParseGameID(id)
	if err != nil {
		return nil, err
	}

	reply, err := r.Services.Game.EditGame(ctx, game.EditGameRequest{
		UserID: userID,
		GameID: gameID,
		Status: game.Aborted,
	})
	if err != nil {
		return &model.GameMutationResponse{
			Code:    int(codes.Internal),
			Success: false,
			Message: "could not abort game",
		}, nil
	}

	return &model.GameMutationResponse{
		Code:    http.StatusOK,
		Success: true,
		Message: "game aborted",
		Game:    resolver.NewGameWithData(r.Services, reply),
	}, nil
}

// User is the resolver for the user field.
func (r *queryResolver) User(ctx context.Context, id *string) (*resolver.User, error) {
	if id != nil {
		userID, err := format.ParseUserID(*id)
		if err != nil {
			return nil, err
		}

		return resolver.NewUser(r.Services, userID), nil
	}

	userID, ok := resolver.GetAuthUserID(ctx)
	if !ok {
		return nil, errors.New("no auth user in context")
	}

	return resolver.NewUser(r.Services, userID), nil
}

// Game is the resolver for the game field.
func (r *queryResolver) Game(ctx context.Context, id string) (*resolver.Game, error) {
	gameID, err := format.ParseGameID(id)
	if err != nil {
		return nil, err
	}

	return resolver.NewGame(r.Services, gameID), nil
}

// OnMoveNew is the resolver for the onMoveNew field.
func (r *subscriptionResolver) OnMoveNew(ctx context.Context, id string) (<-chan *resolver.Move, error) {
	userID, ok := resolver.GetAuthUserID(ctx)
	if !ok {
		return nil, fmt.Errorf("could not parse user from context")
	}

	gameID, err := format.ParseGameID(id)
	if err != nil {
		return nil, err
	}

	mc := make(chan *resolver.Move, 1)
	observers := r.getObserverMap(gameID)
	observers.MoveObservers.Store(userID, &MoveObserver{
		UserID: userID,
		Move:   mc,
	})

	go func() {
		<-ctx.Done()
		// delete observers
		observers.MoveObservers.Delete(userID)
	}()

	log.Printf("[OnMoveNew] %v", mc)
	return mc, nil
}

// Game returns generated.GameResolver implementation.
func (r *Resolver) Game() generated.GameResolver { return &gameResolver{r} }

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Subscription returns generated.SubscriptionResolver implementation.
func (r *Resolver) Subscription() generated.SubscriptionResolver { return &subscriptionResolver{r} }

type gameResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
