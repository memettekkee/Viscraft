package repository

import (
	"database/sql"
	"time"

	"viscraft-backend/model/request"

	"github.com/google/uuid"
)

type Scene struct {
	Id                    string
	ProjectId             string
	UserId                string
	OrderIndex            int
	Prompt                string
	GeneratedPrompt       string
	ReferenceSceneId      string
	UsedUploadedReference bool
	Status                string
	FilePath              string
	FileUrl               string
	ErrorCode             string
	CreatedAt             time.Time
}

type SceneRecord struct {
	Id string
}

type SceneRepository struct {
	db *sql.DB
}

func NewSceneRepository(db *sql.DB) *SceneRepository {
	return &SceneRepository{db: db}
}

func (r *SceneRepository) InsertProcessing(userId string, req request.GenerateSceneRequest, orderIndex int) (string, error) {
	id := uuid.New().String()

	var referenceSceneId interface{}
	if req.ReferenceSceneId != "" {
		referenceSceneId = req.ReferenceSceneId
	} else {
		referenceSceneId = nil
	}

	usedUploadedReference := req.UploadedReferenceImage != ""

	_, err := r.db.Exec(
		`INSERT INTO scenes (id, project_id, user_id, order_index, user_prompt, generated_prompt, reference_scene_id, used_uploaded_reference, status)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, 'processing')`,
		id, req.ProjectId, userId, orderIndex, req.Prompt, req.GeneratedPrompt, referenceSceneId, usedUploadedReference,
	)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (r *SceneRepository) FindById(sceneId string) (*Scene, error) {
	scene := &Scene{}
	var referenceSceneId, generatedPrompt, filePath, fileUrl, errorCode sql.NullString

	err := r.db.QueryRow(
		`SELECT id, project_id, user_id, order_index, user_prompt, generated_prompt, reference_scene_id, used_uploaded_reference, status, file_path, file_url, error_code, created_at
		 FROM scenes WHERE id = $1`,
		sceneId,
	).Scan(
		&scene.Id, &scene.ProjectId, &scene.UserId, &scene.OrderIndex, &scene.Prompt, &generatedPrompt,
		&referenceSceneId, &scene.UsedUploadedReference, &scene.Status, &filePath, &fileUrl, &errorCode, &scene.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	if generatedPrompt.Valid {
		scene.GeneratedPrompt = generatedPrompt.String
	}
	if referenceSceneId.Valid {
		scene.ReferenceSceneId = referenceSceneId.String
	}
	if filePath.Valid {
		scene.FilePath = filePath.String
	}
	if fileUrl.Valid {
		scene.FileUrl = fileUrl.String
	}
	if errorCode.Valid {
		scene.ErrorCode = errorCode.String
	}

	return scene, nil
}

func (r *SceneRepository) FindByIdAndUser(sceneId, userId string) (*Scene, error) {
	scene := &Scene{}
	var referenceSceneId, generatedPrompt, filePath, fileUrl, errorCode sql.NullString

	err := r.db.QueryRow(
		`SELECT id, project_id, user_id, order_index, user_prompt, generated_prompt, reference_scene_id, used_uploaded_reference, status, file_path, file_url, error_code, created_at
		 FROM scenes WHERE id = $1 AND user_id = $2`,
		sceneId, userId,
	).Scan(
		&scene.Id, &scene.ProjectId, &scene.UserId, &scene.OrderIndex, &scene.Prompt, &generatedPrompt,
		&referenceSceneId, &scene.UsedUploadedReference, &scene.Status, &filePath, &fileUrl, &errorCode, &scene.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	if generatedPrompt.Valid {
		scene.GeneratedPrompt = generatedPrompt.String
	}
	if referenceSceneId.Valid {
		scene.ReferenceSceneId = referenceSceneId.String
	}
	if filePath.Valid {
		scene.FilePath = filePath.String
	}
	if fileUrl.Valid {
		scene.FileUrl = fileUrl.String
	}
	if errorCode.Valid {
		scene.ErrorCode = errorCode.String
	}

	return scene, nil
}

func (r *SceneRepository) FindByProjectId(projectId, userId string) ([]Scene, error) {
	rows, err := r.db.Query(
		`SELECT id, project_id, user_id, order_index, user_prompt, generated_prompt,reference_scene_id, used_uploaded_reference, status, file_path, file_url, error_code, created_at
		 FROM scenes WHERE project_id = $1 AND user_id = $2 ORDER BY order_index ASC`,
		projectId, userId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var scenes []Scene
	for rows.Next() {
		var scene Scene
		var referenceSceneId, generatedPrompt, filePath, fileUrl, errorCode sql.NullString

		if err := rows.Scan(
			&scene.Id, &scene.ProjectId, &scene.UserId, &scene.OrderIndex, &scene.Prompt, &generatedPrompt,
			&referenceSceneId, &scene.UsedUploadedReference, &scene.Status, &filePath, &fileUrl, &errorCode, &scene.CreatedAt,
		); err != nil {
			return nil, err
		}

		if generatedPrompt.Valid {
			scene.GeneratedPrompt = generatedPrompt.String
		}
		if referenceSceneId.Valid {
			scene.ReferenceSceneId = referenceSceneId.String
		}
		if filePath.Valid {
			scene.FilePath = filePath.String
		}
		if fileUrl.Valid {
			scene.FileUrl = fileUrl.String
		}
		if errorCode.Valid {
			scene.ErrorCode = errorCode.String
		}

		scenes = append(scenes, scene)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return scenes, nil
}

func (r *SceneRepository) NextOrderIndex(projectId string) (int, error) {
	var nextIndex int
	err := r.db.QueryRow(
		`SELECT COALESCE(MAX(order_index) + 1, 0) FROM scenes WHERE project_id = $1`,
		projectId,
	).Scan(&nextIndex)
	if err != nil {
		return 0, err
	}
	return nextIndex, nil
}

func (r *SceneRepository) UpdateStatus(sceneId, status, errorCode string) error {
	result, err := r.db.Exec(
		`UPDATE scenes SET status = $1, error_code = $2 WHERE id = $3 AND status = 'processing'`,
		status, errorCode, sceneId,
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

func (r *SceneRepository) UpdateCompleted(sceneId, filePath, fileUrl string) error {
	result, err := r.db.Exec(
		`UPDATE scenes SET status = 'completed', file_path = $1, file_url = $2 WHERE id = $3 AND status = 'processing'`,
		filePath, fileUrl, sceneId,
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

func (r *SceneRepository) Delete(sceneId, userId string) error {
	result, err := r.db.Exec(
		`DELETE FROM scenes WHERE id = $1 AND user_id = $2`,
		sceneId, userId,
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

func (r *SceneRepository) FindScenesByProjectId(projectId string) ([]string, error) {
	rows, err := r.db.Query(
		`SELECT id FROM scenes WHERE project_id = $1`,
		projectId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sceneIds []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		sceneIds = append(sceneIds, id)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return sceneIds, nil
}

func (r *SceneRepository) FindScenesByUserId(userId string) ([]SceneRecord, error) {
	rows, err := r.db.Query(
		`SELECT id FROM scenes WHERE user_id = $1`,
		userId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []SceneRecord
	for rows.Next() {
		var rec SceneRecord
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
