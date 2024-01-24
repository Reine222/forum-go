package main

import (
	"database/sql"
	"fmt"

	"html/template"


	"net/http"
	"strconv"
	"time"
	"log"

	// "github.com/gorilla/mux"
    "github.com/gorilla/sessions"
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

var store = sessions.NewCookieStore([]byte("SESSION_SECRET"))
const port = ":8080"


// recuperation des templates (pages html)
var (
	template0  = template.Must(template.ParseFiles("templates/index1.html"))

	template1 = template.Must(template.ParseFiles("templates/page-connexion.html"))
	template2 = template.Must(template.ParseFiles("templates/page-inscription.html"))

	template3 = template.Must(template.ParseFiles("templates/page-create-pub.html"))
	template4 = template.Must(template.ParseFiles("templates/page-categories.html"))
	template5 = template.Must(template.ParseFiles("templates/page-single-topic.html"))
	template6 = template.Must(template.ParseFiles("templates/page-categories-single.html"))
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
	db := dbConn()
	isAuthenticated := false
	var recup_user T_UTILISATEUR_USER
	session, _ := store.Get(r, "session.id")
	fmt.Println("--------------------hh--------------------", session.Values["authenticated"])
	
	if (session.Values["authenticated"] != nil) && session.Values["authenticated"] != false {
		// w.Write([]byte("CONNECTER"))
		isAuthenticated = true
		// Récupérer l'ID de l'utilisateur à partir de la session
		userEmail := session.Values["user"].(string) // Assurez-vous d'adapter cela à votre implémentation
		recup_user, _ = getUserByEmail(db, userEmail)

		fmt.Println("*************************** CONNECTER ******************************", recup_user)
	} else {
		// http.Error(w, "Forbidden", http.StatusForbidden)
		fmt.Println("*************************** DECONNECTER ******************************")
	}

	
	// selDB, err := db.Query("SELECT * FROM T_PUBLICATION_POST ORDER BY post_id DESC")
	selDB, err := db.Query(`
		SELECT 
			T_CATEGORIEPUB_CATPUB.Catpub_id,
			T_CATEGORY_CAT.Cat_id,
			T_CATEGORY_CAT.Cat_name,
			T_PUBLICATION_POST.Post_id,
			T_PUBLICATION_POST.Post_title,
			T_PUBLICATION_POST.Post_date,
			T_PUBLICATION_POST.Post_description,
			T_UTILISATEUR_USER.User_id,
			T_UTILISATEUR_USER.User_name,
			T_UTILISATEUR_USER.User_email,
			T_UTILISATEUR_USER.User_password
		FROM 
			T_CATEGORIEPUB_CATPUB
		INNER JOIN T_CATEGORY_CAT ON T_CATEGORY_CAT.Cat_id = T_CATEGORIEPUB_CATPUB.Cat_id
		INNER JOIN T_PUBLICATION_POST ON T_PUBLICATION_POST.Post_id = T_CATEGORIEPUB_CATPUB.Post_id
		INNER JOIN T_UTILISATEUR_USER ON T_PUBLICATION_POST.User_id = T_UTILISATEUR_USER.User_id;
	`)
	if err != nil {
		fmt.Println("Erreur lors de l'exécution de la requête:", err)
		return
	}
	defer selDB.Close()

	all_pub := []T_CATEGORIEPUB_CATPUB{}

	for selDB.Next() {
		var post T_CATEGORIEPUB_CATPUB
		err := selDB.Scan(
			&post.Catpub_id,
			&post.Cat_id.Cat_id,
			&post.Cat_id.Cat_name,
			&post.Post_id.Post_id,
			&post.Post_id.Post_title,
			&post.Post_id.Post_date,
			&post.Post_id.Post_description,
			&post.Post_id.User_id.User_id,
			&post.Post_id.User_id.User_name,
			&post.Post_id.User_id.User_email,
			&post.Post_id.User_id.User_password,
		)
		if err != nil {
			fmt.Println("Erreur lors du Scan:", err)
			continue
		}

		all_pub = append(all_pub, post)
	}

	fmt.Println("----------------------------------------", all_pub)

	m := map[string]interface{}{
		"data":   all_pub,
		"isAuthenticated": isAuthenticated,
		"user": recup_user,
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
		defer db.Close()
		fmt.Println("************************************** SUCCES *************************************")

		http.Redirect(w, r, "/login", http.StatusSeeOther)

	}

}

// LoginHandler gère le processus d'authentification

func userExists(db *sql.DB, email string) (bool, error) {
	// Exécutez une requête SELECT avec une clause WHERE pour filtrer par nom d'utilisateur
	query := "SELECT COUNT(*) FROM T_UTILISATEUR_USER WHERE user_email = ?"
	var count int

	err := db.QueryRow(query, email).Scan(&count)
	if err != nil {
		return false, err
	}

	// Si le nombre d'utilisateurs correspondant au nom est supérieur à zéro, l'utilisateur existe
	return count > 0, nil
}


func getUserByEmail(db *sql.DB, email string) (T_UTILISATEUR_USER, error) {
	// Exécutez une requête SELECT avec une clause WHERE pour filtrer par email
	query := "SELECT user_id, user_name, user_email, user_password FROM T_UTILISATEUR_USER WHERE user_email = ?"
	var user T_UTILISATEUR_USER

	err := db.QueryRow(query, email).Scan(&user.User_id, &user.User_name, &user.User_email, &user.User_password)
	if err != nil {
		return T_UTILISATEUR_USER{}, err
	}

	return user, nil
}


func LoginHandler(w http.ResponseWriter, r *http.Request) {
	db := dbConn()

	err := r.ParseForm()
	if err != nil {
	   	http.Error(w, "Please pass the data as URL form encoded",http.StatusBadRequest)
	  	return
	}
	email := r.PostForm.Get("email")
	password := r.PostForm.Get("password")

	// Vérifiez si l'utilisateur existe dans la table
	exists, err := userExists(db, email)
	if err != nil {
		log.Fatal(err)
	}
	// originalPassword, ok := users[username]
	if exists {
		fmt.Println("EXIST")
		// Récupérez l'utilisateur par email
		user, err := getUserByEmail(db, email)
		if err != nil {
			log.Fatal(err)
		}
		session, _ := store.Get(r, "session.id")
		errV := bcrypt.CompareHashAndPassword([]byte(user.User_password), []byte(password))
	
		if errV == nil {
			fmt.Println("EXIST")
			session.Values["authenticated"] = true
			session.Values["user"] = user.User_email
			session.Save(r, w)
			fmt.Println("zzzzzzzzzzzzzzzzzzzz", session.Values["authenticated"])
			// Rediriger l'utilisateur vers la page d'accueil après la connexion
			http.Redirect(w, r, "/forum", http.StatusSeeOther)
		} else {
			http.Error(w, "Informations d'identification invalides",http.StatusUnauthorized)
			return
		}
	} else {
		http.Error(w, "L'utilisateur n'a pas été trouvé", http.StatusNotFound)
		return
	}
 
}



// LogoutHandler gère le processus de déconnexion
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
    session, _ := store.Get(r, "session.id")
    session.Values["authenticated"] = false
    session.Save(r, w)
    http.Redirect(w, r, "/forum", http.StatusSeeOther)
}




func formPublication(w http.ResponseWriter, r *http.Request) {
	isAuthenticated := false
	db := dbConn()
	var recup_user T_UTILISATEUR_USER
	session, _ := store.Get(r, "session.id")
	fmt.Println("--------------------hh--------------------", session.Values["authenticated"])
	
	if (session.Values["authenticated"] != nil) && session.Values["authenticated"] != false {
		// w.Write([]byte("CONNECTER"))
		isAuthenticated = true
		// Récupérer l'ID de l'utilisateur à partir de la session
		userEmail := session.Values["user"].(string) // Assurez-vous d'adapter cela à votre implémentation
		recup_user, _ = getUserByEmail(db, userEmail)

		fmt.Println("*************************** CONNECTER ******************************", recup_user)
	} else {
		// http.Error(w, "Forbidden", http.StatusForbidden)
		fmt.Println("*************************** DECONNECTER ******************************")
	}

	m := map[string]interface{}{
		"isAuthenticated": isAuthenticated,
		"user": recup_user,
	}
	err1 := template3.Execute(w, m) // execution de la page index.html
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
	isAuthenticated := false
	
	db := dbConn()
	var recup_user T_UTILISATEUR_USER
	session, _ := store.Get(r, "session.id")
	fmt.Println("--------------------hh--------------------", session.Values["authenticated"])
	
	if (session.Values["authenticated"] != nil) && session.Values["authenticated"] != false {
		// w.Write([]byte("CONNECTER"))
		isAuthenticated = true
		// Récupérer l'ID de l'utilisateur à partir de la session
		userEmail := session.Values["user"].(string) // Assurez-vous d'adapter cela à votre implémentation
		recup_user, _ = getUserByEmail(db, userEmail)

		fmt.Println("*************************** CONNECTER ******************************", recup_user)
	} else {
		// http.Error(w, "Forbidden", http.StatusForbidden)
		fmt.Println("*************************** DECONNECTER ******************************")
	}

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
		"isAuthenticated": isAuthenticated,
		"data":   all_cat,
		"user": recup_user,
	}

	err1 := template4.Execute(w, m) // execution de la page index.html
	if err1 != nil {
		http.Error(w, err1.Error(), http.StatusInternalServerError)
	}
}


func getPostByID(postID int) (T_PUBLICATION_POST, error) {
	db := dbConn()

	var post T_PUBLICATION_POST

	err := db.QueryRow("SELECT T_PUBLICATION_POST.post_id, T_PUBLICATION_POST.post_title, T_PUBLICATION_POST.post_date, T_PUBLICATION_POST.post_description, T_PUBLICATION_POST.user_id, T_UTILISATEUR_USER.user_name, T_UTILISATEUR_USER.user_email, T_UTILISATEUR_USER.user_password FROM T_PUBLICATION_POST INNER JOIN T_UTILISATEUR_USER ON T_PUBLICATION_POST.User_id = T_UTILISATEUR_USER.User_id WHERE T_PUBLICATION_POST.post_id = ?", postID).Scan(
		&post.Post_id, &post.Post_title, &post.Post_date, &post.Post_description, &post.User_id.User_id, &post.User_id.User_name, &post.User_id.User_email, &post.User_id.User_password,
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


// getCategoryByID récupère les informations d'une catégorie par son ID
func getCategoryByID(db *sql.DB, categoryID int) (T_CATEGORY_CAT, error) {
	// Votre requête SQL pour récupérer les informations d'une catégorie par son ID
	query := `
		SELECT cat_id, cat_name
		FROM T_CATEGORY_CAT
		WHERE cat_id = ?;
	`

	var category T_CATEGORY_CAT
	err := db.QueryRow(query, categoryID).Scan(&category.Cat_id, &category.Cat_name)
	if err != nil {
		return T_CATEGORY_CAT{}, err
	}

	return category, nil
}


func getCommentByPostID(postID int) ([]T_COMMENTAIRE_COMMENT, error) {
	db := dbConn()

	rows, err := db.Query("SELECT T_COMMENTAIRE_COMMENT.comment_id, T_COMMENTAIRE_COMMENT.comment_date, T_COMMENTAIRE_COMMENT.comment_description, T_COMMENTAIRE_COMMENT.user_id, T_UTILISATEUR_USER.user_name, T_UTILISATEUR_USER.user_email, T_UTILISATEUR_USER.user_password FROM T_COMMENTAIRE_COMMENT INNER JOIN T_UTILISATEUR_USER ON T_COMMENTAIRE_COMMENT.User_id = T_UTILISATEUR_USER.User_id WHERE T_COMMENTAIRE_COMMENT.post_id = ?", postID)
	if err != nil {
		fmt.Println("Error querying comments:", err)
		return nil, err
	}
	defer rows.Close()

	var comments []T_COMMENTAIRE_COMMENT

	for rows.Next() {
		var comment T_COMMENTAIRE_COMMENT
		err := rows.Scan(
			&comment.Comment_id, &comment.Comment_date, &comment.Comment_description, &comment.User_id.User_id, &comment.User_id.User_name, &comment.User_id.User_email, &comment.User_id.User_password,
		)
		if err != nil {
			fmt.Println("Error scanning comment row:", err)
			return nil, err
		}
		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		fmt.Println("Error retrieving comments:", err)
		return nil, err
	}

	return comments, nil
}



func detailPublication(w http.ResponseWriter, r *http.Request) {
	isAuthenticated := false
	db := dbConn()
	var recup_user T_UTILISATEUR_USER
	session, _ := store.Get(r, "session.id")
	fmt.Println("--------------------hh--------------------", session.Values["authenticated"])
	
	if (session.Values["authenticated"] != nil) && session.Values["authenticated"] != false {
		// w.Write([]byte("CONNECTER"))
		isAuthenticated = true
		// Récupérer l'ID de l'utilisateur à partir de la session
		userEmail := session.Values["user"].(string) // Assurez-vous d'adapter cela à votre implémentation
		recup_user, _ = getUserByEmail(db, userEmail)

		fmt.Println("*************************** CONNECTER ******************************", recup_user)
	} else {
		// http.Error(w, "Forbidden", http.StatusForbidden)
		fmt.Println("*************************** DECONNECTER ******************************")
	}

	indexcat := r.URL.Query().Get("cat") // recuperation de l'id
	id_cat, _ := strconv.Atoi(indexcat)

	index := r.URL.Query().Get("posts") // recuperation de l'id
	id_post, _ := strconv.Atoi(index)

	post, err := getPostByID(id_post)
	catpost, errc := getCategoryByID(db, id_cat)
	comments, errcom := getCommentByPostID(id_post)
	if err != nil && errc != nil && errcom != nil {
		fmt.Println(err)
	}

	fmt.Println("------------------+++++----------------------", comments)
	m := map[string]interface{}{
		"data":   post,
		"catpost": catpost,
		"comments": comments,
		"isAuthenticated": isAuthenticated,
		"user": recup_user,
	}

	err1 := template5.Execute(w, m) // execution de la page index.html
	if err1 != nil {
		http.Error(w, err1.Error(), http.StatusInternalServerError)
	}
}



func detailCatPublication(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	isAuthenticated := false
	var recup_user T_UTILISATEUR_USER
	session, _ := store.Get(r, "session.id")
	fmt.Println("--------------------hh--------------------", session.Values["authenticated"])
	
	if (session.Values["authenticated"] != nil) && session.Values["authenticated"] != false {
		isAuthenticated = true
		// Récupérer l'email de l'utilisateur à partir de la session
		userEmail := session.Values["user"].(string) // Assurez-vous d'adapter cela à votre implémentation
		recup_user, _ = getUserByEmail(db, userEmail)

		fmt.Println("*************************** CONNECTER ******************************", recup_user)
	} else {
		// http.Error(w, "Forbidden", http.StatusForbidden)
		fmt.Println("*************************** DECONNECTER ******************************")
	}

	
	index := r.URL.Query().Get("cat") // recuperation de l'id
	id_cat, _ := strconv.Atoi(index)

	catp, _ := getCategoryByID(db, id_cat)

	selDB, err := db.Query(`
		SELECT 
			T_CATEGORIEPUB_CATPUB.Catpub_id,
			T_CATEGORY_CAT.Cat_id,
			T_CATEGORY_CAT.Cat_name,
			T_PUBLICATION_POST.Post_id,
			T_PUBLICATION_POST.Post_title,
			T_PUBLICATION_POST.Post_date,
			T_PUBLICATION_POST.Post_description,
			T_UTILISATEUR_USER.User_id,
			T_UTILISATEUR_USER.User_name,
			T_UTILISATEUR_USER.User_email,
			T_UTILISATEUR_USER.User_password
		FROM 
			T_CATEGORIEPUB_CATPUB
		INNER JOIN T_CATEGORY_CAT ON T_CATEGORY_CAT.Cat_id = T_CATEGORIEPUB_CATPUB.Cat_id
		INNER JOIN T_PUBLICATION_POST ON T_PUBLICATION_POST.Post_id = T_CATEGORIEPUB_CATPUB.Post_id
		INNER JOIN T_UTILISATEUR_USER ON T_PUBLICATION_POST.User_id = T_UTILISATEUR_USER.User_id
		WHERE
			T_CATEGORY_CAT.Cat_id = ?;
	`, id_cat)
	if err != nil {
		fmt.Println("Erreur lors de l'exécution de la requête:", err)
		return
	}
	defer selDB.Close()

	all_pub := []T_CATEGORIEPUB_CATPUB{}

	for selDB.Next() {
		var post T_CATEGORIEPUB_CATPUB
		err := selDB.Scan(
			&post.Catpub_id,
			&post.Cat_id.Cat_id,
			&post.Cat_id.Cat_name,
			&post.Post_id.Post_id,
			&post.Post_id.Post_title,
			&post.Post_id.Post_date,
			&post.Post_id.Post_description,
			&post.Post_id.User_id.User_id,
			&post.Post_id.User_id.User_name,
			&post.Post_id.User_id.User_email,
			&post.Post_id.User_id.User_password,
		)
		if err != nil {
			fmt.Println("Erreur lors du Scan:", err)
			continue
		}

		all_pub = append(all_pub, post)
	}

	fmt.Println("----------------------------------------", all_pub)

	m := map[string]interface{}{
		"data":   all_pub,
		"catp":   catp,
		"isAuthenticated": isAuthenticated,
		"user": recup_user,
	}

	err1 := template6.Execute(w, m) // execution de la page index.html
	if err1 != nil {
		http.Error(w, err1.Error(), http.StatusInternalServerError)
	}
}



func insertComment(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	session, _ := store.Get(r, "session.id")
	fmt.Println("--------------------hh--------------------", session.Values["authenticated"])

	// Obtenir la date actuelle et l'heure
	// Formater la date pour l'afficher de manière lisible
	currentTime := time.Now()
	formattedDate := currentTime.Format("02 January 2006, 15:04:05")
	
	if (session.Values["authenticated"] != nil) && session.Values["authenticated"] != false {

		userEmail := session.Values["user"].(string) // Assurez-vous d'adapter cela à votre implémentation
		recup_user, _ := getUserByEmail(db, userEmail)

		index := r.URL.Query().Get("post") // recuperation de l'id
		id_post, _ := strconv.Atoi(index)
		poste_id, _ := getPostByID(id_post)

		indexc := r.URL.Query().Get("cat") // recuperation de l'id

		fmt.Println("*************************** CONNECTER ******************************", recup_user)
	
		if r.Method == "POST" {
			comment_description := r.FormValue("comment")
			comment_date := formattedDate
			post_id := poste_id.Post_id
			user_id := recup_user.User_id
			fmt.Println(comment_description)

			if len(comment_description) > 0 {
				
				commentaire := "INSERT INTO T_COMMENTAIRE_COMMENT (comment_date, comment_description, post_id, user_id) VALUES (?,?,?,?)"
				_, errr := db.Exec(commentaire, comment_date, comment_description, post_id, user_id)
				if errr != nil {
					fmt.Println("Erreur lors de l'insertion dans la table users:", errr)
					return
				}

				fmt.Println("Insertion réussie !")
				defer db.Close()
				fmt.Println("************************************** SUCCES *************************************")

			}

			http.Redirect(w, r, "/detail-post?posts=" + index + "&cat=" + indexc, http.StatusSeeOther)

		}
	}

}


func main() {
	http.Handle("/", http.FileServer(http.Dir("static/"))) // gere les fichiers statics
	http.HandleFunc("/forum", home)
	
	http.HandleFunc("/register", formRegister)
	http.HandleFunc("/post_register", insertUsers)

	http.HandleFunc("/login", formLogin)
	http.HandleFunc("/post_login", LoginHandler)
	http.HandleFunc("/logout", LogoutHandler)

	http.HandleFunc("/creat-publication", formPublication)
	http.HandleFunc("/post_publication", insertPublication)
	http.HandleFunc("/post_comment", insertComment)


	http.HandleFunc("/categorie", affichCat)
	http.HandleFunc("/detail-categorie-post", detailCatPublication)
	http.HandleFunc("/detail-post", detailPublication)

	fmt.Println("click sur le lien suivant : http://localhost:8080/forum")
	http.ListenAndServe(port, nil)
}
