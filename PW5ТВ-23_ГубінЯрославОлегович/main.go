package main

import (
    "fmt"
    "log"
    "net/http"
)

func main() {
    InitDB()
    defer db.Close()

    http.HandleFunc("/", index)
    http.HandleFunc("/create", create)
    http.HandleFunc("/update", update)
    http.HandleFunc("/delete", delete)
    http.HandleFunc("/stats", stats)
    http.HandleFunc("/addsession", addSession)


    fmt.Println("Server started at :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
