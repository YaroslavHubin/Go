package main

import (
	"database/sql"
	"encoding/json"
	"html/template"
	"net/http"
)

var tmpl = template.Must(template.ParseGlob("templates/*.html"))

// -------------------------
// Структура для головної сторінки
// -------------------------
type IndexData struct {
	UserName string
	Stations []Station
}

// -------------------------
// Головна сторінка
// -------------------------
func index(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, name, power, slots, location FROM stations")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var stations []Station
	for rows.Next() {
		var s Station
		rows.Scan(&s.ID, &s.Name, &s.Power, &s.Slots, &s.Location)
		stations = append(stations, s)
	}

	session, _ := store.Get(r, "session")
	userName, _ := session.Values["userName"].(string)

	data := IndexData{
		UserName: userName,
		Stations: stations,
	}

	tmpl.ExecuteTemplate(w, "index.html", data)
}

// -------------------------
// CRUD станцій
// -------------------------
func create(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		name := r.FormValue("name")
		power := r.FormValue("power")
		slots := r.FormValue("slots")
		location := r.FormValue("location")

		if name == "" || power == "" || slots == "" {
			tmpl.ExecuteTemplate(w, "index.html", map[string]string{"Error": "Некоректні дані"})
			return
		}

		_, err := db.Exec("INSERT INTO stations(name, power, slots, location) VALUES($1,$2,$3,$4)",
			name, power, slots, location)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func update(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		id := r.FormValue("id")
		name := r.FormValue("name")
		power := r.FormValue("power")
		slots := r.FormValue("slots")
		location := r.FormValue("location")

		if name == "" || power == "" || slots == "" {
			tmpl.ExecuteTemplate(w, "index.html", map[string]string{"Error": "Некоректні дані"})
			return
		}

		_, err := db.Exec("UPDATE stations SET name=$1, power=$2, slots=$3, location=$4 WHERE id=$5",
			name, power, slots, location, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func delete(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	_, err := db.Exec("DELETE FROM stations WHERE id=$1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// -------------------------
// Додавання зарядної сесії
// -------------------------
func addSession(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		stationID := r.FormValue("station_id")
		duration := r.FormValue("duration")
		energy := r.FormValue("energy")

		if stationID == "" || duration == "" || energy == "" {
			tmpl.ExecuteTemplate(w, "index.html", map[string]string{"Error": "Некоректні дані"})
			return
		}

		_, err := db.Exec("INSERT INTO sessions(station_id, duration, energy) VALUES($1, $2, $3)",
			stationID, duration, energy)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	http.Redirect(w, r, "/stats", http.StatusSeeOther)
}

// -------------------------
// Статистика
// -------------------------
func stats(w http.ResponseWriter, r *http.Request) {
	location := r.URL.Query().Get("location")
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
	for rows.Next() {
		var st Stat
		rows.Scan(&st.ID, &st.Name, &st.Location, &st.AvgDuration, &st.MaxEnergy, &st.TotalSessions)
		stats = append(stats, st)
	}
	tmpl.ExecuteTemplate(w, "stats.html", stats)
}

// -------------------------
// Авторизація
// -------------------------
func login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		email := r.FormValue("email")
		password := r.FormValue("password")

		var u User
		err := db.QueryRow("SELECT id, name, email, password, role FROM users WHERE email=$1", email).
			Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.Role)
		if err != nil || u.Password != password {
			tmpl.ExecuteTemplate(w, "login.html", map[string]string{"Error": "Невірний email або пароль"})
			return
		}

		session, _ := store.Get(r, "session")
		session.Values["userEmail"] = u.Email
		session.Values["userName"] = u.Name
		session.Values["userRole"] = u.Role
		session.Save(r, w)

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	tmpl.ExecuteTemplate(w, "login.html", nil)
}

func register(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		name := r.FormValue("name")
		email := r.FormValue("email")
		password := r.FormValue("password")
		_, err := db.Exec("INSERT INTO users(name, email, password) VALUES($1,$2,$3)", name, email, password)
		if err != nil {
			tmpl.ExecuteTemplate(w, "registration.html", map[string]string{"Error": "Помилка: можливо email вже існує"})
			return
		}
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	tmpl.ExecuteTemplate(w, "registration.html", nil)
}

// -------------------------
// Профіль
// -------------------------
func profile(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	user := User{
		Name:  session.Values["userName"].(string),
		Email: session.Values["userEmail"].(string),
		Role:  session.Values["userRole"].(string),
	}
	tmpl.ExecuteTemplate(w, "profile.html", user)
}

// -------------------------
// REST API
// -------------------------
func apiStations(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		rows, err := db.Query("SELECT id, name, power, slots, location FROM stations")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()
		var stations []Station
		for rows.Next() {
			var s Station
			rows.Scan(&s.ID, &s.Name, &s.Power, &s.Slots, &s.Location)
			stations = append(stations, s)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(stations)

	case "POST":
		var s Station
		if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
			http.Error(w, "Некоректні дані", http.StatusBadRequest)
			return
		}
		_, err := db.Exec("INSERT INTO stations(name, power, slots, location) VALUES($1,$2,$3,$4)",
			s.Name, s.Power, s.Slots, s.Location)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}

func apiSessions(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var sess Session
		if err := json.NewDecoder(r.Body).Decode(&sess); err != nil {
			http.Error(w, "Некоректні дані", http.StatusBadRequest)
			return
		}
		_, err := db.Exec("INSERT INTO sessions(station_id, duration, energy) VALUES($1,$2,$3)",
			sess.StationID, sess.Duration, sess.Energy)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}

func apiStats(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`
        SELECT s.id, s.name, s.location,
               AVG(sess.duration) AS avg_duration,
               MAX(sess.energy) AS max_energy,
               COUNT(sess.id) AS total_sessions
        FROM stations s
        LEFT JOIN sessions sess ON s.id = sess.station_id
        GROUP BY s.id`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var stats []Stat
	for rows.Next() {
		var st Stat
		rows.Scan(&st.ID, &st.Name, &st.Location, &st.AvgDuration, &st.MaxEnergy, &st.TotalSessions)
		stats = append(stats, st)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// -------------------------
// Middleware: доступ лише для адміністратора
// -------------------------
func adminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session")
		role, ok := session.Values["userRole"].(string)
		if !ok || role != "admin" {
			http.Error(w, "Доступ заборонено", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
