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
	"go.uber.org/dig"

	"goshop/app/dbs"
	"goshop/app/models"
	"goshop/app/repositories"
	"goshop/app/serializers"
	"goshop/app/services"
	"goshop/config"
	"goshop/pkg/utils"
)

var (
	testContainer *dig.Container
	testRouter    *gin.Engine
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

	testContainer = buildContainer()
	testRouter = initGinEngine(testContainer)

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

func buildContainer() *dig.Container {
	container := dig.New()

	// Inject repositories
	repositories.Inject(container)

	// Inject services
	services.Inject(container)

	// Inject APIs
	Inject(container)

	return container
}
func initGinEngine(container *dig.Container) *gin.Engine {
	app := gin.Default()
	RegisterAPI(app, container)
	return app
}
