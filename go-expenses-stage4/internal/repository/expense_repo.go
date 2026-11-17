package repository

import (
	"database/sql"

	"go-expenses-stage4/internal/models"
)

type ExpenseRepository struct {
	db *sql.DB
}

func NewExpenseRepository(db *sql.DB) *ExpenseRepository {
	return &ExpenseRepository{db: db}
}

func (r *ExpenseRepository) Create(e models.Expense) error {
	_, err := r.db.Exec(
		"INSERT INTO expenses (amount, description, date, category_id) VALUES (?, ?, ?, ?)",
		e.Amount, e.Description, e.Date, e.CategoryID,
	)
	return err
}

// GetAll с фильтрацией и JOIN (защита от N+1)
func (r *ExpenseRepository) GetAll(categoryID *int, from, to string) ([]models.Expense, error) {
	query := `
		SELECT e.id, e.amount, e.description, e.date, e.category_id, COALESCE(c.name, '')
		FROM expenses e
		LEFT JOIN categories c ON e.category_id = c.id
		WHERE 1=1`

	var args []interface{}

	if categoryID != nil {
		query += " AND e.category_id = ?"
		args = append(args, *categoryID)
	}
	if from != "" {
		query += " AND e.date >= ?"
		args = append(args, from)
	}
	if to != "" {
		query += " AND e.date <= ?"
		args = append(args, to)
	}

	query += " ORDER BY e.date DESC, e.id DESC"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.Expense
	for rows.Next() {
		var e models.Expense
		var catID sql.NullInt64

		if err := rows.Scan(&e.ID, &e.Amount, &e.Description, &e.Date, &catID, &e.CategoryName); err != nil {
			return nil, err
		}

		if catID.Valid {
			id := int(catID.Int64)
			e.CategoryID = &id
		}
		result = append(result, e)
	}
	return result, nil
}

func (r *ExpenseRepository) GetByID(id int) (models.Expense, error) {
	var e models.Expense
	var catID sql.NullInt64

	err := r.db.QueryRow(
		"SELECT id, amount, description, date, category_id FROM expenses WHERE id = ?", id,
	).Scan(&e.ID, &e.Amount, &e.Description, &e.Date, &catID)

	if catID.Valid {
		cid := int(catID.Int64)
		e.CategoryID = &cid
	}
	return e, err
}

func (r *ExpenseRepository) Update(e models.Expense) error {
	_, err := r.db.Exec(
		"UPDATE expenses SET amount = ?, description = ?, date = ?, category_id = ? WHERE id = ?",
		e.Amount, e.Description, e.Date, e.CategoryID, e.ID,
	)
	return err
}

func (r *ExpenseRepository) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM expenses WHERE id = ?", id)
	return err
}
