package main

import (
	"fmt"
	"os"

	"github.com/sekruse/adventofcode2025/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v", err)
		os.Exit(1)
	}
}
