package handlers

import (
	"database/sql"
	"ghebant/lbc-api/internal/constants"
	"github.com/gin-gonic/gin"
)

func SetupRouter(db *sql.DB) *gin.Engine {
	router := gin.Default()

	router.GET(constants.HealthPath, Health)
	router.GET(constants.SearchPath, Search)
	router.GET(constants.AdPath, GetAd(db))
	router.GET(constants.AdWithIdPath, GetAd(db))
	router.POST(constants.AdPath, PostAd(db))
	router.PUT(constants.AdWithIdPath, UpdateAd(db))
	router.DELETE(constants.AdWithIdPath, DeleteAd(db))

	return router
}
