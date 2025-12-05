package day05

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"github.com/sekruse/adventofcode2025/day02"
)

func Round1(path string, verbose bool) (int, error) {
	freshProductIDs, productIDs, err := LoadInput(path)
	if err != nil {
		return 0, err
	}
	var freshProductsCount int
	for _, id := range productIDs {
		for _, freshIDs := range freshProductIDs {
			if id >= freshIDs.A && id <= freshIDs.B {
				freshProductsCount++
				break
			}
		}
	}
	return freshProductsCount, nil
}

func Round2(path string, verbose bool) (int64, error) {
	freshProductIDs, _, err := LoadInput(path)
	if err != nil {
		return 0, err
	}
	var stagedFreshProductIDs []*day02.Interval
	for _, origInterval := range freshProductIDs {
		mergedInterval := origInterval
		for i := 0; i < len(stagedFreshProductIDs); {
			stagedInterval := stagedFreshProductIDs[i]
			if origInterval.A > stagedInterval.B || stagedInterval.A > origInterval.B {
				i++
				continue
			}
			if verbose {
				fmt.Printf("Merging %s and %s: ", mergedInterval, stagedInterval)
			}
			mergedInterval = &day02.Interval{
				A: smallest(mergedInterval.A, stagedInterval.A),
				B: greatest(mergedInterval.B, stagedInterval.B),
			}
			if verbose {
				fmt.Printf("%s\n", mergedInterval)
			}
			if i < len(stagedFreshProductIDs)-1 {
				stagedFreshProductIDs[i] = stagedFreshProductIDs[len(stagedFreshProductIDs)-1]
			}
			stagedFreshProductIDs = stagedFreshProductIDs[:len(stagedFreshProductIDs)-1]
		}
		stagedFreshProductIDs = append(stagedFreshProductIDs, mergedInterval)
	}
	var freshProductIDCount int64
	for _, i := range stagedFreshProductIDs {
		freshProductIDCount += i.B - i.A + 1
	}
	return freshProductIDCount, nil
}

func smallest(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func greatest(a, b int64) int64 {
	if a < b {
		return b
	}
	return a
}

func LoadInput(path string) (freshIntervals []*day02.Interval, productIDs []int64, err error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		}
		i, err := day02.ParseInterval(line)
		if err != nil {
			return nil, nil, err
		}
		freshIntervals = append(freshIntervals, i)
	}
	for scanner.Scan() {
		line := scanner.Text()
		id, err := strconv.ParseInt(line, 10, 64)
		if err != nil {
			return nil, nil, err
		}
		productIDs = append(productIDs, id)
	}
	if err := scanner.Err(); err != nil {
		return nil, nil, err
	}
	return freshIntervals, productIDs, nil
}
