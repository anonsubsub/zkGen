package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func ParseGadgetLibCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gadgetlib-list",
		Short: "parse gadget library",
		Run: func(cmd *cobra.Command, args []string) {

			// read folder
			files, err := os.ReadDir("./gadgetlib")
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
