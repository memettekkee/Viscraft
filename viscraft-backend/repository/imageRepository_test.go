package repository

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

// TestFindByPromptHash_ReturnsOnlyCompleted verifies that FindByPromptHash
// returns a completed image when one exists with the given hash.
func TestFindByPromptHash_ReturnsOnlyCompleted(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewImageRepository(db)

	now := time.Now()
	rows := sqlmock.NewRows([]string{
		"id", "project_id", "user_id", "prompt", "prompt_hash",
		"genre", "asset_type", "mood", "status", "file_path", "error_code", "created_at",
	}).AddRow(
		"img-001", "proj-001", "user-001", "a scenic mountain", "hash123",
		"landscape", "background", "calm", "completed", "/images/img-001.png", nil, now,
	)

	mock.ExpectQuery(`SELECT .+ FROM images WHERE prompt_hash = \$1 AND status = 'completed'`).
		WithArgs("hash123").
		WillReturnRows(rows)

	img, err := repo.FindByPromptHash("hash123")
	if err != nil {
		t.Fatalf("FindByPromptHash failed: %v", err)
	}

	if img.Id != "img-001" {
		t.Errorf("expected Id img-001, got %s", img.Id)
	}
	if img.Status != "completed" {
		t.Errorf("expected status completed, got %s", img.Status)
	}
	if img.FilePath != "/images/img-001.png" {
		t.Errorf("expected filePath /images/img-001.png, got %s", img.FilePath)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

// TestFindByPromptHash_NoCompletedImage verifies that FindByPromptHash returns
// sql.ErrNoRows when no completed image matches the given hash.
func TestFindByPromptHash_NoCompletedImage(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewImageRepository(db)

	mock.ExpectQuery(`SELECT .+ FROM images WHERE prompt_hash = \$1 AND status = 'completed'`).
		WithArgs("nonexistent-hash").
		WillReturnError(sql.ErrNoRows)

	_, err = repo.FindByPromptHash("nonexistent-hash")
	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("expected sql.ErrNoRows, got: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

// TestUpdateStatus_RejectsTransitionFromCompleted verifies that UpdateStatus
// returns sql.ErrNoRows when the image is already in a completed state.
func TestUpdateStatus_RejectsTransitionFromCompleted(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewImageRepository(db)

	// Simulate 0 rows affected — image is already completed, so WHERE status='processing' won't match
	mock.ExpectExec(`UPDATE images SET status = \$1, error_code = \$2 WHERE id = \$3 AND status = 'processing'`).
		WithArgs("failed", "ERR_03", "img-completed").
		WillReturnResult(sqlmock.NewResult(0, 0))

	err = repo.UpdateStatus("img-completed", "failed", "ERR_03")
	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("expected sql.ErrNoRows for completed image, got: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

// TestUpdateStatus_RejectsTransitionFromFailed verifies that UpdateStatus
// returns sql.ErrNoRows when the image is already in a failed state.
func TestUpdateStatus_RejectsTransitionFromFailed(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewImageRepository(db)

	// Image is already failed — WHERE status='processing' won't match
	mock.ExpectExec(`UPDATE images SET status = \$1, error_code = \$2 WHERE id = \$3 AND status = 'processing'`).
		WithArgs("completed", "", "img-failed").
		WillReturnResult(sqlmock.NewResult(0, 0))

	err = repo.UpdateStatus("img-failed", "completed", "")
	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("expected sql.ErrNoRows for failed image, got: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

// TestUpdateCompleted_RejectsIfNotProcessing verifies that UpdateCompleted
// returns sql.ErrNoRows when the image is not in processing state.
func TestUpdateCompleted_RejectsIfNotProcessing(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewImageRepository(db)

	// Image is already completed/failed — WHERE status='processing' won't match
	mock.ExpectExec(`UPDATE images SET status = 'completed', file_path = \$1 WHERE id = \$2 AND status = 'processing'`).
		WithArgs("/images/img-done.png", "img-already-done").
		WillReturnResult(sqlmock.NewResult(0, 0))

	err = repo.UpdateCompleted("img-already-done", "/images/img-done.png")
	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("expected sql.ErrNoRows for non-processing image, got: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

// TestUpdateStatus_SucceedsForProcessingImage verifies that UpdateStatus
// succeeds when the image is currently in processing state.
func TestUpdateStatus_SucceedsForProcessingImage(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewImageRepository(db)

	// 1 row affected — image was processing
	mock.ExpectExec(`UPDATE images SET status = \$1, error_code = \$2 WHERE id = \$3 AND status = 'processing'`).
		WithArgs("failed", "ERR_03", "img-processing").
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.UpdateStatus("img-processing", "failed", "ERR_03")
	if err != nil {
		t.Fatalf("expected nil error for processing image, got: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}
