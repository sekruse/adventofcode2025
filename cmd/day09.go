package cmd

import (
	"fmt"

	"github.com/sekruse/adventofcode2025/day09"
	"github.com/spf13/cobra"
)

var day9Round1Cmd = &cobra.Command{
	Use:   "d09r1",
	Short: "Part 1 of day 9.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		res, err := day09.Round1(args[0], verbose)
		if err != nil {
			return err
		}
		fmt.Println(res)
		return nil
	},
}

var day9Round2Cmd = &cobra.Command{
	Use:   "d09r2",
	Short: "Part 2 of day 9.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		res, err := day09.Round2(args[0], verbose)
		if err != nil {
			return err
		}
		fmt.Println(res)
		return nil
	},
}
