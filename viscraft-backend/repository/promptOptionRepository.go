package repository

import "database/sql"

type PromptOption struct {
	Id          string `json:"id"`
	Category    string `json:"category"`
	Label       string `json:"label"`
	PromptValue string `json:"promptValue"`
	SortOrder   int    `json:"sortOrder"`
}

type PromptOptionRepository struct {
	db *sql.DB
}

func NewPromptOptionRepository(db *sql.DB) *PromptOptionRepository {
	return &PromptOptionRepository{db: db}
}

func (r *PromptOptionRepository) FindAll() ([]PromptOption, error) {
	rows, err := r.db.Query(
		`SELECT id, category, label, prompt_value, sort_order FROM prompt_options ORDER BY category, sort_order`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var options []PromptOption
	for rows.Next() {
		var opt PromptOption
		if err := rows.Scan(&opt.Id, &opt.Category, &opt.Label, &opt.PromptValue, &opt.SortOrder); err != nil {
			return nil, err
		}
		options = append(options, opt)
	}
	return options, nil
}

func (r *PromptOptionRepository) FindByCategory(category string) ([]PromptOption, error) {
	rows, err := r.db.Query(
		`SELECT id, category, label, prompt_value, sort_order FROM prompt_options WHERE category = $1 ORDER BY sort_order`,
		category,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var options []PromptOption
	for rows.Next() {
		var opt PromptOption
		if err := rows.Scan(&opt.Id, &opt.Category, &opt.Label, &opt.PromptValue, &opt.SortOrder); err != nil {
			return nil, err
		}
		options = append(options, opt)
	}
	return options, nil
}
