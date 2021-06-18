package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/deividfortuna/sharesies"
	"github.com/go-playground/validator/v10"
	"github.com/robfig/cron/v3"
	"gopkg.in/yaml.v2"

	autoinvest "github.com/deividfortuna/auto-invest-sharesies"
)

func main() {
	scheduler := cron.New(
		cron.WithLogger(
			cron.VerbosePrintfLogger(log.New(os.Stdout, "Scheduler: ", log.LstdFlags))))

	exchange, err := sharesies.New(nil)
	if err != nil {
		log.Fatal(err)
	}

	config, err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}

	err = autoinvest.Run(scheduler, exchange, config)
	if err != nil {
		log.Fatal(err)
	}

	scheduler.Start()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c
	log.Println("Gracefully shutting down...")
	log.Println("Running cleanup tasks...")

	scheduler.Stop()

	log.Println("Successful shutdown.")
}

func loadConfig() (*autoinvest.AutoInvest, error) {
	v := &autoinvest.AutoInvest{}
	f, err := ioutil.ReadFile("config/auto_invest.yml")
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(f, v)
	if err != nil {
		return nil, err
	}

	validate := validator.New()
	err = validate.Struct(v)
	if err != nil {
		return nil, err
	}

	return v, nil
}
