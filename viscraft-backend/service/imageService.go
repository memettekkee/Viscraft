package service

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"viscraft-backend/constant"
	"viscraft-backend/model/request"
	"viscraft-backend/model/response"
	"viscraft-backend/pkg/gemini"
	"viscraft-backend/pkg/logger"
	"viscraft-backend/repository"
)

// GeminiGenerator defines the interface for generating images via the Gemini API.
type GeminiGenerator interface {
	Generate(ctx context.Context, prompt string, refImage *gemini.ReferenceImage) ([]byte, error)
}

// StorageSaver defines the interface for saving image files to storage.
type StorageSaver interface {
	Save(imageId string, data []byte) (string, error)
}

// ProjectOwnershipChecker defines the interface for verifying project ownership.
type ProjectOwnershipChecker interface {
	FindById(projectId, userId string) (*repository.Project, error)
}

// ImageService handles image-related business logic including generation,
// retrieval, listing, and deletion with filesystem cleanup.
type ImageService struct {
	imageRepo    *repository.ImageRepository
	storage      StorageDeleter
	storageSaver StorageSaver
	geminiClient GeminiGenerator
	projectRepo  ProjectOwnershipChecker
}

// NewImageService creates a new ImageService with the required dependencies.
// geminiClient, storageSaver, and projectRepo can be nil if generation is not needed.
func NewImageService(
	imageRepo *repository.ImageRepository,
	storage StorageDeleter,
	geminiClient GeminiGenerator,
	storageSaver StorageSaver,
	projectRepo ProjectOwnershipChecker,
) *ImageService {
	return &ImageService{
		imageRepo:    imageRepo,
		storage:      storage,
		geminiClient: geminiClient,
		storageSaver: storageSaver,
		projectRepo:  projectRepo,
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

// hashPrompt computes a SHA256 hash from the given parts joined by a pipe delimiter.
// Returns a deterministic 64-character hexadecimal string.
func (s *ImageService) hashPrompt(parts ...string) string {
	combined := strings.Join(parts, "|")
	hash := sha256.Sum256([]byte(combined))
	return hex.EncodeToString(hash[:])
}

// blockedWords contains content terms that are not allowed in prompts.
var blockedWords = []string{"nude", "explicit", "nsfw", "gore"}

// validatePrompt checks that the prompt meets length requirements (3-300 characters
// after trimming whitespace) and does not contain any blocked content words.
// Returns nil if valid, or *constant.ErrInvalidPrompt if validation fails.
func (s *ImageService) validatePrompt(prompt string) *constant.AppError {
	trimmed := strings.TrimSpace(prompt)

	if len(trimmed) < 3 {
		return &constant.ErrInvalidPrompt
	}

	if len(trimmed) > 300 {
		return &constant.ErrInvalidPrompt
	}

	lower := strings.ToLower(trimmed)
	for _, word := range blockedWords {
		if strings.Contains(lower, word) {
			return &constant.ErrInvalidPrompt
		}
	}

	return nil
}

// allowedGenres defines the valid genre values (case-sensitive).
var allowedGenres = map[string]bool{
	"fantasy":          true,
	"sci-fi":           true,
	"post-apocalyptic": true,
	"steampunk":        true,
	"horror":           true,
}

// allowedAssetTypes defines the valid assetType values (case-sensitive).
var allowedAssetTypes = map[string]bool{
	"character": true,
	"location":  true,
	"item":      true,
	"creature":  true,
}

// allowedMoods defines the valid mood values (case-sensitive).
var allowedMoods = map[string]bool{
	"dark":       true,
	"epic":       true,
	"mysterious": true,
	"whimsical":  true,
}

// validateFields checks that all required fields are present, non-empty after trim,
// and that genre, assetType, and mood have valid values.
// Returns nil if valid, or the appropriate *constant.AppError if validation fails.
func (s *ImageService) validateFields(req request.GenerateImageRequest) *constant.AppError {
	// Check required fields are non-empty after trim
	if strings.TrimSpace(req.ProjectId) == "" {
		return &constant.AppError{
			Code:       constant.ErrInvalidPrompt.Code,
			Message:    "projectId is required",
			HttpStatus: constant.ErrInvalidPrompt.HttpStatus,
		}
	}
	if strings.TrimSpace(req.Prompt) == "" {
		return &constant.AppError{
			Code:       constant.ErrInvalidPrompt.Code,
			Message:    "prompt is required",
			HttpStatus: constant.ErrInvalidPrompt.HttpStatus,
		}
	}
	if strings.TrimSpace(req.Genre) == "" {
		return &constant.AppError{
			Code:       constant.ErrInvalidPrompt.Code,
			Message:    "genre is required",
			HttpStatus: constant.ErrInvalidPrompt.HttpStatus,
		}
	}
	if strings.TrimSpace(req.AssetType) == "" {
		return &constant.AppError{
			Code:       constant.ErrInvalidPrompt.Code,
			Message:    "assetType is required",
			HttpStatus: constant.ErrInvalidPrompt.HttpStatus,
		}
	}
	if strings.TrimSpace(req.Mood) == "" {
		return &constant.AppError{
			Code:       constant.ErrInvalidPrompt.Code,
			Message:    "mood is required",
			HttpStatus: constant.ErrInvalidPrompt.HttpStatus,
		}
	}

	// Validate genre enum (case-sensitive)
	if !allowedGenres[req.Genre] {
		return &constant.AppError{
			Code:       constant.ErrInvalidPrompt.Code,
			Message:    "invalid genre value",
			HttpStatus: constant.ErrInvalidPrompt.HttpStatus,
		}
	}

	// Validate assetType enum (case-sensitive)
	if !allowedAssetTypes[req.AssetType] {
		return &constant.AppError{
			Code:       constant.ErrInvalidPrompt.Code,
			Message:    "invalid assetType value",
			HttpStatus: constant.ErrInvalidPrompt.HttpStatus,
		}
	}

	// Validate mood enum (case-sensitive)
	if !allowedMoods[req.Mood] {
		return &constant.AppError{
			Code:       constant.ErrInvalidPrompt.Code,
			Message:    "invalid mood value",
			HttpStatus: constant.ErrInvalidPrompt.HttpStatus,
		}
	}

	return nil
}

// verifyProjectOwnership checks that the given projectId belongs to the authenticated user.
// Returns nil if valid, or *constant.AppError if the project doesn't exist or doesn't belong to the user.
func (s *ImageService) verifyProjectOwnership(requestId, projectId, userId string) *constant.AppError {
	if s.projectRepo == nil {
		logger.Error(requestId, "project repository not configured")
		return &constant.ErrInternalServer
	}

	_, err := s.projectRepo.FindById(projectId, userId)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Warn(requestId, "project not found for image generation", "projectId", projectId, "userId", userId)
			return &constant.ErrProjectNotFound
		}
		logger.Error(requestId, "failed to verify project ownership", err)
		return &constant.ErrDatabaseFailed
	}

	return nil
}

// Generate implements the full image generation flow:
// validate prompt → validate fields → verify project ownership → check cache → insert processing record → spawn goroutine → return 202
// Returns a cache hit response (HTTP 200 equivalent) or a new generation response (HTTP 202 equivalent).
func (s *ImageService) Generate(requestId, userId string, req request.GenerateImageRequest) (response.GenerateImageResponse, *constant.AppError, bool) {
	// Step 1: Validate prompt (must pass before any DB/API interaction)
	if err := s.validatePrompt(req.Prompt); err != nil {
		logger.Warn(requestId, "prompt validation failed", "prompt", req.Prompt)
		return response.GenerateImageResponse{}, err, false
	}

	// Step 2: Validate fields (genre, assetType, mood)
	if err := s.validateFields(req); err != nil {
		logger.Warn(requestId, "field validation failed")
		return response.GenerateImageResponse{}, err, false
	}

	// Step 2.5: Validate reference image (if provided)
	var refImageBytes []byte
	var refMimeType string
	if req.ReferenceImage != "" {
		var appErr *constant.AppError
		refImageBytes, refMimeType, appErr = s.validateReferenceImage(req.ReferenceImage)
		if appErr != nil {
			logger.Warn(requestId, "reference image validation failed")
			return response.GenerateImageResponse{}, appErr, false
		}
	}

	// Step 3: Verify project ownership
	if err := s.verifyProjectOwnership(requestId, req.ProjectId, userId); err != nil {
		return response.GenerateImageResponse{}, err, false
	}

	// Step 4: Compute hash and check cache
	hash := s.hashPrompt(req.Prompt, req.Genre, req.AssetType, req.Mood)
	cached, _ := s.imageRepo.FindByPromptHash(hash)
	if cached != nil && cached.Status == "completed" {
		logger.Info(requestId, "cache hit, returning existing image", "imageId", cached.Id)
		data := mapImageToData(cached)
		return response.GenerateImageResponse{
			BaseResponse: response.BaseResponse{
				RequestId: requestId,
				Success:   true,
				Message:   "Image retrieved from cache",
			},
			Data: &data,
		}, nil, true // true = cache hit
	}

	// Step 5: Insert processing record
	imageId, err := s.imageRepo.InsertProcessing(userId, req, hash, req.ReferenceImage != "")
	if err != nil {
		logger.Error(requestId, "failed to insert processing record", err)
		return response.GenerateImageResponse{}, &constant.ErrDatabaseFailed, false
	}

	// Step 6: Spawn async goroutine for generation (pass reference image bytes)
	go s.processGeneration(requestId, imageId, req, refImageBytes, refMimeType)

	// Step 7: Return 202 immediately
	logger.Info(requestId, "image generation started", "imageId", imageId)
	return response.GenerateImageResponse{
		BaseResponse: response.BaseResponse{
			RequestId: requestId,
			Success:   true,
			Message:   "Image generation started",
		},
		Data: &response.ImageData{
			Id:        imageId,
			Status:    "processing",
			Prompt:    req.Prompt,
			Genre:     req.Genre,
			AssetType: req.AssetType,
			Mood:      req.Mood,
		},
	}, nil, false // false = new generation (not cache hit)
}

// processGeneration runs in a background goroutine. It calls the Gemini API,
// saves the result to storage, and updates the database status.
// It always terminates and updates the status exactly once.
// refImageBytes and refMimeType are the decoded reference image data (both empty if no reference image).
func (s *ImageService) processGeneration(requestId, imageId string, req request.GenerateImageRequest, refImageBytes []byte, refMimeType string) {
	// Create independent context with 30s timeout (not tied to HTTP request context)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Build the prompt string for Gemini
	prompt := buildPrompt(req.AssetType, req.Prompt, req.Genre, req.Mood)
	logger.Info(requestId, "calling Gemini API", "imageId", imageId, "prompt", prompt)

	// Construct reference image for Gemini (nil if no reference provided)
	var refImage *gemini.ReferenceImage
	if len(refImageBytes) > 0 {
		refImage = &gemini.ReferenceImage{
			Data:     refImageBytes,
			MimeType: refMimeType,
		}
	}

	// Call Gemini API
	imageBytes, err := s.geminiClient.Generate(ctx, prompt, refImage)

	// Release reference image bytes to help GC
	refImageBytes = nil
	refImage = nil

	if err != nil {
		logger.Error(requestId, "Gemini API call failed", "imageId", imageId, err)
		if ctx.Err() == context.DeadlineExceeded {
			s.imageRepo.UpdateStatus(imageId, "failed", constant.ErrGeminiTimeout.Code)
		} else {
			s.imageRepo.UpdateStatus(imageId, "failed", constant.ErrGeminiBadResponse.Code)
		}
		return
	}

	// Validate response contains actual image data
	if len(imageBytes) == 0 {
		logger.Error(requestId, "Gemini returned empty image data", "imageId", imageId)
		s.imageRepo.UpdateStatus(imageId, "failed", constant.ErrGeminiBadResponse.Code)
		return
	}

	// Save to filesystem
	filePath, err := s.storageSaver.Save(imageId, imageBytes)
	if err != nil {
		logger.Error(requestId, "filesystem save failed", "imageId", imageId, err)
		s.imageRepo.UpdateStatus(imageId, "failed", constant.ErrStorageFailed.Code)
		return
	}

	// Mark completed
	if err := s.imageRepo.UpdateCompleted(imageId, filePath); err != nil {
		logger.Error(requestId, "failed to update image status to completed", "imageId", imageId, err)
		return
	}

	logger.Info(requestId, "image generation completed", "imageId", imageId, "filePath", filePath)
}

// buildPrompt constructs the prompt string sent to the Gemini API.
func buildPrompt(assetType, prompt, genre, mood string) string {
	return fmt.Sprintf("Generate a %s concept art: %s. Style: %s. Mood: %s.", assetType, prompt, genre, mood)
}

// detectMimeType inspects magic bytes of the decoded image data to determine
// the MIME type. Returns the detected MIME type and true if supported,
// or ("", false) if the format is unrecognized or data is too short.
// Detection is based solely on magic bytes, never file extension.
func detectMimeType(data []byte) (string, bool) {
	if len(data) < 12 {
		return "", false
	}

	// JPEG: starts with FF D8 FF
	if data[0] == 0xFF && data[1] == 0xD8 && data[2] == 0xFF {
		return "image/jpeg", true
	}

	// PNG: starts with 89 50 4E 47 (‰PNG)
	if data[0] == 0x89 && data[1] == 0x50 && data[2] == 0x4E && data[3] == 0x47 {
		return "image/png", true
	}

	// WEBP: RIFF....WEBP
	if data[0] == 'R' && data[1] == 'I' && data[2] == 'F' && data[3] == 'F' &&
		data[8] == 'W' && data[9] == 'E' && data[10] == 'B' && data[11] == 'P' {
		return "image/webp", true
	}

	return "", false
}

// validateReferenceImage decodes a base64-encoded reference image, validates
// its size (≤ 5MB), and detects the MIME type from magic bytes.
// Returns decoded bytes, detected MIME type, and nil on success.
// Returns ERR_04 if decoding fails, size exceeds limit, or format is unsupported.
func (s *ImageService) validateReferenceImage(referenceImage string) ([]byte, string, *constant.AppError) {
	// Step 1: Decode base64 (standard encoding first, URL-safe as fallback)
	decoded, err := base64.StdEncoding.DecodeString(referenceImage)
	if err != nil {
		decoded, err = base64.URLEncoding.DecodeString(referenceImage)
		if err != nil {
			return nil, "", &constant.ErrInvalidPrompt
		}
	}

	// Step 2: Check size limit (5MB = 5,242,880 bytes)
	const maxSize = 5 * 1024 * 1024
	if len(decoded) > maxSize {
		return nil, "", &constant.ErrInvalidPrompt
	}

	// Step 3: Detect MIME type from magic bytes
	mimeType, valid := detectMimeType(decoded)
	if !valid {
		return nil, "", &constant.ErrInvalidPrompt
	}

	return decoded, mimeType, nil
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
