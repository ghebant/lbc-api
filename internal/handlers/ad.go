package handlers

import (
	"context"
	"database/sql"
	"ghebant/lbc-api/internal/constants"
	"ghebant/lbc-api/models"
	"github.com/gin-gonic/gin"
	"gopkg.in/errgo.v2/errors"
	"log"
	"net/http"
	"strconv"
	"time"
)

// Todo put in helper
func FindRow(db *sql.DB, idStr, key string) (*sql.Row, int, error) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return nil, http.StatusBadRequest, errors.New("failed to convert id: " + err.Error())
	}

	query := "SELECT * FROM ad WHERE " + key + " = ?"
	res := db.QueryRow(query, id)
	if err != nil {
		return nil, http.StatusNotFound, errors.New("no ad found for the provided id: " + err.Error())
	}

	return res, http.StatusOK, nil
}

func GetAd(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		ads := []models.Ad{}

		// Find ad by ID
		adId := c.Param("id")
		if adId != "" {
			var ad models.Ad

			row, status, err := FindRow(db, adId, constants.AdPrimaryKey)
			if err != nil {
				log.Println(err)
				c.JSON(status, gin.H{"message": err.Error()})
				return
			}

			err = row.Scan(&ad.ID, &ad.Title, &ad.Content, &ad.Category, &ad.CreatedAt, &ad.UpdatedAt)
			if err != nil {
				log.Println("no ad found for the provided id: " + err.Error())
				c.JSON(http.StatusNotFound, gin.H{"message": "no ad found for the provided id"})
				return
			}

			c.JSON(http.StatusOK, ad)
			return
		}

		queryRes, err := db.Query("SELECT * FROM ad")
		if err != nil {
			log.Println("error failed to query db: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "failed to query db"})
			return
		}
		defer queryRes.Close()

		for queryRes.Next() {
			var ad models.Ad

			err = queryRes.Scan(&ad.ID, &ad.Title, &ad.Content, &ad.Category, &ad.CreatedAt, &ad.UpdatedAt)
			if err != nil {
				log.Println("error while scanning ad: " + err.Error())
				c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "failed to scan result"})
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
			return
		}

		insertQuery := "INSERT INTO ad(title, content, category) VALUES (?, ?, ?)"

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		stmt, err := db.PrepareContext(ctx, insertQuery)
		if err != nil {
			log.Printf("Error %s when preparing SQL statement", err)
			return
		}
		defer stmt.Close()
		res, err := stmt.ExecContext(ctx, ad.Title, ad.Content, ad.Category)
		if err != nil {
			log.Printf("Error %s when inserting row into products table", err)
			return
		}
		rows, err := res.RowsAffected()
		if err != nil {
			log.Printf("Error %s when finding rows affected", err)
			return
		}
		log.Printf("%d products created ", rows)

		c.JSON(http.StatusOK, gin.H{"code": 200, "message": "post"})
		//c.IndentedJSON()
	}
}

func UpdateAd(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"code": 200, "message": "edit"})
		//c.IndentedJSON()
	}
}

func DeleteAd(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"code": 200, "message": "delete"})
		//c.IndentedJSON()
	}
}
