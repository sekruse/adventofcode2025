package day03

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func Round1(path string, verbose bool) (int, error) {
	banks, err := LoadBatteryBanks(path)
	if err != nil {
		return 0, err
	}
	var joltage int
	for _, b := range banks {
		j := b.Joltage()
		if verbose {
			fmt.Printf("Joltage for %v: %d\n", b, j)
		}
		joltage += j
	}
	return joltage, nil
}

func Round2(path string, verbose bool) (int64, error) {
	const batteries = 12
	banks, err := LoadBatteryBanks(path)
	if err != nil {
		return 0, err
	}
	var joltage int64
	for _, b := range banks {
		j := b.JoltageN(batteries, verbose)
		if verbose {
			fmt.Printf("Joltage for %v: %d\n", b, j)
		}
		joltage += j
	}
	return joltage, nil
}

type BatteryBank []int

func (b BatteryBank) Joltage() int {
	var pos1, pos2 int
	for i := 1; i < len(b)-1; i++ {
		if b[pos1] < b[i] {
			pos1 = i
		}
	}
	pos2 = pos1 + 1
	for i := pos1 + 2; i < len(b); i++ {
		if b[pos2] < b[i] {
			pos2 = i
		}
	}
	return 10*b[pos1] + b[pos2]
}

func (b BatteryBank) JoltageN(n int, verbose bool) int64 {
	var start int
	var res int64
	for k := 0; k < n; k++ {
		pos := start
		end := len(b) - (n - k - 1)
		for i := start + 1; i < end; i++ {
			if b[pos] < b[i] {
				pos = i
			}
		}
		fmt.Printf("Chose digit %d at %d\n", b[pos], pos)
		res = 10*res + int64(b[pos])
		start = pos + 1
	}
	return res
}

func LoadBatteryBanks(path string) ([]BatteryBank, error) {
	var res []BatteryBank
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
		var bank BatteryBank
		for _, r := range strings.Split(line, "") {
			i, err := strconv.Atoi(r)
			if err != nil {
				return nil, err
			}
			bank = append(bank, i)
		}
		res = append(res, bank)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return res, nil
}
