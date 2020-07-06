package utils

func PrepareResponse(data interface{}, message string, code string) map[string]interface{} {
	result := map[string]interface{}{
		"data":    data,
		"message": message,
	}

	return result
}
