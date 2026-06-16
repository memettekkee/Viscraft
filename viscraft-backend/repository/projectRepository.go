package repository

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// Project represents a row in the projects table.
type Project struct {
	Id          string
	UserId      string
	Name        string
	Description string
	CreatedAt   time.Time
}

// ProjectRepository handles database operations for the projects table.
type ProjectRepository struct {
	db *sql.DB
}

// NewProjectRepository creates a new ProjectRepository with the given database connection.
func NewProjectRepository(db *sql.DB) *ProjectRepository {
	return &ProjectRepository{db: db}
}

// Insert creates a new project and returns the created project.
func (r *ProjectRepository) Insert(userId, name, description string) (*Project, error) {
	id := uuid.New().String()

	var createdAt time.Time
	err := r.db.QueryRow(
		`INSERT INTO projects (id, user_id, name, description) VALUES ($1, $2, $3, $4) RETURNING created_at`,
		id, userId, name, description,
	).Scan(&createdAt)
	if err != nil {
		return nil, err
	}

	return &Project{
		Id:          id,
		UserId:      userId,
		Name:        name,
		Description: description,
		CreatedAt:   createdAt,
	}, nil
}

// FindById retrieves a project by its ID, always filtering by user_id for ownership enforcement.
func (r *ProjectRepository) FindById(projectId, userId string) (*Project, error) {
	project := &Project{}
	err := r.db.QueryRow(
		`SELECT id, user_id, name, description, created_at FROM projects WHERE id = $1 AND user_id = $2`,
		projectId, userId,
	).Scan(&project.Id, &project.UserId, &project.Name, &project.Description, &project.CreatedAt)
	if err != nil {
		return nil, err
	}
	return project, nil
}

// FindByUserId retrieves all projects for a given user.
func (r *ProjectRepository) FindByUserId(userId string) ([]Project, error) {
	rows, err := r.db.Query(
		`SELECT id, user_id, name, description, created_at FROM projects WHERE user_id = $1 ORDER BY created_at DESC`,
		userId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []Project
	for rows.Next() {
		var p Project
		if err := rows.Scan(&p.Id, &p.UserId, &p.Name, &p.Description, &p.CreatedAt); err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return projects, nil
}

// Delete removes a project by its ID with user_id ownership check.
// Cascade delete rules in the database schema handle removal of associated images.
func (r *ProjectRepository) Delete(projectId, userId string) error {
	result, err := r.db.Exec(
		`DELETE FROM projects WHERE id = $1 AND user_id = $2`,
		projectId, userId,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// FindImagesByProjectId retrieves all image IDs belonging to a project.
// Used for filesystem cleanup during project deletion.
func (r *ProjectRepository) FindImagesByProjectId(projectId string) ([]string, error) {
	rows, err := r.db.Query(
		`SELECT id FROM images WHERE project_id = $1`,
		projectId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var imageIds []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		imageIds = append(imageIds, id)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return imageIds, nil
}
