package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/verminio/shopfam/shopping"
)

func UpsertItem(service *shopping.ItemService) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		type request struct {
			Id        *shopping.ItemId `json:"id,omitempty"`
			Name      string           `json:"name"`
			Quantity  string           `json:"quantity"`
			DateAdded int64            `json:"dateAdded"`
		}

		requestData := &request{}

		dec := json.NewDecoder(req.Body)
		if err := dec.Decode(requestData); err != nil {
			http.Error(w, "Invalid Request Body", http.StatusBadRequest)
			return
		}

		item := shopping.New(requestData.Name, requestData.Quantity, time.UnixMilli(requestData.DateAdded))

		id, err := service.Upsert(requestData.Id, item)
		if err != nil {
			log.Printf("Unexpected error: %s", err)
			http.Error(w, "Unexpected error", http.StatusInternalServerError)
			return
		}

		type response struct {
			Id shopping.ItemId `json:"id"`
		}
		w.Header().Set("Content-Type", "application/json")
		enc := json.NewEncoder(w)
		if err := enc.Encode(response{Id: *id}); err != nil {
			log.Printf("Unexpected error: %s", err)
			http.Error(w, "Unexpected error", http.StatusInternalServerError)
			return
		}
	})
}

func ListItems(service *shopping.ItemService) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		items, err := service.List()
		if err != nil {
			log.Printf("Unexpected error: %s", err)
			http.Error(w, "Unexpected error", http.StatusInternalServerError)
			return
		}

		type item struct {
			Id        shopping.ItemId `json:"id"`
			Name      string          `json:"name"`
			Quantity  string          `json:"quantity"`
			DateAdded int64           `json:"dateAdded"`
		}

		rsp := make([]item, len(items))

		for cur, i := range items {
			rsp[cur] = item{
				Id:        i.Id,
				Name:      i.Name,
				Quantity:  i.Quantity,
				DateAdded: i.DateAdded.Unix(),
			}
		}

		w.Header().Set("Content-Type", "application/json")
		enc := json.NewEncoder(w)
		if err := enc.Encode(rsp); err != nil {
			log.Printf("Unexpected error: %s", err)
			http.Error(w, "Unexpected error", http.StatusInternalServerError)
			return
		}
	})
}
