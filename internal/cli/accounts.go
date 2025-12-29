package cli

import (
	"fmt"
	"strings"

	"github.com/rykroon/turtlebunny"
	"github.com/spf13/cobra"
	"lukechampine.com/uint128"
)

func NewCreateAccountCmd() *cobra.Command {
	params := turtlebunny.CreateAccountParams{}

	cmd := &cobra.Command{
		Use:   "create-account",
		Short: "create a new account",
		Args:  requireFilenameArg,
		RunE: func(cmd *cobra.Command, args []string) error {
			filename := args[0]
			client, err := turtlebunny.NewClient(filename)
			if err != nil {
				return err
			}
			defer client.Close()
			err = client.CreateAccount(params)
			if err != nil {
				return err
			}
			return nil
		},
	}

	cmd.Flags().VarP(NewUint128Flag(&params.Id), "id", "i", "id")
	cmd.Flags().Uint32VarP(&params.Ledger, "ledger", "l", 0, "ledger")
	cmd.Flags().Uint16VarP(&params.Code, "code", "c", 0, "code")
	cmd.Flags().BoolVar(
		&params.DebitsMustNotExceedCredits,
		"debits-must-not-exceed-credits",
		false,
		"debits must not exceed credits",
	)

	cmd.Flags().BoolVar(
		&params.CreditsMustNotExceedDebits,
		"credits-must-not-exceed-debits",
		false,
		"credits must not exceed debits",
	)

	cmd.MarkFlagRequired("id")
	cmd.MarkFlagRequired("ledger")
	cmd.MarkFlagRequired("code")
	cmd.MarkFlagsMutuallyExclusive(
		"debits-must-not-exceed-credits", "credits-must-not-exceed-debits",
	)

	return cmd
}

func NewLookupAccountCmd() *cobra.Command {
	ids := []uint128.Uint128{}

	cmd := &cobra.Command{
		Use:   "lookup-accounts",
		Short: "lookup accounts",
		Args:  requireFilenameArg,
		RunE: func(cmd *cobra.Command, args []string) error {
			filename := args[0]
			client, err := turtlebunny.NewClient(filename)
			if err != nil {
				return err
			}
			defer client.Close()
			accounts, err := client.LookupAccounts(ids...)
			if err != nil {
				return err
			}

			fmt.Printf(
				"%15s %15s %15s %7s %5s %20s\n",
				"ID",
				"Debits Posted",
				"Credits Posted",
				"Ledger",
				"Code",
				"Timestamp",
			)
			fmt.Println(strings.Repeat("-", 85))

			for _, account := range accounts {
				fmt.Printf(
					"%15s %15s %15s %7d %5d %20d\n",
					account.Id.String(),
					account.DebitsPosted.String(),
					account.CreditsPosted.String(),
					account.Ledger,
					account.Code,
					account.Timestamp,
				)
			}

			return nil
		},
	}

	Uint128SliceVarP(cmd.Flags(), &ids, "id", "i", []uint128.Uint128{}, "ids")
	cmd.MarkFlagRequired("id")
	return cmd
}
