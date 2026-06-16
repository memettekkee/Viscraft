package repository

import (
	"database/sql"
	"errors"
	"testing"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

// TestInsertHashesPassword verifies that Insert bcrypt-hashes the password
// before storing it (i.e., the returned Password field is not the raw input).
func TestInsertHashesPassword(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db)

	rawPassword := "securepass123"
	user, err := repo.Insert("test@example.com", rawPassword, "Test User")
	if err != nil {
		t.Fatalf("Insert failed: %v", err)
	}

	if user.Password == rawPassword {
		t.Fatal("expected password to be hashed, but raw password was stored")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(rawPassword))
	if err != nil {
		t.Fatalf("stored password is not a valid bcrypt hash of the input: %v", err)
	}
}

// TestInsertReturnsUserFields verifies that Insert returns the user with all expected fields populated.
func TestInsertReturnsUserFields(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db)

	user, err := repo.Insert("alice@example.com", "password123", "Alice")
	if err != nil {
		t.Fatalf("Insert failed: %v", err)
	}

	if user.Id == "" {
		t.Error("expected non-empty Id")
	}
	if user.Email != "alice@example.com" {
		t.Errorf("expected email alice@example.com, got %s", user.Email)
	}
	if user.Name != "Alice" {
		t.Errorf("expected name Alice, got %s", user.Name)
	}
	if user.CreatedAt.IsZero() {
		t.Error("expected non-zero CreatedAt")
	}
}

// TestInsertDuplicateEmail verifies that inserting a user with a duplicate email returns an error.
func TestInsertDuplicateEmail(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db)

	_, err := repo.Insert("dup@example.com", "password123", "First")
	if err != nil {
		t.Fatalf("first Insert failed: %v", err)
	}

	_, err = repo.Insert("dup@example.com", "password456", "Second")
	if err == nil {
		t.Fatal("expected error for duplicate email, got nil")
	}
}

// TestFindByEmail verifies that a user can be retrieved by email after insertion.
func TestFindByEmail(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db)

	inserted, err := repo.Insert("find@example.com", "password123", "Finder")
	if err != nil {
		t.Fatalf("Insert failed: %v", err)
	}

	found, err := repo.FindByEmail("find@example.com")
	if err != nil {
		t.Fatalf("FindByEmail failed: %v", err)
	}

	if found.Id != inserted.Id {
		t.Errorf("expected Id %s, got %s", inserted.Id, found.Id)
	}
	if found.Email != "find@example.com" {
		t.Errorf("expected email find@example.com, got %s", found.Email)
	}
}

// TestFindByEmailNotFound verifies that FindByEmail returns sql.ErrNoRows for a non-existent email.
func TestFindByEmailNotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db)

	_, err := repo.FindByEmail("nonexistent@example.com")
	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("expected sql.ErrNoRows, got: %v", err)
	}
}

// TestFindById verifies that a user can be retrieved by ID after insertion.
func TestFindById(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db)

	inserted, err := repo.Insert("byid@example.com", "password123", "ByID")
	if err != nil {
		t.Fatalf("Insert failed: %v", err)
	}

	found, err := repo.FindById(inserted.Id)
	if err != nil {
		t.Fatalf("FindById failed: %v", err)
	}

	if found.Email != "byid@example.com" {
		t.Errorf("expected email byid@example.com, got %s", found.Email)
	}
	if found.Name != "ByID" {
		t.Errorf("expected name ByID, got %s", found.Name)
	}
}

// TestFindByIdNotFound verifies that FindById returns sql.ErrNoRows for a non-existent ID.
func TestFindByIdNotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db)

	_, err := repo.FindById("00000000-0000-0000-0000-000000000000")
	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("expected sql.ErrNoRows, got: %v", err)
	}
}

// TestDelete verifies that deleting a user makes them unfindable.
func TestDelete(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db)

	inserted, err := repo.Insert("delete@example.com", "password123", "ToDelete")
	if err != nil {
		t.Fatalf("Insert failed: %v", err)
	}

	err = repo.Delete(inserted.Id)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	_, err = repo.FindById(inserted.Id)
	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("expected sql.ErrNoRows after delete, got: %v", err)
	}
}

// TestDeleteNonExistentUser verifies that deleting a non-existent user does not error.
func TestDeleteNonExistentUser(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db)

	err := repo.Delete("00000000-0000-0000-0000-000000000000")
	if err != nil {
		t.Fatalf("expected no error deleting non-existent user, got: %v", err)
	}
}

// setupTestDB connects to a PostgreSQL test database and creates the users table.
// Tests require a running PostgreSQL instance. Set DATABASE_URL env var or use the default.
func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()

	dsn := "host=localhost port=5432 user=postgres password=postgres dbname=viscraft_test sslmode=disable"

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Skipf("skipping test: cannot open database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		t.Skipf("skipping test: cannot connect to database: %v", err)
	}

	// Create users table for tests
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			email       VARCHAR(255) UNIQUE NOT NULL,
			password    VARCHAR(255) NOT NULL,
			name        VARCHAR(255),
			created_at  TIMESTAMP DEFAULT NOW()
		)
	`)
	if err != nil {
		t.Fatalf("failed to create users table: %v", err)
	}

	// Clean up any existing test data
	_, err = db.Exec(`DELETE FROM users`)
	if err != nil {
		t.Fatalf("failed to clean users table: %v", err)
	}

	return db
}
