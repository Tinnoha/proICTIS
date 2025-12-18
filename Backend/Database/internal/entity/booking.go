package entity

import "time"

type Booking struct {
	UserId       int
	EnviromentId int
	BookStart    time.Time
	BookEnd      time.Time
	Status       string
}
