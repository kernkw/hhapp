// Package datamock provides a mockable implementation for data.Store.
package datamock

import (
	"github.com/kernkw/hhapp/internal/schema"
)

// Mock implements data.Store. Methods are implemented by setting
// similarly-named callback fields. Calling a method for which no
// corresponding callback has been set will result in a panic.
type Mock struct {
	CreateUser_      func(schema.User) (int, error)
	GetUser_         func(schema.User) (schema.User, error)
	CreateVenue_     func(schema.Venue) (int, error)
	CreateVenueList_ func(schema.VenueList) (int, error)
	VenueListAdd_    func(schema.VenueListAdd) (int, error)
	CreateMenu_      func(schema.Menu) (int, error)
	AddToMenu_       func(schema.MenuItem) (int, error)
	VenueListGet_    func(schema.VenueList) (schema.VenueList, error)
	VenueByList_     func(schema.VenueList) ([]schema.Venue, error)
	VenuesByList_    func(int) ([]schema.Venue, error)
	VenueGet_        func(schema.Venue) (schema.Venue, error)
	MenuItemsGet_    func(schema.Menu) ([]schema.MenuItem, error)
}

func (s *Mock) CreateUser(u schema.User) (int, error)                      { return s.CreateUser_(u) }
func (s *Mock) GetUser(u schema.User) (schema.User, error)                 { return s.GetUser_(u) }
func (s *Mock) CreateVenue(v schema.Venue) (int, error)                    { return s.CreateVenue_(v) }
func (s *Mock) CreateVenueList(vl schema.VenueList) (int, error)           { return s.CreateVenueList_(vl) }
func (s *Mock) VenueListAdd(vla schema.VenueListAdd) (int, error)          { return s.VenueListAdd_(vla) }
func (s *Mock) CreateMenu(menu schema.Menu) (int, error)                   { return s.CreateMenu_(menu) }
func (s *Mock) AddToMenu(menuItem schema.MenuItem) (int, error)            { return s.AddToMenu_(menuItem) }
func (s *Mock) VenueListGet(vl schema.VenueList) (schema.VenueList, error) { return s.VenueListGet_(vl) }
func (s *Mock) VenuesByList(id int) ([]schema.Venue, error)                { return s.VenuesByList_(id) }
func (s *Mock) VenueGet(v schema.Venue) (schema.Venue, error)              { return s.VenueGet_(v) }
func (s *Mock) MenuItemsGet(m schema.Menu) ([]schema.MenuItem, error)      { return s.MenuItemsGet_(m) }

// func (s *Mock) Close()                                     { return }

// func (s *Mock) transaction(c *sql.DB, fn func(tx *sql.Tx) (bool, error)) error { return nil }
