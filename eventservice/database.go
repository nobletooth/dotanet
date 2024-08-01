package main

import (
	"flag"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

var (
	user            = flag.String("dbuser", "user", "Database user")
	password        = flag.String("dbpassword", "password", "Database password")
	dbname          = flag.String("dbname", "dotanet", "Database name")
	port            = flag.String("dbport", "5432", "Database port")
	host            = flag.String("dbhost", "95.217.125.139", "Database host")
	EventserviceUrl = flag.String("eventserviceurl", ":8082", "Panel port")
	kafkaendpoint   = flag.String("kafkaendpoint", "localhost:9092", "kafka end point")
)

func OpenDbConnection() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		*host, *user, *password, *dbname, *port)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	return DB, nil
}
