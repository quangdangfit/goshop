package main

import (
	"gitlab.com/quangdangfit/gocommon/utils/logger"
	"goshop/dbs"
)

func main() {
	ssl := dbs.Database
	logger.Info("Main", ssl)
}
