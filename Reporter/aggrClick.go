package main

import "gorm.io/gorm"

type aggrClick struct {
	gorm.Model
	AdId       int
	ClickCount int
}
