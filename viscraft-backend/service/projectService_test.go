package service

import (
	"database/sql"
	"errors"
	"strings"
	"testing"
	"time"

	"viscraft-backend/constant"
	"viscraft-backend/model/request"
	"viscraft-backend/repository"

	"github.com/DATA-DOG/go-sqlmock"
)

// --- Mock implementations ---

// mockProjectImageFinder implements ProjectImageFinder for testing.
type mockProjectImageFinder struct {
	imageIds []string
	err      error
}

func (m *mockProjectImageFinder) FindImagesByProjectId(projectId string) ([]string, error) {
	return m.imageIds, m.err
}

// mockStorageDeleter implements StorageDeleter for testing.
type mockStorageDeleter struct {
	deletedIds []string
	err        error
}

func (m *mockStorageDeleter) Delete(imageId string) error {
	m.deletedIds = append(m.deletedIds, imageId)
	return m.err
}

// --- Helper ---

func newProjectServiceWithMockDB(t *testing.T) (*ProjectService, sqlmock.Sqlmock, *sql.DB) {
	t.Helper()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	repo := repository.NewProjectRepository(db)
	svc := NewProjectService(repo, nil, nil)
	return svc, mock, db
}

// --- Name Validation Tests ---

func TestProjectCreateProject_NameEmpty(t *testing.T) {
	svc := NewProjectService(nil, nil, nil)

	req := request.CreateProjectRequest{
		Name:        "",
		Description: "desc",
	}

	_, appErr := svc.CreateProject("req-1", "user-1", req)
	if appErr == nil {
		t.Fatal("expected error for empty name")
	}
	if appErr.Code != constant.ErrValidationFailed.Code {
		t.Errorf("expected error code %s, got %s", constant.ErrValidationFailed.Code, appErr.Code)
	}
}

func TestProjectCreateProject_NameWhitespaceOnly(t *testing.T) {
	svc := NewProjectService(nil, nil, nil)

	req := request.CreateProjectRequest{
		Name:        "   ",
		Description: "desc",
	}

	_, appErr := svc.CreateProject("req-1", "user-1", req)
	if appErr == nil {
		t.Fatal("expected error for whitespace-only name")
	}
	if appErr.Code != constant.ErrValidationFailed.Code {
		t.Errorf("expected error code %s, got %s", constant.ErrValidationFailed.Code, appErr.Code)
	}
}

func TestProjectCreateProject_Name255Chars(t *testing.T) {
	svc, mock, db := newProjectServiceWithMockDB(t)
	defer db.Close()

	name255 := strings.Repeat("a", 255)
	createdAt := time.Now()

	// Mock the INSERT query - the repo generates a UUID so we use AnyArg for the id
	mock.ExpectQuery(`INSERT INTO projects`).
		WithArgs(sqlmock.AnyArg(), "user-1", name255, "desc").
		WillReturnRows(sqlmock.NewRows([]string{"created_at"}).AddRow(createdAt))

	req := request.CreateProjectRequest{
		Name:        name255,
		Description: "desc",
	}

	resp, appErr := svc.CreateProject("req-1", "user-1", req)
	if appErr != nil {
		t.Fatalf("expected no error for 255-char name, got: %s", appErr.Code)
	}
	if !resp.Success {
		t.Fatal("expected success response")
	}
	if resp.Data == nil {
		t.Fatal("expected non-nil data")
	}
	if resp.Data.Name != name255 {
		t.Errorf("expected name of 255 chars, got length %d", len(resp.Data.Name))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet sqlmock expectations: %v", err)
	}
}

func TestProjectCreateProject_Name256Chars(t *testing.T) {
	svc := NewProjectService(nil, nil, nil)

	name256 := strings.Repeat("a", 256)

	req := request.CreateProjectRequest{
		Name:        name256,
		Description: "desc",
	}

	_, appErr := svc.CreateProject("req-1", "user-1", req)
	if appErr == nil {
		t.Fatal("expected error for 256-char name")
	}
	if appErr.Code != constant.ErrValidationFailed.Code {
		t.Errorf("expected error code %s, got %s", constant.ErrValidationFailed.Code, appErr.Code)
	}
}

// --- Ownership Enforcement Tests ---

func TestProjectGetProject_NotFoundReturnsProjectNotFound(t *testing.T) {
	svc, mock, db := newProjectServiceWithMockDB(t)
	defer db.Close()

	// When FindById returns sql.ErrNoRows (user doesn't own project or it doesn't exist)
	mock.ExpectQuery(`SELECT .+ FROM projects WHERE`).
		WithArgs("project-999", "user-1").
		WillReturnError(sql.ErrNoRows)

	req := request.GetProjectRequest{Id: "project-999"}

	_, appErr := svc.GetProject("req-1", "user-1", req)
	if appErr == nil {
		t.Fatal("expected error when project not found")
	}
	if appErr.Code != constant.ErrProjectNotFound.Code {
		t.Errorf("expected error code %s, got %s", constant.ErrProjectNotFound.Code, appErr.Code)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet sqlmock expectations: %v", err)
	}
}

func TestProjectDeleteProject_NotFoundReturnsProjectNotFound(t *testing.T) {
	svc, mock, db := newProjectServiceWithMockDB(t)
	defer db.Close()

	// FindById returns sql.ErrNoRows during ownership check
	mock.ExpectQuery(`SELECT .+ FROM projects WHERE`).
		WithArgs("project-999", "user-1").
		WillReturnError(sql.ErrNoRows)

	req := request.DeleteProjectRequest{Id: "project-999"}

	_, appErr := svc.DeleteProject("req-1", "user-1", req)
	if appErr == nil {
		t.Fatal("expected error when project not found for deletion")
	}
	if appErr.Code != constant.ErrProjectNotFound.Code {
		t.Errorf("expected error code %s, got %s", constant.ErrProjectNotFound.Code, appErr.Code)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet sqlmock expectations: %v", err)
	}
}

// --- Cascade Delete with Image Cleanup Tests ---

func TestProjectDeleteProject_SuccessWithImageCleanup(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	repo := repository.NewProjectRepository(db)

	imageIds := []string{"img-1", "img-2", "img-3"}
	mockFinder := &mockProjectImageFinder{imageIds: imageIds}
	mockStorage := &mockStorageDeleter{}

	svc := NewProjectService(repo, mockFinder, mockStorage)

	createdAt := time.Now()

	// 1. FindById succeeds (ownership check passes)
	mock.ExpectQuery(`SELECT .+ FROM projects WHERE`).
		WithArgs("project-1", "user-1").
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "name", "description", "created_at"}).
			AddRow("project-1", "user-1", "My Project", "desc", createdAt))

	// 2. Delete succeeds (after filesystem cleanup)
	mock.ExpectExec(`DELETE FROM projects WHERE`).
		WithArgs("project-1", "user-1").
		WillReturnResult(sqlmock.NewResult(0, 1))

	req := request.DeleteProjectRequest{Id: "project-1"}

	resp, appErr := svc.DeleteProject("req-1", "user-1", req)
	if appErr != nil {
		t.Fatalf("expected no error, got: %s", appErr.Code)
	}
	if !resp.Success {
		t.Fatal("expected success response")
	}

	// Verify storage.Delete was called for each image
	if len(mockStorage.deletedIds) != len(imageIds) {
		t.Errorf("expected %d storage deletes, got %d", len(imageIds), len(mockStorage.deletedIds))
	}
	for i, id := range imageIds {
		if i < len(mockStorage.deletedIds) && mockStorage.deletedIds[i] != id {
			t.Errorf("expected deleted id %s at index %d, got %s", id, i, mockStorage.deletedIds[i])
		}
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet sqlmock expectations: %v", err)
	}
}

func TestProjectDeleteProject_ContinuesWhenImageFinderFails(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	repo := repository.NewProjectRepository(db)

	// ImageFinder returns an error
	mockFinder := &mockProjectImageFinder{err: errors.New("image lookup failed")}
	mockStorage := &mockStorageDeleter{}

	svc := NewProjectService(repo, mockFinder, mockStorage)

	createdAt := time.Now()

	// FindById succeeds
	mock.ExpectQuery(`SELECT .+ FROM projects WHERE`).
		WithArgs("project-1", "user-1").
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "name", "description", "created_at"}).
			AddRow("project-1", "user-1", "My Project", "desc", createdAt))

	// Delete still proceeds
	mock.ExpectExec(`DELETE FROM projects WHERE`).
		WithArgs("project-1", "user-1").
		WillReturnResult(sqlmock.NewResult(0, 1))

	req := request.DeleteProjectRequest{Id: "project-1"}

	resp, appErr := svc.DeleteProject("req-1", "user-1", req)
	if appErr != nil {
		t.Fatalf("expected no error even when image finder fails, got: %s", appErr.Code)
	}
	if !resp.Success {
		t.Fatal("expected success response")
	}

	// Storage.Delete should NOT have been called since image lookup failed
	if len(mockStorage.deletedIds) != 0 {
		t.Errorf("expected 0 storage deletes when image finder fails, got %d", len(mockStorage.deletedIds))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet sqlmock expectations: %v", err)
	}
}
