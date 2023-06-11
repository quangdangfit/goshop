package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/quangdangfit/gocommon/logger"
	"github.com/quangdangfit/gocommon/validation"

	"goshop/app/dbs"
	"goshop/app/models"
	"goshop/app/repositories"
	"goshop/app/serializers"
	"goshop/app/services"
	"goshop/config"
	"goshop/pkg/utils"
)

var (
	testRouter *gin.Engine
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	setup()
	exitCode := m.Run()
	teardown()

	os.Exit(exitCode)
}

func setup() {
	logger.Initialize(config.TestEnv)

	dbs.Init()

	validator := validation.New()

	userRepo := repositories.NewUserRepository()
	productRepo := repositories.NewProductRepository()
	orderRepo := repositories.NewOrderRepository()

	userSvc := services.NewUserService(userRepo)
	productSvc := services.NewProductService(productRepo)
	orderSvc := services.NewOrderService(orderRepo, productRepo)

	userAPI := NewUserAPI(validator, userSvc)
	productAPI := NewProductAPI(validator, productSvc)
	orderAPI := NewOrderAPI(validator, orderSvc)

	testRouter = initGinEngine(userAPI, productAPI, orderAPI)

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

func initGinEngine(userAPI *UserAPI,
	productAPI *ProductAPI,
	orderAPI *OrderAPI,
) *gin.Engine {
	cfg := config.GetConfig()
	if cfg.Environment == config.ProductionEnv {
		gin.SetMode(gin.ReleaseMode)
	}
	app := gin.Default()
	RegisterAPI(app, userAPI, productAPI, orderAPI)
	return app
}
