package jwt

import (
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"gitlab.com/quangdangfit/gocommon/utils/logger"

	"goshop/utils"
)

const (
	TokenExpiredTime = 300
)

func GenerateToken(payload interface{}) string {
	tokenContent := jwt.MapClaims{
		"payload": payload,
		"exp":     time.Now().Add(time.Second * TokenExpiredTime).Unix(),
	}
	jwtToken := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tokenContent)
	token, err := jwtToken.SignedString([]byte("TokenPassword"))
	if err != nil {
		logger.Error("Failed to generate token: ", err)
		return ""
	}

	return token
}

func ValidateToken(jwtToken string) (map[string]interface{}, error) {
	cleanJWT := strings.Replace(jwtToken, "Bearer ", "", -1)
	tokenData := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(cleanJWT, tokenData, func(token *jwt.Token) (interface{}, error) {
		return []byte("TokenPassword"), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrInvalidKey
	}

	var data map[string]interface{}
	copier.Copy(&data, tokenData["payload"])
	return data, nil
}

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code string

		code = utils.Success
		token := c.GetHeader("Authorization")

		if token == "" {
			code = utils.InvalidParams
			c.JSON(http.StatusUnauthorized, utils.PrepareResponse(nil, "Unauthorized", code))

			c.Abort()
			return
		}

		_, err := ValidateToken(token)
		if err != nil {
			switch err.(*jwt.ValidationError).Errors {
			case jwt.ValidationErrorExpired:
				code = utils.ErrorAuthCheckTokenTimeout
			default:
				code = utils.ErrorAuthCheckTokenFail
			}
		}

		if code != utils.Success {
			c.JSON(http.StatusUnauthorized, utils.PrepareResponse(nil, "Unauthorized", code))

			c.Abort()
			return
		}

		c.Next()
	}
}
