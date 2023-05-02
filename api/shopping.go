package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/verminio/shopfam/shopping"
)

func CreateItem(repository shopping.Repository) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		type createItem struct {
			Id        string `json:"id"`
			Name      string `json:"name"`
			Quantity  string `json:"quantity"`
			DateAdded int64  `json:"dateAdded"`
		}

		created := &createItem{}

		dec := json.NewDecoder(req.Body)
		if err := dec.Decode(created); err != nil {
			http.Error(w, "Invalid Request Body", http.StatusBadRequest)
			return
		}

		item := shopping.New(created.Id, created.Name, created.Quantity, time.UnixMilli(created.DateAdded))

		if err := repository.SaveItem(item); err != nil {
			log.Printf("Unexpected error: %s", err)
			http.Error(w, "Unexpected error", http.StatusInternalServerError)
			return
		}

		http.Error(w, "Created", http.StatusCreated)
	})
}
