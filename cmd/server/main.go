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
	storage := memorystorage.New()

	router := chi.NewRouter()

	listHandler := handler.NewGetListHandler(storage)
	updateHandler := handler.NewUpdateHandler(storage, parser.NewUpdateRequestParser())
	valueHandler := handler.NewValueHandler(storage, parser.NewValueRequestParser())

	router.Handle(shared.RootEndpoint, listHandler)
	router.Handle(shared.UpdateEndpoint, updateHandler)
	router.Handle(shared.ValueEndpoint, valueHandler)

	err := http.ListenAndServe(hostport, router)
	if err != nil {
		panic(err)
	}
}
