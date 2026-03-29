package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte("secret-key"))

func main() {
	InitDB()
	defer db.Close()

	// UI
	http.HandleFunc("/", index)
	http.HandleFunc("/create", create)
	http.HandleFunc("/update", update)
	http.HandleFunc("/delete", delete)
	http.HandleFunc("/stats", stats)
	http.HandleFunc("/addsession", addSession)

	// Auth
	http.HandleFunc("/login", login)
	http.HandleFunc("/register", register)
	http.HandleFunc("/profile", profile)

	// API
	http.HandleFunc("/api/stations", apiStations)
	http.HandleFunc("/api/sessions", apiSessions)
	http.HandleFunc("/api/stats", apiStats)

	// Swagger UI (лише для адміністратора)
	http.Handle("/swagger.yaml", http.FileServer(http.Dir(".")))
	http.Handle("/swagger/", adminOnly(http.StripPrefix("/swagger/", http.FileServer(http.Dir("templates/swagger")))))

	// Static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	fmt.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
