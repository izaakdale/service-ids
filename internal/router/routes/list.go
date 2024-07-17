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
	List(ctx context.Context, pk string, offset uint64, limit int64) ([]datastore.Record, uint64, error)
}

type listReqBody struct {
	Cursor uint64 `json:"cursor"`
	Count  int64  `json:"count"`
}

type listResp struct {
	Cursor uint64             `json:"cursor"`
	Data   []datastore.Record `json:"data"`
}

func List(l Lister) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("listing data")

		var lr listReqBody
		if err := json.NewDecoder(r.Body).Decode(&lr); err != nil {
			log.Println("error decoding request body", err)

			lr.Count = 10
			lr.Cursor = 0
		}

		log.Printf("%+v\n", lr)

		recs, curs, err := l.List(r.Context(), r.PathValue(RouteParamPK), lr.Cursor, lr.Count)
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
