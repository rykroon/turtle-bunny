package cli

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "turtlebunny",
		Short: "",
		Long:  "",
	}

	cmd.AddCommand(
		NewCreateAccountCmd(),
		NewCreateTransferCmd(),
	)

	return cmd
}

func requireFilenameArg(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return errors.New("missing 1 required positional argument: 'FILENAME'")
	}
	if len(args) > 1 {
		return fmt.Errorf("expected 1 positional argument but %d were given", len(args))
	}
	return nil
}
