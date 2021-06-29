package shresiesbot

import (
	"context"
	"errors"

	"github.com/deividfortuna/sharesies"
)

var orderCostBuyType = "order_cost_buy"

func (b *SharesiesBot) buyOrders(ctx context.Context, orders []BuyOrder) error {
	_, err := b.client.Authenticate(ctx, &sharesies.Credentials{
		Username: b.config.Sharesies.Username,
		Password: b.config.Sharesies.Password,
	})
	if err != nil {
		b.logger.Println("Failed to authenticated Sharesies")
		return err
	}

	var costOrders []*sharesies.CostBuyResponse
	for _, v := range orders {
		b.logger.Println("Checking order price for " + v.Reference)
		cb, err := b.client.CostBuy(ctx, v.Id, v.Amount)
		if err != nil {
			return err
		}

		if cb.Type != orderCostBuyType {
			return errors.New(cb.Type)
		}

		costOrders = append(costOrders, cb)
	}

	for _, co := range costOrders {
		_, err = b.client.Buy(ctx, co)
		if err != nil {
			return err
		}
	}

	b.logger.Println("Stonks! We bought everything")
	return nil
}
