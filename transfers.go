package turtlebunny

import (
	"fmt"
	"strings"

	"lukechampine.com/uint128"
)

type CreateTransferParams struct {
	Id              uint128.Uint128
	DebitAccountId  uint128.Uint128
	CreditAccountId uint128.Uint128
	Amount          uint128.Uint128
	UserData128     uint128.Uint128
	UserData64      uint64
	UserData32      uint32
	Ledger          uint32
	Code            uint16
}

func (c *Client) CreateTransfer(params CreateTransferParams) error {
	_, err := c.db.Exec(`
		INSERT INTO transfers (
			id,
			debit_account_id,
			credit_account_id,
			amount,
			user_data_128,
			user_data_64,
			user_data_32,
			ledger,
			code
		)
		VALUES
		(?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		params.Id.String(),
		params.DebitAccountId.String(),
		params.CreditAccountId.String(),
		params.Amount.String(),
		params.UserData128.String(),
		params.UserData64,
		params.UserData32,
		params.Ledger,
		params.Code,
	)

	if err != nil {
		return err
	}
	return nil
}

type Transfer struct {
	Id              uint128.Uint128
	DebitAccountId  uint128.Uint128
	CreditAccountId uint128.Uint128
	Amount          uint128.Uint128
	UserData128     uint128.Uint128
	UserData64      uint64
	UserData32      uint32
	Ledger          uint32
	Code            uint16
	Timestamp       uint64
}

func (c *Client) LookupTransfers(ids ...uint128.Uint128) ([]*Transfer, error) {
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
			NewScannableUint128(&transfer.Id),
			NewScannableUint128(&transfer.DebitAccountId),
			NewScannableUint128(&transfer.CreditAccountId),
			NewScannableUint128(&transfer.Amount),
			NewScannableUint128(&transfer.UserData128),
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
