package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	_ "modernc.org/sqlite"

	"go-expenses-stage4/internal/handlers"
)

func main() {
	db, err := sql.Open("sqlite", "expenses.db")
	if err != nil {
		fmt.Println("Error opening database:", err)
		os.Exit(1)
	}
	defer db.Close()

	// Применяем миграции
	migrations := []string{
		"migrations/001_create_expenses.sql",
		"migrations/002_add_categories.sql",
	}

	for _, file := range migrations {
		migration, err := os.ReadFile(file)
		if err != nil {
			fmt.Println("Error reading migration:", file, err)
			os.Exit(1)
		}
		if _, err := db.Exec(string(migration)); err != nil {
			// Игнорируем ошибку, если колонка уже существует (при повторном запуске)
			fmt.Println("Migration warning for", file, ":", err)
		}
	}

	handlers.Init(db)

	http.HandleFunc("/", handlers.ListExpenses)
	http.HandleFunc("/add", handlers.AddExpense)
	http.HandleFunc("/edit", handlers.EditExpense)
	http.HandleFunc("/delete", handlers.DeleteExpense)

	fmt.Println("Сервер запущен на http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
