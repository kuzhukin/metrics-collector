package server

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kuzhukin/metrics-collector/internal/log"
	"github.com/kuzhukin/metrics-collector/internal/server/endpoint"
	"github.com/kuzhukin/metrics-collector/internal/server/handler"
	"github.com/kuzhukin/metrics-collector/internal/server/parser"
	"github.com/kuzhukin/metrics-collector/internal/server/storage/memorystorage"
)

type MetricServer struct {
	srvr http.Server
	wait chan struct{}
}

func StartNew(config Config) *MetricServer {
	storage := memorystorage.New()

	router := chi.NewRouter()
	listHandler := handler.NewGetListHandler(storage)
	updateHandler := handler.NewUpdateHandler(storage, parser.NewUpdateRequestParser())
	valueHandler := handler.NewValueHandler(storage, parser.NewValueRequestParser())

	router.Handle(endpoint.RootEndpoint, listHandler)
	router.Handle(endpoint.UpdateEndpointJSON, updateHandler)
	router.Handle(endpoint.UpdateEndpoint, updateHandler)
	router.Handle(endpoint.ValueEndpoint, valueHandler)

	server := &MetricServer{
		srvr: http.Server{
			Addr:    config.Hostport,
			Handler: log.LoggingHTTPHandler(router),
		},
		wait: make(chan struct{}),
	}

	go func() {
		defer close(server.wait)

		if err := server.srvr.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Logger.Errorf("Http listen and serve, address=%s, err=%s\n", server.srvr.Addr, err)
		}
	}()

	log.Logger.Infof("Server started hostport=%v", config.Hostport)

	return server
}

func (s *MetricServer) Stop() error {
	log.Logger.Infof("Server stopped")
	return s.srvr.Shutdown(context.Background())
}

func (s *MetricServer) Wait() <-chan struct{} {
	return s.wait
}
