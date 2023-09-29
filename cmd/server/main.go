package main

import (
	"net/http"

	"github.com/kuzhukin/metrics-collector/cmd/server/handler"
	"github.com/kuzhukin/metrics-collector/cmd/server/shared"
	"github.com/kuzhukin/metrics-collector/cmd/server/storage/memorystorage"
)

func main() {
	mux := http.NewServeMux()
	registerHandlers(mux)

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}

func registerHandlers(mux *http.ServeMux) {
	storage := memorystorage.New()
	mux.Handle(shared.UpdateEndpoint, handler.NewUpdateHandler(storage))
}
