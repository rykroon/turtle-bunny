package turtlebunny

import (
	"testing"

	"lukechampine.com/uint128"
)

type createAccountTestCase struct {
	Name          string
	Account       Account
	ExpectedError string
}

func TestCreateAccount(t *testing.T) {
	testCases := []createAccountTestCase{
		{
			Name: "Account Id Must Not Be Zero",
			Account: Account{
				Ledger: 1,
				Code:   1,
			},
			ExpectedError: "id_must_not_be_zero",
		},
		{
			Name: "Id Must Not Be Int Max",
			Account: Account{
				Id:     uint128.Max,
				Ledger: 1,
				Code:   1,
			},
			ExpectedError: "id_must_not_be_int_max",
		},
		{
			Name: "Debits Posted Must Be Zero",
			Account: Account{
				Id:           uint128.From64(1),
				DebitsPosted: uint128.From64(1),
				Ledger:       1,
				Code:         1,
			},
			ExpectedError: "debits_posted_must_be_zero",
		},
		{
			Name: "Credits Posted Must Be Zero",
			Account: Account{
				Id:            uint128.From64(1),
				CreditsPosted: uint128.From64(1),
				Ledger:        1,
				Code:          1,
			},
			ExpectedError: "credits_posted_must_be_zero",
		},
		{
			Name: "Ledger Must Not Be Zero",
			Account: Account{
				Id:     uint128.From64(1),
				Ledger: 0,
				Code:   1,
			},
			ExpectedError: "ledger_must_not_be_zero",
		},
		{
			Name: "Code Must Not Be Zero",
			Account: Account{
				Id:     uint128.From64(1),
				Ledger: 1,
				Code:   0,
			},
			ExpectedError: "code_must_not_be_zero",
		},
		{
			Name: "Flags Are Mutually Exclusive",
			Account: Account{
				Id:                         uint128.From64(1),
				Ledger:                     1,
				Code:                       1,
				DebitsMustNotExceedCredits: true,
				CreditsMustNotExceedDebits: true,
			},
			ExpectedError: "flags_are_mutually_exclusive",
		},
		{
			Name: "Timestamp Must Be Zero",
			Account: Account{
				Id:        uint128.From64(1),
				Ledger:    1,
				Code:      1,
				Imported:  false,
				Timestamp: 1234567890,
			},
			ExpectedError: "timestamp_must_be_zero",
		},
		{
			Name: "Timestamp Must Not Advance",
			Account: Account{
				Id:        uint128.From64(1),
				Ledger:    1,
				Code:      1,
				Imported:  true,
				Timestamp: uint64(unixNano()) + 1e9,
			},
			ExpectedError: "timestamp_must_not_advance",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			client, err := NewClient(":memory:")
			if err != nil {
				t.Fatal(err)
			}
			defer client.Close()

			err = client.Format()
			if err != nil {
				t.Fatal(err)
			}

			err = client.CreateAccount(tc.Account)

			if err == nil {
				t.Fatal("expected error")
			}

			if err.Error() != tc.ExpectedError {
				t.Fatalf("expected %s, got %s", tc.ExpectedError, err.Error())
			}
		})
	}
}

func TestAccountTimestampMustNotRegress(t *testing.T) {
	client, err := NewClient(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	err = client.Format()
	if err != nil {
		t.Fatal(err)
	}

	account1 := Account{
		Id:     uint128.From64(1),
		Ledger: 1,
		Code:   1,
	}
	err = client.CreateAccount(account1)

	if err != nil {
		t.Fatal(err)
	}

	account2 := Account{
		Id:        uint128.From64(2),
		Ledger:    1,
		Code:      1,
		Imported:  true,
		Timestamp: uint64(unixNano()) - 1e9,
	}
	err = client.CreateAccount(account2)

	if err == nil {
		t.Fatal("expected error")
	}

	if err.Error() != "timestamp_must_not_regress" {
		t.Fatalf("expected %s, got %s", "timestamp_must_not_regress", err.Error())
	}

}
