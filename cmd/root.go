package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   os.Args[0],
		Short: "Solution for Advent of Code 2025 puzzles",
	}

	verbose bool
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.AddCommand(day1Round1Cmd)
	rootCmd.AddCommand(day1Round2Cmd)
}
