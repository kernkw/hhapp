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
			"UserFavoriteCreate",
			"POST",
			"/create_user_favorite",
			UserFavoriteCreate(s),
		},
		Route{
			"UserFavoritesList",
			"GET",
			"/user_favorites",
			UserFavoritesList(s),
		},
		Route{
			"UserFavoritesGet",
			"GET",
			"/user_favorites/{venue_id:[0-9]+}/{user_id}",
			UserFavoritesGet(s),
		},
		Route{
			"UserFavoritesRemove",
			"POST",
			"/user_favorite/{id:[0-9]+}",
			UserFavoritesRemove(s),
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
