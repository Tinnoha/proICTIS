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
	"net/url"
	"strconv"
	"strings"

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
	admin := AdminDTO{}

	err := json.NewDecoder(r.Body).Decode(&admin)

	if err != nil {
		HttpError(w, err, http.StatusBadRequest)
		return
	}

	user, err := h.userUseCase.GetAll(admin.AdminId)

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

	err := json.NewDecoder(r.Body).Decode(&emailS)

	if err != nil {

		HttpError(w, err, http.StatusBadRequest)
		return
	}

	user, err := h.userUseCase.GetByEmail(emailS.AdminId, emailS.Email)

	if err != nil {

		if errors.As(err, &usecase.ErrNotFound) {
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
func (a *AuthHandlers) Regist(w http.ResponseWriter, r *http.Request) {
	url := a.CFG.AuthCodeURL("sTate123")

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
		HttpError(w, fmt.Errorf("Failed to decode ya User with  %w", err), http.StatusInternalServerError)
		return
	}

	yauser := entity.User{
		FirstName:  yaUser.FirstName,
		SecondName: yaUser.LastName,
		Email:      yaUser.Email,
		AvatarURL:  yaUser.Avatar,
	}

	createdUsers, err := a.userUseCase.CreateUser([]entity.User{yauser})

	var createdUser entity.User

	if err != nil {
		if strings.Contains(err.Error(), "This is exists") {
			createdUser, err = a.userUseCase.GetByEmailNoAuth(yaUser.Email)
			if err != nil {
				HttpError(w, fmt.Errorf("Failed to get existing user: %w", err), http.StatusInternalServerError)
				return
			}
		} else {
			HttpError(w, fmt.Errorf("Failed to save ya User with  %w", err), http.StatusInternalServerError)
			return
		}
	} else if len(createdUsers) > 0 {
		createdUser = createdUsers[0]
	} else {
		HttpError(w, errors.New("user not created"), http.StatusInternalServerError)
		return
	}

	frontendToken := createdUser.Id.String()
	frontendBaseURL := "http://localhost:8080"

	frontendURL := fmt.Sprintf("%s?token=%s&email=%s&firstName=%s&lastName=%s&avatar=%s&role=%s",
		frontendBaseURL,
		frontendToken,
		url.QueryEscape(yaUser.Email),
		url.QueryEscape(yaUser.FirstName),
		url.QueryEscape(yaUser.LastName),
		url.QueryEscape(yaUser.Avatar),
		url.QueryEscape(string(createdUser.Role)),
	)

	http.Redirect(w, r, frontendURL, http.StatusTemporaryRedirect)

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

	admins := AdminsDTO{}
	err := json.NewDecoder(r.Body).Decode(&admins)

	if err != nil {

		HttpError(w, err, http.StatusBadRequest)
		return
	}

	admin, err := h.userUseCase.MakeAdmin(admins.FirstAdmin, admins.SecondAdmin)

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

	admin, err := h.userUseCase.MakeSuperAdmin(admins.FirstAdmin, admins.SecondAdmin)

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

func (h *UserHandlers) CreateLink(w http.ResponseWriter, r *http.Request) {
	user := UserIdDTO{}

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		HttpError(w, err, http.StatusBadRequest)
		return
	}

	link, err := h.userUseCase.VKCreateLink(user.Id)

	if err != nil {
		HttpError(w, err, http.StatusBadRequest)
		return
	}

	Link := LinkDTO{Link: link}

	b, err := json.MarshalIndent(Link, "", "    ")

	if err != nil {

		HttpError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		fmt.Println("Error to write answer to http: ", err)
	}
}

func (h *UserHandlers) ConnectVk(w http.ResponseWriter, r *http.Request) {
	vk := VkConnection{}

	err := json.NewDecoder(r.Body).Decode(&vk)
	if err != nil {
		HttpError(w, err, http.StatusBadRequest)
		return
	}

	err = h.userUseCase.VKConnect(vk.VkId, vk.Token)

	if err != nil {
		HttpError(w, err, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *UserHandlers) GetUserByVkId(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("vkid")
	if id == "" {
		HttpError(w, errors.New("vkid query param is required"), http.StatusBadRequest)
		return
	}

	vkID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		HttpError(w, errors.New("invalid vk_id format"), http.StatusBadRequest)
		return
	}

	user, err := h.userUseCase.GetByVkId(int(vkID))

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
