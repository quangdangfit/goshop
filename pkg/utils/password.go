package utils

import (
	"github.com/quangdangfit/gocommon/logger"
	"golang.org/x/crypto/bcrypt"
)

func HashAndSalt(pass []byte) string {
	hashed, err := bcrypt.GenerateFromPassword(pass, bcrypt.MinCost)
	if err != nil {
		logger.Error("Failed to generate password: ", err)
		return ""
	}

	return string(hashed)
}
