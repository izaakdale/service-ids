package routes

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/izaakdale/service-ids/internal/datastore"
)

type Inserter interface {
	Insert(ctx context.Context, rec datastore.IDRecord) error
}

type postBody struct {
	Data string `json:"data"`
}

func Post(i Inserter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("posting data")
		var pb postBody
		if err := json.NewDecoder(r.Body).Decode(&pb); err != nil {
			log.Println("decode error", err)
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}
		rec := datastore.IDRecord{
			Keys: datastore.Keys{
				PK: r.PathValue(RouteParamPK),
				SK: r.PathValue(RouteParamSK),
			},
			ID: pb.Data,
		}
		if err := i.Insert(r.Context(), rec); err != nil {
			log.Println("error storing record", err)
			http.Error(w, "failed to store ID record", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(rec); err != nil {
			log.Println("encoder error", err)
			http.Error(w, "unknown error occurred", http.StatusInternalServerError)
			return
		}
	}
}
