// Code generated by mockery v2.36.0. DO NOT EDIT.

package mockparser

import (
	http "net/http"

	metric "github.com/kuzhukin/metrics-collector/internal/metric"
	mock "github.com/stretchr/testify/mock"
)

// RequestParser is an autogenerated mock type for the RequestParser type
type RequestParser struct {
	mock.Mock
}

// Parse provides a mock function with given fields: r
func (_m *RequestParser) Parse(r *http.Request) (*metric.Metric, error) {
	ret := _m.Called(r)

	var r0 *metric.Metric
	var r1 error
	if rf, ok := ret.Get(0).(func(*http.Request) (*metric.Metric, error)); ok {
		return rf(r)
	}
	if rf, ok := ret.Get(0).(func(*http.Request) *metric.Metric); ok {
		r0 = rf(r)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*metric.Metric)
		}
	}

	if rf, ok := ret.Get(1).(func(*http.Request) error); ok {
		r1 = rf(r)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewRequestParser creates a new instance of RequestParser. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewRequestParser(t interface {
	mock.TestingT
	Cleanup(func())
}) *RequestParser {
	mock := &RequestParser{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
