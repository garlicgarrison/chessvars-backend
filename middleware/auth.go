package middleware

import (
	"context"
	"log"
	"net/http"

	"firebase.google.com/go/v4/auth"
	"github.com/garlicgarrison/chessvars-backend/pkg/format"
)

type ContextKey string

const (
	AUTH_USER_CONTEXT_KEY       ContextKey = "AUTH_USER"
	AUTH_USER_EMAIL_CONTEXT_KEY ContextKey = "AUTH_USER_EMAIL"
)

type Auth struct {
	client *auth.Client
	next   http.Handler
}

func NewAuth(client *auth.Client, next http.Handler) *Auth {
	return &Auth{client: client, next: next}
}

func (a *Auth) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// w.Write([]byte(fmt.Sprintf("headers %v", r.Header)))
	id := r.Header.Get("Authorization")
	if id == "" {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("missing Authorization header -- must have a valid id token"))
		return
	}

	token, err := a.client.VerifyIDToken(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error in verifying token"))
		log.Printf("error in verifying token: %s\n", err)
		return
	}

	userID := format.NewUserIDFromIdentifer(token.UID)
	email := token.Firebase.Identities["email"].([]interface{})[0].(string)
	ctx := context.WithValue(r.Context(), AUTH_USER_CONTEXT_KEY, userID)
	ctx = context.WithValue(ctx, AUTH_USER_EMAIL_CONTEXT_KEY, email)
	request := r.WithContext(ctx)
	a.next.ServeHTTP(w, request)
}
