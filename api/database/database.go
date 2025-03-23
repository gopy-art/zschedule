package database

import (
	"fmt"
	"os"
	"zschedule/cmd"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DatabaseConfiguration struct {
	Address  string `json:"address"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	DBname   string `json:"database_name"`
	DB       *gorm.DB
}

func NewDatabaseConnection() (*DatabaseConfiguration, error) {
	// Load environment variables
	err := godotenv.Load(cmd.EnvFile)
	if err != nil {
		return nil, err
	}

	instance := new(DatabaseConfiguration)
	if os.Getenv("DATABASE_ADDRESS") == "" {
		return nil, fmt.Errorf("DATABASE_ADDRESS is empty")
	} else {
		instance.Address = os.Getenv("DATABASE_ADDRESS")
	}
	if os.Getenv("DATABASE_PORT") == "" {
		return nil, fmt.Errorf("DATABASE_PORT is empty")
	} else {
		instance.Port = os.Getenv("DATABASE_PORT")
	}
	if os.Getenv("DATABASE_USERNAME") == "" {
		return nil, fmt.Errorf("DATABASE_USERNAME is empty")
	} else {
		instance.Username = os.Getenv("DATABASE_USERNAME")
	}
	if os.Getenv("DATABASE_PASSWORD") == "" {
		return nil, fmt.Errorf("DATABASE_PASSWORD is empty")
	} else {
		instance.Password = os.Getenv("DATABASE_PASSWORD")
	}
	if os.Getenv("DATABASE_NAME") == "" {
		return nil, fmt.Errorf("DATABASE_NAME is empty")
	} else {
		instance.DBname = os.Getenv("DATABASE_NAME")
	}

	return instance, nil
}

func (d *DatabaseConfiguration) Connection() error {
	// Construct connection string
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		d.Address,
		d.Username,
		d.Password,
		d.DBname,
		d.Port,
	)

	// Connect to the database
	var err error
	if d.DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{}); err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	return nil
}

func (d *DatabaseConfiguration) CreateTables(tables ...any) error {
	return d.DB.AutoMigrate(tables...)
}