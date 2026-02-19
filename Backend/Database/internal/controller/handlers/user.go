package handlers

import (
	"context"
	"database/internal/entity"
	"database/internal/usecase"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
)

type UserHandlers struct {
	userUseCase usecase.UserUseCase
}

func NewUserHandlers(userUseCase usecase.UserUseCase) *UserHandlers {
	return &UserHandlers{userUseCase: userUseCase}
}

type AuthHandlers struct {
	CFG         *oauth2.Config
	userUseCase usecase.UserUseCase
}

func NewAuthHandlers(cfg *oauth2.Config, userUseCase usecase.UserUseCase) *AuthHandlers {
	return &AuthHandlers{CFG: cfg, userUseCase: userUseCase}
}

func (h *UserHandlers) GetAll(w http.ResponseWriter, r *http.Request) {
	user, err := h.userUseCase.GetAll()

	if err != nil {
		if errors.As(err, &usecase.ErrNotFound) {
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
func (h *UserHandlers) GetUserById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

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
func (h *UserHandlers) GetUserByEmail(w http.ResponseWriter, r *http.Request) {
	emailS := EmailDTO{}
	fmt.Println("Popla")

	err := json.NewDecoder(r.Body).Decode(&emailS)

	if err != nil {
		fmt.Println("WWOWOWOW")
		HttpError(w, err, http.StatusBadRequest)
		return
	}

	user, err := h.userUseCase.GetByEmail(emailS.Email)

	if err != nil {
		fmt.Println("Popchung")
		if errors.As(err, usecase.ErrNotFound) {
			HttpError(w, err, http.StatusNotFound)
			return
		}
		HttpError(w, err, http.StatusInternalServerError)
		return
	}

	b, err := json.MarshalIndent(user, "", "    ")

	if err != nil {
		fmt.Println("Popchunsssss")
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
func (a *AuthHandlers) Regist(w http.ResponseWriter, r *http.Request) {
	url := a.CFG.AuthCodeURL("sTate123")
	fmt.Println("url", url)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (a *AuthHandlers) RegistCallback(w http.ResponseWriter, r *http.Request) {
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
		body, _ := io.ReadAll(resp.Body)
		HttpError(w, fmt.Errorf("Fail from Yandex server %s", string(body)), http.StatusInternalServerError)
		return
	}

	yaUser := YandexUser{}
	if err := json.NewDecoder(resp.Body).Decode(&yaUser); err != nil {
		HttpError(w, fmt.Errorf("Failed to decode ya User with data: %w", err), http.StatusInternalServerError)
		return
	}

	fmt.Println("Avatar:", yaUser.Avatar, "First", yaUser.FirstName, "email", yaUser.Email, "last", yaUser.LastName)

	user := []entity.User{}
	yauser := entity.User{
		FirstName:  yaUser.FirstName,
		SecondName: yaUser.LastName,
		Email:      yaUser.Email,
		AvatarURL:  yaUser.Avatar,
	}

	user = append(user, yauser)

	_, err = a.userUseCase.CreateUser(user)
	if err != nil {
		HttpError(w, fmt.Errorf("Failed to save ya User with data: %w", err), http.StatusInternalServerError)
		return
	}
	url := "http://localhost:8080/"

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)

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
func (h *UserHandlers) MakeAdmin(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Admin")
	admins := AdminsDTO{}
	err := json.NewDecoder(r.Body).Decode(&admins)

	if err != nil {
		fmt.Println("Admin First")
		HttpError(w, err, http.StatusBadRequest)
		return
	}

	admin, err := h.userUseCase.MakeAdmin(admins.SecondAdmin)

	if err != nil {
		if errors.As(err, &usecase.ErrNotFound) {
			fmt.Println("Admin Second")
			HttpError(w, err, http.StatusNotFound)
			return
		}
		fmt.Println("Admin Third")
		HttpError(w, err, http.StatusInternalServerError)
		return
	}

	b, err := json.MarshalIndent(admin, "", "    ")

	fmt.Print("dECOEDEDEDEDED")

	if err != nil {
		fmt.Println("Admin Fourth")
		HttpError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		fmt.Println("Error to write answer to http: ", err)
	}
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
func (h *UserHandlers) MakeSuperAdmin(w http.ResponseWriter, r *http.Request) {
	admins := AdminsDTO{}
	err := json.NewDecoder(r.Body).Decode(&admins)

	if err != nil {
		HttpError(w, err, http.StatusBadRequest)
		return
	}

	admin, err := h.userUseCase.MakeSuperAdmin(admins.SecondAdmin)

	if err != nil {
		if errors.As(err, &usecase.ErrNotFound) {
			HttpError(w, err, http.StatusNotFound)
			return
		}
		HttpError(w, err, http.StatusInternalServerError)
		return
	}

	b, err := json.MarshalIndent(admin, "", "    ")

	if err != nil {
		HttpError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		fmt.Println("Error to write answer to http: ", err)
	}
}
