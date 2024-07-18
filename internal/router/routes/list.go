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
	List(ctx context.Context, pk string) ([]datastore.Record, uint64, error)
}

type listResp struct {
	Cursor uint64             `json:"cursor"`
	Data   []datastore.Record `json:"data"`
}

func List(l Lister) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("listing data")

		recs, curs, err := l.List(r.Context(), r.PathValue(RouteParamPK))
		if err != nil {
			if errors.Is(err, datastore.ErrNotFound) {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			log.Println("error listing records", err)
			http.Error(w, "failed to list records", http.StatusInternalServerError)
			return
		}
		resp := listResp{
			Cursor: curs,
			Data:   recs,
		}

		w.Header().Set("Content-Type", "application/json")
		if err = json.NewEncoder(w).Encode(resp); err != nil {
			log.Println("encoder error", err)
			http.Error(w, "unknown error occurred", http.StatusInternalServerError)
			return
		}
	}
}
