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
GetBooksByEquipmentId(EquipmentId uuid.UUID) ([]entity.Booking, error)

Book(UserId uuid.UUID, EquipmentId uuid.UUID, Start time.Time, End time.Time) (entity.Booking, error)

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

func (b *bookingRepo) GetAllBooks() ([]entity.Booking, error) {
	rows, err := b.db.Query(`SELECT id, user_id, equipment_id, book_start, book_end, status FROM proICTIS_booking`)

	if err != nil {
		return []entity.Booking{}, err
	}

	defer rows.Close()

	bookings := []entity.Booking{}

	for rows.Next() {
		book := entity.Booking{}

		err := rows.Scan(
			&book.ID,
			&book.UserId,
			&book.EquipmentId,
			&book.BookStart,
			&book.BookEnd,
			&book.Status,
		)

		if err != nil {
			return []entity.Booking{}, err
		}

		bookings = append(bookings, book)
	}

	if err := rows.Err(); err != nil {
		return []entity.Booking{}, err
	}

	return bookings, nil
}

func (b *bookingRepo) GetBooksByUserId(userId uuid.UUID) ([]entity.Booking, error)
func (b *bookingRepo) GetBooksByEquipmentId(EquipmentId uuid.UUID) ([]entity.Booking, error)
func (b *bookingRepo) Book(UserId uuid.UUID, EquipmentId uuid.UUID, Start time.Time, End time.Time) (entity.Booking, error)
func (b *bookingRepo) AcceptBooking(BookingId uuid.UUID) (entity.Booking, error)
func (b *bookingRepo) EditStatusBooking(BookingId uuid.UUID, status string) (entity.Booking, error)
func (b *bookingRepo) DeleteBooking(BookingId uuid.UUID) error
