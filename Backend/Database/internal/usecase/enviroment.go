package usecase

import (
	"database/internal/entity"
	"database/sql"
	"errors"

	"github.com/gofrs/uuid"
)

type EnviromentRepository interface {
	GetAll() ([]entity.Enviroment, error)
	GetByType(TypeOfEnviroment string) ([]entity.Enviroment, error)
	GetById(id uuid.UUID) (entity.Enviroment, error)
	GetTypes() ([]entity.TypeOfEnviroment, error)

	Add(enviroment entity.Enviroment) (entity.Enviroment, error)

	Edit(enviroment entity.Enviroment, id uuid.UUID) (entity.Enviroment, error)
	SetActive(id uuid.UUID, active bool) error

	Delete(id uuid.UUID) error
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

func (uc *EnviromentUsecase) GetAll() ([]entity.Enviroment, error) {
	return uc.EnviromentRepo.GetAll()
}

func (uc *EnviromentUsecase) GetByType(TypeOfEnviroment string) ([]entity.Enviroment, error) {
	vr, err := uc.EnviromentRepo.GetByType(TypeOfEnviroment)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []entity.Enviroment{}, ErrNotFound
		}
		return []entity.Enviroment{}, ErrInntenal(err)
	}

	return vr, nil
}

func (uc *EnviromentUsecase) GetById(id uuid.UUID) (entity.Enviroment, error) {
	vr, err := uc.EnviromentRepo.GetById(id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Enviroment{}, ErrNotFound
		}
		return entity.Enviroment{}, ErrInntenal(err)
	}

	return vr, nil
}

func (uc *EnviromentUsecase) GetTypes() ([]entity.TypeOfEnviroment, error) {
	return uc.EnviromentRepo.GetTypes()
}

func (uc *EnviromentUsecase) Add(Enviroment []entity.Enviroment) ([]entity.Enviroment, error) {
	rez := []entity.Enviroment{}
	for _, vr := range Enviroment {
		id, err := uuid.NewV4()

		if err != nil {
			return []entity.Enviroment{}, ErrInntenal(err)
		}

		vr.Id = id

		added, err := uc.EnviromentRepo.Add(vr)

		if err != nil {
			return []entity.Enviroment{}, ErrInntenal(err)
		}

		rez = append(rez, added)
	}
	return rez, nil
}

func (uc *EnviromentUsecase) Edit(id uuid.UUID, enviroment entity.Enviroment) (entity.Enviroment, error) {
	enviroment.Id = id

	glasses, err := uc.EnviromentRepo.Edit(enviroment, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Enviroment{}, ErrNotFound
		}
		return entity.Enviroment{}, ErrInntenal(err)
	}

	return glasses, nil
}

func (uc *EnviromentUsecase) Delete(id uuid.UUID) error {
	err := uc.EnviromentRepo.Delete(id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
		return ErrInntenal(err)
	}

	return nil
}

func (uc *EnviromentUsecase) SetActive(id uuid.UUID, active bool) error {
	err := uc.EnviromentRepo.SetActive(id, active)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
		return ErrInntenal(err)
	}

	return nil
}
