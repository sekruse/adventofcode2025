package day10

import (
	"bufio"
	"fmt"
	"iter"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

var (
	lightsRE  = regexp.MustCompile(`\[([.#]*)\]`)
	buttonsRE = regexp.MustCompile(`\(([0-9,]+)\)`)
	joltageRE = regexp.MustCompile(`\{([0-9,]+)\}`)
)

func Round1(path string, verbose bool) (int, error) {
	machines, err := LoadMachines(path)
	if err != nil {
		return 0, err
	}
	var res int
Machines:
	for _, m := range machines {
		if verbose {
			fmt.Printf("Trying to create the following lights: %s\n", m.lights)
		}
		for b := 1; b <= m.width; b++ {
			for c := range Cube(len(m.buttons), b) {
				if verbose {
					fmt.Printf("Pressing: ")
				}
				var lights Bitset
				for _, i := range c {
					if verbose {
						fmt.Printf("%s ", m.buttons[i])
					}
					lights ^= m.buttons[i]
				}
				if verbose {
					fmt.Printf("and got %s.\n", lights)
				}
				if lights == m.lights {
					res += b
					continue Machines
				}
			}
		}
		return 0, fmt.Errorf("could not find button combo for %v", m)
	}
	return res, nil
}

func Round2(path string, verbose bool) (int, error) {
	machines, err := LoadMachines(path)
	if err != nil {
		return 0, err
	}
	var res int
	type state struct {
		presses int
		pressVector Vector
		joltage Vector
		distance float64
	}
	Machines:
	for _, m := range machines {
		if verbose {
			fmt.Printf("Trying to meet the following joltage requirements: %s\n", m.joltageRequirements)
		}
		initState := &state {
				presses: 0,
				pressVector: NewVector(len(m.buttons2)),
				joltage: NewVector(m.width),
			}
		initState.distance, _ = initState.joltage.Distance(m.joltageRequirements)
		states := []*state{initState}
		shortestLeadUps := make(map[string]int)
		for len(states) > 0 {
			// Pop the top state.
			s := states[len(states)-1]
			states = states[:len(states)-1]
			for ib, b := range m.buttons2 {
				next := &state{
					presses: s.presses+1,
					pressVector: slices.Clone(s.pressVector),
					joltage: s.joltage.Add(b),
				}
				next.pressVector[ib] += 1
				dist, ok := next.joltage.Distance(m.joltageRequirements)
				if !ok {
					continue
				}
				if dist == 0 {
					fmt.Printf("Solution for %+v:\n\t%v\n", m, next) 
					res += next.presses
					continue Machines
				}
				key := next.joltage.String()
				slu, ok := shortestLeadUps[key]
				if ok && slu <= next.presses {
					continue
				}
				shortestLeadUps[key] = next.presses
				next.distance = dist
				states = append(states, next)
			}
			slices.SortFunc(states, func(a, b *state) int {
				if b.distance < a.distance { return -1 }
				if b.distance > a.distance { return 1 }
				return 0	
			})
		}
		return 0, fmt.Errorf("did not find any solution for %+v", m)
	}
	return res, nil
}

type Bitset int64

func (b Bitset) Set(i int) Bitset {
	return b | (1 << i)
}

func (b Bitset) Xor(c Bitset) Bitset {
	return b ^ c
}

func (b Bitset) Cardinality() int {
	var c int
	for b != 0 {
		if b & 1 == 1 {
			c++
		}
		b = b >> 1
	}
return c
}

func (b Bitset) String() string {
	return fmt.Sprintf("%b", b)
}

type Vector []int

func minLen(v, w Vector) int {
	if len(v) < len(w) {
		return len(v)
	}
	return len(w)
}

func NewVector(l int) Vector {
	return make([]int, l)
}

func (v Vector) Add(w Vector) Vector {
	l := minLen(v, w)
	res := make([]int, l)
	for i := 0; i < l; i++ {
		res[i] = v[i] + w[i]
	}
	return res
}

func (v Vector) Compare(w Vector) (eq, oob bool) {
	eq = true
	l := minLen(v, w)
	for i := 0; i < l; i++ {
		if v[i] < w[i] {
			eq = false
			continue
		}
		if v[i] > w[i] {
			eq = false
			oob = true
			return eq, oob
		}
	}
	return eq, oob
}

func (v Vector) Distance(w Vector) (dist float64, ok bool) {
	l := minLen(v, w)
	for i := 0; i < l; i++ {
		d := float64(w[i] - v[i])
		if d < 0 {
			return -1, false
		}
		dist += d * d
	}
	return dist, true
}



func (v Vector) String() string {
	return fmt.Sprintf("% d", v)
}

type Machine struct {
	width               int
	lights              Bitset
	buttons             []Bitset
	buttons2 []Vector
	joltageRequirements Vector
}

func ParseMachine(line string) (*Machine, error) {
	var m Machine
	// Parse lights.
	lights := lightsRE.FindStringSubmatch(line)
	if lights == nil {
		return nil, fmt.Errorf("no lights in %q", line)
	}
	m.width = len(lights[1])
	for i, r := range lights[1] {
		if r == '#' {
			m.lights = m.lights.Set(i)
		}
	}
	// Parse buttons.
	buttonSpecs := buttonsRE.FindAllStringSubmatch(line, -1)
	if buttonSpecs == nil {
		return nil, fmt.Errorf("no buttons in %q", line)
	}
	for _, bs := range buttonSpecs {
		var button Bitset
		button2 := NewVector(m.width)
		for _, b := range strings.Split(bs[1], ",") {
			i, err := strconv.Atoi(b)
			if err != nil {
				return nil, fmt.Errorf("not a number in %q in %q", bs, line)
			}
			button = button.Set(i)
			button2[i] = 1
		}
		m.buttons = append(m.buttons, button)
		m.buttons2 = append(m.buttons2, button2)
	}
	// Parse joltage.
	joltageSpec := joltageRE.FindStringSubmatch(line)
	if joltageSpec == nil {
		return nil, fmt.Errorf("no joltage requirements in %q", line)
	}
	m.joltageRequirements = NewVector(m.width)
	for i, js := range strings.Split(joltageSpec[1], ",") {
		j, err := strconv.Atoi(js)
		if err != nil {
			return nil, fmt.Errorf("not a number in %q in %q", js, line)
		}
		m.joltageRequirements[i] = j
	}
	return &m, nil
}

func LoadMachines(path string) ([]*Machine, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	var machines []*Machine
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		m, err := ParseMachine(line)
		if err != nil {
			return nil, err
		}
		machines = append(machines, m)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return machines, nil
}

func Cube(n, dim int) iter.Seq[[]int] {
	if dim > 1 {
		return func(yield func([]int) bool) {
			for v := range Cube(n, dim-1) {
				for k := v[len(v)-1] + 1; k < n; k++ {
					w := v[:]
					if !yield(append(w, k)) {
						return
					}
				}
			}
		}
	}
	return func(yield func([]int) bool) {
		for i := 0; i < n; i++ {
			if !yield([]int{i}) {
				return
			}
		}
	}
}
