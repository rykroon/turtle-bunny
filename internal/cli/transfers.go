package cli

import (
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
	ids := []turtlebunny.Uint128{}

	cmd := &cobra.Command{
		Use:   "lookup-transfers",
		Short: "Lookup Transfers",
		Args:  requireFilenameArg,
		RunE: func(cmd *cobra.Command, args []string) error {
			filename := args[0]
			client, err := turtlebunny.NewClient(filename)
			if err != nil {
				return err
			}
			defer client.Close()
			transfers, err := client.LookupTransfers(ids...)
			if err != nil {
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

			for _, transfer := range transfers {
				fmt.Printf(
					"%7s %17s %18s %10s %6d %5d %20d\n",
					transfer.Id.String(),
					transfer.DebitAccountId.String(),
					transfer.CreditAccountId.String(),
					transfer.Amount.String(),
					transfer.Ledger,
					transfer.Code,
					transfer.Timestamp,
				)

			}
			return nil
		},
	}

	Uint128SliceVarP(cmd.Flags(), &ids, "id", "i", []turtlebunny.Uint128{}, "ids")
	cmd.MarkFlagRequired("id")

	return cmd
}
