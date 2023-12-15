package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"

	_ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library
)

// ----------------------------------> mes tables structures -------------------------------------------------------------------------/

type T_CATEGORY_CAT struct {
	Cat_id   int
	Cat_name string
}

type T_COMMENTAIRE_COMMENT struct {
	Comment_id          int
	Comment_date        string
	Comment_description string
	User_id             T_UTILISATEUR_USER
	Post_id             T_PUBLICATION_POST
}

type T_PUBLICATION_POST struct {
	Post_id          int
	Post_title       string
	Post_date        string
	Post_description string
	User_id          T_UTILISATEUR_USER
}

type T_CATEGORIEPUB_CATPUB struct {
	Catpub_id int
	Cat_id    T_CATEGORY_CAT
	Post_id   T_PUBLICATION_POST
}

type T_UTILISATEUR_USER struct {
	User_id       int
	User_name     string
	User_email    string
	User_password string
}

type T_LIKECOMMENTAIRE_LIKECOM struct {
	Likecom_id        int
	Likecom_date      string
	Likecom_champlike int
	Comment_id        T_COMMENTAIRE_COMMENT
	User_id           T_UTILISATEUR_USER
}

type T_LIKEPUBLICATION_PUB struct {
	Likepub_id        int
	Likepub_date      string
	Likepub_champlike int
	Post_id           T_PUBLICATION_POST
	User_id           T_UTILISATEUR_USER
}

// ----------------------------------> mes tables structures -------------------------------------------------------------------------/

const port = ":8080"

const (
	cookieName    = "session-cookie"
	authenticated = "authenticated"
)

// recuperation des templates (pages html)
var (
	template00 = template.Must(template.ParseFiles("index.html"))
	template0  = template.Must(template.ParseFiles("index1.html"))

	template1 = template.Must(template.ParseFiles("page-connexion.html"))
	template2 = template.Must(template.ParseFiles("page-inscription.html"))

	template33 = template.Must(template.ParseFiles("page-create-pub1.html"))
	template44 = template.Must(template.ParseFiles("page-categories1.html"))
	template55 = template.Must(template.ParseFiles("page-single-topic1.html"))

	template3 = template.Must(template.ParseFiles("page-create-pub.html"))
	template4 = template.Must(template.ParseFiles("page-categories.html"))
	template5 = template.Must(template.ParseFiles("page-single-topic.html"))
)

// Connect to database
func dbConn() (db *sql.DB) {
	db, err := sql.Open("sqlite3", "forum.db")
	if err != nil {
		panic(err.Error())
	}
	return db
}

func home0(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	// selDB, err := db.Query("SELECT * FROM T_PUBLICATION_POST ORDER BY post_id DESC")
	selDB, err := db.Query("SELECT * FROM T_PUBLICATION_POST")
	if err != nil {
		fmt.Println(err.Error())
	}
	post := T_PUBLICATION_POST{}
	all_pub := []T_PUBLICATION_POST{}
	for selDB.Next() {
		var post_id, user_id int
		var post_title, post_date, post_description string
		err = selDB.Scan(&post_id, &post_title, &post_date, &post_description, &user_id)
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println("+++++++++++++++++++++++++++++++", post.User_id.User_id)

		post.Post_id = post_id
		post.Post_title = post_title
		post.Post_date = post_date
		post.Post_description = post_description
		post.User_id.User_id = user_id
		all_pub = append(all_pub, post)

	}

	m := map[string]interface{}{
		"data": all_pub,
	}

	err1 := template00.Execute(w, m) // execution de la page index.html
	if err1 != nil {
		http.Error(w, err1.Error(), http.StatusInternalServerError)
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(cookieName)
	if err != nil || cookie.Value != authenticated {
		// Si l'utilisateur n'est pas connecté, le rediriger vers la page de connexion
		fmt.Println("*************************** DECONNECTER ******************************")
	} else {
		fmt.Println("*************************** CONNECTER ******************************")
	}

	db := dbConn()
	// selDB, err := db.Query("SELECT * FROM T_PUBLICATION_POST ORDER BY post_id DESC")
	selDB, err := db.Query("SELECT * FROM T_PUBLICATION_POST")
	if err != nil {
		fmt.Println(err.Error())
	}
	post := T_PUBLICATION_POST{}
	all_pub := []T_PUBLICATION_POST{}
	for selDB.Next() {
		var post_id, user_id int
		var post_title, post_date, post_description string
		err = selDB.Scan(&post_id, &post_title, &post_date, &post_description, &user_id)
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println("+++++++++++++++++++++++++++++++", post.User_id.User_id)

		post.Post_id = post_id
		post.Post_title = post_title
		post.Post_date = post_date
		post.Post_description = post_description
		post.User_id.User_id = user_id
		all_pub = append(all_pub, post)

	}

	fmt.Println("----------------------------------------", all_pub[0])

	m := map[string]interface{}{
		"data":   all_pub,
		"cookie": cookie,
	}

	err1 := template0.Execute(w, m) // execution de la page index.html
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

func formPublication(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(cookieName)
	if err != nil || cookie.Value != authenticated {
		// Si l'utilisateur n'est pas connecté, le rediriger vers la page de connexion
		fmt.Println("*************************** DECONNECTER ******************************")
	} else {
		fmt.Println("*************************** CONNECTER ******************************")
	}

	m := map[string]interface{}{
		"cookie": cookie,
	}
	err1 := template3.Execute(w, m) // execution de la page index.html
	if err1 != nil {
		http.Error(w, err1.Error(), http.StatusInternalServerError)
	}
}

func formPublication1(w http.ResponseWriter, r *http.Request) {
	err1 := template33.Execute(w, nil) // execution de la page index.html
	if err1 != nil {
		http.Error(w, err1.Error(), http.StatusInternalServerError)
	}
}

func datePub() string {
	// Récupérer la date et l'heure actuelle
	heureActuelle := time.Now()

	// Afficher la date et l'heure
	fmt.Println("Date et heure actuelles :", heureActuelle)

	// Vous pouvez également formater la date et l'heure selon vos besoins
	formatPersonnalise := "02/01/2006 15:04:05"
	dateFormatee := heureActuelle.Format(formatPersonnalise)
	fmt.Println("Date et heure formatées :", dateFormatee)

	return dateFormatee
}

func insertPublication(w http.ResponseWriter, r *http.Request) {
	// Vérifier si la méthode HTTP est POST
	if r.Method == "POST" {
		// Récupérer les valeurs du formulaire
		titre := r.FormValue("titre")
		description := r.FormValue("content")
		date := datePub()
		categories := r.FormValue("catego")
		user_email := r.Header.Get("x-email")
		fmt.Println("///////////////////////////////////////////////", titre, description, date, categories, user_email)
		// Se connecter à la base de données
		db := dbConn()
		/*
			// Insérer la publication dans la table T_PUBLICATION_POST
			postQuery := "INSERT INTO T_PUBLICATION_POST (titre, date, description, foreign_key_col) VALUES (?, ?, ?, ?)"
			result, err := db.Exec(postQuery, titre, date, description, userID)
			if err != nil {
				fmt.Println("Erreur lors de l'insertion dans la table T_PUBLICATION_POST:", err)
				http.Error(w, "Erreur lors de l'insertion dans la table T_PUBLICATION_POST", http.StatusInternalServerError)
				return
			}

			// Récupérer l'ID de la publication insérée
			postID, _ := result.LastInsertId()

			// Insérer la relation Many-to-Many dans la table T_CATEGORIEPUB_CATPUB
			catQuery := "INSERT INTO T_CATEGORIEPUB_CATPUB (cat_id, post_id) VALUES ((SELECT id FROM T_CATEGORY_CAT WHERE id = ?), ?)"
			_, err = db.Exec(catQuery, catID, postID)
			if err != nil {
				fmt.Println("Erreur lors de l'insertion dans la table T_CATEGORIEPUB_CATPUB:", err)
				http.Error(w, "Erreur lors de l'insertion dans la table T_CATEGORIEPUB_CATPUB", http.StatusInternalServerError)
				return
			}
		*/
		fmt.Println("Insertion réussie !")

		defer db.Close()
		fmt.Println("************************************** SUCCES *************************************")
		http.Redirect(w, r, "/forum", 302)
	}
}

func affichCat(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(cookieName)
	if err != nil || cookie.Value != authenticated {
		// Si l'utilisateur n'est pas connecté, le rediriger vers la page de connexion
		fmt.Println("*************************** DECONNECTER ******************************")
	} else {
		fmt.Println("*************************** CONNECTER ******************************")
	}

	db := dbConn()
	selDB, err := db.Query("SELECT * FROM T_CATEGORY_CAT ORDER BY cat_id DESC")
	if err != nil {
		panic(err.Error())
	}
	cat := T_CATEGORY_CAT{}
	all_cat := []T_CATEGORY_CAT{}
	for selDB.Next() {
		var cat_id int
		var cat_name string
		err = selDB.Scan(&cat_id, &cat_name)
		if err != nil {
			panic(err.Error())
		}
		cat.Cat_id = cat_id
		cat.Cat_name = cat_name
		all_cat = append(all_cat, cat)

	}
	fmt.Println(all_cat)

	m := map[string]interface{}{
		"cookie": cookie,
		"data":   all_cat,
	}

	err1 := template4.Execute(w, m) // execution de la page index.html
	if err1 != nil {
		http.Error(w, err1.Error(), http.StatusInternalServerError)
	}
}

func affichCat1(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	selDB, err := db.Query("SELECT * FROM T_CATEGORY_CAT ORDER BY cat_id DESC")
	if err != nil {
		panic(err.Error())
	}
	cat := T_CATEGORY_CAT{}
	all_cat := []T_CATEGORY_CAT{}
	for selDB.Next() {
		var cat_id int
		var cat_name string
		err = selDB.Scan(&cat_id, &cat_name)
		if err != nil {
			panic(err.Error())
		}
		cat.Cat_id = cat_id
		cat.Cat_name = cat_name
		all_cat = append(all_cat, cat)

	}
	fmt.Println(all_cat)

	m := map[string]interface{}{
		"data": all_cat,
	}

	err1 := template44.Execute(w, m) // execution de la page index.html
	if err1 != nil {
		http.Error(w, err1.Error(), http.StatusInternalServerError)
	}
}

func getPostByID(postID int) (T_PUBLICATION_POST, error) {
	db := dbConn()

	var post T_PUBLICATION_POST

	err := db.QueryRow("SELECT post_id, post_title, post_date, post_description, user_id FROM T_PUBLICATION_POST WHERE post_id = ?", postID).Scan(
		&post.Post_id, &post.Post_title, &post.Post_date, &post.Post_description, &post.User_id.User_id,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("Post not found")
			return T_PUBLICATION_POST{}, nil
		}
		fmt.Println("Error retrieving post:", err)
		return T_PUBLICATION_POST{}, err
	}

	return post, nil
}

func detailPublication(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(cookieName)
	if err != nil || cookie.Value != authenticated {
		// Si l'utilisateur n'est pas connecté, le rediriger vers la page de connexion
		fmt.Println("*************************** DECONNECTER ******************************")
	} else {
		fmt.Println("*************************** CONNECTER ******************************")
	}

	index := r.URL.Query().Get("posts") // recuperation de l'id
	id_lang, _ := strconv.Atoi(index)
	post, err := getPostByID(id_lang)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("------------------+++++----------------------", post.Post_title)
	m := map[string]interface{}{
		"data":   post,
		"cookie": cookie,
	}

	err1 := template5.Execute(w, m) // execution de la page index.html
	if err1 != nil {
		http.Error(w, err1.Error(), http.StatusInternalServerError)
	}
}

func detailPublication1(w http.ResponseWriter, r *http.Request) {
	index := r.URL.Query().Get("posts") // recuperation de l'id
	id_lang, _ := strconv.Atoi(index)
	post, err := getPostByID(id_lang)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("------------------+++++----------------------", post.Post_title)
	m := map[string]interface{}{
		"data": post,
	}

	err1 := template55.Execute(w, m) // execution de la page index.html
	if err1 != nil {
		http.Error(w, err1.Error(), http.StatusInternalServerError)
	}
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("static/"))) // gere les fichiers statics
	http.HandleFunc("/forum", home)
	http.HandleFunc("/forum0", home0)
	http.HandleFunc("/register", formRegister)
	http.HandleFunc("/post_register", insertUsers)

	http.HandleFunc("/login", formLogin)
	http.HandleFunc("/post_login", LoginHandler)
	http.HandleFunc("/logout", LogoutHandler)

	http.HandleFunc("/creat-publication0", formPublication1)
	http.HandleFunc("/detail-post0", detailPublication1)
	http.HandleFunc("/categorie0", affichCat1)

	http.HandleFunc("/creat-publication", formPublication)
	http.HandleFunc("/post_publication", insertPublication)

	http.HandleFunc("/categorie", affichCat)
	http.HandleFunc("/detail-post", detailPublication)

	fmt.Println("click sur le lien suivant : http://localhost:8080/forum0")
	http.ListenAndServe(port, nil)
}
