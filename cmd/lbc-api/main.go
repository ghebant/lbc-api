package main

import (
	"context"
	"database/sql"
	"ghebant/lbc-api/internal/constants"
	"ghebant/lbc-api/internal/handlers"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/errgo.v2/errors"
	"log"
	"os"
	"time"
)

// TODO
// - Goland supprime les dependancies indirect du go.mod wtf
// - J'arrive pas a me co a la db mysql

func main() {
	router := gin.Default()
	db, err := InitConnectionDB()
	defer db.Close()

	if err != nil {
		log.Fatalf("Failed connect to database: %s", err)
	}

	log.Println("Connected to DB successfully")

	err = CreateAdTable(db)
	if err != nil {
		log.Fatalf("Failed to initialize to database: %s", err)
	}

	router.GET(constants.HealthPath, handlers.Health)
	router.GET(constants.AdPath, handlers.GetAd(db))
	router.GET(constants.AdWithIdPath, handlers.GetAd(db))
	router.POST(constants.AdPath, handlers.PostAd(db))
	router.PUT(constants.AdWithIdPath, handlers.UpdateAd(db))
	router.DELETE(constants.AdWithIdPath, handlers.DeleteAd(db))

	log.Println("Running !")

	err = router.Run(os.Getenv("PORT"))
	if err != nil {
		log.Fatalln(err)
	}
}

// TODO MOVE
func InitConnectionDB() (*sql.DB, error) {
	// Open
	db, err := sql.Open("mysql", "root:root@tcp(db:3306)/lbc-api_development?parseTime=true")
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
		return errors.New("Failed to create table ad: " + err.Error())
	}

	return nil
}
