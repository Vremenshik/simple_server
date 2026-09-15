package main

import (
    "fmt"
    "net/http"
)

func main() {
    // Регистрируем маршруты
    http.HandleFunc("/", home)
    http.HandleFunc("/about", about)
    http.HandleFunc("/ping", ping)

    // Запускаем сервер
    fmt.Println("Сервер запущен на http://localhost:8080")
    http.ListenAndServe(":8080", nil)
}

// Главная страница
func home(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Добро пожаловать в приложение «Учёт личных трат»!")
}

// Страница "О проекте"
func about(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Это веб-приложение для учёта личных расходов.")
    fmt.Fprintln(w, "Стек: Go + net/http + html/template + SQLite")
}

// Проверка работы сервера (только GET)
func ping(w http.ResponseWriter, r *http.Request) {
    if r.Method != "GET" {
        http.Error(w, "Метод не разрешён", http.StatusMethodNotAllowed)
        return
    }
    fmt.Fprintln(w, "pong")
}