package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"time"

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

const (
	cookieName    = "session-cookie"
	authenticated = "authenticated"
)

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
	/*
		// Vérifier si l'utilisateur est connecté en vérifiant l'existence du cookie
		cookie, err := r.Cookie(cookieName)
		if err != nil || cookie.Value != authenticated {
			// Si l'utilisateur n'est pas connecté, le rediriger vers la page de connexion
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}


		// L'utilisateur est connecté, afficher la page d'accueil
		w.Write([]byte("Bienvenue sur la page d'accueil!"))
	*/
	cookie, err := r.Cookie(cookieName)
	if err != nil || cookie.Value != authenticated {
		// Si l'utilisateur n'est pas connecté, le rediriger vers la page de connexion
		fmt.Println("*************************** DECONNECTER ******************************")
	} else {
		fmt.Println("*************************** CONNECTER ******************************")
	}

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
		cookie, err := r.Cookie(cookieName)
		if err != nil || cookie.Value != authenticated {
			// Si l'utilisateur n'est pas connecté, le rediriger vers la page de connexion
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

	}

	defer db.Close()
	fmt.Println("************************************** SUCCES *************************************")
}

// LoginHandler gère le processus d'authentification
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Authentification réussie, définir un cookie pour marquer l'utilisateur comme connecté
	expiration := time.Now().Add(24 * time.Hour) // Expire dans 24 heures
	cookie := http.Cookie{Name: cookieName, Value: authenticated, Expires: expiration}
	http.SetCookie(w, &cookie)

	// Rediriger l'utilisateur vers la page d'accueil après la connexion
	http.Redirect(w, r, "/forum", http.StatusSeeOther)
}

// LogoutHandler gère le processus de déconnexion
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Déconnexion de l'utilisateur, supprimer le cookie
	cookie := http.Cookie{Name: cookieName, Value: "", Expires: time.Now().Add(-time.Hour)}
	http.SetCookie(w, &cookie)

	// Rediriger l'utilisateur vers la page d'accueil après la déconnexion
	http.Redirect(w, r, "/forum", http.StatusSeeOther)
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("static/"))) // gere les fichiers statics
	http.HandleFunc("/forum", home)
	http.HandleFunc("/register", formRegister)
	http.HandleFunc("/post_register", insertUsers)

	http.HandleFunc("/login", formLogin)
	http.HandleFunc("/post_login", LoginHandler)
	http.HandleFunc("/logout", LogoutHandler)

	fmt.Println("click sur le lien suivant : http://localhost:8080/forum")
	http.ListenAndServe(port, nil)
}
