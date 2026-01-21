package repository

import (
	"database/internal/entity"
	"database/internal/usecase"
	"database/sql"
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

func (b *bookingRepo) GetBooksByUserId(userId uuid.UUID) ([]entity.Booking, error) {
	rows, err := b.db.Query(`
	SELECT id, user_id, equipment_id, book_start, book_end, status FROM proICTIS_booking 
	WHERE user_id = $1`, userId)
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

func (b *bookingRepo) GetBooksByEquipmentId(EquipmentId uuid.UUID) ([]entity.Booking, error) {
	rows, err := b.db.Query(`
	SELECT id, user_id, equipment_id, book_start, book_end, status FROM proICTIS_booking 
	WHERE equipment_id = $1`, EquipmentId)
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

func (b *bookingRepo) Book(UserId uuid.UUID, EquipmentId uuid.UUID, Start time.Time, End time.Time) (entity.Booking, error) {
	var using bool

	err := b.db.QueryRow(`	
	SELECT EXISTS(
	SELECT 1 FROM proICTIS_booking
	WHERE equipment_id = $1 and
	book_start < $2 and
	book_end > $3
	)`, EquipmentId, End, Start).Scan(&using)

	if err != nil {
		return entity.Booking{}, err
	}

	if using {
		return entity.Booking{}, usecase.ErrThisExists("booking", "time")
	}

	id, err := uuid.NewV4()

	if err != nil {
		return entity.Booking{}, err
	}

	book := entity.Booking{}

	err = b.db.QueryRow(`INSERT INTO proICTIS_booking 
	(id, user_id,equipment_id,book_start,book_end,status )
	VALUES ($1,$2,$3,$4,$5,$6)
	RETURNING id, user_id,equipment_id,book_start,book_end,status`,
		id, UserId, EquipmentId, Start, End, "Waiting answer",
	).Scan(
		&book.ID,
		&book.UserId,
		&book.EquipmentId,
		&book.BookStart,
		&book.BookEnd,
		&book.Status,
	)

	if err != nil {
		return entity.Booking{}, err
	}

	return book, nil
}

func (b *bookingRepo) EditStatusBooking(BookingId uuid.UUID, status string) (entity.Booking, error) {
	book := entity.Booking{}

	err := b.db.QueryRow(`
		UPDATE proICTIS_booking 
		SET status = $1 
		WHERE id = $2
		RETURNING id, user_id,equipment_id,book_start,book_end,status`,
		status, BookingId).Scan(
		&book.ID,
		&book.UserId,
		&book.EquipmentId,
		&book.BookStart,
		&book.BookEnd,
		&book.Status,
	)

	if err != nil {
		return entity.Booking{}, err
	}

	return book, nil
}

func (b *bookingRepo) DeleteBooking(BookingId uuid.UUID) error {
	rez, err := b.db.Exec(`DELETE from proICTIS_booking where id = $1`, BookingId)
	if err != nil {
		return err
	}

	c, err := rez.RowsAffected()

	if err != nil {
		return err
	}

	if c == 0 {
		return sql.ErrNoRows
	}

	return nil
}
