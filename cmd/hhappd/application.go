package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/kelseyhightower/envconfig"
	"github.com/kernkw/hhapp/internal/config"
	"github.com/kernkw/hhapp/internal/data"
	"github.com/kernkw/hhapp/internal/route"
)

var (
	cfg = new(config.Config)
)

func init() {
	err := envconfig.Process("HHAPP", cfg)
	if err != nil {
		log.Println("failed to load environment:", err)
		os.Exit(1)
	}
}

func main() {
	db, err := data.NewStore(cfg)
	if err != nil {
		panic(err)
	}

	router := route.NewRouter(db)
	// bind := fmt.Sprintf("%s:%d", cfg.Addr, cfg.Port)
	bind := fmt.Sprintf("%s:%d", "localhost", cfg.Port)
	log.Printf("serving http on %s", bind)

	log.Fatal(http.ListenAndServe(bind, router))
}
