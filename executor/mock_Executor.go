package executor

import "github.com/stretchr/testify/mock"

type MockExecutor struct {
	mock.Mock
}

// Init provides a mock function with given fields: config
func (_m *MockExecutor) Init(config *Config) error {
	ret := _m.Called(config)

	var r0 error
	if rf, ok := ret.Get(0).(func(*Config) error); ok {
		r0 = rf(config)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ScaleUp provides a mock function with given fields: instances
func (_m *MockExecutor) ScaleUp(instances int) error {
	ret := _m.Called(instances)

	var r0 error
	if rf, ok := ret.Get(0).(func(int) error); ok {
		r0 = rf(instances)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ScaleDown provides a mock function with given fields: agentsToKill
func (_m *MockExecutor) ScaleDown(agentsToKill []string) error {
	ret := _m.Called(agentsToKill)

	var r0 error
	if rf, ok := ret.Get(0).(func([]string) error); ok {
		r0 = rf(agentsToKill)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ManagedAgents provides a mock function with given fields:
func (_m *MockExecutor) ManagedAgents() ([]string, error) {
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
