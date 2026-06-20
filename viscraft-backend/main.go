package main

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"viscraft-backend/config"
	"viscraft-backend/handler"
	"viscraft-backend/middleware"
	"viscraft-backend/pkg/imagegen"
	"viscraft-backend/pkg/router"
	"viscraft-backend/pkg/storage"
	"viscraft-backend/repository"
	"viscraft-backend/service"

	"github.com/gin-gonic/gin"
)

// sceneFinderAdapter bridges repository.SceneRepository to the service.ImageFinder
// interface. The repository returns []repository.SceneRecord while the service
// expects []service.ImageRecord — both have an identical shape (Id string) but are
// distinct types to avoid circular imports.
type sceneFinderAdapter struct {
	repo *repository.SceneRepository
}

func (a *sceneFinderAdapter) FindImagesByUserId(userId string) ([]service.ImageRecord, error) {
	records, err := a.repo.FindScenesByUserId(userId)
	if err != nil {
		return nil, err
	}

	result := make([]service.ImageRecord, len(records))
	for i, r := range records {
		result[i] = service.ImageRecord{Id: r.Id}
	}
	return result, nil
}

func main() {
	// Load configuration from environment / .env file
	cfg := config.Load()

	// Initialize database connection pool
	db := config.InitDB(cfg)
	defer db.Close()

	// Initialize external clients
	imagegenClient := imagegen.NewClient(cfg.PollinationsAPIKey)

	// Set up storage paths
	tempBasePath := filepath.Join(filepath.Dir(cfg.StoragePath), "temp")
	if err := os.MkdirAll(tempBasePath, 0755); err != nil {
		log.Fatalf("Failed to create temp storage directory: %v", err)
	}
	localStorage := storage.NewLocalStorage(cfg.StoragePath, tempBasePath, cfg.StoragePublicURL, cfg.StorageTempPublicURL)

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	projectRepo := repository.NewProjectRepository(db)
	sceneRepo := repository.NewSceneRepository(db)
	promptOptionRepo := repository.NewPromptOptionRepository(db)

	// Create the adapter that bridges SceneRepository → service.ImageFinder
	sceneFinderAdapt := &sceneFinderAdapter{repo: sceneRepo}

	// Initialize services with dependency injection
	userService := service.NewUserService(
		userRepo,
		sceneFinderAdapt,
		localStorage,
		cfg.JWTSecret,
		cfg.JWTExpiry,
	)

	projectService := service.NewProjectService(
		projectRepo,
		projectRepo, // satisfies service.ProjectImageFinder directly
		localStorage,
	)

	sceneService := service.NewSceneService(
		sceneRepo,
		projectRepo,    // satisfies service.SceneProjectFinder
		localStorage,   // satisfies service.SceneStorageSaver
		localStorage,   // satisfies service.SceneStorageDeleter
		imagegenClient, // satisfies service.SceneImageGenerator
	)

	// Initialize controllers
	userController := handler.NewUserController(userService)
	projectController := handler.NewProjectController(projectService)
	sceneController := handler.NewSceneController(sceneService)
	promptOptionController := handler.NewPromptOptionController(promptOptionRepo)
	healthController := handler.NewHealthController()

	// Create Gin engine with middleware chain: CORS → RequestID → Timeout → Recovery
	r := gin.New()
	r.Use(middleware.CORS())
	r.Use(middleware.RequestID())
	r.Use(middleware.Timeout(30 * time.Second))
	r.Use(middleware.Recovery())

	// Serve static image files — in Docker this is handled by nginx,
	// but this fallback ensures local development works without nginx.
	r.Static("/storage/images", cfg.StoragePath)
	r.Static("/storage/temp", tempBasePath)

	// Register /scenes/generate separately with JWTAuth + RateLimit middleware
	authMiddleware := middleware.JWTAuth(cfg.JWTSecret)
	rateLimitMiddleware := middleware.RateLimit(5, 60*time.Second)
	r.POST("/scenes/generate", authMiddleware, rateLimitMiddleware, sceneController.Generate)

	// Public endpoint for prompt options (no auth required)
	r.POST("/prompt-options", promptOptionController.ListByCategory)

	// Register all other controllers via the router registry
	err := router.Register(
		r,
		authMiddleware,
		userController,
		projectController,
		sceneController,
		healthController,
	)
	if err != nil {
		log.Fatalf("Failed to register routes: %v", err)
	}

	// Start server
	log.Printf("Viscraft backend starting on port %s", cfg.AppPort)
	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
