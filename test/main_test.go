package test

import (
	"fmt"
	"os"
	"testing"

	"github.com/quangdangfit/gocommon/logger"
	"go.uber.org/dig"

	"goshop/app"
)

var (
	container *dig.Container
)

func TestMain(m *testing.M) {
	logger.Initialize("test")
	container = app.BuildContainer()

	setup()
	code := m.Run()
	teardown()

	os.Exit(code)
}

func setup() {
	fmt.Println("============> Setup for testing")
	clearAll()
	createData()
}

func teardown() {
	fmt.Println("============> Teardown")
	clearAll()
}

func createData() {
}

func clearAll() {
}
