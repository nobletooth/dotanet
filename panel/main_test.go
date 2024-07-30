package main

import (
	"common"
	"github.com/gin-gonic/gin"
	"github.com/nobletooth/dotanet/panel/advertiser"
	"github.com/nobletooth/dotanet/panel/database"
	"github.com/nobletooth/dotanet/panel/publisher"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func setupTestDatabase() {
	err := database.NewDatabase()
	if err != nil {
		panic("failed to connect to test database")
	}

	err = database.DB.AutoMigrate(&publisher.Publisher{})
	if err != nil {
		panic("failed to migrate test database")
	}
	err = database.DB.AutoMigrate(&advertiser.Ad{})
	if err != nil {
		panic("failed to migrate test database")
	}
	err = database.DB.AutoMigrate(&common.ClickedEvent{}, &common.ViewedEvent{})
	if err != nil {
		panic("failed to migrate test database")
	}
}

func TestEventServiceHandler(t *testing.T) {
	setupTestDatabase()
	defer func() {
		if err := database.Close(); err != nil {
			t.Fatalf("failed to close test database: %v", err)
		}
	}()

	router := gin.Default()
	router.POST("/eventservice", eventServerHandler)

	testEvent := `{"pubId":"1", "adId":"1", "isClicked":true, "time":"2024-07-24T00:00:00Z"}`
	req, _ := http.NewRequest("POST", "/eventservice", strings.NewReader(testEvent))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Event processed successfully")

	var clickedEvent common.ClickedEvent
	result := database.DB.Where("pid = ? AND ad_id = ?", "1", "101").First(&clickedEvent)
	assert.NoError(t, result.Error)
	assert.Equal(t, "1", clickedEvent.Pid)
	assert.Equal(t, "1", clickedEvent.AdId)
	assert.NotZero(t, clickedEvent.Time)
	////////////////////////////////////////////////////////////////////////////////////////
	testEvent2 := `{"pubId":"1", "adId":"1", "isClicked":false, "time":"2024-07-24T00:00:00Z"}`
	req, _ = http.NewRequest("POST", "/eventservice", strings.NewReader(testEvent2))
	req.Header.Set("Content-Type", "application/json")

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Event processed successfully")

	var clickedEvent2 common.ClickedEvent
	result = database.DB.Where("pid = ? AND ad_id = ?", "1", "101").First(&clickedEvent)
	assert.NoError(t, result.Error)
	assert.Equal(t, "1", clickedEvent2.Pid)
	assert.Equal(t, "1", clickedEvent2.AdId)
	assert.NotZero(t, clickedEvent2.Time)
}

func TestHomeHandler(t *testing.T) {
	router := gin.Default()
	router.GET("/", homeHandler)

	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "<html>")
}
