package users

import "github.com/garlicgarrison/chessvars-backend/pkg/format"

type Gender string

const (
	GENDER_MALE      Gender = "MALE"
	GENDER_FEMALE    Gender = "FEMALE"
	GENDER_NONBINARY Gender = "NONBINARY"
	GENDER_UNKNOWN   Gender = "UNKNOWN"
)

type UserDocument struct {
	UserID   format.UserID `firestore:"user_id"`
	Username string        `firestore:"username"`
	Gender   Gender        `firestore:"gender"`
	Elo      int           `firestore:"elo"`
}
