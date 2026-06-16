package service

import (
	"database/sql"
	"fmt"
	"time"

	"viscraft-backend/constant"
	"viscraft-backend/model/request"
	"viscraft-backend/model/response"
	"viscraft-backend/pkg/logger"
	"viscraft-backend/repository"
)

// ImageService handles image-related business logic including retrieval,
// listing, and deletion with filesystem cleanup.
type ImageService struct {
	imageRepo *repository.ImageRepository
	storage   StorageDeleter
}

// NewImageService creates a new ImageService with the required dependencies.
func NewImageService(imageRepo *repository.ImageRepository, storage StorageDeleter) *ImageService {
	return &ImageService{
		imageRepo: imageRepo,
		storage:   storage,
	}
}

// GetImage retrieves an image by ID, verifying ownership via user_id filter.
// If the image status is "completed", it includes the fileUrl in the response.
func (s *ImageService) GetImage(requestId, userId string, req request.GetImageRequest) (response.GetImageResponse, *constant.AppError) {
	img, err := s.imageRepo.FindById(req.Id, userId)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Warn(requestId, "image not found", "imageId", req.Id, "userId", userId)
			return response.GetImageResponse{}, &constant.ErrImageNotFound
		}
		logger.Error(requestId, "failed to find image", err)
		return response.GetImageResponse{}, &constant.ErrDatabaseFailed
	}

	data := mapImageToData(img)

	logger.Info(requestId, "image retrieved", "imageId", img.Id)

	return response.GetImageResponse{
		BaseResponse: response.BaseResponse{
			RequestId: requestId,
			Success:   true,
			Message:   "Image retrieved successfully",
		},
		Data: &data,
	}, nil
}

// ListImages retrieves all images for a project, verifying ownership via user_id filter.
// Results are ordered by creation date descending. Completed images include fileUrl.
func (s *ImageService) ListImages(requestId, userId string, req request.ListImagesRequest) (response.ListImagesResponse, *constant.AppError) {
	images, err := s.imageRepo.FindByProjectId(req.ProjectId, userId)
	if err != nil {
		logger.Error(requestId, "failed to list images", err)
		return response.ListImagesResponse{}, &constant.ErrDatabaseFailed
	}

	data := make([]response.ImageData, 0, len(images))
	for i := range images {
		data = append(data, mapImageToData(&images[i]))
	}

	logger.Info(requestId, "images listed", "projectId", req.ProjectId, "userId", userId, "count", len(data))

	return response.ListImagesResponse{
		BaseResponse: response.BaseResponse{
			RequestId: requestId,
			Success:   true,
			Message:   "Images retrieved successfully",
		},
		Data: data,
	}, nil
}

// DeleteImage verifies ownership and status before deleting an image.
// Rejects deletion if the image is currently being processed (status="processing").
// Removes both the database record and the stored file.
func (s *ImageService) DeleteImage(requestId, userId string, req request.DeleteImageRequest) (response.DeleteImageResponse, *constant.AppError) {
	// Check image exists and belongs to user
	img, err := s.imageRepo.FindById(req.Id, userId)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Warn(requestId, "image not found for deletion", "imageId", req.Id, "userId", userId)
			return response.DeleteImageResponse{}, &constant.ErrImageNotFound
		}
		logger.Error(requestId, "failed to find image for deletion", err)
		return response.DeleteImageResponse{}, &constant.ErrDatabaseFailed
	}

	// Reject deletion of images currently being processed
	if img.Status == "processing" {
		logger.Warn(requestId, "cannot delete image in processing state", "imageId", req.Id)
		processingErr := constant.AppError{
			Code:       constant.ErrValidationFailed.Code,
			Message:    "Cannot delete image while processing",
			HttpStatus: 422,
		}
		return response.DeleteImageResponse{}, &processingErr
	}

	// Delete database record
	if err := s.imageRepo.Delete(req.Id, userId); err != nil {
		logger.Error(requestId, "failed to delete image from database", err)
		return response.DeleteImageResponse{}, &constant.ErrDatabaseFailed
	}

	// Delete file from storage (best-effort, log errors but don't fail)
	if s.storage != nil {
		if err := s.storage.Delete(req.Id); err != nil {
			logger.Warn(requestId, "failed to delete image file", "imageId", req.Id, "error", err.Error())
		}
	}

	logger.Info(requestId, "image deleted", "imageId", req.Id, "userId", userId)

	return response.DeleteImageResponse{
		BaseResponse: response.BaseResponse{
			RequestId: requestId,
			Success:   true,
			Message:   "Image deleted successfully",
		},
	}, nil
}

// mapImageToData converts a repository Image to a response ImageData.
// If the image status is "completed", it includes the fileUrl.
func mapImageToData(img *repository.Image) response.ImageData {
	data := response.ImageData{
		Id:        img.Id,
		Status:    img.Status,
		Prompt:    img.Prompt,
		Genre:     img.Genre,
		AssetType: img.AssetType,
		Mood:      img.Mood,
		CreatedAt: img.CreatedAt.Format(time.RFC3339),
	}

	if img.Status == "completed" {
		data.FileUrl = fmt.Sprintf("/storage/images/%s.png", img.Id)
	}

	if img.ErrorCode != "" {
		data.ErrorCode = img.ErrorCode
	}

	return data
}
