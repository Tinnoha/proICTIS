package usecase

import (
	"database/internal/entity"
	"database/sql"
	"errors"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

type UserRepository interface {
	GetAll() ([]entity.User, error)
	GetById(id uuid.UUID) (entity.User, error)
	GetByEmail(email string) (entity.User, error)
	CreateUser(entity.User) (entity.User, error)
	MakeAdmin(id uuid.UUID) (entity.User, error)
	MakeSuperAdmin(id uuid.UUID) (entity.User, error)
}

type UserUseCase struct {
	UserRepo UserRepository
}

func NewUserUseCase(
	userRepo UserRepository,
) *UserUseCase {
	return &UserUseCase{
		UserRepo: userRepo,
	}
}

func (uc *UserUseCase) GetAll() ([]entity.User, error) {
	vasya, err := uc.UserRepo.GetAll()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []entity.User{}, ErrNotFound
		} else {
			return []entity.User{}, ErrInntenal(err)
		}
	}

	return vasya, nil
}

func (uc *UserUseCase) GetById(id uuid.UUID) (entity.User, error) {
	vasya, err := uc.UserRepo.GetById(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.User{}, ErrNotFound
		} else {
			return entity.User{}, ErrInntenal(err)
		}
	}

	return vasya, nil
}

func (uc *UserUseCase) GetByEmail(email string) (entity.User, error) {
	vasya, err := uc.UserRepo.GetByEmail(email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.User{}, ErrNotFound
		}
		return entity.User{}, ErrInntenal(err)
	}

	return vasya, nil
}

func (uc *UserUseCase) CreateUser(vasy []entity.User) ([]entity.User, error) {
	rezult := []entity.User{}

	for _, annya := range vasy {
		Id, err := uuid.NewV4()
		annya.Id = Id
		if err != nil {
			return []entity.User{}, ErrInntenal(err)
		}

		annya.Role = "student"

		vasya, err := uc.UserRepo.CreateUser(annya)

		if err != nil {
			var pgerr *pgconn.PgError
			if errors.As(err, &pgerr) && pgerr.Code == "23505" {
				return []entity.User{}, ErrThisExists("email", annya.Email)
			} else {
				return []entity.User{}, ErrInntenal(err)
			}
		}

		rezult = append(rezult, vasya)
	}

	return rezult, nil
}

func (uc *UserUseCase) MakeAdmin(id uuid.UUID) (entity.User, error) {

	vasya, err := uc.UserRepo.MakeAdmin(id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.User{}, ErrNotFound
		}
		return entity.User{}, ErrInntenal(err)
	}

	return vasya, nil
}

func (uc *UserUseCase) MakeSuperAdmin(id uuid.UUID) (entity.User, error) {

	vasya, err := uc.UserRepo.MakeSuperAdmin(id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.User{}, ErrNotFound
		}
		return entity.User{}, ErrInntenal(err)
	}

	return vasya, nil
}
