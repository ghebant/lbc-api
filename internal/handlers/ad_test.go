package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"ghebant/lbc-api/internal/constants"
	"ghebant/lbc-api/internal/helpers"
	"ghebant/lbc-api/models"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetAd(t *testing.T) {
	db, err := helpers.SetupDbForTesting()
	defer db.Close()
	assert.NoError(t, err)

	mocks := []models.Ad{
		{
			Title:    "testAd1",
			Content:  "testContent1",
			Category: "testCategory1",
		},
		{
			Title:    "testAd2",
			Content:  "testContent2",
			Category: "testCategory2",
		},
		{
			Title:    "testAd3",
			Content:  "testContent3",
			Category: "testCategory3",
		},
	}

	err = helpers.InsertAds(mocks, db)
	assert.NoError(t, err, "failed to insert mock ads")

	router := SetupRouter(db)

	tests := []struct {
		testTitle      string
		route          string
		expectedStatus int
		expectResult   bool
		expected       []models.Ad
	}{
		{"No ads", constants.AdPath, http.StatusOK, false, []models.Ad{}},
		{"Get all ads", constants.AdPath, http.StatusOK, true, mocks},
		{"Get ad with valid id", fmt.Sprintf("%s/%d", constants.AdPath, mocks[1].ID), http.StatusOK, true, []models.Ad{mocks[1]}},
		{"Get ad with bad id", constants.AdPath + "/90", http.StatusNotFound, false, []models.Ad{}},
	}

	for i := range tests {
		// Call to GET ad handler
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, tests[i].route, nil)
		router.ServeHTTP(w, req)

		assert.Equalf(t, tests[i].expectedStatus, w.Code, "test %s failed", tests[i].testTitle)

		if tests[i].expectResult {
			var got []models.Ad

			// If get ad with id, we expect only one ad in response
			if strings.Contains(tests[i].route, constants.AdPath+"/") {
				var ad models.Ad
				err = json.NewDecoder(w.Body).Decode(&ad)
				got = append(got, ad)
			} else {
				err = json.NewDecoder(w.Body).Decode(&got)
			}

			assert.NoError(t, err, "failed to read body from request")
			assert.Equalf(t, len(tests[i].expected), len(got), "test %s failed", tests[i].testTitle)

			for j := range tests[i].expected {
				assert.Equalf(t, tests[i].expected[j].Title, got[j].Title, "test %s failed", tests[i].testTitle)
				assert.Equalf(t, tests[i].expected[j].Content, got[j].Content, "test %s failed", tests[i].testTitle)
				assert.Equalf(t, tests[i].expected[j].Category, got[j].Category, "test %s failed", tests[i].testTitle)
				assert.Falsef(t, got[j].CreatedAt.IsZero(), "test %s failed", tests[i].testTitle)
				assert.Falsef(t, got[j].UpdatedAt.IsZero(), "test %s failed", tests[i].testTitle)
			}
		}
	}
}

func TestPostAd(t *testing.T) {
	db, err := helpers.SetupDbForTesting()
	defer db.Close()
	assert.NoError(t, err)

	router := SetupRouter(db)

	tests := []struct {
		testTitle      string
		route          string
		expectedStatus int
		expectResult   bool
		expected       models.Ad
	}{
		{"Post valid ad", constants.AdPath, http.StatusCreated, true, models.Ad{
			Title:    "testAd_post1",
			Content:  "testContent_post1",
			Category: "testCategory_post1",
		}},
		{"Post ad title is missing", constants.AdPath, http.StatusBadRequest, false, models.Ad{
			Title:    "",
			Content:  "testContent_post1",
			Category: "category",
		}},
		{"Post ad content is missing", constants.AdPath, http.StatusBadRequest, false, models.Ad{
			Title:    "testAd_post1",
			Content:  "",
			Category: "category",
		}},
		{"Post ad category is missing", constants.AdPath, http.StatusBadRequest, false, models.Ad{
			Title:    "testAd_post1",
			Content:  "testContent_post1",
			Category: "",
		}},
	}

	for i := range tests {
		// Call to POST ad handler
		w := httptest.NewRecorder()

		body, err := json.Marshal(tests[i].expected)
		assert.NoError(t, err, "failed to marshall expected ad")

		req, _ := http.NewRequest(http.MethodPost, tests[i].route, bytes.NewBuffer(body))
		router.ServeHTTP(w, req)

		assert.Equalf(t, tests[i].expectedStatus, w.Code, "test %s failed", tests[i].testTitle)

		if tests[i].expectResult {
			var got models.Ad

			err = json.NewDecoder(w.Body).Decode(&got)
			assert.NoError(t, err, "failed to read body from request")

			assert.Equalf(t, tests[i].expected.Title, got.Title, "test %s failed", tests[i].testTitle)
			assert.Equalf(t, tests[i].expected.Content, got.Content, "test %s failed", tests[i].testTitle)
			assert.Equalf(t, tests[i].expected.Category, got.Category, "test %s failed", tests[i].testTitle)
			assert.Falsef(t, got.CreatedAt.IsZero(), "test %s failed", tests[i].testTitle)
			assert.Falsef(t, got.UpdatedAt.IsZero(), "test %s failed", tests[i].testTitle)
		}
	}
}

func TestUpdateAd(t *testing.T) {
	db, err := helpers.SetupDbForTesting()
	defer db.Close()
	assert.NoError(t, err)

	mocks := []models.Ad{
		{
			Title:    "testAd1",
			Content:  "testContent1",
			Category: "testCategory1",
		},
		{
			Title:    "testAd2",
			Content:  "testContent2",
			Category: "testCategory2",
		},
		{
			Title:    "testAd3",
			Content:  "testContent3",
			Category: "testCategory3",
		},
	}

	err = helpers.InsertAds(mocks, db)
	assert.NoError(t, err, "failed to insert mock ads")

	router := SetupRouter(db)

	tests := []struct {
		testTitle      string
		route          string
		expectedStatus int
		expectResult   bool
		expected       models.Ad
	}{
		{"Put valid ad", fmt.Sprintf("%s/%d", constants.AdPath, mocks[0].ID), http.StatusOK, true, models.Ad{
			Title:    "test put",
			Content:  "testContent_post1",
			Category: "testCategory_post1",
		}},
		{"Put ad title is missing", fmt.Sprintf("%s/%d", constants.AdPath, mocks[0].ID), http.StatusBadRequest, false, models.Ad{
			Title:    "",
			Content:  "testContent_post1",
			Category: "category",
		}},
		{"Put ad content is missing", fmt.Sprintf("%s/%d", constants.AdPath, mocks[0].ID), http.StatusBadRequest, false, models.Ad{
			Title:    "testAd_post1",
			Content:  "",
			Category: "category",
		}},
		{"Put ad category is missing", fmt.Sprintf("%s/%d", constants.AdPath, mocks[0].ID), http.StatusBadRequest, false, models.Ad{
			Title:    "testAd_post1",
			Content:  "testContent_post1",
			Category: "",
		}},
		{"Put ad id not exist", fmt.Sprintf("%s/%d", constants.AdPath, 111), http.StatusNotFound, false, models.Ad{
			Title:    "test",
			Content:  "test",
			Category: "test",
		}},
		{"Put ad invalid id", constants.AdPath + "/a", http.StatusBadRequest, false, models.Ad{
			Title:    "test",
			Content:  "test",
			Category: "test",
		}},
	}

	for i := range tests {
		// Call to POST ad handler
		w := httptest.NewRecorder()

		body, err := json.Marshal(tests[i].expected)
		assert.NoError(t, err, "failed to marshall expected ad")

		req, _ := http.NewRequest(http.MethodPut, tests[i].route, bytes.NewBuffer(body))
		router.ServeHTTP(w, req)

		assert.Equalf(t, tests[i].expectedStatus, w.Code, "test %s failed", tests[i].testTitle)

		if tests[i].expectResult {
			var got models.Ad

			err = json.NewDecoder(w.Body).Decode(&got)
			assert.NoError(t, err, "failed to read body from request")

			assert.Equalf(t, tests[i].expected.Title, got.Title, "test %s failed", tests[i].testTitle)
			assert.Equalf(t, tests[i].expected.Content, got.Content, "test %s failed", tests[i].testTitle)
			assert.Equalf(t, tests[i].expected.Category, got.Category, "test %s failed", tests[i].testTitle)
			assert.Falsef(t, got.CreatedAt.IsZero(), "test %s failed", tests[i].testTitle)
			assert.Truef(t, got.UpdatedAt.After(got.CreatedAt), "test %s failed", tests[i].testTitle)
		}
	}
}

func TestDeleteAd(t *testing.T) {
	db, err := helpers.SetupDbForTesting()
	defer db.Close()
	assert.NoError(t, err)

	mocks := []models.Ad{
		{
			Title:    "testAd1",
			Content:  "testContent1",
			Category: "testCategory1",
		},
		{
			Title:    "testAd2",
			Content:  "testContent2",
			Category: "testCategory2",
		},
		{
			Title:    "testAd3",
			Content:  "testContent3",
			Category: "testCategory3",
		},
	}

	err = helpers.InsertAds(mocks, db)
	assert.NoError(t, err, "failed to insert mock ads")

	router := SetupRouter(db)

	tests := []struct {
		testTitle       string
		route           string
		expectedStatus  int
		shouldBeDeleted bool
		adId            int
	}{
		{"Delete valid ad", fmt.Sprintf("%s/%d", constants.AdPath, mocks[0].ID), http.StatusOK, true, mocks[0].ID},
		{"Delete ad invalid id", fmt.Sprintf("%s/%d", constants.AdPath, 111), http.StatusNotFound, false, 0},
	}

	for i := range tests {
		// Call to POST ad handler
		w := httptest.NewRecorder()

		req, _ := http.NewRequest(http.MethodDelete, tests[i].route, nil)
		router.ServeHTTP(w, req)

		assert.Equalf(t, tests[i].expectedStatus, w.Code, "test %s failed", tests[i].testTitle)

		if tests[i].shouldBeDeleted {
			_, err := helpers.FindAdById(db, tests[i].adId)
			assert.Equal(t, sql.ErrNoRows, err)
		}
	}
}
