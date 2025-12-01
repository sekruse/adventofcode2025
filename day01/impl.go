package day01

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

const (
	dialInit      int = 50
	dialPositions int = 100
)

func Round1(path string, verbose bool) (int, error) {
	instructions, err := LoadSafeInstructions(path)
	if err != nil {
		return 0, err
	}
	// Run the simulation.
	counter := 0
	dial := 50
	for _, i := range instructions {
		dial = (dial + int(i.Direction)*i.Clicks + dialPositions) % dialPositions
		if dial == 0 {
			counter++
		}
		if verbose {
			fmt.Printf("Turn: %5s -> Dial: %3d. Counter: %5d.\n", i, dial, counter)
		}
	}
	return counter, nil
}

func Round2(path string, verbose bool) (int, error) {
	instructions, err := LoadSafeInstructions(path)
	if err != nil {
		return 0, err
	}
	// Run the simulation.
	counter := 0
	dial := 50
	for _, i := range instructions {
		// Every full revolution necessarily passes the 0 exactly once.
		fullRevolutions := i.Clicks / dialPositions
		counter += fullRevolutions
		// We check if the remaining motion either crosses the 0 or lands on it.
		remainingClicks := i.Clicks - fullRevolutions*dialPositions
		if remainingClicks > 0 {
			dialWasZero := dial == 0
			overshootDial := dial + int(i.Direction)*remainingClicks
			dial = (overshootDial + dialPositions) % dialPositions
			if dial == 0 || (!dialWasZero && dial != overshootDial) {
				counter++
			}
		}
		if verbose {
			fmt.Printf("Turn: %5s -> Dial: %3d. Counter: %5d.\n", i, dial, counter)
		}
	}
	return counter, nil
}

type Direction int

const (
	Left  Direction = -1
	Right Direction = 1
)

var safeInstructionRegex = regexp.MustCompile(`^([LR])(\d+)$`)

type SafeInstruction struct {
	Direction Direction
	Clicks    int
	code      string
}

func (s *SafeInstruction) String() string {
	return s.code
}

func ParseSafeInstruction(code string) (*SafeInstruction, error) {
	res := SafeInstruction{
		code: code,
	}
	matches := safeInstructionRegex.FindStringSubmatch(code)
	if matches == nil {
		return nil, fmt.Errorf("unexpected safe instruction: %q", code)
	}
	switch matches[1] {
	case "L":
		res.Direction = Left
	case "R":
		res.Direction = Right
	}
	var err error
	res.Clicks, err = strconv.Atoi(matches[2])
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func LoadSafeInstructions(path string) ([]*SafeInstruction, error) {
	var res []*SafeInstruction
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		si, err := ParseSafeInstruction(line) // Println will add back the final '\n'
		if err != nil {
			return nil, err
		}
		res = append(res, si)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return res, nil
}
