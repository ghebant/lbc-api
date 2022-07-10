package helpers

import (
	"database/sql"
	"ghebant/lbc-api/internal/constants"
	"ghebant/lbc-api/models"
	"gopkg.in/errgo.v2/errors"
	"log"
)

func InsertAds(ads []models.Ad, db *sql.DB) error {
	stmt, err := db.Prepare("INSERT INTO ad(title, content, category) VALUES (?, ?, ?)")
	if err != nil {
		return errors.New("failed to create sql prepare statement: " + err.Error())
	}

	for i := range ads {
		res, err := stmt.Exec(ads[i].Title, ads[i].Content, ads[i].Category)
		if err != nil {
			return errors.New("failed to insert mock ad in db: " + err.Error())
		}
		lastId, _ := res.LastInsertId()
		ads[i].ID = int(lastId)

		if ads[i].Automobile != nil {
			_, err = InsertAutomobile(int(lastId), ads[i].Automobile, db)
			if err != nil {
				return errors.New("failed to insert mock automobile in db: " + err.Error())
			}
		}
	}

	return nil
}

func FindAdById(db *sql.DB, id int) (models.Ad, error) {
	ad := models.Ad{}

	query := "SELECT * FROM ad WHERE " + constants.AdPrimaryKey + " = ?"
	row := db.QueryRow(query, id)

	err := row.Scan(&ad.ID, &ad.Title, &ad.Content, &ad.Category, &ad.CreatedAt, &ad.UpdatedAt)
	if err != nil {
		return models.Ad{}, err
	}

	switch ad.Category {
	case constants.Automobile:
		automobile, err := FindAutomobileByAdId(db, ad.ID)
		if err != nil {
			log.Println("failed to retrieve automobile from ad: " + err.Error())
		}

		ad.Automobile = &automobile
	}

	return ad, nil
}
