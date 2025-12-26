package turtlebunny

import (
	"fmt"
	"strings"

	"github.com/rykroon/turtlebunny/internal/uint128x"
)

type CreateAccountParams struct {
	Id                         Uint128
	Ledger                     uint32
	Code                       uint16
	DebitsMustNotExceedCredits bool
	CreditsMustNotExceedDebits bool
}

func (c *Client) CreateAccount(params *CreateAccountParams) error {
	_, err := c.db.Exec(`
		INSERT INTO accounts
		(id, ledger, code, debits_must_not_exceed_credits, credits_must_not_exceed_debits)
		VALUES
		(?, ?, ?, ?, ?)
	`, params.Id.String(), params.Ledger, params.Code, params.DebitsMustNotExceedCredits, params.CreditsMustNotExceedDebits)

	if err != nil {
		return err
	}
	return nil
}

type Account struct {
	Id                         Uint128
	DebitsPosted               Uint128
	CreditsPosted              Uint128
	Ledger                     uint32
	Code                       uint16
	DebitsMustNotExceedCredits bool
	CreditsMustNotExceedDebits bool
	Timestamp                  uint64
}

func (c *Client) LookupAccounts(ids ...Uint128) ([]*Account, error) {
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
			uint128x.NewScannableUint128(&account.Id),
			uint128x.NewScannableUint128(&account.DebitsPosted),
			uint128x.NewScannableUint128(&account.CreditsPosted),
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
