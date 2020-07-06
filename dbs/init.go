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
	connectionPath := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=%s",
		config.Config.Database.Host, config.Config.Database.Port, config.Config.Database.User,
		config.Config.Database.Name, config.Config.Database.Password, config.Config.Database.SSLMode)

	database, err := gorm.Open("postgres", connectionPath)
	if err != nil {
		logger.Fatal("Cannot connect to database", err)
	}

	// Set up connection pool
	database.DB().SetMaxIdleConns(20)
	database.DB().SetMaxOpenConns(200)
	Database = database
}
