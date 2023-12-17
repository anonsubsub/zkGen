package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	// "log"
	cmds "transpiler/commands"
)

func Commands() *cobra.Command {

	// create new cobra command
	cmd := &cobra.Command{
		Use:   "zkGen",
		Short: "\nWelcome,\n\nzkGen is a command-line tool to transpile the zkGen json DSL into secure computation circuits.\n",
	}

	// version command
	cmd.AddCommand(newVersionCommand())

	// zkpolicy
	cmd.AddCommand(cmds.PolicyGetCommand())

	// transpiler
	cmd.AddCommand(cmds.PolicyTranspileCommand())
	cmd.AddCommand(cmds.PolicyListCommand())

	// testing
	cmd.AddCommand(cmds.TestCircuitCommand()) // solidity flag to generate solidity code

	// gadgetlib
	cmd.AddCommand(cmds.ParseGadgetLibCommand())

	return cmd
}

func newVersionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version number of zkGen CMD toolkit.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("zkGen version 0.1")
		},
	}

	return cmd
}

func main() {

	// start command-line toolkit
	cmd := Commands()
	if err := cmd.Execute(); err != nil {
		os.Exit(0)
	}
}
