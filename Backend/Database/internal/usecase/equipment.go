package usecase

import (
	"database/internal/entity"
	"database/sql"
	"errors"

	"github.com/gofrs/uuid"
)

type EquipmentRepository interface {
	GetAll() ([]entity.Equipment, error)
	GetByType(TypeOfEquipment string) ([]entity.Equipment, error)
	GetById(id uuid.UUID) (entity.Equipment, error)
	GetTypes() ([]entity.TypeOfEquipment, error)

	Add(equipment entity.Equipment) (entity.Equipment, error)

	Edit(equipment entity.Equipment, id uuid.UUID) (entity.Equipment, error)
	SetActive(id uuid.UUID, active bool) error

	Delete(id uuid.UUID) error
}

type EquipmentUsecase struct {
	EquipmentRepo EquipmentRepository
	BookingRepo   BookingRepository
}

func NewEquipmentUsecase(
	bookRepo BookingRepository,
	equipmentRepo EquipmentRepository,
) *EquipmentUsecase {
	return &EquipmentUsecase{
		BookingRepo:   bookRepo,
		EquipmentRepo: equipmentRepo,
	}
}

func (uc *EquipmentUsecase) GetAll() ([]entity.Equipment, error) {
	return uc.EquipmentRepo.GetAll()
}

func (uc *EquipmentUsecase) GetByType(TypeOfEquipment string) ([]entity.Equipment, error) {
	vr, err := uc.EquipmentRepo.GetByType(TypeOfEquipment)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []entity.Equipment{}, ErrNotFound
		}
		return []entity.Equipment{}, ErrInntenal(err)
	}

	return vr, nil
}

func (uc *EquipmentUsecase) GetById(id uuid.UUID) (entity.Equipment, error) {
	vr, err := uc.EquipmentRepo.GetById(id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Equipment{}, ErrNotFound
		}
		return entity.Equipment{}, ErrInntenal(err)
	}

	return vr, nil
}

func (uc *EquipmentUsecase) GetTypes() ([]entity.TypeOfEquipment, error) {
	return uc.EquipmentRepo.GetTypes()
}

func (uc *EquipmentUsecase) Add(Equipment []entity.Equipment) ([]entity.Equipment, error) {
	rez := []entity.Equipment{}
	for _, vr := range Equipment {
		id, err := uuid.NewV4()

		if err != nil {
			return []entity.Equipment{}, ErrInntenal(err)
		}

		vr.Id = id

		added, err := uc.EquipmentRepo.Add(vr)

		if err != nil {
			return []entity.Equipment{}, ErrInntenal(err)
		}

		rez = append(rez, added)
	}
	return rez, nil
}

func (uc *EquipmentUsecase) Edit(id uuid.UUID, equipment entity.Equipment) (entity.Equipment, error) {
	equipment.Id = id

	glasses, err := uc.EquipmentRepo.Edit(equipment, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Equipment{}, ErrNotFound
		}
		return entity.Equipment{}, ErrInntenal(err)
	}

	return glasses, nil
}

func (uc *EquipmentUsecase) Delete(id uuid.UUID) error {
	err := uc.EquipmentRepo.Delete(id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
		return ErrInntenal(err)
	}

	return nil
}

func (uc *EquipmentUsecase) SetActive(id uuid.UUID, active bool) error {
	err := uc.EquipmentRepo.SetActive(id, active)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
		return ErrInntenal(err)
	}

	return nil
}
