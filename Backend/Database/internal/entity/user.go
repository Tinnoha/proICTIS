package entity

import "time"

type User struct {
	FirstName     string `json:"first_name"`
	SecondName    string `json:"second_name"`
	Email         string `json:"email"`
	AvatarURL     string `json:"avatar_url"`
	Role          string `json:"role"`
	TokenProvider int    `json:"tokenProvider"`

	CreatedAt time.Time `json:"created_at"`
}
