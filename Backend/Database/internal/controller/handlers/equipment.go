package handlers

import (
	"database/internal/entity"
	"database/internal/usecase"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
)

type EquipmentHandlers struct {
	equipmentUseCase usecase.EquipmentUsecase
}

func NewEquipmentHandlers(equipmentUseCase usecase.EquipmentUsecase) *EquipmentHandlers {
	return &EquipmentHandlers{equipmentUseCase: equipmentUseCase}
}

/*
pattern: /equipment
method:  GET
info:    no

succeed:
  - status code:   200 OK
  - response body: JSON represent all equipment info

failed:
  - status code:   500, ...
  - response body: JSON with error + time
*/

func (h *EquipmentHandlers) GetAllEquipment(w http.ResponseWriter, r *http.Request) {
	equipment, err := h.equipmentUseCase.GetAll()

	if err != nil {
		HttpError(w, err, http.StatusInternalServerError)
		return
	}

	b, err := json.MarshalIndent(equipment, "", "    ")
	if err != nil {
		HttpError(w, err, http.StatusInternalServerError)
		return
	}

	if _, err = w.Write(b); err != nil {
		fmt.Println("Error to write info :(")
	}

}

/*
pattern: /equipment/{type}
method:  GET
info:    Query

succeed:
  - status code:   200 OK
  - response body: JSON represent equipment info with that type

failed:
  - status code:   400, 404, 500, ...
  - response body: JSON with error + time
*/

func (h *EquipmentHandlers) GetEquipmentByType(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	Type := vars["type"]
	if Type == "" {
		HttpError(w, errors.New("Invalid data in query"), http.StatusBadRequest)
		return
	}

	enviroments, err := h.equipmentUseCase.GetByType(Type)

	if err != nil {
		if errors.As(err, usecase.ErrNotFound) {
			HttpError(w, err, http.StatusNotFound)
			return
		} else {
			HttpError(w, err, http.StatusInternalServerError)
			return
		}
	}

	b, err := json.MarshalIndent(enviroments, "", "    ")
	if err != nil {
		HttpError(w, err, http.StatusInternalServerError)
		return
	}

	if _, err = w.Write(b); err != nil {
		fmt.Println("Error to write info :(")
	}
}

/*
pattern: /equipment/{id}
method:  GET
info:    Query

succeed:
  - status code:   200 OK
  - response body: JSON represent all equipment info with that id

failed:
  - status code:   400, 404, 500, ...
  - response body: JSON with error + time
*/

func (h *EquipmentHandlers) GetEquipmentById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		HttpError(w, errors.New("Invalid data in query"), http.StatusBadRequest)
		return
	}

	uuid, err := uuid.FromString(id)

	if err != nil {
		HttpError(w, err, http.StatusBadRequest)
		return
	}

	equipment, err := h.equipmentUseCase.GetById(uuid)
	if err != nil {
		if errors.As(err, usecase.ErrNotFound) {
			HttpError(w, err, http.StatusNotFound)
			return
		} else {
			HttpError(w, err, http.StatusInternalServerError)
			return
		}
	}

	b, err := json.MarshalIndent(equipment, "", "    ")
	if err != nil {
		HttpError(w, err, http.StatusInternalServerError)
		return
	}

	if _, err = w.Write(b); err != nil {
		fmt.Println("Error to write info :(")
	}
}

/*
pattern: /equipment
method:  POST
info:    JSON with data about equipment

succeed:
  - status code:   201 Created
  - response body: JSON represent all equipment info with that id

failed:
  - status code:   400, 409, 500, ...
  - response body: JSON with error + time
*/

func (h *EquipmentHandlers) CreateEquipment(w http.ResponseWriter, r *http.Request) {
	equipments := entity.Equipments{}
	err := json.NewDecoder(r.Body).Decode(&equipments)

	if err != nil {
		HttpError(w, err, http.StatusBadRequest)
		return
	}

	result, err := h.equipmentUseCase.Add(equipments)
	if err != nil {
		if errors.As(err, usecase.ErrThisExist) {
			HttpError(w, err, http.StatusBadRequest)
			return
		} else {
			HttpError(w, err, http.StatusInternalServerError)
			return
		}
	}

	b, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		HttpError(w, err, http.StatusInternalServerError)
		return
	}

	if _, err = w.Write(b); err != nil {
		fmt.Println("Error to write info :(")
	}
}

/*
pattern: /equipment/{EquipmentId}
method:  PUT
info:    URL + json

succeed:
  - status code:   200 OK
  - response body: JSON represent all equipment info with that id

failed:
  - status code:   400, 409, 500, ...
  - response body: JSON with error + time
*/

func (h *EquipmentHandlers) EditEquipment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["EquipmentId"]

	print("id", id)
	if id == "" {
		HttpError(w, errors.New("NO ID"), http.StatusBadRequest)
		return
	}

	editor := entity.Equipment{}
	err := json.NewDecoder(r.Body).Decode(&editor)

	if err != nil {
		HttpError(w, err, http.StatusBadRequest)
		return
	}

	uuid, err := uuid.FromString(id)

	if err != nil {
		HttpError(w, err, http.StatusBadRequest)
		return
	}

	result, err := h.equipmentUseCase.Edit(uuid, editor)

	if err != nil {
		if errors.As(err, usecase.ErrNotFound) {
			HttpError(w, err, http.StatusNotFound)
			return
		} else {
			HttpError(w, err, http.StatusInternalServerError)
			return
		}
	}

	b, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		HttpError(w, err, http.StatusInternalServerError)
		return
	}

	if _, err = w.Write(b); err != nil {
		fmt.Println("Error to write info :(")
	}
}

/*
pattern: /equipment/status/{EquipmentId}
method:  PUT
info:    URL + json

succeed:
  - status code:   200 OK
  - response body: JSON represent all equipment info with that id

failed:
  - status code:   400, 409, 500, ...
  - response body: JSON with error + time
*/

func (h *EquipmentHandlers) EditStatusEquipment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["EquipmentId"]
	if id == "" {
		HttpError(w, errors.New("NO ID"), http.StatusBadRequest)
		return
	}

	uuid, err := uuid.FromString(id)

	if err != nil {
		HttpError(w, err, http.StatusBadRequest)
		return
	}

	active := BoolDTO{}
	err = json.NewDecoder(r.Body).Decode(&active)

	if err != nil {
		HttpError(w, err, http.StatusBadRequest)
		return
	}

	err = h.equipmentUseCase.SetActive(uuid, active.Active)
	if err != nil {
		if errors.As(err, usecase.ErrNotFound) {
			HttpError(w, err, http.StatusNotFound)
			return
		} else {
			HttpError(w, err, http.StatusInternalServerError)
			return
		}
	}
}

/*
pattern: /equipment/{EquipmentId}
method:  DELETE
info:    URL

succeed:
  - status code:   204 No content
  - response body: JSON represent all equipment info with that id

failed:
  - status code:   400, 409, 500, ...
  - response body: JSON with error + time
*/

func (h *EquipmentHandlers) DeleteEquipment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["EquipmentId"]
	if id == "" {
		HttpError(w, errors.New("NO ID"), http.StatusBadRequest)
		return
	}

	uuid, err := uuid.FromString(id)

	if err != nil {
		HttpError(w, err, http.StatusBadRequest)
		return
	}

	err = h.equipmentUseCase.Delete(uuid)
	if err != nil {
		if errors.As(err, usecase.ErrNotFound) {
			HttpError(w, err, http.StatusNotFound)
			return
		} else {
			HttpError(w, err, http.StatusInternalServerError)
			return
		}
	}
}
