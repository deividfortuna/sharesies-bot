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
	CostSell(ctx context.Context, foundId string, shareAmount float64) (*sharesies.CostSellResponse, error)
	Sell(ctx context.Context, sellBuy *sharesies.CostSellResponse) (*sharesies.ProfileResponse, error)
}

type SharesiesBot struct {
	scheduler Scheduler
	client    ExchangeClient
	config    *AutoInvest
	logger    *log.Logger
}

func New(scheduler Scheduler, client ExchangeClient, config *AutoInvest) *SharesiesBot {
	var logger = log.New(os.Stdout, "SharesiesBot: ", log.LstdFlags)

	return &SharesiesBot{
		scheduler,
		client,
		config,
		logger,
	}
}

func (b *SharesiesBot) Run() error {
	ctx := context.Background()

	if b.config.Buy != nil {
		_, err := b.scheduler.AddFunc(b.config.Buy.Scheduler, func() {
			err := b.buyOrders(ctx, b.config.Buy.Orders)
			if err != nil {
				log.Fatal(err)
			}
		})

		if err != nil {
			return err
		}
	}

	if b.config.Balance != nil {
		_, err := b.scheduler.AddFunc(b.config.Balance.Scheduler, func() {
			err := b.balance(ctx, b.config.Balance.Holds)
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
