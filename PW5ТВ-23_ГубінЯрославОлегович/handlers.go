package main

import (
	"database/sql"  // пакет для роботи з SQL-базами даних
	"html/template" // пакет для роботи з HTML-шаблонами
	"net/http"      // стандартна бібліотека для роботи з HTTP-запитами
)

// створюємо глобальну змінну шаблонів, які будуть братися з папки templates
var tmpl = template.Must(template.ParseGlob("templates/*.html"))

// -------------------------
// Хендлер для головної сторінки (CRUD для станцій)
// -------------------------
func index(w http.ResponseWriter, r *http.Request) {
	// виконуємо SQL-запит для отримання всіх станцій
	rows, err := db.Query("SELECT id, name, power, slots, location FROM stations")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close() // закриваємо результат після використання

	var stations []Station
	// перебираємо всі рядки результату
	for rows.Next() {
		var s Station
		rows.Scan(&s.ID, &s.Name, &s.Power, &s.Slots, &s.Location)
		stations = append(stations, s)
	}

	// передаємо список станцій у шаблон index.html
	tmpl.ExecuteTemplate(w, "index.html", stations)
}

// -------------------------
// Хендлер для створення нової станції
// -------------------------
func create(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		// отримуємо дані з форми
		name := r.FormValue("name")
		power := r.FormValue("power")
		slots := r.FormValue("slots")
		location := r.FormValue("location")

		// вставляємо новий запис у таблицю stations
		_, err := db.Exec("INSERT INTO stations(name, power, slots, location) VALUES($1,$2,$3,$4)",
			name, power, slots, location)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	// після вставки перенаправляємо на головну сторінку
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// -------------------------
// Хендлер для оновлення даних станції
// -------------------------
func update(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		// отримуємо дані з форми
		id := r.FormValue("id")
		name := r.FormValue("name")
		power := r.FormValue("power")
		slots := r.FormValue("slots")
		location := r.FormValue("location")

		// оновлюємо запис у таблиці stations
		_, err := db.Exec("UPDATE stations SET name=$1, power=$2, slots=$3, location=$4 WHERE id=$5",
			name, power, slots, location, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	// після оновлення перенаправляємо на головну сторінку
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// -------------------------
// Хендлер для видалення станції
// -------------------------
func delete(w http.ResponseWriter, r *http.Request) {
	// отримуємо id станції з параметра URL
	id := r.URL.Query().Get("id")
	// видаляємо запис з таблиці stations
	_, err := db.Exec("DELETE FROM stations WHERE id=$1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// після видалення перенаправляємо на головну сторінку
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// -------------------------
// Хендлер для додавання нової зарядної сесії
// -------------------------
func addSession(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		// отримуємо дані з форми
		stationID := r.FormValue("station_id")
		duration := r.FormValue("duration")
		energy := r.FormValue("energy")

		// вставляємо новий запис у таблицю sessions
		_, err := db.Exec("INSERT INTO sessions(station_id, duration, energy) VALUES($1, $2, $3)",
			stationID, duration, energy)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	// після вставки перенаправляємо на сторінку статистики
	http.Redirect(w, r, "/stats", http.StatusSeeOther)
}

// -------------------------
// Хендлер для перегляду статистики станцій
// -------------------------
func stats(w http.ResponseWriter, r *http.Request) {
	// отримуємо параметр фільтрації по локації
	location := r.URL.Query().Get("location")

	// SQL-запит з JOIN та агрегатними функціями
	query := `
        SELECT s.id, s.name, s.location,
               AVG(sess.duration) AS avg_duration,
               MAX(sess.energy) AS max_energy,
               COUNT(sess.id) AS total_sessions
        FROM stations s
        LEFT JOIN sessions sess ON s.id = sess.station_id
    `
	var rows *sql.Rows
	var err error

	// якщо задана локація — додаємо WHERE
	if location != "" {
		query += " WHERE s.location = $1 GROUP BY s.id"
		rows, err = db.Query(query, location)
	} else {
		query += " GROUP BY s.id"
		rows, err = db.Query(query)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var stats []Stat
	// зчитуємо результати у структуру Stat
	for rows.Next() {
		var st Stat
		rows.Scan(&st.ID, &st.Name, &st.Location, &st.AvgDuration, &st.MaxEnergy, &st.TotalSessions)
		stats = append(stats, st)
	}

	// передаємо дані у шаблон stats.html
	tmpl.ExecuteTemplate(w, "stats.html", stats)
}
