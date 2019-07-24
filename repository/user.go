package repository

import (
	"database/sql"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"pickabooker/models"
	"pickabooker/utils"

	"github.com/jmoiron/sqlx"
)

func GetUsers(db *sqlx.DB) (users []models.User, err error) {
	rows, err := db.Queryx("SELECT login, password FROM users ")
	defer rows.Close()

	for rows.Next() {
		var user models.User
		err := rows.StructScan(&user)
		users = append(users, user)
		if err != nil {
			return []models.User{}, err
		}
	}

	return users, err
}

func AddNewUser(db *sqlx.DB, log, pass string) (uid int, exists bool) {
	err := db.QueryRow("SELECT EXISTS (SELECT * FROM users WHERE login = $1)", log).Scan(&exists)
	utils.LogError(fmt.Errorf("Trouble discovering if user exists in database. %v", err))
	fmt.Println(log, exists)
	if exists {
		return
	}
	err = db.QueryRow("INSERT INTO users (login, password) VALUES ($1, $2) RETURNING id", log, pass).Scan(&uid)
	utils.LogError(fmt.Errorf("Can't insert user to database: %v", err))
	return
}

func IsAbookInPersonal(db *sqlx.DB, abookID, userID int) bool {
	var exists bool
	err := db.QueryRow("SELECT EXISTS (SELECT * FROM orders WHERE abookid = $1 AND userid = $2)", abookID, userID).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		utils.LogError(fmt.Errorf("error checking if row exists %v", err))
	}
	return exists
}

func EqualUsers (user1, user2 models.User) bool {

	if user1.Login != user2.Login {
		fmt.Println("Not equal by name")
		return false
	}

	err := bcrypt.CompareHashAndPassword([]byte(user1.Password), []byte(user2.Password))
	if err != nil {
		fmt.Println("Not equal by pass")
		return false
	}

	return true
}
