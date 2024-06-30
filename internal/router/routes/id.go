package routes

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/izaakdale/service-ids/internal/datastore"
)

var (
	RouteParamPK = "pk"
	RouteParamSK = "sk"
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
	Data string `json:"data"`
}

func Post(i Inserter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("posting data")
		var pb postBody
		if err := json.NewDecoder(r.Body).Decode(&pb); err != nil {
			log.Println("decode error")
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

type Lister interface {
	List(ctx context.Context, pk string) ([]datastore.IDRecord, error)
}

func List(l Lister) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("listing data")
		recs, err := l.List(r.Context(), r.PathValue(RouteParamPK))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err = json.NewEncoder(w).Encode(recs); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
