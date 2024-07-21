package advertiser

import (
	"flag"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	DB *gorm.DB
}

var (
	user     string
	password string
	dbname   string
	port     string
	host     string
)

func init() {
	flag.StringVar(&user, "dbuser", "postgres", "Database user")
	flag.StringVar(&password, "dbpassword", "Ala.13495782", "Database password")
	flag.StringVar(&dbname, "dbname", "todo", "Database name")
	flag.StringVar(&port, "dbport", "5432", "Database port")
	flag.StringVar(&host, "dbhost", "localhost", "Database host")
}

func NewDatabase() (*Database, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password, dbname, port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return &Database{DB: db}, nil
}

func (d *Database) AutoMigrate(entity interface{}) error {
	return d.DB.AutoMigrate(entity)
}

func (d *Database) Close() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
