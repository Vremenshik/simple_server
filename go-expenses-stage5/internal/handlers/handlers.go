package handlers

import (
	"html/template"
	"net/http"
	"strconv"

	"go-expenses-stage5/internal/auth"
	"go-expenses-stage5/internal/models"
	"go-expenses-stage5/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	expenseRepo  *repository.ExpenseRepository
	categoryRepo *repository.CategoryRepository
	userRepo     *repository.UserRepository
	sessions     *auth.SessionManager
	templates    *template.Template
}

func New(expRepo *repository.ExpenseRepository, catRepo *repository.CategoryRepository, usrRepo *repository.UserRepository, sm *auth.SessionManager) *Handler {
	tmpl := template.Must(template.ParseGlob("web/templates/*.html"))
	return &Handler{
		expenseRepo: expRepo, categoryRepo: catRepo, userRepo: usrRepo, sessions: sm, templates: tmpl,
	}
}

// --- АВТОРИЗАЦИЯ ---

func (h *Handler) ShowRegister(w http.ResponseWriter, r *http.Request) {
	h.templates.ExecuteTemplate(w, "register.html", nil)
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	email := r.FormValue("email")
	password := r.FormValue("password")

	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	_, err := h.userRepo.Create(email, string(hash))
	if err != nil {
		h.templates.ExecuteTemplate(w, "register.html", map[string]string{"Error": "Email уже занят"})
		return
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (h *Handler) ShowLogin(w http.ResponseWriter, r *http.Request) {
	h.templates.ExecuteTemplate(w, "login.html", nil)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	email := r.FormValue("email")
	password := r.FormValue("password")

	user, err := h.userRepo.GetByEmail(email)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		h.templates.ExecuteTemplate(w, "login.html", map[string]string{"Error": "Неверный email или пароль"})
		return
	}

	// Создаем сессию и ставим cookie
	sessionID := h.sessions.Create(user.ID)
	h.sessions.SetCookie(w, sessionID)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err == nil {
		h.sessions.Delete(cookie.Value)
	}
	h.sessions.ClearCookie(w)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// --- ТРАТЫ ---

func (h *Handler) ListExpenses(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	q := r.URL.Query()

	var catID *int64
	if catStr := q.Get("category"); catStr != "" {
		if id, err := strconv.ParseInt(catStr, 10, 64); err == nil {
			catID = &id
		}
	}

	expenses, err := h.expenseRepo.GetAll(userID, catID, q.Get("from"), q.Get("to"))
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	categories, _ := h.categoryRepo.GetAll()

	selectedCat := int64(0)
	if catID != nil {
		selectedCat = *catID
	}

	data := map[string]interface{}{
		"Expenses": expenses, "Categories": categories, "SelectedCategory": selectedCat,
		"FromDate": q.Get("from"), "ToDate": q.Get("to"),
	}
	h.templates.ExecuteTemplate(w, "list.html", data)
}

func (h *Handler) AddExpense(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		categories, _ := h.categoryRepo.GetAll()
		h.templates.ExecuteTemplate(w, "add.html", categories)
		return
	}

	userID := getUserID(r)
	r.ParseForm()
	amount, _ := strconv.ParseFloat(r.FormValue("amount"), 64)

	var catID *int64
	if cStr := r.FormValue("category_id"); cStr != "" {
		if id, err := strconv.ParseInt(cStr, 10, 64); err == nil {
			catID = &id
		}
	}

	expense := &models.Expense{
		Amount: amount, Description: r.FormValue("description"),
		Date: r.FormValue("date"), UserID: userID, CategoryID: catID,
	}

	if err := h.expenseRepo.Create(expense); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) EditExpense(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	id, _ := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)

	if r.Method == "GET" {
		expense, err := h.expenseRepo.GetByID(userID, id)
		if err != nil {
			http.Error(w, "Трата не найдена", 404)
			return
		}

		categories, _ := h.categoryRepo.GetAll()
		var currentCatID int64 = 0
		if expense.CategoryID != nil {
			currentCatID = *expense.CategoryID
		}

		data := map[string]interface{}{
			"Expense": expense, "Categories": categories, "CurrentCatID": currentCatID,
		}
		h.templates.ExecuteTemplate(w, "edit.html", data)
		return
	}

	r.ParseForm()
	amount, _ := strconv.ParseFloat(r.FormValue("amount"), 64)
	var catID *int64
	if cStr := r.FormValue("category_id"); cStr != "" {
		if cid, err := strconv.ParseInt(cStr, 10, 64); err == nil {
			catID = &cid
		}
	}

	expense := &models.Expense{
		ID: id, Amount: amount, Description: r.FormValue("description"),
		Date: r.FormValue("date"), UserID: userID, CategoryID: catID,
	}

	if err := h.expenseRepo.Update(expense); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) DeleteExpense(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	id, _ := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)

	if err := h.expenseRepo.Delete(userID, id); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
