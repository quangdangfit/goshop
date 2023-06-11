package test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"goshop/app/dbs"
	"goshop/app/models"
	"goshop/app/serializers"
)

// Get Product Detail
// =================================================================================================

func TestProductAPI_GetProductByIDSuccess(t *testing.T) {
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

	// clean up
	dbs.Database.Where("1 = 1").Delete(&models.Product{})
}

func TestProductAPI_GetProductByIDNotFound(t *testing.T) {
	writer := makeRequest("GET", "/api/v1/products/notfound", nil, accessToken())
	var response map[string]map[string]string
	_ = json.Unmarshal(writer.Body.Bytes(), &response)
	assert.Equal(t, http.StatusNotFound, writer.Code)
	assert.Equal(t, "Not found", response["error"]["message"])
}

// Get List Products
// =================================================================================================

func TestProductAPI_ListProductsSuccess(t *testing.T) {
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

	// clean up
	dbs.Database.Where("1 = 1").Delete(&models.Product{})
}

func TestProductAPI_ListProductsNotFound(t *testing.T) {
	writer := makeRequest("GET", "/api/v1/products", nil, accessToken())
	var res serializers.ListProductRes
	parseResponseResult(writer.Body.Bytes(), &res)
	assert.Equal(t, http.StatusOK, writer.Code)
	assert.Equal(t, 0, len(res.Products))
}

func TestProductAPI_ListProductsInvalidFieldType(t *testing.T) {
	writer := makeRequest("GET", "/api/v1/products?page=a", nil, accessToken())
	var response map[string]map[string]string
	_ = json.Unmarshal(writer.Body.Bytes(), &response)
	assert.Equal(t, http.StatusBadRequest, writer.Code)
	assert.Equal(t, "Invalid parameters", response["error"]["message"])
}

func TestProductAPI_ListProductsFindByNameSuccess(t *testing.T) {
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

	// clean up
	dbs.Database.Where("1 = 1").Delete(&models.Product{})
}

func TestProductAPI_ListProductsFindByNameNotFound(t *testing.T) {
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

	// clean up
	dbs.Database.Where("1 = 1").Delete(&models.Product{})
}

func TestProductAPI_ListProductsFindByCodeSuccess(t *testing.T) {
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

	// clean up
	dbs.Database.Where("1 = 1").Delete(&models.Product{})
}

func TestProductAPI_ListProductsFindByCodeNotFound(t *testing.T) {
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

	// clean up
	dbs.Database.Where("1 = 1").Delete(&models.Product{})
}

func TestProductAPI_ListProductsWithPagination(t *testing.T) {
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

	// clean up
	dbs.Database.Where("1 = 1").Delete(&models.Product{})
}

func TestProductAPI_ListProductsWithOrder(t *testing.T) {
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

	writer := makeRequest("GET", "/api/v1/products?orderBy=name&orderDesc=true", nil, accessToken())
	var res serializers.ListProductRes
	parseResponseResult(writer.Body.Bytes(), &res)
	assert.Equal(t, http.StatusOK, writer.Code)
	assert.Equal(t, int64(3), res.Pagination.Total)
	assert.Equal(t, int64(1), res.Pagination.CurrentPage)
	assert.Equal(t, int64(1), res.Pagination.TotalPage)
	assert.Equal(t, int64(20), res.Pagination.Limit)
	assert.Equal(t, 3, len(res.Products))
	assert.Equal(t, "test-product-1", res.Products[0].Name)
	assert.Equal(t, "test-product-1", res.Products[0].Description)
	assert.Equal(t, float64(1), res.Products[0].Price)

	// clean up
	dbs.Database.Where("1 = 1").Delete(&models.Product{})
}

// Create Product
// =================================================================================================

func TestProductAPI_CreateProductSuccess(t *testing.T) {
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

	// clean up
	dbs.Database.Where("1 = 1").Delete(&models.Product{})
}

func TestProductAPI_CreateProductInvalidFieldType(t *testing.T) {
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

	// clean up
	dbs.Database.Where("1 = 1").Delete(&models.Product{})
}

// Update Product
// =================================================================================================

func TestProductAPI_UpdateProductSuccess(t *testing.T) {
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

	// clean up
	dbs.Database.Where("1 = 1").Delete(&models.Product{})
}

func TestProductAPI_UpdateProductInvalidFieldType(t *testing.T) {
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

	// clean up
	dbs.Database.Where("1 = 1").Delete(&models.Product{})
}

func TestProductAPI_UpdateProductPriceLessThanZero(t *testing.T) {
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

	// clean up
	dbs.Database.Where("1 = 1").Delete(&models.Product{})
}

func TestProductAPI_UpdateProductNotFound(t *testing.T) {
	update := &serializers.UpdateProductReq{
		Price: 1,
	}
	writer := makeRequest("PUT", "/api/v1/products/notfound", update, accessToken())
	var response map[string]map[string]string
	_ = json.Unmarshal(writer.Body.Bytes(), &response)
	assert.Equal(t, http.StatusInternalServerError, writer.Code)
	assert.Equal(t, "Something went wrong", response["error"]["message"])
}
