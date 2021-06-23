package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/deividfortuna/sharesies"
	"github.com/go-playground/validator/v10"
	"github.com/robfig/cron/v3"
	"gopkg.in/yaml.v2"

	shresiesbot "github.com/deividfortuna/sharesies-bot"
)

const defaultTimeZone = "Pacific/Auckland"

func main() {
	tz := os.Getenv("TZ")
	if tz == "" {
		tz = defaultTimeZone
	}

	timezone, _ := time.LoadLocation(tz)

	scheduler := cron.New(
		cron.WithLocation(timezone),
		cron.WithLogger(
			cron.VerbosePrintfLogger(log.New(os.Stdout, "Scheduler: ", log.LstdFlags))))

	exchange, err := sharesies.New(nil)
	if err != nil {
		log.Fatal(err)
	}

	filepath := os.Getenv("CONFIG_FILE")
	if filepath == "" {
		filepath = "config/auto_invest.yml"
	}
	config, err := loadConfig(filepath)
	if err != nil {
		log.Fatal(err)
	}

	bot := shresiesbot.New(scheduler, exchange, config)
	err = bot.Run()
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

func loadConfig(filePath string) (*shresiesbot.AutoInvest, error) {
	v := &shresiesbot.AutoInvest{}
	f, err := ioutil.ReadFile(filePath)
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
