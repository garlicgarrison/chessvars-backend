package users

import "github.com/garlicgarrison/chessvars-backend/pkg/format"

type UserDocument struct {
	UserID   format.UserID `firestore:"user_id"`
	Username string        `firestore:"username"`
}
