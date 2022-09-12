package users

import "github.com/pafkiuq/backend/pkg/format"

type Preference string

const (
	PREF_STRAIGHT Preference = "STRAIGHT"
	PREF_GAY      Preference = "GAY"
	PREF_TRANS    Preference = "TRANS"
)

type Gender string

const (
	GENDER_MALE      Gender = "MALE"
	GENDER_FEMALE    Gender = "FEMALE"
	GENDER_NONBINARY Gender = "NONBINARY"
	GENDER_UNKNOWN   Gender = "UNKNOWN"
)

type UserDocument struct {
	UserID      format.UserID `firestore:"user_id"`
	Username    string        `firestore:"username"`
	Bio         string        `firestore:"bio"`
	Preferences []Preference  `firestore:"preferences"`
	Gender      Gender        `firestore:"gender"`
}
