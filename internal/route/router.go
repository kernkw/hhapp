package route

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kernkw/hhapp/internal/data"
	"github.com/kernkw/hhapp/internal/event"
)

func NewRouter(db *data.Store) *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	routes := getRoutes(db)
	for _, route := range routes {
		var handler http.Handler

		handler = route.HandlerFunc
		handler = event.Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)

	}

	return router
}
