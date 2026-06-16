package main

import (
	"log"
	"time"

	"viscraft-backend/config"
	"viscraft-backend/handler"
	"viscraft-backend/middleware"
	"viscraft-backend/pkg/gemini"
	"viscraft-backend/pkg/router"
	"viscraft-backend/pkg/storage"
	"viscraft-backend/repository"
	"viscraft-backend/service"

	"github.com/gin-gonic/gin"
)

// imageFinderAdapter bridges repository.ImageRepository to the service.ImageFinder
// interface. The repository returns []repository.ImageRecord while the service
// expects []service.ImageRecord — both have an identical shape (Id string) but are
// distinct types to avoid circular imports.
type imageFinderAdapter struct {
	repo *repository.ImageRepository
}

func (a *imageFinderAdapter) FindImagesByUserId(userId string) ([]service.ImageRecord, error) {
	records, err := a.repo.FindImagesByUserId(userId)
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
	geminiClient := gemini.NewClient(cfg.GeminiAPIKey, cfg.GeminiModel)
	localStorage := storage.NewLocalStorage(cfg.StoragePath)

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	projectRepo := repository.NewProjectRepository(db)
	imageRepo := repository.NewImageRepository(db)

	// Create the adapter that bridges ImageRepository → service.ImageFinder
	imgFinderAdapter := &imageFinderAdapter{repo: imageRepo}

	// Initialize services with dependency injection
	userService := service.NewUserService(
		userRepo,
		imgFinderAdapter,
		localStorage,
		cfg.JWTSecret,
		cfg.JWTExpiry,
	)

	projectService := service.NewProjectService(
		projectRepo,
		projectRepo, // satisfies service.ProjectImageFinder directly
		localStorage,
	)

	imageService := service.NewImageService(
		imageRepo,
		localStorage, // StorageDeleter
		geminiClient, // GeminiGenerator
		localStorage, // StorageSaver
		projectRepo,  // ProjectOwnershipChecker
	)

	// Initialize controllers
	userController := handler.NewUserController(userService)
	projectController := handler.NewProjectController(projectService)
	imageController := handler.NewImageController(imageService)
	healthController := handler.NewHealthController()

	// Create Gin engine with middleware chain: RequestID → Timeout → Recovery
	r := gin.New()
	r.Use(middleware.RequestID())
	r.Use(middleware.Timeout(30 * time.Second))
	r.Use(middleware.Recovery())

	// Register /images/generate separately with JWTAuth + RateLimit middleware
	authMiddleware := middleware.JWTAuth(cfg.JWTSecret)
	rateLimitMiddleware := middleware.RateLimit(5, 60*time.Second)
	r.POST("/images/generate", authMiddleware, rateLimitMiddleware, imageController.Generate)

	// Register all other controllers via the router registry
	err := router.Register(
		r,
		authMiddleware,
		userController,
		projectController,
		imageController,
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
