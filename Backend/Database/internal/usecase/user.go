package usecase

import (
	"database/internal/entity"

	"github.com/gofrs/uuid"
)

type UserRepository interface {
	GetById(id uuid.UUID) (entity.User, error)
	GetByEmail(email string) (entity.User, error)
	CreateUser(entity.User) (entity.User, error)
	MakeAdmin(id uuid.UUID) (entity.User, error)
	MakeSuperAdmin(id uuid.UUID) (entity.User, error)
	IsExistsById(id uuid.UUID) bool
	IsExistsByEmail(email string) bool
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

func (uc *UserUseCase) GetById(id uuid.UUID) (entity.User, error) {
	exists := uc.UserRepo.IsExistsById(id)

	if !exists {
		return entity.User{}, ErrNotFound
	}

	vasya, err := uc.UserRepo.GetById(id)
	if err != nil {
		return entity.User{}, ErrInntenal(err)
	}

	return vasya, nil
}

func (uc *UserUseCase) GetByEmail(email string) (entity.User, error) {
	exists := uc.UserRepo.IsExistsByEmail(email)

	if !exists {
		return entity.User{}, ErrNotFound
	}

	vasya, err := uc.UserRepo.GetByEmail(email)
	if err != nil {
		return entity.User{}, ErrInntenal(err)
	}

	return vasya, nil
}

func (uc *UserUseCase) CreateUser(vasy []entity.User) ([]entity.User, error) {
	rezult := []entity.User{}

	for _, annya := range vasy {

		exists := uc.UserRepo.IsExistsByEmail(annya.Email)

		if exists {
			return []entity.User{}, ErrThisExists("email", annya.Email)
		}

		vasya, err := uc.UserRepo.CreateUser(annya)

		if err != nil {
			return []entity.User{}, ErrInntenal(err)
		}

		vasya.Id, err = uuid.NewV4()

		if err != nil {
			return []entity.User{}, ErrInntenal(err)
		}

		rezult = append(rezult, vasya)
	}

	return vasy, nil
}

func (uc *UserUseCase) MakeAdmin(id uuid.UUID) (entity.User, error) {
	exists := uc.UserRepo.IsExistsById(id)

	if !exists {
		return entity.User{}, ErrNotFound
	}

	vasya, err := uc.UserRepo.MakeAdmin(id)

	if err != nil {
		return entity.User{}, ErrInntenal(err)
	}

	return vasya, nil
}

func (uc *UserUseCase) MakeSuperAdmin(id uuid.UUID) (entity.User, error) {
	exists := uc.UserRepo.IsExistsById(id)

	if !exists {
		return entity.User{}, ErrNotFound
	}

	vasya, err := uc.UserRepo.MakeSuperAdmin(id)

	if err != nil {
		return entity.User{}, ErrInntenal(err)
	}

	return vasya, nil
}

func (uc *UserUseCase) IsExistsById(id uuid.UUID) bool {
	return uc.UserRepo.IsExistsById(id)
}

func (uc *UserUseCase) IsExistsByEmail(email string) bool {
	return uc.IsExistsByEmail(email)
}
