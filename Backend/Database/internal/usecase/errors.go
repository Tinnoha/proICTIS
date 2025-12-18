package usecase

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound  = errors.New("Is not found")
	ErrIntenal   = errors.New("Something get wrong")
	ErrThisExist = errors.New("This is exists")
)

func ErrThisExists(that string, takenname string) error {
	return fmt.Errorf(ErrThisExist.Error(), ":", that, takenname)
}

func ErrInntenal(err error) error {
	return fmt.Errorf(ErrIntenal.Error(), err.Error())
}
