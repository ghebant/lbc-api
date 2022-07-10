package handlers

import (
	"ghebant/lbc-api/internal/constants"
	"ghebant/lbc-api/internal/matching"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func Search(c *gin.Context) {
	queryParam := c.Query("input")

	if len(queryParam) <= 0 {
		log.Println("input is empty")
		c.JSON(http.StatusBadRequest, gin.H{"message": "an input must be provided"})
		return
	}

	match := matching.FindBestMatch(queryParam, &constants.Vehicles)

	c.JSON(http.StatusOK, gin.H{"brand": match.Brand, "model": match.Model})
}
