package service

import (
	"context"
	"database/sql"
	"encoding/base64"
	"os"
	"strings"
	"time"

	"viscraft-backend/constant"
	"viscraft-backend/model/request"
	"viscraft-backend/model/response"
	"viscraft-backend/pkg/logger"
	"viscraft-backend/repository"
)

type SceneImageGenerator interface {
	Generate(ctx context.Context, prompt string, referenceImage []byte) ([]byte, error)
}

type SceneStorageSaver interface {
	Save(id string, data []byte) (filePath string, fileUrl string, err error)
}

type SceneStorageDeleter interface {
	Delete(id string) error
}

type SceneProjectFinder interface {
	FindById(projectId, userId string) (*repository.Project, error)
}

type SceneService struct {
	repository     *repository.SceneRepository
	projectRepo    SceneProjectFinder
	storage        SceneStorageSaver
	storageDeleter SceneStorageDeleter
	imagegenClient SceneImageGenerator
}

func NewSceneService(
	repo *repository.SceneRepository,
	projectRepo SceneProjectFinder,
	storage SceneStorageSaver,
	storageDeleter SceneStorageDeleter,
	imagegenClient SceneImageGenerator,
) *SceneService {
	return &SceneService{
		repository:     repo,
		projectRepo:    projectRepo,
		storage:        storage,
		storageDeleter: storageDeleter,
		imagegenClient: imagegenClient,
	}
}

var sceneBlockedWords = []string{"nude", "explicit", "nsfw", "gore"}

func (s *SceneService) validatePrompt(prompt string) *constant.AppError {
	trimmed := strings.TrimSpace(prompt)
	if len(trimmed) < 3 {
		return &constant.ErrInvalidPrompt
	}
	if len(trimmed) > 300 {
		return &constant.ErrInvalidPrompt
	}
	lower := strings.ToLower(trimmed)
	for _, word := range sceneBlockedWords {
		if strings.Contains(lower, word) {
			return &constant.ErrInvalidPrompt
		}
	}
	return nil
}

func (s *SceneService) Generate(requestId, userId string, req request.GenerateSceneRequest) (response.GenerateSceneResponse, *constant.AppError) {
	if err := s.validatePrompt(req.Prompt); err != nil {
		logger.Warn(requestId, "prompt validation failed", "prompt", req.Prompt)
		return response.GenerateSceneResponse{}, err
	}

	_, findErr := s.projectRepo.FindById(req.ProjectId, userId)
	if findErr != nil {
		if findErr == sql.ErrNoRows {
			return response.GenerateSceneResponse{}, &constant.ErrProjectNotFound
		}
		logger.Error(requestId, "failed to verify project ownership", findErr)
		return response.GenerateSceneResponse{}, &constant.ErrDatabaseFailed
	}

	var referenceImageBytes []byte
	if req.UploadedReferenceImage != "" {
		decoded, appErr := s.validateUploadedReference(req.UploadedReferenceImage)
		if appErr != nil {
			return response.GenerateSceneResponse{}, appErr
		}
		referenceImageBytes = decoded
	}

	if req.ReferenceSceneId != "" && referenceImageBytes == nil {
		refScene, err := s.repository.FindById(req.ReferenceSceneId)
		if err == nil && refScene.Status == "completed" && refScene.FilePath != "" {
			refBytes, readErr := os.ReadFile(refScene.FilePath)
			if readErr == nil {
				referenceImageBytes = refBytes
			}
		}
	}

	orderIndex, err := s.repository.NextOrderIndex(req.ProjectId)
	if err != nil {
		logger.Error(requestId, "failed to calculate next order index", err)
		return response.GenerateSceneResponse{}, &constant.ErrDatabaseFailed
	}

	sceneId, err := s.repository.InsertProcessing(userId, req, orderIndex)
	if err != nil {
		logger.Error(requestId, "failed to insert processing record", err)
		return response.GenerateSceneResponse{}, &constant.ErrDatabaseFailed
	}

	generatedPrompt := req.Prompt
	if req.GeneratedPrompt != "" {
		generatedPrompt = req.GeneratedPrompt
	}

	go s.processGeneration(requestId, sceneId, generatedPrompt, referenceImageBytes)

	logger.Info(requestId, "ad shot generation started", "sceneId", sceneId, "orderIndex", orderIndex)

	return response.GenerateSceneResponse{
		BaseResponse: response.BaseResponse{
			RequestId: requestId,
			Success:   true,
			Message:   "Generation started",
		},
		Data: &response.SceneData{
			Id:         sceneId,
			OrderIndex: orderIndex,
			Prompt:     req.Prompt,
			Status:     "processing",
		},
	}, nil
}

func (s *SceneService) processGeneration(requestId, sceneId, prompt string, referenceImage []byte) {
	timeout := 60 * time.Second
	if referenceImage != nil {
		timeout = 90 * time.Second
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	logger.Info(requestId, "calling Pollinations API", "sceneId", sceneId, "hasReference", referenceImage != nil)
	imageBytes, err := s.imagegenClient.Generate(ctx, prompt, referenceImage)
	if err != nil {
		logger.Error(requestId, "Pollinations API call failed", "sceneId", sceneId, err)
		if ctx.Err() == context.DeadlineExceeded {
			s.repository.UpdateStatus(sceneId, "failed", constant.ErrGeminiTimeout.Code)
		} else if strings.HasPrefix(err.Error(), "content_policy_violation:") {
			s.repository.UpdateStatus(sceneId, "failed", constant.ErrContentPolicy.Code)
		} else {
			s.repository.UpdateStatus(sceneId, "failed", constant.ErrGeminiBadResponse.Code)
		}
		return
	}

	if len(imageBytes) == 0 {
		s.repository.UpdateStatus(sceneId, "failed", constant.ErrGeminiBadResponse.Code)
		return
	}

	filePath, fileUrl, err := s.storage.Save(sceneId, imageBytes)
	if err != nil {
		s.repository.UpdateStatus(sceneId, "failed", constant.ErrStorageFailed.Code)
		return
	}

	if err := s.repository.UpdateCompleted(sceneId, filePath, fileUrl); err != nil {
		logger.Error(requestId, "failed to update status to completed", "sceneId", sceneId, err)
		return
	}

	logger.Info(requestId, "ad shot generation completed", "sceneId", sceneId)
}

func (s *SceneService) validateUploadedReference(base64Data string) ([]byte, *constant.AppError) {
	decoded, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		decoded, err = base64.URLEncoding.DecodeString(base64Data)
		if err != nil {
			appErr := constant.AppError{Code: constant.ErrValidationFailed.Code, Message: "Invalid base64 encoding", HttpStatus: 422}
			return nil, &appErr
		}
	}
	const maxSize = 5 * 1024 * 1024
	if len(decoded) > maxSize {
		appErr := constant.AppError{Code: constant.ErrValidationFailed.Code, Message: "Reference image exceeds 5MB", HttpStatus: 422}
		return nil, &appErr
	}
	return decoded, nil
}

func (s *SceneService) GetScene(requestId, userId string, req request.GetSceneRequest) (response.GetSceneResponse, *constant.AppError) {
	scene, err := s.repository.FindByIdAndUser(req.Id, userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return response.GetSceneResponse{}, &constant.AppError{Code: constant.ErrImageNotFound.Code, Message: "Not found", HttpStatus: 404}
		}
		return response.GetSceneResponse{}, &constant.ErrDatabaseFailed
	}
	data := mapSceneToData(scene)
	return response.GetSceneResponse{
		BaseResponse: response.BaseResponse{RequestId: requestId, Success: true, Message: "Retrieved"},
		Data:         &data,
	}, nil
}

func (s *SceneService) ListScenes(requestId, userId string, req request.ListScenesRequest) (response.ListScenesResponse, *constant.AppError) {
	_, findErr := s.projectRepo.FindById(req.ProjectId, userId)
	if findErr != nil {
		if findErr == sql.ErrNoRows {
			return response.ListScenesResponse{}, &constant.ErrProjectNotFound
		}
		return response.ListScenesResponse{}, &constant.ErrDatabaseFailed
	}
	scenes, err := s.repository.FindByProjectId(req.ProjectId, userId)
	if err != nil {
		return response.ListScenesResponse{}, &constant.ErrDatabaseFailed
	}
	data := make([]response.SceneData, 0, len(scenes))
	for i := range scenes {
		data = append(data, mapSceneToData(&scenes[i]))
	}
	return response.ListScenesResponse{
		BaseResponse: response.BaseResponse{RequestId: requestId, Success: true, Message: "Retrieved"},
		Data:         data,
	}, nil
}

func (s *SceneService) DeleteScene(requestId, userId string, req request.DeleteSceneRequest) (response.DeleteSceneResponse, *constant.AppError) {
	err := s.repository.Delete(req.Id, userId)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Warn(requestId, "scene not found for deletion", "sceneId", req.Id, "userId", userId)
			return response.DeleteSceneResponse{}, &constant.AppError{Code: constant.ErrImageNotFound.Code, Message: "Not found", HttpStatus: 404}
		}
		logger.Error(requestId, "failed to delete scene", "sceneId", req.Id, err)
		return response.DeleteSceneResponse{}, &constant.ErrDatabaseFailed
	}
	if s.storageDeleter != nil {
		s.storageDeleter.Delete(req.Id)
	}
	logger.Info(requestId, "ad shot deleted", "sceneId", req.Id, "userId", userId)
	return response.DeleteSceneResponse{
		BaseResponse: response.BaseResponse{RequestId: requestId, Success: true, Message: "Deleted"},
	}, nil
}

func mapSceneToData(scene *repository.Scene) response.SceneData {
	data := response.SceneData{
		Id:              scene.Id,
		ProjectId:       scene.ProjectId,
		OrderIndex:      scene.OrderIndex,
		Prompt:          scene.Prompt,
		GeneratedPrompt: scene.GeneratedPrompt,
		Status:          scene.Status,
		CreatedAt:       scene.CreatedAt.Format(time.RFC3339),
	}
	if scene.Status == "completed" {
		data.FileUrl = scene.FileUrl
	}
	if scene.ErrorCode != "" {
		data.ErrorCode = scene.ErrorCode
	}
	return data
}
