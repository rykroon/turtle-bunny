package main

import (
	"fmt"
	"os"

	"github.com/rykroon/turtlebunny/internal/cli"
)

func main() {
	cmd := cli.NewRootCmd()
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
