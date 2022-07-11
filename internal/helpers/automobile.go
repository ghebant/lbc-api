package helpers

import (
	"database/sql"
	"ghebant/lbc-api/internal/constants"
	"ghebant/lbc-api/models"
	"gopkg.in/errgo.v2/errors"
	"time"
)

func InsertAutomobile(adId int, auto *models.Automobile, db *sql.DB) (int, error) {
	res, err := db.Exec("INSERT INTO automobile(ad_id, brand, model) VALUES (?, ?, ?)", adId, auto.Brand, auto.Model)
	if err != nil {
		return -1, err
	}

	lastId, _ := res.LastInsertId()

	return int(lastId), nil
}

func InsertAndReturnAutomobile(adId int, automobile *models.Automobile, db *sql.DB) (models.Automobile, error) {
	if automobile == nil {
		return models.Automobile{}, errors.New("failed to insert automobile: automobile is nil")
	}

	autoMobileId, err := InsertAutomobile(adId, automobile, db)
	if err != nil {
		return models.Automobile{}, errors.New("error failed to insert automobile in db: " + err.Error())
	}

	res, err := FindAutomobileById(db, autoMobileId)
	if err != nil {
		return models.Automobile{}, errors.New("failed to retrieve created automobile: " + err.Error())
	}

	return res, nil
}

func FindAutomobileById(db *sql.DB, id int) (models.Automobile, error) {
	auto := models.Automobile{}

	query := "SELECT * FROM automobile WHERE " + constants.AutomobilePrimaryKey + " = ?"
	row := db.QueryRow(query, id)

	err := row.Scan(&auto.ID, &auto.AdId, &auto.Brand, &auto.Model, &auto.CreatedAt, &auto.UpdatedAt)
	if err != nil {
		return models.Automobile{}, err
	}

	return auto, nil
}

func FindAutomobileByAdId(db *sql.DB, id int) (models.Automobile, error) {
	auto := models.Automobile{}

	query := "SELECT * FROM automobile WHERE ad_id = ?"
	row := db.QueryRow(query, id)

	err := row.Scan(&auto.ID, &auto.AdId, &auto.Brand, &auto.Model, &auto.CreatedAt, &auto.UpdatedAt)
	if err != nil {
		return models.Automobile{}, err
	}

	return auto, nil
}

func UpdateAutomobile(adId int, automobile *models.Automobile, db *sql.DB) error {
	res, err := db.Exec("UPDATE automobile SET brand = ?, model = ?, updated_at = ? WHERE ad_id = ?", automobile.Brand, automobile.Model, time.Now(), adId)
	if err != nil {
		return err
	}

	row, _ := res.RowsAffected()
	if row <= 0 {
		_, err := InsertAutomobile(adId, automobile, db)
		if err != nil {
			return errors.New("error failed to insert automobile in db: " + err.Error())
		}
	}

	return nil
}
