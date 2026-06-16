package repository

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// User represents a row in the users table.
type User struct {
	Id        string
	Email     string
	Password  string
	Name      string
	CreatedAt time.Time
}

// UserRepository handles database operations for the users table.
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new UserRepository with the given database connection.
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Insert creates a new user with a bcrypt-hashed password and returns the created user.
// The provided password is hashed before storage — raw passwords are never persisted.
func (r *UserRepository) Insert(email, password, name string) (*User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	id := uuid.New().String()

	var createdAt time.Time
	err = r.db.QueryRow(
		`INSERT INTO users (id, email, password, name) VALUES ($1, $2, $3, $4) RETURNING created_at`,
		id, email, string(hashedPassword), name,
	).Scan(&createdAt)
	if err != nil {
		return nil, err
	}

	return &User{
		Id:        id,
		Email:     email,
		Password:  string(hashedPassword),
		Name:      name,
		CreatedAt: createdAt,
	}, nil
}

// FindByEmail retrieves a user by their email address.
// Returns the user with the hashed password for login comparison.
func (r *UserRepository) FindByEmail(email string) (*User, error) {
	user := &User{}
	err := r.db.QueryRow(
		`SELECT id, email, password, name, created_at FROM users WHERE email = $1`,
		email,
	).Scan(&user.Id, &user.Email, &user.Password, &user.Name, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// FindById retrieves a user by their unique ID.
// Used for auth context lookup after JWT validation.
func (r *UserRepository) FindById(id string) (*User, error) {
	user := &User{}
	err := r.db.QueryRow(
		`SELECT id, email, password, name, created_at FROM users WHERE id = $1`,
		id,
	).Scan(&user.Id, &user.Email, &user.Password, &user.Name, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// Delete removes a user by their ID. Cascade delete rules in the database
// schema handle removal of associated projects and images.
func (r *UserRepository) Delete(userId string) error {
	_, err := r.db.Exec(`DELETE FROM users WHERE id = $1`, userId)
	return err
}
