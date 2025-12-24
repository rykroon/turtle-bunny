package turtlebunny

import "strconv"

type CreateTransferParams struct {
	Id              uint64
	DebitAccountId  uint64
	CreditAccountId uint64
	Amount          uint64
	Ledger          uint32
	Code            uint16
}

func (c *Client) CreateTransfer(params CreateTransferParams) error {
	idStr := strconv.Itoa(int(params.Id))
	debitAccountIdStr := strconv.Itoa(int(params.DebitAccountId))
	creditAccountIdStr := strconv.Itoa(int(params.CreditAccountId))
	amountStr := strconv.Itoa(int(params.Amount))

	_, err := c.db.Exec(`
		INSERT INTO transfers
		(id, debit_account_id, credit_account_id, amount, ledger, code)
		VALUES
		(?, ?, ?, ?, ?, ?)
	`, idStr, debitAccountIdStr, creditAccountIdStr, amountStr, params.Ledger, params.Code)

	if err != nil {
		return err
	}
	return nil
}
