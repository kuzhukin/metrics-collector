package parser

import (
	"errors"
	"net/http"

	"github.com/kuzhukin/metrics-collector/internal/server/metric"
)

var ErrMetricNameIsNotFound error = errors.New("metric name isn't found")
var ErrBadMetricKind error = errors.New("bad metric kind")

//go:generate mockery --name=RequestParser --outpkg=mockparser --output=mockparser
type RequestParser interface {
	Parse(r *http.Request) (*metric.Metric, error)
}
