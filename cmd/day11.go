package cmd

import (
	"fmt"

	"github.com/sekruse/adventofcode2025/day11"
	"github.com/spf13/cobra"
)

var day11Round1Cmd = &cobra.Command{
	Use:   "d11r1",
	Short: "Part 1 of day 11.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		res, err := day11.Round1(args[0], verbose)
		if err != nil {
			return err
		}
		fmt.Println(res)
		return nil
	},
}

var day11Round2Cmd = &cobra.Command{
	Use:   "d11r2",
	Short: "Part 2 of day 11.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		res, err := day11.Round2(args[0], verbose)
		if err != nil {
			return err
		}
		fmt.Println(res)
		return nil
	},
}
