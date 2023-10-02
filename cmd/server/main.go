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
	router.Route("/", func(r chi.Router) {
		r.Handle(shared.UpdateEndpoint, handler.NewUpdateHandler(storage, parser.NewUpdateRequestParser()))
		r.Handle(shared.ValueEndpoint, handler.NewValueHandler(storage, parser.NewValueRequestParser()))
	})

	err := http.ListenAndServe(hostport, router)
	if err != nil {
		panic(err)
	}
}
