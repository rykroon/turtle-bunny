package cli

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/rykroon/turtlebunny"
	"github.com/spf13/cobra"
)

func NewCreateTransferCmd() *cobra.Command {
	params := turtlebunny.CreateTransferParams{}

	cmd := &cobra.Command{
		Use:   "create-transfer",
		Short: "Create Transfer",
		Args:  requireFilenameArg,
		RunE: func(cmd *cobra.Command, args []string) error {
			filename := args[0]
			client, err := turtlebunny.NewClient(filename)
			if err != nil {
				return err
			}
			defer client.Close()
			err = client.CreateTransfer(params)
			if err != nil {
				return err
			}
			return nil
		},
	}

	cmd.Flags().Uint64VarP(&params.Id, "id", "i", 0, "id")
	cmd.Flags().Uint64VarP(&params.DebitAccountId, "debit-account-id", "D", 0, "debit account id")
	cmd.Flags().Uint64VarP(&params.CreditAccountId, "credit-account-id", "C", 0, "credit account id")
	cmd.Flags().Uint64VarP(&params.Amount, "amount", "a", 0, "amount")
	cmd.Flags().Uint32VarP(&params.Ledger, "ledger", "l", 0, "ledger")
	cmd.Flags().Uint16VarP(&params.Code, "code", "c", 0, "code")

	cmd.MarkFlagRequired("id")
	cmd.MarkFlagRequired("debit-account-id")
	cmd.MarkFlagRequired("credit-account-id")
	cmd.MarkFlagRequired("amount")
	cmd.MarkFlagRequired("ledger")
	cmd.MarkFlagRequired("code")

	return cmd
}

func NewLookupTransferCmd() *cobra.Command {
	var id uint64 = 0

	cmd := &cobra.Command{
		Use:   "lookup-transfer",
		Short: "Lookup Transfer",
		Args:  requireFilenameArg,
		RunE: func(cmd *cobra.Command, args []string) error {
			filename := args[0]
			client, err := turtlebunny.NewClient(filename)
			if err != nil {
				return err
			}
			defer client.Close()
			transfer, err := client.LookupTransfer(id)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return errors.New("transfer not found")
				}
				return err
			}

			fmt.Printf(
				"%7s %17s %18s %10s %6s %5s %20s\n",
				"ID",
				"Debit Account Id",
				"Credit Account Id",
				"Amount",
				"Ledger",
				"Code",
				"Timestamp",
			)
			fmt.Println(strings.Repeat("-", 90))
			fmt.Printf(
				"%7d %17d %18d %10d %6d %5d %20d\n",
				transfer.Id,
				transfer.DebitAccountId,
				transfer.CreditAccountId,
				transfer.Amount,
				transfer.Ledger,
				transfer.Code,
				transfer.Timestamp,
			)
			return nil
		},
	}

	cmd.Flags().Uint64VarP(&id, "id", "i", 0, "id")
	cmd.MarkFlagRequired("id")

	return cmd
}
