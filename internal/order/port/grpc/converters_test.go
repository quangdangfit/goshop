package grpc

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"goshop/internal/order/model"
)

func TestOrderInfoFromModel(t *testing.T) {
	assert.Nil(t, orderInfoFromModel(nil))
	got := orderInfoFromModel(&model.Order{
		ID: "o1", Code: "C", UserID: "u", TotalPrice: 12.5, Status: model.OrderStatusNew,
		Lines: []*model.OrderLine{
			{ProductID: "p", Quantity: 2, Price: 4, Product: &model.Product{Name: "n"}},
			{ProductID: "p2", Quantity: 1, Price: 1},
		},
	})
	assert.Equal(t, "o1", got.Id)
	assert.Equal(t, float32(12.5), got.TotalPrice)
	assert.Len(t, got.Lines, 2)
	assert.Equal(t, "n", got.Lines[0].ProductName)
	assert.Equal(t, "", got.Lines[1].ProductName)
}

func TestOrdersInfoFromModel(t *testing.T) {
	out := ordersInfoFromModel([]*model.Order{{ID: "a"}, nil})
	assert.Len(t, out, 2)
	assert.Equal(t, "a", out[0].Id)
	assert.Nil(t, out[1])
}
