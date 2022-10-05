package users

import (
	"time"

	"github.com/garlicgarrison/chessvars-backend/pkg/format"
)

type UserDocument struct {
	UserID    format.UserID `firestore:"user_id"`
	Email     string        `firestore:"email"`
	Username  string        `firestore:"username"`
	CreatedAt time.Time     `firestore:"created_at"`
}
