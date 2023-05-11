package main

import (
	"embed"
	"flag"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/verminio/shopfam/api"
	"github.com/verminio/shopfam/server"
	"github.com/verminio/shopfam/shopping"

	"github.com/go-chi/chi/v5"
)

//go:embed html
var content embed.FS

//go:embed db/migrations/*.sql
var migrations embed.FS

var listenAddr string = "localhost:3000"
var dbUrl string = "/tmp/shopfam/dev.db"

func init() {
	flag.StringVar(&listenAddr, "listen-address", envOrString("LISTEN_ADDR", listenAddr), "HTTP Listen Address")
	flag.StringVar(&dbUrl, "db-url", envOrString("DB_URL", dbUrl), "Sqlite database path")
}

func envOrString(key string, def string) string {
	if val, found := os.LookupEnv(key); found {
		return val
	}

	return def
}

func main() {
	flag.Parse()

	d, err := server.DB(dbUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer d.Close()

	if err = server.RunMigrations(d, migrations); err != nil {
		log.Fatal(err)
	}

	log.Println("Starting HTTP Server...")

	files := fs.FS(content)
	html, _ := fs.Sub(files, "html")
	fileServer := http.FileServer(http.FS(html))
	itemService := shopping.NewItemService(shopping.NewRepository(d))

	router := chi.NewRouter()
	router.Put("/api/items", api.UpsertItem(itemService))
	router.Get("/api/items", api.ListItems(itemService))
	router.Delete("/api/items/{itemId}", api.DeleteItem(itemService))
	router.Handle("/*", fileServer)

	err = http.ListenAndServe(listenAddr, router)

	if err != nil {
		log.Fatal(err)
	}
}
