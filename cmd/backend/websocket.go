package main

import (
	"context"
	"fmt"

	"firebase.google.com/go/v4/auth"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/garlicgarrison/chessvars-backend/middleware"
	"github.com/garlicgarrison/chessvars-backend/pkg/format"
)

func initWebsocket(ctx context.Context, client *auth.Client, payload transport.InitPayload) (context.Context, error) {
	id := payload.Authorization()

	token, err := client.VerifyIDToken(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("[initWebsocket] -- could not verify token")
	}

	userID := format.NewUserIDFromIdentifer(token.UID)
	ctxNew := context.WithValue(ctx, middleware.AUTH_USER_CONTEXT_KEY, userID)

	return ctxNew, nil
}
