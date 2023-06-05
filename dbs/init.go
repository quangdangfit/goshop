package dbs

import (
	"github.com/jinzhu/gorm"
	"github.com/quangdangfit/gocommon/logger"

	"goshop/config"
)

var Database *gorm.DB

func Init() {
	cfg := config.GetConfig()
	database, err := gorm.Open("postgres", cfg.DatabaseURI)
	if err != nil {
		logger.Fatal("Cannot connect to database", err)
	}

	// Set up connection pool
	database.DB().SetMaxIdleConns(20)
	database.DB().SetMaxOpenConns(200)
	Database = database
}
