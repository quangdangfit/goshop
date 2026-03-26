package response

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"goshop/pkg/config"
)

func TestJSON(t *testing.T) {
	tests := []struct {
		name   string
		status int
		data   interface{}
	}{
		{
			name:   "with data",
			status: http.StatusOK,
			data:   map[string]string{"key": "value"},
		},
		{
			name:   "nil data",
			status: http.StatusNotFound,
			data:   nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			JSON(c, tc.status, tc.data)

			assert.Equal(t, tc.status, w.Code)

			var resp Response
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			assert.NoError(t, err)
		})
	}
}

func TestError(t *testing.T) {
	tests := []struct {
		name       string
		status     int
		err        error
		message    string
		production bool
	}{
		{
			name:       "non-production includes debug info",
			status:     http.StatusBadRequest,
			err:        fmt.Errorf("something went wrong"),
			message:    "Bad Request",
			production: false,
		},
		{
			name:       "production hides debug info",
			status:     http.StatusInternalServerError,
			err:        fmt.Errorf("internal error"),
			message:    "Internal Server Error",
			production: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			if tc.production {
				origEnv := config.GetConfig().Environment
				config.GetConfig().Environment = config.ProductionEnv
				defer func() { config.GetConfig().Environment = origEnv }()
			}

			Error(c, tc.status, tc.err, tc.message)

			assert.Equal(t, tc.status, w.Code)

			var resp Response
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			assert.NoError(t, err)
		})
	}
}
