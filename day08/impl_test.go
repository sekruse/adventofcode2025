package day08

import (
	"path/filepath"
	"testing"
)

func TestResults(t *testing.T) {
	testFilePath := filepath.Join("testdata", "example.txt")
	t.Run("Round 1", func(t *testing.T) {
		const want = 40
		got, err := Round1(testFilePath, 10, 3, false)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if got != want {
			t.Errorf("got %d, want %d", got, want)
		}
	})
	t.Run("Round 2", func(t *testing.T) {
		const want = 25272
		got, err := Round2(testFilePath, false)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if got != want {
			t.Errorf("got %d, want %d", got, want)
		}
	})
}
