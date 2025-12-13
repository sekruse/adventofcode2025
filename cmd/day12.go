package cmd

import (
	"fmt"

	"github.com/sekruse/adventofcode2025/day12"
	"github.com/spf13/cobra"
)

var day12Round1Cmd = &cobra.Command{
	Use:   "d12r1",
	Short: "Part 1 of day 12.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		res, err := day12.Round1(args[0], verbose)
		if err != nil {
			return err
		}
		fmt.Println(res)
		return nil
	},
}
