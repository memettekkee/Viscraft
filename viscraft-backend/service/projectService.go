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

type ProjectSceneFinder interface {
	FindScenesByProjectId(projectId string) ([]string, error)
}

type ProjectService struct {
	projectRepo *repository.ProjectRepository
	sceneFinder ProjectSceneFinder
	storage     StorageDeleter
}

func NewProjectService(
	projectRepo *repository.ProjectRepository,
	sceneFinder ProjectSceneFinder,
	storage StorageDeleter,
) *ProjectService {
	return &ProjectService{
		projectRepo: projectRepo,
		sceneFinder: sceneFinder,
		storage:     storage,
	}
}

func (s *ProjectService) CreateProject(requestId string, userId string, req request.CreateProjectRequest) (response.CreateProjectResponse, *constant.AppError) {
	name := strings.TrimSpace(req.Name)
	if len(name) == 0 || len(name) > 255 {
		logger.Warn(requestId, "project name validation failed", "name", req.Name)
		return response.CreateProjectResponse{}, &constant.ErrValidationFailed
	}

	project, err := s.projectRepo.Insert(userId, name, req.Description, req.ProductCategory, req.VisualStyle)
	if err != nil {
		logger.Error(requestId, "failed to insert project", err)
		return response.CreateProjectResponse{}, &constant.ErrDatabaseFailed
	}

	logger.Info(requestId, "campaign created", "projectId", project.Id, "userId", userId)

	return response.CreateProjectResponse{
		BaseResponse: response.BaseResponse{RequestId: requestId, Success: true, Message: "Campaign created successfully"},
		Data:         mapProjectToData(project),
	}, nil
}

func (s *ProjectService) GetProject(requestId string, userId string, req request.GetProjectRequest) (response.GetProjectResponse, *constant.AppError) {
	project, err := s.projectRepo.FindById(req.Id, userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return response.GetProjectResponse{}, &constant.ErrProjectNotFound
		}
		logger.Error(requestId, "failed to find project", err)
		return response.GetProjectResponse{}, &constant.ErrDatabaseFailed
	}

	return response.GetProjectResponse{
		BaseResponse: response.BaseResponse{RequestId: requestId, Success: true, Message: "Retrieved"},
		Data:         mapProjectToData(project),
	}, nil
}

func (s *ProjectService) ListProjects(requestId string, userId string) (response.ListProjectsResponse, *constant.AppError) {
	projects, err := s.projectRepo.FindByUserId(userId)
	if err != nil {
		logger.Error(requestId, "failed to list projects", err)
		return response.ListProjectsResponse{}, &constant.ErrDatabaseFailed
	}

	data := make([]response.ProjectData, 0, len(projects))
	for i := range projects {
		data = append(data, *mapProjectToData(&projects[i]))
	}

	return response.ListProjectsResponse{
		BaseResponse: response.BaseResponse{RequestId: requestId, Success: true, Message: "Retrieved"},
		Data:         data,
	}, nil
}

func (s *ProjectService) DeleteProject(requestId string, userId string, req request.DeleteProjectRequest) (response.DeleteProjectResponse, *constant.AppError) {
	_, err := s.projectRepo.FindById(req.Id, userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return response.DeleteProjectResponse{}, &constant.ErrProjectNotFound
		}
		return response.DeleteProjectResponse{}, &constant.ErrDatabaseFailed
	}

	if s.sceneFinder != nil && s.storage != nil {
		sceneIds, findErr := s.sceneFinder.FindScenesByProjectId(req.Id)
		if findErr == nil {
			for _, sceneId := range sceneIds {
				s.storage.Delete(sceneId)
			}
		}
	}

	if err := s.projectRepo.Delete(req.Id, userId); err != nil {
		if err == sql.ErrNoRows {
			return response.DeleteProjectResponse{}, &constant.ErrProjectNotFound
		}
		return response.DeleteProjectResponse{}, &constant.ErrDatabaseFailed
	}

	logger.Info(requestId, "campaign deleted", "projectId", req.Id, "userId", userId)

	return response.DeleteProjectResponse{
		BaseResponse: response.BaseResponse{RequestId: requestId, Success: true, Message: "Deleted"},
	}, nil
}

func mapProjectToData(p *repository.Project) *response.ProjectData {
	return &response.ProjectData{
		Id:              p.Id,
		Name:            p.Name,
		Description:     p.Description,
		ProductCategory: p.ProductCategory,
		VisualStyle:     p.VisualStyle,
		CreatedAt:       p.CreatedAt.Format(time.RFC3339),
	}
}
