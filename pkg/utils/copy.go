package utils

import (
	"encoding/json"
	"fmt"
)

func Copy(dest interface{}, src interface{}) error {
	data, err := json.Marshal(src)
	if err != nil {
		return fmt.Errorf("copy marshal: %w", err)
	}
	if err := json.Unmarshal(data, dest); err != nil {
		return fmt.Errorf("copy unmarshal: %w", err)
	}
	return nil
}
