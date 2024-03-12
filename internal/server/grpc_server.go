package server

import (
	"context"
	"fmt"
	"net"

	"github.com/kuzhukin/metrics-collector/internal/metric"
	pb "github.com/kuzhukin/metrics-collector/internal/proto"
	"github.com/kuzhukin/metrics-collector/internal/server/config"
	"github.com/kuzhukin/metrics-collector/internal/server/storage"
	"google.golang.org/grpc"
)

var _ pb.MetricsServiceServer = &GRPCMetricServer{}

type GRPCMetricServer struct {
	pb.UnimplementedMetricsServiceServer
	storage storage.Storage
	server  *grpc.Server
}

func NewGrpcServer(storage storage.Storage, config *config.Config) (*GRPCMetricServer, error) {
	listen, err := net.Listen("tcp", config.Hostport)
	if err != nil {
		return nil, fmt.Errorf("grpc listen err=%w", err)
	}

	s := grpc.NewServer()

	grpcMetricServer := &GRPCMetricServer{storage: storage, server: s}

	pb.RegisterMetricsServiceServer(s, grpcMetricServer)

	if err := s.Serve(listen); err != nil {
		return nil, fmt.Errorf("serve err=%w", err)
	}

	return grpcMetricServer, nil
}

func (s *GRPCMetricServer) Stop() {
	s.server.Stop()
}

func (s *GRPCMetricServer) BatchUpdate(ctx context.Context, req *pb.BatchUpdateRequest) (*pb.BatchUpdateResponse, error) {
	metrics := make([]*metric.Metric, 0, len(req.Metric))

	for _, m := range req.Metric {
		metrics = append(metrics, &metric.Metric{ID: m.Id, Delta: &m.Delta, Value: &m.Value, Type: metric.Kind(m.Type)})
	}

	if err := s.storage.BatchUpdate(ctx, metrics); err != nil {
		return nil, fmt.Errorf("batch update err %s", err)
	}

	return &pb.BatchUpdateResponse{}, nil
}
