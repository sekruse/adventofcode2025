package cmd

import (
	"fmt"

	"github.com/sekruse/adventofcode2025/day04"
	"github.com/spf13/cobra"
)

var day4Round1Cmd = &cobra.Command{
	Use:   "d04r1",
	Short: "Part 1 of day 4.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		res, err := day04.Round1(args[0], verbose)
		if err != nil {
			return err
		}
		fmt.Println(res)
		return nil
	},
}

var day4Round2Cmd = &cobra.Command{
	Use:   "d04r2",
	Short: "Part 2 of day 4.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		res, err := day04.Round2(args[0], verbose)
		if err != nil {
			return err
		}
		fmt.Println(res)
		return nil
	},
}
