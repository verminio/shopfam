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
)

//go:embed html/*
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
	router := server.Router()
	router.RegisterFS("", "html", fs.FS(content))
	router.HandleFunc(http.MethodPut, "/api/items", api.CreateItem(shopping.NewRepository(d)))

	err = http.ListenAndServe(listenAddr, router)

	if err != nil {
		log.Fatal(err)
	}
}
