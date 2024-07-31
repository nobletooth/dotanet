package main

import (
	"gorm.io/gorm"
)

type aggrImpression struct {
	gorm.Model
	AdId            int
	ImpressionCount int
}
