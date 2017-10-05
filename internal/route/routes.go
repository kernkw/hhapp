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
			"VenueListCreate",
			"POST",
			"/create_venue_list",
			VenueListCreate(s),
		},
		Route{
			"VenueListAdd",
			"POST",
			"/venue_list_add",
			VenueListAdd(s),
		},
		Route{
			"VenueListGet",
			"GET",
			"/venue_list",
			VenueListGet(s),
		},
		Route{
			"VenueGet",
			"GET",
			"/venue/{id:[0-9]+}",
			VenueGet(s),
		},
		Route{
			"MenuItemAdd",
			"POST",
			"/add_menu_item",
			MenuItemAdd(s),
		},
		Route{
			"MenuItemsGet",
			"GET",
			"/menu_items",
			MenuItemsGet(s),
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
