package repository

import (
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func NewDatabase() *sqlx.DB {
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}

	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "5432"
	}

	user := os.Getenv("DB_USER")
	if user == "" {
		user = "postgres"
	}

	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		password = "0000"
	}

	dbname := os.Getenv("DB_NAME")
	if dbname == "" {
		dbname = "5432"
	}

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sqlx.Connect("pgx", connStr)
	if err != nil {
		fmt.Println("Error with connection to DB:", err)
		return nil
	}

	return db
}

func CreateTables(db *sqlx.DB) error {
	err := db.Ping()

	if err != nil {
		fmt.Println("Error with connetcion with db")
		return err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS proICTIS_user(
		id UUID PRIMARY KEY,
		first_name varchar(255),
		second_name varchar(255),
		email varchar(255) UNIQUE,
		avatar_URL varchar(255),
		role varchar(255),
		token_provider integer
		)
	`)

	if err != nil {
		fmt.Println("Error with creating table proICTIS_user")
		return err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS proICTIS_type_of_enviroment(
		id UUID PRIMARY KEY,
		name varchar(255) NOT NULL UNIQUE
		)
	`)

	if err != nil {
		fmt.Println("Error with creating table proICTIS_type_of_enviroment")
		return err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS proICTIS_enviroment(
		id UUID PRIMARY KEY,
		name varchar(255),
		description TEXT,
		photo_url varchar(255),
		type_id UUID NOT NULL REFERENCES proICTIS_type_of_enviroment(id) ON DELETE RESTRICT,
		auditory varchar(255),
		is_active boolean
		)
	`)

	if err != nil {
		fmt.Println("Error with creating table proICTIS_enviroment")
		return err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS proICTIS_booking(
		id UUID PRIMARY KEY,
		user_id UUID not null references proICTIS_user(id) on delete cascade,
		enviroment_id UUID NOT NULL REFERENCES proICTIS_enviroment(id) on delete cascade,
		book_start TIMESTAMPTZ,
		book_end TIMESTAMPTZ,
		status varchar(255)
		)
	`)

	if err != nil {
		fmt.Println("Error with creating table proICTIS_booking")
		return err
	}

	return nil

}
