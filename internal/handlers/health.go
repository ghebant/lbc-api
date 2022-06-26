package handlers

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func Health(c *gin.Context) {
	// TODO Remove
	log.Println("/health !")
	c.Status(http.StatusOK)
}
