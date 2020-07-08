package utils

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/copier"
	"gitlab.com/quangdangfit/gocommon/utils/logger"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

const (
	TokenExpiredTime = 3600
)

func GenerateToken(payload interface{}) string {
	tokenContent := jwt.MapClaims{
		"payload": payload,
		"expiry":  time.Now().Add(time.Second * TokenExpiredTime).Unix(),
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
	_, err := jwt.ParseWithClaims(cleanJWT, tokenData, func(token *jwt.Token) (interface{}, error) {
		return []byte("TokenPassword"), nil
	})

	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	copier.Copy(&data, tokenData["payload"])
	return data, nil
}

func HashAndSalt(pass []byte) string {
	hashed, err := bcrypt.GenerateFromPassword(pass, bcrypt.MinCost)
	if err != nil {
		logger.Error("Failed to generate password: ", err)
		return ""
	}

	return string(hashed)
}
