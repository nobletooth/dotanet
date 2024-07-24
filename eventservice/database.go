package main

import (
	"flag"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

var (
	user             string
	password         string
	dbname           string
	port             string
	host             string
	EventservicePort string
)

func init() {
	flag.StringVar(&user, "dbuser", "postgres", "Database user")
	flag.StringVar(&password, "dbpassword", "Ala.13495782", "Database password")
	flag.StringVar(&dbname, "dbname", "todo", "Database name")
	flag.StringVar(&port, "dbport", "5432", "Database port")
	flag.StringVar(&host, "dbhost", "localhost", "Database host")
	flag.StringVar(&EventservicePort, "eventserviceport", "8082", "eventservice port")
}

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
