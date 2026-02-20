package controller

import (
	"database/internal/controller/handlers"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type HTTPServer struct {
	authHandler      handlers.AuthHandlers
	userHandler      handlers.UserHandlers
	equipmentHandler handlers.EquipmentHandlers
	bookHandler      handlers.BookingHandlers
}

func NewHTTPServer(
	authHandler handlers.AuthHandlers,
	userHandler handlers.UserHandlers,
	equipmentHandler handlers.EquipmentHandlers,
	bookHandler handlers.BookingHandlers,
) HTTPServer {
	return HTTPServer{
		authHandler:      authHandler,
		userHandler:      userHandler,
		equipmentHandler: equipmentHandler,
		bookHandler:      bookHandler,
	}
}

func (s *HTTPServer) Run() {
	fmt.Println("We start HTTP SERVER")
	router := mux.NewRouter()

	router.Path("/Regist").HandlerFunc(s.authHandler.Regist)
	router.Path("/callback").HandlerFunc(s.authHandler.RegistCallback)
	//WORK

	router.Path("/User/email").Methods("GET").HandlerFunc(s.userHandler.GetUserByEmail)
	router.Path("/User/admin").Methods("PATCH").HandlerFunc(s.userHandler.MakeAdmin)
	router.Path("/User/SuperAdmin").Methods("PATCH").HandlerFunc(s.userHandler.MakeSuperAdmin)
	router.Path("/User/{id}").Methods("GET").HandlerFunc(s.userHandler.GetUserById)
	router.Path("/User").Methods("GET").HandlerFunc(s.userHandler.GetAll)
	//WORK

	router.Path("/Equipment").Methods("GET").HandlerFunc(s.equipmentHandler.GetAllEquipment)
	router.Path("/Equipment/type/{type}").Methods("GET").HandlerFunc(s.equipmentHandler.GetEquipmentByType)
	router.Path("/Equipment/id/{id}").Methods("GET").HandlerFunc(s.equipmentHandler.GetEquipmentById)
	router.Path("/Equipment").Methods("POST").HandlerFunc(s.equipmentHandler.CreateEquipment)
	router.Path("/Equipment/{EquipmentId}").Methods("PATCH").HandlerFunc(s.equipmentHandler.EditEquipment)
	router.Path("/Equipment/status/{EquipmentId}").Methods("PATCH").HandlerFunc(s.equipmentHandler.EditStatusEquipment)
	router.Path("/Equipment/{EquipmentId}").Methods("DELETE").HandlerFunc(s.equipmentHandler.DeleteEquipment)
	// WORK

	router.Path("/Types").Methods("GET")
	router.Path("/Types").Methods("POST")
	router.Path("/Types/{id}").Methods("PATCH")
	router.Path("/Types/{id}").Methods("DELETE")

	router.Path("/Booking").Methods("GET").HandlerFunc(s.bookHandler.GetAllBooking)
	router.Path("/Booking/user/{id}").Methods("GET").HandlerFunc(s.bookHandler.GetAllBooking)
	router.Path("/Booking/equipment/{id}").Methods("GET").HandlerFunc(s.bookHandler.GetBookingByEquipmentId)
	router.Path("/Booking").Methods("POST").HandlerFunc(s.bookHandler.CreateBooking)
	router.Path("/Booking/{id}").Methods("PUT").HandlerFunc(s.bookHandler.AcceptOrCancelBooking)
	router.Path("/Booking/return/{id}").Methods("PUT").HandlerFunc(s.bookHandler.ReturnEquipment)
	// WORK
	http.ListenAndServe(":8080", router)
	fmt.Println("We finish HTTP SERVER")
}
