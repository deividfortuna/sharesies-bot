package shresiesbot

import (
	"context"
	"errors"
	"testing"

	"github.com/deividfortuna/sharesies"
	"github.com/stretchr/testify/assert"

	"github.com/deividfortuna/auto-invest-sharesies/mocks"
)

var mockOrder = &BuyOrder{
	Reference: "reference",
	Id:        "mock_id",
	Amount:    100.00,
}

var mockOrders = []BuyOrder{
	*mockOrder,
}

func Test_Buy(t *testing.T) {
	ctx := context.Background()

	sharesiesMock := &mocks.ExchangeClient{}
	costBuyResponseMock := &sharesies.CostBuyResponse{
		FundID:    mockOrder.Id,
		TotalCost: "100.00",
		Type:      "order_cost_buy",
	}

	sharesiesMock.On("Authenticate", ctx, &sharesies.Credentials{Username: "username", Password: "passowrd"}).Return(nil, nil)
	sharesiesMock.On("CostBuy", ctx, mockOrder.Id, mockOrder.Amount).Return(costBuyResponseMock, nil)
	sharesiesMock.On("Buy", ctx, costBuyResponseMock).Return(nil, nil)

	err := buyOrders(ctx, sharesiesMock, &Credentials{
		Username: "username",
		Password: "passowrd",
	}, mockOrders)

	assert.Nil(t, err)
	sharesiesMock.AssertExpectations(t)
}

func Test_Buy_Fail_Auth(t *testing.T) {
	ctx := context.Background()
	errAuthFailed := errors.New("auth_failed")

	sharesiesMock := &mocks.ExchangeClient{}
	sharesiesMock.On("Authenticate", ctx, &sharesies.Credentials{Username: "username", Password: "passowrd"}).Return(nil, errAuthFailed)

	err := buyOrders(ctx, sharesiesMock, &Credentials{
		Username: "username",
		Password: "passowrd",
	}, mockOrders)

	assert.Error(t, err, errAuthFailed)
	sharesiesMock.AssertExpectations(t)
}

func Test_Buy_Fail_Get_Price_Generic(t *testing.T) {
	ctx := context.Background()
	errGeneric := errors.New("generic_error")

	sharesiesMock := &mocks.ExchangeClient{}

	sharesiesMock.On("Authenticate", ctx, &sharesies.Credentials{Username: "username", Password: "passowrd"}).Return(nil, nil)
	sharesiesMock.On("CostBuy", ctx, mockOrder.Id, mockOrder.Amount).Return(nil, errGeneric)

	err := buyOrders(ctx, sharesiesMock, &Credentials{
		Username: "username",
		Password: "passowrd",
	}, mockOrders)

	assert.ErrorIs(t, err, errGeneric)
	sharesiesMock.AssertExpectations(t)
}

func Test_Buy_Fail_Get_Price(t *testing.T) {
	ctx := context.Background()

	sharesiesMock := &mocks.ExchangeClient{}
	costBuyResponseMock := &sharesies.CostBuyResponse{
		FundID:    mockOrder.Id,
		TotalCost: "100.00",
		Type:      "failed",
	}

	sharesiesMock.On("Authenticate", ctx, &sharesies.Credentials{Username: "username", Password: "passowrd"}).Return(nil, nil)
	sharesiesMock.On("CostBuy", ctx, mockOrder.Id, mockOrder.Amount).Return(costBuyResponseMock, nil)

	err := buyOrders(ctx, sharesiesMock, &Credentials{
		Username: "username",
		Password: "passowrd",
	}, mockOrders)

	assert.Error(t, err)
	sharesiesMock.AssertExpectations(t)
}

func Test_Buy_Fail_Buy(t *testing.T) {
	ctx := context.Background()
	errBuy := errors.New("fail_buying")

	sharesiesMock := &mocks.ExchangeClient{}
	costBuyResponseMock := &sharesies.CostBuyResponse{
		FundID:    mockOrder.Id,
		TotalCost: "100.00",
		Type:      "order_cost_buy",
	}

	sharesiesMock.On("Authenticate", ctx, &sharesies.Credentials{Username: "username", Password: "passowrd"}).Return(nil, nil)
	sharesiesMock.On("CostBuy", ctx, mockOrder.Id, mockOrder.Amount).Return(costBuyResponseMock, nil)
	sharesiesMock.On("Buy", ctx, costBuyResponseMock).Return(nil, errBuy)

	err := buyOrders(ctx, sharesiesMock, &Credentials{
		Username: "username",
		Password: "passowrd",
	}, mockOrders)

	assert.ErrorIs(t, err, errBuy)
	sharesiesMock.AssertExpectations(t)
}
