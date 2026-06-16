package config

import (
	"os"
	"testing"
	"time"
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
