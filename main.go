package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
var err error

const dbname = "/test"

func main() {
	db, err = sql.Open("mysql", dbname)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	http.HandleFunc("/", homepage)
	log.Fatal(http.ListenAndServe(":3000", nil))
}

func homepage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func login(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.ServeFile(w, r, "login.html")
		return
	}
}
