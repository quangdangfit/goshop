package test

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/quangdangfit/gocommon/validation"
	"github.com/stretchr/testify/assert"

	"goshop/app/api"
	"goshop/app/dbs"
	"goshop/app/models"
	"goshop/app/serializers"
	"goshop/app/services"
	"goshop/mocks"
)

// Get Product Detail
// =================================================================================================

func TestProductAPI_GetProductByIDSuccess(t *testing.T) {
	defer cleanData()

	p := models.Product{
		Name:        "test-product",
		Description: "test-product",
		Price:       1,
	}
	dbs.Database.Create(&p)

	writer := makeRequest("GET", fmt.Sprintf("/api/v1/products/%s", p.ID), nil, accessToken())
	var res models.Product
	parseResponseResult(writer.Body.Bytes(), &res)
	assert.Equal(t, http.StatusOK, writer.Code)
	assert.Equal(t, "test-product", res.Name)
	assert.Equal(t, "test-product", res.Description)
	assert.Equal(t, float64(1), res.Price)
}

func TestProductAPI_GetProductByIDSuccessFromCache(t *testing.T) {
	defer cleanData()

	p := models.Product{
		Name:        "test-product",
		Description: "test-product",
		Price:       1,
	}
	dbs.Database.Create(&p)

	writer := makeRequest("GET", fmt.Sprintf("/api/v1/products/%s", p.ID), nil, accessToken())
	var res models.Product
	parseResponseResult(writer.Body.Bytes(), &res)
	assert.Equal(t, http.StatusOK, writer.Code)
	assert.Equal(t, "test-product", res.Name)
	assert.Equal(t, "test-product", res.Description)
	assert.Equal(t, float64(1), res.Price)

	writer = makeRequest("GET", fmt.Sprintf("/api/v1/products/%s", p.ID), nil, accessToken())
	parseResponseResult(writer.Body.Bytes(), &res)
	assert.Equal(t, http.StatusOK, writer.Code)
	assert.Equal(t, "test-product", res.Name)
	assert.Equal(t, "test-product", res.Description)
	assert.Equal(t, float64(1), res.Price)
}

func TestProductAPI_GetProductByIDNotFound(t *testing.T) {
	defer cleanData()
	writer := makeRequest("GET", "/api/v1/products/notfound", nil, accessToken())
	var response map[string]map[string]string
	_ = json.Unmarshal(writer.Body.Bytes(), &response)
	assert.Equal(t, http.StatusNotFound, writer.Code)
	assert.Equal(t, "Not found", response["error"]["message"])
}

// Get List Products
// =================================================================================================

func TestProductAPI_ListProductsSuccess(t *testing.T) {
	defer cleanData()

	p := models.Product{
		Name:        "test-product",
		Description: "test-product",
		Price:       1,
	}
	dbs.Database.Create(&p)

	writer := makeRequest("GET", "/api/v1/products", nil, accessToken())
	var res serializers.ListProductRes
	parseResponseResult(writer.Body.Bytes(), &res)
	assert.Equal(t, http.StatusOK, writer.Code)
	assert.Equal(t, int64(1), res.Pagination.Total)
	assert.Equal(t, int64(1), res.Pagination.CurrentPage)
	assert.Equal(t, int64(1), res.Pagination.TotalPage)
	assert.Equal(t, int64(20), res.Pagination.Limit)
	assert.Equal(t, 1, len(res.Products))
	assert.Equal(t, "test-product", res.Products[0].Name)
	assert.Equal(t, "test-product", res.Products[0].Description)
	assert.Equal(t, float64(1), res.Products[0].Price)
}

func TestProductAPI_ListProductsSuccessFromCache(t *testing.T) {
	defer cleanData()

	p := models.Product{
		Name:        "test-product",
		Description: "test-product",
		Price:       1,
	}
	dbs.Database.Create(&p)

	writer := makeRequest("GET", "/api/v1/products", nil, accessToken())
	var res serializers.ListProductRes
	parseResponseResult(writer.Body.Bytes(), &res)
	assert.Equal(t, http.StatusOK, writer.Code)
	assert.Equal(t, int64(1), res.Pagination.Total)
	assert.Equal(t, int64(1), res.Pagination.CurrentPage)
	assert.Equal(t, int64(1), res.Pagination.TotalPage)
	assert.Equal(t, int64(20), res.Pagination.Limit)
	assert.Equal(t, 1, len(res.Products))
	assert.Equal(t, "test-product", res.Products[0].Name)
	assert.Equal(t, "test-product", res.Products[0].Description)
	assert.Equal(t, float64(1), res.Products[0].Price)

	writer = makeRequest("GET", "/api/v1/products", nil, accessToken())
	parseResponseResult(writer.Body.Bytes(), &res)
	assert.Equal(t, http.StatusOK, writer.Code)
	assert.Equal(t, int64(1), res.Pagination.Total)
	assert.Equal(t, int64(1), res.Pagination.CurrentPage)
	assert.Equal(t, int64(1), res.Pagination.TotalPage)
	assert.Equal(t, int64(20), res.Pagination.Limit)
	assert.Equal(t, 1, len(res.Products))
	assert.Equal(t, "test-product", res.Products[0].Name)
	assert.Equal(t, "test-product", res.Products[0].Description)
	assert.Equal(t, float64(1), res.Products[0].Price)
}

func TestProductAPI_ListProductsNotFound(t *testing.T) {
	defer cleanData()
	writer := makeRequest("GET", "/api/v1/products", nil, accessToken())
	var res serializers.ListProductRes
	parseResponseResult(writer.Body.Bytes(), &res)
	assert.Equal(t, http.StatusOK, writer.Code)
	assert.Equal(t, 0, len(res.Products))
}

func TestProductAPI_ListProductsInvalidFieldType(t *testing.T) {
	defer cleanData()
	writer := makeRequest("GET", "/api/v1/products?page=a", nil, accessToken())
	var response map[string]map[string]string
	_ = json.Unmarshal(writer.Body.Bytes(), &response)
	assert.Equal(t, http.StatusBadRequest, writer.Code)
	assert.Equal(t, "Invalid parameters", response["error"]["message"])
}

func TestProductAPI_ListProductsFindByNameSuccess(t *testing.T) {
	defer cleanData()

	p := models.Product{
		Name:        "test-product",
		Description: "test-product",
		Price:       1,
	}
	dbs.Database.Create(&p)

	writer := makeRequest("GET", "/api/v1/products?name=test", nil, accessToken())
	var res serializers.ListProductRes
	parseResponseResult(writer.Body.Bytes(), &res)
	assert.Equal(t, http.StatusOK, writer.Code)
	assert.Equal(t, int64(1), res.Pagination.Total)
	assert.Equal(t, int64(1), res.Pagination.CurrentPage)
	assert.Equal(t, int64(1), res.Pagination.TotalPage)
	assert.Equal(t, int64(20), res.Pagination.Limit)
	assert.Equal(t, 1, len(res.Products))
	assert.Equal(t, "test-product", res.Products[0].Name)
	assert.Equal(t, "test-product", res.Products[0].Description)
	assert.Equal(t, float64(1), res.Products[0].Price)
}

func TestProductAPI_ListProductsFindByNameNotFound(t *testing.T) {
	defer cleanData()

	p := models.Product{
		Name:        "test-product",
		Description: "test-product",
		Price:       1,
	}
	dbs.Database.Create(&p)

	writer := makeRequest("GET", "/api/v1/products?name=notfound", nil, accessToken())
	var res serializers.ListProductRes
	parseResponseResult(writer.Body.Bytes(), &res)
	assert.Equal(t, http.StatusOK, writer.Code)
	assert.Equal(t, 0, len(res.Products))
}

func TestProductAPI_ListProductsFindByCodeSuccess(t *testing.T) {
	defer cleanData()

	p := models.Product{
		Name:        "test-product",
		Description: "test-product",
		Price:       1,
	}
	dbs.Database.Create(&p)

	writer := makeRequest("GET", fmt.Sprintf("/api/v1/products?code=%s", p.Code), nil, accessToken())
	var res serializers.ListProductRes
	parseResponseResult(writer.Body.Bytes(), &res)
	assert.Equal(t, http.StatusOK, writer.Code)
	assert.Equal(t, int64(1), res.Pagination.Total)
	assert.Equal(t, int64(1), res.Pagination.CurrentPage)
	assert.Equal(t, int64(1), res.Pagination.TotalPage)
	assert.Equal(t, int64(20), res.Pagination.Limit)
	assert.Equal(t, 1, len(res.Products))
	assert.Equal(t, "test-product", res.Products[0].Name)
	assert.Equal(t, "test-product", res.Products[0].Description)
	assert.Equal(t, float64(1), res.Products[0].Price)
}

func TestProductAPI_ListProductsFindByCodeNotFound(t *testing.T) {
	defer cleanData()

	p := models.Product{
		Name:        "test-product",
		Description: "test-product",
		Price:       1,
	}
	dbs.Database.Create(&p)

	writer := makeRequest("GET", "/api/v1/products?code=notfound", nil, accessToken())
	var res serializers.ListProductRes
	parseResponseResult(writer.Body.Bytes(), &res)
	assert.Equal(t, http.StatusOK, writer.Code)
	assert.Equal(t, 0, len(res.Products))
}

func TestProductAPI_ListProductsWithPagination(t *testing.T) {
	defer cleanData()

	p1 := models.Product{
		Name:        "test-product-1",
		Description: "test-product-1",
		Price:       1,
	}
	dbs.Database.Create(&p1)

	p2 := models.Product{
		Name:        "test-product-2",
		Description: "test-product-2",
		Price:       2,
	}
	dbs.Database.Create(&p2)

	p3 := models.Product{
		Name:        "test-product-3",
		Description: "test-product-3",
		Price:       3,
	}
	dbs.Database.Create(&p3)

	writer := makeRequest("GET", "/api/v1/products?page=2&limit=2", nil, accessToken())
	var res serializers.ListProductRes
	parseResponseResult(writer.Body.Bytes(), &res)
	assert.Equal(t, http.StatusOK, writer.Code)
	assert.Equal(t, int64(3), res.Pagination.Total)
	assert.Equal(t, int64(2), res.Pagination.CurrentPage)
	assert.Equal(t, int64(2), res.Pagination.TotalPage)
	assert.Equal(t, int64(2), res.Pagination.Limit)
	assert.Equal(t, 1, len(res.Products))
	assert.Equal(t, "test-product-3", res.Products[0].Name)
	assert.Equal(t, "test-product-3", res.Products[0].Description)
	assert.Equal(t, float64(3), res.Products[0].Price)
}

func TestProductAPI_ListProductsWithOrder(t *testing.T) {
	defer cleanData()

	p1 := models.Product{
		Name:        "test-product-1",
		Description: "test-product-1",
		Price:       1,
	}
	dbs.Database.Create(&p1)

	p2 := models.Product{
		Name:        "test-product-2",
		Description: "test-product-2",
		Price:       2,
	}
	dbs.Database.Create(&p2)

	p3 := models.Product{
		Name:        "test-product-3",
		Description: "test-product-3",
		Price:       3,
	}
	dbs.Database.Create(&p3)

	writer := makeRequest("GET", "/api/v1/products?order_by=name&order_desc=true", nil, accessToken())
	var res serializers.ListProductRes
	parseResponseResult(writer.Body.Bytes(), &res)
	assert.Equal(t, http.StatusOK, writer.Code)
	assert.Equal(t, int64(3), res.Pagination.Total)
	assert.Equal(t, int64(1), res.Pagination.CurrentPage)
	assert.Equal(t, int64(1), res.Pagination.TotalPage)
	assert.Equal(t, int64(20), res.Pagination.Limit)
	assert.Equal(t, 3, len(res.Products))
	assert.Equal(t, "test-product-3", res.Products[0].Name)
	assert.Equal(t, "test-product-3", res.Products[0].Description)
	assert.Equal(t, float64(3), res.Products[0].Price)
}

func TestProductAPI_ListProductsFail(t *testing.T) {
	defer cleanData()

	p := models.Product{
		Name:        "test-product",
		Description: "test-product",
		Price:       1,
	}
	dbs.Database.Create(&p)

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockRepo := mocks.NewMockIProductRepository(mockCtrl)

	productSvc := services.NewProductService(mockRepo)
	mockTestProductAPI := api.NewProductAPI(validation.New(), testRedis, productSvc)
	mockTestRouter = initGinEngine(testUserAPI, mockTestProductAPI, testOrderAPI)

	mockRepo.EXPECT().ListProducts(gomock.Any(), gomock.Any()).Return(nil, nil, errors.New("update order fail")).Times(1)

	writer := makeMockRequest("GET", "/api/v1/products", nil, accessToken())
	var response map[string]map[string]string
	_ = json.Unmarshal(writer.Body.Bytes(), &response)
	assert.Equal(t, http.StatusInternalServerError, writer.Code)
	assert.Equal(t, "Something went wrong", response["error"]["message"])
}

// Create Product
// =================================================================================================

func TestProductAPI_CreateProductSuccess(t *testing.T) {
	defer cleanData()

	p := &serializers.CreateProductReq{
		Name:        "test-product",
		Description: "test-product",
		Price:       1,
	}
	writer := makeRequest("POST", "/api/v1/products", p, accessToken())
	var res models.Product
	parseResponseResult(writer.Body.Bytes(), &res)
	assert.Equal(t, http.StatusOK, writer.Code)
	assert.Equal(t, "test-product", res.Name)
	assert.Equal(t, "test-product", res.Description)
	assert.Equal(t, float64(1), res.Price)
}

func TestProductAPI_CreateProductInvalidFieldType(t *testing.T) {
	defer cleanData()

	p := map[string]interface{}{
		"name":        "test-product",
		"description": "test-product",
		"price":       "1",
	}
	writer := makeRequest("POST", "/api/v1/products", p, accessToken())
	var response map[string]map[string]string
	_ = json.Unmarshal(writer.Body.Bytes(), &response)
	assert.Equal(t, http.StatusBadRequest, writer.Code)
	assert.Equal(t, "Invalid parameters", response["error"]["message"])
}

func TestProductAPI_CreateProductMissingName(t *testing.T) {
	defer cleanData()

	p := &serializers.CreateProductReq{
		Description: "test-product",
		Price:       1,
	}
	writer := makeRequest("POST", "/api/v1/products", p, accessToken())
	var response map[string]map[string]string
	_ = json.Unmarshal(writer.Body.Bytes(), &response)
	assert.Equal(t, http.StatusBadRequest, writer.Code)
	assert.Equal(t, "Invalid parameters", response["error"]["message"])
}

func TestProductAPI_CreateProductMissingDescription(t *testing.T) {
	defer cleanData()

	p := &serializers.CreateProductReq{
		Name:  "test-product",
		Price: 1,
	}
	writer := makeRequest("POST", "/api/v1/products", p, accessToken())
	var response map[string]map[string]string
	_ = json.Unmarshal(writer.Body.Bytes(), &response)
	assert.Equal(t, http.StatusBadRequest, writer.Code)
	assert.Equal(t, "Invalid parameters", response["error"]["message"])
}

func TestProductAPI_CreateProductPriceLessThanZero(t *testing.T) {
	defer cleanData()

	p := &serializers.CreateProductReq{
		Name:        "test-product",
		Description: "test-product",
		Price:       -1,
	}
	writer := makeRequest("POST", "/api/v1/products", p, accessToken())
	var response map[string]map[string]string
	_ = json.Unmarshal(writer.Body.Bytes(), &response)
	assert.Equal(t, http.StatusBadRequest, writer.Code)
	assert.Equal(t, "Invalid parameters", response["error"]["message"])
}

func TestProductAPI_CreateProductPriceEqualZero(t *testing.T) {
	defer cleanData()

	p := &serializers.CreateProductReq{
		Name:        "test-product",
		Description: "test-product",
		Price:       0,
	}
	writer := makeRequest("POST", "/api/v1/products", p, accessToken())
	var response map[string]map[string]string
	_ = json.Unmarshal(writer.Body.Bytes(), &response)
	assert.Equal(t, http.StatusBadRequest, writer.Code)
	assert.Equal(t, "Invalid parameters", response["error"]["message"])
}

func TestProductAPI_CreateProductDuplicateName(t *testing.T) {
	defer cleanData()

	p := models.Product{
		Name:        "test-product",
		Description: "test-product",
		Price:       1,
	}
	dbs.Database.Create(&p)

	writer := makeRequest("POST", "/api/v1/products", p, accessToken())
	var response map[string]map[string]string
	_ = json.Unmarshal(writer.Body.Bytes(), &response)
	assert.Equal(t, http.StatusInternalServerError, writer.Code)
	assert.Equal(t, "Something went wrong", response["error"]["message"])
}

// Update Product
// =================================================================================================

func TestProductAPI_UpdateProductSuccess(t *testing.T) {
	defer cleanData()

	p := models.Product{
		Name:        "test-product",
		Description: "test-product",
		Price:       1,
	}
	dbs.Database.Create(&p)

	update := &serializers.UpdateProductReq{
		Name: "update-test-product",
	}
	writer := makeRequest("PUT", fmt.Sprintf("/api/v1/products/%s", p.ID), update, accessToken())
	var res models.Product
	parseResponseResult(writer.Body.Bytes(), &res)
	assert.Equal(t, http.StatusOK, writer.Code)
	assert.Equal(t, "update-test-product", res.Name)
	assert.Equal(t, "test-product", res.Description)
	assert.Equal(t, float64(1), res.Price)
}

func TestProductAPI_UpdateProductInvalidFieldType(t *testing.T) {
	defer cleanData()

	p := models.Product{
		Name:        "test-product",
		Description: "test-product",
		Price:       1,
	}
	dbs.Database.Create(&p)

	update := map[string]interface{}{
		"price": "1",
	}
	writer := makeRequest("PUT", fmt.Sprintf("/api/v1/products/%s", p.ID), update, accessToken())
	var response map[string]map[string]string
	_ = json.Unmarshal(writer.Body.Bytes(), &response)
	assert.Equal(t, http.StatusBadRequest, writer.Code)
	assert.Equal(t, "Invalid parameters", response["error"]["message"])
}

func TestProductAPI_UpdateProductPriceLessThanZero(t *testing.T) {
	defer cleanData()

	p := models.Product{
		Name:        "test-product",
		Description: "test-product",
		Price:       1,
	}
	dbs.Database.Create(&p)

	update := &serializers.UpdateProductReq{
		Price: -1,
	}
	writer := makeRequest("PUT", fmt.Sprintf("/api/v1/products/%s", p.ID), update, accessToken())
	var response map[string]map[string]string
	_ = json.Unmarshal(writer.Body.Bytes(), &response)
	assert.Equal(t, http.StatusBadRequest, writer.Code)
	assert.Equal(t, "Invalid parameters", response["error"]["message"])
}

func TestProductAPI_UpdateProductNotFound(t *testing.T) {
	defer cleanData()
	update := &serializers.UpdateProductReq{
		Price: 1,
	}
	writer := makeRequest("PUT", "/api/v1/products/notfound", update, accessToken())
	var response map[string]map[string]string
	_ = json.Unmarshal(writer.Body.Bytes(), &response)
	assert.Equal(t, http.StatusInternalServerError, writer.Code)
	assert.Equal(t, "Something went wrong", response["error"]["message"])
}

func TestProductAPI_UpdateProductFail(t *testing.T) {
	defer cleanData()

	p := models.Product{
		Name:        "test-product",
		Description: "test-product",
		Price:       1,
	}
	dbs.Database.Create(&p)

	update := &serializers.UpdateProductReq{
		Name: "update-test-product",
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockRepo := mocks.NewMockIProductRepository(mockCtrl)

	productSvc := services.NewProductService(mockRepo)
	mockTestProductAPI := api.NewProductAPI(validation.New(), testRedis, productSvc)
	mockTestRouter = initGinEngine(testUserAPI, mockTestProductAPI, testOrderAPI)

	mockRepo.EXPECT().GetProductByID(gomock.Any(), gomock.Any()).Return(&models.Product{}, nil).Times(1)
	mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(errors.New("update order fail")).Times(1)

	writer := makeMockRequest("PUT", fmt.Sprintf("/api/v1/products/%s", p.ID), update, accessToken())
	var response map[string]map[string]string
	_ = json.Unmarshal(writer.Body.Bytes(), &response)
	assert.Equal(t, http.StatusInternalServerError, writer.Code)
	assert.Equal(t, "Something went wrong", response["error"]["message"])
}
