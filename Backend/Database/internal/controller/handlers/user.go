package handlers

import (
	"net/http"
)

type userHandlers struct{}

func NewUserHandlers() *userHandlers {
	return &userHandlers{}
}

/*
pattern: /user/{id}
method:  GET
info:    URL

succeed:
  - status code:   200 OK
  - response body: JSON represent user with that ID

failed:
  - status code:   400, 404, 500, ...
  - response body: JSON with error + time
*/
func (h *userHandlers) GetUserById(w http.ResponseWriter, r *http.Request) {

}

/*
pattern: /user/email
method:  GET
info:    JSON with email user

succeed:
  - status code:   200 OK
  - response body: JSON represent user with that email

failed:
  - status code:   400, 404, 500, ...
  - response body: JSON with error + time
*/
func (h *userHandlers) GetUserByEmail(w http.ResponseWriter, r *http.Request) {

}

/*
pattern: /user
method:  POST
info:    JSON with user information

succeed:
  - status code:   201 Created
  - response body: JSON represent created user

failed:
  - status code:   400, 409, 500, ...
  - response body: JSON with error + time
*/
func (h *userHandlers) CreateUser(w http.ResponseWriter, r *http.Request) {

}

/*
pattern: /user/admin
method:  PUT
info:    JSON with SuperAdmin user.id and new admin id

succeed:
  - status code:   200 OK
  - response body: JSON represent edited user

failed:
  - status code:   400, 404, 500, ...
  - response body: JSON with error + time
*/
func (h *userHandlers) MakeAdmin(w http.ResponseWriter, r *http.Request) {

}

/*
pattern: /user/superadmin
method:  PUT
info:    JSON with SuperAdmin user.id and new SuperAdmin id

succeed:
  - status code:   200 OK
  - response body: JSON represent edited user

failed:
  - status code:   400, 404, 500, ...
  - response body: JSON with error + time
*/
func (h *userHandlers) MakeSuperAdmin(w http.ResponseWriter, r *http.Request) {

}
