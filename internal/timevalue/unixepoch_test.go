package timevalue

import (
	"database/sql"
	"math"
	"slices"
	"testing"

	sqlite "github.com/mattn/go-sqlite3"
)

func newClient() (*sql.DB, error) {
	if !slices.Contains(sql.Drivers(), "sqlite3_unixepoch") {
		sql.Register(
			"sqlite3_unixepoch",
			&sqlite.SQLiteDriver{
				ConnectHook: func(conn *sqlite.SQLiteConn) error {
					if err := conn.RegisterFunc("unixepoch_custom", UnixEpoch, false); err != nil {
						return err
					}
					return nil
				},
			})
	}

	db, err := sql.Open("sqlite3_unixepoch", ":memory:")
	if err != nil {
		return nil, err
	}
	return db, nil
}

func TestUnixEpochNoArgs(t *testing.T) {
	db, err := newClient()
	if err != nil {
		t.Error(err)
	}

	var r1, r2 any
	err = db.QueryRow("SELECT unixepoch(), unixepoch_custom()").Scan(&r1, &r2)
	if err != nil {
		t.Error(err)
	}

	if r1 != r2 {
		t.Errorf("not equivalent: %v, %v", r1, r2)
	}
}

func TestUnixEpochOneArgs(t *testing.T) {
	db, err := newClient()
	if err != nil {
		t.Error(err)
	}

	stmt, err := db.Prepare("SELECT unixepoch(?), unixepoch_custom(?)")
	if err != nil {
		t.Error(err)
	}

	testCases := []struct {
		Name string
		Arg  any
	}{
		{"Now", "now"},
		{"Subsec", "subsec"},
		{"Subsecond", "subsecond"},
		{"Date", "2025-01-01"},
		{"DateTime", "2025-01-01 12:00:00"},
		{"DateTimeSubSec", "2025-01-01 12:00:00.123456"},
		{"Hour-Min", "12:30"},
		{"Time", "12:30:00"},
		{"TimeSubSec", "12:30:00.123456"},
		{"integer", 0},
		{"float", 1.23},
		{"not valid", "Hello World"},
	}

	for _, tc := range testCases {
		if tc.Arg == "subsec" || tc.Arg == "subsecond" {
			var r1, r2 float64
			err := stmt.QueryRow(tc.Arg, tc.Arg).Scan(&r1, &r2)
			if err != nil {
				t.Error(err)
			}

			diff := math.Abs(r1 - r2)
			if diff > .001 {
				t.Errorf("test %s failed: r1=%v, r2=%v, diff=%v", tc.Name, r1, r2, diff)
			}
		} else {
			var r1, r2 any
			err := stmt.QueryRow(tc.Arg, tc.Arg).Scan(&r1, &r2)
			if err != nil {
				t.Error(err)
			}
			if r1 != r2 {
				t.Errorf("test %s failed: r1=%v, r2=%v", tc.Name, r1, r2)
			}
		}
	}
}
