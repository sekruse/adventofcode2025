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
	for _, m := range machines {
		if verbose {
			fmt.Printf("Trying to meet the following joltage requirements: %s\n", m.joltageRequirements)
		}
		slices.SortFunc(m.buttons, func(a, b Bitset) int { return b.Cardinality() - a.Cardinality() })
		presses, ok := meetJoltageRequirements(m, 0, NewVector(m.width), m.joltageRequirements, NewVector(len(m.buttons)), verbose)
		if !ok {
			return 0, fmt.Errorf("could not find button presses")
		}
		res += presses
		break
	}
	return res, nil
}

func meetJoltageRequirements(m *Machine, button int, v, req, pv Vector, verbose bool) (int, bool) {
	if button >= len(m.buttons) {
		return 0, false
	}
	var buttonPresses int
	var totalPresses []int
	for {
		fmt.Printf("Presses: %s, joltage: %s\n", pv, v)
		eq, oob := v.Compare(req)
		if eq {
			totalPresses = append(totalPresses, buttonPresses)
		}
		if eq || oob {
			break
		}
		p, ok := meetJoltageRequirements(m, button + 1, v, req, pv, verbose)
		if ok {
			totalPresses = append(totalPresses, p + buttonPresses)
		}
		buttonPresses++
		pv[button]++
		v = v.Add(m.buttons2[button])
	}
	if len(totalPresses) == 0 {
		return 0, false
	}
	slices.Sort(totalPresses)
	return totalPresses[0], true
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
