package turtlebunny

import (
	"fmt"
	"os"
	"testing"
)

const insertAccountQuery string = `
	INSERT INTO accounts (
		id,
		debits_posted,
		credits_posted,
		user_data_128,
		user_data_64,
		user_data_32,
		ledger,
		code,
		debits_must_not_exceed_credits,
		credits_must_not_exceed_debits
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
`

func TestUint128(t *testing.T) {
	client, err := NewClient("test.db")
	if err != nil {
		t.Error(err)
	}
	defer client.Close()

	err = client.Format()
	if err != nil {
		t.Error(err)
	}

	// negative integer should fail
	_, err = client.db.Exec(
		insertAccountQuery,
		"-123", "0", "0", "0", "0", 0, 1, 1, false, false,
	)
	if err == nil {
		t.Errorf("expected error")
	}

	// decimal should fail
	_, err = client.db.Exec(
		insertAccountQuery,
		"1.23", "0", "0", "0", "0", 0, 1, 1, false, false,
	)
	if err == nil {
		t.Errorf("expected error")
	}

	// leading zeros should fail
	_, err = client.db.Exec(
		insertAccountQuery,
		"0123", "0", "0", "0", "0", 0, 1, 1, false, false,
	)
	if err == nil {
		t.Errorf("expected error")
	}

	// non numeric should fail
	_, err = client.db.Exec(
		insertAccountQuery,
		"Hello World", "0", "0", "0", "0", 0, 1, 1, false, false,
	)
	if err == nil {
		t.Errorf("expected error")
	}

	// overflow should fail
	_, err = client.db.Exec(
		insertAccountQuery,
		"999999999999999999999999999999999999999", "0", "0", "0", "0", 0, 1, 1, false, false,
	)
	if err == nil {
		t.Errorf("expected error")
	}
}

func TestIdMustNotBeZero(t *testing.T) {
	client, err := NewClient("test.db")
	if err != nil {
		t.Error(err)
	}
	defer client.Close()

	err = client.Format()
	if err != nil {
		t.Error(err)
	}

	// negative integer should fail
	_, err = client.db.Exec(
		insertAccountQuery,
		"0", "0", "0", "0", "0", 0, 1, 1, false, false,
	)
	if err == nil {
		t.Errorf("expected error")
	}

	if err.Error() != "id_must_not_be_zero" {
		t.Errorf("expected %s, got %s", "id_must_not_be_zero", err.Error())
	}
}

func TestIdMustNotBeIntMax(t *testing.T) {
	client, err := NewClient("test.db")
	if err != nil {
		t.Error(err)
	}
	defer client.Close()

	err = client.Format()
	if err != nil {
		t.Error(err)
	}

	_, err = client.db.Exec(
		insertAccountQuery,
		"340282366920938463463374607431768211455", "0", "0", "0", "0", 1, 1, 1, false, false,
	)
	if err == nil {
		t.Errorf("expected error")
	}

	if err.Error() != "id_must_not_be_int_max" {
		t.Errorf("expected %s, got %s", "id_must_not_be_int_max", err.Error())
	}
}

func TestFlagsAreMutuallyExclusive(t *testing.T) {
	client, err := NewClient("test.db")
	if err != nil {
		t.Error(err)
	}
	defer client.Close()

	err = client.Format()
	if err != nil {
		t.Error(err)
	}

	_, err = client.db.Exec(
		insertAccountQuery,
		"1", "0", "0", "0", "0", 1, 1, 1, true, true,
	)
	if err == nil {
		t.Errorf("expected error")
	}

	if err.Error() != "flags_are_mutually_exclusive" {
		t.Errorf("expected %s, got %s", "flags_are_mutually_exclusive", err.Error())
	}
}

func TestDebitsPostedMustBeZero(t *testing.T) {
	client, err := NewClient("test.db")
	if err != nil {
		t.Error(err)
	}
	defer client.Close()

	err = client.Format()
	if err != nil {
		t.Error(err)
	}

	_, err = client.db.Exec(
		insertAccountQuery,
		"1", "100", "0", "0", "0", 1, 1, 1, false, false,
	)
	if err == nil {
		t.Errorf("expected error")
	}

	if err.Error() != "debits_posted_must_be_zero" {
		t.Errorf("expected %s, got %s", "debits_posted_must_be_zero", err.Error())
	}
}

func TestCreditsPostedMustBeZero(t *testing.T) {
	client, err := NewClient("test.db")
	if err != nil {
		t.Error(err)
	}
	defer client.Close()

	err = client.Format()
	if err != nil {
		t.Error(err)
	}

	_, err = client.db.Exec(
		insertAccountQuery,
		"1", "0", "100", "0", "0", 0, 1, 1, false, false,
	)
	if err == nil {
		t.Errorf("expected error")
	}

	if err.Error() != "credits_posted_must_be_zero" {
		t.Errorf("expected %s, got %s", "credits_posted_must_be_zero", err.Error())
	}
}

func TestLedgerMustNotBeZero(t *testing.T) {
	client, err := NewClient("test.db")
	if err != nil {
		t.Error(err)
	}
	defer client.Close()

	err = client.Format()
	if err != nil {
		t.Error(err)
	}

	_, err = client.db.Exec(
		insertAccountQuery,
		"1", "0", "0", "0", "0", 0, 0, 1, false, false,
	)

	if err.Error() != "ledger_must_not_be_zero" {
		t.Errorf("expected %s, got %s", "ledger_must_not_be_zero", err.Error())
	}
}

func TestCodeMustNotBeZero(t *testing.T) {
	client, err := NewClient("test.db")
	if err != nil {
		t.Error(err)
	}
	defer client.Close()

	err = client.Format()
	if err != nil {
		t.Error(err)
	}

	_, err = client.db.Exec(
		insertAccountQuery,
		"1", "0", "0", "0", "0", 0, 1, 0, false, false,
	)

	if err.Error() != "code_must_not_be_zero" {
		t.Errorf("expected %s, got %s", "code_must_not_be_zero", err.Error())
	}
}

func TestAccountUpdate(t *testing.T) {
	os.Remove("./test.db")
	client, err := NewClient("test.db")
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
			credits_must_not_exceed_debits
		) VALUES (1, 0, 0, 0, 0, 0, 1, 1, false, false)
	`)

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
