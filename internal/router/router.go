package router

import (
	"fmt"
	"net/http"

	"github.com/izaakdale/ittp"
	"github.com/izaakdale/service-ids/internal/router/routes"
)

type Datastore interface {
	routes.Fetcher
	routes.Inserter
	routes.Lister
}

func New(d Datastore) http.Handler {
	mux := ittp.NewServeMux()
	mux.Get("/_/live", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	mux.Get("/_/ready", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	mux.Get(fmt.Sprintf("/{%s}/{%s}", routes.RouteParamPK, routes.RouteParamSK), routes.Get(d))
	mux.Post(fmt.Sprintf("/{%s}/{%s}", routes.RouteParamPK, routes.RouteParamSK), routes.Post(d))
	mux.Get(fmt.Sprintf("/{%s}", routes.RouteParamPK), routes.List(d))

	return mux
}
