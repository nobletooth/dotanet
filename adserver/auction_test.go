package main

import (
	"common"
	"github.com/gin-gonic/gin"
	"net/http"
	"testing"
)

// Mock Random Generator
type mockRandomGenerator struct {
	float64Result float64
	intnResult    int
}

func (m *mockRandomGenerator) Float64() float64 {
	return m.float64Result
}

func (m *mockRandomGenerator) Intn(n int) int {
	return m.intnResult
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/ad/:pubID", GetAdHandler)
	return r
}

func TestGetAdHandler(t *testing.T) {
	tests := []struct {
		name          string
		allAds        []common.AdWithMetrics
		randFloat     float64
		randInt       int
		expectedCode  int
		expectedTitle string
	}{
		{
			name:          "No ads available",
			allAds:        []common.AdWithMetrics{},
			expectedCode:  http.StatusNotFound,
			expectedTitle: "",
		},
	}
}
