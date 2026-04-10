package handlers

import (
	"database/internal/usecase"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
)

type BookingHandlers struct {
	bookingUsecase usecase.BookingUseCase
}

func NewBookingHandlers(bookingUsecase usecase.BookingUseCase) *BookingHandlers {
	return &BookingHandlers{bookingUsecase: bookingUsecase}
}

/*
pattern: /booking
method:  GET
info:    no

succeed:
  - status code:   200 OK
  - response body: JSON represent all bookings info

failed:
  - status code:   500, ...
  - response body: JSON with error + time
*/

func (h *BookingHandlers) GetAllBooking(w http.ResponseWriter, r *http.Request) {
	rents, err := h.bookingUsecase.GetAllBooks()
	if err != nil {
		HttpError(w, err, http.StatusInternalServerError)
		return
	}

	b, err := json.MarshalIndent(rents, "", "    ")

	if err != nil {
		HttpError(w, err, http.StatusInternalServerError)
		return
	}

	if _, err = w.Write(b); err != nil {
		fmt.Println("Error to write answer")
	}
}

/*
pattern: /booking/?UserId={UserId}
method:  GET
info:    Query

succeed:
  - status code:   200 OK
  - response body: JSON represent all bookings info

failed:
  - status code:   400, 404, 500, ...
  - response body: JSON with error + time
*/

func (h *BookingHandlers) GetBookingByUserId(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		HttpError(w, errors.New("No query params !"), http.StatusBadRequest)
		return
	}

	uuid, err := uuid.FromString(id)
	if err != nil {
		HttpError(w, err, http.StatusBadRequest)
		return
	}

	rents, err := h.bookingUsecase.GetBooksByUserId(uuid)

	if err != nil {
		HttpError(w, err, http.StatusBadRequest)
		return
	}

	b, err := json.MarshalIndent(rents, "", "    ")

	if err != nil {
		HttpError(w, err, http.StatusInternalServerError)
		return
	}

	if _, err = w.Write(b); err != nil {
		fmt.Println("Error to write answer")
	}
}

/*
pattern: /booking/?EquipmentId={EnviromtId}
method:  GET
info:    Query

succeed:
  - status code:   200 OK
  - response body: JSON represent all bookings info

failed:
  - status code:   400, 404, 500, ...
  - response body: JSON with error + time
*/

func (h *BookingHandlers) GetBookingByEquipmentId(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		HttpError(w, errors.New("No query params !"), http.StatusBadRequest)
		return
	}

	uuid, err := uuid.FromString(id)
	if err != nil {
		HttpError(w, err, http.StatusBadRequest)
		return
	}

	rents, err := h.bookingUsecase.GetBooksByEquipmentId(uuid)

	if err != nil {
		HttpError(w, err, http.StatusBadRequest)
		return
	}

	b, err := json.MarshalIndent(rents, "", "    ")

	if err != nil {
		HttpError(w, err, http.StatusInternalServerError)
		return
	}

	if _, err = w.Write(b); err != nil {
		fmt.Println("Error to write answer")
	}
}

/*
pattern: /booking
method:  Post
info:    json

succeed:
  - status code:   201 Created
  - response body: JSON represent all bookings info

failed:
  - status code:   400, 409, 500, ...
  - response body: JSON with error + time
*/

func (h *BookingHandlers) CreateBooking(w http.ResponseWriter, r *http.Request) {
	bok := BookDTO{}
	err := json.NewDecoder(r.Body).Decode(&bok)

	if err != nil {
		print(1)
		HttpError(w, err, http.StatusBadRequest)
		return
	}

	booking, err := h.bookingUsecase.Book(bok.UserId, bok.EnviromtId, bok.Start, bok.End)
	if err != nil {
		print(2)
		HttpError(w, err, http.StatusBadRequest)
		return
	}

	b, err := json.MarshalIndent(booking, "", "    ")

	if err != nil {
		HttpError(w, err, http.StatusInternalServerError)
		return
	}

	if _, err = w.Write(b); err != nil {
		fmt.Println("Error to write answer")
	}
}

/*
pattern: /booking/{id}
method:  PUT
info:    json

succeed:
  - status code:   201 Created
  - response body: JSON represent all bookings info

failed:
  - status code:   400, 404, 500, ...
  - response body: JSON with error + time
*/

func (h *BookingHandlers) AcceptOrCancelBooking(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		fmt.Println("Wpw1")
		HttpError(w, errors.New("No query params !"), http.StatusBadRequest)
		return
	}

	uuid, err := uuid.FromString(id)
	if err != nil {
		fmt.Println("Wpww2")
		HttpError(w, err, http.StatusBadRequest)
		return
	}

	status := StatusDTO{}
	err = json.NewDecoder(r.Body).Decode(&status)

	if err != nil {
		fmt.Println("Wpwwe3")
		HttpError(w, err, http.StatusBadRequest)
		return
	}

	res, err := h.bookingUsecase.EditStatusBooking(status.AdminId, uuid, status.Status)

	if err != nil {
		fmt.Println("Wpwwe65")
		HttpError(w, err, http.StatusBadRequest)
		return
	}

	b, err := json.MarshalIndent(res, "", "    ")

	if err != nil {
		HttpError(w, err, http.StatusInternalServerError)
		return
	}

	if _, err = w.Write(b); err != nil {
		fmt.Println("Error to write answer")
	}
}

/*
pattern: /booking/{id}
method:  PUT
info:    json

succeed:
  - status code:   201 Created
  - response body: JSON represent all bookings info

failed:
  - status code:   400, 404, 500, ...
  - response body: JSON with error + time
*/

func (h *BookingHandlers) ReturnEquipment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		HttpError(w, errors.New("No query params !"), http.StatusBadRequest)
		return
	}

	uuid, err := uuid.FromString(id)
	if err != nil {
		HttpError(w, err, http.StatusBadRequest)
		return
	}

	admin := AdminDTO{}

	err = json.NewDecoder(r.Body).Decode(&admin)

	if err != nil {
		HttpError(w, err, http.StatusBadRequest)
		return
	}

	err = h.bookingUsecase.DeleteBooking(admin.AdminId, uuid)
	if err != nil {
		HttpError(w, err, http.StatusBadRequest)
		return
	}

}
