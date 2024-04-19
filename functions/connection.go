package funcs

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBConnector func(dbConnStr string) (*gorm.DB, error)

func OpenDatabaseConnection(dbConnStr string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dbConnStr), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func ConnectToDb(
	dbConnStr string,
	dbConnector DBConnector,
	retryDelay time.Duration,
	maxRetries int,
) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	for i := 1; i <= maxRetries; i++ {
		db, err = dbConnector(dbConnStr)
		if err == nil {
			return db, nil
		}

		// Log error and retry
		fmt.Printf("Failed to connect to database (attempt %d): %s\n", i, err)
		if i < maxRetries {
			fmt.Printf("Retrying in %s...\n", retryDelay)
			time.Sleep(retryDelay)
		}
	}

	return nil, fmt.Errorf("failed to connect to database after %d attempts", maxRetries)
}

func GetConnection(viper *viper.Viper) (*string, error) {

	// Read database connection parameters
	databaseHost := viper.GetString("DATABASE_HOST")
	databasePort := viper.GetInt("DATABASE_PORT")
	databaseUsername := viper.GetString("DATABASE_USERNAME")
	databasePassword := viper.GetString("DATABASE_PASSWORD")
	databaseDBName := viper.GetString("DATABASE_DBNAME")
	if databaseDBName == "" ||
		databaseHost == "" ||
		databasePassword == "" ||
		databasePort == 0 ||
		databaseUsername == "" {
		return nil, fmt.Errorf("no database specified")
	}
	// Construct database connection string
	dbConnStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable ",
		databaseHost, fmt.Sprint(databasePort), databaseUsername, databasePassword, databaseDBName)
	fmt.Println("dbConnStr", dbConnStr)
	return &dbConnStr, nil
}
