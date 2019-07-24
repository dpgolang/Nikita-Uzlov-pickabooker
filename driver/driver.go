package driver

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"os"
	"pickabooker/utils"
)

func ConnectDB() *sqlx.DB {
	var (
		dbhost     = os.Getenv("POSTGRES_HOST")
		dbport     = os.Getenv("POSTGRES_PORT")
		dbuser     = os.Getenv("POSTGRES_USER")
		dbpassword = os.Getenv("POSTGRES_PASS")
		dbname     = os.Getenv("POSTGRES_NAME")
	)

	var err error
	var db *sqlx.DB

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbhost, dbport, dbuser, dbpassword, dbname)

	db, err = sqlx.Open("postgres", psqlInfo)
	utils.LogError(err)

	err = db.Ping()
	utils.LogError(err)
	fmt.Println("Successfully connected to pg_docker database.")
	return db
}
