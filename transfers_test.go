package turtlebunny

import (
	"testing"

	"lukechampine.com/uint128"
)

type transferCase struct {
	Name          string
	DbName        string
	Account1      Account
	Account2      Account
	Transfer      Transfer
	ExpectedError string
}

func TestCreateTransfers(t *testing.T) {
	cases := []transferCase{
		{
			Name: "Transfer Id Must Not Be Zero",
			Account1: Account{
				Id:     uint128.From64(1),
				Ledger: 1,
				Code:   1,
			},
			Account2: Account{
				Id:     uint128.From64(2),
				Ledger: 1,
				Code:   1,
			},
			Transfer: Transfer{
				Id:              uint128.Zero,
				DebitAccountId:  uint128.From64(1),
				CreditAccountId: uint128.From64(2),
				Amount:          uint128.From64(100),
				Ledger:          1,
				Code:            1,
			},
			ExpectedError: "id_must_not_be_zero",
		},
		{
			Name: "Transfer Id Must Not Be Int Max",
			Account1: Account{
				Id:     uint128.From64(1),
				Ledger: 1,
				Code:   1,
			},
			Account2: Account{
				Id:     uint128.From64(2),
				Ledger: 1,
				Code:   1,
			},
			Transfer: Transfer{
				Id:              uint128.Max,
				DebitAccountId:  uint128.From64(1),
				CreditAccountId: uint128.From64(2),
				Amount:          uint128.From64(100),
				Ledger:          1,
				Code:            1,
			},
			ExpectedError: "id_must_not_be_int_max",
		},
		{
			Name: "Debit Account Id Must Not Be Zero",
			Account1: Account{
				Id:     uint128.From64(1),
				Ledger: 1,
				Code:   1,
			},
			Account2: Account{
				Id:     uint128.From64(2),
				Ledger: 1,
				Code:   1,
			},
			Transfer: Transfer{
				Id:              uint128.From64(1),
				DebitAccountId:  uint128.Zero,
				CreditAccountId: uint128.From64(2),
				Amount:          uint128.From64(100),
				Ledger:          1,
				Code:            1,
			},
			ExpectedError: "debit_account_id_must_not_be_zero",
		},
		{
			Name: "Debit Account Id Must Not Int Max",
			Account1: Account{
				Id:     uint128.From64(1),
				Ledger: 1,
				Code:   1,
			},
			Account2: Account{
				Id:     uint128.From64(2),
				Ledger: 1,
				Code:   1,
			},
			Transfer: Transfer{
				Id:              uint128.From64(1),
				DebitAccountId:  uint128.Max,
				CreditAccountId: uint128.From64(2),
				Amount:          uint128.From64(100),
				Ledger:          1,
				Code:            1,
			},
			ExpectedError: "debit_account_id_must_not_be_int_max",
		},
		{
			Name: "Credit Account Id Must Not Be Zero",
			Account1: Account{
				Id:     uint128.From64(1),
				Ledger: 1,
				Code:   1,
			},
			Account2: Account{
				Id:     uint128.From64(2),
				Ledger: 1,
				Code:   1,
			},
			Transfer: Transfer{
				Id:              uint128.From64(1),
				DebitAccountId:  uint128.From64(1),
				CreditAccountId: uint128.Zero,
				Amount:          uint128.From64(100),
				Ledger:          1,
				Code:            1,
			},
			ExpectedError: "credit_account_id_must_not_be_zero",
		},
		{
			Name: "Credit Account Id Must Not Int Max",
			Account1: Account{
				Id:     uint128.From64(1),
				Ledger: 1,
				Code:   1,
			},
			Account2: Account{
				Id:     uint128.From64(2),
				Ledger: 1,
				Code:   1,
			},
			Transfer: Transfer{
				Id:              uint128.From64(1),
				DebitAccountId:  uint128.From64(1),
				CreditAccountId: uint128.Max,
				Amount:          uint128.From64(100),
				Ledger:          1,
				Code:            1,
			},
			ExpectedError: "credit_account_id_must_not_be_int_max",
		},
		{
			Name: "Accounts Must Be Different",
			Account1: Account{
				Id:     uint128.From64(1),
				Ledger: 1,
				Code:   1,
			},
			Account2: Account{
				Id:     uint128.From64(2),
				Ledger: 1,
				Code:   1,
			},
			Transfer: Transfer{
				Id:              uint128.From64(1),
				DebitAccountId:  uint128.From64(1),
				CreditAccountId: uint128.From64(1),
				Amount:          uint128.From64(100),
				Ledger:          1,
				Code:            1,
			},
			ExpectedError: "accounts_must_be_different",
		},
		{
			Name: "Debit Account Not Found",
			Account1: Account{
				Id:     uint128.From64(1),
				Ledger: 1,
				Code:   1,
			},
			Account2: Account{
				Id:     uint128.From64(2),
				Ledger: 1,
				Code:   1,
			},
			Transfer: Transfer{
				Id:              uint128.From64(1),
				DebitAccountId:  uint128.From64(3),
				CreditAccountId: uint128.From64(2),
				Amount:          uint128.From64(100),
				Ledger:          1,
				Code:            1,
			},
			ExpectedError: "debit_account_not_found",
		},
		{
			Name: "Credit Account Not Found",
			Account1: Account{
				Id:     uint128.From64(1),
				Ledger: 1,
				Code:   1,
			},
			Account2: Account{
				Id:     uint128.From64(2),
				Ledger: 1,
				Code:   1,
			},
			Transfer: Transfer{
				Id:              uint128.From64(1),
				DebitAccountId:  uint128.From64(1),
				CreditAccountId: uint128.From64(3),
				Amount:          uint128.From64(100),
				Ledger:          1,
				Code:            1,
			},
			ExpectedError: "credit_account_not_found",
		},
		{
			Name: "Accounts Must Have The Same Ledger",
			Account1: Account{
				Id:     uint128.From64(1),
				Ledger: 1,
				Code:   1,
			},
			Account2: Account{
				Id:     uint128.From64(2),
				Ledger: 2,
				Code:   1,
			},
			Transfer: Transfer{
				Id:              uint128.From64(1),
				DebitAccountId:  uint128.From64(1),
				CreditAccountId: uint128.From64(2),
				Amount:          uint128.From64(100),
				Ledger:          1,
				Code:            1,
			},
			ExpectedError: "accounts_must_have_the_same_ledger",
		},
		{
			Name: "Transfer Must Have The Same Ledger As Accounts",
			Account1: Account{
				Id:     uint128.From64(1),
				Ledger: 1,
				Code:   1,
			},
			Account2: Account{
				Id:     uint128.From64(2),
				Ledger: 1,
				Code:   1,
			},
			Transfer: Transfer{
				Id:              uint128.From64(1),
				DebitAccountId:  uint128.From64(1),
				CreditAccountId: uint128.From64(2),
				Amount:          uint128.From64(100),
				Ledger:          2,
				Code:            1,
			},
			ExpectedError: "transfer_must_have_the_same_ledger_as_accounts",
		},
		{
			Name: "Ledger Must Not Be Zero",
			Account1: Account{
				Id:     uint128.From64(1),
				Ledger: 1,
				Code:   1,
			},
			Account2: Account{
				Id:     uint128.From64(2),
				Ledger: 1,
				Code:   1,
			},
			Transfer: Transfer{
				Id:              uint128.From64(1),
				DebitAccountId:  uint128.From64(1),
				CreditAccountId: uint128.From64(2),
				Amount:          uint128.From64(100),
				Ledger:          0,
				Code:            1,
			},
			ExpectedError: "ledger_must_not_be_zero",
		},
		{
			Name: "Code Must Not Be Zero",
			Account1: Account{
				Id:     uint128.From64(1),
				Ledger: 1,
				Code:   1,
			},
			Account2: Account{
				Id:     uint128.From64(2),
				Ledger: 1,
				Code:   1,
			},
			Transfer: Transfer{
				Id:              uint128.From64(1),
				DebitAccountId:  uint128.From64(1),
				CreditAccountId: uint128.From64(2),
				Amount:          uint128.From64(100),
				Ledger:          1,
				Code:            0,
			},
			ExpectedError: "code_must_not_be_zero",
		},
		{
			Name: "Timestamp Must Be Zero",
			Account1: Account{
				Id:     uint128.From64(1),
				Ledger: 1,
				Code:   1,
			},
			Account2: Account{
				Id:     uint128.From64(2),
				Ledger: 1,
				Code:   1,
			},
			Transfer: Transfer{
				Id:              uint128.From64(1),
				DebitAccountId:  uint128.From64(1),
				CreditAccountId: uint128.From64(1),
				Amount:          uint128.From64(100),
				Ledger:          1,
				Code:            1,
				Timestamp:       1234567890,
			},
			ExpectedError: "timestamp_must_be_zero",
		},
		{
			Name: "Timestamp Must Not Advance",
			Account1: Account{
				Id:     uint128.From64(1),
				Ledger: 1,
				Code:   1,
			},
			Account2: Account{
				Id:     uint128.From64(2),
				Ledger: 1,
				Code:   1,
			},
			Transfer: Transfer{
				Id:              uint128.From64(1),
				DebitAccountId:  uint128.From64(1),
				CreditAccountId: uint128.From64(2),
				Amount:          uint128.From64(100),
				Ledger:          1,
				Code:            1,
				Imported:        true,
				Timestamp:       unixNano() + 5e9,
			},
			ExpectedError: "timestamp_must_not_advance",
		},
		{
			Name: "Exceeds Credits",
			Account1: Account{
				Id:                         uint128.From64(1),
				Ledger:                     1,
				Code:                       1,
				DebitsMustNotExceedCredits: true,
			},
			Account2: Account{
				Id:     uint128.From64(2),
				Ledger: 1,
				Code:   1,
			},
			Transfer: Transfer{
				Id:              uint128.From64(1),
				DebitAccountId:  uint128.From64(1),
				CreditAccountId: uint128.From64(2),
				Amount:          uint128.From64(100),
				Ledger:          1,
				Code:            1,
			},
			ExpectedError: "exceeds_credits",
		},
		{
			Name: "Exceeds Debits",
			Account1: Account{
				Id:     uint128.From64(1),
				Ledger: 1,
				Code:   1,
			},
			Account2: Account{
				Id:                         uint128.From64(2),
				Ledger:                     1,
				Code:                       1,
				CreditsMustNotExceedDebits: true,
			},
			Transfer: Transfer{
				Id:              uint128.From64(1),
				DebitAccountId:  uint128.From64(1),
				CreditAccountId: uint128.From64(2),
				Amount:          uint128.From64(100),
				Ledger:          1,
				Code:            1,
			},
			ExpectedError: "exceeds_debits",
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			if tc.DbName == "" {
				tc.DbName = ":memory:"
			}
			client, err := NewClient(tc.DbName)
			if err != nil {
				t.Fatal(err)
			}
			defer client.Close()

			err = client.Format()
			if err != nil {
				t.Fatal(err)
			}

			err = client.CreateAccount(tc.Account1)
			if err != nil {
				t.Fatalf("failed to create account: %s", err.Error())
			}

			err = client.CreateAccount(tc.Account2)
			if err != nil {
				t.Fatalf("failed to create account: %s", err.Error())
			}

			err = client.CreateTransfer(tc.Transfer)
			if err == nil {
				t.Fatal("expected error")
			}

			if err.Error() != tc.ExpectedError {
				t.Fatalf("expected %s, got %s", tc.ExpectedError, err.Error())
			}
		})
	}
}
