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
}

func (s *Mock) CreateUser(u schema.User) (int, error)            { return s.CreateUser_(u) }
func (s *Mock) GetUser(u schema.User) (schema.User, error)       { return s.GetUser_(u) }
func (s *Mock) CreateVenue(v schema.Venue) (int, error)          { return s.CreateVenue_(v) }
func (s *Mock) CreateVenueList(vl schema.VenueList) (int, error) { return s.CreateVenueList_(vl) }

// func (s *Mock) Close()                                     { return }

// func (s *Mock) transaction(c *sql.DB, fn func(tx *sql.Tx) (bool, error)) error { return nil }
