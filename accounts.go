package turtlebunny

type CreateAccountParams struct {
	Id                         uint64
	Ledger                     uint32
	Code                       uint16
	DebitsMustNotExceedCredits bool
	CreditsMustNotExceedDebits bool
}

func (c *Client) CreateAccount(params CreateAccountParams) error {
	_, err := c.db.Exec(`
		INSERT INTO accounts
		(id, ledger, code, debits_must_not_exceed_credits, credits_must_not_exceed_debits)
		VALUES
		(?, ?, ?, ?, ?)
	`, params.Id, params.Ledger, params.Code, params.DebitsMustNotExceedCredits, params.CreditsMustNotExceedDebits)

	if err != nil {
		return err
	}
	return nil
}

type Account struct {
	Id                         uint64
	DebitsPosted               uint64
	CreditsPosted              uint64
	Ledger                     uint32
	Code                       uint16
	DebitsMustNotExceedCredits bool
	CreditsMustNotExceedDebits bool
	Timestamp                  uint64
}

func (c *Client) LookupAccount(id uint64) {

}
