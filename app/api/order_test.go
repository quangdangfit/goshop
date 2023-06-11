package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"goshop/app/dbs"
	"goshop/app/models"
	"goshop/app/serializers"
	"goshop/pkg/jtoken"
)

// Place Order
// =================================================================================================

func TestOrderAPI_PlaceOrderSuccess(t *testing.T) {
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

	req := &serializers.PlaceOrderReq{
		Lines: []serializers.PlaceOrderLineReq{
			{
				ProductID: p1.ID,
				Quantity:  2,
			},
			{
				ProductID: p2.ID,
				Quantity:  3,
			},
		},
	}
	writer := makeRequest("POST", "/api/v1/orders", req, accessToken())
	var res serializers.Order
	parseResponseResult(writer.Body.Bytes(), &res)
	assert.Equal(t, http.StatusOK, writer.Code)
	assert.Equal(t, "new", res.Status)
	assert.Equal(t, float64(8), res.TotalPrice)
	assert.Equal(t, 2, len(res.Lines))
	assert.Equal(t, req.Lines[0].ProductID, res.Lines[0].Product.ID)
	assert.Equal(t, req.Lines[0].Quantity, res.Lines[0].Quantity)
	assert.Equal(t, float64(2), res.Lines[0].Price)

	assert.Equal(t, req.Lines[1].ProductID, res.Lines[1].Product.ID)
	assert.Equal(t, req.Lines[1].Quantity, res.Lines[1].Quantity)
	assert.Equal(t, float64(6), res.Lines[1].Price)

	// clean up
	dbs.Database.Where("1 = 1").Delete(&models.OrderLine{})
	dbs.Database.Where("1 = 1").Delete(&models.Product{})
	dbs.Database.Where("1 = 1").Delete(&models.Order{})
}

func TestOrderAPI_PlaceOrderInvalidFieldType(t *testing.T) {
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

	req := map[string]interface{}{
		"lines": []map[string]interface{}{
			{
				"product_id": p1.ID,
				"quantity":   2,
			},
			{
				"product_id": p2.ID,
				"quantity":   "1",
			},
		},
	}
	writer := makeRequest("POST", "/api/v1/orders", req, accessToken())
	var response map[string]map[string]string
	_ = json.Unmarshal(writer.Body.Bytes(), &response)
	assert.Equal(t, http.StatusBadRequest, writer.Code)
	assert.Equal(t, "Invalid parameters", response["error"]["message"])

	// clean up
	dbs.Database.Where("1 = 1").Delete(&models.OrderLine{})
	dbs.Database.Where("1 = 1").Delete(&models.Product{})
	dbs.Database.Where("1 = 1").Delete(&models.Order{})
}
func TestOrderAPI_PlaceOrderInvalidMissProductID(t *testing.T) {
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

	req := &serializers.PlaceOrderReq{
		Lines: []serializers.PlaceOrderLineReq{
			{
				Quantity: 2,
			},
			{
				ProductID: p2.ID,
				Quantity:  3,
			},
		},
	}
	writer := makeRequest("POST", "/api/v1/orders", req, accessToken())
	var response map[string]map[string]string
	_ = json.Unmarshal(writer.Body.Bytes(), &response)
	assert.Equal(t, http.StatusBadRequest, writer.Code)
	assert.Equal(t, "Invalid parameters", response["error"]["message"])

	// clean up
	dbs.Database.Where("1 = 1").Delete(&models.OrderLine{})
	dbs.Database.Where("1 = 1").Delete(&models.Product{})
	dbs.Database.Where("1 = 1").Delete(&models.Order{})
}
func TestOrderAPI_PlaceOrderInvalidMissQuantity(t *testing.T) {
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

	req := &serializers.PlaceOrderReq{
		Lines: []serializers.PlaceOrderLineReq{
			{
				ProductID: p1.ID,
				Quantity:  2,
			},
			{
				ProductID: p2.ID,
			},
		},
	}
	writer := makeRequest("POST", "/api/v1/orders", req, accessToken())
	var response map[string]map[string]string
	_ = json.Unmarshal(writer.Body.Bytes(), &response)
	assert.Equal(t, http.StatusBadRequest, writer.Code)
	assert.Equal(t, "Invalid parameters", response["error"]["message"])

	// clean up
	dbs.Database.Where("1 = 1").Delete(&models.OrderLine{})
	dbs.Database.Where("1 = 1").Delete(&models.Product{})
	dbs.Database.Where("1 = 1").Delete(&models.Order{})
}
func TestOrderAPI_PlaceOrderInvalidProductNotFound(t *testing.T) {
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

	req := &serializers.PlaceOrderReq{
		Lines: []serializers.PlaceOrderLineReq{
			{
				ProductID: p1.ID,
				Quantity:  2,
			},
			{
				ProductID: "notfound",
				Quantity:  1,
			},
		},
	}
	writer := makeRequest("POST", "/api/v1/orders", req, accessToken())
	var response map[string]map[string]string
	_ = json.Unmarshal(writer.Body.Bytes(), &response)
	assert.Equal(t, http.StatusInternalServerError, writer.Code)
	assert.Equal(t, "Something went wrong", response["error"]["message"])

	// clean up
	dbs.Database.Where("1 = 1").Delete(&models.OrderLine{})
	dbs.Database.Where("1 = 1").Delete(&models.Product{})
	dbs.Database.Where("1 = 1").Delete(&models.Order{})
}
func TestOrderAPI_PlaceOrderUnauthorized(t *testing.T) {
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

	req := &serializers.PlaceOrderReq{
		Lines: []serializers.PlaceOrderLineReq{
			{
				ProductID: p1.ID,
				Quantity:  2,
			},
			{
				ProductID: p2.ID,
				Quantity:  3,
			},
		},
	}

	writer := makeRequest("POST", "/api/v1/orders", req, "")
	assert.Equal(t, http.StatusUnauthorized, writer.Code)

	// clean up
	dbs.Database.Where("1 = 1").Delete(&models.OrderLine{})
	dbs.Database.Where("1 = 1").Delete(&models.Product{})
	dbs.Database.Where("1 = 1").Delete(&models.Order{})
}

// Get Order Detail
// =================================================================================================

func TestOrderAPI_GetOrderByIDSuccess(t *testing.T) {
	u := models.User{
		Email:    "test1@gmail.com",
		Password: "test123456",
	}
	dbs.Database.Create(&u)
	token := jtoken.GenerateAccessToken(map[string]interface{}{"id": u.ID})

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

	o := models.Order{
		UserID: u.ID,
		Lines: []*models.OrderLine{
			{
				ProductID: p1.ID,
				Quantity:  2,
			},
			{
				ProductID: p2.ID,
				Quantity:  3,
			},
		},
	}
	dbs.Database.Create(&o)

	writer := makeRequest("GET", fmt.Sprintf("/api/v1/orders/%s", o.ID), nil, token)
	var res serializers.Order
	parseResponseResult(writer.Body.Bytes(), &res)
	assert.Equal(t, http.StatusOK, writer.Code)
	assert.Equal(t, "new", res.Status)
	assert.Equal(t, 2, len(res.Lines))
	assert.Equal(t, o.Lines[0].ProductID, res.Lines[0].Product.ID)
	assert.Equal(t, o.Lines[0].Quantity, res.Lines[0].Quantity)

	assert.Equal(t, o.Lines[1].ProductID, res.Lines[1].Product.ID)
	assert.Equal(t, o.Lines[1].Quantity, res.Lines[1].Quantity)

	// clean up
	dbs.Database.Where("1 = 1").Delete(&models.OrderLine{})
	dbs.Database.Where("1 = 1").Delete(&models.Product{})
	dbs.Database.Where("1 = 1").Delete(&models.Order{})
}

func TestOrderAPI_GetOrderByIDNotFound(t *testing.T) {
	writer := makeRequest("GET", "/api/v1/orders/notfound", nil, accessToken())
	var response map[string]map[string]string
	_ = json.Unmarshal(writer.Body.Bytes(), &response)
	assert.Equal(t, http.StatusNotFound, writer.Code)
	assert.Equal(t, "Not found", response["error"]["message"])
}

// Cancel Order
// =================================================================================================

func TestOrderAPI_CancelOrderSuccess(t *testing.T) {
	u := models.User{
		Email:    "test1@test.com",
		Password: "test123456",
	}
	dbs.Database.Create(&u)
	token := jtoken.GenerateAccessToken(map[string]interface{}{"id": u.ID})

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

	o := models.Order{
		UserID: u.ID,
		Lines: []*models.OrderLine{
			{
				ProductID: p1.ID,
				Quantity:  2,
			},
			{
				ProductID: p2.ID,
				Quantity:  3,
			},
		},
	}
	dbs.Database.Create(&o)

	writer := makeRequest("PUT", fmt.Sprintf("/api/v1/orders/%s/cancel", o.ID), nil, token)
	assert.Equal(t, http.StatusOK, writer.Code)

	// clean up
	dbs.Database.Where("1 = 1").Delete(&models.OrderLine{})
	dbs.Database.Where("1 = 1").Delete(&models.Product{})
	dbs.Database.Where("1 = 1").Delete(&models.Order{})
	dbs.Database.Where("1 = 1").Delete(u)
}

func TestOrderAPI_CancelOrderNotFound(t *testing.T) {
	writer := makeRequest("PUT", "/api/v1/orders/notfound/cancel", nil, accessToken())
	var response map[string]map[string]string
	_ = json.Unmarshal(writer.Body.Bytes(), &response)
	assert.Equal(t, http.StatusInternalServerError, writer.Code)
	assert.Equal(t, "Something went wrong", response["error"]["message"])
}

func TestOrderAPI_CancelOrderStatusDone(t *testing.T) {
	u := models.User{
		Email:    "test1@test.com",
		Password: "test123456",
	}
	dbs.Database.Create(&u)
	token := jtoken.GenerateAccessToken(map[string]interface{}{"id": u.ID})

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

	o := models.Order{
		UserID: u.ID,
		Lines: []*models.OrderLine{
			{
				ProductID: p1.ID,
				Quantity:  2,
			},
			{
				ProductID: p2.ID,
				Quantity:  3,
			},
		},
		Status: models.OrderStatusDone,
	}
	dbs.Database.Create(&o)

	writer := makeRequest("PUT", fmt.Sprintf("/api/v1/orders/%s/cancel", o.ID), nil, token)
	var response map[string]map[string]string
	_ = json.Unmarshal(writer.Body.Bytes(), &response)
	assert.Equal(t, http.StatusInternalServerError, writer.Code)
	assert.Equal(t, "Something went wrong", response["error"]["message"])

	// clean up
	dbs.Database.Where("1 = 1").Delete(&models.OrderLine{})
	dbs.Database.Where("1 = 1").Delete(&models.Product{})
	dbs.Database.Where("1 = 1").Delete(&models.Order{})
	dbs.Database.Where("1 = 1").Delete(u)
}

func TestOrderAPI_CancelOrderStatusCancelled(t *testing.T) {
	u := models.User{
		Email:    "test1@test.com",
		Password: "test123456",
	}
	dbs.Database.Create(&u)
	token := jtoken.GenerateAccessToken(map[string]interface{}{"id": u.ID})

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

	o := models.Order{
		UserID: u.ID,
		Lines: []*models.OrderLine{
			{
				ProductID: p1.ID,
				Quantity:  2,
			},
			{
				ProductID: p2.ID,
				Quantity:  3,
			},
		},
		Status: models.OrderStatusCancelled,
	}
	dbs.Database.Create(&o)

	writer := makeRequest("PUT", fmt.Sprintf("/api/v1/orders/%s/cancel", o.ID), nil, token)
	var response map[string]map[string]string
	_ = json.Unmarshal(writer.Body.Bytes(), &response)
	assert.Equal(t, http.StatusInternalServerError, writer.Code)
	assert.Equal(t, "Something went wrong", response["error"]["message"])

	// clean up
	dbs.Database.Where("1 = 1").Delete(&models.OrderLine{})
	dbs.Database.Where("1 = 1").Delete(&models.Product{})
	dbs.Database.Where("1 = 1").Delete(&models.Order{})
	dbs.Database.Where("1 = 1").Delete(u)
}

func TestOrderAPI_CancelOrderNotMine(t *testing.T) {
	u := models.User{
		Email:    "test1@test.com",
		Password: "test123456",
	}
	dbs.Database.Create(&u)

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

	o := models.Order{
		UserID: u.ID,
		Lines: []*models.OrderLine{
			{
				ProductID: p1.ID,
				Quantity:  2,
			},
			{
				ProductID: p2.ID,
				Quantity:  3,
			},
		},
		Status: models.OrderStatusNew,
	}
	dbs.Database.Create(&o)

	writer := makeRequest("PUT", fmt.Sprintf("/api/v1/orders/%s/cancel", o.ID), nil, accessToken())
	var response map[string]map[string]string
	_ = json.Unmarshal(writer.Body.Bytes(), &response)
	assert.Equal(t, http.StatusInternalServerError, writer.Code)
	assert.Equal(t, "Something went wrong", response["error"]["message"])

	// clean up
	dbs.Database.Where("1 = 1").Delete(&models.OrderLine{})
	dbs.Database.Where("1 = 1").Delete(&models.Product{})
	dbs.Database.Where("1 = 1").Delete(&models.Order{})
	dbs.Database.Where("1 = 1").Delete(u)
}

// List My Orders
// =================================================================================================

func TestOrderAPI_ListProductsSuccess(t *testing.T) {
	u := models.User{
		Email:    "test1@test.com",
		Password: "test123456",
	}
	dbs.Database.Create(&u)
	token := jtoken.GenerateAccessToken(map[string]interface{}{"id": u.ID})

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

	o1 := models.Order{
		UserID: u.ID,
		Lines: []*models.OrderLine{
			{
				ProductID: p1.ID,
				Quantity:  5,
			},
		},
		Status: models.OrderStatusDone,
	}
	dbs.Database.Create(&o1)

	o2 := models.Order{
		UserID: u.ID,
		Lines: []*models.OrderLine{
			{
				ProductID: p2.ID,
				Quantity:  3,
			},
		},
		Status: models.OrderStatusNew,
	}
	dbs.Database.Create(&o2)

	writer := makeRequest("GET", "/api/v1/orders", nil, token)
	var res serializers.ListOrderRes
	parseResponseResult(writer.Body.Bytes(), &res)
	assert.Equal(t, http.StatusOK, writer.Code)
	assert.Equal(t, int64(2), res.Pagination.Total)
	assert.Equal(t, int64(1), res.Pagination.CurrentPage)
	assert.Equal(t, int64(1), res.Pagination.TotalPage)
	assert.Equal(t, int64(20), res.Pagination.Limit)
	assert.Equal(t, 2, len(res.Orders))
	assert.Equal(t, o1.Lines[0].ProductID, res.Orders[0].Lines[0].Product.ID)
	assert.Equal(t, o1.Lines[0].Quantity, res.Orders[0].Lines[0].Quantity)

	assert.Equal(t, o2.Lines[0].ProductID, res.Orders[1].Lines[0].Product.ID)
	assert.Equal(t, o2.Lines[0].Quantity, res.Orders[1].Lines[0].Quantity)

	// clean up
	dbs.Database.Where("1 = 1").Delete(&models.OrderLine{})
	dbs.Database.Where("1 = 1").Delete(&models.Product{})
	dbs.Database.Where("1 = 1").Delete(&models.Order{})
	dbs.Database.Where("1 = 1").Delete(u)
}

func TestOrderAPI_ListProductsNotFound(t *testing.T) {
	writer := makeRequest("GET", "/api/v1/orders", nil, accessToken())
	var res serializers.ListProductRes
	parseResponseResult(writer.Body.Bytes(), &res)
	assert.Equal(t, http.StatusOK, writer.Code)
	assert.Equal(t, 0, len(res.Products))
}

func TestOrderAPI_ListProductsInvalidFieldType(t *testing.T) {
	writer := makeRequest("GET", "/api/v1/orders?page=a", nil, accessToken())
	var response map[string]map[string]string
	_ = json.Unmarshal(writer.Body.Bytes(), &response)
	assert.Equal(t, http.StatusBadRequest, writer.Code)
	assert.Equal(t, "Invalid parameters", response["error"]["message"])
}

func TestOrderAPI_ListMyOrdersFindByStatusSuccess(t *testing.T) {
	u := models.User{
		Email:    "test1@test.com",
		Password: "test123456",
	}
	dbs.Database.Create(&u)
	token := jtoken.GenerateAccessToken(map[string]interface{}{"id": u.ID})

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

	o1 := models.Order{
		UserID: u.ID,
		Lines: []*models.OrderLine{
			{
				ProductID: p1.ID,
				Quantity:  5,
			},
		},
		Status: models.OrderStatusDone,
	}
	dbs.Database.Create(&o1)

	o2 := models.Order{
		UserID: u.ID,
		Lines: []*models.OrderLine{
			{
				ProductID: p2.ID,
				Quantity:  3,
			},
		},
		Status: models.OrderStatusNew,
	}
	dbs.Database.Create(&o2)

	writer := makeRequest("GET", "/api/v1/orders?status=new", nil, token)
	var res serializers.ListOrderRes
	parseResponseResult(writer.Body.Bytes(), &res)
	assert.Equal(t, http.StatusOK, writer.Code)
	assert.Equal(t, int64(1), res.Pagination.Total)
	assert.Equal(t, int64(1), res.Pagination.CurrentPage)
	assert.Equal(t, int64(1), res.Pagination.TotalPage)
	assert.Equal(t, int64(20), res.Pagination.Limit)
	assert.Equal(t, 1, len(res.Orders))
	assert.Equal(t, o2.Lines[0].ProductID, res.Orders[0].Lines[0].Product.ID)
	assert.Equal(t, o2.Lines[0].Quantity, res.Orders[0].Lines[0].Quantity)

	// clean up
	dbs.Database.Where("1 = 1").Delete(&models.OrderLine{})
	dbs.Database.Where("1 = 1").Delete(&models.Product{})
	dbs.Database.Where("1 = 1").Delete(&models.Order{})
	dbs.Database.Where("1 = 1").Delete(u)
}

func TestOrderAPI_ListProductsFindByStatusNotFound(t *testing.T) {
	u := models.User{
		Email:    "test1@test.com",
		Password: "test123456",
	}
	dbs.Database.Create(&u)
	token := jtoken.GenerateAccessToken(map[string]interface{}{"id": u.ID})

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

	o1 := models.Order{
		UserID: u.ID,
		Lines: []*models.OrderLine{
			{
				ProductID: p1.ID,
				Quantity:  5,
			},
		},
		Status: models.OrderStatusDone,
	}
	dbs.Database.Create(&o1)

	o2 := models.Order{
		UserID: u.ID,
		Lines: []*models.OrderLine{
			{
				ProductID: p2.ID,
				Quantity:  3,
			},
		},
		Status: models.OrderStatusNew,
	}
	dbs.Database.Create(&o2)

	writer := makeRequest("GET", "/api/v1/orders?status=cancelled", nil, token)
	var res serializers.ListOrderRes
	parseResponseResult(writer.Body.Bytes(), &res)
	assert.Equal(t, http.StatusOK, writer.Code)
	assert.Equal(t, 0, len(res.Orders))

	// clean up
	dbs.Database.Where("1 = 1").Delete(&models.OrderLine{})
	dbs.Database.Where("1 = 1").Delete(&models.Product{})
	dbs.Database.Where("1 = 1").Delete(&models.Order{})
	dbs.Database.Where("1 = 1").Delete(u)
}

func TestOrderAPI_ListProductsFindByCodeSuccess(t *testing.T) {
	u := models.User{
		Email:    "test1@test.com",
		Password: "test123456",
	}
	dbs.Database.Create(&u)
	token := jtoken.GenerateAccessToken(map[string]interface{}{"id": u.ID})

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

	o1 := models.Order{
		UserID: u.ID,
		Lines: []*models.OrderLine{
			{
				ProductID: p1.ID,
				Quantity:  5,
			},
		},
		Status: models.OrderStatusDone,
	}
	dbs.Database.Create(&o1)

	o2 := models.Order{
		UserID: u.ID,
		Lines: []*models.OrderLine{
			{
				ProductID: p2.ID,
				Quantity:  3,
			},
		},
		Status: models.OrderStatusNew,
	}
	dbs.Database.Create(&o2)

	writer := makeRequest("GET", fmt.Sprintf("/api/v1/orders?code=%s", o1.Code), nil, token)
	var res serializers.ListOrderRes
	parseResponseResult(writer.Body.Bytes(), &res)
	assert.Equal(t, http.StatusOK, writer.Code)
	assert.Equal(t, int64(1), res.Pagination.Total)
	assert.Equal(t, int64(1), res.Pagination.CurrentPage)
	assert.Equal(t, int64(1), res.Pagination.TotalPage)
	assert.Equal(t, int64(20), res.Pagination.Limit)
	assert.Equal(t, 1, len(res.Orders))
	assert.Equal(t, o1.Lines[0].ProductID, res.Orders[0].Lines[0].Product.ID)
	assert.Equal(t, o1.Lines[0].Quantity, res.Orders[0].Lines[0].Quantity)

	// clean up
	dbs.Database.Where("1 = 1").Delete(&models.OrderLine{})
	dbs.Database.Where("1 = 1").Delete(&models.Product{})
	dbs.Database.Where("1 = 1").Delete(&models.Order{})
	dbs.Database.Where("1 = 1").Delete(u)
}

func TestOrderAPI_ListProductsFindByCodeNotFound(t *testing.T) {
	u := models.User{
		Email:    "test1@test.com",
		Password: "test123456",
	}
	dbs.Database.Create(&u)
	token := jtoken.GenerateAccessToken(map[string]interface{}{"id": u.ID})

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

	o1 := models.Order{
		UserID: u.ID,
		Lines: []*models.OrderLine{
			{
				ProductID: p1.ID,
				Quantity:  5,
			},
		},
		Status: models.OrderStatusDone,
	}
	dbs.Database.Create(&o1)

	o2 := models.Order{
		UserID: u.ID,
		Lines: []*models.OrderLine{
			{
				ProductID: p2.ID,
				Quantity:  3,
			},
		},
		Status: models.OrderStatusNew,
	}
	dbs.Database.Create(&o2)

	writer := makeRequest("GET", "/api/v1/orders?code=notfound", nil, token)
	var res serializers.ListOrderRes
	parseResponseResult(writer.Body.Bytes(), &res)
	assert.Equal(t, http.StatusOK, writer.Code)
	assert.Equal(t, 0, len(res.Orders))

	// clean up
	dbs.Database.Where("1 = 1").Delete(&models.OrderLine{})
	dbs.Database.Where("1 = 1").Delete(&models.Product{})
	dbs.Database.Where("1 = 1").Delete(&models.Order{})
	dbs.Database.Where("1 = 1").Delete(u)
}

func TestOrderAPI_ListProductsWithPagination(t *testing.T) {
	u := models.User{
		Email:    "test1@test.com",
		Password: "test123456",
	}
	dbs.Database.Create(&u)
	token := jtoken.GenerateAccessToken(map[string]interface{}{"id": u.ID})

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

	o1 := models.Order{
		UserID: u.ID,
		Lines: []*models.OrderLine{
			{
				ProductID: p1.ID,
				Quantity:  5,
			},
		},
		Status: models.OrderStatusDone,
	}
	dbs.Database.Create(&o1)

	o2 := models.Order{
		UserID: u.ID,
		Lines: []*models.OrderLine{
			{
				ProductID: p2.ID,
				Quantity:  3,
			},
		},
		Status: models.OrderStatusNew,
	}
	dbs.Database.Create(&o2)

	writer := makeRequest("GET", "/api/v1/orders?page=2&limit=1", nil, token)
	var res serializers.ListOrderRes
	parseResponseResult(writer.Body.Bytes(), &res)
	assert.Equal(t, http.StatusOK, writer.Code)
	assert.Equal(t, int64(2), res.Pagination.Total)
	assert.Equal(t, int64(2), res.Pagination.CurrentPage)
	assert.Equal(t, int64(2), res.Pagination.TotalPage)
	assert.Equal(t, int64(1), res.Pagination.Limit)
	assert.Equal(t, 1, len(res.Orders))
	assert.Equal(t, o2.Lines[0].ProductID, res.Orders[0].Lines[0].Product.ID)
	assert.Equal(t, o2.Lines[0].Quantity, res.Orders[0].Lines[0].Quantity)

	// clean up
	dbs.Database.Where("1 = 1").Delete(&models.OrderLine{})
	dbs.Database.Where("1 = 1").Delete(&models.Product{})
	dbs.Database.Where("1 = 1").Delete(&models.Order{})
	dbs.Database.Where("1 = 1").Delete(u)
}

func TestOrderAPI_ListProductsWithOrder(t *testing.T) {
	u := models.User{
		Email:    "test1@test.com",
		Password: "test123456",
	}
	dbs.Database.Create(&u)
	token := jtoken.GenerateAccessToken(map[string]interface{}{"id": u.ID})

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

	o1 := models.Order{
		UserID: u.ID,
		Lines: []*models.OrderLine{
			{
				ProductID: p1.ID,
				Quantity:  5,
			},
		},
		Status: models.OrderStatusDone,
	}
	dbs.Database.Create(&o1)

	o2 := models.Order{
		UserID: u.ID,
		Lines: []*models.OrderLine{
			{
				ProductID: p2.ID,
				Quantity:  3,
			},
		},
		Status: models.OrderStatusNew,
	}
	dbs.Database.Create(&o2)

	writer := makeRequest("GET", "/api/v1/orders?order_by=created_at&order_desc=true", nil, token)
	var res serializers.ListOrderRes
	parseResponseResult(writer.Body.Bytes(), &res)
	assert.Equal(t, http.StatusOK, writer.Code)
	assert.Equal(t, int64(2), res.Pagination.Total)
	assert.Equal(t, int64(1), res.Pagination.CurrentPage)
	assert.Equal(t, int64(1), res.Pagination.TotalPage)
	assert.Equal(t, int64(20), res.Pagination.Limit)
	assert.Equal(t, 2, len(res.Orders))
	assert.Equal(t, o2.Lines[0].ProductID, res.Orders[0].Lines[0].Product.ID)
	assert.Equal(t, o2.Lines[0].Quantity, res.Orders[0].Lines[0].Quantity)

	assert.Equal(t, o1.Lines[0].ProductID, res.Orders[1].Lines[0].Product.ID)
	assert.Equal(t, o1.Lines[0].Quantity, res.Orders[1].Lines[0].Quantity)

	// clean up
	dbs.Database.Where("1 = 1").Delete(&models.OrderLine{})
	dbs.Database.Where("1 = 1").Delete(&models.Product{})
	dbs.Database.Where("1 = 1").Delete(&models.Order{})
	dbs.Database.Where("1 = 1").Delete(u)
}

func TestOrderAPI_ListProductsNotMine(t *testing.T) {
	u := models.User{
		Email:    "test1@test.com",
		Password: "test123456",
	}
	dbs.Database.Create(&u)

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

	o1 := models.Order{
		UserID: u.ID,
		Lines: []*models.OrderLine{
			{
				ProductID: p1.ID,
				Quantity:  5,
			},
		},
		Status: models.OrderStatusDone,
	}
	dbs.Database.Create(&o1)

	o2 := models.Order{
		UserID: u.ID,
		Lines: []*models.OrderLine{
			{
				ProductID: p2.ID,
				Quantity:  3,
			},
		},
		Status: models.OrderStatusNew,
	}
	dbs.Database.Create(&o2)

	writer := makeRequest("GET", "/api/v1/orders?code=notfound", nil, accessToken())
	var res serializers.ListOrderRes
	parseResponseResult(writer.Body.Bytes(), &res)
	assert.Equal(t, http.StatusOK, writer.Code)
	assert.Equal(t, 0, len(res.Orders))

	// clean up
	dbs.Database.Where("1 = 1").Delete(&models.OrderLine{})
	dbs.Database.Where("1 = 1").Delete(&models.Product{})
	dbs.Database.Where("1 = 1").Delete(&models.Order{})
	dbs.Database.Where("1 = 1").Delete(u)
}
