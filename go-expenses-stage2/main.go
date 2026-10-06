package main

import (
	"fmt"
	"net/http"

	"go-expenses/internal/handlers"
)

func main() {
	http.HandleFunc("/", handlers.ListExpenses)
	http.HandleFunc("/add", handlers.AddExpense)

	fmt.Println("Server started on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
