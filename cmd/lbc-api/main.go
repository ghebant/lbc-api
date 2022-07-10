package main

import (
	"ghebant/lbc-api/internal/handlers"
	"ghebant/lbc-api/internal/helpers"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
)

func main() {
	db, err := helpers.InitDB()
	defer db.Close()

	if err != nil {
		log.Fatalf("Failed connect to database: %s", err)
	}

	log.Println("Connected to DB successfully")

	err = helpers.CreateAdTable(db)
	if err != nil {
		log.Fatalf("Failed to initialize to database: %s", err)
	}

	err = helpers.CreateCategoryTable(db)
	if err != nil {
		log.Fatalf("Failed to initialize to database: %s", err)
	}

	err = helpers.InsertCategories(db)
	if err != nil {
		log.Fatalf("Failed to initialize to database: %s", err)
	}

	router := handlers.SetupRouter(db)
	log.Println("Running !")

	err = router.Run(os.Getenv("PORT"))
	if err != nil {
		log.Fatalln(err)
	}
}
