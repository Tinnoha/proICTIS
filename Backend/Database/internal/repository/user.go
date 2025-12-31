package repository

import (
	"database/internal/entity"

	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
)

/*
type UserRepository interface {
	GetById(id uuid.UUID) (entity.User, error)
	GetByEmail(email string) (entity.User, error)
	CreateUser(entity.User) (entity.User, error)
	MakeAdmin(id uuid.UUID) (entity.User, error)
	MakeSuperAdmin(id uuid.UUID) (entity.User, error)
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
		`UPDATE users
		set role='Admin'
		WHERE id = $1
		RETURNING first_name, second_name, email,avatar_url,role,token_provider
		`, id).Scan(
		&user.Id,
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
		`UPDATE users
		set role='Super_Admin'
		WHERE id = $1
		RETURNING first_name, second_name, email,avatar_url,role,token_provider
		`, id).Scan(
		&user.Id,
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
