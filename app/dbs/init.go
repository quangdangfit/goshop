package dbs

import (
	"github.com/quangdangfit/gocommon/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"

	"goshop/config"
)

var Database *gorm.DB

func Init() {
	cfg := config.GetConfig()
	database, err := gorm.Open(postgres.Open(cfg.DatabaseURI), &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogger.Warn),
	})
	if err != nil {
		logger.Fatal("Cannot connect to database", err)
	}

	// Set up connection pool
	sqlDB, err := database.DB()
	if err != nil {
		logger.Fatal("Cannot connect to database", err)
	}
	sqlDB.SetMaxIdleConns(20)
	sqlDB.SetMaxOpenConns(200)
	Database = database
}
