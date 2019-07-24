package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
	"github.com/subosito/gotenv"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"pickabooker/models"
	"pickabooker/repository"
	"pickabooker/utils"
)

type UserController struct {
	DB *sqlx.DB
}
type AuthUsers struct {
	models.User
	Auth bool
}

var store *sessions.CookieStore

func init() {
	err := gotenv.Load()
	utils.LogError(err)
	key := []byte(os.Getenv("SESSION_KEY"))
	store = sessions.NewCookieStore(key)
}

func (c *UserController) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "cookie-name")
		session.Options.MaxAge = 300

		var unknownU models.User

		err := json.NewDecoder(r.Body).Decode(&unknownU)
		utils.LogError(err)

		users, err := repository.GetUsers(c.DB)
		utils.LogError(err)

		for _, user := range users {
			if repository.EqualUsers(user, unknownU) {
				session.Values["authenticated"] = true
				session.Values["id"] = user.ID
				err := session.Save(r, w)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(fmt.Sprintf("Logged as %d", user.ID))
				return
			}
		}

		http.Error(w, "No such user.", http.StatusUnauthorized)
	}
}

func (c *UserController) Logout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "cookie-name")
		session.Options.MaxAge = -1
		session.Values["authenticated"] = false
		err := session.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(fmt.Sprintf("Logged out"))
	}
}

func (c *UserController) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var udata models.User
		err := json.NewDecoder(r.Body).Decode(&udata)
		utils.LogError(err)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		hashPass, err := bcrypt.GenerateFromPassword([]byte(udata.Password), 8)
		userID, exists := repository.AddNewUser(c.DB, udata.Login, string(hashPass))
		if exists {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode("User with this login already exists")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(&userID)
	}
}
