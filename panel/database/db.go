package database

import (
	"flag"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

var (
	user     = flag.String("dbuser", "user", "Database user")
	password = flag.String("dbpassword", "password", "Database password")
	dbname   = flag.String("dbname", "dotanet", "Database name")
	port     = flag.String("dbport", "5432", "Database port")
	host     = flag.String("dbhost", "95.217.125.139", "Database host") //95.217.125.139
	PanelUrl = flag.String("panelurl", ":8085", "Panel port")
)

func NewDatabase() error {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s",
		*host, *user, *password, *dbname, *port)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	return nil
}

func Close() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
