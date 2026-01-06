package repository

import (
	"database/internal/entity"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
)

/*
GetAllBooks() ([]entity.Booking, error)
GetBooksByUserId(userId uuid.UUID) ([]entity.Booking, error)
GetBooksByEnviromentId(EnviromentId uuid.UUID) ([]entity.Booking, error)

Book(UserId uuid.UUID, EnviromentId uuid.UUID, Start time.Time, End time.Time) (entity.Booking, error)

AcceptBooking(BookingId uuid.UUID) (entity.Booking, error)

EditStatusBooking(BookingId uuid.UUID, status string) (entity.Booking, error)
DeleteBooking(BookingId uuid.UUID) error
*/

type bookingRepo struct {
	db sqlx.DB
}

func NewBookingRepo(db sqlx.DB) *bookingRepo {
	return &bookingRepo{
		db: db,
	}
}

func (b *bookingRepo) GetAllBooks() ([]entity.Booking, error)
func (b *bookingRepo) GetBooksByUserId(userId uuid.UUID) ([]entity.Booking, error)
func (b *bookingRepo) GetBooksByEnviromentId(EnviromentId uuid.UUID) ([]entity.Booking, error)
func (b *bookingRepo) Book(UserId uuid.UUID, EnviromentId uuid.UUID, Start time.Time, End time.Time) (entity.Booking, error)
func (b *bookingRepo) AcceptBooking(BookingId uuid.UUID) (entity.Booking, error)
func (b *bookingRepo) EditStatusBooking(BookingId uuid.UUID, status string) (entity.Booking, error)
func (b *bookingRepo) DeleteBooking(BookingId uuid.UUID) error
