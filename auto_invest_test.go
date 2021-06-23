package shresiesbot_test

import (
	"errors"
	"testing"

	"github.com/robfig/cron/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	shresiesbot "github.com/deividfortuna/sharesies-bot"
	"github.com/deividfortuna/sharesies-bot/mocks"
)

func Test_Run(t *testing.T) {
	mockScheduler := mocks.Scheduler{}
	mockExchage := mocks.ExchangeClient{}

	mockConfig := &shresiesbot.AutoInvest{
		Buy: &shresiesbot.BuyConfiguration{
			Scheduler: "MY_SCHEDULER",
			Orders:    []shresiesbot.BuyOrder{},
		},
	}

	entryId := new(cron.EntryID)
	mockScheduler.On("AddFunc", "MY_SCHEDULER", mock.Anything).Return(*entryId, nil)

	bot := shresiesbot.New(&mockScheduler, &mockExchage, mockConfig)
	err := bot.Run()

	assert.Nil(t, err)
	mockScheduler.AssertExpectations(t)
}

func Test_Run_Buy_Fail(t *testing.T) {
	mockScheduler := mocks.Scheduler{}
	mockExchage := mocks.ExchangeClient{}

	errFailAddFunc := errors.New("fail_add_func")

	mockConfig := &shresiesbot.AutoInvest{
		Buy: &shresiesbot.BuyConfiguration{
			Scheduler: "MY_SCHEDULER",
			Orders:    []shresiesbot.BuyOrder{},
		},
	}

	entryId := new(cron.EntryID)
	mockScheduler.On("AddFunc", "MY_SCHEDULER", mock.Anything).Return(*entryId, errFailAddFunc)

	bot := shresiesbot.New(&mockScheduler, &mockExchage, mockConfig)
	err := bot.Run()

	assert.ErrorIs(t, err, errFailAddFunc)
	mockScheduler.AssertExpectations(t)
}
