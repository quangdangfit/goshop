package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"

	"goshop/pkg/jtoken"
	"goshop/pkg/utils"
)

func JWTAuth() gin.HandlerFunc {
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

		payload, err := jtoken.ValidateToken(token)
		if err != nil {
			switch err.(*jwt.ValidationError).Errors {
			case jwt.ValidationErrorExpired:
				code = utils.ErrorAuthCheckTokenTimeout
			default:
				code = utils.ErrorAuthCheckTokenFail
			}
		}

		if code != utils.Success || payload == nil || payload["type"] != jtoken.AccessTokenType {
			c.JSON(http.StatusUnauthorized, utils.PrepareResponse(nil, "Unauthorized", code))
			c.Abort()
			return
		}

		c.Set("userId", payload["id"])
		c.Next()
	}
}

func JWTRefresh() gin.HandlerFunc {
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

		payload, err := jtoken.ValidateToken(token)
		if err != nil {
			switch err.(*jwt.ValidationError).Errors {
			case jwt.ValidationErrorExpired:
				code = utils.ErrorAuthCheckTokenTimeout
			default:
				code = utils.ErrorAuthCheckTokenFail
			}
		}

		if code != utils.Success || payload == nil || payload["type"] != jtoken.RefreshTokenType {
			c.JSON(http.StatusUnauthorized, utils.PrepareResponse(nil, "Unauthorized", code))

			c.Abort()
			return
		}

		c.Set("userId", payload["id"])
		c.Next()
	}
}
