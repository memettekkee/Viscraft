package storage

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

// LocalStorage manages image files on the local filesystem.
type LocalStorage struct {
	basePath          string
	tempBasePath      string
	publicBaseUrl     string
	publicTempBaseUrl string
}

func NewLocalStorage(basePath, tempBasePath, publicBaseUrl, publicTempBaseUrl string) *LocalStorage {
	return &LocalStorage{
		basePath:          basePath,
		tempBasePath:      tempBasePath,
		publicBaseUrl:     publicBaseUrl,
		publicTempBaseUrl: publicTempBaseUrl,
	}
}

func (ls *LocalStorage) Save(id string, data []byte) (filePath string, fileUrl string, err error) {
	filePath = ls.GetPath(id)

	if err := os.MkdirAll(ls.basePath, 0755); err != nil {
		return "", "", fmt.Errorf("failed to create storage directory: %w", err)
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return "", "", fmt.Errorf("failed to write image file: %w", err)
	}

	fileUrl = ls.publicBaseUrl + "/" + id + ".png"
	return filePath, fileUrl, nil
}

func (ls *LocalStorage) SaveTemp(data []byte) (filePath string, fileUrl string, err error) {
	id := uuid.New().String()
	filename := id + ".png"
	filePath = filepath.Join(ls.tempBasePath, filename)

	if err := os.MkdirAll(ls.tempBasePath, 0755); err != nil {
		return "", "", fmt.Errorf("failed to create temp storage directory: %w", err)
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return "", "", fmt.Errorf("failed to write temp image file: %w", err)
	}

	fileUrl = ls.publicTempBaseUrl + "/" + filename
	return filePath, fileUrl, nil
}

func (ls *LocalStorage) DeleteTemp(filePath string) error {
	err := os.Remove(filePath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete temp file: %w", err)
	}
	return nil
}

func (ls *LocalStorage) Delete(imageId string) error {
	filePath := ls.GetPath(imageId)

	err := os.Remove(filePath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete image file: %w", err)
	}

	return nil
}

func (ls *LocalStorage) GetPath(imageId string) string {
	return filepath.Join(ls.basePath, imageId+".png")
}
