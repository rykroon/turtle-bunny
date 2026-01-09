package turtlebunny

import (
	"fmt"
	"testing"
	"time"
)

func TestAccountUpdate(t *testing.T) {
	client, err := NewClient(":memory:")
	if err != nil {
		t.Error(err)
	}
	defer client.Close()

	err = client.Format()
	if err != nil {
		t.Error(err)
	}

	_, err = client.db.Exec(
		`INSERT INTO accounts (
			id,
			debits_posted,
			credits_posted,
			user_data_128,
			user_data_64,
			user_data_32,
			ledger,
			code,
			debits_must_not_exceed_credits,
			credits_must_not_exceed_debits,
			timestamp
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, 1, 0, 0, 0, 0, 0, 1, 1, false, false, time.Now().UnixNano()) // change to UnixNano once nano seconds are supported.

	if err != nil {
		t.Error(err)
	}

	testCases := []struct {
		Field string
		Value any
	}{
		{Field: "id", Value: "2"},
		{Field: "user_data_128", Value: "123"},
		{Field: "user_data_64", Value: "123"},
		{Field: "user_data_32", Value: 123},
		{Field: "ledger", Value: 2},
		{Field: "code", Value: 2},
		{Field: "debits_must_not_exceed_credits", Value: true},
		{Field: "credits_must_not_exceed_debits", Value: true},
		{Field: "timestamp", Value: "1234567890"},
	}

	for _, tc := range testCases {
		query := fmt.Sprintf("UPDATE accounts SET %s = ? WHERE id = 1", tc.Field)
		_, err = client.db.Exec(query, tc.Value)
		if err == nil {
			t.Error("expected error")
		}

		if err.Error() != "account_cannot_be_modified" {
			t.Errorf("expected %s, got %s", "account_cannot_be_modified", err.Error())
		}
	}

}
