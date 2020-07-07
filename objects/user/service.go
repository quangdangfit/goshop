package user

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"gitlab.com/quangdangfit/gocommon/utils/logger"
	"goshop/utils"
	"net/http"
)

type Service interface {
	//Login(c *gin.Context)
	Register(c *gin.Context)
	//GetUserByID(c *gin.Context)
}

type service struct {
	repo Repository
}

func NewService() Service {
	return &service{repo: NewRepository()}
}

//func (s *service) Login(c *gin.Context) {
//	var reqBody LoginRequest
//	if err := c.ShouldBindJSON(&reqBody); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//
//	user, err := s.repo.Login(reqBody.Username, reqBody.Password)
//	if err != nil {
//		logger.Error(err.Error())
//		c.JSON(http.StatusBadRequest, utils.PrepareResponse(nil, err.Error(), ""))
//		return
//	}
//
//	var res UserResponse
//	copier.Copy(&res, &user)
//	c.JSON(http.StatusOK, utils.PrepareResponse(res, "OK", ""))
//}

func (s *service) validate(r RegisterRequest) bool {
	return utils.Validate(
		[]utils.Validation{
			{Value: r.Username, Valid: "username"},
			{Value: r.Email, Valid: "email"},
			{Value: r.Password, Valid: "password"},
		})
}

func (s *service) Register(c *gin.Context) {
	var reqBody RegisterRequest
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	valid := s.validate(reqBody)
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "token is invalid"})
		return
	}

	user, err := s.repo.Register(reqBody)
	if err != nil {
		logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, utils.PrepareResponse(nil, err.Error(), ""))
		return
	}

	var res UserResponse
	copier.Copy(&res, &user)
	res.Extra = map[string]interface{}{
		"token": utils.GenerateToken(user),
	}
	c.JSON(http.StatusOK, utils.PrepareResponse(res, "OK", ""))
}

//func (s *service) GetUserByID(c *gin.Context) {
//	userUUID := c.Param("uuid")
//	auth := c.GetHeader("Authorization")
//
//	user, err := s.repo.GetUser(userUUID, auth)
//	if err != nil {
//		logger.Error(err.Error())
//		c.JSON(http.StatusBadRequest, utils.PrepareResponse(nil, err.Error(), ""))
//		return
//	}
//	var res UserResponse
//	copier.Copy(&res, &user)
//	c.JSON(http.StatusOK, utils.PrepareResponse(res, "OK", ""))
//}
