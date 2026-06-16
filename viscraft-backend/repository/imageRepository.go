package repository

import (
	"database/sql"
	"time"

	"viscraft-backend/model/request"

	"github.com/google/uuid"
)

// Image represents a row in the images table.
type Image struct {
	Id        string
	ProjectId string
	UserId    string
	Prompt    string
	PromptHash string
	Genre     string
	AssetType string
	Mood      string
	Status    string
	FilePath  string
	ErrorCode string
	CreatedAt time.Time
}

// ImageRecord holds the minimal fields needed for cleanup operations.
// Used by FindImagesByUserId for user deletion cleanup.
type ImageRecord struct {
	Id string
}

// ImageRepository handles database operations for the images table.
type ImageRepository struct {
	db *sql.DB
}

// NewImageRepository creates a new ImageRepository with the given database connection.
func NewImageRepository(db *sql.DB) *ImageRepository {
	return &ImageRepository{db: db}
}

// FindById retrieves an image by its ID, always filtering by user_id for ownership enforcement.
func (r *ImageRepository) FindById(imageId, userId string) (*Image, error) {
	img := &Image{}
	var filePath, errorCode, promptHash sql.NullString

	err := r.db.QueryRow(
		`SELECT id, project_id, user_id, prompt, prompt_hash, genre, asset_type, mood, status, file_path, error_code, created_at
		 FROM images WHERE id = $1 AND user_id = $2`,
		imageId, userId,
	).Scan(
		&img.Id, &img.ProjectId, &img.UserId, &img.Prompt, &promptHash,
		&img.Genre, &img.AssetType, &img.Mood, &img.Status, &filePath, &errorCode, &img.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	if filePath.Valid {
		img.FilePath = filePath.String
	}
	if errorCode.Valid {
		img.ErrorCode = errorCode.String
	}
	if promptHash.Valid {
		img.PromptHash = promptHash.String
	}

	return img, nil
}

// FindByProjectId retrieves all images for a project, ordered by created_at DESC.
// Always filters by user_id for ownership enforcement.
func (r *ImageRepository) FindByProjectId(projectId, userId string) ([]Image, error) {
	rows, err := r.db.Query(
		`SELECT id, project_id, user_id, prompt, prompt_hash, genre, asset_type, mood, status, file_path, error_code, created_at
		 FROM images WHERE project_id = $1 AND user_id = $2 ORDER BY created_at DESC`,
		projectId, userId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var images []Image
	for rows.Next() {
		var img Image
		var filePath, errorCode, promptHash sql.NullString

		if err := rows.Scan(
			&img.Id, &img.ProjectId, &img.UserId, &img.Prompt, &promptHash,
			&img.Genre, &img.AssetType, &img.Mood, &img.Status, &filePath, &errorCode, &img.CreatedAt,
		); err != nil {
			return nil, err
		}

		if filePath.Valid {
			img.FilePath = filePath.String
		}
		if errorCode.Valid {
			img.ErrorCode = errorCode.String
		}
		if promptHash.Valid {
			img.PromptHash = promptHash.String
		}

		images = append(images, img)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return images, nil
}

// FindByPromptHash performs a cache lookup by prompt hash.
// Returns only completed images (status = 'completed').
func (r *ImageRepository) FindByPromptHash(hash string) (*Image, error) {
	img := &Image{}
	var filePath, errorCode, promptHash sql.NullString

	err := r.db.QueryRow(
		`SELECT id, project_id, user_id, prompt, prompt_hash, genre, asset_type, mood, status, file_path, error_code, created_at
		 FROM images WHERE prompt_hash = $1 AND status = 'completed' LIMIT 1`,
		hash,
	).Scan(
		&img.Id, &img.ProjectId, &img.UserId, &img.Prompt, &promptHash,
		&img.Genre, &img.AssetType, &img.Mood, &img.Status, &filePath, &errorCode, &img.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	if filePath.Valid {
		img.FilePath = filePath.String
	}
	if errorCode.Valid {
		img.ErrorCode = errorCode.String
	}
	if promptHash.Valid {
		img.PromptHash = promptHash.String
	}

	return img, nil
}

// InsertProcessing inserts a new image record with status="processing" and returns the generated image ID.
// The userId parameter is the authenticated user's ID for ownership enforcement.
func (r *ImageRepository) InsertProcessing(userId string, req request.GenerateImageRequest, hash string) (string, error) {
	id := uuid.New().String()

	_, err := r.db.Exec(
		`INSERT INTO images (id, project_id, user_id, prompt, prompt_hash, genre, asset_type, mood, status)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, 'processing')`,
		id, req.ProjectId, userId, req.Prompt, hash, req.Genre, req.AssetType, req.Mood,
	)
	if err != nil {
		return "", err
	}

	return id, nil
}

// UpdateStatus updates the status of an image only if the current status is "processing".
// This prevents overwriting terminal states (completed, failed).
func (r *ImageRepository) UpdateStatus(imageId, status, errorCode string) error {
	result, err := r.db.Exec(
		`UPDATE images SET status = $1, error_code = $2 WHERE id = $3 AND status = 'processing'`,
		status, errorCode, imageId,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// UpdateCompleted sets the image status to "completed" and stores the file path.
// Only updates if the current status is "processing".
func (r *ImageRepository) UpdateCompleted(imageId, filePath string) error {
	result, err := r.db.Exec(
		`UPDATE images SET status = 'completed', file_path = $1 WHERE id = $2 AND status = 'processing'`,
		filePath, imageId,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// Delete removes an image record with user_id ownership check.
func (r *ImageRepository) Delete(imageId, userId string) error {
	result, err := r.db.Exec(
		`DELETE FROM images WHERE id = $1 AND user_id = $2`,
		imageId, userId,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// FindImagesByUserId retrieves all image records for a given user.
// Used for filesystem cleanup during user deletion.
func (r *ImageRepository) FindImagesByUserId(userId string) ([]ImageRecord, error) {
	rows, err := r.db.Query(
		`SELECT id FROM images WHERE user_id = $1`,
		userId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []ImageRecord
	for rows.Next() {
		var rec ImageRecord
		if err := rows.Scan(&rec.Id); err != nil {
			return nil, err
		}
		records = append(records, rec)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return records, nil
}
