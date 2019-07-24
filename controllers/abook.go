package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"net/http"
	"pickabooker/models"
	"pickabooker/repository"
	"pickabooker/utils"
	"strconv"
)

type AbookController struct {
	DB *sqlx.DB
}

func (c *AbookController) GetAbooks() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var abooks []models.Abook
		var err error

		if keys, ok := r.URL.Query()["duration"]; ok {
			abooks, err = repository.GetAbooksLess(c.DB, keys[0])
			utils.LogError(err)
		} else {
			abooks, err = repository.GetAbooks(c.DB)
			utils.LogError(err)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(abooks)
	}
}

func (c *AbookController) GetBestsellers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var bestsellers []models.Bestsellers

		bestsellers, err := repository.GetBestsellers(c.DB)
		utils.LogError(err)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(bestsellers)
	}
}

func (c *AbookController) PickAbooker() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "cookie-name")
		uid := session.Values["id"]

		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			fmt.Println(ok, auth)
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		} else {

			userID := uid.(int)
			keys, ok := r.URL.Query()["id"]

			if !ok || len(keys[0]) < 1 {
				utils.LogError(errors.New("Url Param 'id' is missing"))
				return
			}

			abookID, err := strconv.Atoi(keys[0])

			if err != nil {
				utils.LogError(errors.New("id is not an integer number"))
			}

			if repository.IsAbookInPersonal(c.DB, abookID, userID) {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode("audio-book is picked already")
				return
			}

			if abookPicked := repository.PickAbook(c.DB, abookID, userID); abookPicked {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode("audio-book is picked")

			} else {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode("audio-book is not picked")
			}
		}
	}
}

func (c *UserController) PersonalAbooks() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "cookie-name")

		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		var userAbooks []models.Abook
		userID := session.Values["id"].(int)

		userAbooks, err := repository.GetUserAbooks(c.DB, userID)
		utils.LogError(err)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(userAbooks)
	}
}
