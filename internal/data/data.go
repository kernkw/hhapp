package data

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jpillora/backoff"
	"github.com/kernkw/hhapp/internal/config"
	"github.com/kernkw/hhapp/internal/schema"
)

const (
	mysql = "mysql"
)

type Database interface {
	CreateUser(user schema.User) (int, error)
	CreateUserFavorite(userFav schema.UserFavorite) (int, error)
	UserFavoritesList(u schema.UserFavorite) ([]schema.Venue, error)
	UserFavoritesGet(u schema.UserFavorite) (schema.Venue, error)
	UserFavoritesDelete(id int) error
	GetUser(user schema.User) (schema.User, error)
	CreateVenue(venue schema.Venue) (int, error)
	CreateVenueList(venueList schema.VenueList) (int, error)
	VenueListAdd(vla schema.VenueListAdd) (int, error)
	CreateMenu(menu schema.Menu) (int, error)
	AddToMenu(menuItem schema.MenuItem) (int, error)
	VenueListGet(vl schema.VenueList) (schema.VenueList, error)
	VenuesByList(id int) ([]schema.Venue, error)
	VenueGet(v schema.Venue) (schema.Venue, error)
	MenuItemsGet(m schema.Menu) ([]schema.MenuItem, error)
}

func NewStore(cfg *config.Config) (*Store, error) {
	dsn := "%s:%s@tcp(%s:%d)/%s?parseTime=true"
	conn := fmt.Sprintf(dsn, cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
	fmt.Println(conn)
	db, err := sql.Open(mysql, conn)
	if err != nil {
		return nil, err
	}

	store := &Store{db: db}

	return store, nil
}

// Store is the database connection.
type Store struct {
	db *sql.DB
}

// Close closes all database related connections.
func (s *Store) Close() {
	s.db.Close()
}

var ErrDuplicateEntry = errors.New("duplicate entry")
var ErrNotFound = errors.New("no matching records found")

func (s *Store) CreateUser(user schema.User) (int, error) {
	var id int
	err := s.transaction(s.db, func(tx *sql.Tx) (bool, error) {
		q := `INSERT INTO user (username, password, email, created_at) VALUES (?, ?, ?, ?)`
		res, err := tx.Exec(q, user.UserName, user.Password, user.Email, time.Now().UTC())
		if err != nil && strings.Contains(err.Error(), "Duplicate entry") {
			return true, ErrDuplicateEntry
		}
		resID, err := res.LastInsertId()
		id = int(resID)
		return false, err
	})
	return id, err
}

func (s *Store) CreateUserFavorite(userFav schema.UserFavorite) (int, error) {
	var id int
	err := s.transaction(s.db, func(tx *sql.Tx) (bool, error) {
		q := `INSERT INTO user_favorites (user_id, venue_id, created_at) VALUES (?, ?, ?)`
		res, err := tx.Exec(q, userFav.UserID, userFav.VenueID, time.Now().UTC())
		if err != nil {
			return true, err
		}
		resID, err := res.LastInsertId()
		id = int(resID)
		return false, err
	})
	return id, err
}

func (s *Store) UserFavoritesList(u schema.UserFavorite) ([]schema.Venue, error) {
	var venues []schema.Venue
	err := s.transaction(s.db, func(tx *sql.Tx) (bool, error) {
		query := `SELECT v.id, v.name, v.address, v.address2, v.city, v.state, v.zip, v.country, v.image from user_favorites as uf
					JOIN venue as v on uf.venue_id = v.id
					WHERE uf.user_id = ?`
		rows, err := tx.Query(query, u.UserID)
		if err != nil {
			return false, err
		}
		for rows.Next() {
			var venue schema.Venue
			err := rows.Scan(&venue.ID, &venue.Name, &venue.Address, &venue.Address2, &venue.City, &venue.State, &venue.Zip, &venue.Country, &venue.Image)
			if err != nil {
				return false, err
			}
			venues = append(venues, venue)
		}

		return false, err
	})

	return venues, err
}

func (s *Store) UserFavoritesGet(u schema.UserFavorite) (schema.Venue, error) {
	var venue schema.Venue
	err := s.transaction(s.db, func(tx *sql.Tx) (bool, error) {
		query := `SELECT uf.id, v.name, v.address, v.address2, v.city, v.state, v.zip, v.country, v.image from user_favorites as uf
					JOIN venue as v on uf.venue_id = v.id
					WHERE uf.user_id = ? AND uf.venue_id = ?`
		row := tx.QueryRow(query, u.UserID, u.VenueID)
		err := row.Scan(&venue.ID, &venue.Name, &venue.Address, &venue.Address2, &venue.City, &venue.State, &venue.Zip, &venue.Country, &venue.Image)
		if err != nil {
			return false, err
		}
		return false, err
	})

	return venue, err
}
func (s *Store) UserFavoritesDelete(id int) error {
	err := s.transaction(s.db, func(tx *sql.Tx) (bool, error) {
		query := `DELETE FROM user_favorites WHERE id = ?`
		_, err := tx.Exec(query, id)
		if err != nil {
			return false, err
		}
		return false, err
	})

	return err
}

func (s *Store) GetUser(user schema.User) (schema.User, error) {
	u := schema.User{}
	err := s.transaction(s.db, func(tx *sql.Tx) (bool, error) {
		row := tx.QueryRow(`SELECT id, username, password, email, first_name, last_name FROM user WHERE username=?`, user.UserName)
		row.Scan(&u.ID, &u.UserName, &u.Password, &u.Email, &u.FirstName, &u.LastName)
		return false, nil
	})
	return u, err
}

func (s *Store) CreateVenue(venue schema.Venue) (int, error) {
	var id int
	err := s.transaction(s.db, func(tx *sql.Tx) (bool, error) {
		q := `INSERT INTO venue (name, address, address2, city, state, zip, country, image, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
		fmt.Println(fmt.Sprintf("%+v", venue))
		res, err := tx.Exec(q, venue.Name, venue.Address, venue.Address2, venue.City, venue.State, venue.Zip, venue.Country, venue.Image, time.Now().UTC())
		if err != nil && strings.Contains(err.Error(), "Duplicate entry") {
			return true, ErrDuplicateEntry
		}
		resID, err := res.LastInsertId()
		id = int(resID)
		return false, err
	})

	return id, err
}

func (s *Store) CreateVenueList(venueList schema.VenueList) (int, error) {
	var id int
	err := s.transaction(s.db, func(tx *sql.Tx) (bool, error) {
		q := `INSERT INTO venue_list (name, created_at) VALUES (?, ?)`
		res, err := tx.Exec(q, venueList.Name, time.Now().UTC())
		if err != nil && strings.Contains(err.Error(), "Duplicate entry") {
			return true, ErrDuplicateEntry
		}
		resID, err := res.LastInsertId()
		id = int(resID)
		return false, err
	})

	return id, err
}

func (s *Store) VenueListAdd(vla schema.VenueListAdd) (int, error) {
	var id int
	venueList := schema.VenueList{ID: vla.VenueListID, Name: vla.VenueListName}
	vl, err := s.VenueListGet(venueList)
	if err != nil {
		return 0, fmt.Errorf("venue list %s not found", vla.VenueListName)
	}
	venue := schema.Venue{ID: vla.VenueID, Name: vla.VenueName}
	v, err := s.VenueGet(venue)
	if err != nil {
		return 0, fmt.Errorf("venue %s not found", vla.VenueName)
	}
	err = s.transaction(s.db, func(tx *sql.Tx) (bool, error) {
		q := `INSERT INTO venue_lists (venue_id, venue_list_id, created_at) VALUES (?, ?, ?)`
		res, err := tx.Exec(q, v.ID, vl.ID, time.Now().UTC())
		if err != nil && strings.Contains(err.Error(), "Duplicate entry") {
			return true, ErrDuplicateEntry
		}
		resID, err := res.LastInsertId()
		id = int(resID)
		return false, err
	})

	return id, err
}

func (s *Store) VenueListGet(vl schema.VenueList) (schema.VenueList, error) {
	var query, svalue string
	var venueList schema.VenueList
	switch {
	case vl.ID != 0:
		query = `SELECT id, name FROM venue_list WHERE id = ?`
		svalue = strconv.Itoa(vl.ID)
	case vl.Name != "":
		query = `SELECT id, name FROM venue_list WHERE name = ?`
		svalue = vl.Name
	default:
		return venueList, errors.New("no venue list id or name provided")
	}

	err := s.transaction(s.db, func(tx *sql.Tx) (bool, error) {
		row := tx.QueryRow(query, svalue)
		err := row.Scan(&venueList.ID, &venueList.Name)
		if err == sql.ErrNoRows {
			return true, ErrNotFound
		}
		return false, err
	})

	return venueList, err
}

func (s *Store) VenuesByList(id int) ([]schema.Venue, error) {
	var venues []schema.Venue
	err := s.transaction(s.db, func(tx *sql.Tx) (bool, error) {
		query := `SELECT v.id, v.name, v.address, v.address2, v.city, v.state, v.zip, v.country, v.image from venue_lists as vl
					JOIN venue as v on vl.venue_id = v.id
					WHERE vl.venue_list_id = ?`
		rows, err := tx.Query(query, id)
		if err != nil {
			return false, err
		}
		for rows.Next() {
			var venue schema.Venue
			err := rows.Scan(&venue.ID, &venue.Name, &venue.Address, &venue.Address2, &venue.City, &venue.State, &venue.Zip, &venue.Country, &venue.Image)
			if err != nil {
				return false, err
			}
			venues = append(venues, venue)
		}

		return false, err
	})

	return venues, err
}

func (s *Store) VenueGet(v schema.Venue) (schema.Venue, error) {
	var query, svalue string
	var venue schema.Venue
	switch {
	case v.ID != 0:
		query = `SELECT id, name, address, address2, city, state, zip, country, image FROM venue WHERE id = ?`
		svalue = strconv.Itoa(v.ID)
	case v.Name != "":
		query = `SELECT id, name, address, address2, city, state, zip, country, image FROM venue WHERE name = ?`
		svalue = v.Name
	default:
		return venue, errors.New("no venue id or name provided")
	}

	err := s.transaction(s.db, func(tx *sql.Tx) (bool, error) {
		row := tx.QueryRow(query, svalue)
		err := row.Scan(&venue.ID, &venue.Name, &venue.Address, &venue.Address2, &venue.City, &venue.State, &venue.Zip, &venue.Country, &venue.Image)
		if err == sql.ErrNoRows {
			return true, ErrNotFound
		}
		return false, err
	})

	return venue, err
}

func (s *Store) MenuItemsGet(m schema.Menu) ([]schema.MenuItem, error) {
	var menuItems []schema.MenuItem
	err := s.transaction(s.db, func(tx *sql.Tx) (bool, error) {
		query := `SELECT mi.id, m.id, mi.category, mi.price, mi.description
					FROM menu as m
					JOIN menu_item as mi on m.id = mi.menu_id
					WHERE m.venue_id = ?`
		rows, err := tx.Query(query, m.VenueID)
		if err != nil {
			return false, err
		}
		for rows.Next() {
			var mi schema.MenuItem
			err := rows.Scan(&mi.ID, &mi.MenuID, &mi.Category, &mi.Price, &mi.Description)
			if err != nil {
				return false, err
			}
			menuItems = append(menuItems, mi)
		}

		return false, err
	})

	return menuItems, err
}

func (s *Store) CreateMenu(menu schema.Menu) (int, error) {
	var id int
	err := s.transaction(s.db, func(tx *sql.Tx) (bool, error) {
		q := `INSERT INTO menu (venue_id, created_at) VALUES (?, ?)`
		res, err := tx.Exec(q, menu.VenueID, time.Now().UTC())
		if err != nil && strings.Contains(err.Error(), "Duplicate entry") {
			return true, ErrDuplicateEntry
		}
		resID, err := res.LastInsertId()
		id = int(resID)
		return false, err
	})

	return id, err
}

func (s *Store) AddToMenu(menuItem schema.MenuItem) (int, error) {
	var id int
	fmt.Printf("MenuItem: %+v", menuItem)
	err := s.transaction(s.db, func(tx *sql.Tx) (bool, error) {
		q := `INSERT INTO menu_item (menu_id, category, price, description, created_at) VALUES (?, ?, ?, ?, ?)`
		res, err := tx.Exec(q, menuItem.MenuID, menuItem.Category, menuItem.Price, menuItem.Description, time.Now().UTC())
		if err != nil {
			if strings.Contains(err.Error(), "Duplicate entry") {
				return true, ErrDuplicateEntry
			}
			return true, err
		}
		resID, err := res.LastInsertId()
		id = int(resID)
		return false, err
	})

	return id, err
}

var retryN int64 = 3

func (s *Store) transaction(c *sql.DB, fn func(tx *sql.Tx) (bool, error)) error {
	return retry(retryN, func() (bool, error) {
		tx, err := c.Begin()
		if err != nil {
			return false, err
		}

		bypass, err := fn(tx)
		if err != nil {
			tx.Rollback()
			return bypass, err
		}

		return bypass, tx.Commit()
	})
}

func retry(maxRetry int64, fn func() (bool, error)) error {
	backoff := backoff.Backoff{
		Jitter: true,
		Factor: 1.25,
		Min:    100 * time.Millisecond,
		Max:    10 * time.Second,
	}
	retry := maxRetry + 1
	for {
		bypass, err := fn()
		if bypass {
			return err
		}

		retry--
		if err != nil {
			if retry == 0 && maxRetry >= 0 {
				return fmt.Errorf("Maximum number of retries exceeded: %+v", err)
			}
			time.Sleep(backoff.Duration())
			continue
		}
		return nil
	}
}

func newConnection(driver, dataSource string) (*sql.DB, error) {
	return sql.Open(driver, dataSource)
}
