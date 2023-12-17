package commands

import (
	"errors"
	t "transpiler/templates"
	p "transpiler/zkpolicy"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func TestCircuitCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "zkpolicy-test",
		Short: "test created circuit via ZKP system (use flag -solidity to generate solidity code for circuit verification)",
		RunE: func(cmd *cobra.Command, args []string) error {

			// check for input arguments
			if len(args) < 1 {
				return errors.New("provide policy filename without extension")
			}
			policyName := args[0]

			solidityFlag := ""
			if len(args) > 1 {
				solidityFlag = args[1]
			}

			// fetch policy
			zkPolicy, err := p.ParseZkPolicy(policyName)
			if err != nil {
				log.Error().Err(err).Msg("p.ParseZkPolicy()")
				return err
			}

			// test circuit solidity
			err = t.TestCircuit(zkPolicy.Name, solidityFlag)
			if err != nil {
				log.Error().Err(err).Msg("t.TestCircuit()")
				return err
			}

			return nil
		},
	}

	return cmd
}
