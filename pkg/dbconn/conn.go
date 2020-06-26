package dbconn

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	// pq driver must be imported
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"time"
)

// ConnectionConfig for db
type ConnectionConfig struct {
	Host         string
	User         string
	Pass         string
	Port         string
	DB           string
	Timeout      time.Duration
	MaxOpenConns int
	MaxIdleConns int
}

// GetConnection returns a database connection, remember to close it
//
// "A DB instance is not a connection, but an abstraction representing a
// Database... It maintains a connection pool internally"
// Open vs Connect
// "open a DB and connect at the same time;
// for instance, in order to catch configuration issues during your
// initialization phase"
// http://jmoiron.github.io/sqlx/#connecting
//
// "...use DB.SetMaxOpenConns to set the maximum size of the pool
// ...set the maximum idle size with DB.SetMaxIdleConns"
// http://jmoiron.github.io/sqlx/#connectionPool
//
func GetConnection(c *ConnectionConfig) (db *sqlx.DB, err error) {
	if c.MaxOpenConns == 0 {
		c.MaxOpenConns = 2
	}
	if c.MaxIdleConns == 0 {
		c.MaxIdleConns = 2
	}

	var dataSourceName string
	if c.Pass == "" {
		dataSourceName = fmt.Sprintf(
			"host=%s port=%s user=%s dbname=%s sslmode=disable",
			c.Host, c.Port, c.User, c.DB)
	} else {
		dataSourceName = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			c.Host, c.Port, c.User, c.Pass, c.DB)
	}

	// Timeout if connection takes too long
	// https://gobyexample.com/timeouts
	if c.Timeout == 0 {
		c.Timeout = time.Second * 3
	}
	type dbMsg struct {
		DB  *sqlx.DB
		Err error
	}
	dbChan := make(chan dbMsg, 1)
	go func() {
		db, err = sqlx.Connect("postgres", dataSourceName)
		dbChan <- dbMsg{
			DB:  db,
			Err: errors.WithStack(err),
		}
	}()

	select {
	case res := <-dbChan:
		if res.Err != nil {
			return nil, res.Err
		}
		db = res.DB
		db.SetMaxOpenConns(c.MaxOpenConns)
		db.SetMaxIdleConns(c.MaxIdleConns)
		return db, nil
	case <-time.After(c.Timeout):
		return nil, errors.WithStack(fmt.Errorf("timeout connecting to db"))
	}
}

// DB returns a new or existing connection with predefined config
type DB func() (db *sqlx.DB, err error)
