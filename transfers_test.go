package turtlebunny

import (
	"testing"

	"lukechampine.com/uint128"
)

type transferCase struct {
	Name          string
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
	}

	for _, tc := range cases {
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
