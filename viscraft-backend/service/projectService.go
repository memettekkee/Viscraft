package service

import (
	"database/sql"
	"strings"
	"time"

	"viscraft-backend/constant"
	"viscraft-backend/model/request"
	"viscraft-backend/model/response"
	"viscraft-backend/pkg/logger"
	"viscraft-backend/repository"
)

// ProjectImageFinder defines the minimal interface for finding images by project ID.
// Used during project deletion to locate files for filesystem cleanup.
type ProjectImageFinder interface {
	FindImagesByProjectId(projectId string) ([]string, error)
}

// ProjectService handles project-related business logic including creation,
// retrieval, listing, and deletion with filesystem cleanup.
type ProjectService struct {
	projectRepo  *repository.ProjectRepository
	imageFinder  ProjectImageFinder
	storage      StorageDeleter
}

// NewProjectService creates a new ProjectService with the required dependencies.
// imageFinder and storage can be nil if filesystem cleanup is not needed (they are
// checked before use during deletion).
func NewProjectService(
	projectRepo *repository.ProjectRepository,
	imageFinder ProjectImageFinder,
	storage StorageDeleter,
) *ProjectService {
	return &ProjectService{
		projectRepo: projectRepo,
		imageFinder: imageFinder,
		storage:     storage,
	}
}

// CreateProject validates the input and creates a new project for the authenticated user.
func (s *ProjectService) CreateProject(requestId string, userId string, req request.CreateProjectRequest) (response.CreateProjectResponse, *constant.AppError) {
	// Validate name: trim and check length 1-255
	name := strings.TrimSpace(req.Name)
	if len(name) == 0 || len(name) > 255 {
		logger.Warn(requestId, "project name validation failed", "name", req.Name)
		return response.CreateProjectResponse{}, &constant.ErrValidationFailed
	}

	project, err := s.projectRepo.Insert(userId, name, req.Description)
	if err != nil {
		logger.Error(requestId, "failed to insert project", err)
		return response.CreateProjectResponse{}, &constant.ErrDatabaseFailed
	}

	logger.Info(requestId, "project created", "projectId", project.Id, "userId", userId)

	return response.CreateProjectResponse{
		BaseResponse: response.BaseResponse{
			RequestId: requestId,
			Success:   true,
			Message:   "Project created successfully",
		},
		Data: &response.ProjectData{
			Id:          project.Id,
			Name:        project.Name,
			Description: project.Description,
			CreatedAt:   project.CreatedAt.Format(time.RFC3339),
		},
	}, nil
}

// GetProject retrieves a project by ID, verifying ownership via user_id filter.
func (s *ProjectService) GetProject(requestId string, userId string, req request.GetProjectRequest) (response.GetProjectResponse, *constant.AppError) {
	project, err := s.projectRepo.FindById(req.Id, userId)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Warn(requestId, "project not found", "projectId", req.Id, "userId", userId)
			return response.GetProjectResponse{}, &constant.ErrProjectNotFound
		}
		logger.Error(requestId, "failed to find project", err)
		return response.GetProjectResponse{}, &constant.ErrDatabaseFailed
	}

	logger.Info(requestId, "project retrieved", "projectId", project.Id)

	return response.GetProjectResponse{
		BaseResponse: response.BaseResponse{
			RequestId: requestId,
			Success:   true,
			Message:   "Project retrieved successfully",
		},
		Data: &response.ProjectData{
			Id:          project.Id,
			Name:        project.Name,
			Description: project.Description,
			CreatedAt:   project.CreatedAt.Format(time.RFC3339),
		},
	}, nil
}

// ListProjects returns all projects belonging to the authenticated user.
func (s *ProjectService) ListProjects(requestId string, userId string) (response.ListProjectsResponse, *constant.AppError) {
	projects, err := s.projectRepo.FindByUserId(userId)
	if err != nil {
		logger.Error(requestId, "failed to list projects", err)
		return response.ListProjectsResponse{}, &constant.ErrDatabaseFailed
	}

	data := make([]response.ProjectData, 0, len(projects))
	for _, p := range projects {
		data = append(data, response.ProjectData{
			Id:          p.Id,
			Name:        p.Name,
			Description: p.Description,
			CreatedAt:   p.CreatedAt.Format(time.RFC3339),
		})
	}

	logger.Info(requestId, "projects listed", "userId", userId, "count", len(data))

	return response.ListProjectsResponse{
		BaseResponse: response.BaseResponse{
			RequestId: requestId,
			Success:   true,
			Message:   "Projects retrieved successfully",
		},
		Data: data,
	}, nil
}

// DeleteProject verifies ownership, triggers filesystem cleanup for project images,
// and deletes the project. Database cascade handles image record removal.
func (s *ProjectService) DeleteProject(requestId string, userId string, req request.DeleteProjectRequest) (response.DeleteProjectResponse, *constant.AppError) {
	// Verify the project exists and belongs to the user before cleanup
	_, err := s.projectRepo.FindById(req.Id, userId)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Warn(requestId, "project not found for deletion", "projectId", req.Id, "userId", userId)
			return response.DeleteProjectResponse{}, &constant.ErrProjectNotFound
		}
		logger.Error(requestId, "failed to find project for deletion", err)
		return response.DeleteProjectResponse{}, &constant.ErrDatabaseFailed
	}

	// Perform filesystem cleanup if dependencies are available
	if s.imageFinder != nil && s.storage != nil {
		imageIds, findErr := s.imageFinder.FindImagesByProjectId(req.Id)
		if findErr != nil {
			logger.Error(requestId, "failed to find project images for cleanup", findErr)
			// Continue with deletion even if image lookup fails
		} else {
			for _, imageId := range imageIds {
				if delErr := s.storage.Delete(imageId); delErr != nil {
					logger.Warn(requestId, "failed to delete image file", "imageId", imageId, "error", delErr.Error())
				}
			}
			logger.Info(requestId, "filesystem cleanup completed", "projectId", req.Id, "imageCount", len(imageIds))
		}
	}

	// Delete project from database (cascade handles image DB records)
	if err := s.projectRepo.Delete(req.Id, userId); err != nil {
		if err == sql.ErrNoRows {
			logger.Warn(requestId, "project not found during delete", "projectId", req.Id)
			return response.DeleteProjectResponse{}, &constant.ErrProjectNotFound
		}
		logger.Error(requestId, "failed to delete project", err)
		return response.DeleteProjectResponse{}, &constant.ErrDatabaseFailed
	}

	logger.Info(requestId, "project deleted", "projectId", req.Id, "userId", userId)

	return response.DeleteProjectResponse{
		BaseResponse: response.BaseResponse{
			RequestId: requestId,
			Success:   true,
			Message:   "Project deleted successfully",
		},
	}, nil
}
