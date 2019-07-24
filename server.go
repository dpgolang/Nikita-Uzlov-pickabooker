package main

import (
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
	"os"
	"pickabooker/controllers"
	"pickabooker/driver"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/subosito/gotenv"
)

func init() {
	err := gotenv.Load()
	if err != nil {
		log.Println(err)
	}
}

var DB *sqlx.DB

func main() {
	DB = driver.ConnectDB()
	defer DB.Close()

	r := mux.NewRouter()

	ABcontroller := controllers.AbookController{DB}
	Ucontroller := controllers.UserController{DB}

	r.HandleFunc("/register", Ucontroller.Register()).Methods("POST")
	r.HandleFunc("/login", Ucontroller.Login()).Methods("PUT")
	r.HandleFunc("/logout", Ucontroller.Logout()).Methods("PUT")
	r.HandleFunc("/personal", Ucontroller.PersonalAbooks()).Methods("GET")

	r.HandleFunc("/abooks", ABcontroller.GetAbooks()).Methods("GET")
	r.HandleFunc("/abooks/{id}", ABcontroller.PickAbooker()).Methods("PUT")

	r.HandleFunc("/bestsellers", ABcontroller.GetBestsellers()).Methods("GET")

	port := os.Getenv("SERVER_PORT")

	loggedRouter := handlers.LoggingHandler(os.Stdout, r)
	log.Println("listening on " + port)
	log.Fatal(http.ListenAndServe((":" + port), loggedRouter))
}
