package turtlebunny

import (
	"fmt"
	"testing"

	"lukechampine.com/uint128"
)

func TestAccountUpdate(t *testing.T) {
	client, err := NewClient(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	err = client.Format()
	if err != nil {
		t.Fatal(err)
	}

	err = client.CreateAccount(Account{
		Id:     uint128.From64(1),
		Ledger: 1,
		Code:   1,
	})
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		Name  string
		Field string
		Value any
	}{
		{Name: "Update Account Id", Field: "id", Value: "2"},
		{Name: "Update Account User Data 128", Field: "user_data_128", Value: "123"},
		{Name: "Update Account User Data 64", Field: "user_data_64", Value: "123"},
		{Name: "Update Account User Data 32", Field: "user_data_32", Value: 123},
		{Name: "Update Account Ledger", Field: "ledger", Value: 2},
		{Name: "Update Account Code", Field: "code", Value: 2},
		{Name: "Update Debits Must Not Exceed Credits", Field: "debits_must_not_exceed_credits", Value: true},
		{Name: "Update Credits Must Not Exceed Debits", Field: "credits_must_not_exceed_debits", Value: true},
		{Name: "Update Imported", Field: "imported", Value: true},
		{Name: "Update Account Timestamp", Field: "timestamp", Value: "1234567890"},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			query := fmt.Sprintf("UPDATE accounts SET %s = ? WHERE id = 1", tc.Field)
			_, err = client.db.Exec(query, tc.Value)
			if err == nil {
				t.Fatal("expected error")
			}

			if err.Error() != "account_cannot_be_modified" {
				t.Fatalf("expected %s, got %s", "account_cannot_be_modified", err.Error())
			}
		})

	}

}
