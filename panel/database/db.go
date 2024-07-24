package database

import (
	"flag"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

var (
	user      = flag.String("dbuser", "postgres", "Database user")
	password  = flag.String("dbpassword", "Ala.13495782", "Database password")
	dbname    = flag.String("dbname", "todo", "Database name")
	port      = flag.String("dbport", "5432", "Database port")
	host      = flag.String("dbhost", "localhost", "Database host")
	PanelPort = flag.String("panelport", "8081", "Panel port")
)

// var species = flag.String("species", "gopher", "the species we are studying")
func NewDatabase() error {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
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
