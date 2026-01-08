package repository

import (
	"database/internal/entity"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
)

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
	enviromentId, err := uuid.NewV4()

	if err != nil {
		return entity.Enviroment{}, err
	}

	enviroment.Id = enviromentId
	var id uuid.UUID

	id, err = uuid.NewV4()

	if err != nil {
		return entity.Enviroment{}, err
	}

	err = e.db.QueryRow(`INSERT INTO proICTIS_type_of_enviroment 
	(id, name) values ($1, $2)
	ON CONFLICT(name) DO NOTHING
	RETURNING id`, id, enviroment.TypeOfEnviroment).Scan(&id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = e.db.QueryRow(`SELECT id from proICTIS_type_of_enviroment where name = $1`, enviroment.TypeOfEnviroment).Scan(&id)

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

func (e *enviromentRepo) Edit(enviroment entity.Enviroment, ID uuid.UUID) (entity.Enviroment, error) {
	pat := `UPDATE proICTIS_enviroment`

	coloms := []string{}
	args := []any{}

	if enviroment.Name != "" {
		coloms = append(coloms, "name")
		args = append(args, enviroment.Name)
	}
	if enviroment.Description != "" {
		coloms = append(coloms, "description")
		args = append(args, enviroment.Description)
	}
	if enviroment.PhotoURL != "" {
		coloms = append(coloms, "photo_url")
		args = append(args, enviroment.PhotoURL)
	}
	if enviroment.Auditory != "" {
		coloms = append(coloms, "auditory")
		args = append(args, enviroment.Auditory)
	}
	if enviroment.TypeOfEnviroment != "" {
		id, err := uuid.NewV4()

		if err != nil {
			return entity.Enviroment{}, err
		}

		err = e.db.QueryRow(`INSERT INTO proICTIS_type_of_enviroment 
		(id, name) values ($1, $2)
		ON CONFLICT(name) DO NOTHING
		RETURNING id`, id, enviroment.TypeOfEnviroment).Scan(&id)

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				err = e.db.QueryRow(`SELECT id from proICTIS_type_of_enviroment where name = $1`, enviroment.TypeOfEnviroment).Scan(&id)

				if err != nil {
					return entity.Enviroment{}, err
				}
			} else {
				return entity.Enviroment{}, err
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

	editedEnviroment := entity.Enviroment{}
	args = append(args, ID)

	_, err := e.db.Exec(pat, args...)

	if err != nil {
		return entity.Enviroment{}, err
	}

	err = e.db.QueryRow(`select 
	e.id,
	e.name,
	e.description,
	e.photo_url,
	e.auditory,
	t.name
	from proICTIS_enviroment e
	join proICTIS_type_of_enviroment t on e.type_id = t.id
	where e.id = $1`, ID).Scan(
		&editedEnviroment.Id,
		&editedEnviroment.Name,
		&editedEnviroment.Description,
		&editedEnviroment.PhotoURL,
		&editedEnviroment.Auditory,
		&editedEnviroment.TypeOfEnviroment,
	)

	if err != nil {
		return entity.Enviroment{}, err
	}

	return editedEnviroment, nil
}

func (e *enviromentRepo) SetActive(id uuid.UUID, active bool) error {
	_, err := e.db.Exec(`update proICTIS_enviroment set is_active = $1 where id = $2`, active, id)
	if err != nil {
		return err
	}
	return nil
}

func (e *enviromentRepo) Delete(id uuid.UUID) error {
	_, err := e.db.Exec(`delete from proICTIS_enviroment where id = $1`, id)
	if err != nil {
		return err
	}
	return nil
}
