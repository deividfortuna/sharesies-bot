package shresiesbot

import (
	"context"
	"errors"
	"github.com/stretchr/testify/mock"
	"log"
	"testing"

	"github.com/deividfortuna/sharesies"
	"github.com/deividfortuna/sharesies-bot/mocks"
	"github.com/stretchr/testify/assert"
)

func Test_Balance(t *testing.T) {
	ctx := context.Background()

	holds := []Hold{
		{Id: "HOLD_1", Weight: 50},
		{Id: "HOLD_2", Weight: 50},
	}

	costBuyMock := &sharesies.CostBuyResponse{
		FundID:    "HOLD_2",
		TotalCost: "13.79",
		Type:      "order_cost_buy",
	}

	costSellMock := &sharesies.CostSellResponse{
		FundID: "HOLD_1",
		Request: &sharesies.OrderSell{
			Type:        "share_market",
			ShareAmount: "0.2161780843392381",
		},
		Type: "order_cost_sell",
	}

	sharesiesMock := &mocks.ExchangeClient{}
	sharesiesMock.On("Authenticate", ctx, &sharesies.Credentials{Username: "username", Password: "password"}).Return(&sharesies.ProfileResponse{
		Portfolio: []*sharesies.Portfolio{
			{FundID: "HOLD_1", Value: "63.79", Shares: "1"},
			{FundID: "HOLD_2", Value: "36.21", Shares: "1"},
		},
	}, nil)
	sharesiesMock.On("CostBuy", ctx, "HOLD_2", 13.79).Once().Return(costBuyMock, nil)
	sharesiesMock.On("Buy", ctx, costBuyMock).Once().Return(nil, nil)
	sharesiesMock.On("CostSell", ctx, "HOLD_1", 0.2161780843392381).Once().Return(costSellMock, nil)
	sharesiesMock.On("Sell", ctx, mock.Anything).Once().Return(nil, nil)

	bot := &SharesiesBot{
		scheduler: &mocks.Scheduler{},
		client:    sharesiesMock,
		config: &AutoInvest{
			Sharesies: &Credentials{
				Username: "username",
				Password: "password",
			},
		},
		logger: log.Default(),
	}

	err := bot.balance(ctx, holds)

	assert.Nil(t, err)
	sharesiesMock.AssertExpectations(t)
}

func Test_Balance_Authenticate_Fail(t *testing.T) {
	ctx := context.Background()

	errFailAuthenticate := errors.New("authenticate_fail")

	sharesiesMock := &mocks.ExchangeClient{}
	sharesiesMock.On("Authenticate", ctx, &sharesies.Credentials{Username: "username", Password: "password"}).Return(nil, errFailAuthenticate)

	bot := &SharesiesBot{
		scheduler: &mocks.Scheduler{},
		client:    sharesiesMock,
		config: &AutoInvest{
			Sharesies: &Credentials{
				Username: "username",
				Password: "password",
			},
		},
		logger: log.Default(),
	}

	err := bot.balance(ctx, []Hold{})

	assert.ErrorIs(t, err, errFailAuthenticate)
	sharesiesMock.AssertExpectations(t)
}
