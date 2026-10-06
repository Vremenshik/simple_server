package repository

import (
	"database/sql"

	"go-expenses/internal/models"
)

// ExpenseRepository — обёртка над базой данных
type ExpenseRepository struct {
	db *sql.DB
}

// NewExpenseRepository — создать репозиторий
func NewExpenseRepository(db *sql.DB) *ExpenseRepository {
	return &ExpenseRepository{db: db}
}

// Create — добавить трату
func (r *ExpenseRepository) Create(e models.Expense) error {
	_, err := r.db.Exec(
		"INSERT INTO expenses (amount, description, date) VALUES (?, ?, ?)",
		e.Amount, e.Description, e.Date,
	)
	return err
}

// GetAll — получить все траты
func (r *ExpenseRepository) GetAll() ([]models.Expense, error) {
	rows, err := r.db.Query("SELECT id, amount, description, date FROM expenses ORDER BY id DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.Expense
	for rows.Next() {
		var e models.Expense
		if err := rows.Scan(&e.ID, &e.Amount, &e.Description, &e.Date); err != nil {
			return nil, err
		}
		result = append(result, e)
	}
	return result, nil
}

// GetByID — получить одну трату по ID
func (r *ExpenseRepository) GetByID(id int) (models.Expense, error) {
	var e models.Expense
	err := r.db.QueryRow(
		"SELECT id, amount, description, date FROM expenses WHERE id = ?", id,
	).Scan(&e.ID, &e.Amount, &e.Description, &e.Date)
	return e, err
}

// Update — обновить трату
func (r *ExpenseRepository) Update(e models.Expense) error {
	_, err := r.db.Exec(
		"UPDATE expenses SET amount = ?, description = ?, date = ? WHERE id = ?",
		e.Amount, e.Description, e.Date, e.ID,
	)
	return err
}

// Delete — удалить трату
func (r *ExpenseRepository) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM expenses WHERE id = ?", id)
	return err
}
