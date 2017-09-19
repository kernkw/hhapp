package data

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jpillora/backoff"
	"github.com/kernkw/hhapp/internal/schema"
)

const (
	mysql = "mysql"
)

type Database interface {
	CreateUser(user schema.User) (int, error)
	GetUser(user schema.User) (schema.User, error)
	CreateVenue(venue schema.Venue) (int, error)
}

func NewStore() (*Store, error) {
	conn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?parseTime=true", "root", "", "localhost", "happy_hour")

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

var ErrDuplicateEntry = errors.New("Duplicate entry")

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
		q := `INSERT INTO venue (name, address, address2, city, state, zip, country, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
		res, err := tx.Exec(q, venue.Name, venue.Address, venue.Address2, venue.City, venue.State, venue.Zip, venue.Country, time.Now().UTC())
		if err != nil && strings.Contains(err.Error(), "Duplicate entry") {
			return true, ErrDuplicateEntry
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
