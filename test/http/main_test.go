package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/quangdangfit/gocommon/logger"
	"github.com/quangdangfit/gocommon/validation"

	orderModel "goshop/internal/order/model"
	productModel "goshop/internal/product/model"
	httpServer "goshop/internal/server/http"
	"goshop/internal/user/dto"
	userModel "goshop/internal/user/model"
	"goshop/pkg/config"
	"goshop/pkg/dbs"
	"goshop/pkg/redis"
	"goshop/pkg/utils"
)

var (
	testRouter *gin.Engine
	dbTest     dbs.IDatabase
	testCache  redis.IRedis
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	setup()
	exitCode := m.Run()
	teardown()

	os.Exit(exitCode)
}

func setup() {
	cfg := config.LoadConfig()
	logger.Initialize(config.ProductionEnv)

	var err error
	dbTest, err = dbs.NewDatabase(cfg.DatabaseURI)
	if err != nil {
		logger.Fatal("Cannot connect to database", err)
	}

	err = dbTest.AutoMigrate(&userModel.User{}, &productModel.Product{}, orderModel.Order{}, orderModel.OrderLine{})
	if err != nil {
		logger.Fatal("Database migration fail", err)
	}

	validator := validation.New()
	testCache = redis.New(redis.Config{
		Address:  cfg.RedisURI,
		Password: cfg.RedisPassword,
		Database: cfg.RedisDB,
	})

	server := httpServer.NewServer(validator, dbTest, testCache)
	_ = server.MapRoutes()
	testRouter = server.GetEngine()

	dbTest.Create(context.Background(), &userModel.User{
		Email:    "test@test.com",
		Password: "test123456",
	})
}

func teardown() {
	migrator := dbTest.GetDB().Migrator()
	migrator.DropTable(&userModel.User{}, &productModel.Product{}, &orderModel.Order{}, &orderModel.OrderLine{})
}

func makeRequest(method, url string, body interface{}, token string) *httptest.ResponseRecorder {
	requestBody, _ := json.Marshal(body)
	request, _ := http.NewRequest(method, url, bytes.NewBuffer(requestBody))
	if token != "" {
		request.Header.Add("Authorization", "Bearer "+token)
	}
	writer := httptest.NewRecorder()
	testRouter.ServeHTTP(writer, request)
	return writer
}

func accessToken() string {
	user := dto.LoginReq{
		Email:    "test@test.com",
		Password: "test123456",
	}

	writer := makeRequest("POST", "/api/v1/auth/login", user, "")
	var response map[string]map[string]string
	_ = json.Unmarshal(writer.Body.Bytes(), &response)
	return response["result"]["access_token"]
}

func refreshToken() string {
	user := dto.LoginReq{
		Email:    "test@test.com",
		Password: "test123456",
	}

	writer := makeRequest("POST", "/api/v1/auth/login", user, "")
	var response map[string]map[string]string
	_ = json.Unmarshal(writer.Body.Bytes(), &response)
	return response["result"]["refresh_token"]
}

func parseResponseResult(resData []byte, result interface{}) {
	var response map[string]interface{}
	_ = json.Unmarshal(resData, &response)
	utils.Copy(result, response["result"])
}

func cleanData(records ...interface{}) {
	dbTest.GetDB().Where("1 = 1").Delete(&orderModel.OrderLine{})
	dbTest.GetDB().Where("1 = 1").Delete(&productModel.Product{})
	dbTest.GetDB().Where("1 = 1").Delete(&orderModel.Order{})

	for _, record := range records {
		dbTest.Delete(context.Background(), record)
	}

	testCache.RemovePattern("*")
}
