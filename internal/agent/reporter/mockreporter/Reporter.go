// Code generated by mockery v2.35.2. DO NOT EDIT.

package mockreporter

import mock "github.com/stretchr/testify/mock"

// Reporter is an autogenerated mock type for the Reporter type
type Reporter struct {
	mock.Mock
}

// Report provides a mock function with given fields: gaugeMetrics, counterMetrics
func (_m *Reporter) Report(gaugeMetrics map[string]float64, counterMetrics map[string]int64) {
	_m.Called(gaugeMetrics, counterMetrics)
}

// NewReporter creates a new instance of Reporter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewReporter(t interface {
	mock.TestingT
	Cleanup(func())
}) *Reporter {
	mock := &Reporter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
