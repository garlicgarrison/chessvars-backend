package main

import (
	"context"
	"fmt"
	"log"

	"firebase.google.com/go/v4/auth"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/garlicgarrison/chessvars-backend/middleware"
	"github.com/garlicgarrison/chessvars-backend/pkg/format"
)

// type ContextKey string

// const (
// 	AUTH_USER_CONTEXT_KEY       ContextKey = "AUTH_USER"
// 	AUTH_USER_EMAIL_CONTEXT_KEY ContextKey = "AUTH_USER_EMAIL"
// )

func initWebsocket(ctx context.Context, client *auth.Client, payload transport.InitPayload) (context.Context, error) {
	id := payload.Authorization()

	token, err := client.VerifyIDToken(ctx, id)
	log.Printf("token and err %v %s", token, err)
	if err != nil {
		return nil, fmt.Errorf("[initWebsocket] -- could not verify token")
	}

	userID := format.NewUserIDFromIdentifer(token.UID)
	log.Printf("userID %s", userID)
	ctxNew := context.WithValue(ctx, middleware.AUTH_USER_CONTEXT_KEY, userID)

	return ctxNew, nil
}
