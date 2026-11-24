package models

type CategoryStat struct {
	Name  string  `json:"name"`
	Total float64 `json:"total"`
}

type MonthStat struct {
	Month string  `json:"month"`
	Total float64 `json:"total"`
}
