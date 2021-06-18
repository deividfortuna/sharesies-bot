package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/deividfortuna/sharesies"
	"github.com/robfig/cron/v3"
)

var logger = log.New(os.Stdout, "SharesiesBot: ", log.LstdFlags)

func main() {
	ctx := context.Background()

	cr := cron.New(
		cron.WithLogger(
			cron.VerbosePrintfLogger(log.New(os.Stdout, "Scheduler: ", log.LstdFlags))))

	a := &AutoInvest{}
	c := &Credentials{}

	err := load("config/credentials.yml", c)
	checkErr(err)
	err = load("config/auto_invest.yml", a)
	checkErr(err)

	s, err := sharesies.New(nil)
	checkErr(err)

	cr.AddFunc(a.Scheduler, func() {
		err = buyOrders(ctx, s, c, a.Buy)
		checkErr(err)
	})

	cr.Start()

	channel := make(chan os.Signal, 1)
	signal.Notify(channel, os.Interrupt, syscall.SIGTERM)

	<-channel
	logger.Println("Gracefully shutting down...")
	logger.Println("Running cleanup tasks...")

	cr.Stop()

	logger.Println("Successful shutdown.")
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
