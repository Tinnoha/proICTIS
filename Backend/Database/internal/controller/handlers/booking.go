package handlers

import "net/http"

type bookingHandlers struct{}

func NewBookingHandlers() *bookingHandlers {
	return &bookingHandlers{}
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

func (h *bookingHandlers) GetAllBooking(w http.ResponseWriter, r *http.Request) {

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

func (h *bookingHandlers) GetBookingByUserId(w http.ResponseWriter, r *http.Request) {

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

func (h *bookingHandlers) GetBookingByEquipmentId(w http.ResponseWriter, r *http.Request) {

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

func (h *bookingHandlers) CreateBooking(w http.ResponseWriter, r *http.Request) {

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

func (h *bookingHandlers) AcceptOrCancelBooking(w http.ResponseWriter, r *http.Request) {

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

func (h *bookingHandlers) ReturnEquipment(w http.ResponseWriter, r *http.Request) {

}
