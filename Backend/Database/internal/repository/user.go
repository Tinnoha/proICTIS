package repository

import (
	"database/internal/entity"
	"errors"
	"fmt"
	"time"

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
	_, err := u.db.Exec(`
	INSERT INTO proICTIS_user 
	(id,
	first_name, 
	second_name, 
	email, 
	avatar_url,
	role, 
	token_provider) 
	VALUES($1,$2,$3,$4,$5,$6,$7)`,
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

func (u *userRepo) ChangeRole(id uuid.UUID, role string) (entity.User, error) {
	user := entity.User{Id: id}

	err := u.db.QueryRow(
		`UPDATE PROICTIS_user
		set role=$1
		WHERE id = $2
		RETURNING first_name, second_name, email,avatar_url,role,token_provider
		`, role, id).Scan(
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

// func (u *userRepo) MakeAdmin(id uuid.UUID) (entity.User, error) {
// 	user := entity.User{Id: id}

// 	err := u.db.QueryRow(
// 		`UPDATE PROICTIS_user
// 		set role='Admin'
// 		WHERE id = $1
// 		RETURNING first_name, second_name, email,avatar_url,role,token_provider
// 		`, id).Scan(
// 		&user.FirstName,
// 		&user.SecondName,
// 		&user.Email,
// 		&user.AvatarURL,
// 		&user.Role,
// 		&user.TokenProvider,
// 	)

// 	if err != nil {
// 		return entity.User{}, err
// 	}

// 	return user, nil
// }

// func (u *userRepo) MakeSuperAdmin(id uuid.UUID) (entity.User, error) {
// 	user := entity.User{Id: id}

// 	err := u.db.QueryRow(
// 		`UPDATE PROICTIS_user
// 		set role='Super_Admin'
// 		WHERE id = $1
// 		RETURNING first_name, second_name, email,avatar_url,role,token_provider
// 		`, id).Scan(
// 		&user.FirstName,
// 		&user.SecondName,
// 		&user.Email,
// 		&user.AvatarURL,
// 		&user.Role,
// 		&user.TokenProvider,
// 	)

// 	if err != nil {
// 		return entity.User{}, err
// 	}

// 	return user, nil
// }

func (u *userRepo) IsAdmin(id uuid.UUID) (bool, error) {
	role := ""
	err := u.db.QueryRow(`SELECT role FROM proICTIS_user where id = $1`, id).Scan(&role)

	if err != nil {
		return false, err
	}

	if role == "Admin" || role == "Super_Admin" {
		return true, nil
	} else {
		return false, nil
	}
}

func (u *userRepo) IsSuperAdmin(id uuid.UUID) (bool, error) {
	role := ""
	err := u.db.QueryRow(`SELECT role FROM proICTIS_user where id = $1`, id).Scan(&role)

	if err != nil {
		return false, err
	}

	if role == "Super_Admin" {
		return true, nil
	} else {
		return false, nil
	}
}

func (u *userRepo) CreateToken(
	userId uuid.UUID,
	token string,
	timeExpire time.Time,
) (string, error) {
	fmt.Println("We in Create Link repo")
	sqlRequest := `
		UPDATE proICTIS_user 
		SET vk_token = $2,
			time_token = $3
		WHERE id = $1;
		`

	pat, err := u.db.Exec(sqlRequest, userId, token, timeExpire)

	if c, _ := pat.RowsAffected(); c == 0 {
		return "", errors.New("user_id is not exists")
	}

	if err != nil {
		return "", err
	}

	token = fmt.Sprintf("https://vk.com/write-237660555?ref=%s", token)

	fmt.Println("We finish cteate link in repo")

	return token, nil
}

func (u *userRepo) ConnectVK(vk_token uuid.UUID, vkId int) error {
	sqlRequest := `
	SELECT time_token
	FROM proICTIS_user
	WHERE vk_token = $1
	`

	row := u.db.QueryRow(sqlRequest, vk_token)

	if row.Err() != nil {
		return row.Err()
	}

	tokenTime := time.Time{}

	err := row.Scan(&tokenTime)

	if err != nil {
		return err
	}

	if time.Now().After(tokenTime) {
		return errors.New("time expired, try again")
	}

	sqlRequest = `
	UPDATE proICTIS_user
	SET vk_token = NULL,
		time_token = NULL,
		vk_id = $1
	WHERE vk_token = $2 AND vk_id IS NULL; 
	`

	pat, err := u.db.Exec(sqlRequest, vkId, vk_token)
	if err != nil {
		return err
	}

	if c, _ := pat.RowsAffected(); c == 0 {
		return errors.New("user_id is not exists")
	}

	return nil
}

func (u *userRepo) GetByVkId(VkId int) (entity.User, error) {
	vasya := entity.User{}
	err := u.db.QueryRow(`Select id, first_name, second_name, email, avatar_url,role, token_provider 
	from proICTIS_user 
	where vk_id = $1`, VkId).Scan(
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
