package turtlebunny

import (
	"fmt"
	"strings"

	"lukechampine.com/uint128"
)

type CreateAccountParams struct {
	Id                         uint128.Uint128
	UserData128                uint128.Uint128
	UserData64                 uint64
	UserData32                 uint32
	Ledger                     uint32
	Code                       uint16
	DebitsMustNotExceedCredits bool
	CreditsMustNotExceedDebits bool
}

type Account struct {
	Id                         uint128.Uint128
	DebitsPosted               uint128.Uint128
	CreditsPosted              uint128.Uint128
	UserData128                uint128.Uint128
	UserData64                 uint64
	UserData32                 uint32
	Ledger                     uint32
	Code                       uint16
	DebitsMustNotExceedCredits bool
	CreditsMustNotExceedDebits bool
	Timestamp                  uint64
}

func (c *Client) CreateAccount(params *CreateAccountParams) error {
	_, err := c.db.Exec(`
		INSERT INTO accounts (
			id,
			user_data_128,
			user_data_64,
			user_data_32,
			ledger,
			code,
			debits_must_not_exceed_credits,
			credits_must_not_exceed_debits
		)
		VALUES
		(?, ?, ?, ?, ?, ?, ?, ?)
	`,
		params.Id.String(),
		params.UserData128.String(),
		params.UserData64,
		params.UserData32,
		params.Ledger,
		params.Code,
		params.DebitsMustNotExceedCredits,
		params.CreditsMustNotExceedDebits,
	)

	if err != nil {
		return err
	}
	return nil
}

func (c *Client) LookupAccounts(ids ...uint128.Uint128) ([]*Account, error) {
	placeholders := make([]string, len(ids))
	args := make([]any, len(ids))
	for i, id := range ids {
		placeholders[i] = "?"
		args[i] = id.String()
	}

	query := fmt.Sprintf(`
		SELECT
			id,
			debits_posted,
			credits_posted,
			ledger,
			code,
			debits_must_not_exceed_credits,
			credits_must_not_exceed_debits,
			timestamp
		FROM accounts
		WHERE id IN (%s)
	`, strings.Join(placeholders, ","))

	rows, err := c.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := []*Account{}
	for rows.Next() {
		account := &Account{}
		err := rows.Scan(
			&scannableUint128{&account.Id},
			&scannableUint128{&account.DebitsPosted},
			&scannableUint128{&account.CreditsPosted},
			&account.Ledger,
			&account.Code,
			&account.DebitsMustNotExceedCredits,
			&account.CreditsMustNotExceedDebits,
			&account.Timestamp,
		)

		if err != nil {
			return nil, err
		}

		result = append(result, account)

	}
	return result, nil
}
