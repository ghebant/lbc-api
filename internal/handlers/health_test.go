package handlers

import (
	"ghebant/lbc-api/internal/constants"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealth(t *testing.T) {
	router := SetupRouter(nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, constants.HealthPath, nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
