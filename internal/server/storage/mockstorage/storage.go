// Code generated by mockery v2.36.0. DO NOT EDIT.

package mockstorage

import (
	metric "github.com/kuzhukin/metrics-collector/internal/metric"
	mock "github.com/stretchr/testify/mock"
)

// Storage is an autogenerated mock type for the Storage type
type Storage struct {
	mock.Mock
}

// BatchUpdate provides a mock function with given fields: m
func (_m *Storage) BatchUpdate(m []*metric.Metric) error {
	ret := _m.Called(m)

	var r0 error
	if rf, ok := ret.Get(0).(func([]*metric.Metric) error); ok {
		r0 = rf(m)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Get provides a mock function with given fields: kind, name
func (_m *Storage) Get(kind metric.Kind, name string) (*metric.Metric, error) {
	ret := _m.Called(kind, name)

	var r0 *metric.Metric
	var r1 error
	if rf, ok := ret.Get(0).(func(metric.Kind, string) (*metric.Metric, error)); ok {
		return rf(kind, name)
	}
	if rf, ok := ret.Get(0).(func(metric.Kind, string) *metric.Metric); ok {
		r0 = rf(kind, name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*metric.Metric)
		}
	}

	if rf, ok := ret.Get(1).(func(metric.Kind, string) error); ok {
		r1 = rf(kind, name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// List provides a mock function with given fields:
func (_m *Storage) List() ([]*metric.Metric, error) {
	ret := _m.Called()

	var r0 []*metric.Metric
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]*metric.Metric, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []*metric.Metric); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*metric.Metric)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: m
func (_m *Storage) Update(m *metric.Metric) error {
	ret := _m.Called(m)

	var r0 error
	if rf, ok := ret.Get(0).(func(*metric.Metric) error); ok {
		r0 = rf(m)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewStorage creates a new instance of Storage. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewStorage(t interface {
	mock.TestingT
	Cleanup(func())
}) *Storage {
	mock := &Storage{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
