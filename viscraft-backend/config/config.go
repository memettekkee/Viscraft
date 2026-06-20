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
	AppPort              string
	AppEnv               string
	DBHost               string
	DBPort               string
	DBName               string
	DBUser               string
	DBPassword           string
	PollinationsAPIKey   string
	StoragePath          string
	StoragePublicURL     string
	StorageTempPublicURL string
	JWTSecret            string
	JWTExpiry            time.Duration
}

func Load() *Config {
	_ = godotenv.Load()

	cfg := &Config{
		AppPort:              getEnv("APP_PORT", "8080"),
		AppEnv:               getEnv("APP_ENV", "development"),
		DBHost:               getEnv("DB_HOST", "localhost"),
		DBPort:               getEnv("DB_PORT", "5432"),
		DBName:               getEnv("DB_NAME", "viscraft"),
		DBUser:               getEnv("DB_USER", "postgres"),
		DBPassword:           getEnv("DB_PASSWORD", ""),
		PollinationsAPIKey:   os.Getenv("POLLINATIONS_API_KEY"),
		StoragePath:          os.Getenv("STORAGE_PATH"),
		StoragePublicURL:     os.Getenv("STORAGE_PUBLIC_URL"),
		StorageTempPublicURL: os.Getenv("STORAGE_TEMP_PUBLIC_URL"),
		JWTSecret:            os.Getenv("JWT_SECRET"),
	}

	// Validate required fields
	if cfg.PollinationsAPIKey == "" {
		log.Fatal("POLLINATIONS_API_KEY environment variable is required")
	}
	if cfg.StoragePath == "" {
		log.Fatal("STORAGE_PATH environment variable is required")
	}
	if cfg.StoragePublicURL == "" {
		log.Fatal("STORAGE_PUBLIC_URL environment variable is required")
	}
	if cfg.StorageTempPublicURL == "" {
		log.Fatal("STORAGE_TEMP_PUBLIC_URL environment variable is required")
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
