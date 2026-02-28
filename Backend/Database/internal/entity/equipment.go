package entity

import "github.com/gofrs/uuid"

type Equipment struct {
	Id              uuid.UUID `json:"id"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	PhotoURL        string    `json:"photo_url"`
	TypeOfEquipment string    `json:"type"`
	Auditory        string    `json:"auditory"`
	IsActive        bool      `json:"is_active"`
}

type Equipments struct {
	AdminId uuid.UUID   `json:"admin_id"`
	Tovars  []Equipment `json:"tovars"`
}

type TypeOfEquipment struct {
	Id   uuid.UUID
	Name string `json:"name"`
}
