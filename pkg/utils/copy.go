package utils

import (
	"encoding/json"

	"gitlab.com/quangdangfit/gocommon/utils/logger"
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
