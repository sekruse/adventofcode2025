package day12

import (
	"path/filepath"
	"testing"
)

func TestResults(t *testing.T) {
	t.Run("Round 1", func(t *testing.T) {
		testFilePath := filepath.Join("testdata", "example.txt")
		const want = 2
		got, err := Round1(testFilePath, false)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if got != want {
			t.Errorf("got %d, want %d", got, want)
		}
	})
}
