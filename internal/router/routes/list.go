package routes

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/izaakdale/service-ids/internal/datastore"
)

type Lister interface {
	List(ctx context.Context, pk string) ([]datastore.Record, error)
}

func List(l Lister) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("listing data")
		recs, err := l.List(r.Context(), r.PathValue(RouteParamPK))
		if err != nil {
			if errors.Is(err, datastore.ErrNotFound) {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			log.Println("error listing records", err)
			http.Error(w, "failed to list records", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err = json.NewEncoder(w).Encode(recs); err != nil {
			log.Println("encoder error", err)
			http.Error(w, "unknown error occurred", http.StatusInternalServerError)
			return
		}
	}
}
