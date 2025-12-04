package day04

import (
	"bufio"
	"fmt"
	"os"
)

func Round1(path string, verbose bool) (int, error) {
	plan, err := LoadFloorPlan(path)
	if err != nil {
		return 0, err
	}
	var clearPaperRolls int
	for y, row := range plan {
		for x, tile := range row {
			if tile != Tile_PaperRoll {
				continue
			}
			var neighbors int
			for _, offset := range neighborMask {
				x2, y2 := x+offset.X, y+offset.Y
				if x2 < 0 || x2 >= len(row) || y2 < 0 || y2 >= len(plan) {
					continue
				}
				if plan[y2][x2] != Tile_Clear {
					neighbors++
				}
			}
			if neighbors < 4 {
				clearPaperRolls++
			}
		}
	}
	return clearPaperRolls, nil
}

func Round2(path string, verbose bool) (int, error) {
	plan, err := LoadFloorPlan(path)
	if err != nil {
		return 0, err
	}
	var clearPaperRolls int
	for {
		prevClearPaperRolls := clearPaperRolls
		for y, row := range plan {
			for x, tile := range row {
				if tile != Tile_PaperRoll {
					continue
				}
				var neighbors int
				for _, offset := range neighborMask {
					x2, y2 := x+offset.X, y+offset.Y
					if x2 < 0 || x2 >= len(row) || y2 < 0 || y2 >= len(plan) {
						continue
					}
					if plan[y2][x2] != Tile_Clear {
						neighbors++
					}
				}
				if neighbors < 4 {
					clearPaperRolls++
					plan[y][x] = Tile_Clear
				}
			}
		}
		if prevClearPaperRolls == clearPaperRolls {
			break
		}
	}
	return clearPaperRolls, nil
}

type Tile int

const (
	Tile_Clear     Tile = 0
	Tile_PaperRoll Tile = 1
)

type FloorPlan [][]Tile

func LoadFloorPlan(path string) (FloorPlan, error) {
	var plan FloorPlan
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
		if len(plan) > 0 && len(plan[0]) != len(line) {
			return nil, fmt.Errorf("unexpected line length on floor plan: %q", line)
		}
		var row []Tile
		for _, r := range line {
			switch r {
			case '.':
				row = append(row, Tile_Clear)
			case '@':
				row = append(row, Tile_PaperRoll)
			default:
				return nil, fmt.Errorf("unexpected character on floor plan: %q", r)
			}
		}
		plan = append(plan, row)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return plan, nil
}

type Offset struct {
	X, Y int
}

var neighborMask []*Offset = []*Offset{
	{-1, -1},
	{-1, 0},
	{-1, 1},
	{0, -1},
	{0, 1},
	{1, -1},
	{1, 0},
	{1, 1},
}
