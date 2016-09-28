package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
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

	http.HandleFunc("/login", login)
	http.HandleFunc("/signup", signupPage)
	http.HandleFunc("/", homepage)
	log.Fatal(http.ListenAndServe(":3000", nil))
}

func homepage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func login(res http.ResponseWriter, req *http.Request) {
	log.Print("logging in")
	if req.Method != "POST" {
		log.Print("not a POST - redirecting")
		http.ServeFile(res, req, "login.html")
		return
	}

	username, password := req.FormValue("username"), req.FormValue("password")
	var dbUsername, dbPassword string

	select1 := "SELECT username, password FROM users WHERE username=?"
	err := db.QueryRow(select1, username).Scan(&dbUsername, &dbPassword)
	if err != nil {
		log.Print("db error - redirecting")
		http.Redirect(res, req, "/login", 301)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(password))
	if err != nil {
		log.Print("bcrypt error - redirecting")
		http.Redirect(res, req, "/login", 301)
		return
	}

	log.Print("success")
	res.Write([]byte("Hello " + dbUsername))
}

func signupPage(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.ServeFile(res, req, "signup.html")
	}

	username, password := req.FormValue("username"), req.FormValue("password")
	var user string

	dbErr := db.QueryRow("SELECT username FROM users WHERE username=?", username).Scan(&user)

	switch {
	case dbErr == sql.ErrNoRows:
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(res, "Server error, unable to create your account", 500)
			return
		}

		_, err = db.Exec("INSERT INTO users(username, password) VALUES(?, ?)", username, hashedPassword)
		if err != nil {
			http.Error(res, "Server error, unable to create your account.", 500)
			return
		}
	case dbErr != nil:
		http.Error(res, "Server error, unable to create your account", 500)
		return
	default:
		http.Redirect(res, req, "/", 301)
	}

}
