package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/quangdangfit/gocommon/logger"
	"github.com/quangdangfit/gocommon/redis"
	"github.com/quangdangfit/gocommon/validation"

	"goshop/internal/app/api"
	"goshop/internal/app/dbs"
	"goshop/internal/app/models"
	"goshop/internal/app/repositories"
	"goshop/internal/app/serializers"
	"goshop/internal/app/services"
	"goshop/internal/config"
	"goshop/internal/pkg/utils"
)

var (
	testRouter     *gin.Engine
	mockTestRouter *gin.Engine
	testUserAPI    *api.UserAPI
	testProductAPI *api.ProductAPI
	testOrderAPI   *api.OrderAPI
	testRedis      redis.IRedis
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	setup()
	exitCode := m.Run()
	teardown()

	os.Exit(exitCode)
}

func setup() {
	cfg := config.GetConfig()
	logger.Initialize(config.TestEnv)

	dbs.Init()

	validator := validation.New()
	testRedis = redis.New(redis.Config{
		Address:  cfg.RedisURI,
		Password: cfg.RedisPassword,
		Database: cfg.RedisDB,
	})

	userRepo := repositories.NewUserRepository()
	productRepo := repositories.NewProductRepository()
	orderRepo := repositories.NewOrderRepository()

	userSvc := services.NewUserService(userRepo)
	productSvc := services.NewProductService(productRepo)
	orderSvc := services.NewOrderService(orderRepo, productRepo)

	testUserAPI = api.NewUserAPI(validator, userSvc)
	testProductAPI = api.NewProductAPI(validator, testRedis, productSvc)
	testOrderAPI = api.NewOrderAPI(validator, orderSvc)

	testRouter = initGinEngine(testUserAPI, testProductAPI, testOrderAPI)

	dbs.Database.Create(&models.User{
		Email:    "test@test.com",
		Password: "test123456",
	})
}

func teardown() {
	migrator := dbs.Database.Migrator()
	migrator.DropTable(&models.User{}, &models.Product{}, &models.Order{}, &models.OrderLine{})
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
	user := serializers.LoginReq{
		Email:    "test@test.com",
		Password: "test123456",
	}

	writer := makeRequest("POST", "/auth/login", user, "")
	var response map[string]map[string]string
	_ = json.Unmarshal(writer.Body.Bytes(), &response)
	return response["result"]["access_token"]
}

func refreshToken() string {
	user := serializers.LoginReq{
		Email:    "test@test.com",
		Password: "test123456",
	}

	writer := makeRequest("POST", "/auth/login", user, "")
	var response map[string]map[string]string
	_ = json.Unmarshal(writer.Body.Bytes(), &response)
	return response["result"]["refresh_token"]
}

func parseResponseResult(resData []byte, result interface{}) {
	var response map[string]interface{}
	_ = json.Unmarshal(resData, &response)
	utils.Copy(result, response["result"])
}

func initGinEngine(
	userAPI *api.UserAPI,
	productAPI *api.ProductAPI,
	orderAPI *api.OrderAPI,
) *gin.Engine {
	cfg := config.GetConfig()
	if cfg.Environment == config.ProductionEnv {
		gin.SetMode(gin.ReleaseMode)
	}
	app := gin.Default()
	api.RegisterAPI(app, userAPI, productAPI, orderAPI)
	return app
}

func makeMockRequest(method, url string, body interface{}, token string) *httptest.ResponseRecorder {
	requestBody, _ := json.Marshal(body)
	request, _ := http.NewRequest(method, url, bytes.NewBuffer(requestBody))
	if token != "" {
		request.Header.Add("Authorization", "Bearer "+token)
	}
	writer := httptest.NewRecorder()
	mockTestRouter.ServeHTTP(writer, request)
	return writer
}

func cleanData(records ...interface{}) {
	dbs.Database.Where("1 = 1").Delete(&models.OrderLine{})
	dbs.Database.Where("1 = 1").Delete(&models.Product{})
	dbs.Database.Where("1 = 1").Delete(&models.Order{})

	for _, record := range records {
		dbs.Database.Delete(record)
	}

	testRedis.RemovePattern("*")
}
