package cmd

import (
	"fmt"

	"github.com/sekruse/adventofcode2025/day08"
	"github.com/spf13/cobra"
)

var day8Round1Cmd = &cobra.Command{
	Use:   "d08r1",
	Short: "Part 1 of day 8.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		res, err := day08.Round1(args[0], 1000, 3, verbose)
		if err != nil {
			return err
		}
		fmt.Println(res)
		return nil
	},
}

var day8Round2Cmd = &cobra.Command{
	Use:   "d08r2",
	Short: "Part 2 of day 8.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		res, err := day08.Round2(args[0], verbose)
		if err != nil {
			return err
		}
		fmt.Println(res)
		return nil
	},
}
