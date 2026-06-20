package repository

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

var AllowedProductCategories = map[string]bool{
	"general":     true,
	"food":        true,
	"beverage":    true,
	"cosmetics":   true,
	"fashion":     true,
	"electronics": true,
	"home":        true,
}

type Project struct {
	Id              string
	UserId          string
	Name            string
	Description     string
	ProductCategory string
	VisualStyle     string
	CreatedAt       time.Time
}

func (p *Project) Genre() string    { return p.ProductCategory }
func (p *Project) ArtStyle() string { return p.VisualStyle }

type ProjectRepository struct {
	db *sql.DB
}

func NewProjectRepository(db *sql.DB) *ProjectRepository {
	return &ProjectRepository{db: db}
}

func (r *ProjectRepository) Insert(userId, name, description, productCategory, visualStyle string) (*Project, error) {
	if productCategory == "" {
		productCategory = "general"
	}

	id := uuid.New().String()
	var createdAt time.Time
	err := r.db.QueryRow(
		`INSERT INTO projects (id, user_id, name, description, product_category, visual_style) VALUES ($1, $2, $3, $4, $5, $6) RETURNING created_at`,
		id, userId, name, description, productCategory, visualStyle,
	).Scan(&createdAt)
	if err != nil {
		return nil, err
	}

	return &Project{
		Id:              id,
		UserId:          userId,
		Name:            name,
		Description:     description,
		ProductCategory: productCategory,
		VisualStyle:     visualStyle,
		CreatedAt:       createdAt,
	}, nil
}

func (r *ProjectRepository) FindById(projectId, userId string) (*Project, error) {
	project := &Project{}
	err := r.db.QueryRow(
		`SELECT id, user_id, name, description, product_category, visual_style, created_at FROM projects WHERE id = $1 AND user_id = $2`,
		projectId, userId,
	).Scan(&project.Id, &project.UserId, &project.Name, &project.Description, &project.ProductCategory, &project.VisualStyle, &project.CreatedAt)
	if err != nil {
		return nil, err
	}
	return project, nil
}

func (r *ProjectRepository) FindByUserId(userId string) ([]Project, error) {
	rows, err := r.db.Query(
		`SELECT id, user_id, name, description, product_category, visual_style, created_at FROM projects WHERE user_id = $1 ORDER BY created_at DESC`,
		userId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []Project
	for rows.Next() {
		var p Project
		if err := rows.Scan(&p.Id, &p.UserId, &p.Name, &p.Description, &p.ProductCategory, &p.VisualStyle, &p.CreatedAt); err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}
	return projects, rows.Err()
}

func (r *ProjectRepository) Delete(projectId, userId string) error {
	result, err := r.db.Exec(`DELETE FROM projects WHERE id = $1 AND user_id = $2`, projectId, userId)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *ProjectRepository) FindScenesByProjectId(projectId string) ([]string, error) {
	rows, err := r.db.Query(`SELECT id FROM scenes WHERE project_id = $1`, projectId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}
