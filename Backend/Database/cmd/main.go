package main

import (
	"context"
	"database/internal/config"
	"database/internal/controller"
	"database/internal/controller/handlers"
	"database/internal/repository"
	"database/internal/usecase"
	"fmt"
	"time"
)

func main() {
	cfgRedis := config.RedisCondig{
		Addr:        "localhost:6379",
		Password:    "test1234",
		User:        "testuser",
		DB:          0,
		MaxRetries:  5,
		DialTimeout: 10 * time.Second,
		Timeout:     5 * time.Second,
	}
	db := repository.NewDatabase()
	red, err := config.NewRedisConfig(context.Background(), cfgRedis)
	if err != nil {
		panic(err)
	}

	err = repository.CreateTables(db)
	if err != nil {
		panic(err)
	}
	fmt.Println("We create db")

	bookingRepo := repository.NewBookingRepo(*db)
	equipmentRepo := repository.NewEquipmentRepo(db, red)
	userRepo := repository.NewUserRepo(*db)
	fmt.Println("We create repos")

	bookingUseCase := usecase.NewBooknigUseCase(userRepo, bookingRepo, equipmentRepo)
	equipmentUseCase := usecase.NewEquipmentUsecase(bookingRepo, equipmentRepo, userRepo)
	userUseCase := usecase.NewUserUseCase(userRepo)
	cfg := config.NewHernyaOauthConfig()
	if cfg == nil {
		panic("wow")
	}
	fmt.Println("We create usecases")

	config := cfg.ToOauth()

	bookingHandler := handlers.NewBookingHandlers(*bookingUseCase)
	eqiupmentHandler := handlers.NewEquipmentHandlers(*equipmentUseCase)
	userHandler := handlers.NewUserHandlers(*userUseCase)
	authHandler := handlers.NewAuthHandlers(&config, *userUseCase)
	fmt.Println("We create handlers")

	server := controller.NewHTTPServer(*authHandler, *userHandler, *eqiupmentHandler, *bookingHandler)

	fmt.Println("We create server")
	server.Run()
}
