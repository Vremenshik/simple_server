package repository

import (
	"database/sql"
	"go-expenses-stage6/internal/models"
)

type ExpenseRepository struct{ db *sql.DB }

func NewExpenseRepository(db *sql.DB) *ExpenseRepository { return &ExpenseRepository{db: db} }

func (r *ExpenseRepository) Create(e *models.Expense) error {
	_, err := r.db.Exec(
		"INSERT INTO expenses (amount, description, date, user_id, category_id) VALUES (?, ?, ?, ?, ?)",
		e.Amount, e.Description, e.Date, e.UserID, e.CategoryID,
	)
	return err
}

func (r *ExpenseRepository) GetAll(userID int64, categoryID *int64, from, to string) ([]models.Expense, error) {
	query := `
		SELECT e.id, e.amount, e.description, e.date, e.user_id, e.category_id, COALESCE(c.name, '')
		FROM expenses e
		LEFT JOIN categories c ON e.category_id = c.id
		WHERE e.user_id = ?`

	args := []interface{}{userID}

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
		rows.Scan(&e.ID, &e.Amount, &e.Description, &e.Date, &e.UserID, &catID, &e.CategoryName)
		if catID.Valid {
			id := catID.Int64
			e.CategoryID = &id
		}
		result = append(result, e)
	}
	return result, nil
}

func (r *ExpenseRepository) GetByID(userID, id int64) (models.Expense, error) {
	var e models.Expense
	var catID sql.NullInt64
	err := r.db.QueryRow(
		"SELECT id, amount, description, date, user_id, category_id FROM expenses WHERE id = ? AND user_id = ?", id, userID,
	).Scan(&e.ID, &e.Amount, &e.Description, &e.Date, &e.UserID, &catID)
	if catID.Valid {
		cid := catID.Int64
		e.CategoryID = &cid
	}
	return e, err
}

func (r *ExpenseRepository) Update(e *models.Expense) error {
	_, err := r.db.Exec(
		"UPDATE expenses SET amount = ?, description = ?, date = ?, category_id = ? WHERE id = ? AND user_id = ?",
		e.Amount, e.Description, e.Date, e.CategoryID, e.ID, e.UserID,
	)
	return err
}

func (r *ExpenseRepository) Delete(userID, id int64) error {
	_, err := r.db.Exec("DELETE FROM expenses WHERE id = ? AND user_id = ?", id, userID)
	return err
}

// === НОВОЕ ДЛЯ 6 ЭТАПА: СТАТИСТИКА ===

func (r *ExpenseRepository) GetStatsByCategory(userID int64) ([]models.CategoryStat, error) {
	query := `
		SELECT COALESCE(c.name, 'Без категории') as name, SUM(e.amount) as total
		FROM expenses e
		LEFT JOIN categories c ON e.category_id = c.id
		WHERE e.user_id = ?
		GROUP BY e.category_id
		ORDER BY total DESC`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []models.CategoryStat
	for rows.Next() {
		var s models.CategoryStat
		if err := rows.Scan(&s.Name, &s.Total); err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}
	return stats, nil
}

func (r *ExpenseRepository) GetStatsByMonth(userID int64) ([]models.MonthStat, error) {
	// substr(date, 1, 7) берет 'YYYY-MM' из 'YYYY-MM-DD'
	query := `
		SELECT substr(date, 1, 7) as month, SUM(amount) as total
		FROM expenses
		WHERE user_id = ?
		GROUP BY month
		ORDER BY month`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []models.MonthStat
	for rows.Next() {
		var s models.MonthStat
		if err := rows.Scan(&s.Month, &s.Total); err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}
	return stats, nil
}
