package utils

import (
	"encoding/json"

	"github.com/quangdangfit/gocommon/logger"
)

func Copy(dest interface{}, src interface{}) error {
	data, err := json.Marshal(src)
	if err != nil {
		logger.Error("Failed to marshal data")
		return err
	}

	json.Unmarshal(data, dest)

	return nil
}
