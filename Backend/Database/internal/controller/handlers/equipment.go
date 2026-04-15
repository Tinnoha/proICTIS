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
	fileUseCase      usecase.FileStorageUsecase
}

func NewEquipmentHandlers(
	equipmentUseCase usecase.EquipmentUsecase,
	fileUsecase usecase.FileStorageUsecase,
) *EquipmentHandlers {
	return &EquipmentHandlers{
		equipmentUseCase: equipmentUseCase,
		fileUseCase:      fileUsecase,
	}
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
		if errors.As(err, &usecase.ErrThisExist) {
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

	editor := EquipmentEdit{}
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

	result, err := h.equipmentUseCase.Edit(editor.AdminId, uuid, editor.Equipment)

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

	err = h.equipmentUseCase.SetActive(active.AdminId, uuid, active.Active)
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
	admin := AdminDTO{}

	err = json.NewDecoder(r.Body).Decode(&admin)

	if err != nil {
		HttpError(w, err, http.StatusBadRequest)
		return
	}

	err = h.equipmentUseCase.Delete(admin.AdminId, uuid)
	if err != nil {
		if errors.As(err, &usecase.ErrNotFound) {
			HttpError(w, err, http.StatusNotFound)
			return
		} else {
			HttpError(w, err, http.StatusInternalServerError)
			return
		}
	}
}

func (h *EquipmentHandlers) GetTypes(w http.ResponseWriter, r *http.Request) {
	types, err := h.equipmentUseCase.GetTypes()

	if err != nil {
		HttpError(w, err, http.StatusInternalServerError)
		return
	}

	b, err := json.MarshalIndent(types, "", "    ")

	if err != nil {
		HttpError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		fmt.Println("We cannot to give answer")
	}
}

func (h *EquipmentHandlers) AddType(w http.ResponseWriter, r *http.Request) {
	types := TypesDTO{}

	err := json.NewDecoder(r.Body).Decode(&types)

	if err != nil {
		HttpError(w, err, http.StatusBadRequest)
		return
	}

	result, err := h.equipmentUseCase.AddTypes(types.AdminId, types.Types)

	if err != nil {
		HttpError(w, err, http.StatusInternalServerError)
		return
	}

	b, err := json.MarshalIndent(result, "", "    ")

	if err != nil {
		HttpError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write(b); err != nil {
		fmt.Println("We cannot to give answer")
	}
}

func (h *EquipmentHandlers) EditType(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, ok := vars["id"]

	if !ok {
		HttpError(w, errors.New("WITHOUT ID"), http.StatusBadRequest)
		return
	}

	uuid, err := uuid.FromString(id)

	if err != nil {
		HttpError(w, errors.New("id is not uuid"), http.StatusBadRequest)
		return
	}

	admin := TypeDTO{}
	err = json.NewDecoder(r.Body).Decode(&admin)

	if err != nil {
		HttpError(w, err, http.StatusBadRequest)
		return
	}

	admin.Type.Id = uuid

	ttype, err := h.equipmentUseCase.EditType(admin.AdminId, admin.Type)

	if err != nil {
		HttpError(w, err, http.StatusBadRequest)
		return
	}

	b, err := json.MarshalIndent(ttype, "", "    ")

	if err != nil {
		HttpError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write(b); err != nil {
		fmt.Println("We cannot to give answer")
	}
}

func (h *EquipmentHandlers) DeleteType(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, ok := vars["id"]

	if !ok {
		HttpError(w, errors.New("WITHOUT ID"), http.StatusBadRequest)
		return
	}

	uuid, err := uuid.FromString(id)

	if err != nil {
		HttpError(w, errors.New("id is not uuid"), http.StatusBadRequest)
		return
	}

	admin := AdminDTO{}
	err = json.NewDecoder(r.Body).Decode(&admin)

	if err != nil {
		HttpError(w, err, http.StatusBadRequest)
		return
	}

	err = h.equipmentUseCase.DeleteType(admin.AdminId, uuid)

	if err != nil {
		HttpError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *EquipmentHandlers) UploadImage(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		HttpError(w, err, http.StatusBadRequest)
		return
	}

	adminIDStr := r.FormValue("admin_id")
	adminID, err := uuid.FromString(adminIDStr)
	if err != nil {
		HttpError(w, errors.New("invalid admin_id"), http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		HttpError(w, err, http.StatusBadRequest)
		return
	}
	defer file.Close()

	url, err := h.fileUseCase.Save(adminID, file, header.Filename)
	if err != nil {
		HttpError(w, err, http.StatusInternalServerError)
		return
	}

	// 6. Отдаём ответ
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(UploadDTO{URL: url})
}
