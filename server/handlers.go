package main

import (
	"encoding/json"
	"net/http"
	"time"
)

func newEntry(w http.ResponseWriter, r *http.Request) {
	var e Entry
	err := json.NewDecoder(r.Body).Decode(&e)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if e.Content == "" || e.UserEmail == "" || e.Time.IsZero() {
		http.Error(w, "missing field", http.StatusBadRequest)
		return
	}
	_, err = dbConnection.Exec("INSERT INTO Entries (userEmail, content, time) VALUES (?,?,?)", e.UserEmail, e.Content, e.Time.Format("2006-01-02 15:04:05"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func getUserEntries(w http.ResponseWriter, r *http.Request) {
	email := r.PathValue("userEmail")
	rows, err := dbConnection.Query("SELECT * FROM Entries WHERE userEmail=?", email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var entriesList []Entry
	var entry Entry
	for rows.Next() {
		err := rows.Scan(&entry.UserEmail, &entry.Content, &entry.Time)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		entriesList = append(entriesList, entry)
	}
	jsonResult, err := json.Marshal(entriesList)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(jsonResult)
}

func timeCheck(w http.ResponseWriter, r *http.Request) {
	cutOff := time.Now().Add(-24 * time.Hour)
	email := r.PathValue("userEmail")
	rows, err := dbConnection.Query("SELECT * FROM Entries WHERE userEmail=? AND time > ? ", email, cutOff)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	jsonResult, err := json.Marshal(Completed{IsCompleted: rows.Next()})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(jsonResult)
}
