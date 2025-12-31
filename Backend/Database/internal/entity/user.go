package entity

import (
	"github.com/gofrs/uuid"
)

type User struct {
	Id            uuid.UUID `json:"id"`
	FirstName     string    `json:"first_name"`
	SecondName    string    `json:"second_name"`
	Email         string    `json:"email"`
	AvatarURL     string    `json:"avatar_url"`
	Role          string    `json:"role"`
	TokenProvider int       `json:"tokenProvider"`
}
