package cli

import (
	"github.com/rykroon/turtlebunny"
	"github.com/spf13/cobra"
)

func NewCreateTransferCmd() *cobra.Command {
	params := turtlebunny.CreateTransferParams{}

	cmd := &cobra.Command{
		Use:   "create-transfer",
		Short: "create transfer",
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

	return cmd
}
