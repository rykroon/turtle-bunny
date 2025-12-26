package turtlebunny

import (
	"github.com/rykroon/turtlebunny/internal/uint128x"
	"lukechampine.com/uint128"
)

type CreateAccountParams struct {
	Id                         uint128.Uint128
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
	Id                         uint128.Uint128
	DebitsPosted               uint128.Uint128
	CreditsPosted              uint128.Uint128
	Ledger                     uint32
	Code                       uint16
	DebitsMustNotExceedCredits bool
	CreditsMustNotExceedDebits bool
	Timestamp                  uint64
}

func (c *Client) LookupAccount(id uint128.Uint128) (*Account, error) {
	account := &Account{}
	err := c.db.QueryRow(`
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
		WHERE id = ?
	`, id.String()).Scan(
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
	return account, nil
}
