package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kuzhukin/metrics-collector/cmd/server/handler"
	"github.com/kuzhukin/metrics-collector/cmd/server/parser"
	"github.com/kuzhukin/metrics-collector/internal/shared"
	"github.com/kuzhukin/metrics-collector/internal/storage/memorystorage"
)

const hostport = ":8080"

func main() {
	router := chi.NewRouter()
	router.Handle(shared.UpdateEndpoint+"{kind}/{name}/{value}", createUpdateHander())

	err := http.ListenAndServe(hostport, router)
	if err != nil {
		panic(err)
	}
}

func createUpdateHander() *handler.UpdateHandler {
	storage := memorystorage.New()
	parser := parser.NewRequestParser()
	return handler.NewUpdateHandler(storage, parser)
}
