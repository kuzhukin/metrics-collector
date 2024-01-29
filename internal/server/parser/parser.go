// module declares interfaces for parsing single requests - RequestParser
// and interface for parsing batch requests - BatchRequestParser
package parser

import (
	"errors"
	"net/http"

	"github.com/kuzhukin/metrics-collector/internal/metric"
)

var ErrMetricNameIsNotFound error = errors.New("metric name isn't found")
var ErrBadMetricKind error = errors.New("bad metric kind")

//go:generate mockery --name=RequestParser --filename=parser.go --outpkg=mockparser --output=mockparser
type RequestParser interface {
	Parse(r *http.Request) (*metric.Metric, error)
}

//go:generate mockery --name=BatchRequestParser --filename=batch_parser.go --outpkg=mockparser --output=mockparser
type BatchRequestParser interface {
	BatchParse(r *http.Request) ([]*metric.Metric, error)
}
