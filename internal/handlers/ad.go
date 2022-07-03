package handlers

import (
	"database/sql"
	"ghebant/lbc-api/internal/helpers"
	"ghebant/lbc-api/models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"time"
)

func GetAd(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		ads := []models.Ad{}

		// TODO everywhere ?
		c.Header("Content-Type", "application/json")

		// Find ad by ID
		idStr := c.Param("id")
		if idStr != "" {
			var ad models.Ad

			adId, err := strconv.Atoi(idStr)
			if err != nil {
				log.Println("failed to convert id:", err)
				c.JSON(http.StatusBadRequest, gin.H{"message": "failed to convert id"})
				return
			}

			ad, err = helpers.FindAdById(db, adId)
			if err != nil {
				log.Println("no ad found for the provided id:", err)
				c.JSON(http.StatusNotFound, gin.H{"message": "no ad found for the provided id"})
				return
			}

			c.JSON(http.StatusOK, ad)
			return
		}

		queryRes, err := db.Query("SELECT * FROM ad")
		if err != nil {
			log.Println("error failed to query db:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Cannot fetch ads"})
			return
		}
		defer queryRes.Close()

		for queryRes.Next() {
			var ad models.Ad

			err = queryRes.Scan(&ad.ID, &ad.Title, &ad.Content, &ad.Category, &ad.CreatedAt, &ad.UpdatedAt)
			if err != nil {
				log.Println("no ad found:", err)
				c.JSON(http.StatusNotFound, gin.H{"message": "no ad found"})
				return
			}

			ads = append(ads, ad)
		}

		c.JSON(http.StatusOK, ads)
	}
}

func PostAd(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		ad := models.Ad{}

		err := c.ShouldBindJSON(&ad)
		if err != nil {
			log.Println("failed to read body", err)
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		res, err := db.Exec("INSERT INTO ad(title, content, category) VALUES (?, ?, ?)", ad.Title, ad.Content, ad.Category)
		if err != nil {
			log.Println("error failed to insert ad in db:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "error failed to insert ad in db"})
			return
		}

		// Return created ad
		lastId, _ := res.LastInsertId()

		ad, err = helpers.FindAdById(db, int(lastId))
		if err != nil {
			log.Println("failed to retrieve created ad:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to retrieve created ad"})
			return
		}

		c.JSON(http.StatusCreated, ad)
	}
}

func UpdateAd(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		ad := models.Ad{}

		idStr := c.Param("id")
		adId, err := strconv.Atoi(idStr)
		if err != nil {
			log.Println("failed to convert id:", err)
			c.JSON(http.StatusBadRequest, gin.H{"message": "failed to convert id"})
			return
		}

		err = c.ShouldBindJSON(&ad)
		if err != nil {
			log.Println("failed to read body", err)
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		res, err := db.Exec("UPDATE ad SET title = ?, content = ?, category = ?, updated_at = ? WHERE ad_id = ?", ad.Title, ad.Content, ad.Category, time.Now(), adId)
		if err != nil {
			log.Println("error failed to update ad in db:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "error failed to update ad in db"})
			return
		}

		row, _ := res.RowsAffected()
		if row <= 0 {
			log.Println("failed to update ad: not found :", err)
			c.JSON(http.StatusNotFound, gin.H{"message": "failed to update ad: not found"})
			return
		}

		// Return created ad
		ad, err = helpers.FindAdById(db, adId)
		if err != nil {
			log.Println("failed to retrieve updated ad:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to retrieve updated ad"})
			return
		}

		c.JSON(http.StatusOK, ad)
	}
}

func DeleteAd(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		adId, err := strconv.Atoi(idStr)
		if err != nil {
			log.Println("failed to convert id:", err)
			c.JSON(http.StatusBadRequest, gin.H{"message": "failed to convert id"})
			return
		}

		res, err := db.Exec("DELETE FROM ad WHERE ad_id = ?", adId)
		if err != nil {
			log.Println("error failed to delete ad in db:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "error failed to delete ad in db"})
			return
		}

		row, _ := res.RowsAffected()
		if row <= 0 {
			log.Println("failed to delete ad: not found :", err)
			c.JSON(http.StatusNotFound, gin.H{"message": "failed to delete ad: not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "deleted"})
	}
}
