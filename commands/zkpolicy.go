package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	p "transpiler/zkpolicy"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func PolicyGetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "zkpolicy-get",
		Short: "returns a zkpolicy according to the provided filename.",
		RunE: func(cmd *cobra.Command, args []string) error {

			// check for config name as input argument
			if len(args) < 1 {
				return errors.New("provide policy filename without extension")
			}
			policyName := args[0]

			// fetch policy
			zkPolicy, err := p.ParseZkPolicy(policyName)
			if err != nil {
				log.Error().Err(err).Msg("p.ParseZkPolicy()")
				return err
			}

			// pretty print json string
			s, err := json.MarshalIndent(zkPolicy, "", "\t")
			if err != nil {
				log.Error().Err(err).Msg("json.MashalIndent()")
				return err
			}

			// print to console
			fmt.Print(string(s))
			fmt.Println()

			return nil
		},
	}

	return cmd
}

func PolicyListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "zkpolicy-list",
		Short: "list available zkPolicies.",
		Run: func(cmd *cobra.Command, args []string) {

			// read folder
			files, err := os.ReadDir("./zkpolicy")
			if err != nil {
				log.Error().Err(err).Msg("os.ReadDir()")
			}

			// print filename if not a directory
			for _, file := range files {
				if !file.IsDir() && strings.Contains(file.Name(), ".json") {
					fmt.Println(strings.Split(file.Name(), ".json")[0])
				}
			}
		},
	}

	return cmd
}
