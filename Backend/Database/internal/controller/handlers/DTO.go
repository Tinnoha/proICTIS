package handlers

import (
	"database/internal/entity"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gofrs/uuid"
)

type AdminDTO struct {
	AdminId uuid.UUID `json:"admin_id"`
}

type TypesDTO struct {
	AdminId uuid.UUID                `json:"admin_id"`
	Types   []entity.TypeOfEquipment `json:"types"`
}

type TypeDTO struct {
	AdminId uuid.UUID              `json:"admin_id"`
	Type    entity.TypeOfEquipment `json:"type"`
}

type EquipmentEdit struct {
	AdminId   uuid.UUID        `json:"admin_id"`
	Equipment entity.Equipment `json:"Equipment"`
}

type EmailDTO struct {
	AdminId uuid.UUID `json:"admin_id"`
	Email   string    `json:"email"`
}

type AdminsDTO struct {
	FirstAdmin  uuid.UUID `json:"first_admin"`
	SecondAdmin uuid.UUID `json:"second_admin"`
}

type BoolDTO struct {
	AdminId uuid.UUID `json:"admin_id"`
	Active  bool      `json:"active"`
}

type BookDTO struct {
	UserId     uuid.UUID `json:"user_id"`
	EnviromtId uuid.UUID `json:"enviromt_id"`
	Start      time.Time `json:"start"`
	End        time.Time `json:"end"`
}

type StatusDTO struct {
	AdminId uuid.UUID `json:"admin_id"`
	Status  string    `json:"status"`
}

type YandexUser struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"default_email"`
	Avatar    string `json:"default_avatar_id"`
}

type errDTO struct {
	Err  string    `json:"Error"`
	Time time.Time `json:"Time"`
}

type UploadDTO struct {
	URL string `json:"url"`
}

func HttpError(w http.ResponseWriter, err error, status int) {
	errdto := errDTO{
		Err:  err.Error(),
		Time: time.Now(),
	}

	b, err := json.MarshalIndent(errdto, "", "    ")

	if err != nil {
		panic("LOl" + err.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(b)
}
