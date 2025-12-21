package cli

import (
	"github.com/rykroon/turtlebunny"
	"github.com/spf13/cobra"
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

	cmd.Flags().Uint64VarP(&params.Id, "id", "i", 0, "id")
	cmd.Flags().Uint32VarP(&params.Ledger, "ledger", "l", 0, "ledger")
	cmd.Flags().Uint16VarP(&params.Code, "code", "c", 0, "code")
	cmd.Flags().BoolVar(&params.DebitsMustNotExceedCredits, "debits-must-not-exceed-credits", false, "debits must not exceed credits")
	cmd.Flags().BoolVar(&params.CreditsMustNotExceedDebits, "credits-must-not-exceed-debits", false, "credits must not exceed debits")

	return cmd
}
