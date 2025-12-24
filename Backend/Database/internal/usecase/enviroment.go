package usecase

import (
	"database/internal/entity"

	"github.com/gofrs/uuid"
)

type EnviromentRepository interface {
	GetAll() []entity.Enviroment
	GetByType(TypeOfEnviroment string) ([]entity.Enviroment, error)
	GetById(id uuid.UUID) (entity.Enviroment, error)
	Add(enviroment entity.Enviroment) (entity.Enviroment, error)
	AddType(TypeOfEnviroment string) error
	Edit(enviroment entity.Enviroment) (entity.Enviroment, error)
	Delete(id uuid.UUID) error
	IsExists(id uuid.UUID) bool
	NameIsTaken(name string) bool
	TypeIsExsits(TypeOfEnviroment string) bool
	SetActive(id uuid.UUID, active bool) (entity.Enviroment, error)
}

type EnviromentUsecase struct {
	EnviromentRepo EnviromentRepository
	BookingRepo    BookingRepository
}

func NewEnviromentUsecase(
	bookRepo BookingRepository,
	enviromentRepo EnviromentRepository,
) *EnviromentUsecase {
	return &EnviromentUsecase{
		BookingRepo:    bookRepo,
		EnviromentRepo: enviromentRepo,
	}
}

func (uc *EnviromentUsecase) GetAll() []entity.Enviroment {
	return uc.EnviromentRepo.GetAll()
}

func (uc *EnviromentUsecase) GetByType(TypeOfEnviroment string) ([]entity.Enviroment, error) {
	exists := uc.EnviromentRepo.TypeIsExsits(TypeOfEnviroment)

	if !exists {
		return []entity.Enviroment{}, ErrNotFound
	}

	vr, err := uc.EnviromentRepo.GetByType(TypeOfEnviroment)

	if err != nil {
		return []entity.Enviroment{}, ErrInntenal(err)
	}

	return vr, nil
}

func (uc *EnviromentUsecase) GetById(id uuid.UUID) (entity.Enviroment, error) {
	exists := uc.EnviromentRepo.IsExists(id)

	if !exists {
		return entity.Enviroment{}, ErrNotFound
	}

	vr, err := uc.EnviromentRepo.GetById(id)

	if err != nil {
		return entity.Enviroment{}, ErrInntenal(err)
	}

	return vr, nil
}

func (uc *EnviromentUsecase) Add(Enviroment []entity.Enviroment) ([]entity.Enviroment, error) {
	for _, vr := range Enviroment {
		taken := uc.EnviromentRepo.NameIsTaken(vr.Name)

		if taken {
			return []entity.Enviroment{}, ErrThisExists("name", vr.Name)
		}

		TypeIsExsits := uc.EnviromentRepo.TypeIsExsits(vr.TypeOfEnviroment)

		if !TypeIsExsits {
			err := uc.AddType(vr.TypeOfEnviroment)

			if err != nil {
				return []entity.Enviroment{}, ErrInntenal(err)
			}

		}

		_, err := uc.EnviromentRepo.Add(vr)

		if err != nil {
			return []entity.Enviroment{}, ErrInntenal(err)
		}
	}
	return Enviroment, nil
}

func (uc *EnviromentUsecase) AddType(TypeOfEnviroment string) error {
	TypeIsExsits := uc.EnviromentRepo.TypeIsExsits(TypeOfEnviroment)

	if !TypeIsExsits {
		err := uc.AddType(TypeOfEnviroment)

		if err != nil {
			return ErrInntenal(err)
		}

	} else {
		return ErrThisExists("type", TypeOfEnviroment)
	}
	return nil
}

func (uc *EnviromentUsecase) Edit(id uuid.UUID, enviroment entity.Enviroment) (entity.Enviroment, error) {
	exists := uc.EnviromentRepo.IsExists(id)

	if !exists {
		return entity.Enviroment{}, ErrNotFound
	}

	glasses, err := uc.EnviromentRepo.Edit(enviroment)

	if err != nil {
		return entity.Enviroment{}, ErrInntenal(err)
	}

	return glasses, nil
}

func (uc *EnviromentUsecase) Delete(id uuid.UUID) error {
	exists := uc.EnviromentRepo.IsExists(id)

	if !exists {
		return ErrNotFound
	}

	err := uc.EnviromentRepo.Delete(id)

	if err != nil {
		return ErrInntenal(err)
	}

	return nil
}

func (uc *EnviromentUsecase) IsExists(id uuid.UUID) bool {
	return uc.EnviromentRepo.IsExists(id)
}

func (uc *EnviromentUsecase) TypeIsExsits(TypeOfEnviroment string) bool {
	return uc.EnviromentRepo.TypeIsExsits(TypeOfEnviroment)
}

func (uc *EnviromentUsecase) SetActive(id uuid.UUID, active bool) (entity.Enviroment, error) {
	exists := uc.EnviromentRepo.IsExists(id)

	if !exists {
		return entity.Enviroment{}, ErrNotFound
	}

	activ, err := uc.EnviromentRepo.SetActive(id, active)

	if err != nil {
		return entity.Enviroment{}, ErrInntenal(err)
	}

	return activ, nil
}
