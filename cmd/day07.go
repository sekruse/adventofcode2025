package cmd

import (
	"fmt"

	"github.com/sekruse/adventofcode2025/day06"
	"github.com/spf13/cobra"
)

var day6Round1Cmd = &cobra.Command{
	Use:   "d06r1",
	Short: "Part 1 of day 6.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		res, err := day06.Round1(args[0], verbose)
		if err != nil {
			return err
		}
		fmt.Println(res)
		return nil
	},
}

var day6Round2Cmd = &cobra.Command{
	Use:   "d06r2",
	Short: "Part 2 of day 6.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		res, err := day06.Round2(args[0], verbose)
		if err != nil {
			return err
		}
		fmt.Println(res)
		return nil
	},
}
