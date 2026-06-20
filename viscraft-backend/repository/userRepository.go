package repository

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id            string
	Email         string
	Password      string
	Name          string
	CreatedAt     time.Time
	TourCompleted bool
}

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

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
		Id:            id,
		Email:         email,
		Password:      string(hashedPassword),
		Name:          name,
		CreatedAt:     createdAt,
		TourCompleted: false,
	}, nil
}

func (r *UserRepository) FindByEmail(email string) (*User, error) {
	user := &User{}
	err := r.db.QueryRow(
		`SELECT id, email, password, name, created_at, tour_completed FROM users WHERE email = $1`,
		email,
	).Scan(&user.Id, &user.Email, &user.Password, &user.Name, &user.CreatedAt, &user.TourCompleted)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) FindById(id string) (*User, error) {
	user := &User{}
	err := r.db.QueryRow(
		`SELECT id, email, password, name, created_at, tour_completed FROM users WHERE id = $1`,
		id,
	).Scan(&user.Id, &user.Email, &user.Password, &user.Name, &user.CreatedAt, &user.TourCompleted)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) CompleteTour(userId string) error {
	_, err := r.db.Exec(`UPDATE users SET tour_completed = TRUE WHERE id = $1`, userId)
	return err
}

func (r *UserRepository) Delete(userId string) error {
	_, err := r.db.Exec(`DELETE FROM users WHERE id = $1`, userId)
	return err
}
