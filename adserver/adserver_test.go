package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"common"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetAdsHandlerUnit(t *testing.T) {
	testAds := []common.AdWithMetrics{
		{Id: 1, Title: "Test Ad 1"},
		{Id: 2, Title: "Test Ad 2"},
	}
	allAds = testAds

	// Create a mock Gin context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	GetAdsHandler(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response []common.AdWithMetrics
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, testAds, response)
}
