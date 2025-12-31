package entity

import "github.com/gofrs/uuid"

type Enviroment struct {
	Id               uuid.UUID `json:"id"`
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	PhotoURL         string    `json:"photo_url"`
	TypeOfEnviroment string    `json:"type"`
	Auditory         string    `json:"auditory"`
	IsActive         bool      `json:"is_active"`
}

type TypeOfEnviroment struct {
	Id   uuid.UUID
	Name string
}
