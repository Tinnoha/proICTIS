package main

import (
	"database/internal/config"
	"database/internal/controller"
	"database/internal/controller/handlers"
	"database/internal/repository"
	"database/internal/usecase"
	"fmt"
)

func main() {
	db := repository.NewDatabase()
	err := repository.CreateTables(db)
	if err != nil {
		panic(err)
	}
	fmt.Println("We create db")

	bookingRepo := repository.NewBookingRepo(*db)
	equipmentRepo := repository.NewУquipmentRepo(*db)
	userRepo := repository.NewUserRepo(*db)
	fmt.Println("We create repos")

	bookingUseCase := usecase.NewBooknigUseCase(userRepo, bookingRepo, equipmentRepo)
	equipmentUseCase := usecase.NewEquipmentUsecase(bookingRepo, equipmentRepo)
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
