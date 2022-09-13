// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"fmt"
	"io"
	"strconv"
)

type MutationResponse interface {
	IsMutationResponse()
}

type BasicMutationResponse struct {
	Code    int    `json:"code"`
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func (BasicMutationResponse) IsMutationResponse() {}

type Game struct {
	ID        string   `json:"id"`
	Moves     []string `json:"moves"`
	PlayerOne *User    `json:"playerOne"`
	PlayerTwo *User    `json:"playerTwo"`
	Winner    *User    `json:"winner"`
	Draw      *bool    `json:"draw"`
	Aborted   *bool    `json:"aborted"`
	Timestamp *string  `json:"timestamp"`
}

type GameMutationResponse struct {
	Code    int    `json:"code"`
	Success bool   `json:"success"`
	Message string `json:"message"`
	Game    *Game  `json:"game"`
}

func (GameMutationResponse) IsMutationResponse() {}

type Pagination struct {
	Cursor *string `json:"cursor"`
	Limit  *int    `json:"limit"`
}

type User struct {
	ID       string  `json:"id"`
	Exists   *bool   `json:"exists"`
	Username *string `json:"username"`
	Elo      *int    `json:"elo"`
}

type UserEditInput struct {
	Username *string `json:"username"`
	Bio      *string `json:"bio"`
}

type UserMutationResponse struct {
	Code    int    `json:"code"`
	Success bool   `json:"success"`
	Message string `json:"message"`
	User    *User  `json:"user"`
}

func (UserMutationResponse) IsMutationResponse() {}

type Users struct {
	Users []*User `json:"users"`
	Next  *string `json:"next"`
}

type GameType string

const (
	GameTypeJanggi GameType = "JANGGI"
	GameTypeShogi  GameType = "SHOGI"
)

var AllGameType = []GameType{
	GameTypeJanggi,
	GameTypeShogi,
}

func (e GameType) IsValid() bool {
	switch e {
	case GameTypeJanggi, GameTypeShogi:
		return true
	}
	return false
}

func (e GameType) String() string {
	return string(e)
}

func (e *GameType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = GameType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid GameType", str)
	}
	return nil
}

func (e GameType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
