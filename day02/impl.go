package day02

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	dialInit      int = 50
	dialPositions int = 100
)

func Round1(path string, verbose bool) (int64, error) {
	intervals, err := LoadIntervals(path)
	if err != nil {
		return 0, err
	}
	var sum int64
	for _, i := range intervals {
		// Any number with an odd number of digits cannot be a repeated pattern.
		low, lowLen, err := nextEvenDigits(i.A, false)
		if err != nil {
			return 0, fmt.Errorf(`could not "round up" %d to an even-digited number`, i.A)
		}
		high, highLen, err := nextEvenDigits(i.B, true)
		if err != nil {
			fmt.Fprintf(os.Stderr, `could not "round down" %d to an even-digited number`, i.B)
			continue
		}
		if verbose {
			fmt.Printf("Clamped %s down to %d-%d\n", i, low, high)
		}
		if high < low {
			if verbose {
				fmt.Println("Skipping empty interval")
			}
			continue
		}
		if lowLen != highLen {
			return 0, fmt.Errorf("cannot find invalid digits in intervals with different orders of magnitude: %s", i)
		}
		// Narrowed down to an interval of 2*n digits, we can now simply consider all n-prefixes in the
		// interval and check if the corresponding repeated number is still in the interval.
		var suffixSum int64
		lowPrefix, lowSuffix := split(low)
		highPrefix, highSuffix := split(high)
		if verbose {
			fmt.Printf("Decomposition: %d/%d - %d/%d\n", lowPrefix, lowSuffix, highPrefix, highSuffix)
		}
		// Check if lowPrefix, when repeated, is in the interval.
		if lowPrefix >= lowSuffix && (lowPrefix < highPrefix || lowPrefix <= highSuffix) {
			suffixSum += lowPrefix
			if verbose {
				fmt.Printf("Found invalid ID at the start of the interval: %d %d\n", lowPrefix, lowPrefix)
			}
		}
		// Check if highPrefix, when repeated, is in the interval.
		if highPrefix > lowPrefix && highPrefix <= highSuffix {
			suffixSum += highPrefix
			if verbose {
				fmt.Printf("Found invalid ID at the end of the interval: %d %d\n", highPrefix, highPrefix)
			}
		}
		// Every prefix between lowPrefix and highPrefix yields a repeated pattern.
		// We can use Gauss's famous trick to compute their sum.
		if lowPrefix < highPrefix-1 {
			suffixSum += highPrefix*(highPrefix-1)/2 - lowPrefix*(lowPrefix+1)/2
			if verbose {
				fmt.Printf("Found sum of invalid IDs in the middle of the intervals: %d %d\n", suffixSum, suffixSum)
			}
		}
		// Because all suffixes have the same length, we can postpone and bundle the prefix calculation.
		sum += suffixSum + exp(lowLen/2)*suffixSum
	}
	return sum, nil
}

func nextEvenDigits(n int64, down bool) (res int64, length int, err error) {
	if n <= 0 {
		return 0, 0, fmt.Errorf("cannot make an even-digited number of %d", n)
	}
	s := fmt.Sprintf("%d", n)
	if len(s)%2 == 0 {
		return n, len(s), nil
	}
	base := exp(len(s) - 1)
	if down {
		if base == 1 {
			return 0, 0, fmt.Errorf("there is no even-digited number below %q", n)
		}
		return base - 1, len(s) - 1, nil
	}
	return base * 10, len(s) + 1, nil
}

func exp(k int) int64 {
	if k < 0 {
		panic(fmt.Sprintf("cannot compute 10^%d as an integer", k))
	}
	var res int64 = 1
	for i := 1; i <= k; i++ {
		res *= 10
	}
	return res
}

func split(n int64) (prefix, suffix int64) {
	s := fmt.Sprintf("%d", n)
	prefix, _ = strconv.ParseInt(s[:len(s)/2], 10, 64)
	suffix, _ = strconv.ParseInt(s[len(s)/2:], 10, 64)
	return prefix, suffix
}

func Round2(path string, verbose bool) (int64, error) {
	intervals, err := LoadIntervals(path)
	if err != nil {
		return 0, err
	}
	var sum int64
	for _, interval := range intervals {
		elis := splitIntoEquilengthIntervals(interval)
		if verbose {
			fmt.Printf("Split %s into: %s\n", interval, elis)
		}
		for _, eli := range elis {
			// Step 2: Iterate all possible token lengths (divisible by interval length).
			width := len(fmt.Sprintf("%d", eli.A))
			// Deduplicate patterns, such as 2222 that would be found at token lengths 1 and 2
			invalidIDs := make(map[int64]struct{})
			for tokenLen := 1; tokenLen <= width/2; tokenLen++ {
				if width%tokenLen != 0 {
					continue
				}
				// Step 3: Test the lowest possible invalid ID, e.g., starting from 123456 and token length 2, test 12 12 12.
				tokenA := eli.A / exp(width-tokenLen)
				patternA := createPattern(tokenA, tokenLen, width/tokenLen)
				if patternA >= eli.A && patternA <= eli.B {
					if verbose {
						fmt.Printf("Adding start pattern in %s for token length %d: %d times %d = %d\n", eli, tokenLen, width/tokenLen, tokenA, patternA)
					}
					invalidIDs[patternA] = struct{}{}
				}
				// Step 4: Test the lowest possible invalid ID, e.g., ending at 153344 and token length 2, test 15 15 15
				tokenB := eli.B / exp(width-tokenLen)
				patternB := createPattern(tokenB, tokenLen, width/tokenLen)
				if tokenA < tokenB && patternB >= eli.A && patternB <= eli.B {
					if verbose {
						fmt.Printf("End pattern in %s for token length %d: %d times %d = %d\n", eli, tokenLen, width/tokenLen, tokenB, patternB)
					}
					invalidIDs[patternB] = struct{}{}
				}
				// Step 5: Collect all possible invalid IDs in between the two above, e.g., 13 13 13 and 14 14 14.
				if tokenB-tokenA > 1 {
					if verbose {
						fmt.Printf("Adding %d more inner patterns between %d and %d\n", tokenB-tokenA-1, patternA, patternB)
					}
					for innerToken := tokenA + 1; innerToken < tokenB; innerToken++ {
						invalidIDs[createPattern(innerToken, tokenLen, width/tokenLen)] = struct{}{}
					}
				}
			}
			// Commit the IDs to the sum after the entire interval has been processed.
			for invalidID := range invalidIDs {
				sum += invalidID
			}
		}
	}
	return sum, nil
}

// splitIntoEquilengthIntervals splits the given interval into smallest number of consecutive intervals
// such that the start and end of the interval have the same number of digits.
func splitIntoEquilengthIntervals(in *Interval) []*Interval {
	sa := fmt.Sprintf("%d", in.A)
	sb := fmt.Sprintf("%d", in.B)
	if len(sa) == len(sb) {
		return []*Interval{in}
	}
	var res []*Interval
	start := in.A
	end := exp(len(sa)) - 1
	for start <= in.B {
		res = append(res, &Interval{A: start, B: end})
		start = end + 1
		end = start * 10
		if end > in.B {
			end = in.B
		}
	}
	return res
}

func createPattern(token int64, tokenLen, reps int) int64 {
	repFactor := exp(tokenLen)
	var res int64
	for i := 0; i < reps; i++ {
		res = repFactor*res + token
	}
	return res
}

type Interval struct {
	A, B int64
}

func (s *Interval) String() string {
	return fmt.Sprintf("%d-%d", s.A, s.B)
}

func ParseInterval(code string) (*Interval, error) {
	var res Interval
	literals := strings.Split(code, "-")
	if len(literals) != 2 {
		return nil, fmt.Errorf("unexpected interval format: %q", code)
	}
	var err error
	res.A, err = strconv.ParseInt(literals[0], 10, 64)
	if err != nil {
		return nil, err
	}
	res.B, err = strconv.ParseInt(literals[1], 10, 64)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func LoadIntervals(path string) ([]*Interval, error) {
	var res []*Interval
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	codes := strings.Split(strings.TrimSpace(string(data)), ",")
	for _, code := range codes {
		i, err := ParseInterval(code) // Println will add back the final '\n'
		if err != nil {
			return nil, err
		}
		res = append(res, i)
	}
	return res, nil
}
