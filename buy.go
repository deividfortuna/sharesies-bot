package main

import (
	"context"
	"errors"

	"github.com/deividfortuna/sharesies"
)

var orderCostBuyType = "order_cost_buy"

func buyOrders(ctx context.Context, creds *Credentials, orders []BuyOrder) error {
	s, err := sharesies.New(nil)
	if err != nil {
		return err
	}

	_, err = s.Authenticate(ctx, &sharesies.Credentials{
		Username: creds.Username,
		Password: creds.Password,
	})
	if err != nil {
		logger.Println("Failed to authenticated Sharesies")
		return err
	}

	costOrders := []*sharesies.CostBuyResponse{}
	for _, v := range orders {
		logger.Println("Checking order price for " + v.Reference)
		cb, err := s.CostBuy(ctx, v.Id, v.Amount)
		if err != nil {
			return err
		}

		if cb.Type != orderCostBuyType {
			return errors.New(cb.Type)
		}

		costOrders = append(costOrders, cb)
	}

	for _, co := range costOrders {
		_, err = s.Buy(ctx, co)
		if err != nil {
			return err
		}
	}

	logger.Printf("Stonks! We bought everything")
	return nil
}
