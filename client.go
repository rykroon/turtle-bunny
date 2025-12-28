package turtlebunny

import (
	"database/sql"
	_ "embed"
	"regexp"
	"slices"

	sqlite "github.com/mattn/go-sqlite3"
)

//go:embed format.sql
var formatSql string

type Client struct {
	db *sql.DB
}

func newDriver() *sqlite.SQLiteDriver {
	return &sqlite.SQLiteDriver{
		ConnectHook: func(conn *sqlite.SQLiteConn) error {
			if err := conn.RegisterFunc("regexp", regexp.MatchString, true); err != nil {
				return err
			}
			if err := conn.RegisterFunc("decimal", toDecimal, true); err != nil {
				return err
			}
			if err := conn.RegisterFunc("decimal_add", decimalAdd, true); err != nil {
				return err
			}
			if err := conn.RegisterFunc("decimal_sub", decimalSub, true); err != nil {
				return err
			}
			if err := conn.RegisterFunc("decimal_cmp", decimalCmp, true); err != nil {
				return err
			}
			return nil
		},
	}
}

func NewClient(filename string) (*Client, error) {
	if !slices.Contains(sql.Drivers(), "sqlite3_turtlebunny") {
		sql.Register("sqlite3_turtlebunny", newDriver())
	}

	db, err := sql.Open("sqlite3_turtlebunny", filename)
	if err != nil {
		return nil, err
	}

	return &Client{db: db}, nil
}

func (c *Client) Close() error {
	return c.db.Close()
}

func (c *Client) Format() error {
	_, err := c.db.Exec(formatSql) // format the database if it isn't already.
	if err != nil {
		return err
	}
	return nil
}
