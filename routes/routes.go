package routes

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/http"
	"github.com/gorilla/mux"
	"../middleware"
	"../models"
	"../sessions"
	"../utils"
	"os"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter()
	r. HandleFunc("/", middleware.AuthRequired(indexGetHandler)).Methods("GET")
	r. HandleFunc("/", middleware.AuthRequired(indexPostHandler)).Methods("POST")
	r. HandleFunc("/login", loginGetHandler).Methods("GET")
	r. HandleFunc("/login", loginPostHandler).Methods("POST")
	r. HandleFunc("/logout", logoutGetHandler).Methods("GET")
	r. HandleFunc("/register", registerGetHandler).Methods("GET")
	r. HandleFunc("/register", registerPostHandler).Methods("POST")
	fs := http.FileServer(http.Dir("./static/"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))
	r. HandleFunc("/{username}", middleware.AuthRequired(userGetHandler)).Methods("GET")
	return r
}

func indexGetHandler(w http.ResponseWriter, r *http.Request) {
	updates, err := models.GetAllUpdates()
	if err != nil {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
		utils.InternalServerError(w)
		return
	}
	utils.ExecuteTemplate(w, "index.html", struct{
		Title string
		Updates []*models.Update
		DisplayForm bool
	} {
		Title: "All updates",
		Updates: updates,
		DisplayForm: true,
	})
}

func indexPostHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := sessions.Store.Get(r, "session")
	untypeduserId := session.Values["user_id"]
	userId, ok := untypeduserId.(int64)
	if !ok {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
		utils.InternalServerError(w)
        return
	}
	r.ParseForm()
	body := r.PostForm.Get("update")
	err := models.PostUpdate(userId, body)
	if err != nil {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
		utils.InternalServerError(w)
        return
	}
	http.Redirect(w, r, "/", 302)
}

func userGetHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := sessions.Store.Get(r, "session")
	untypeduserId := session.Values["user_id"]
	currentUserId, ok := untypeduserId.(int64)
	if !ok {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
		utils.InternalServerError(w)
        return
	}
	vars := mux.Vars(r)
	username := vars["username"]
	user, err := models.GetUserByUsername(username)
	if err != nil {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
		utils.InternalServerError(w)
		return
	}

	userId, err := user.GetId()
	if err != nil {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
		utils.InternalServerError(w)
		return
	}
	updates, err := models.GetUpdates(userId)
	if err != nil {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
		utils.InternalServerError(w)
		return
	}
	utils.ExecuteTemplate(w, "index.html", struct{
		Title string
		Updates []*models.Update
		DisplayForm bool
	} {
		Title: username,
		Updates: updates,
		DisplayForm: currentUserId == userId,
	})
}

func loginGetHandler(w http.ResponseWriter, r *http.Request) {
	utils.ExecuteTemplate(w, "login.html", nil)
}

func loginPostHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.PostForm.Get("username")
	password := r.PostForm.Get("password")
	user, err := models.AuthenticateUser(username, password)
	if err != nil {
		switch err {
		case models.ErrUserNotFound:
			log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
			utils.ExecuteTemplate(w, "login.html", "Unknown user")
			log.Error().Msgf("error: %s\n", err.Error())
		case models.ErrInvalidLogin:
			log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
			utils.ExecuteTemplate(w, "login.html", "Invalid login")
		default:
			log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
			utils.InternalServerError(w)
		}
		return
	}
	userId, err := user.GetId()
	if err != nil {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
		utils.InternalServerError(w)
		return
	}
	session, _ := sessions.Store.Get(r, "session")
	session.Values["user_id"] = userId
	session.Save(r, w)
	http.Redirect(w, r, "/", 302)
}

func logoutGetHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := sessions.Store.Get(r, "session")
	delete(session.Values, "user_id")
	session.Save(r, w)
	http.Redirect(w, r, "/login", 302)
}

func registerGetHandler(w http.ResponseWriter, r *http.Request) {
	utils.ExecuteTemplate(w, "register.html", nil)
}

func registerPostHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.PostForm.Get("username")
	password := r.PostForm.Get("password")
	err := models.RegisterUser(username, password)
	if err == models.ErrUsernameTaken {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
		utils.ExecuteTemplate(w, "register.html", "Username taken")
		return
	} else if err != nil {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
		utils.InternalServerError(w)
        return
	}
	http.Redirect(w, r, "/login", 302)
}
