package main

import (
	"flag"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

var (
	user             = flag.String("dbuser", "postgres", "Database user")
	password         = flag.String("dbpassword", "Ala.13495782", "Database password")
	dbname           = flag.String("dbname", "todo", "Database name")
	port             = flag.String("dbport", "5432", "Database port")
	host             = flag.String("dbhost", "localhost", "Database host")
	EventservicePort = flag.String("panelport", "8082", "Panel port")
)

func OpenDbConnection() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password, dbname, port)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	return DB, nil
}
