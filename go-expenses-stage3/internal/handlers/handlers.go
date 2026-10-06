package handlers

import (
	"database/sql"
	"html/template"
	"net/http"
	"strconv"

	"go-expenses/internal/models"
	"go-expenses/internal/repository"
)

var (
	repo      *repository.ExpenseRepository
	templates *template.Template
)

func Init(db *sql.DB) {
	repo = repository.NewExpenseRepository(db)
	templates = template.Must(template.ParseGlob("web/templates/*.html"))
}

func ListExpenses(w http.ResponseWriter, r *http.Request) {
	expenses, err := repo.GetAll()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	templates.ExecuteTemplate(w, "list.html", expenses)
}

func AddExpense(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		templates.ExecuteTemplate(w, "add.html", nil)
		return
	}

	r.ParseForm()
	amount, _ := strconv.ParseFloat(r.FormValue("amount"), 64)

	expense := models.Expense{
		Amount:      amount,
		Description: r.FormValue("description"),
		Date:        r.FormValue("date"),
	}

	if err := repo.Create(expense); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func EditExpense(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))

	if r.Method == "GET" {
		expense, err := repo.GetByID(id)
		if err != nil {
			http.Error(w, "Трата не найдена", 404)
			return
		}
		templates.ExecuteTemplate(w, "edit.html", expense)
		return
	}

	r.ParseForm()
	amount, _ := strconv.ParseFloat(r.FormValue("amount"), 64)

	expense := models.Expense{
		ID:          id,
		Amount:      amount,
		Description: r.FormValue("description"),
		Date:        r.FormValue("date"),
	}

	if err := repo.Update(expense); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func DeleteExpense(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	if err := repo.Delete(id); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
