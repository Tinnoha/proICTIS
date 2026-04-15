package usecase

import (
	"database/internal/repository"
	"errors"
	"io"
	"path/filepath"

	"github.com/gofrs/uuid"
)

type FileStorage interface {
	Save(file io.Reader, path string) (url string, err error)
	Delete(path string) error
}

type FileStorageUsecase struct {
	FileStorageLocal repository.LocalStorage
	UserRepo         UserRepository
}

func NewFileStorageUsecase(
	file repository.LocalStorage,
	userRepo UserRepository,
) *FileStorageUsecase {
	return &FileStorageUsecase{
		FileStorageLocal: file,
		UserRepo:         userRepo,
	}
}

func (u *FileStorageUsecase) Save(
	adminID uuid.UUID,
	file io.Reader,
	originalFilename string,
) (string, error) {
	isAdmin, err := u.UserRepo.IsAdmin(adminID)
	if err != nil {
		return "", err
	}
	if !isAdmin {
		return "", errors.New("You are not a admin ! ! ! ! ! !")
	}

	ext := filepath.Ext(originalFilename)
	id, err := uuid.NewV4()
	if err != nil {
		return "", ErrInntenal(err)
	}
	randomName := id.String() + ext
	relativePath := "equipment/" + randomName

	url, err := u.FileStorageLocal.Save(file, relativePath)
	if err != nil {
		return "", err
	}

	return url, nil
}

func (u *FileStorageUsecase) Delete(
	adminID uuid.UUID,
	url string,
) error {
	isAdmin, err := u.UserRepo.IsAdmin(adminID)
	if err != nil {
		return err
	}
	if !isAdmin {
		return errors.New("You are not a admin ! ! ! ! ! !")
	}

	return u.FileStorageLocal.Delete(url)
}
