// package server - Module implements HTTP server for collecting metrics
package server

import (
	"context"
	"fmt"
	"net/http"
	"net/netip"

	"github.com/go-chi/chi/v5"
	"github.com/kuzhukin/metrics-collector/internal/server/config"
	"github.com/kuzhukin/metrics-collector/internal/server/endpoint"
	"github.com/kuzhukin/metrics-collector/internal/server/handler"
	"github.com/kuzhukin/metrics-collector/internal/server/handler/middleware"
	"github.com/kuzhukin/metrics-collector/internal/server/parser"
	"github.com/kuzhukin/metrics-collector/internal/server/storage"
	"github.com/kuzhukin/metrics-collector/internal/server/storage/dbstorage"
	"github.com/kuzhukin/metrics-collector/internal/server/storage/filestorage"
	"github.com/kuzhukin/metrics-collector/internal/server/storage/memorystorage"
	"github.com/kuzhukin/metrics-collector/internal/zlog"
)

type MetricServer struct {
	// HTTP server and router
	srvr http.Server
	// GRPC server
	grpc *GRPCMetricServer
	// channel for waiting of server shutdown
	wait chan struct{}
}

// StartNew - creates and starts HTTP server
func StartNew() (*MetricServer, error) {
	config, err := config.MakeConfig()
	if err != nil {
		return nil, fmt.Errorf("make config, err=%w", err)
	}

	server, err := createServer(&config)
	if err != nil {
		return nil, fmt.Errorf("create server, err=%w", err)
	}

	server.startHTTPServer()

	zlog.Logger.Infof("Server started config=%+v", config)

	return server, nil
}

func createServer(config *config.Config) (*MetricServer, error) {
	var err error
	var storage storage.Storage
	var dbStorage *dbstorage.DBStorage

	if config.Storage.DatabaseDSN != "" {
		dbStorage, err = dbstorage.StartNew(config.Storage.DatabaseDSN)
		if err != nil {
			return nil, fmt.Errorf("new db storage, err=%w", err)
		}

		storage = dbStorage
	} else if config.Storage.FilePath != "" {
		storage, err = filestorage.New(config.Storage)
		if err != nil {
			return nil, fmt.Errorf("new file storage, err=%w", err)
		}
	} else {
		storage = memorystorage.New()
	}

	requestsParser := parser.New()

	listHandler := handler.NewGetListHandler(storage)
	updateHandler := handler.NewUpdateHandler(storage, requestsParser)
	valueHandler := handler.NewValueHandler(storage, requestsParser)
	pingHandler := handler.NewPingHandler(dbStorage)
	batchUpdateHandler := handler.NewBatchUpdateHandler(storage, requestsParser)

	router := chi.NewRouter()

	if config.TrustedSubnet != "" {
		trustedPrefix := netip.MustParsePrefix(config.TrustedSubnet)
		checker := middleware.NewIPChecker(trustedPrefix)

		router.Use(checker.CheckIPHandler)
	}

	if config.EnableLogger {
		router.Use(middleware.LoggingHTTPHandler)
	}

	if config.CryptoKey != "" {
		decrypter, err := middleware.NewDecryptHTTPHandler(config.CryptoKey)
		if err != nil {
			return nil, fmt.Errorf("new decrypt handler, err=%w", err)
		}

		router.Use(decrypter)
	}

	if config.SingnatureKey != "" {
		middleware.InitSignHandlers(config.SingnatureKey)
		router.Use(middleware.SignCheckHandler)
		router.Use(middleware.SignCreateHandler)
	}

	router.Use(middleware.CompressingHTTPHandler)

	router.Handle(endpoint.RootEndpoint, listHandler)
	router.Handle(endpoint.UpdateEndpoint, updateHandler)
	router.Handle(endpoint.UpdateEndpointJSON, updateHandler)
	router.Handle(endpoint.ValueEndpoint, valueHandler)
	router.Handle(endpoint.ValueEndpointJSON, valueHandler)
	router.Handle(endpoint.PingEndpoint, pingHandler)
	router.Handle(endpoint.BatchUpdateEndpointJSON, batchUpdateHandler)

	metricServer := &MetricServer{
		srvr: http.Server{
			Addr:    config.Hostport,
			Handler: router,
		},
		wait: make(chan struct{}),
	}

	if config.UseGRPC {
		grpcServer, err := NewGrpcServer(storage, config)
		if err != nil {
			return nil, fmt.Errorf("new grpc server err %w", err)
		}

		metricServer.grpc = grpcServer

	}

	return metricServer, nil
}

func (s *MetricServer) startHTTPServer() {
	go func() {
		defer close(s.wait)

		if err := s.srvr.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zlog.Logger.Errorf("Http listen and serve, address=%s, err=%s\n", s.srvr.Addr, err)
		}
	}()
}

// Stop server
func (s *MetricServer) Stop() error {
	zlog.Logger.Infof("Server stopped")

	if s.grpc != nil {
		s.grpc.Stop()
	}

	return s.srvr.Shutdown(context.Background())
}

// Wait shutdown the server
func (s *MetricServer) Wait() <-chan struct{} {
	return s.wait
}
