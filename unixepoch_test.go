package turtlebunny

import (
	"database/sql"
	"reflect"
	"testing"

	sqlite "github.com/mattn/go-sqlite3"
)

func newTestDriver() *sqlite.SQLiteDriver {
	return &sqlite.SQLiteDriver{
		ConnectHook: func(conn *sqlite.SQLiteConn) error {
			if err := conn.RegisterFunc("custom_unixepoch", unixEpoch, false); err != nil {
				return err
			}
			return nil
		},
	}
}

func TestUnixEpoch(t *testing.T) {
	sql.Register("sqlite3_unixepoch", newTestDriver())

	db, err := sql.Open("sqlite3_unixepoch", ":memory:")
	if err != nil {
		t.Error(err)
	}

	stmt, err := db.Prepare("SELECT unixepoch(?), custom_unixepoch(?)")

	testCases := []struct {
		Arg1 any
		Arg2 any
	}{
		{"now", "now"},
		{"2026-01-01", "2026-01-01"},
		{"2026-01-01 01:23", "2026-01-01 01:23"},
		{"2026-01-01T01:23", "2026-01-01T01:23"},
		{"2026-01-01 01:23:45", "2026-01-01 01:23:45"},
		{"2026-01-01T01:23:45", "2026-01-01T01:23:45"},
		{"2026-01-01 01:23:45.123", "2026-01-01 01:23:45.123"},
		{"2026-01-01T01:23:45.123", "2026-01-01T01:23:45.123"},
		{"01:23", "01:23"},
		{100, 100},
	}

	for _, tc := range testCases {
		var result1 any
		var result2 any
		err = stmt.QueryRow(tc.Arg1, tc.Arg2).Scan(&result1, &result2)
		if err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(result1, result2) {
			t.Errorf("results not equivalent:\n unixepoch(%v) -> %v, custom_unixepoch(%v) -> %v", tc.Arg1, result1, tc.Arg2, result2)
		}

	}

}
