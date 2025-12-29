package cli

import (
	"fmt"
	"strings"

	"github.com/rykroon/turtlebunny"
	"github.com/spf13/cobra"
	"lukechampine.com/uint128"
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

	cmd.Flags().VarP(NewUint128Flag(&params.Id), "id", "i", "id")
	cmd.Flags().VarP(
		NewUint128Flag(&params.DebitAccountId), "debit-account-id", "D", "debit account id",
	)

	cmd.Flags().VarP(
		NewUint128Flag(&params.CreditAccountId), "credit-account-id", "C", "credit account id",
	)

	cmd.Flags().VarP(NewUint128Flag(&params.Amount), "amount", "a", "amount")
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
	ids := []uint128.Uint128{}

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

	Uint128SliceVarP(cmd.Flags(), &ids, "id", "i", []uint128.Uint128{}, "ids")
	cmd.MarkFlagRequired("id")

	return cmd
}
