package main

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

func Connect(path string) (*sql.DB, error) {
	conn, err := sql.Open("sqlite", ("file:" + path))
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func main() {
	dbConnection, err := Connect("./database.db")
	if err != nil {
		panic(err)
	}
	dbConnection.Ping()
	fmt.Println("hello server")

}
