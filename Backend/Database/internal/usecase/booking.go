package usecase

import (
	"database/internal/entity"
	"time"
)

type BookingRepository interface {
	GetBooksByUserId(id int) ([]entity.Booking, error)
	GetBooksByEnviromentId(TypeOfEnviroment string) ([]entity.Booking, error)
	Book(UserId int, EnviromentId int, Start time.Time, End time.Time) (entity.Booking, error)
	IsTaken(EnviromentId int, start time.Time, end time.Time) bool
	ReturnEnviroment(EnviromentId int) (entity.Booking, error)
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

func (uc *BookingUseCase) GetBooksByUserId(id int) ([]entity.Booking, error) {
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

func (uc *BookingUseCase) GetBooksByEnviromentName(TypeOfEnviroment string) ([]entity.Booking, error) {
	exists := uc.EnviromentRepo.TypeIsExsits(TypeOfEnviroment)

	if !exists {
		return []entity.Booking{}, ErrNotFound
	}

	books, err := uc.BookingRepo.GetBooksByEnviromentId(TypeOfEnviroment)

	if err != nil {
		return []entity.Booking{}, ErrInntenal(err)
	}

	return books, err
}

func (uc *BookingUseCase) Book(UserId int, EnviromentId int, start time.Time, end time.Time) (entity.Booking, error) {
	exists := uc.EnviromentRepo.IsExists(EnviromentId)

	if !exists {
		return entity.Booking{}, ErrNotFound
	}

	exists = uc.UserRepo.IsExistsById(UserId)

	if !exists {
		return entity.Booking{}, ErrNotFound
	}

	taken := uc.BookingRepo.IsTaken(EnviromentId, start, end)

	if taken {
		return entity.Booking{}, ErrThisExists("booking", string(EnviromentId))
	}

	booking, err := uc.BookingRepo.Book(UserId, EnviromentId, start, end)

	if err != nil {
		return entity.Booking{}, ErrInntenal(err)
	}

	return booking, nil
}

func (uc *BookingUseCase) ReturnEnviroment(EnviromentId int) (entity.Booking, error) {

}

func (uc *BookingUseCase) IsTaken(EnviromentId int, start time.Time, end time.Time) bool {

}
