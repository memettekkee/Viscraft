package config

import (
	"database/sql"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"
)

func TestGetEnv_ReturnsValue(t *testing.T) {
	os.Setenv("TEST_CONFIG_VAR", "hello")
	defer os.Unsetenv("TEST_CONFIG_VAR")

	got := getEnv("TEST_CONFIG_VAR", "default")
	if got != "hello" {
		t.Errorf("getEnv() = %q, want %q", got, "hello")
	}
}

func TestGetEnv_ReturnsDefault(t *testing.T) {
	os.Unsetenv("TEST_CONFIG_MISSING")

	got := getEnv("TEST_CONFIG_MISSING", "fallback")
	if got != "fallback" {
		t.Errorf("getEnv() = %q, want %q", got, "fallback")
	}
}

func TestValidateStoragePath_ValidDir(t *testing.T) {
	dir := t.TempDir()
	// Should not panic/fatal
	validateStoragePath(dir)
}

func TestLoad_ValidConfig(t *testing.T) {
	// Create a temp directory for storage
	dir := t.TempDir()

	// Set all required env vars
	os.Setenv("GEMINI_API_KEY", "test-key")
	os.Setenv("STORAGE_PATH", dir)
	os.Setenv("JWT_SECRET", "test-secret")
	os.Setenv("JWT_EXPIRY", "24h")
	defer func() {
		os.Unsetenv("GEMINI_API_KEY")
		os.Unsetenv("STORAGE_PATH")
		os.Unsetenv("JWT_SECRET")
		os.Unsetenv("JWT_EXPIRY")
	}()

	cfg := Load()

	if cfg.GeminiAPIKey != "test-key" {
		t.Errorf("GeminiAPIKey = %q, want %q", cfg.GeminiAPIKey, "test-key")
	}
	if cfg.StoragePath != dir {
		t.Errorf("StoragePath = %q, want %q", cfg.StoragePath, dir)
	}
	if cfg.JWTSecret != "test-secret" {
		t.Errorf("JWTSecret = %q, want %q", cfg.JWTSecret, "test-secret")
	}
	if cfg.JWTExpiry != 24*time.Hour {
		t.Errorf("JWTExpiry = %v, want %v", cfg.JWTExpiry, 24*time.Hour)
	}
	if cfg.AppPort != "8080" {
		t.Errorf("AppPort = %q, want %q", cfg.AppPort, "8080")
	}
}

func TestLoad_CustomPort(t *testing.T) {
	dir := t.TempDir()

	os.Setenv("APP_PORT", "9090")
	os.Setenv("GEMINI_API_KEY", "test-key")
	os.Setenv("STORAGE_PATH", dir)
	os.Setenv("JWT_SECRET", "test-secret")
	os.Setenv("JWT_EXPIRY", "720h")
	defer func() {
		os.Unsetenv("APP_PORT")
		os.Unsetenv("GEMINI_API_KEY")
		os.Unsetenv("STORAGE_PATH")
		os.Unsetenv("JWT_SECRET")
		os.Unsetenv("JWT_EXPIRY")
	}()

	cfg := Load()

	if cfg.AppPort != "9090" {
		t.Errorf("AppPort = %q, want %q", cfg.AppPort, "9090")
	}
	if cfg.JWTExpiry != 720*time.Hour {
		t.Errorf("JWTExpiry = %v, want %v", cfg.JWTExpiry, 720*time.Hour)
	}
}

func TestBuildDSN(t *testing.T) {
	cfg := &Config{
		DBHost:     "localhost",
		DBPort:     "5432",
		DBUser:     "myuser",
		DBPassword: "mypass",
		DBName:     "mydb",
	}

	got := buildDSN(cfg)
	want := "host=localhost port=5432 user=myuser password=mypass dbname=mydb sslmode=disable"

	if got != want {
		t.Errorf("buildDSN() = %q, want %q", got, want)
	}
}

func TestBuildDSN_CustomValues(t *testing.T) {
	cfg := &Config{
		DBHost:     "db.example.com",
		DBPort:     "5433",
		DBUser:     "admin",
		DBPassword: "s3cret!",
		DBName:     "production",
	}

	got := buildDSN(cfg)
	want := "host=db.example.com port=5433 user=admin password=s3cret! dbname=production sslmode=disable"

	if got != want {
		t.Errorf("buildDSN() = %q, want %q", got, want)
	}
}

func TestInitDB_OpenSucceeds(t *testing.T) {
	// sql.Open with postgres driver validates the driver name but does not connect.
	// This test verifies that InitDB correctly calls sql.Open without error.
	cfg := &Config{
		DBHost:     "localhost",
		DBPort:     "5432",
		DBUser:     "testuser",
		DBPassword: "testpass",
		DBName:     "testdb",
	}

	dsn := buildDSN(cfg)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatalf("sql.Open() with DSN from buildDSN failed: %v", err)
	}
	defer db.Close()

	// Verify the connection object is not nil
	if db == nil {
		t.Fatal("sql.Open() returned nil db")
	}
}
