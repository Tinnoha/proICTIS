package repository

import (
	"database/internal/entity"
	"database/sql"
	"errors"

	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
)

/*
GetAll() ([]entity.Enviroment, error)
GetByType(TypeOfEnviroment string) ([]entity.Enviroment, error)
GetById(id uuid.UUID) (entity.Enviroment, error)
GetTypes() ([]entity.TypeOfEnviroment, error)

Add(enviroment entity.Enviroment) (entity.Enviroment, error)

Edit(enviroment entity.Enviroment) (entity.Enviroment, error)
SetActive(id uuid.UUID, active bool) (entity.Enviroment, error)

Delete(id uuid.UUID) error

`CREATE TABLE IF NOT EXISTS proICTIS_type_of_enviroment(
	id UUID PRIMARY KEY,
	name varchar(255) NOT NULL UNIQUE

`CREATE TABLE IF NOT EXISTS proICTIS_enviroment(
	id UUID PRIMARY KEY,
	name varchar(255),
	description TEXT,
	photo_url varchar(255),
	type_id UUID NOT NULL REFERENCES proICTIS_type_of_enviroment(id) ON DELETE RESTRICT,
	auditory varchar(255),
	is_active boolean
	)
`)
*/

type enviromentRepo struct {
	db sqlx.DB
}

func NewEnviromentRepo(db sqlx.DB) *enviromentRepo {
	return &enviromentRepo{
		db: db,
	}
}

func (e *enviromentRepo) GetAll() ([]entity.Enviroment, error) {
	rows, err := e.db.Query(`SELECT 
		proICTIS_enviroment.id,
		proICTIS_enviroment.name,
		proICTIS_enviroment.description,
		proICTIS_enviroment.photo_url,
		proICTIS_enviroment.auditory,
		proICTIS_enviroment.is_active,
		proICTIS_type_of_enviroment.name
		FROM proICTIS_enviroment
		join proICTIS_type_of_enviroment on proICTIS_type_of_enviroment.id=proICTIS_enviroment.type_id
		`)

	if err != nil {
		return []entity.Enviroment{}, err
	}

	defer rows.Close()

	enviroments := []entity.Enviroment{}

	for rows.Next() {
		enviroment := entity.Enviroment{}

		err = rows.Scan(
			&enviroment.Id,
			&enviroment.Name,
			&enviroment.Description,
			&enviroment.PhotoURL,
			&enviroment.Auditory,
			&enviroment.IsActive,
			&enviroment.TypeOfEnviroment,
		)

		if err != nil {
			return []entity.Enviroment{}, err
		}

		enviroments = append(enviroments, enviroment)

	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return enviroments, nil

}

func (e *enviromentRepo) GetByType(TypeOfEnviroment string) ([]entity.Enviroment, error) {
	rows, err := e.db.Query(`SELECT 
		proICTIS_enviroment.id,
		proICTIS_enviroment.name,
		proICTIS_enviroment.description,
		proICTIS_enviroment.photo_url,
		proICTIS_enviroment.auditory,
		proICTIS_enviroment.is_active,
		proICTIS_type_of_enviroment.name
		FROM proICTIS_enviroment
		join proICTIS_type_of_enviroment on proICTIS_type_of_enviroment.id=proICTIS_enviroment.type_id
		where proICTIS_type_of_enviroment.name = $1
		`, TypeOfEnviroment)

	if err != nil {
		return []entity.Enviroment{}, err
	}

	defer rows.Close()

	enviroments := []entity.Enviroment{}

	for rows.Next() {
		enviroment := entity.Enviroment{}

		err = rows.Scan(
			&enviroment.Id,
			&enviroment.Name,
			&enviroment.Description,
			&enviroment.PhotoURL,
			&enviroment.Auditory,
			&enviroment.IsActive,
			&enviroment.TypeOfEnviroment,
		)

		if err != nil {
			return []entity.Enviroment{}, err
		}

		enviroments = append(enviroments, enviroment)

	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if len(enviroments) == 0 {
		return nil, sql.ErrNoRows
	}

	return enviroments, nil

}

func (e *enviromentRepo) GetById(id uuid.UUID) (entity.Enviroment, error) {
	enviroment := entity.Enviroment{}

	err := e.db.QueryRow(`SELECT 
		proICTIS_enviroment.id,
		proICTIS_enviroment.name,
		proICTIS_enviroment.description,
		proICTIS_enviroment.photo_url,
		proICTIS_enviroment.auditory,
		proICTIS_enviroment.is_active,
		proICTIS_type_of_enviroment.name
		FROM proICTIS_enviroment
		join proICTIS_type_of_enviroment on proICTIS_type_of_enviroment.id=proICTIS_enviroment.type_id
		where proICTIS_enviroment.id = $1
		`, id).Scan(
		&enviroment.Id,
		&enviroment.Name,
		&enviroment.Description,
		&enviroment.PhotoURL,
		&enviroment.Auditory,
		&enviroment.IsActive,
		&enviroment.TypeOfEnviroment,
	)

	if err != nil {
		return entity.Enviroment{}, err
	}

	return enviroment, nil
}

func (e *enviromentRepo) GetTypes() ([]entity.TypeOfEnviroment, error) {
	rows, err := e.db.Query(`SELECT 
		id,
		name
		FROM proICTIS_type_of_enviroment
		`)

	if err != nil {
		return []entity.TypeOfEnviroment{}, err
	}

	defer rows.Close()

	enviroments := []entity.TypeOfEnviroment{}

	for rows.Next() {
		TypeOfEnviroment := entity.TypeOfEnviroment{}

		err = rows.Scan(
			&TypeOfEnviroment.Id,
			&TypeOfEnviroment.Name,
		)

		if err != nil {
			return []entity.TypeOfEnviroment{}, err
		}

		enviroments = append(enviroments, TypeOfEnviroment)

	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return enviroments, nil

}

func (e *enviromentRepo) Add(enviroment entity.Enviroment) (entity.Enviroment, error) {
	var id uuid.UUID

	err := e.db.QueryRow(`select id from proICTIS_type_of_enviroment where name = $1`, enviroment.TypeOfEnviroment).Scan(&id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			id, err = uuid.NewV4()

			if err != nil {
				return entity.Enviroment{}, err
			}

			_, err = e.db.Exec(`INSERT INTO proICTIS_type_of_enviroment (id, name) values ($1, $2)`, id, enviroment.TypeOfEnviroment)

			if err != nil {
				return entity.Enviroment{}, err
			}
		} else {
			return entity.Enviroment{}, err
		}
	}

	_, err = e.db.Exec(`INSERT INTO proICTIS_enviroment (id,name,description,photo_url,auditory,is_active,type_id) values($1,$2,$3,$4,$5,$6,$7)`,
		enviroment.Id,
		enviroment.Name,
		enviroment.Description,
		enviroment.PhotoURL,
		enviroment.Auditory,
		enviroment.IsActive,
		id,
	)
	if err != nil {
		return entity.Enviroment{}, err
	}

	return enviroment, nil
}
