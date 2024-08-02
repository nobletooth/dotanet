package main

import (
	"github.com/gin-gonic/gin"
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
