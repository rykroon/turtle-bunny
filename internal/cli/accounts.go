package cli

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/rykroon/turtlebunny"
	"github.com/spf13/cobra"
)

func NewCreateAccountCmd() *cobra.Command {
	params := &turtlebunny.CreateAccountParams{}

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

	cmd.Flags().Uint64VarP(&params.Id, "id", "i", 0, "id")
	cmd.Flags().Uint32VarP(&params.Ledger, "ledger", "l", 0, "ledger")
	cmd.Flags().Uint16VarP(&params.Code, "code", "c", 0, "code")
	cmd.Flags().BoolVar(&params.DebitsMustNotExceedCredits, "debits-must-not-exceed-credits", false, "debits must not exceed credits")
	cmd.Flags().BoolVar(&params.CreditsMustNotExceedDebits, "credits-must-not-exceed-debits", false, "credits must not exceed debits")

	cmd.MarkFlagRequired("id")
	cmd.MarkFlagRequired("ledger")
	cmd.MarkFlagRequired("code")
	cmd.MarkFlagsMutuallyExclusive("debits-must-not-exceed-credits", "credits-must-not-exceed-debits")

	return cmd
}

func NewLookupAccountCmd() *cobra.Command {
	var id uint64 = 0

	cmd := &cobra.Command{
		Use:   "lookup-account",
		Short: "lookup account",
		Args:  requireFilenameArg,
		RunE: func(cmd *cobra.Command, args []string) error {
			filename := args[0]
			client, err := turtlebunny.NewClient(filename)
			if err != nil {
				return err
			}
			defer client.Close()
			account, err := client.LookupAccount(id)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return errors.New("account not found")
				}
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
			fmt.Printf(
				"%15d %15d %15d %7d %5d %20d\n",
				account.Id,
				account.DebitsPosted,
				account.CreditsPosted,
				account.Ledger,
				account.Code,
				account.Timestamp,
			)
			return nil
		},
	}

	cmd.Flags().Uint64VarP(&id, "id", "i", 0, "id")
	cmd.MarkFlagRequired("id")
	return cmd
}
