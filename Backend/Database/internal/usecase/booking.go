package usecase

import (
	"database/internal/entity"
	"database/sql"
	"errors"
	"time"

	"github.com/gofrs/uuid"
)

type BookingRepository interface {
	GetAllBooks() ([]entity.Booking, error)                                  // +
	GetBooksByUserId(userId uuid.UUID) ([]entity.Booking, error)             // +
	GetBooksByEnviromentId(EnviromentId uuid.UUID) ([]entity.Booking, error) // +

	Book(UserId uuid.UUID, EnviromentId uuid.UUID, Start time.Time, End time.Time) (entity.Booking, error) // +

	AcceptBooking(BookingId uuid.UUID) (entity.Booking, error)

	EditStatusBooking(BookingId uuid.UUID, status string) (entity.Booking, error)
	DeleteBooking(BookingId uuid.UUID) error
}

type BookingUseCase struct {
	UserRepo       UserRepository
	BookingRepo    BookingRepository
	EnviromentRepo EnviromentRepository
}

func NewBooknigUseCase(
	userRepo UserRepository,
	bookRepo BookingRepository,
	enviromentRepo EnviromentRepository,
) *BookingUseCase {
	return &BookingUseCase{
		UserRepo:       userRepo,
		BookingRepo:    bookRepo,
		EnviromentRepo: enviromentRepo,
	}
}

func (uc *BookingUseCase) GetAllBooks() ([]entity.Booking, error) {
	books, err := uc.BookingRepo.GetAllBooks()

	if err != nil {
		return []entity.Booking{}, ErrInntenal(err)
	}

	return books, nil
}

func (uc *BookingUseCase) GetBooksByUserId(id uuid.UUID) ([]entity.Booking, error) {
	books, err := uc.BookingRepo.GetBooksByUserId(id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []entity.Booking{}, ErrNotFound
		}
		return []entity.Booking{}, ErrInntenal(err)
	}

	return books, err
}

func (uc *BookingUseCase) GetBooksByEnviromentId(EnviromentId uuid.UUID) ([]entity.Booking, error) {
	books, err := uc.BookingRepo.GetBooksByEnviromentId(EnviromentId)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []entity.Booking{}, ErrNotFound
		}
		return []entity.Booking{}, ErrInntenal(err)
	}

	return books, err
}

func (uc *BookingUseCase) Book(UserId uuid.UUID, EnviromentId uuid.UUID, start time.Time, end time.Time) (entity.Booking, error) {
	booking, err := uc.BookingRepo.Book(UserId, EnviromentId, start, end)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Booking{}, ErrNotFound
		} else if errors.Is(err, ErrThisExist) {
			return entity.Booking{}, ErrThisExists("enviroment", EnviromentId.String())
		}
		return entity.Booking{}, ErrInntenal(err)
	}

	if booking.Status != "Wait" {
		return entity.Booking{}, ErrInntenal(errors.New("No good value from Book"))
	}

	return booking, nil
}

func (uc *BookingUseCase) AcceptBooking(BookingId uuid.UUID) (entity.Booking, error) {
	arenda, err := uc.BookingRepo.EditStatusBooking(BookingId, "Wait Time for Booking")

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Booking{}, ErrNotFound
		}
		return entity.Booking{}, ErrInntenal(err)
	}

	return arenda, nil
}

func (uc *BookingUseCase) DeleteBooking(BookingId uuid.UUID) error {
	err := uc.BookingRepo.DeleteBooking(BookingId)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
		return ErrInntenal(err)
	}

	return nil
}
