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

func GetID(f Fetcher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("getting ID")

		rec, err := f.Fetch(r.Context(), datastore.Keys{
			PK: r.PathValue("id"),
			SK: r.PathValue("type"),
		})
		if err != nil {
			if errors.Is(err, datastore.ErrNotFound) {
				log.Println("not found during populate")
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			log.Println("other error")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		log.Println("transmitting record")
		if err = json.NewEncoder(w).Encode(rec); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

type Inserter interface {
	Insert(ctx context.Context, rec datastore.IDRecord) error
}

type postBody struct {
	ID string `json:"type_id"`
}

func PostID(i Inserter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("posting ID")
		var pb postBody
		if err := json.NewDecoder(r.Body).Decode(&pb); err != nil {
			log.Println("decode error")
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}
		rec := datastore.IDRecord{
			Keys: datastore.Keys{
				PK: r.PathValue("id"),
				SK: r.PathValue("type"),
			},
			ID: pb.ID,
		}
		if err := i.Insert(r.Context(), rec); err != nil {
			log.Println("error storing record")
			http.Error(w, "failed to store ID record", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(rec); err != nil {
			log.Println("encoder error")
			http.Error(w, "unknown error occurred", http.StatusInternalServerError)
			return
		}
	}
}
