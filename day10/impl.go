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
	"time"
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
		presses, ok := minPressesForJoltageRequirements(m, verbose)
		if !ok {
			return 0, fmt.Errorf("could not determine button presses for %+v", m)
		}
		res += presses
	}
	return res, nil
}

func minPressesForJoltageRequirements(m *Machine, verbose bool) (int, bool) {
	type comboCount struct {
		index        int
		combinations int
	}
	var combos []comboCount
	buttonsByLevel := make([][]Vector, m.width)
	if verbose {
		fmt.Printf("Trying to meet the following joltage requirements: %s\n", m.joltageRequirements)
	}
	for i := 0; i < m.width; i++ {
		if m.joltageRequirements[i] == 0 {
			continue
		}
		// Count the number of buttons wired to this level.
	Buttons:
		for _, b := range m.buttons2 {
			if b[i] == 0 {
				continue Buttons
			}
			for j, v := range b {
				if v > m.joltageRequirements[j] {
					continue Buttons
				}
			}
			buttonsByLevel[i] = append(buttonsByLevel[i], b)
		}
		// Calculate how many combinations of button presses would lead to the desired joltage level.
		if len(buttonsByLevel[i]) == 0 {
			return 0, false
		}
		combinations := 1
		for b := 2; b <= len(buttonsByLevel[i]); b++ {
			combinations *= (m.joltageRequirements[i] + b - 1)
			combinations /= (b - 1)
		}
		if verbose {
			fmt.Printf(
				"Joltage requirement %d is set to %d and is wired to %d buttons. So there are %d combinations for it.\n",
				i, m.joltageRequirements[i], len(buttonsByLevel[i]), combinations)
		}
		combos = append(combos, comboCount{
			index:        i,
			combinations: combinations,
		})
	}
	// If there's no more combos, we're done.
	if len(combos) == 0 {
		return 0, true
	}
	// Pick the joltage level with the fewest options to set it.
	slices.SortFunc(combos, func(a, b comboCount) int {
		return a.combinations - b.combinations
	})
	i := combos[0].index
	buttons := buttonsByLevel[i]
	if verbose {
		fmt.Printf("Picked joltage requirement %d for which there are %d combinations with %d buttons\n", i, combos[0].combinations, len(buttons))
	}
	time.Sleep(0 * time.Second)
	// Iterate all options to reach the joltage level at index i.
	var found bool
	var minPresses int
Candidates:
	for p := range Partitions(m.joltageRequirements[i], len(buttons)) {
		if verbose {
			fmt.Printf("Pressing the following buttons:")
			for k, c := range p {
				fmt.Printf(" %dx % v", c, buttons[k])
			}
			fmt.Println()
		}
		mPrime := m.Clone()
		for b, times := range p {
			mPrime.joltageRequirements = mPrime.joltageRequirements.Add(buttons[b].Scale(-times))
		}
		// Check that this option is valid and didn't overshoot on any other joltage requirements.
		for _, jr := range mPrime.joltageRequirements {
			if jr < 0 {
				continue Candidates
			}
		}
		// Recursively run this algorithm on the remaining unmet joltage requirements.
		presses, ok := minPressesForJoltageRequirements(mPrime, verbose)
		if !ok {
			continue
		}
		if !found || presses+m.joltageRequirements[i] < minPresses {
			found = true
			minPresses = presses + m.joltageRequirements[i]
		}
	}
	return minPresses, found
}

type Bitset int64

func (b Bitset) Set(i int) Bitset {
	return b | (1 << i)
}

func (b Bitset) Has(i int) bool {
	return (b & (1 << i)) != 0
}

func (b Bitset) Xor(c Bitset) Bitset {
	return b ^ c
}

func (b Bitset) Cardinality() int {
	var c int
	for b != 0 {
		if b&1 == 1 {
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

func (v Vector) Scale(scalar int) Vector {
	res := make([]int, len(v))
	for i := 0; i < len(v); i++ {
		res[i] = v[i] * scalar
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
	buttons2            []Vector
	joltageRequirements Vector
}

func (m *Machine) Clone() *Machine {
	return &Machine{
		width:               m.width,
		lights:              m.lights,
		buttons:             m.buttons,
		buttons2:            m.buttons2,
		joltageRequirements: m.joltageRequirements,
	}
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

func Partitions(sum int, parts int) iter.Seq[[]int] {
	if parts == 1 {
		return func(yield func([]int) bool) {
			yield([]int{sum})
		}
	}
	return func(yield func([]int) bool) {
		for i := 0; i <= sum; i++ {
			for p := range Partitions(sum-i, parts-1) {
				if !yield(append(p[:], i)) {
					return
				}
			}
		}
	}
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
