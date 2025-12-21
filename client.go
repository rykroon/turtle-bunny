package turtlebunny

import (
	"database/sql"
	_ "embed"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed format.sql
var formatSql string

type Client struct {
	db *sql.DB
}

func NewClient(filename string) (*Client, error) {
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(formatSql) // format the database if it isn't already.
	if err != nil {
		db.Close()
		return nil, err
	}
	return &Client{db: db}, nil
}

func (c *Client) Close() {
	c.db.Close()
}
