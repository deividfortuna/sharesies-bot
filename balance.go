package shresiesbot

import (
	"context"
	"strconv"

	"github.com/deividfortuna/sharesies"
)

type holdPrice struct {
	Value    float64
	Quantity float64
}

func (b *SharesiesBot) balance(ctx context.Context, holds []Hold) error {
	p, err := b.client.Authenticate(ctx, &sharesies.Credentials{
		Username: b.config.Sharesies.Username,
		Password: b.config.Sharesies.Password,
	})
	if err != nil {
		return err
	}

	portfolioHolds := map[string]holdPrice{}
	portfolioValue := 0.00

	for _, h := range holds {
		for _, v := range p.Portfolio {
			if h.Id == v.FundID {
				value, err := strconv.ParseFloat(v.Value, 64)
				quantity, err := strconv.ParseFloat(v.Shares, 64)
				if err != nil {
					return err
				}
				portfolioHolds[v.FundID] = holdPrice{Value: value, Quantity: quantity}
				portfolioValue = portfolioValue + value
			}
		}
	}

	var bOrders []BuyOrder
	var sOrders []SellOrder

	for _, h := range holds {
		v := (h.Weight * portfolioValue) / 100

		if v > portfolioHolds[h.Id].Value {
			bOrders = append(bOrders, BuyOrder{
				Id:        h.Id,
				Reference: h.Reference,
				Amount:    v - portfolioHolds[h.Id].Value,
			})
		}

		if v < portfolioHolds[h.Id].Value {
			difference := portfolioHolds[h.Id].Value - v
			sharePrice := portfolioHolds[h.Id].Value / portfolioHolds[h.Id].Quantity
			shareQuantity := difference / sharePrice

			sOrders = append(sOrders, SellOrder{
				Id:        h.Id,
				Reference: h.Reference,
				Shares:    shareQuantity,
			})
		}
	}

	b.logger.Println("buy orders: ", bOrders)
	err = b.buyOrders(ctx, bOrders)
	if err != nil {
		return err
	}

	err = b.sellOrders(ctx, sOrders)
	if err != nil {
		return err
	}

	return nil
}
