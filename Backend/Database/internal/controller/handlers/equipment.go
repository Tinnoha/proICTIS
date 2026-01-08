package handlers

import "net/http"

type equipmentHandlers struct{}

func NewEquipmentHandlers() *equipmentHandlers {
	return &equipmentHandlers{}
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

func (h *equipmentHandlers) GetAllEquipment(w http.ResponseWriter, r *http.Request) {

}

/*
pattern: /equipment/?type={type}
method:  GET
info:    Query

succeed:
  - status code:   200 OK
  - response body: JSON represent equipment info with that type

failed:
  - status code:   400, 404, 500, ...
  - response body: JSON with error + time
*/

func (h *equipmentHandlers) GetEquipmentByType(w http.ResponseWriter, r *http.Request) {

}

/*
pattern: /equipment/?id={id}
method:  GET
info:    Query

succeed:
  - status code:   200 OK
  - response body: JSON represent all equipment info with that id

failed:
  - status code:   400, 404, 500, ...
  - response body: JSON with error + time
*/

func (h *equipmentHandlers) GetEquipmentById(w http.ResponseWriter, r *http.Request) {

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

func (h *equipmentHandlers) CreateEquipment(w http.ResponseWriter, r *http.Request) {

}

/*
pattern: /equipment/type
method:  POST
info:    JSON with type of equipment

succeed:
  - status code:   201 Created
  - response body: JSON represent all equipment info with that id

failed:
  - status code:   400, 409, 500, ...
  - response body: JSON with error + time
*/

func (h *equipmentHandlers) CreateTypeOfEquipment(w http.ResponseWriter, r *http.Request) {

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

func (h *equipmentHandlers) EditEquipment(w http.ResponseWriter, r *http.Request) {

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

func (h *equipmentHandlers) EditStatusEquipment(w http.ResponseWriter, r *http.Request) {

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

func (h *equipmentHandlers) DeleteEquipment(w http.ResponseWriter, r *http.Request) {

}
