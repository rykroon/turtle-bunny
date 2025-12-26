package turtlebunny

import (
	"database/sql"
	_ "embed"

	sqlite "github.com/mattn/go-sqlite3"
)

//go:embed format.sql
var formatSql string

type Client struct {
	db *sql.DB
}

func NewClient(filename string) (*Client, error) {
	sql.Register("sqlite3_custom", &sqlite.SQLiteDriver{
		ConnectHook: func(conn *sqlite.SQLiteConn) error {
			if err := conn.RegisterFunc("decimal_add", decimalAdd, true); err != nil {
				return err
			}
			if err := conn.RegisterFunc("decimal_sub", decimalSub, true); err != nil {
				return err
			}
			if err := conn.RegisterFunc("is_uint128", isUint128, true); err != nil {
				return err
			}
			if err := conn.RegisterFunc("is_uint64", isUint64, true); err != nil {
				return err
			}
			if err := conn.RegisterFunc("unix_nano", unixNano, true); err != nil {
				return err
			}
			return nil
		},
	})

	db, err := sql.Open("sqlite3_custom", filename)
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
