package helpers

import (
	"database/sql"
	"gopkg.in/errgo.v2/errors"
)

func SetupDbForTesting() (*sql.DB, error) {
	db, err := InitDB()
	if err != nil {
		return nil, errors.New("failed connect to database: " + err.Error())
	}

	db.QueryRow("DELETE FROM ad;")

	err = CreateAdTable(db)
	if err != nil {
		return nil, errors.New("failed to initialize to database: " + err.Error())
	}

	return db, nil
}
