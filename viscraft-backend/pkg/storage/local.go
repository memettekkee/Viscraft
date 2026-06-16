package storage

import (
	"fmt"
	"os"
	"path/filepath"
)

// LocalStorage manages image files on the local filesystem.
type LocalStorage struct {
	basePath string
}

// NewLocalStorage creates a new LocalStorage instance with the given base directory path.
func NewLocalStorage(basePath string) *LocalStorage {
	return &LocalStorage{basePath: basePath}
}

// Save writes image data to the filesystem as basePath/imageId.png and returns the full file path.
func (ls *LocalStorage) Save(imageId string, data []byte) (string, error) {
	filePath := ls.GetPath(imageId)

	if err := os.MkdirAll(ls.basePath, 0755); err != nil {
		return "", fmt.Errorf("failed to create storage directory: %w", err)
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return "", fmt.Errorf("failed to write image file: %w", err)
	}

	return filePath, nil
}

// Delete removes the image file from the filesystem. Returns nil if the file
// does not exist (idempotent). Only returns an error for unexpected failures.
func (ls *LocalStorage) Delete(imageId string) error {
	filePath := ls.GetPath(imageId)

	err := os.Remove(filePath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete image file: %w", err)
	}

	return nil
}

// GetPath constructs and returns the file path for the given image ID.
func (ls *LocalStorage) GetPath(imageId string) string {
	return filepath.Join(ls.basePath, imageId+".png")
}
