package handlers

import (
	"encoding/json"
	"net/http"
	"time"
)

type EmailDTO struct {
	Email string `json:"email"`
}

type errDTO struct {
	err  error     `json:"Error"`
	time time.Time `json:"Time"`
}

func HttpError(w http.ResponseWriter, err error, status int) {
	errdto := errDTO{
		err:  err,
		time: time.Now(),
	}

	b, err := json.MarshalIndent(errdto, "", "    ")

	if err != nil {
		panic(err)
	}

	http.Error(w, string(b), status)
}
