package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"goshop/internal/order/dto"
	"goshop/internal/order/model"
	productModel "goshop/internal/product/model"
	userModel "goshop/internal/user/model"
	"goshop/pkg/jtoken"
)

// Place Order
// =================================================================================================

func TestOrderAPI_PlaceOrderSuccess(t *testing.T) {
	defer cleanData()

	p1 := productModel.Product{
		Name:        "test-product-1",
		Description: "test-product-1",
		Price:       1,
	}
	dbTest.Create(context.Background(), &p1)

	p2 := productModel.Product{
		Name:        "test-product-2",
		Description: "test-product-2",
		Price:       2,
	}
	dbTest.Create(context.Background(), &p2)

	req := &dto.PlaceOrderReq{
		Lines: []dto.PlaceOrderLineReq{
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
	var res dto.Order
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
}

func TestOrderAPI_PlaceOrderInvalidFieldType(t *testing.T) {
	defer cleanData()

	p1 := productModel.Product{
		Name:        "test-product-1",
		Description: "test-product-1",
		Price:       1,
	}
	dbTest.Create(context.Background(), &p1)

	p2 := productModel.Product{
		Name:        "test-product-2",
		Description: "test-product-2",
		Price:       2,
	}
	dbTest.Create(context.Background(), &p2)

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
}

func TestOrderAPI_PlaceOrderInvalidMissProductID(t *testing.T) {
	defer cleanData()

	p1 := productModel.Product{
		Name:        "test-product-1",
		Description: "test-product-1",
		Price:       1,
	}
	dbTest.Create(context.Background(), &p1)

	p2 := productModel.Product{
		Name:        "test-product-2",
		Description: "test-product-2",
		Price:       2,
	}
	dbTest.Create(context.Background(), &p2)

	req := &dto.PlaceOrderReq{
		Lines: []dto.PlaceOrderLineReq{
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
	assert.Equal(t, http.StatusInternalServerError, writer.Code)
	assert.Equal(t, "Something went wrong", response["error"]["message"])
}

func TestOrderAPI_PlaceOrderInvalidMissQuantity(t *testing.T) {
	defer cleanData()

	p1 := productModel.Product{
		Name:        "test-product-1",
		Description: "test-product-1",
		Price:       1,
	}
	dbTest.Create(context.Background(), &p1)

	p2 := productModel.Product{
		Name:        "test-product-2",
		Description: "test-product-2",
		Price:       2,
	}
	dbTest.Create(context.Background(), &p2)

	req := &dto.PlaceOrderReq{
		Lines: []dto.PlaceOrderLineReq{
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
	assert.Equal(t, http.StatusInternalServerError, writer.Code)
	assert.Equal(t, "Something went wrong", response["error"]["message"])
}

func TestOrderAPI_PlaceOrderInvalidProductNotFound(t *testing.T) {
	defer cleanData()

	p1 := productModel.Product{
		Name:        "test-product-1",
		Description: "test-product-1",
		Price:       1,
	}
	dbTest.Create(context.Background(), &p1)

	p2 := productModel.Product{
		Name:        "test-product-2",
		Description: "test-product-2",
		Price:       2,
	}
	dbTest.Create(context.Background(), &p2)

	req := &dto.PlaceOrderReq{
		Lines: []dto.PlaceOrderLineReq{
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
}

func TestOrderAPI_PlaceOrderUnauthorized(t *testing.T) {
	defer cleanData()

	p1 := productModel.Product{
		Name:        "test-product-1",
		Description: "test-product-1",
		Price:       1,
	}
	dbTest.Create(context.Background(), &p1)

	p2 := productModel.Product{
		Name:        "test-product-2",
		Description: "test-product-2",
		Price:       2,
	}
	dbTest.Create(context.Background(), &p2)

	req := &dto.PlaceOrderReq{
		Lines: []dto.PlaceOrderLineReq{
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
}

// Get Order Detail
// =================================================================================================

func TestOrderAPI_GetOrderByIDSuccess(t *testing.T) {
	u := userModel.User{
		Email:    "test1@gmail.com",
		Password: "test123456",
	}
	dbTest.Create(context.Background(), &u)
	token := jtoken.GenerateAccessToken(map[string]interface{}{"id": u.ID})
	defer cleanData(&u)

	p1 := productModel.Product{
		Name:        "test-product-1",
		Description: "test-product-1",
		Price:       1,
	}
	dbTest.Create(context.Background(), &p1)

	p2 := productModel.Product{
		Name:        "test-product-2",
		Description: "test-product-2",
		Price:       2,
	}
	dbTest.Create(context.Background(), &p2)

	o := model.Order{
		UserID: u.ID,
		Lines: []*model.OrderLine{
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
	dbTest.Create(context.Background(), &o)

	writer := makeRequest("GET", fmt.Sprintf("/api/v1/orders/%s", o.ID), nil, token)
	var res dto.Order
	parseResponseResult(writer.Body.Bytes(), &res)
	assert.Equal(t, http.StatusOK, writer.Code)
	assert.Equal(t, "new", res.Status)
	assert.Equal(t, 2, len(res.Lines))
	assert.Equal(t, o.Lines[0].ProductID, res.Lines[0].Product.ID)
	assert.Equal(t, o.Lines[0].Quantity, res.Lines[0].Quantity)

	assert.Equal(t, o.Lines[1].ProductID, res.Lines[1].Product.ID)
	assert.Equal(t, o.Lines[1].Quantity, res.Lines[1].Quantity)
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
	u := userModel.User{
		Email:    "test1@test.com",
		Password: "test123456",
	}
	dbTest.Create(context.Background(), &u)
	token := jtoken.GenerateAccessToken(map[string]interface{}{"id": u.ID})
	defer cleanData(&u)

	p1 := productModel.Product{
		Name:        "test-product-1",
		Description: "test-product-1",
		Price:       1,
	}
	dbTest.Create(context.Background(), &p1)

	p2 := productModel.Product{
		Name:        "test-product-2",
		Description: "test-product-2",
		Price:       2,
	}
	dbTest.Create(context.Background(), &p2)

	o := model.Order{
		UserID: u.ID,
		Lines: []*model.OrderLine{
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
	dbTest.Create(context.Background(), &o)

	writer := makeRequest("PUT", fmt.Sprintf("/api/v1/orders/%s/cancel", o.ID), nil, token)
	assert.Equal(t, http.StatusOK, writer.Code)
}

func TestOrderAPI_CancelOrderNotFound(t *testing.T) {
	writer := makeRequest("PUT", "/api/v1/orders/notfound/cancel", nil, accessToken())
	var response map[string]map[string]string
	_ = json.Unmarshal(writer.Body.Bytes(), &response)
	assert.Equal(t, http.StatusInternalServerError, writer.Code)
	assert.Equal(t, "Something went wrong", response["error"]["message"])
}

func TestOrderAPI_CancelOrderStatusDone(t *testing.T) {
	u := userModel.User{
		Email:    "test1@test.com",
		Password: "test123456",
	}
	dbTest.Create(context.Background(), &u)
	token := jtoken.GenerateAccessToken(map[string]interface{}{"id": u.ID})
	defer cleanData(&u)

	p1 := productModel.Product{
		Name:        "test-product-1",
		Description: "test-product-1",
		Price:       1,
	}
	dbTest.Create(context.Background(), &p1)

	p2 := productModel.Product{
		Name:        "test-product-2",
		Description: "test-product-2",
		Price:       2,
	}
	dbTest.Create(context.Background(), &p2)

	o := model.Order{
		UserID: u.ID,
		Lines: []*model.OrderLine{
			{
				ProductID: p1.ID,
				Quantity:  2,
			},
			{
				ProductID: p2.ID,
				Quantity:  3,
			},
		},
		Status: model.OrderStatusDone,
	}
	dbTest.Create(context.Background(), &o)

	writer := makeRequest("PUT", fmt.Sprintf("/api/v1/orders/%s/cancel", o.ID), nil, token)
	var response map[string]map[string]string
	_ = json.Unmarshal(writer.Body.Bytes(), &response)
	assert.Equal(t, http.StatusInternalServerError, writer.Code)
	assert.Equal(t, "Something went wrong", response["error"]["message"])
}

func TestOrderAPI_CancelOrderStatusCancelled(t *testing.T) {
	u := userModel.User{
		Email:    "test1@test.com",
		Password: "test123456",
	}
	dbTest.Create(context.Background(), &u)
	token := jtoken.GenerateAccessToken(map[string]interface{}{"id": u.ID})
	defer cleanData(&u)

	p1 := productModel.Product{
		Name:        "test-product-1",
		Description: "test-product-1",
		Price:       1,
	}
	dbTest.Create(context.Background(), &p1)

	p2 := productModel.Product{
		Name:        "test-product-2",
		Description: "test-product-2",
		Price:       2,
	}
	dbTest.Create(context.Background(), &p2)

	o := model.Order{
		UserID: u.ID,
		Lines: []*model.OrderLine{
			{
				ProductID: p1.ID,
				Quantity:  2,
			},
			{
				ProductID: p2.ID,
				Quantity:  3,
			},
		},
		Status: model.OrderStatusCancelled,
	}
	dbTest.Create(context.Background(), &o)

	writer := makeRequest("PUT", fmt.Sprintf("/api/v1/orders/%s/cancel", o.ID), nil, token)
	var response map[string]map[string]string
	_ = json.Unmarshal(writer.Body.Bytes(), &response)
	assert.Equal(t, http.StatusInternalServerError, writer.Code)
	assert.Equal(t, "Something went wrong", response["error"]["message"])
}

func TestOrderAPI_CancelOrderNotMine(t *testing.T) {
	u := userModel.User{
		Email:    "test1@test.com",
		Password: "test123456",
	}
	dbTest.Create(context.Background(), &u)
	defer cleanData(&u)

	p1 := productModel.Product{
		Name:        "test-product-1",
		Description: "test-product-1",
		Price:       1,
	}
	dbTest.Create(context.Background(), &p1)

	p2 := productModel.Product{
		Name:        "test-product-2",
		Description: "test-product-2",
		Price:       2,
	}
	dbTest.Create(context.Background(), &p2)

	o := model.Order{
		UserID: u.ID,
		Lines: []*model.OrderLine{
			{
				ProductID: p1.ID,
				Quantity:  2,
			},
			{
				ProductID: p2.ID,
				Quantity:  3,
			},
		},
		Status: model.OrderStatusNew,
	}
	dbTest.Create(context.Background(), &o)

	writer := makeRequest("PUT", fmt.Sprintf("/api/v1/orders/%s/cancel", o.ID), nil, accessToken())
	var response map[string]map[string]string
	_ = json.Unmarshal(writer.Body.Bytes(), &response)
	assert.Equal(t, http.StatusInternalServerError, writer.Code)
	assert.Equal(t, "Something went wrong", response["error"]["message"])
}

// List My Orders
// =================================================================================================

func TestOrderAPI_ListOrdersSuccess(t *testing.T) {
	u := userModel.User{
		Email:    "test1@test.com",
		Password: "test123456",
	}
	dbTest.Create(context.Background(), &u)
	token := jtoken.GenerateAccessToken(map[string]interface{}{"id": u.ID})
	defer cleanData(&u)

	p1 := productModel.Product{
		Name:        "test-product-1",
		Description: "test-product-1",
		Price:       1,
	}
	dbTest.Create(context.Background(), &p1)

	p2 := productModel.Product{
		Name:        "test-product-2",
		Description: "test-product-2",
		Price:       2,
	}
	dbTest.Create(context.Background(), &p2)

	o1 := model.Order{
		UserID: u.ID,
		Lines: []*model.OrderLine{
			{
				ProductID: p1.ID,
				Quantity:  5,
			},
		},
		Status: model.OrderStatusDone,
	}
	dbTest.Create(context.Background(), &o1)

	o2 := model.Order{
		UserID: u.ID,
		Lines: []*model.OrderLine{
			{
				ProductID: p2.ID,
				Quantity:  3,
			},
		},
		Status: model.OrderStatusNew,
	}
	dbTest.Create(context.Background(), &o2)

	writer := makeRequest("GET", "/api/v1/orders", nil, token)
	var res dto.ListOrderRes
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
}

func TestOrderAPI_ListOrdersNotFound(t *testing.T) {
	defer cleanData()

	writer := makeRequest("GET", "/api/v1/orders", nil, accessToken())
	var res dto.ListOrderRes
	parseResponseResult(writer.Body.Bytes(), &res)
	assert.Equal(t, http.StatusOK, writer.Code)
	assert.Equal(t, 0, len(res.Orders))
}

func TestOrderAPI_ListOrdersInvalidFieldType(t *testing.T) {
	writer := makeRequest("GET", "/api/v1/orders?page=a", nil, accessToken())
	var response map[string]map[string]string
	_ = json.Unmarshal(writer.Body.Bytes(), &response)
	assert.Equal(t, http.StatusBadRequest, writer.Code)
	assert.Equal(t, "Invalid parameters", response["error"]["message"])
}

func TestOrderAPI_ListMyOrdersFindByStatusSuccess(t *testing.T) {
	u := userModel.User{
		Email:    "test1@test.com",
		Password: "test123456",
	}
	dbTest.Create(context.Background(), &u)
	token := jtoken.GenerateAccessToken(map[string]interface{}{"id": u.ID})
	defer cleanData(&u)

	p1 := productModel.Product{
		Name:        "test-product-1",
		Description: "test-product-1",
		Price:       1,
	}
	dbTest.Create(context.Background(), &p1)

	p2 := productModel.Product{
		Name:        "test-product-2",
		Description: "test-product-2",
		Price:       2,
	}
	dbTest.Create(context.Background(), &p2)

	o1 := model.Order{
		UserID: u.ID,
		Lines: []*model.OrderLine{
			{
				ProductID: p1.ID,
				Quantity:  5,
			},
		},
		Status: model.OrderStatusDone,
	}
	dbTest.Create(context.Background(), &o1)

	o2 := model.Order{
		UserID: u.ID,
		Lines: []*model.OrderLine{
			{
				ProductID: p2.ID,
				Quantity:  3,
			},
		},
		Status: model.OrderStatusNew,
	}
	dbTest.Create(context.Background(), &o2)

	writer := makeRequest("GET", "/api/v1/orders?status=new", nil, token)
	var res dto.ListOrderRes
	parseResponseResult(writer.Body.Bytes(), &res)
	assert.Equal(t, http.StatusOK, writer.Code)
	assert.Equal(t, int64(1), res.Pagination.Total)
	assert.Equal(t, int64(1), res.Pagination.CurrentPage)
	assert.Equal(t, int64(1), res.Pagination.TotalPage)
	assert.Equal(t, int64(20), res.Pagination.Limit)
	assert.Equal(t, 1, len(res.Orders))
	assert.Equal(t, o2.Lines[0].ProductID, res.Orders[0].Lines[0].Product.ID)
	assert.Equal(t, o2.Lines[0].Quantity, res.Orders[0].Lines[0].Quantity)
}

func TestOrderAPI_ListOrdersFindByStatusNotFound(t *testing.T) {
	u := userModel.User{
		Email:    "test1@test.com",
		Password: "test123456",
	}
	dbTest.Create(context.Background(), &u)
	token := jtoken.GenerateAccessToken(map[string]interface{}{"id": u.ID})
	defer cleanData(&u)

	p1 := productModel.Product{
		Name:        "test-product-1",
		Description: "test-product-1",
		Price:       1,
	}
	dbTest.Create(context.Background(), &p1)

	p2 := productModel.Product{
		Name:        "test-product-2",
		Description: "test-product-2",
		Price:       2,
	}
	dbTest.Create(context.Background(), &p2)

	o1 := model.Order{
		UserID: u.ID,
		Lines: []*model.OrderLine{
			{
				ProductID: p1.ID,
				Quantity:  5,
			},
		},
		Status: model.OrderStatusDone,
	}
	dbTest.Create(context.Background(), &o1)

	o2 := model.Order{
		UserID: u.ID,
		Lines: []*model.OrderLine{
			{
				ProductID: p2.ID,
				Quantity:  3,
			},
		},
		Status: model.OrderStatusNew,
	}
	dbTest.Create(context.Background(), &o2)

	writer := makeRequest("GET", "/api/v1/orders?status=cancelled", nil, token)
	var res dto.ListOrderRes
	parseResponseResult(writer.Body.Bytes(), &res)
	assert.Equal(t, http.StatusOK, writer.Code)
	assert.Equal(t, 0, len(res.Orders))
}

func TestOrderAPI_ListOrdersFindByCodeSuccess(t *testing.T) {
	u := userModel.User{
		Email:    "test1@test.com",
		Password: "test123456",
	}
	dbTest.Create(context.Background(), &u)
	token := jtoken.GenerateAccessToken(map[string]interface{}{"id": u.ID})
	defer cleanData(&u)

	p1 := productModel.Product{
		Name:        "test-product-1",
		Description: "test-product-1",
		Price:       1,
	}
	dbTest.Create(context.Background(), &p1)

	p2 := productModel.Product{
		Name:        "test-product-2",
		Description: "test-product-2",
		Price:       2,
	}
	dbTest.Create(context.Background(), &p2)

	o1 := model.Order{
		UserID: u.ID,
		Lines: []*model.OrderLine{
			{
				ProductID: p1.ID,
				Quantity:  5,
			},
		},
		Status: model.OrderStatusDone,
	}
	dbTest.Create(context.Background(), &o1)

	o2 := model.Order{
		UserID: u.ID,
		Lines: []*model.OrderLine{
			{
				ProductID: p2.ID,
				Quantity:  3,
			},
		},
		Status: model.OrderStatusNew,
	}
	dbTest.Create(context.Background(), &o2)

	writer := makeRequest("GET", fmt.Sprintf("/api/v1/orders?code=%s", o1.Code), nil, token)
	var res dto.ListOrderRes
	parseResponseResult(writer.Body.Bytes(), &res)
	assert.Equal(t, http.StatusOK, writer.Code)
	assert.Equal(t, int64(1), res.Pagination.Total)
	assert.Equal(t, int64(1), res.Pagination.CurrentPage)
	assert.Equal(t, int64(1), res.Pagination.TotalPage)
	assert.Equal(t, int64(20), res.Pagination.Limit)
	assert.Equal(t, 1, len(res.Orders))
	assert.Equal(t, o1.Lines[0].ProductID, res.Orders[0].Lines[0].Product.ID)
	assert.Equal(t, o1.Lines[0].Quantity, res.Orders[0].Lines[0].Quantity)
}

func TestOrderAPI_ListOrdersFindByCodeNotFound(t *testing.T) {
	u := userModel.User{
		Email:    "test1@test.com",
		Password: "test123456",
	}
	dbTest.Create(context.Background(), &u)
	token := jtoken.GenerateAccessToken(map[string]interface{}{"id": u.ID})
	defer cleanData(&u)

	p1 := productModel.Product{
		Name:        "test-product-1",
		Description: "test-product-1",
		Price:       1,
	}
	dbTest.Create(context.Background(), &p1)

	p2 := productModel.Product{
		Name:        "test-product-2",
		Description: "test-product-2",
		Price:       2,
	}
	dbTest.Create(context.Background(), &p2)

	o1 := model.Order{
		UserID: u.ID,
		Lines: []*model.OrderLine{
			{
				ProductID: p1.ID,
				Quantity:  5,
			},
		},
		Status: model.OrderStatusDone,
	}
	dbTest.Create(context.Background(), &o1)

	o2 := model.Order{
		UserID: u.ID,
		Lines: []*model.OrderLine{
			{
				ProductID: p2.ID,
				Quantity:  3,
			},
		},
		Status: model.OrderStatusNew,
	}
	dbTest.Create(context.Background(), &o2)

	writer := makeRequest("GET", "/api/v1/orders?code=notfound", nil, token)
	var res dto.ListOrderRes
	parseResponseResult(writer.Body.Bytes(), &res)
	assert.Equal(t, http.StatusOK, writer.Code)
	assert.Equal(t, 0, len(res.Orders))
}

func TestOrderAPI_ListOrdersWithPagination(t *testing.T) {
	u := userModel.User{
		Email:    "test1@test.com",
		Password: "test123456",
	}
	dbTest.Create(context.Background(), &u)
	token := jtoken.GenerateAccessToken(map[string]interface{}{"id": u.ID})
	defer cleanData(&u)

	p1 := productModel.Product{
		Name:        "test-product-1",
		Description: "test-product-1",
		Price:       1,
	}
	dbTest.Create(context.Background(), &p1)

	p2 := productModel.Product{
		Name:        "test-product-2",
		Description: "test-product-2",
		Price:       2,
	}
	dbTest.Create(context.Background(), &p2)

	o1 := model.Order{
		UserID: u.ID,
		Lines: []*model.OrderLine{
			{
				ProductID: p1.ID,
				Quantity:  5,
			},
		},
		Status: model.OrderStatusDone,
	}
	dbTest.Create(context.Background(), &o1)

	o2 := model.Order{
		UserID: u.ID,
		Lines: []*model.OrderLine{
			{
				ProductID: p2.ID,
				Quantity:  3,
			},
		},
		Status: model.OrderStatusNew,
	}
	dbTest.Create(context.Background(), &o2)

	writer := makeRequest("GET", "/api/v1/orders?page=2&limit=1", nil, token)
	var res dto.ListOrderRes
	parseResponseResult(writer.Body.Bytes(), &res)
	assert.Equal(t, http.StatusOK, writer.Code)
	assert.Equal(t, int64(2), res.Pagination.Total)
	assert.Equal(t, int64(2), res.Pagination.CurrentPage)
	assert.Equal(t, int64(2), res.Pagination.TotalPage)
	assert.Equal(t, int64(1), res.Pagination.Limit)
	assert.Equal(t, 1, len(res.Orders))
	assert.Equal(t, o2.Lines[0].ProductID, res.Orders[0].Lines[0].Product.ID)
	assert.Equal(t, o2.Lines[0].Quantity, res.Orders[0].Lines[0].Quantity)
}

func TestOrderAPI_ListOrdersWithOrder(t *testing.T) {
	u := userModel.User{
		Email:    "test1@test.com",
		Password: "test123456",
	}
	dbTest.Create(context.Background(), &u)
	token := jtoken.GenerateAccessToken(map[string]interface{}{"id": u.ID})
	defer cleanData(&u)

	p1 := productModel.Product{
		Name:        "test-product-1",
		Description: "test-product-1",
		Price:       1,
	}
	dbTest.Create(context.Background(), &p1)

	p2 := productModel.Product{
		Name:        "test-product-2",
		Description: "test-product-2",
		Price:       2,
	}
	dbTest.Create(context.Background(), &p2)

	o1 := model.Order{
		UserID: u.ID,
		Lines: []*model.OrderLine{
			{
				ProductID: p1.ID,
				Quantity:  5,
			},
		},
		Status: model.OrderStatusDone,
	}
	dbTest.Create(context.Background(), &o1)

	o2 := model.Order{
		UserID: u.ID,
		Lines: []*model.OrderLine{
			{
				ProductID: p2.ID,
				Quantity:  3,
			},
		},
		Status: model.OrderStatusNew,
	}
	dbTest.Create(context.Background(), &o2)

	writer := makeRequest("GET", "/api/v1/orders?order_by=created_at&order_desc=true", nil, token)
	var res dto.ListOrderRes
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
}

func TestOrderAPI_GetMyOrdersNotMine(t *testing.T) {
	u := userModel.User{
		Email:    "test1@test.com",
		Password: "test123456",
	}
	dbTest.Create(context.Background(), &u)
	defer cleanData(&u)

	p1 := productModel.Product{
		Name:        "test-product-1",
		Description: "test-product-1",
		Price:       1,
	}
	dbTest.Create(context.Background(), &p1)

	p2 := productModel.Product{
		Name:        "test-product-2",
		Description: "test-product-2",
		Price:       2,
	}
	dbTest.Create(context.Background(), &p2)

	o1 := model.Order{
		UserID: u.ID,
		Lines: []*model.OrderLine{
			{
				ProductID: p1.ID,
				Quantity:  5,
			},
		},
		Status: model.OrderStatusDone,
	}
	dbTest.Create(context.Background(), &o1)

	o2 := model.Order{
		UserID: u.ID,
		Lines: []*model.OrderLine{
			{
				ProductID: p2.ID,
				Quantity:  3,
			},
		},
		Status: model.OrderStatusNew,
	}
	dbTest.Create(context.Background(), &o2)

	writer := makeRequest("GET", "/api/v1/orders?code=notfound", nil, accessToken())
	var res dto.ListOrderRes
	parseResponseResult(writer.Body.Bytes(), &res)
	assert.Equal(t, http.StatusOK, writer.Code)
	assert.Equal(t, 0, len(res.Orders))
}
