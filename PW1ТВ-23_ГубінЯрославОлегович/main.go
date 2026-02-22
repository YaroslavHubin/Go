package main

import (
	"fmt"
	"html"
	"log"
	"net/http"
)

// logParams — універсальна функція для обробки та логування параметрів URL.
// Вона повертає карту (map) з усіма параметрами, які користувач передав у запиті.
// Також тут є обмеження довжини значення (до 50 символів), щоб уникнути засмічення логів.
func logParams(r *http.Request) map[string]string {
	params := r.URL.Query()           // Отримуємо всі параметри з URL
	result := make(map[string]string) // Створюємо карту для збереження параметрів

	for key, values := range params { // Перебираємо всі ключі та їхні значення
		for _, v := range values {
			// Якщо значення занадто довге — обрізаємо його
			if len(v) > 50 {
				v = v[:50] + "..."
			}
			// Логуємо параметр у консоль разом із адресою користувача
			log.Printf("Параметр: %s = %q (від %s)", key, v, r.RemoteAddr)
			result[key] = v
		}
	}
	return result
}

// stationHandler — головна сторінка зарядної станції.
// Вона показує тип станції (звичайна, fast тощо) і меню навігації.
func stationHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Запит: %s %s від %s", r.Method, r.URL.Path, r.RemoteAddr)
	params := logParams(r) // Логуємо всі параметри

	// Отримуємо параметр "station"
	stationType := params["station"]
	if stationType == "" {
		stationType = "звичайна"
	}
	// Екрануємо значення, щоб уникнути XSS
	stationTypeSafe := html.EscapeString(stationType)

	// Генеруємо HTML-сторінку
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, `
        <!DOCTYPE html>
        <html lang="uk">
        <head><meta charset="UTF-8"><title>Зарядна станція</title></head>
        <body>
            <h1>Зарядна станція для електромобілів</h1>
            <p>Тип станції: <b>%s</b></p>
            <ul>
                <li><a href="/pricing">Тарифи</a></li>
                <li><a href="/status">Статус портів</a></li>
                <li><a href="/support">Підтримка</a></li>
                <li><a href="/feedback">Зворотній зв'язок</a></li>
            </ul>
        </body>
        </html>
    `, stationTypeSafe)
}

// pricingHandler — сторінка з тарифами.
// Тут немає параметрів, але ми все одно викликаємо logParams для безпеки.
func pricingHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Запит: %s %s від %s", r.Method, r.URL.Path, r.RemoteAddr)
	logParams(r)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, `
        <!DOCTYPE html>
        <html lang="uk">
        <head><meta charset="UTF-8"><title>Тарифи</title></head>
        <body>
            <h1>Тарифи зарядної станції</h1>
            <p>Звичайна зарядка: 5 грн/кВт·год</p>
            <p>Швидка зарядка: 8 грн/кВт·год</p>
            <p><a href="/station">Назад</a></p>
        </body>
        </html>
    `)
}

// statusHandler — сторінка для перегляду статусу портів.
// Використовує параметр "port" (?port=1).
func statusHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Запит: %s %s від %s", r.Method, r.URL.Path, r.RemoteAddr)
	params := logParams(r)

	port := params["port"]
	if port == "" {
		port = "усі"
	}
	// Екрануємо значення параметра
	portSafe := html.EscapeString(port)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, `
        <!DOCTYPE html>
        <html lang="uk">
        <head><meta charset="UTF-8"><title>Статус портів</title></head>
        <body>
            <h1>Статус зарядних портів</h1>
            <p>Відображається інформація для: %s</p>
            <p>(Наприклад: <code>?port=1</code>)</p>
            <p><a href="/station">Назад</a></p>
        </body>
        </html>
    `, portSafe)
}

// supportHandler — сторінка служби підтримки.
// Просто показує контактні дані.
func supportHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Запит: %s %s від %s", r.Method, r.URL.Path, r.RemoteAddr)
	logParams(r)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, `
        <!DOCTYPE html>
        <html lang="uk">
        <head><meta charset="UTF-8"><title>Підтримка</title></head>
        <body>
            <h1>Служба підтримки</h1>
            <p>Телефон: +380-XX-XXX-XX-XX</p>
            <p>Email: support_charging_station@gmail.com</p>
            <p><a href="/station">Назад</a></p>
        </body>
        </html>
    `)
}

// feedbackHandler — сторінка для зворотного зв'язку.
// Використовує параметр "msg" (?msg=Привіт).
func feedbackHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Запит: %s %s від %s", r.Method, r.URL.Path, r.RemoteAddr)
	params := logParams(r)

	message := params["msg"]
	if message == "" {
		message = "немає повідомлення"
	}
	// Екрануємо повідомлення
	messageSafe := html.EscapeString(message)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, `
        <!DOCTYPE html>
        <html lang="uk">
        <head><meta charset="UTF-8"><title>Зворотній зв'язок</title></head>
        <body>
            <h1>Зворотній зв'язок</h1>
            <p>Ваше повідомлення: %s</p>
            <p>(Наприклад: <code>?msg=Привіт</code>)</p>
            <p><a href="/station">Назад</a></p>
        </body>
        </html>
    `, messageSafe)
}

// main — точка входу програми.
// Реєструє всі хендлери та запускає сервер на порту 8080.
func main() {
	http.HandleFunc("/station", stationHandler)
	http.HandleFunc("/pricing", pricingHandler)
	http.HandleFunc("/status", statusHandler)
	http.HandleFunc("/support", supportHandler)
	http.HandleFunc("/feedback", feedbackHandler)

	fmt.Println("Сервер запущено на http://localhost:8080/station")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
