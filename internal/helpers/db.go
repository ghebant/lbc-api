package helpers

import (
	"context"
	"database/sql"
	"fmt"
	"ghebant/lbc-api/internal/constants"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/errgo.v2/errors"
	"os"
	"time"
)

func InitDB() (*sql.DB, error) {
	mysqlUser := os.Getenv("MYSQL_ROOT_USER")
	mysqlPwd := os.Getenv("MYSQL_PASSWORD")
	mysqlDb := os.Getenv("MYSQL_DATABASE")
	mysqlHost := os.Getenv("MYSQL_HOST")
	mysqlPort := os.Getenv("MYSQL_PORT")

	mysqlURI := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", mysqlUser, mysqlPwd, mysqlHost, mysqlPort, mysqlDb)

	// Open
	db, err := sql.Open(constants.MysqlDriver, mysqlURI)
	if err != nil {
		return nil, errors.New("Failed to open connection to database: " + err.Error())
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	//Ping db
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err = db.PingContext(ctx)
	if err != nil {
		return nil, errors.New("Failed to ping database: " + err.Error())
	}

	return db, nil
}

func CreateTables(db *sql.DB) error {
	tableQueries := []string{
		"CREATE TABLE IF NOT EXISTS ad(ad_id int primary key auto_increment, " +
			"title text, " +
			"content text, " +
			"category int, " +
			"created_at TIMESTAMP default CURRENT_TIMESTAMP, " +
			"updated_at TIMESTAMP default CURRENT_TIMESTAMP)",

		"CREATE TABLE IF NOT EXISTS automobile(automobile_id int primary key auto_increment, " +
			"ad_id int, " +
			"brand text, " +
			"model text, " +
			"created_at TIMESTAMP default CURRENT_TIMESTAMP, " +
			"updated_at TIMESTAMP default CURRENT_TIMESTAMP)",
	}

	for i := range tableQueries {
		_, err := db.Exec(tableQueries[i])
		if err != nil {
			return errors.New("Failed to create table: " + err.Error())
		}
	}

	return nil
}

func SetupDbForTesting() (*sql.DB, error) {
	db, err := InitDB()
	if err != nil {
		return nil, errors.New("failed connect to database: " + err.Error())
	}

	db.QueryRow("DELETE FROM ad;")

	err = CreateTables(db)
	if err != nil {
		return nil, errors.New("failed to initialize to database: " + err.Error())
	}

	return db, nil
}
