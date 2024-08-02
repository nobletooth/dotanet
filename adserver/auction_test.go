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
		{
			name: "Only new ads available",
			allAds: []common.AdWithMetrics{
				{AdInfo: common.AdInfo{Id: 1, Title: "New Ad 1", Price: 1}, ImpressionCount: 1},
				{AdInfo: common.AdInfo{Id: 2, Title: "New Ad 2", Price: 1.5}, ImpressionCount: 2},
				{AdInfo: common.AdInfo{Id: 3, Title: "New Ad 3", Price: 2}, ImpressionCount: 1},
			},
			randFloat:     0.1,
			randInt:       1,
			expectedCode:  http.StatusOK,
			expectedTitle: "New Ad 2",
		},
		{
			name: "Only experienced ads available",
			allAds: []common.AdWithMetrics{
				{AdInfo: common.AdInfo{Id: 4, Title: "Experienced Ad 1", Price: 1}, ImpressionCount: 5, ClickCount: 1},
				{AdInfo: common.AdInfo{Id: 5, Title: "Experienced Ad 2", Price: 2}, ImpressionCount: 10, ClickCount: 3},
				{AdInfo: common.AdInfo{Id: 6, Title: "Experienced Ad 3", Price: 1.5}, ImpressionCount: 7, ClickCount: 2},
			},
			randFloat:     0.9,
			randInt:       2,
			expectedCode:  http.StatusOK,
			expectedTitle: "Experienced Ad 3",
		},
		{
			name: "Both new and experienced ads available, select new",
			allAds: []common.AdWithMetrics{
				{AdInfo: common.AdInfo{Id: 1, Title: "New Ad 1", Price: 1}, ImpressionCount: 1},
				{AdInfo: common.AdInfo{Id: 2, Title: "Experienced Ad 1", Price: 1}, ImpressionCount: 5, ClickCount: 1},
				{AdInfo: common.AdInfo{Id: 3, Title: "Experienced Ad 2", Price: 2}, ImpressionCount: 10, ClickCount: 3},
				{AdInfo: common.AdInfo{Id: 4, Title: "New Ad 2", Price: 1.5}, ImpressionCount: 2},
			},
			randFloat:     0.1,
			randInt:       0,
			expectedCode:  http.StatusOK,
			expectedTitle: "New Ad 1",
		},
		{
			name: "Both new and experienced ads available, select experienced",
			allAds: []common.AdWithMetrics{
				{AdInfo: common.AdInfo{Id: 1, Title: "New Ad 1", Price: 1}, ImpressionCount: 1},
				{AdInfo: common.AdInfo{Id: 2, Title: "Experienced Ad 1", Price: 1}, ImpressionCount: 5, ClickCount: 1},
				{AdInfo: common.AdInfo{Id: 3, Title: "Experienced Ad 2", Price: 2}, ImpressionCount: 10, ClickCount: 3},
				{AdInfo: common.AdInfo{Id: 4, Title: "New Ad 2", Price: 1.5}, ImpressionCount: 2},
			},
			randFloat:     0.25,
			randInt:       1,
			expectedCode:  http.StatusOK,
			expectedTitle: "Experienced Ad 1",
		},
		{
			name: "Higher ctr*Price ad selected",
			allAds: []common.AdWithMetrics{
				{AdInfo: common.AdInfo{Id: 7, Title: "Low CTR Ad 1", Price: 1}, ImpressionCount: 10, ClickCount: 1},
				{AdInfo: common.AdInfo{Id: 8, Title: "High CTR Ad 1", Price: 1}, ImpressionCount: 5, ClickCount: 2},
				{AdInfo: common.AdInfo{Id: 9, Title: "High CTR Ad 2", Price: 1.5}, ImpressionCount: 5, ClickCount: 2},
			},
			randFloat:     0.5,
			randInt:       2,
			expectedCode:  http.StatusOK,
			expectedTitle: "High CTR Ad 2",
		},
		{
			name: "Higher price ad with fewer clicks selected",
			allAds: []common.AdWithMetrics{
				{AdInfo: common.AdInfo{Id: 10, Title: "Low Price High Clicks Ad", Price: 1}, ImpressionCount: 20, ClickCount: 5},
				{AdInfo: common.AdInfo{Id: 11, Title: "High Price Low Clicks Ad", Price: 5}, ImpressionCount: 10, ClickCount: 1},
				{AdInfo: common.AdInfo{Id: 12, Title: "Medium Price Medium Clicks Ad", Price: 1.5}, ImpressionCount: 15, ClickCount: 3},
			},
			randFloat:     0.5,
			randInt:       1,
			expectedCode:  http.StatusOK,
			expectedTitle: "High Price Low Clicks Ad",
		},
	}
}
