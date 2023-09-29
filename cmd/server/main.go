package main

import (
	"net/http"

	"github.com/kuzhukin/metrics-collector/cmd/server/handler"
	"github.com/kuzhukin/metrics-collector/internal/shared"
	"github.com/kuzhukin/metrics-collector/internal/storage/memorystorage"
)

const hostport = ":8080"

func main() {
	mux := http.NewServeMux()
	registerHandlers(mux)

	err := http.ListenAndServe(hostport, mux)
	if err != nil {
		panic(err)
	}
}

func registerHandlers(mux *http.ServeMux) {
	storage := memorystorage.New()
	mux.Handle(shared.UpdateEndpoint, handler.NewUpdateHandler(storage))
}
