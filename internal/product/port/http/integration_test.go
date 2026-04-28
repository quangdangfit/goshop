package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/quangdangfit/gocommon/logger"
	"github.com/quangdangfit/gocommon/validation"

	orderModel "goshop/internal/order/model"
	productModel "goshop/internal/product/model"
	httpServer "goshop/internal/server/http"
	domain "goshop/internal/user/domain"
	userModel "goshop/internal/user/model"
	"goshop/pkg/config"
	"goshop/pkg/dbs"
	"goshop/pkg/redis"
	"goshop/pkg/utils"
)

var (
	testRouter       *gin.Engine
	dbTest           dbs.Database
	testCache        redis.Redis
	integrationReady bool
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	integrationReady = setup()
	exitCode := m.Run()
	if integrationReady {
		teardown()
	}
	os.Exit(exitCode)
}

func setup() bool {
	cfg := config.LoadConfig()
	logger.Initialize(config.ProductionEnv)

	var err error
	dbTest, err = dbs.NewDatabase(cfg.DatabaseURI)
	if err != nil {
		return false
	}

	err = dbTest.AutoMigrate(
		&userModel.User{}, &userModel.Address{}, &userModel.Wishlist{},
		&productModel.Category{}, &productModel.Product{}, &productModel.Review{},
		orderModel.Coupon{}, orderModel.Order{}, orderModel.OrderLine{},
	)
	if err != nil {
		return false
	}

	conn, dialErr := net.DialTimeout("tcp", cfg.RedisURI, 2*time.Second)
	if dialErr != nil {
		return false
	}
	_ = conn.Close()

	validator := validation.New()
	testCache = redis.New(redis.Config{
		Address:  cfg.RedisURI,
		Password: cfg.RedisPassword,
		Database: cfg.RedisDB,
	})

	server := httpServer.NewServer(validator, dbTest, testCache)
	_ = server.MapRoutes()
	testRouter = server.GetEngine()

	_ = dbTest.Create(context.Background(), &userModel.User{
		Email:    "test@test.com",
		Password: "test123456",
	})

	_ = dbTest.Create(context.Background(), &userModel.User{
		Email:    "admin@test.com",
		Password: "admin123456",
		Role:     userModel.UserRoleAdmin,
	})

	return true
}

func teardown() {
	migrator := dbTest.GetDB().Migrator()
	_ = migrator.DropTable(
		&userModel.User{}, &userModel.Address{}, &userModel.Wishlist{},
		&productModel.Category{}, &productModel.Product{}, &productModel.Review{},
		&orderModel.Coupon{}, &orderModel.Order{}, &orderModel.OrderLine{},
	)
}

func requireIntegration(t *testing.T) {
	t.Helper()
	if !integrationReady {
		t.Skip("integration test skipped: DB/Redis not available")
	}
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
	user := domain.LoginReq{
		Email:    "test@test.com",
		Password: "test123456",
	}

	writer := makeRequest("POST", "/api/v1/auth/login", user, "")
	var response map[string]map[string]string
	_ = json.Unmarshal(writer.Body.Bytes(), &response)
	return response["result"]["access_token"]
}

func parseResponseResult(resData []byte, result interface{}) {
	var response map[string]interface{}
	_ = json.Unmarshal(resData, &response)
	_ = utils.Copy(result, response["result"])
}

func adminAccessToken() string {
	user := domain.LoginReq{
		Email:    "admin@test.com",
		Password: "admin123456",
	}

	writer := makeRequest("POST", "/api/v1/auth/login", user, "")
	var response map[string]map[string]string
	_ = json.Unmarshal(writer.Body.Bytes(), &response)
	return response["result"]["access_token"]
}

func cleanData() {
	dbTest.GetDB().Where("1 = 1").Delete(&orderModel.OrderLine{})
	dbTest.GetDB().Where("1 = 1").Delete(&productModel.Product{})
	dbTest.GetDB().Where("1 = 1").Delete(&orderModel.Order{})

	_ = testCache.RemovePattern("*")
}
