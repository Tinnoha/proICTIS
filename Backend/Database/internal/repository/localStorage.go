package repository

import (
	"io"
	"os"
	"path/filepath"
	"strings"
)

type LocalStorage struct {
	path string
	url  string
}

func NewLocalStorage(path, url string) *LocalStorage {
	return &LocalStorage{
		path: path,
		url:  url,
	}
}

func (s *LocalStorage) Save(file io.Reader, plusPath string) (string, error) {
	fullPath := filepath.Join(s.path, plusPath)

	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}

	dst, err := os.Create(fullPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return "", err
	}

	url := s.url + "/" + strings.ReplaceAll(plusPath, "\\", "/")
	return url, nil
}

func (s *LocalStorage) Delete(url string) error {
	if !strings.HasPrefix(url, s.url+"/") {
		return nil
	}

	relativePath := strings.TrimPrefix(url, s.url+"/")
	fullPath := filepath.Join(s.path, relativePath)

	err := os.Remove(fullPath)
	if os.IsNotExist(err) {
		return nil
	}
	return err
}
