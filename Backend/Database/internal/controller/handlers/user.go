package handlers

import (
	"context"
	"database/internal/usecase"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/gofrs/uuid"
	"golang.org/x/oauth2"
)

type userHandlers struct {
	userUseCase usecase.UserUseCase
}

func NewUserHandlers() *userHandlers {
	return &userHandlers{}
}

type authHandlers struct {
	CFG *oauth2.Config
}

func NewAuthHandlers() *authHandlers {
	return &authHandlers{}
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
	id := r.PathValue("UserID")

	if id == "" {
		HttpError(w, errors.New("Id is empty"), http.StatusBadRequest)
		return
	}

	uuid, err := uuid.FromString(id)

	if err != nil {
		HttpError(w, err, http.StatusBadRequest)
		return
	}

	user, err := h.userUseCase.GetById(uuid)

	if err != nil {
		if errors.As(err, usecase.ErrNotFound) {
			HttpError(w, err, http.StatusNotFound)
			return
		} else {
			HttpError(w, err, http.StatusInternalServerError)
			return
		}
	}

	b, err := json.MarshalIndent(user, "", "    ")

	if err != nil {
		HttpError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		fmt.Println("Error to write answer to http: ", err)
	}

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
	emailS := EmailDTO{}

	err := json.NewDecoder(r.Body).Decode(&emailS)

	if err != nil {
		HttpError(w, err, http.StatusBadRequest)
		return
	}

	user, err := h.userUseCase.GetByEmail(emailS.Email)

	if err != nil {
		if errors.As(err, usecase.ErrNotFound) {
			HttpError(w, err, http.StatusNotFound)
			return
		}
		HttpError(w, err, http.StatusInternalServerError)
		return
	}

	b, err := json.MarshalIndent(user, "", "    ")

	if err != nil {
		HttpError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		fmt.Println("Error to write answer to http: ", err)
	}
}

/*
pattern: /user
method:  POST
info:    JSON with user information

succeed:
  - status code:   201 Created
  - response body: JSON represent created user

failed:
  - status code:   401, 409, 500, ...
  - response body: JSON with error + time
*/
func (a *authHandlers) Regist(w http.ResponseWriter, r *http.Request) {
	url := a.CFG.AuthCodeURL("sTate123")
	http.Redirect(w, r, url, http.StatusFound)
}

func (a *authHandlers) RegistCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	if code == "" {
		HttpError(w, errors.New("no code from ya"), http.StatusInternalServerError)
		return
	}

	token, err := a.CFG.Exchange(context.Background(), code)

	if err != nil {
		HttpError(w, errors.New("token exchange failed"), http.StatusBadRequest)
		return
	}

	client := a.CFG.Client(context.Background(), token)

	resp, err := client.Get("https://login.yandex.ru/info?format=json")
	if err != nil {
		HttpError(w, errors.New("failed to get user data"), http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			HttpError(w, fmt.Errorf("Fail from Yandex server %w, status %s", err, re), http.StatusInternalServerError)
			return
		}

	}

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
