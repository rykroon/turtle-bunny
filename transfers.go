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

type Transfer struct {
	Id              uint64
	DebitAccountId  uint64
	CreditAccountId uint64
	Amount          uint64
	UserData128     uint64
	UserData64      uint64
	UserData32      uint32
	Ledger          uint32
	Code            uint16
	Timestamp       uint64
}

func (c *Client) LookupTransfer(id uint64) (*Transfer, error) {
	transfer := &Transfer{}
	err := c.db.QueryRow(`
		SELECT
			id,
			debit_account_id,
			credit_account_id,
			amount,
			user_data_128,
			user_data_64,
			user_data_32,
			ledger,
			code,
			timestamp
		FROM transfers
		WHERE id = ?
	`, id).Scan(
		&transfer.Id,
		&transfer.DebitAccountId,
		&transfer.CreditAccountId,
		&transfer.Amount,
		&transfer.UserData128,
		&transfer.UserData64,
		&transfer.UserData32,
		&transfer.Ledger,
		&transfer.Code,
		&transfer.Timestamp,
	)
	if err != nil {
		return nil, err
	}
	return transfer, nil
}
