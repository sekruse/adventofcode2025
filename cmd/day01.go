package cmd

import (
	"fmt"

	"github.com/sekruse/adventofcode2025/day01"
	"github.com/spf13/cobra"
)

var day1Round1Cmd = &cobra.Command{
	Use:   "d01r1",
	Short: "Part 1 of day 01.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		res, err := day01.Round1(args[0], verbose)
		if err != nil {
			return err
		}
		fmt.Println(res)
		return nil
	},
}

var day1Round2Cmd = &cobra.Command{
	Use:   "d01r2",
	Short: "Part 2 of day 2.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		res, err := day01.Round2(args[0], verbose)
		if err != nil {
			return err
		}
		fmt.Println(res)
		return nil
	},
}
