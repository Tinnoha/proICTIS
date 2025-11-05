package server

import (
	"context"
	"fmt"

	_ "github.com/jackc/pgx"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type databases struct {
	postgr *sqlx.DB
	redis  *redis.Client
}

func NewDatabases() *databases {
	connectStr := "user=postgres password=0000 dbname=postgres sslmode=disable"

	postgr, err := sqlx.Connect("postgres", connectStr)

	if err != nil {
		fmt.Println(err)
	}

	rediska := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	pong, err := rediska.Ping(context.Background()).Result()

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Подключение верно: ", pong)

	return &databases{
		postgr: postgr,
		redis:  rediska,
	}

}
