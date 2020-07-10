package services

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/quangdangfit/gocommon/utils/logger"
)

type Service struct {
}

func (s *Service) BindQuery(c *gin.Context, queryModel interface{}, query map[string]interface{}) {
	if err := c.ShouldBindQuery(queryModel); err != nil {
		logger.Error("Failed to parse request query: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	data, _ := json.Marshal(queryModel)
	json.Unmarshal(data, &query)

	query["active"] = query["active"] == "true"
}
