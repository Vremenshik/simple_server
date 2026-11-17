package repository

import (
	"database/sql"
	"go-expenses-stage5/internal/models"
)

type CategoryRepository struct{ db *sql.DB }

func NewCategoryRepository(db *sql.DB) *CategoryRepository { return &CategoryRepository{db: db} }

func (r *CategoryRepository) GetAll() ([]models.Category, error) {
	rows, err := r.db.Query("SELECT id, name FROM categories ORDER BY name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cats []models.Category
	for rows.Next() {
		var c models.Category
		rows.Scan(&c.ID, &c.Name)
		cats = append(cats, c)
	}
	return cats, nil
}
