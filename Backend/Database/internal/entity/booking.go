package entity

import (
	"time"

	"github.com/gofrs/uuid"
)

type Booking struct {
	ID          uuid.UUID
	UserId      uuid.UUID
	EquipmentId uuid.UUID
	BookStart   time.Time
	BookEnd     time.Time
	Status      string
}
