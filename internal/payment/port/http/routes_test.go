package http

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/quangdangfit/gocommon/validation"
	"github.com/stretchr/testify/require"

	"goshop/pkg/config"
	dbsMocks "goshop/pkg/dbs/mocks"
)

func TestRoutesRegistersEndpoints(t *testing.T) {
	// Routes() reads pkg/config; ensure it's loaded so we don't panic on a nil pointer.
	_ = config.GetConfig()
	gin.SetMode(gin.TestMode)
	r := gin.New()
	g := r.Group("/api/v1")
	Routes(g, dbsMocks.NewDatabase(t), validation.New())

	paths := map[string]bool{}
	for _, ri := range r.Routes() {
		paths[ri.Method+" "+ri.Path] = true
	}
	require.True(t, paths["POST /api/v1/orders/:id/payment-intent"])
	require.True(t, paths["POST /api/v1/webhooks/stripe"])
	require.True(t, paths["GET /api/v1/config/public"])
}
