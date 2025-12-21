package turtlebunny

type CreateTransferParams struct {
	Id              uint64
	DebitAccountId  uint64
	CreditAccountId uint64
	Amount          uint64
	Ledger          uint32
	Code            uint16
}

func (c *Client) CreateTransfer(params CreateTransferParams) error {
	_, err := c.db.Exec(`
		INSERT INTO transfers
		(id, debit_account_id, credit_account_id, amount, ledger, code)
		VALUES
		(?, ?, ?, ?, ?, ?)
	`, params.Id, params.DebitAccountId, params.CreditAccountId, params.Amount, params.Ledger, params.Code)

	if err != nil {
		return err
	}
	return nil
}
