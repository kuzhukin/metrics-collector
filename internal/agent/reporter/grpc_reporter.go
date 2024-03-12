package reporter

import (
	"context"
	"time"

	"github.com/kuzhukin/metrics-collector/internal/agent/config"
	pb "github.com/kuzhukin/metrics-collector/internal/proto"
	"github.com/kuzhukin/metrics-collector/internal/zlog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type grpcReporterImpl struct {
	hostport string
}

func newGRPCReporter(config config.Config) (*grpcReporterImpl, error) {
	return &grpcReporterImpl{
		hostport: config.Hostport,
	}, nil
}

func (r *grpcReporterImpl) Report(gaugeMetrics map[string]float64, counterMetrics map[string]int64) {
	metrics := preparePbMetric(gaugeMetrics, counterMetrics)

	conn, err := grpc.Dial(":3200", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		zlog.Logger.Errorf("grpc dial err %s", err)

		return
	}
	defer conn.Close()

	c := pb.NewMetricsServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	resp, err := c.BatchUpdate(ctx, &pb.BatchUpdateRequest{Metric: metrics})
	if err != nil {
		zlog.Logger.Errorf("batch update err %s", err)
	}

	if resp.Error != "" {
		zlog.Logger.Errorf("batch update upload err %s", err)
	}

}
func preparePbMetric(gaugeMetrics map[string]float64, counterMetrics map[string]int64) []*pb.Metric {
	metrics := make([]*pb.Metric, 0, len(gaugeMetrics)+len(counterMetrics))

	for name, gauge := range gaugeMetrics {
		metrics = append(metrics, &pb.Metric{Id: name, Type: "gauge", Value: gauge})
	}

	for name, counter := range counterMetrics {
		metrics = append(metrics, &pb.Metric{Id: name, Type: "counter", Delta: counter})
	}

	return metrics
}
