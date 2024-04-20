package main

import (
	"database/sql"
	"net/http"

	_ "modernc.org/sqlite"
)

type Entry struct {
	Content   string `json:"content"`
	UserEmail string `json:"userEmail"`
	Time      string `json:"time"`
}

func Connect(path string) (*sql.DB, error) {
	conn, err := sql.Open("sqlite", ("file:" + path))
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func newEntry(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("hi create"))
}

func getUserEntries(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hi user entries"))
}

func timeCheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("time check"))
}

func main() {
	dbConnection, err := Connect("./database.db")
	if err != nil {
		panic(err)
	}
	dbConnection.Ping()

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
