package handlers

import (
	"database/sql"
	"html/template"
	"net/http"
	"strconv"

	"go-expenses-stage4/internal/models"
	"go-expenses-stage4/internal/repository"
)

var (
	repo      *repository.ExpenseRepository
	catRepo   *repository.CategoryRepository
	templates *template.Template
)

func Init(db *sql.DB) {
	repo = repository.NewExpenseRepository(db)
	catRepo = repository.NewCategoryRepository(db)
	templates = template.Must(template.ParseGlob("web/templates/*.html"))
}

func ListExpenses(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	var catID *int
	if catStr := q.Get("category"); catStr != "" {
		if id, err := strconv.Atoi(catStr); err == nil {
			catID = &id
		}
	}
	from := q.Get("from")
	to := q.Get("to")

	expenses, err := repo.GetAll(catID, from, to)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	categories, err := catRepo.GetAll()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	selectedCat := 0
	if catID != nil {
		selectedCat = *catID
	}

	data := struct {
		Expenses         []models.Expense
		Categories       []models.Category
		SelectedCategory int
		FromDate         string
		ToDate           string
	}{
		Expenses:         expenses,
		Categories:       categories,
		SelectedCategory: selectedCat,
		FromDate:         from,
		ToDate:           to,
	}

	templates.ExecuteTemplate(w, "list.html", data)
}

func AddExpense(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		categories, _ := catRepo.GetAll()
		templates.ExecuteTemplate(w, "add.html", categories)
		return
	}

	r.ParseForm()
	amount, _ := strconv.ParseFloat(r.FormValue("amount"), 64)

	var catID *int
	if cStr := r.FormValue("category_id"); cStr != "" {
		if id, err := strconv.Atoi(cStr); err == nil {
			catID = &id
		}
	}

	expense := models.Expense{
		Amount:      amount,
		Description: r.FormValue("description"),
		Date:        r.FormValue("date"),
		CategoryID:  catID,
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
		categories, _ := catRepo.GetAll()

		var currentCatID int
		if expense.CategoryID != nil {
			currentCatID = *expense.CategoryID
		}

		data := struct {
			Expense      models.Expense
			Categories   []models.Category
			CurrentCatID int
		}{
			Expense:      expense,
			Categories:   categories,
			CurrentCatID: currentCatID,
		}

		templates.ExecuteTemplate(w, "edit.html", data)
		return
	}

	r.ParseForm()
	amount, _ := strconv.ParseFloat(r.FormValue("amount"), 64)

	var catID *int
	if cStr := r.FormValue("category_id"); cStr != "" {
		if cid, err := strconv.Atoi(cStr); err == nil {
			catID = &cid
		}
	}

	expense := models.Expense{
		ID:          id,
		Amount:      amount,
		Description: r.FormValue("description"),
		Date:        r.FormValue("date"),
		CategoryID:  catID,
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
