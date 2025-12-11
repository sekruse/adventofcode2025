package day11

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"
	"time"
)

const (
	you string = "you"
	out string = "out"
	svr string = "svr"
	dac string = "dac"
	fft string = "fft"
)

func Round1(path string, verbose bool) (int, error) {
	graph, err := LoadGraph(path)
	if err != nil {
		return 0, err
	}
	// Traverse the graph.
	var pathsCount int
	leads := []string{you}
	for len(leads) > 0 {
		src := leads[len(leads)-1]
		if verbose {
			fmt.Printf("Following %q:", src)
		}
		leads = leads[:len(leads)-1]
		for _, dst := range graph[src] {
			if verbose {
				fmt.Printf(" -> %q", dst)
			}
			if dst == out {
				pathsCount++
				continue
			}
			leads = append(leads, dst)
		}
		if verbose {
			fmt.Println()
		}
	}
	return pathsCount, nil
}

func Round2(path string, verbose bool) (int, error) {
	graph, err := LoadGraph(path)
	if err != nil {
		return 0, err
	}
	// Traverse the graph.
	return traverse(graph, svr, false, false, make(map[cacheKey]int), verbose), nil
}

type cacheKey struct {
	node                   string
	visitedDAC, visitedFFT bool
}

func traverse(graph map[string][]string, src string, visitedDAC, visitedFFT bool, cache map[cacheKey]int, verbose bool) int {
	k := cacheKey{
		node:       src,
		visitedFFT: visitedFFT,
		visitedDAC: visitedDAC,
	}
	if cached, ok := cache[k]; ok {
		if verbose {
			fmt.Printf("cache[%+v] -> %d\n", k, cached)
		}
		return cached
	}
	var res int
	for _, dst := range graph[src] {
		if verbose {
			fmt.Printf("Exploring %q -> %q\n", src, dst)
		}
		switch dst {
		case out:
			if visitedDAC && visitedFFT {
				res++
			}
		default:
			res += traverse(graph, dst, visitedDAC || dst == dac, visitedFFT || dst == fft, cache, verbose)
		}
	}
	if verbose {
		fmt.Printf("cache[%+v] <- %d\n", k, res)
	}
	cache[k] = res
	return res
}

func Round2a(path string, verbose bool) (int, error) {
	graph, err := LoadGraph(path)
	if err != nil {
		return 0, err
	}
	// Traverse the graph.
	var pathsCount int
	leads := []*lead{{next: svr}}
	for len(leads) > 0 {
		l := leads[len(leads)-1]
		if verbose {
			fmt.Printf("Stack size %d. Following %s ->", len(leads), l)
		}
		leads = leads[:len(leads)-1]
		for _, dst := range graph[l.next] {
			if verbose {
				fmt.Printf(" %q", dst)
			}
			if dst == out {
				if l.dac && l.fft {
					pathsCount++
				}
				continue
			}
			nextLead := &lead{
				path: append(slices.Clone(l.path), l.next),
				next: dst,
				dac:  l.dac || dst == dac,
				fft:  l.fft || dst == fft,
			}
			leads = append(leads, nextLead)
		}
		if verbose {
			fmt.Println()
		}
		time.Sleep(1 * time.Millisecond)
	}
	return pathsCount, nil
}

type lead struct {
	next     string
	dac, fft bool
	path     []string
}

func (l *lead) String() string {
	return fmt.Sprintf("%s>>>%s", strings.Join(l.path, " > "), l.next)

}

func LoadGraph(path string) (map[string][]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	graph := make(map[string][]string)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		tokens := strings.Split(line, ": ")
		if len(tokens) != 2 {
			return nil, fmt.Errorf("cannot parse line: %q", line)
		}
		key := tokens[0]
		vals := strings.Split(tokens[1], " ")
		graph[key] = vals
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return graph, nil
}
