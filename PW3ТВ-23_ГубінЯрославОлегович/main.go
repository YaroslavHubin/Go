package main

import (
    "fmt"
    "html/template"
    "log"
    "net/http"
    "strconv"
    "strings"
    "time"
)

// Структура, яка описує зарядну станцію
type Station struct {
    UserName   string    // ім’я користувача
    Name       string    // назва станції
    Address    string    // адреса станції
    Ports      int       // кількість портів
    PowerKW    float64   // потужність у кВт
    Installed  time.Time // дата встановлення
}

// Зріз для збереження всіх доданих станцій
var stations []Station

func main() {
    // Реєструємо обробники маршрутів
    http.HandleFunc("/", formHandler)       // головна сторінка з формою
    http.HandleFunc("/submit", submitHandler) // обробка даних після відправки форми

    port := ":8080" // порт, на якому працює сервер
    fmt.Printf("Server started at http://localhost%s\n", port) // повідомлення про запуск сервера
    log.Fatal(http.ListenAndServe(port, nil)) // запуск сервера, у разі помилки — завершення програми
}

// Обробник головної сторінки з формою
func formHandler(w http.ResponseWriter, r *http.Request) {
    tmpl := template.Must(template.ParseFiles("templates/form.html")) // завантажуємо шаблон форми
    tmpl.Execute(w, nil) // рендеримо шаблон без даних
}

// Обробник форми після її відправки
func submitHandler(w http.ResponseWriter, r *http.Request) {
    // Перевіряємо, чи запит був методом POST
    if r.Method != http.MethodPost {
        http.Redirect(w, r, "/", http.StatusSeeOther) // якщо ні — перенаправляємо назад на форму
        return
    }

    // Отримуємо дані з форми
    userName := r.FormValue("username")
    name := r.FormValue("name")
    address := r.FormValue("address")
    portsStr := r.FormValue("ports")
    powerStr := r.FormValue("power")
    dateStr := r.FormValue("date")

    // Логування введених даних + IP користувача
    clientIP := r.RemoteAddr
    log.Printf("Request received from IP=%s: UserName=%s, Name=%s, Address=%s, Ports=%s, Power=%s, Date=%s",
        clientIP, userName, name, address, portsStr, powerStr, dateStr)

    // Масив для збереження помилок
    errors := []string{}
    if userName == "" {
        errors = append(errors, "Ім’я користувача обов’язкове")
    }
    if name == "" {
        errors = append(errors, "Назва станції обов’язкова")
    }
    if address == "" {
        errors = append(errors, "Адреса обов’язкова")
    }

    // Перетворення кількості портів у число
    ports, err := strconv.Atoi(portsStr)
    if err != nil || ports <= 0 {
        errors = append(errors, "Кількість портів має бути > 0")
    }

    // Перетворення потужності у число з плаваючою точкою
    power, err := strconv.ParseFloat(powerStr, 64)
    if err != nil || power <= 0 {
        errors = append(errors, "Потужність має бути > 0")
    }

    // Перетворення дати у формат time.Time
    installed, err := time.Parse("2006-01-02", dateStr)
    if err != nil {
        errors = append(errors, "Дата має бути у форматі YYYY-MM-DD")
    }

    // Якщо є помилки — повертаємо користувача назад на форму з повідомленням
    if len(errors) > 0 {
        tmpl := template.Must(template.ParseFiles("templates/form.html"))
        tmpl.Execute(w, strings.Join(errors, "; ")) // показуємо всі помилки одним рядком
        return
    }

    // Створюємо нову станцію з даних користувача
    station := Station{userName, name, address, ports, power, installed}
    stations = append(stations, station) // додаємо її у зріз

    // Виводимо сторінку успіху з даними станції
    tmpl := template.Must(template.ParseFiles("templates/success.html"))
    tmpl.Execute(w, station)
}