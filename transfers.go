package turtlebunny

import (
	"fmt"
	"strings"
)

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
	Id              Uint128
	DebitAccountId  Uint128
	CreditAccountId Uint128
	Amount          Uint128
	UserData128     Uint128
	UserData64      uint64
	UserData32      uint32
	Ledger          uint32
	Code            uint16
	Timestamp       uint64
}

func (c *Client) LookupTransfers(ids ...Uint128) ([]*Transfer, error) {
	placeholders := make([]string, len(ids))
	args := make([]any, len(ids))
	for i, id := range ids {
		placeholders[i] = "?"
		args[i] = id.String()
	}

	query := fmt.Sprintf(`
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
		WHERE id IN (%s)
	`, strings.Join(placeholders, ","),
	)

	rows, err := c.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := []*Transfer{}
	for rows.Next() {
		transfer := &Transfer{}
		err := rows.Scan(
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

		result = append(result, transfer)
	}

	return result, nil
}
