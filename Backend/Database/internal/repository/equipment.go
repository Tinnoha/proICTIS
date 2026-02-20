package repository

import (
	"context"
	"database/internal/entity"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type equipmentRepo struct {
	db Baza
}

func NewУquipmentRepo(postr *sqlx.DB, redis *redis.Client) *equipmentRepo {
	baza := Baza{post: postr, redis: redis}
	return &equipmentRepo{
		db: baza,
	}
}

func (e *equipmentRepo) GetAll() ([]entity.Equipment, error) {

	b, err := e.db.redis.Get(context.Background(), "Equipment").Result()
	var equipmen []entity.Equipment
	if err == nil {
		if err := json.Unmarshal([]byte(b), &equipmen); err == nil {
			return equipmen, nil
		}
	}

	if err == nil {
		return equipmen, nil
	}

	rows, err := e.db.post.Query(`SELECT 
		proICTIS_equipment.id,
		proICTIS_equipment.name,
		proICTIS_equipment.description,
		proICTIS_equipment.photo_url,
		proICTIS_equipment.auditory,
		proICTIS_equipment.is_active,
		proICTIS_type_of_equipment.name
		FROM proICTIS_equipment
		join proICTIS_type_of_equipment on proICTIS_type_of_equipment.id=proICTIS_equipment.type_id
		`)

	if err != nil {
		return []entity.Equipment{}, err
	}

	defer rows.Close()

	equipments := []entity.Equipment{}

	for rows.Next() {
		equipment := entity.Equipment{}

		err = rows.Scan(
			&equipment.Id,
			&equipment.Name,
			&equipment.Description,
			&equipment.PhotoURL,
			&equipment.Auditory,
			&equipment.IsActive,
			&equipment.TypeOfEquipment,
		)

		if err != nil {
			return []entity.Equipment{}, err
		}

		equipments = append(equipments, equipment)

	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	redisData, err := json.Marshal(equipments)

	if err != nil {
		return nil, err
	}

	if err := e.db.redis.Set(context.Background(), "Equipment", redisData, time.Minute).Err(); err != nil {
		fmt.Println("Error to save redis data", err)
	}

	return equipments, nil

}

func (e *equipmentRepo) GetByType(TypeOfEquipment string) ([]entity.Equipment, error) {
	rows, err := e.db.post.Query(`SELECT 
		proICTIS_equipment.id,
		proICTIS_equipment.name,
		proICTIS_equipment.description,
		proICTIS_equipment.photo_url,
		proICTIS_equipment.auditory,
		proICTIS_equipment.is_active,
		proICTIS_type_of_equipment.name
		FROM proICTIS_equipment
		join proICTIS_type_of_equipment on proICTIS_type_of_equipment.id=proICTIS_equipment.type_id
		where proICTIS_type_of_equipment.name = $1
		`, TypeOfEquipment)

	if err != nil {
		return []entity.Equipment{}, err
	}

	defer rows.Close()

	equipments := []entity.Equipment{}

	for rows.Next() {
		equipment := entity.Equipment{}

		err = rows.Scan(
			&equipment.Id,
			&equipment.Name,
			&equipment.Description,
			&equipment.PhotoURL,
			&equipment.Auditory,
			&equipment.IsActive,
			&equipment.TypeOfEquipment,
		)

		if err != nil {
			return []entity.Equipment{}, err
		}

		equipments = append(equipments, equipment)

	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if len(equipments) == 0 {
		return nil, sql.ErrNoRows
	}

	return equipments, nil

}

func (e *equipmentRepo) GetById(id uuid.UUID) (entity.Equipment, error) {
	equipment := entity.Equipment{}

	err := e.db.post.QueryRow(`SELECT 
		proICTIS_equipment.id,
		proICTIS_equipment.name,
		proICTIS_equipment.description,
		proICTIS_equipment.photo_url,
		proICTIS_equipment.auditory,
		proICTIS_equipment.is_active,
		proICTIS_type_of_equipment.name
		FROM proICTIS_equipment
		join proICTIS_type_of_equipment on proICTIS_type_of_equipment.id=proICTIS_equipment.type_id
		where proICTIS_equipment.id = $1
		`, id).Scan(
		&equipment.Id,
		&equipment.Name,
		&equipment.Description,
		&equipment.PhotoURL,
		&equipment.Auditory,
		&equipment.IsActive,
		&equipment.TypeOfEquipment,
	)

	if err != nil {
		return entity.Equipment{}, err
	}

	return equipment, nil
}

func (e *equipmentRepo) GetTypes() ([]entity.TypeOfEquipment, error) {
	rows, err := e.db.post.Query(`SELECT 
		id,
		name
		FROM proICTIS_type_of_equipment
		`)

	if err != nil {
		return []entity.TypeOfEquipment{}, err
	}

	defer rows.Close()

	equipments := []entity.TypeOfEquipment{}

	for rows.Next() {
		TypeOfEquipment := entity.TypeOfEquipment{}

		err = rows.Scan(
			&TypeOfEquipment.Id,
			&TypeOfEquipment.Name,
		)

		if err != nil {
			return []entity.TypeOfEquipment{}, err
		}

		equipments = append(equipments, TypeOfEquipment)

	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return equipments, nil

}

func (e *equipmentRepo) Add(equipment entity.Equipment) (entity.Equipment, error) {
	equipmentId, err := uuid.NewV4()

	if err != nil {
		return entity.Equipment{}, err
	}

	equipment.Id = equipmentId
	var id uuid.UUID

	id, err = uuid.NewV4()

	if err != nil {
		return entity.Equipment{}, err
	}

	err = e.db.post.QueryRow(`INSERT INTO proICTIS_type_of_equipment 
	(id, name) values ($1, $2)
	ON CONFLICT(name) DO NOTHING
	RETURNING id`, id, equipment.TypeOfEquipment).Scan(&id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = e.db.post.QueryRow(`SELECT id from proICTIS_type_of_equipment where name = $1`, equipment.TypeOfEquipment).Scan(&id)

			if err != nil {
				return entity.Equipment{}, err
			}
		} else {
			return entity.Equipment{}, err
		}
	}

	_, err = e.db.post.Exec(`INSERT INTO proICTIS_equipment (id,name,description,photo_url,auditory,is_active,type_id) values($1,$2,$3,$4,$5,$6,$7)`,
		equipment.Id,
		equipment.Name,
		equipment.Description,
		equipment.PhotoURL,
		equipment.Auditory,
		equipment.IsActive,
		id,
	)
	if err != nil {
		return entity.Equipment{}, err
	}

	err = e.db.redis.Del(context.Background(), "Equipment").Err()

	if err != nil {
		return entity.Equipment{}, err
	}

	return equipment, nil
}

func (e *equipmentRepo) Edit(equipment entity.Equipment, ID uuid.UUID) (entity.Equipment, error) {
	pat := `UPDATE proICTIS_equipment`

	coloms := []string{}
	args := []any{}

	if equipment.Name != "" {
		coloms = append(coloms, "name")
		args = append(args, equipment.Name)
	}
	if equipment.Description != "" {
		coloms = append(coloms, "description")
		args = append(args, equipment.Description)
	}
	if equipment.PhotoURL != "" {
		coloms = append(coloms, "photo_url")
		args = append(args, equipment.PhotoURL)
	}
	if equipment.Auditory != "" {
		coloms = append(coloms, "auditory")
		args = append(args, equipment.Auditory)
	}
	if equipment.TypeOfEquipment != "" {
		id, err := uuid.NewV4()

		if err != nil {
			return entity.Equipment{}, err
		}

		err = e.db.post.QueryRow(`INSERT INTO proICTIS_type_of_equipment 
		(id, name) values ($1, $2)
		ON CONFLICT(name) DO NOTHING
		RETURNING id`, id, equipment.TypeOfEquipment).Scan(&id)

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				err = e.db.post.QueryRow(`SELECT id from proICTIS_type_of_equipment where name = $1`, equipment.TypeOfEquipment).Scan(&id)

				if err != nil {
					return entity.Equipment{}, err
				}
			} else {
				return entity.Equipment{}, err
			}
		}

		coloms = append(coloms, "type_id")
		args = append(args, id)
	}

	if len(coloms) != len(args) {
		panic(errors.New("NON PREDICTED"))
	}

	setClas := make([]string, len(coloms))

	for k, v := range coloms {
		setClas[k] = fmt.Sprintf("%s = $%d", v, k+1)
	}

	pat += fmt.Sprintf(` set %s where id = $%d`,
		strings.Join(setClas, ", "), len(coloms)+1)

	editedEquipment := entity.Equipment{}
	args = append(args, ID)

	rez, err := e.db.post.Exec(pat, args...)

	if err != nil {
		return entity.Equipment{}, err
	}

	c, err := rez.RowsAffected()

	if err != nil {
		return entity.Equipment{}, err
	}

	if c == 0 {
		return entity.Equipment{}, sql.ErrNoRows
	}

	err = e.db.post.QueryRow(`select 
	e.id,
	e.name,
	e.description,
	e.photo_url,
	e.auditory,
	t.name
	from proICTIS_equipment e
	join proICTIS_type_of_equipment t on e.type_id = t.id
	where e.id = $1`, ID).Scan(
		&editedEquipment.Id,
		&editedEquipment.Name,
		&editedEquipment.Description,
		&editedEquipment.PhotoURL,
		&editedEquipment.Auditory,
		&editedEquipment.TypeOfEquipment,
	)

	if err != nil {
		return entity.Equipment{}, err
	}

	return editedEquipment, nil
}

func (e *equipmentRepo) SetActive(id uuid.UUID, active bool) error {
	rez, err := e.db.post.Exec(`update proICTIS_equipment set is_active = $1 where id = $2`, active, id)
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

func (e *equipmentRepo) Delete(id uuid.UUID) error {
	rez, err := e.db.post.Exec(`delete from proICTIS_equipment where id = $1`, id)
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
