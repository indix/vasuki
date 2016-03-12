package scalar

import "github.com/stretchr/testify/mock"

import "github.com/ashwanthkumar/go-gocd"

import "github.com/ind9/vasuki/executor"

type MockScalar struct {
	mock.Mock
}

// config provides a mock function with given fields:
func (_m *MockScalar) config() *Config {
	ret := _m.Called()

	var r0 *Config
	if rf, ok := ret.Get(0).(func() *Config); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*Config)
		}
	}

	return r0
}

// client provides a mock function with given fields:
func (_m *MockScalar) client() gocd.Client {
	ret := _m.Called()

	var r0 gocd.Client
	if rf, ok := ret.Get(0).(func() gocd.Client); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(gocd.Client)
	}

	return r0
}

// Demand provides a mock function with given fields:
func (_m *MockScalar) Demand() (int, error) {
	ret := _m.Called()

	var r0 int
	if rf, ok := ret.Get(0).(func() int); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Supply provides a mock function with given fields:
func (_m *MockScalar) Supply() (int, error) {
	ret := _m.Called()

	var r0 int
	if rf, ok := ret.Get(0).(func() int); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Executor provides a mock function with given fields:
func (_m *MockScalar) Executor() executor.Executor {
	ret := _m.Called()

	var r0 executor.Executor
	if rf, ok := ret.Get(0).(func() executor.Executor); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(executor.Executor)
	}

	return r0
}

// IdleAgents provides a mock function with given fields:
func (_m *MockScalar) IdleAgents() ([]string, error) {
	ret := _m.Called()

	var r0 []string
	if rf, ok := ret.Get(0).(func() []string); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ComputeScaleUp provides a mock function with given fields: demand, supply
func (_m *MockScalar) ComputeScaleUp(demand int, supply int) (int, error) {
	ret := _m.Called(demand, supply)

	var r0 int
	if rf, ok := ret.Get(0).(func(int, int) int); ok {
		r0 = rf(demand, supply)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int, int) error); ok {
		r1 = rf(demand, supply)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ComputeScaleDown provides a mock function with given fields: demand, supply, idleAgents
func (_m *MockScalar) ComputeScaleDown(demand int, supply int, idleAgents int) (int, error) {
	ret := _m.Called(demand, supply, idleAgents)

	var r0 int
	if rf, ok := ret.Get(0).(func(int, int, int) int); ok {
		r0 = rf(demand, supply, idleAgents)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int, int, int) error); ok {
		r1 = rf(demand, supply, idleAgents)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
