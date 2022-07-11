package handlers

import (
	"database/sql"
	"ghebant/lbc-api/internal/constants"
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

		// Find all ads
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

			switch ad.Category {
			case constants.Automobile:
				automobile, err := helpers.FindAutomobileByAdId(db, ad.ID)
				if err != nil {
					log.Println("failed to retrieve automobile from ad: " + err.Error())
				}

				ad.Automobile = &automobile
			}

			ads = append(ads, ad)
		}

		c.JSON(http.StatusOK, ads)
	}
}

func PostAd(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		payload := models.Ad{}

		err := c.ShouldBindJSON(&payload)
		if err != nil {
			log.Println("failed to read body", err)
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		errMsg := ""
		errorCategory := false

		switch payload.Category {
		case constants.Automobile:
			if payload.Automobile == nil {
				errorCategory = true
				errMsg = "automobile cannot be nil"
			}
		case constants.RealEstate:
			if payload.RealEstate == nil {
				errorCategory = true
				errMsg = "real estate cannot be nil"
			}
		case constants.Job:
			if payload.Job == nil {
				errorCategory = true
				errMsg = "job cannot be nil"
			}
		default:
			errorCategory = true
			errMsg = "wrong category provided"
		}

		if errorCategory {
			log.Println(errMsg)
			c.JSON(http.StatusBadRequest, gin.H{"message": "failed to create ad: " + errMsg})
			return
		}

		txn, err := db.Begin()
		if err != nil {
			log.Println("error failed to start transaction:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "error failed to start transaction"})
			return
		}

		// Insert Ad
		res, err := db.Exec("INSERT INTO ad(title, content, category) VALUES (?, ?, ?)", payload.Title, payload.Content, payload.Category)
		if err != nil {
			log.Println("error failed to insert ad in db:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "error failed to insert ad in db"})
			return
		}

		lastId, _ := res.LastInsertId()

		ad, err := helpers.FindAdById(db, int(lastId))
		if err != nil {
			log.Println("failed to retrieve created ad:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to retrieve created ad"})
			return
		}

		switch payload.Category {
		case constants.Automobile:
			autoMobile, err := helpers.InsertAndReturnAutomobile(ad.ID, payload.Automobile, db)
			if err != nil {
				log.Println(err.Error())
				txn.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			}

			ad.Automobile = &autoMobile
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

		txn, err := db.Begin()
		if err != nil {
			log.Println("error failed to start transaction:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "error failed to start transaction"})
			return
		}

		switch ad.Category {
		case constants.Automobile:
			if ad.Automobile == nil {
				log.Println("automobile cannot be nil")
				c.JSON(http.StatusBadRequest, gin.H{"message": "failed to update ad: automobile cannot be nil"})
				return
			}

			// Update car
			err = helpers.UpdateAutomobile(adId, ad.Automobile, db)
			if err != nil {
				log.Println("failed to update automobile:", err.Error())
				txn.Rollback()
				c.JSON(http.StatusBadRequest, gin.H{"message": "failed to update automobile: " + err.Error()})
				return
			}

		case constants.RealEstate:
			if ad.RealEstate == nil {
				log.Println("real estate cannot be nil")
				c.JSON(http.StatusBadRequest, gin.H{"message": "failed to update ad: real estate cannot be nil"})
				return
			}
		case constants.Job:
			if ad.Job == nil {
				log.Println("job cannot be nil")
				c.JSON(http.StatusBadRequest, gin.H{"message": "failed to update ad: job cannot be nil"})
				return
			}
		default:
			log.Println("wrong category provided")
			c.JSON(http.StatusBadRequest, gin.H{"message": "failed to update ad: wrong category provided"})
			return
		}

		res, err := db.Exec("UPDATE ad SET title = ?, content = ?, category = ?, updated_at = ? WHERE ad_id = ?", ad.Title, ad.Content, ad.Category, time.Now(), adId)
		if err != nil {
			log.Println("error failed to update ad in db:", err)
			txn.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"message": "error failed to update ad in db"})
			return
		}

		row, _ := res.RowsAffected()
		if row <= 0 {
			log.Println("failed to update ad: not found :", err)
			txn.Rollback()
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

		// Delete automobile
		res, err = db.Exec("DELETE FROM automobile WHERE ad_id = ?", adId)
		if err != nil {
			log.Println("error failed to delete automobile from ad in db:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "error failed to delete automobile from ad in db"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "deleted"})
	}
}
