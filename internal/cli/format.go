package cli

import (
	"github.com/rykroon/turtlebunny"
	"github.com/spf13/cobra"
)

func NewFormatCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "format",
		Short: "",
		Args:  requireFilenameArg,
		RunE: func(cmd *cobra.Command, args []string) error {
			filename := args[0]
			client, err := turtlebunny.NewClient(filename)
			if err != nil {
				return err
			}
			defer client.Close()
			err = client.Format()
			if err != nil {
				return err
			}
			return nil
		},
	}

	return cmd
}
