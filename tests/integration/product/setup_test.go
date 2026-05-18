//go:build integration

package tests_product

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"

	orderModel "goshop/internal/order/model"
	productModel "goshop/internal/product/model"
	userDomain "goshop/internal/user/domain"
	userModel "goshop/internal/user/model"
	"goshop/pkg/dbs"
	"goshop/pkg/redis"
	"goshop/pkg/utils"
	"goshop/tests/testutil"
)

var (
	testRouter *gin.Engine
	dbTest     dbs.Database
	testCache  redis.Redis
)

func TestMain(m *testing.M) {
	env, err := testutil.NewHTTPEnv(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "integration setup failed: %v\n", err)
		os.Exit(1)
	}
	testRouter = env.Engine
	dbTest = env.DB
	testCache = env.Cache

	_ = dbTest.Create(context.Background(), &userModel.User{
		Email:    "test@test.com",
		Password: "test123456",
	})
	_ = dbTest.Create(context.Background(), &userModel.User{
		Email:    "admin@test.com",
		Password: "admin123456",
		Role:     userModel.UserRoleAdmin,
	})

	code := m.Run()
	env.Cleanup()
	os.Exit(code)
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

func loginAs(email, password string) string {
	writer := makeRequest("POST", "/api/v1/auth/login", userDomain.LoginReq{Email: email, Password: password}, "")
	var response map[string]map[string]string
	_ = json.Unmarshal(writer.Body.Bytes(), &response)
	return response["result"]["access_token"]
}

func accessToken() string      { return loginAs("test@test.com", "test123456") }
func adminAccessToken() string { return loginAs("admin@test.com", "admin123456") }

func parseResponseResult(resData []byte, result interface{}) {
	var response map[string]interface{}
	_ = json.Unmarshal(resData, &response)
	_ = utils.Copy(result, response["result"])
}

func cleanData() {
	dbTest.GetDB().Where("1 = 1").Delete(&orderModel.OrderLine{})
	dbTest.GetDB().Where("1 = 1").Delete(&productModel.Product{})
	dbTest.GetDB().Where("1 = 1").Delete(&orderModel.Order{})

	_ = testCache.RemovePattern("*")
}
