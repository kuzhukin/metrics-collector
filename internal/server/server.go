package server

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kuzhukin/metrics-collector/internal/server/endpoint"
	"github.com/kuzhukin/metrics-collector/internal/server/handler"
	"github.com/kuzhukin/metrics-collector/internal/server/parser"
	"github.com/kuzhukin/metrics-collector/internal/server/storage/memorystorage"
)

const hostport = ":8080"

func Run() error {
	storage := memorystorage.New()

	router := chi.NewRouter()

	listHandler := handler.NewGetListHandler(storage)
	updateHandler := handler.NewUpdateHandler(storage, parser.NewUpdateRequestParser())
	valueHandler := handler.NewValueHandler(storage, parser.NewValueRequestParser())

	router.Handle(endpoint.RootEndpoint, listHandler)
	router.Handle(endpoint.UpdateEndpoint, updateHandler)
	router.Handle(endpoint.ValueEndpoint, valueHandler)

	err := http.ListenAndServe(hostport, router)
	if err != nil {
		return fmt.Errorf("http listen and server, hostport=%s, err=%w", hostport, err)
	}

	return nil
}
