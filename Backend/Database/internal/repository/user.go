package repository

import (
	"database/internal/entity"

	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
)

/*
type User struct {
	Id            uuid.UUID `json:"id"`
	FirstName     string    `json:"first_name"`
	SecondName    string    `json:"second_name"`
	Email         string    `json:"email"`
	AvatarURL     string    `json:"avatar_url"`
	Role          string    `json:"role"`
	TokenProvider int       `json:"tokenProvider"`
}


type UserRepository interface {
	GetById(id uuid.UUID) (entity.User, error)
	GetByEmail(email string) (entity.User, error)
	CreateUser(entity.User) (entity.User, error)
	MakeAdmin(id uuid.UUID) (entity.User, error)
	MakeSuperAdmin(id uuid.UUID) (entity.User, error)
	CreateUser(vasya entity.User) entity.User
}
*/

type userRepo struct {
	db sqlx.DB
}

func NewUserRepo(db sqlx.DB) *userRepo {
	return &userRepo{
		db: db,
	}
}

func (u *userRepo) GetAll() ([]entity.User, error) {
	rows, err := u.db.Query(`SELECT id, first_name, second_name, email, avatar_url, role, token_provider FROM proICTIS_user`)
	if err != nil {
		return []entity.User{}, err
	}

	defer rows.Close()

	res := []entity.User{}
	for rows.Next() {
		vasya := entity.User{}

		err := rows.Scan(
			&vasya.Id,
			&vasya.FirstName,
			&vasya.SecondName,
			&vasya.Email,
			&vasya.AvatarURL,
			&vasya.Role,
			&vasya.TokenProvider,
		)

		if err != nil {
			return []entity.User{}, err
		}

		res = append(res, vasya)
	}

	if err := rows.Err(); err != nil {
		return []entity.User{}, err
	}

	return res, nil
}

func (u *userRepo) GetById(id uuid.UUID) (entity.User, error) {
	vasya := entity.User{}
	err := u.db.QueryRow(`Select id, first_name, second_name, email, avatar_url,role, token_provider 
	from proICTIS_user 
	where id = $1`, id).Scan(
		&vasya.Id,
		&vasya.FirstName,
		&vasya.SecondName,
		&vasya.Email,
		&vasya.AvatarURL,
		&vasya.Role,
		&vasya.TokenProvider,
	)

	if err != nil {
		return entity.User{}, err
	}

	return vasya, nil
}

func (u *userRepo) GetByEmail(email string) (entity.User, error) {
	vasya := entity.User{}
	err := u.db.QueryRow(`Select id, first_name, second_name, email, avatar_url,role, token_provider 
	from proICTIS_user 
	where email = $1`, email).Scan(
		&vasya.Id,
		&vasya.FirstName,
		&vasya.SecondName,
		&vasya.Email,
		&vasya.AvatarURL,
		&vasya.Role,
		&vasya.TokenProvider,
	)

	if err != nil {
		return entity.User{}, err
	}

	return vasya, nil
}

func (u *userRepo) CreateUser(user entity.User) (entity.User, error) {
	_, err := u.db.Exec(`INSERT INTO proICTIS_user (id, first_name, second_name, email, avatar_url,role, token_provider) VALUES($1,$2,$3,$4,$5,$6,$7)`,
		user.Id,
		user.FirstName,
		user.SecondName,
		user.Email,
		user.AvatarURL,
		user.Role,
		user.TokenProvider,
	)

	if err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func (u *userRepo) MakeAdmin(id uuid.UUID) (entity.User, error) {
	user := entity.User{}

	err := u.db.QueryRow(
		`UPDATE PROICTIS_user
		set role='Admin'
		WHERE id = $1
		RETURNING first_name, second_name, email,avatar_url,role,token_provider
		`, id).Scan(
		&user.FirstName,
		&user.SecondName,
		&user.Email,
		&user.AvatarURL,
		&user.Role,
		&user.TokenProvider,
	)

	if err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func (u *userRepo) MakeSuperAdmin(id uuid.UUID) (entity.User, error) {
	user := entity.User{}

	err := u.db.QueryRow(
		`UPDATE PROICTIS_user
		set role='Super_Admin'
		WHERE id = $1
		RETURNING first_name, second_name, email,avatar_url,role,token_provider
		`, id).Scan(
		&user.FirstName,
		&user.SecondName,
		&user.Email,
		&user.AvatarURL,
		&user.Role,
		&user.TokenProvider,
	)

	if err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func (u *userRepo) IsAdmin(id uuid.UUID) (bool, error) {
	row, err := u.db.Query(`SELECT role FROM proICTIS_use where id = $1`, id)

	if err != nil {
		return false, err
	}

	role := ""
	err = row.Scan(&role)

	if err != nil {
		return false, err
	}

	if role == "admin" {
		return true, nil
	} else {
		return false, nil
	}
}

func (u *userRepo) IsSuperAdmin(id uuid.UUID) (bool, error) {
	row, err := u.db.Query(`SELECT role FROM proICTIS_use where id = $1`, id)

	if err != nil {
		return false, err
	}

	role := ""
	err = row.Scan(&role)

	if err != nil {
		return false, err
	}

	if role == "Super_Admin" {
		return true, nil
	} else {
		return false, nil
	}
}
