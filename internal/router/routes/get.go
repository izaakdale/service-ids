package routes

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/izaakdale/service-ids/internal/datastore"
)

type Fetcher interface {
	Fetch(ctx context.Context, keys datastore.Keys) (*datastore.IDRecord, error)
}

func Get(f Fetcher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("getting data")

		rec, err := f.Fetch(r.Context(), datastore.Keys{
			PK: r.PathValue(RouteParamPK),
			SK: r.PathValue(RouteParamSK),
		})
		if err != nil {
			if errors.Is(err, datastore.ErrNotFound) {
				log.Println("not found during populate", err)
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			log.Println("other error", err)
			http.Error(w, "unknown error occurred", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		log.Println("transmitting record")
		if err = json.NewEncoder(w).Encode(rec); err != nil {
			http.Error(w, "unknown error occurred", http.StatusInternalServerError)
			return
		}
	}
}
