package cmd

import (
	"fmt"

	"github.com/sekruse/adventofcode2025/day05"
	"github.com/spf13/cobra"
)

var day5Round1Cmd = &cobra.Command{
	Use:   "d05r1",
	Short: "Part 1 of day 5.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		res, err := day05.Round1(args[0], verbose)
		if err != nil {
			return err
		}
		fmt.Println(res)
		return nil
	},
}

var day5Round2Cmd = &cobra.Command{
	Use:   "d05r2",
	Short: "Part 2 of day 5.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		res, err := day05.Round2(args[0], verbose)
		if err != nil {
			return err
		}
		fmt.Println(res)
		return nil
	},
}
