package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	_ "modernc.org/sqlite"

	"go-expenses-stage5/internal/auth"
	"go-expenses-stage5/internal/handlers"
	"go-expenses-stage5/internal/repository"
)

func main() {
	db, err := sql.Open("sqlite", "expenses.db")
	if err != nil {
		fmt.Println("DB Error:", err)
		os.Exit(1)
	}
	defer db.Close()

	migrations := []string{"migrations/001_init.sql", "migrations/002_categories.sql"}
	for _, file := range migrations {
		content, err := os.ReadFile(file)
		if err != nil {
			fmt.Println("Migration read error:", err)
			os.Exit(1)
		}
		if _, err := db.Exec(string(content)); err != nil {
			fmt.Println("Migration warning:", err)
		}
	}

	expRepo := repository.NewExpenseRepository(db)
	catRepo := repository.NewCategoryRepository(db)
	usrRepo := repository.NewUserRepository(db)
	sm := auth.NewSessionManager()

	h := handlers.New(expRepo, catRepo, usrRepo, sm)

	// Маршруты авторизации (публичные)
	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			h.ShowRegister(w, r)
		} else {
			h.Register(w, r)
		}
	})
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			h.ShowLogin(w, r)
		} else {
			h.Login(w, r)
		}
	})
	http.HandleFunc("/logout", h.Logout)

	// Защищенные маршруты (только для авторизованных)
	http.HandleFunc("/", h.AuthRequired(h.ListExpenses))
	http.HandleFunc("/add", h.AuthRequired(h.AddExpense))
	http.HandleFunc("/edit", h.AuthRequired(h.EditExpense))
	http.HandleFunc("/delete", h.AuthRequired(h.DeleteExpense))

	fmt.Println("Сервер запущен на http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
