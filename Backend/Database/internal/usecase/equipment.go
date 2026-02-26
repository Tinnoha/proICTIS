package usecase

import (
	"database/internal/entity"
	"database/sql"
	"errors"
	"fmt"

	"github.com/gofrs/uuid"
)

type EquipmentRepository interface {
	GetAll() ([]entity.Equipment, error)
	GetByType(TypeOfEquipment string) ([]entity.Equipment, error)
	GetById(id uuid.UUID) (entity.Equipment, error)

	GetTypes() ([]entity.TypeOfEquipment, error)                     // -+
	AddType(entity.TypeOfEquipment) (entity.TypeOfEquipment, error)  // -+
	EditType(entity.TypeOfEquipment) (entity.TypeOfEquipment, error) // -+
	DeleteType(id uuid.UUID) error                                   // -+

	Add(equipment entity.Equipment) (entity.Equipment, error)

	Edit(equipment entity.Equipment, id uuid.UUID) (entity.Equipment, error)
	SetActive(id uuid.UUID, active bool) error

	Delete(id uuid.UUID) error
}

type EquipmentUsecase struct {
	EquipmentRepo EquipmentRepository
	BookingRepo   BookingRepository
	UserRepo      UserRepository
}

func NewEquipmentUsecase(
	bookRepo BookingRepository,
	equipmentRepo EquipmentRepository,
	userRepo UserRepository,
) *EquipmentUsecase {
	return &EquipmentUsecase{
		BookingRepo:   bookRepo,
		EquipmentRepo: equipmentRepo,
		UserRepo:      userRepo,
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

func (uc *EquipmentUsecase) AddTypes(AdminId uuid.UUID, types []entity.TypeOfEquipment) ([]entity.TypeOfEquipment, error) {
	admin, err := uc.UserRepo.IsAdmin(AdminId)
	if err != nil {
		return []entity.TypeOfEquipment{}, err
	}

	if !admin {
		fmt.Println(AdminId, "Try to fatatl our service!")
		return []entity.TypeOfEquipment{}, errors.New("You are not a admin ! ! ! ! ! !")
	}
	resultat := []entity.TypeOfEquipment{}

	for _, typee := range types {
		id, err := uuid.NewV4()

		if err != nil {
			return []entity.TypeOfEquipment{}, ErrInntenal(err)
		}

		typee.Id = id

		added, err := uc.EquipmentRepo.AddType(typee)

		if err != nil {
			return []entity.TypeOfEquipment{}, ErrInntenal(err)
		}

		resultat = append(resultat, added)
	}

	return resultat, nil

}

func (uc *EquipmentUsecase) EditType(AdminId uuid.UUID, tip entity.TypeOfEquipment) (entity.TypeOfEquipment, error) {
	admin, err := uc.UserRepo.IsAdmin(AdminId)
	if err != nil {
		return entity.TypeOfEquipment{}, err
	}

	if !admin {
		fmt.Println(AdminId, "Try to fatatl our service!")
		return entity.TypeOfEquipment{}, errors.New("You are not a admin ! ! ! ! ! !")
	}
	edited, err := uc.EquipmentRepo.EditType(tip)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.TypeOfEquipment{}, ErrNotFound
		}
		return entity.TypeOfEquipment{}, ErrInntenal(err)
	}

	return edited, nil
}

func (uc *EquipmentUsecase) DeleteType(AdminId, id uuid.UUID) error {
	admin, err := uc.UserRepo.IsAdmin(AdminId)
	if err != nil {
		return err
	}

	if !admin {
		fmt.Println(AdminId, "Try to fatatl our service!")
		return errors.New("You are not a admin ! ! ! ! ! !")
	}
	err = uc.EquipmentRepo.DeleteType(id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}

		return ErrInntenal(err)
	}
	return nil
}

func (uc *EquipmentUsecase) Add(Equipment entity.Equipments) ([]entity.Equipment, error) {
	admin, err := uc.UserRepo.IsAdmin(Equipment.AdminId)
	if err != nil {
		return []entity.Equipment{}, err
	}

	if !admin {
		fmt.Println(AdminId, "Try to fatatl our service!")
		return []entity.Equipment{}, errors.New("You are not a admin ! ! ! ! ! !")
	}
	rez := []entity.Equipment{}
	for _, vr := range Equipment.Tovars {
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

func (uc *EquipmentUsecase) Edit(AdminId, id uuid.UUID, equipment entity.Equipment) (entity.Equipment, error) {
	admin, err := uc.UserRepo.IsAdmin(AdminId)
	if err != nil {
		return entity.Equipment{}, err
	}

	if !admin {
		fmt.Println(AdminId, "Try to fatatl our service!")
		return entity.Equipment{}, errors.New("You are not a admin ! ! ! ! ! !")
	}
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

func (uc *EquipmentUsecase) Delete(AdminId, id uuid.UUID) error {
	admin, err := uc.UserRepo.IsAdmin(AdminId)
	if err != nil {
		return err
	}

	if !admin {
		fmt.Println(AdminId, "Try to fatatl our service!")
		return errors.New("You are not a admin ! ! ! ! ! !")
	}
	err = uc.EquipmentRepo.Delete(id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
		return ErrInntenal(err)
	}

	return nil
}

func (uc *EquipmentUsecase) SetActive(AdminId, id uuid.UUID, active bool) error {
	admin, err := uc.UserRepo.IsAdmin(AdminId)
	if err != nil {
		return err
	}

	if !admin {
		fmt.Println(AdminId, "Try to fatatl our service!")
		return errors.New("You are not a admin ! ! ! ! ! !")
	}
	err = uc.EquipmentRepo.SetActive(id, active)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
		return ErrInntenal(err)
	}

	return nil
}
