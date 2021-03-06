// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	cron "github.com/robfig/cron/v3"
	mock "github.com/stretchr/testify/mock"
)

// Scheduler is an autogenerated mock type for the Scheduler type
type Scheduler struct {
	mock.Mock
}

// AddFunc provides a mock function with given fields: spec, cmd
func (_m *Scheduler) AddFunc(spec string, cmd func()) (cron.EntryID, error) {
	ret := _m.Called(spec, cmd)

	var r0 cron.EntryID
	if rf, ok := ret.Get(0).(func(string, func()) cron.EntryID); ok {
		r0 = rf(spec, cmd)
	} else {
		r0 = ret.Get(0).(cron.EntryID)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, func()) error); ok {
		r1 = rf(spec, cmd)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
