package route

import (
	"net/http"

	"github.com/kernkw/hhapp/internal/data"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func getRoutes(s *data.Store) Routes {
	routes := Routes{
		Route{
			"VenueCreate",
			"POST",
			"/create_venue",
			VenueCreate(s),
		},
		Route{
			"AccountCreate",
			"POST",
			"/create_account",
			UserCreate(s),
		},
		Route{
			"UserLogin",
			"POST",
			"/authenticate",
			UserLogin(s),
		},
	}
	return routes
}
