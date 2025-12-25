package handlers

import "net/http"

type enviromentHandlers struct{}

func NewEnviromentHandlers() *enviromentHandlers {
	return &enviromentHandlers{}
}

/*
pattern: /enviroment
method:  GET
info:    no

succeed:
  - status code:   200 OK
  - response body: JSON represent all enviroment info

failed:
  - status code:   500, ...
  - response body: JSON with error + time
*/

func (h *enviromentHandlers) GetAllEnviroment(w http.ResponseWriter, r *http.Request) {

}

/*
pattern: /enviroment/?type={type}
method:  GET
info:    Query

succeed:
  - status code:   200 OK
  - response body: JSON represent enviroment info with that type

failed:
  - status code:   400, 404, 500, ...
  - response body: JSON with error + time
*/

func (h *enviromentHandlers) GetEnviromentByType(w http.ResponseWriter, r *http.Request) {

}

/*
pattern: /enviroment/?id={id}
method:  GET
info:    Query

succeed:
  - status code:   200 OK
  - response body: JSON represent all enviroment info with that id

failed:
  - status code:   400, 404, 500, ...
  - response body: JSON with error + time
*/

func (h *enviromentHandlers) GetEnviromentById(w http.ResponseWriter, r *http.Request) {

}

/*
pattern: /enviroment
method:  POST
info:    JSON with data about enviroment

succeed:
  - status code:   201 Created
  - response body: JSON represent all enviroment info with that id

failed:
  - status code:   400, 409, 500, ...
  - response body: JSON with error + time
*/

func (h *enviromentHandlers) CreateEnviroment(w http.ResponseWriter, r *http.Request) {

}

/*
pattern: /enviroment/type
method:  POST
info:    JSON with type of enviroment

succeed:
  - status code:   201 Created
  - response body: JSON represent all enviroment info with that id

failed:
  - status code:   400, 409, 500, ...
  - response body: JSON with error + time
*/

func (h *enviromentHandlers) CreateTypeOfEnviroment(w http.ResponseWriter, r *http.Request) {

}

/*
pattern: /enviroment/{EnviromentId}
method:  PUT
info:    URL + json

succeed:
  - status code:   200 OK
  - response body: JSON represent all enviroment info with that id

failed:
  - status code:   400, 409, 500, ...
  - response body: JSON with error + time
*/

func (h *enviromentHandlers) EditEnviroment(w http.ResponseWriter, r *http.Request) {

}

/*
pattern: /enviroment/status/{EnviromentId}
method:  PUT
info:    URL + json

succeed:
  - status code:   200 OK
  - response body: JSON represent all enviroment info with that id

failed:
  - status code:   400, 409, 500, ...
  - response body: JSON with error + time
*/

func (h *enviromentHandlers) EditStatusEnviroment(w http.ResponseWriter, r *http.Request) {

}

/*
pattern: /enviroment/{EnviromentId}
method:  DELETE
info:    URL

succeed:
  - status code:   204 No content
  - response body: JSON represent all enviroment info with that id

failed:
  - status code:   400, 409, 500, ...
  - response body: JSON with error + time
*/

func (h *enviromentHandlers) DeleteEnviroment(w http.ResponseWriter, r *http.Request) {

}
