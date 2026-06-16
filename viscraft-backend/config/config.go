package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// Config holds all environment configuration for the application.
type Config struct {
	AppPort      string
	AppEnv       string
	DBHost       string
	DBPort       string
	DBName       string
	DBUser       string
	DBPassword   string
	GeminiAPIKey string
	GeminiModel  string
	StoragePath  string
	JWTSecret    string
	JWTExpiry    time.Duration
}

// Load reads environment variables from .env (if present) and returns a validated Config.
// It fails with a descriptive error if required fields are missing or invalid.
func Load() *Config {
	// Load .env file if it exists; ignore error if file is absent
	_ = godotenv.Load()

	cfg := &Config{
		AppPort:      getEnv("APP_PORT", "8080"),
		AppEnv:       getEnv("APP_ENV", "development"),
		DBHost:       getEnv("DB_HOST", "localhost"),
		DBPort:       getEnv("DB_PORT", "5432"),
		DBName:       getEnv("DB_NAME", "viscraft"),
		DBUser:       getEnv("DB_USER", "postgres"),
		DBPassword:   getEnv("DB_PASSWORD", ""),
		GeminiAPIKey: os.Getenv("GEMINI_API_KEY"),
		GeminiModel:  getEnv("GEMINI_MODEL", "gemini-2.0-flash-preview-image-generation"),
		StoragePath:  os.Getenv("STORAGE_PATH"),
		JWTSecret:    os.Getenv("JWT_SECRET"),
	}

	// Validate required fields
	if cfg.GeminiAPIKey == "" {
		log.Fatal("GEMINI_API_KEY environment variable is required")
	}
	if cfg.StoragePath == "" {
		log.Fatal("STORAGE_PATH environment variable is required")
	}
	if cfg.JWTSecret == "" {
		log.Fatal("JWT_SECRET environment variable is required")
	}

	// Parse JWT_EXPIRY as a Go duration string (e.g., "24h", "720h")
	jwtExpiryStr := os.Getenv("JWT_EXPIRY")
	if jwtExpiryStr == "" {
		log.Fatal("JWT_EXPIRY environment variable is required")
	}
	duration, err := time.ParseDuration(jwtExpiryStr)
	if err != nil {
		log.Fatalf("JWT_EXPIRY is not a valid duration (e.g., \"24h\", \"720h\"): %v", err)
	}
	cfg.JWTExpiry = duration

	// Validate STORAGE_PATH exists and is writable
	validateStoragePath(cfg.StoragePath)

	return cfg
}

// buildDSN constructs a PostgreSQL connection string from the provided config.
func buildDSN(cfg *Config) string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)
}

// InitDB creates a database connection pool using the provided config.
// It verifies connectivity with a ping and configures pool settings.
func InitDB(cfg *Config) *sql.DB {
	dsn := buildDSN(cfg)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to open database connection: %v", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Verify connectivity
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	return db
}

// validateStoragePath checks that the directory exists and is writable.
func validateStoragePath(path string) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			log.Fatalf("STORAGE_PATH directory does not exist: %s", path)
		}
		log.Fatalf("Failed to stat STORAGE_PATH: %v", err)
	}

	if !info.IsDir() {
		log.Fatalf("STORAGE_PATH is not a directory: %s", path)
	}

	// Check writability by creating and removing a temporary file
	tmpFile := path + "/.viscraft_write_test"
	f, err := os.Create(tmpFile)
	if err != nil {
		log.Fatalf("STORAGE_PATH is not writable: %s (error: %v)", path, err)
	}
	f.Close()
	os.Remove(tmpFile)
}

// getEnv returns the value of an environment variable or a default if not set.
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
