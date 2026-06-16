package service

import (
	"context"
	"database/sql"
	"math/rand"
	"strings"
	"sync/atomic"
	"testing"
	"testing/quick"
	"time"

	"viscraft-backend/constant"
	"viscraft-backend/model/request"
	"viscraft-backend/pkg/gemini"
	"viscraft-backend/repository"
)

// mockGeminiClient tracks whether Generate was called.
type mockGeminiClient struct {
	callCount atomic.Int64
}

func (m *mockGeminiClient) Generate(ctx context.Context, prompt string, refImage *gemini.ReferenceImage) ([]byte, error) {
	m.callCount.Add(1)
	return []byte("fake-image"), nil
}

// mockImageRepo tracks whether InsertProcessing was called.
type mockImageRepo struct {
	insertCount atomic.Int64
}

func (m *mockImageRepo) FindById(imageId, userId string) (*repository.Image, error) {
	return nil, sql.ErrNoRows
}

func (m *mockImageRepo) FindByProjectId(projectId, userId string) ([]repository.Image, error) {
	return nil, nil
}

func (m *mockImageRepo) FindByPromptHash(hash string) (*repository.Image, error) {
	return nil, sql.ErrNoRows
}

func (m *mockImageRepo) InsertProcessing(userId string, req request.GenerateImageRequest, hash string, usedReferenceImage bool) (string, error) {
	m.insertCount.Add(1)
	return "img-123", nil
}

func (m *mockImageRepo) UpdateStatus(imageId, status, errorCode string) error {
	return nil
}

func (m *mockImageRepo) UpdateCompleted(imageId, filePath string) error {
	return nil
}

func (m *mockImageRepo) Delete(imageId, userId string) error {
	return nil
}

func (m *mockImageRepo) FindImagesByUserId(userId string) ([]repository.ImageRecord, error) {
	return nil, nil
}

// mockStorageSaver tracks whether Save was called.
type mockStorageSaver struct {
	saveCount atomic.Int64
}

func (m *mockStorageSaver) Save(imageId string, data []byte) (string, error) {
	m.saveCount.Add(1)
	return "/storage/" + imageId + ".png", nil
}

func (m *mockStorageSaver) Delete(imageId string) error {
	return nil
}

// mockProjectRepo that always confirms ownership.
type mockProjectRepo struct{}

func (m *mockProjectRepo) FindById(projectId, userId string) (*repository.Project, error) {
	return &repository.Project{
		Id:        projectId,
		UserId:    userId,
		Name:      "Test Project",
		CreatedAt: time.Now(),
	}, nil
}

// generateInvalidPrompt creates a prompt that will fail validation.
// It randomly chooses one of: too short, too long, or contains a blocked word.
func generateInvalidPrompt(r *rand.Rand) string {
	strategy := r.Intn(3)
	switch strategy {
	case 0:
		// Too short (0-2 chars after trim)
		length := r.Intn(3) // 0, 1, or 2
		return strings.Repeat("a", length)
	case 1:
		// Too long (301+ chars)
		length := 301 + r.Intn(200)
		return strings.Repeat("x", length)
	case 2:
		// Contains blocked word
		blocked := blockedWords[r.Intn(len(blockedWords))]
		prefix := strings.Repeat("a", 3+r.Intn(10))
		return prefix + " " + blocked + " " + prefix
	}
	return ""
}

// TestProperty_PromptValidationGate verifies Property 3: Prompt Validation Gate.
// For any invalid prompt, the system SHALL NOT:
// - Write to the database (no InsertProcessing call)
// - Spawn a goroutine (no Gemini API call)
// - Call the Gemini API
// Validates: Requirements 6.4
func TestProperty_PromptValidationGate(t *testing.T) {
	gemini := &mockGeminiClient{}
	storage := &mockStorageSaver{}
	projectRepo := &mockProjectRepo{}

	// We need a real ImageRepository to pass, but since validation fails before
	// any repo call, we use a mock that would track calls if they happened.
	// For this test, we create the service with nil imageRepo since the
	// validation should fail before any repo interaction.
	svc := &ImageService{
		imageRepo:    nil, // Would panic if accessed - proving no DB interaction
		storage:      storage,
		storageSaver: storage,
		geminiClient: gemini,
		projectRepo:  projectRepo,
	}

	f := func(seed int64) bool {
		r := rand.New(rand.NewSource(seed))
		invalidPrompt := generateInvalidPrompt(r)

		// Reset counters
		gemini.callCount.Store(0)
		storage.saveCount.Store(0)

		req := request.GenerateImageRequest{
			ProjectId: "project-123",
			Prompt:    invalidPrompt,
			Genre:     "fantasy",
			AssetType: "character",
			Mood:      "dark",
		}

		_, appErr, _ := svc.Generate("test-request-id", "user-123", req)

		// Validation must fail
		if appErr == nil {
			t.Logf("expected error for invalid prompt %q, got nil", invalidPrompt)
			return false
		}

		// Error code must be ERR_04 (invalid prompt)
		if appErr.Code != constant.ErrInvalidPrompt.Code {
			t.Logf("expected error code %s, got %s for prompt %q",
				constant.ErrInvalidPrompt.Code, appErr.Code, invalidPrompt)
			return false
		}

		// No Gemini call should have been made
		if gemini.callCount.Load() != 0 {
			t.Logf("Gemini API was called for invalid prompt %q", invalidPrompt)
			return false
		}

		// No storage save should have been made
		if storage.saveCount.Load() != 0 {
			t.Logf("Storage save was called for invalid prompt %q", invalidPrompt)
			return false
		}

		return true
	}

	cfg := &quick.Config{MaxCount: 200}
	if err := quick.Check(f, cfg); err != nil {
		t.Errorf("Property 3 (Prompt Validation Gate) violated: %v", err)
	}
}

// TestProperty_PromptValidationGate_FieldValidation verifies that invalid fields
// also prevent side effects (no DB insert, no Gemini call, no storage write).
func TestProperty_PromptValidationGate_FieldValidation(t *testing.T) {
	gemini := &mockGeminiClient{}
	storage := &mockStorageSaver{}
	projectRepo := &mockProjectRepo{}

	svc := &ImageService{
		imageRepo:    nil, // Would panic if accessed
		storage:      storage,
		storageSaver: storage,
		geminiClient: gemini,
		projectRepo:  projectRepo,
	}

	invalidRequests := []struct {
		name string
		req  request.GenerateImageRequest
	}{
		{
			"invalid genre",
			request.GenerateImageRequest{
				ProjectId: "project-1",
				Prompt:    "a valid prompt here",
				Genre:     "invalid-genre",
				AssetType: "character",
				Mood:      "dark",
			},
		},
		{
			"invalid assetType",
			request.GenerateImageRequest{
				ProjectId: "project-1",
				Prompt:    "a valid prompt here",
				Genre:     "fantasy",
				AssetType: "invalid-type",
				Mood:      "dark",
			},
		},
		{
			"invalid mood",
			request.GenerateImageRequest{
				ProjectId: "project-1",
				Prompt:    "a valid prompt here",
				Genre:     "fantasy",
				AssetType: "character",
				Mood:      "invalid-mood",
			},
		},
		{
			"empty genre",
			request.GenerateImageRequest{
				ProjectId: "project-1",
				Prompt:    "a valid prompt here",
				Genre:     "",
				AssetType: "character",
				Mood:      "dark",
			},
		},
		{
			"empty assetType",
			request.GenerateImageRequest{
				ProjectId: "project-1",
				Prompt:    "a valid prompt here",
				Genre:     "fantasy",
				AssetType: "",
				Mood:      "dark",
			},
		},
		{
			"empty mood",
			request.GenerateImageRequest{
				ProjectId: "project-1",
				Prompt:    "a valid prompt here",
				Genre:     "fantasy",
				AssetType: "character",
				Mood:      "",
			},
		},
	}

	for _, tt := range invalidRequests {
		t.Run(tt.name, func(t *testing.T) {
			gemini.callCount.Store(0)
			storage.saveCount.Store(0)

			_, appErr, _ := svc.Generate("test-request-id", "user-123", tt.req)

			if appErr == nil {
				t.Errorf("expected error for %s, got nil", tt.name)
				return
			}

			if appErr.Code != constant.ErrInvalidPrompt.Code {
				t.Errorf("expected error code %s, got %s", constant.ErrInvalidPrompt.Code, appErr.Code)
			}

			if gemini.callCount.Load() != 0 {
				t.Errorf("Gemini API was called for invalid request: %s", tt.name)
			}

			if storage.saveCount.Load() != 0 {
				t.Errorf("Storage save was called for invalid request: %s", tt.name)
			}
		})
	}
}
