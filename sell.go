package shresiesbot

import (
	"context"
	"errors"
	"github.com/deividfortuna/sharesies"
)

var orderSellBuyType = "order_cost_sell"

type SellOrder struct {
	Id        string
	Reference string
	Shares    float64
}

func (b *SharesiesBot) sellOrders(ctx context.Context, orders []SellOrder) error {
	_, err := b.client.Authenticate(ctx, &sharesies.Credentials{
		Username: b.config.Sharesies.Username,
		Password: b.config.Sharesies.Password,
	})
	if err != nil {
		b.logger.Println("Failed to authenticated Sharesies")
		return err
	}

	var sellOrders []*sharesies.CostSellResponse
	for _, v := range orders {
		b.logger.Println("Checking order price for " + v.Reference)
		cs, err := b.client.CostSell(ctx, v.Id, v.Shares)
		if err != nil {
			return err
		}

		if cs.Type != orderSellBuyType {
			return errors.New(cs.Type)
		}

		sellOrders = append(sellOrders, cs)
	}

	for _, so := range sellOrders {
		_, err = b.client.Sell(ctx, so)
		if err != nil {
			return err
		}
	}

	b.logger.Println("Stonks! We sold everything")
	return nil
}
