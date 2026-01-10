package turtlebunny

import (
	"errors"
	"fmt"
	"strings"

	"lukechampine.com/uint128"
)

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
	Imported        bool
	Timestamp       uint64
}

func (c *Client) CreateTransfer(transfer Transfer) error {
	if !transfer.Imported {
		if transfer.Timestamp == 0 {
			transfer.Timestamp = unixNano()
		} else {
			return errors.New("timestamp_must_be_zero")
		}
	}
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
			code,
			imported,
			timestamp
		)
		VALUES
		(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		transfer.Id.String(),
		transfer.DebitAccountId.String(),
		transfer.CreditAccountId.String(),
		transfer.Amount.String(),
		transfer.UserData128.String(),
		transfer.UserData64,
		transfer.UserData32,
		transfer.Ledger,
		transfer.Code,
		transfer.Imported,
		transfer.Timestamp,
	)

	if err != nil {
		return err
	}
	return nil
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
			imported,
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
			&transfer.Imported,
			&transfer.Timestamp,
		)
		if err != nil {
			return nil, err
		}

		result = append(result, transfer)
	}

	return result, nil
}
