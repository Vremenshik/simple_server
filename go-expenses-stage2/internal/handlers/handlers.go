package handlers

import (
	"html/template"
	"net/http"
	"strconv"

	"go-expenses/internal/models"
)

var expenses []models.Expense
var nextID = 1

var templates = template.Must(template.ParseGlob("web/templates/*.html"))

func ListExpenses(w http.ResponseWriter, r *http.Request) {
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
		ID:          nextID,
		Amount:      amount,
		Description: r.FormValue("description"),
		Date:        r.FormValue("date"),
	}
	nextID++

	expenses = append(expenses, expense)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
