package shresiesbot

import (
	"context"
	"log"
	"os"

	"github.com/deividfortuna/sharesies"
	"github.com/robfig/cron/v3"
)

type Scheduler interface {
	AddFunc(spec string, cmd func()) (cron.EntryID, error)
}

type ExchangeClient interface {
	Authenticate(ctx context.Context, creds *sharesies.Credentials) (*sharesies.ProfileResponse, error)
	CostBuy(ctx context.Context, fundId string, amount float64) (*sharesies.CostBuyResponse, error)
	Buy(ctx context.Context, costBuy *sharesies.CostBuyResponse) (*sharesies.ProfileResponse, error)
}

var logger = log.New(os.Stdout, "SharesiesBot: ", log.LstdFlags)

func Run(cr Scheduler, s ExchangeClient, conf *AutoInvest) error {
	ctx := context.Background()

	if conf.Buy != nil {
		_, err := cr.AddFunc(conf.Buy.Scheduler, func() {
			err := buyOrders(ctx, s, conf.Sharesies, conf.Buy.Orders)
			if err != nil {
				log.Fatal(err)
			}
		})

		if err != nil {
			return err
		}
	}

	return nil
}
