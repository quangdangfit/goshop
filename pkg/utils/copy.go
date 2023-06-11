package utils

import (
	"encoding/json"
)

func Copy(dest interface{}, src interface{}) {
	data, _ := json.Marshal(src)
	_ = json.Unmarshal(data, dest)
}
