package database

import (
	"flag"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

var (
	user      = flag.String("dbuser", "1", "Database user")
	password  = flag.String("dbpassword", "2", "Database password")
	dbname    = flag.String("dbname", "3", "Database name")
	port      = flag.String("dbport", "4", "Database port")
	host      = flag.String("dbhost", "5", "Database host")
	PanelUrl = flag.String("panelurl", "localhost:8081", "Panel url")
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
