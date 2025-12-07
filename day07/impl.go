package day07

import (
	"bufio"
	"fmt"
	"os"
)

func Round1(path string, verbose bool) (int, error) {
	plan, err := LoadPlan(path)
	if err != nil {
		return 0, err
	}
	// Horizontal cross-section of the beam.
	var splitCount int
	beamSlice := make([]bool, len(plan.tiles[0]))
	beamSlice[plan.start.X] = true
	for y := plan.start.Y; y < plan.Dim().Y; y++ {
		nextBeamSlice := make([]bool, len(plan.tiles[0]))
		for x := 0; x < plan.Dim().X; x++ {
			if !beamSlice[x] {
				continue
			}
			if plan.tiles[y][x] == Tile_Splitter {
				splitCount++
				if x > 0 {
					nextBeamSlice[x-1] = true
				}
				if x < plan.Dim().X - 1 {
					nextBeamSlice[x+1] = true
				}
				continue
			}
			nextBeamSlice[x] = true
		}
		beamSlice = nextBeamSlice
	}
	return splitCount, nil
}

func Round2(path string, verbose bool) (int, error) {
	plan, err := LoadPlan(path)
	if err != nil {
		return 0, err
	}
	// In contrast to Round1, the slice keeps track of the paths that lead to the beam passing through.
	beamSlice := make([]int, len(plan.tiles[0]))
	beamSlice[plan.start.X] = 1
	for y := plan.start.Y; y < plan.Dim().Y; y++ {
		nextBeamSlice := make([]int, len(plan.tiles[0]))
		for x := 0; x < plan.Dim().X; x++ {
			paths := beamSlice[x]
			if paths == 0 {
				continue
			}
			if plan.tiles[y][x] == Tile_Splitter {
				if x > 0 {
					nextBeamSlice[x-1] += paths
				}
				if x < plan.Dim().X - 1 {
					nextBeamSlice[x+1] += paths
				}
				continue
			}
			nextBeamSlice[x] += paths
		}
		beamSlice = nextBeamSlice
	}
	var paths int
	for _, v := range beamSlice {
		paths += v
	}
	return paths, nil
}

type Tile int

const (
	Tile_Empty Tile = iota
	Tile_Start
	Tile_Splitter
)

func NewTile(char rune) (Tile, error) {
	switch char {
	case '.':
		return Tile_Empty, nil
	case 'S':
		return Tile_Start, nil
	case '^':
		return Tile_Splitter, nil
	default:
		return Tile_Empty, fmt.Errorf("not a valid tile character: %q", char)
	}
}

type Point2D struct {
	X, Y int
}

type Plan struct {
	tiles [][]Tile
	start *Point2D
}

func (p *Plan) Dim() *Point2D {
	dim := Point2D{ Y: len(p.tiles) }
	if dim.Y > 0 {
		dim.X = len(p.tiles[0])
	}
	return &dim
}

func (p *Point2D) String() string {
	return fmt.Sprintf("(%d, %d)", p.X, p.Y)
}

func LoadPlan(path string) (*Plan, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	var plan Plan
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		var row []Tile
		for _, ch := range line {
			t, err := NewTile(ch)
			if err != nil {
				return nil, err
			}
			if t == Tile_Start {
				p := Point2D{
					X: len(row),
					Y: len(plan.tiles),
				}
				if plan.start != nil {
					return nil, fmt.Errorf("plan contains second start point at %s", &p)
				}
				plan.start = &p
			}
			row = append(row, t)
		}
		if dim := plan.Dim(); dim.Y > 0 && dim.X != len(row) {
			return nil, fmt.Errorf("inconsistent plan width: %q", row)
		}
		plan.tiles = append(plan.tiles, row)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return &plan, nil
}
