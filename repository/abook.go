package repository

import (
	"github.com/jmoiron/sqlx"
	"pickabooker/models"
)

func GetAbooks(db *sqlx.DB) (abooks []models.Abook, err error) {
	rows, err := db.Queryx("SELECT * FROM audiobooks ORDER BY duration DESC")
	defer rows.Close()

	for rows.Next() {
		var abook models.Abook
		err := rows.StructScan(&abook)
		abooks = append(abooks, abook)
		if err != nil {
			return []models.Abook{}, err
		}
	}

	return abooks, err
}

func GetAbooksLess(db *sqlx.DB, duration string) (vacants []models.Abook, err error) {

	rows, err := db.Queryx("SELECT * FROM audiobooks where duration <= $1 order by duration DESC", duration)
	defer rows.Close()

	for rows.Next() {
		var vacant models.Abook
		err := rows.StructScan(&vacant)
		vacants = append(vacants, vacant)
		if err != nil {
			return []models.Abook{}, err
		}
	}

	return
}

func GetBestsellers(db *sqlx.DB) (bestsellers []models.Bestsellers, err error) {
	rows, err := db.Queryx("SELECT copies, id, title, duration, author, narrator, price FROM audiobooks RIGHT JOIN sales ON audiobooks.id = sales.abookid ORDER BY copies DESC;")
	defer rows.Close()

	for rows.Next() {
		var bestseller models.Bestsellers
		err := rows.StructScan(&bestseller)
		bestsellers = append(bestsellers, bestseller)
		if err != nil {
			return []models.Bestsellers{}, err
		}
	}

	return
}

func PickAbook(db *sqlx.DB, abookID, userID int) bool {
	result, err := db.Exec("INSERT INTO sales (abookid, copies) VALUES ($1, 1) ON CONFLICT (abookid) DO UPDATE SET copies = sales.copies + 1", abookID)
	if err != nil {
		return false
	}

	rowsUpd, err := result.RowsAffected()
	if err != nil || rowsUpd != 1 {
		return false
	}

	result, err = db.Exec("INSERT INTO orders (userid, abookid) VALUES ($1, $2)", userID, abookID)
	if err != nil {
		return false
	}

	rowsUpd, err = result.RowsAffected()
	if err != nil || rowsUpd != 1 {
		return false
	}

	return true
}

func GetUserAbooks(db *sqlx.DB, userID int) (userAbooks []models.Abook, err error) {
	rows, err := db.Queryx("SELECT * FROM audiobooks WHERE id IN (SELECT abookid FROM orders WHERE userid = $1)", userID)
	defer rows.Close()

	for rows.Next() {
		var uabook models.Abook
		err = rows.StructScan(&uabook)
		userAbooks = append(userAbooks, uabook)
		if err != nil {
			return []models.Abook{}, err
		}
	}

	return

}
