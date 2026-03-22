package main

import (
    "database/sql"
    "log"

    _ "github.com/lib/pq"
)

var db *sql.DB

func InitDB() {
    var err error
    db, err = sql.Open("postgres", "user=postgres password=1234 dbname=evstations sslmode=disable")
    if err != nil {
        log.Fatal(err)
    }
}
