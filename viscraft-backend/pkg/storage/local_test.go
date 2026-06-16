package storage

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewLocalStorage(t *testing.T) {
	ls := NewLocalStorage("/tmp/test-images")
	if ls == nil {
		t.Fatal("expected non-nil LocalStorage")
	}
	if ls.basePath != "/tmp/test-images" {
		t.Errorf("expected basePath /tmp/test-images, got %s", ls.basePath)
	}
}

func TestGetPath(t *testing.T) {
	ls := NewLocalStorage("/app/storage/images")
	path := ls.GetPath("abc-123")
	expected := filepath.Join("/app/storage/images", "abc-123.png")
	if path != expected {
		t.Errorf("expected %s, got %s", expected, path)
	}
}

func TestSave(t *testing.T) {
	dir := t.TempDir()
	ls := NewLocalStorage(dir)

	data := []byte{0x89, 0x50, 0x4E, 0x47} // PNG magic bytes
	path, err := ls.Save("img-001", data)
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	expected := filepath.Join(dir, "img-001.png")
	if path != expected {
		t.Errorf("expected path %s, got %s", expected, path)
	}

	// Verify file was written with correct content
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read saved file: %v", err)
	}
	if len(content) != len(data) {
		t.Errorf("expected %d bytes, got %d", len(data), len(content))
	}
	for i, b := range content {
		if b != data[i] {
			t.Errorf("byte mismatch at index %d: expected %x, got %x", i, data[i], b)
			break
		}
	}
}

func TestSave_CreatesDirectory(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "nested", "dir")
	ls := NewLocalStorage(dir)

	data := []byte("test image data")
	path, err := ls.Save("img-002", data)
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Error("expected file to exist after Save")
	}
}

func TestDelete(t *testing.T) {
	dir := t.TempDir()
	ls := NewLocalStorage(dir)

	// Create a file first
	data := []byte("image data")
	_, err := ls.Save("img-003", data)
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Delete it
	err = ls.Delete("img-003")
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// Verify file is gone
	path := ls.GetPath("img-003")
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Error("expected file to be deleted")
	}
}

func TestDelete_Idempotent(t *testing.T) {
	dir := t.TempDir()
	ls := NewLocalStorage(dir)

	// Delete a file that doesn't exist — should return nil
	err := ls.Delete("nonexistent-image")
	if err != nil {
		t.Errorf("expected nil error for missing file, got: %v", err)
	}
}

func TestDelete_IdempotentAfterDelete(t *testing.T) {
	dir := t.TempDir()
	ls := NewLocalStorage(dir)

	// Create and delete a file, then delete again
	data := []byte("image data")
	ls.Save("img-004", data)
	ls.Delete("img-004")

	// Second delete should succeed silently
	err := ls.Delete("img-004")
	if err != nil {
		t.Errorf("expected nil error on second delete, got: %v", err)
	}
}
