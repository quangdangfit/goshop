package http

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	dbsMocks "goshop/pkg/dbs/mocks"
)

func TestRoutesRegistersEndpoints(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	g := r.Group("/api/v1")
	Routes(g, dbsMocks.NewDatabase(t))

	paths := map[string]bool{}
	for _, ri := range r.Routes() {
		paths[ri.Method+" "+ri.Path] = true
	}
	require.True(t, paths["GET /api/v1/me/notification-preferences"])
	require.True(t, paths["PUT /api/v1/me/notification-preferences"])
}
