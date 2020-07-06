package dbs

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"gitlab.com/quangdangfit/gocommon/utils/logger"
	"goshop/config"
)

var Database *gorm.DB

func init() {
	dbConfig := config.Config.Database
	connectionPath := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=%s",
		dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Name, dbConfig.Password, dbConfig.SSLMode)

	database, err := gorm.Open("postgres", connectionPath)
	if err != nil {
		logger.Fatal("Cannot connect to database", err)
	}

	// Set up connection pool
	database.DB().SetMaxIdleConns(20)
	database.DB().SetMaxOpenConns(200)
	Database = database
}
