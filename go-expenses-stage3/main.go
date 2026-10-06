package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	_ "modernc.org/sqlite"

	"go-expenses/internal/handlers"
)

func main() {
	db, err := sql.Open("sqlite", "expenses.db")
	if err != nil {
		fmt.Println("Error opening database:", err)
		os.Exit(1)
	}
	defer db.Close()

	migration, err := os.ReadFile("migrations/001_create_expenses.sql")
	if err != nil {
		fmt.Println("Error reading migration:", err)
		os.Exit(1)
	}
	if _, err := db.Exec(string(migration)); err != nil {
		fmt.Println("Error running migration:", err)
		os.Exit(1)
	}

	handlers.Init(db)

	http.HandleFunc("/", handlers.ListExpenses)
	http.HandleFunc("/add", handlers.AddExpense)
	http.HandleFunc("/edit", handlers.EditExpense)
	http.HandleFunc("/delete", handlers.DeleteExpense)

	fmt.Println("Server started on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
