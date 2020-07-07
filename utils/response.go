package utils

const (
	Success       = "200"
	Error         = "500"
	InvalidParams = "400"

	ErrorExistProduct      = "10001"
	ErrorExistProductFail  = "10002"
	ErrorNotExistProduct   = "10003"
	ErrorGetProductsFail   = "10004"
	ErrorCountProductFail  = "10005"
	ErrorAddProductFail    = "10006"
	ErrorEditProductFail   = "10007"
	ErrorDeleteProductFail = "10008"
	ErrorExportProductFail = "10009"
	ErrorImportProductFail = "10010"

	ErrorNotExistUser       = "20001"
	ErrorCheckExistUserFail = "20002"
	ErrorAddUserFail        = "20003"
	ErrorDeleteUserFail     = "20004"
	ErrorEditUserFail       = "20005"
	ErrorCountUserFail      = "20006"
	ErrorGetUsersFail       = "20007"
	ErrorGetUserFail        = "20008"
	ErrorGenUserPosterFail  = "20009"

	ErrorAuthCheckTokenFail    = "30001"
	ErrorAuthCheckTokenTimeout = "30002"
	ErrorAuthToken             = "30003"
	ErrorAuth                  = "30004"
)

func PrepareResponse(data interface{}, message string, code string) map[string]interface{} {
	result := map[string]interface{}{
		"data":    data,
		"message": message,
		"code":    code,
	}

	return result
}
