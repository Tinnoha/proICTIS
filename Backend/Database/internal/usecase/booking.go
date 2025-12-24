package usecase

import (
	"database/internal/entity"
	"errors"
	"time"

	"github.com/gofrs/uuid"
)

type BookingRepository interface {
	GetAllBooks() ([]entity.Booking, error)                                  // +
	GetBooksByUserId(userId uuid.UUID) ([]entity.Booking, error)             // +
	GetBooksByEnviromentId(EnviromentId uuid.UUID) ([]entity.Booking, error) // +

	Book(UserId uuid.UUID, EnviromentId uuid.UUID, Start time.Time, End time.Time) (entity.Booking, error) // +

	IsTaken(EnviromentId uuid.UUID, start time.Time, end time.Time) (bool, error)

	AcceptBooking(BookingId uuid.UUID) (entity.Booking, error)
	CancelBooking(BookingId uuid.UUID) error
	ReturnEnviroment(BookingId uuid.UUID) error
	DeleteBooking(BookingId uuid.UUID) error

	EditStatusBooking(BookingId uuid.UUID, status string) (entity.Booking, error)

	IsExists(BoookingId uuid.UUID) bool
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
	exists := uc.UserRepo.IsExistsById(id)

	if !exists {
		return []entity.Booking{}, ErrNotFound
	}

	books, err := uc.BookingRepo.GetBooksByUserId(id)

	if err != nil {
		return []entity.Booking{}, ErrInntenal(err)
	}

	return books, err
}

func (uc *BookingUseCase) GetBooksByEnviromentId(EnviromentId uuid.UUID) ([]entity.Booking, error) {
	books, err := uc.BookingRepo.GetBooksByEnviromentId(EnviromentId)

	if err != nil {
		return []entity.Booking{}, ErrInntenal(err)
	}

	return books, err
}

func (uc *BookingUseCase) Book(UserId uuid.UUID, EnviromentId uuid.UUID, start time.Time, end time.Time) (entity.Booking, error) {
	exists := uc.EnviromentRepo.IsExists(EnviromentId)

	if !exists {
		return entity.Booking{}, ErrNotFound
	}

	exists = uc.UserRepo.IsExistsById(UserId)

	if !exists {
		return entity.Booking{}, ErrNotFound
	}

	taken, err := uc.IsTaken(EnviromentId, start, end)

	if err != nil {
		return entity.Booking{}, ErrInntenal(err)
	}

	if taken {
		return entity.Booking{}, ErrThisExists("booking", EnviromentId.String())
	}

	booking, err := uc.BookingRepo.Book(UserId, EnviromentId, start, end)

	if err != nil {
		return entity.Booking{}, ErrInntenal(err)
	}

	if booking.Status != "Wait" {
		return entity.Booking{}, ErrInntenal(errors.New("No good value from Book"))
	}

	return booking, nil
}

func (uc *BookingUseCase) IsTaken(EnviromentId uuid.UUID, start time.Time, end time.Time) (bool, error) {
	exists := uc.EnviromentRepo.IsExists(EnviromentId)

	if !exists {
		return false, ErrNotFound
	}

	Timer := start.Compare(end)

	if Timer != -1 {
		return false, ErrWrongData
	}

	return uc.BookingRepo.IsTaken(EnviromentId, start, end)
}

func (uc *BookingUseCase) AcceptBooking(BookingId uuid.UUID) (entity.Booking, error) {
	exists := uc.BookingRepo.IsExists(BookingId)

	if !exists {
		return entity.Booking{}, ErrNotFound
	}

	arenda, err := uc.BookingRepo.EditStatusBooking(BookingId, "Wait Time for Booking")

	if err != nil {
		return entity.Booking{}, ErrInntenal(err)
	}

	return arenda, nil
}

func (uc *BookingUseCase) CancelBooking(BookingId uuid.UUID) error {
	exists := uc.BookingRepo.IsExists(BookingId)

	if !exists {
		return ErrNotFound
	}

	err := uc.BookingRepo.DeleteBooking(BookingId)

	if err != nil {
		return ErrInntenal(err)
	}

	return nil
}

func (uc *BookingUseCase) ReturnEnviroment(BookingId uuid.UUID) error {
	exists := uc.BookingRepo.IsExists(BookingId)

	if !exists {
		return ErrNotFound
	}

	err := uc.DeleteBooking(BookingId)

	return err
}

func (uc *BookingUseCase) DeleteBooking(BookingId uuid.UUID) error {
	exists := uc.BookingRepo.IsExists(BookingId)

	if !exists {
		return ErrNotFound
	}

	err := uc.BookingRepo.DeleteBooking(BookingId)

	if err != nil {
		return ErrInntenal(err)
	}

	return nil
}

func (uc *BookingUseCase) IsExists(BoookingId uuid.UUID) bool {
	return uc.BookingRepo.IsExists(BoookingId)
}
