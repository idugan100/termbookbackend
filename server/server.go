package main

import (
	"database/sql"
	"net/http"
	"time"

	_ "modernc.org/sqlite"
)

type Entry struct {
	Content   string    `json:"content"`
	UserEmail string    `json:"userEmail"`
	Time      time.Time `json:"time"`
}

type Completed struct {
	IsCompleted bool `json:"isCompleted"`
}

var dbConnection *sql.DB

func main() {
	var err error
	dbConnection, err = Connect("./database.db")
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /newentry", newEntry)
	mux.HandleFunc("GET /entries/{userEmail}", getUserEntries)
	mux.HandleFunc("GET /timecheck/{userEmail}", timeCheck)

	s := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	panic(s.ListenAndServe())
}

func Connect(path string) (*sql.DB, error) {
	conn, err := sql.Open("sqlite", ("file:" + path))
	if err != nil {
		return nil, err
	}
	return conn, nil
}
