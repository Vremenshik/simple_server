package models

type Expense struct {
	ID           int64
	Amount       float64
	Description  string
	Date         string
	UserID       int64
	CategoryID   *int64
	CategoryName string
}
