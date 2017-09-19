package main

import (
	"log"
	"net/http"

	"github.com/kernkw/hhapp/internal/data"
	"github.com/kernkw/hhapp/internal/route"
)

func main() {
	db, err := data.NewStore()
	if err != nil {
		panic(err)
	}

	router := route.NewRouter(db)

	log.Fatal(http.ListenAndServe(":8080", router))
}
