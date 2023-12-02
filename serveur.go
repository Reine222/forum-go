package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	_ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library
)

// ----------------------------------> mes tables structures -------------------------------------------------------------------------/

// structure qui gere les artistes
type T_UTILISATEUR_USER struct {
	User_id       int
	User_name     string
	User_email    string
	User_password string
}

// structure qui gere les artistes
type T_PUBLICATION_POST struct {
	Post_id          int
	Post_date        string
	Post_description string
	User_id          int
}

// structure qui gere les artistes
type T_COMMENTAIRE_COMMENT struct {
	Comment_id          int
	Comment_date        string
	Comment_description string
	Post_id             int
	User_id             int
}

// structure qui gere les artistes
type T_REACTION_REACT struct {
	React_id   int
	React_date string
	React_type string
	User_id    int
	Post_id    int
	Comment_id int
}

// ----------------------------------> mes tables structures -------------------------------------------------------------------------/

const port = ":8080"

// recuperation des templates (pages html)
var (
	template0 = template.Must(template.ParseFiles("./template/index.html"))
	template1 = template.Must(template.ParseFiles("./template/login.html"))
	template2 = template.Must(template.ParseFiles("./template/register.html"))
)

// Connect to database
func dbConn() (db *sql.DB) {
	db, err := sql.Open("sqlite3", "forum.db")
	if err != nil {
		panic(err.Error())
	}
	return db
}

func home(w http.ResponseWriter, r *http.Request) {
	err1 := template0.Execute(w, nil) // execution de la page index.html
	if err1 != nil {
		http.Error(w, err1.Error(), http.StatusInternalServerError)
	}
}

func formLogin(w http.ResponseWriter, r *http.Request) {
	err1 := template1.Execute(w, nil) // execution de la page index.html
	if err1 != nil {
		http.Error(w, err1.Error(), http.StatusInternalServerError)
	}
}

func formRegister(w http.ResponseWriter, r *http.Request) {
	err1 := template2.Execute(w, nil) // execution de la page index.html
	if err1 != nil {
		http.Error(w, err1.Error(), http.StatusInternalServerError)
	}
}

func insertUsers(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		user_name := r.FormValue("nom")
		user_email := r.FormValue("email")
		pass := r.FormValue("password")
		passs := []byte(pass)
		// Hashing the password
		user_password, err := bcrypt.GenerateFromPassword(passs, bcrypt.DefaultCost)
		if err != nil {
			panic(err)
		}
		fmt.Println(user_name, user_email, user_password)
		register := "INSERT INTO T_UTILISATEUR_USER (user_name, user_email, user_password) VALUES (?,?,?)"
		_, errr := db.Exec(register, user_name, user_email, user_password)
		if errr != nil {
			fmt.Println("Erreur lors de l'insertion dans la table users:", err)
			return
		}

		fmt.Println("Insertion réussie !")

	}

	defer db.Close()
	fmt.Println("************************************** SUCCES *************************************")
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("static/"))) // gere les fichiers statics
	http.HandleFunc("/forum", home)
	http.HandleFunc("/register", formRegister)
	http.HandleFunc("/post_register", insertUsers)

	http.HandleFunc("/login", formLogin)

	fmt.Println("click sur le lien suivant : http://localhost:8080/forum")
	http.ListenAndServe(port, nil)
}
