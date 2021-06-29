package shresiesbot

import (
	"context"
	"errors"
	"log"
	"testing"

	"github.com/deividfortuna/sharesies"
	"github.com/stretchr/testify/assert"

	"github.com/deividfortuna/sharesies-bot/mocks"
)

var mockConfiguration = &AutoInvest{
	Buy: &BuyConfiguration{
		Orders: mockOrders,
	},
	Sharesies: &Credentials{
		Username: "username",
		Password: "password",
	},
}

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

	sharesiesMock.On("Authenticate", ctx, &sharesies.Credentials{Username: "username", Password: "password"}).Return(nil, nil)
	sharesiesMock.On("CostBuy", ctx, mockOrder.Id, mockOrder.Amount).Return(costBuyResponseMock, nil)
	sharesiesMock.On("Buy", ctx, costBuyResponseMock).Return(nil, nil)

	bot := &SharesiesBot{
		scheduler: &mocks.Scheduler{},
		client:    sharesiesMock,
		config:    mockConfiguration,
		logger:    log.Default(),
	}

	err := bot.buyOrders(ctx, mockOrders)

	assert.Nil(t, err)
	sharesiesMock.AssertExpectations(t)
}

func Test_Buy_Fail_Auth(t *testing.T) {
	ctx := context.Background()
	errAuthFailed := errors.New("auth_failed")

	sharesiesMock := &mocks.ExchangeClient{}
	sharesiesMock.On("Authenticate", ctx, &sharesies.Credentials{Username: "username", Password: "password"}).Return(nil, errAuthFailed)

	bot := &SharesiesBot{
		scheduler: &mocks.Scheduler{},
		client:    sharesiesMock,
		config:    mockConfiguration,
		logger:    log.Default(),
	}

	err := bot.buyOrders(ctx, mockOrders)

	assert.Error(t, err, errAuthFailed)
	sharesiesMock.AssertExpectations(t)
}

func Test_Buy_Fail_Get_Price_Generic(t *testing.T) {
	ctx := context.Background()
	errGeneric := errors.New("generic_error")

	sharesiesMock := &mocks.ExchangeClient{}

	sharesiesMock.On("Authenticate", ctx, &sharesies.Credentials{Username: "username", Password: "password"}).Return(nil, nil)
	sharesiesMock.On("CostBuy", ctx, mockOrder.Id, mockOrder.Amount).Return(nil, errGeneric)

	bot := &SharesiesBot{
		scheduler: &mocks.Scheduler{},
		client:    sharesiesMock,
		config:    mockConfiguration,
		logger:    log.Default(),
	}

	err := bot.buyOrders(ctx, mockOrders)

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

	sharesiesMock.On("Authenticate", ctx, &sharesies.Credentials{Username: "username", Password: "password"}).Return(nil, nil)
	sharesiesMock.On("CostBuy", ctx, mockOrder.Id, mockOrder.Amount).Return(costBuyResponseMock, nil)

	bot := &SharesiesBot{
		scheduler: &mocks.Scheduler{},
		client:    sharesiesMock,
		config:    mockConfiguration,
		logger:    log.Default(),
	}

	err := bot.buyOrders(ctx, mockOrders)

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

	sharesiesMock.On("Authenticate", ctx, &sharesies.Credentials{Username: "username", Password: "password"}).Return(nil, nil)
	sharesiesMock.On("CostBuy", ctx, mockOrder.Id, mockOrder.Amount).Return(costBuyResponseMock, nil)
	sharesiesMock.On("Buy", ctx, costBuyResponseMock).Return(nil, errBuy)

	bot := &SharesiesBot{
		scheduler: &mocks.Scheduler{},
		client:    sharesiesMock,
		config:    mockConfiguration,
		logger:    log.Default(),
	}

	err := bot.buyOrders(ctx, mockOrders)

	assert.ErrorIs(t, err, errBuy)
	sharesiesMock.AssertExpectations(t)
}
