package entity

import (
	"time"

	"github.com/gofrs/uuid"
)

type Booking struct {
	ID           uuid.UUID
	UserId       int
	EnviromentId int
	BookStart    time.Time
	BookEnd      time.Time
	Status       string
	mes          string
}
