package shresiesbot

import (
	"context"
	"strconv"

	"github.com/deividfortuna/sharesies"
)

func (b *SharesiesBot) balance(ctx context.Context, holds []Hold) error {
	p, err := b.client.Authenticate(ctx, &sharesies.Credentials{
		Username: b.config.Sharesies.Username,
		Password: b.config.Sharesies.Password,
	})
	if err != nil {
		return err
	}

	portfolioHolds := map[string]float64{}
	portfolioValue := 0.00

	for _, h := range holds {
		for _, v := range p.Portfolio {
			if h.Id == v.FundID {
				value, err := strconv.ParseFloat(v.Value, 64)
				if err != nil {
					return err
				}
				portfolioHolds[v.FundID] = value
				portfolioValue = portfolioValue + value
			}
		}
	}

	var bOrders []BuyOrder

	for _, h := range holds {
		v := (h.Weight * portfolioValue) / 100

		if v > portfolioHolds[h.Id] {
			bOrders = append(bOrders, BuyOrder{
				Id:        h.Id,
				Reference: h.Reference,
				Amount:    v - portfolioHolds[h.Id],
			})
		}

		if v < portfolioHolds[h.Id] {
			// TODO: Implement sell proportional shares
			b.logger.Println("sell", portfolioHolds[h.Id]-v)
		}
	}

	b.logger.Println("buy orders: ", bOrders)
	return b.buyOrders(ctx, bOrders)
}
