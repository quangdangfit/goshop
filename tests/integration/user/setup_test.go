//go:build integration

package tests_user

import (
	"bytes"
	"context"
	cryptoRand "crypto/rand"
	"encoding/hex"
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

func loginTokens() map[string]string {
	user := userDomain.LoginReq{
		Email:    "test@test.com",
		Password: "test123456",
	}
	writer := makeRequest("POST", "/api/v1/auth/login", user, "")
	var response map[string]map[string]string
	_ = json.Unmarshal(writer.Body.Bytes(), &response)
	return response["result"]
}

func accessToken() string  { return loginTokens()["access_token"] }
func refreshToken() string { return loginTokens()["refresh_token"] }

// randSuffix produces a short unique suffix so tests that create users /
// coupons in parallel don't collide on unique indexes. Crypto-grade random
// is overkill here — time-based monotonic IDs would do — but rand keeps the
// helper self-contained.
func randSuffix() string {
	var b [6]byte
	_, _ = cryptoRand.Read(b[:])
	return hex.EncodeToString(b[:])
}

func cleanData(records ...interface{}) {
	dbTest.GetDB().Where("1 = 1").Delete(&orderModel.OrderLine{})
	dbTest.GetDB().Where("1 = 1").Delete(&productModel.Product{})
	dbTest.GetDB().Where("1 = 1").Delete(&orderModel.Order{})

	for _, record := range records {
		_ = dbTest.Delete(context.Background(), record)
	}

	_ = testCache.RemovePattern("*")
}
