package router

import (
	"net/http"

	"github.com/izaakdale/ittp"
	"github.com/izaakdale/service-ids/internal/router/routes"
)

type Datastore interface {
	routes.Fetcher
	routes.Inserter
}

func New(d Datastore) http.Handler {
	mux := ittp.NewServeMux()
	mux.Get("/_/live", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	mux.Get("/_/ready", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	mux.Get("/id/{id}/{type}", routes.GetID(d))
	mux.Post("/id/{id}/{type}", routes.PostID(d))

	return mux
}
