package helpers

import (
	"context"
	"database/sql"
	"fmt"
	"ghebant/lbc-api/internal/constants"
	"ghebant/lbc-api/models"
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

func CreateAdTable(db *sql.DB) error {
	createTableQuery := "CREATE TABLE IF NOT EXISTS ad(ad_id int primary key auto_increment, " +
		"title text, " +
		"content text, " +
		"category text, " +
		"created_at TIMESTAMP default CURRENT_TIMESTAMP, " +
		"updated_at TIMESTAMP default CURRENT_TIMESTAMP)"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := db.ExecContext(ctx, createTableQuery)
	if err != nil {
		return errors.New("Failed to create ad table: " + err.Error())
	}

	return nil
}

func CreateCategoryTable(db *sql.DB) error {
	createTableQuery := "CREATE TABLE IF NOT EXISTS category(category_id int primary key auto_increment, " +
		"brand text, " +
		"model text, " +
		"created_at TIMESTAMP default CURRENT_TIMESTAMP, " +
		"updated_at TIMESTAMP default CURRENT_TIMESTAMP)"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := db.ExecContext(ctx, createTableQuery)
	if err != nil {
		return errors.New("Failed to create category table: " + err.Error())
	}

	return nil
}

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
	}

	return nil
}

func InsertCategories(db *sql.DB) error {
	stmt, err := db.Prepare("INSERT INTO category(brand, model) VALUES (?, ?)")
	if err != nil {
		return errors.New("failed to create sql prepare statement: " + err.Error())
	}

	for brand, model := range constants.Vehicles {
		_, err := stmt.Exec(brand, model)
		if err != nil {
			return errors.New("failed to insert vehicles in db: " + err.Error())
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

	return ad, nil
}

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
