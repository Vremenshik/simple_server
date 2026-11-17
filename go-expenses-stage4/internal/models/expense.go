package models

type Expense struct {
	ID           int
	Amount       float64
	Description  string
	Date         string // Строка, так как в БД хранится TEXT
	CategoryID   *int
	CategoryName string
}
