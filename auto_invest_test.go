package shresiesbot_test

import (
	"errors"
	"testing"

	"github.com/robfig/cron/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	autoinvest "github.com/deividfortuna/auto-invest-sharesies"
	"github.com/deividfortuna/auto-invest-sharesies/mocks"
)

func Test_Run(t *testing.T) {
	mockScheduler := mocks.Scheduler{}
	mockExchage := mocks.ExchangeClient{}

	mockConfig := &autoinvest.AutoInvest{
		Buy: &autoinvest.BuyConfiguration{
			Scheduler: "MY_SCHEDULER",
			Orders:    []autoinvest.BuyOrder{},
		},
	}

	entryId := new(cron.EntryID)
	mockScheduler.On("AddFunc", "MY_SCHEDULER", mock.Anything).Return(*entryId, nil)

	err := autoinvest.Run(&mockScheduler, &mockExchage, mockConfig)

	assert.Nil(t, err)
	mockScheduler.AssertExpectations(t)
}

func Test_Run_Buy_Fail(t *testing.T) {
	mockScheduler := mocks.Scheduler{}
	mockExchage := mocks.ExchangeClient{}

	errFailAddFunc := errors.New("fail_add_func")

	mockConfig := &autoinvest.AutoInvest{
		Buy: &autoinvest.BuyConfiguration{
			Scheduler: "MY_SCHEDULER",
			Orders:    []autoinvest.BuyOrder{},
		},
	}

	entryId := new(cron.EntryID)
	mockScheduler.On("AddFunc", "MY_SCHEDULER", mock.Anything).Return(*entryId, errFailAddFunc)

	err := autoinvest.Run(&mockScheduler, &mockExchage, mockConfig)

	assert.ErrorIs(t, err, errFailAddFunc)
	mockScheduler.AssertExpectations(t)
}
