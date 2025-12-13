package day12

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/gdamore/tcell/v3"
)

var (
	defStyle    = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	shapeStyles = []tcell.Style{
		defStyle.Background(tcell.ColorBlue).Foreground(tcell.ColorWhite),
		defStyle.Background(tcell.ColorRed).Foreground(tcell.ColorWhite),
		defStyle.Background(tcell.ColorGreen).Foreground(tcell.ColorWhite),
		defStyle.Background(tcell.ColorViolet).Foreground(tcell.ColorWhite),
		defStyle.Background(tcell.ColorYellow).Foreground(tcell.ColorBlack),
		defStyle.Background(tcell.ColorWhite).Foreground(tcell.ColorBlack),
	}
	shapeMarkers = []rune{'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z'}
)

func Round1(path string, verbose bool) (int, error) {
	shapes, regions, err := LoadData(path)
	if err != nil {
		return 0, err
	}
	screen, err := tcell.NewScreen()
	if err != nil {
		return 0, err
	}
	defer screen.Fini()
	if err := screen.Init(); err != nil {
		return 0, err
	}
	stop := make(chan struct{})
	go func() {
		for {
			ev := <-screen.EventQ()
			switch ev := ev.(type) {
			case *tcell.EventKey:
				if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
					stop <- struct{}{}
					close(stop)
					return
				}
			}
		}
	}()
	var res int
	for _, reg := range regions {
		if solve(reg, shapes, screen, stop) {
			res++
		}
		select {
		case <-stop:
			return 0, fmt.Errorf("interrupted")
		case <-time.After(100 * time.Millisecond):
		}
	}
	return res, nil
}

func solve(r *Region, shapes []*Shape, screen tcell.Screen, stop <-chan struct{}) bool {
	var shapeVariants [][]*Shape
	for _, s := range shapes {
		shapeVariants = append(shapeVariants, s.Variants())
	}
	attempt := NewAttempt(r)
	var round int
	for {
		round++
		// List shapes that could be added.
		var pendingShapes []*Shape
		for i, c := range attempt.shapeCounts {
			if c > 0 {
				pendingShapes = append(pendingShapes, shapeVariants[i]...)
			}
		}
		if len(pendingShapes) == 0 {
			return true
		}
		// Find possible placements for the shapes.
		var possiblePlacements []*ShapePlacement
		for _, shape := range pendingShapes {
			Rows:
			for x := 0; x < r.width; x++ {
				for y := 0; y < r.height; y++ {
					sp := attempt.Place(shape, x, y)
					if sp == nil {
						continue
					}
					// if !drawAndPause(attempt, fmt.Sprintf("Round %d", round), screen, 100 * time.Millisecond, stop) {
					// 	return false
					// }
					attempt.Remove()
					possiblePlacements = append(possiblePlacements, sp)
					if sp.overlap > 0 {
						continue
					}
					if y == 0 {
						break Rows
					}
					break
				}
			}
		}
		if len(possiblePlacements) == 0 {
			break
		}
		// Pick a placement that maximizes overlap.
		slices.SortFunc(possiblePlacements, func(a, b *ShapePlacement) int {
			return b.overlap - a.overlap
		})
		sp := possiblePlacements[0]
		if test := attempt.Place(sp.shape, sp.x, sp.y); test == nil {
			panic(fmt.Sprintf("failed to place shape %d at (%d, %d)", sp.shape.index, sp.x, sp.y))
		}
		// if !drawAndPause(attempt, fmt.Sprintf("Round %d: Done", round), screen, 0 * time.Millisecond, stop) {
		//	return false
		// }
	}
	if !drawAndPause(attempt, fmt.Sprintf("Round %d: Done", round), screen, 0 * time.Millisecond, stop) {
		return false
	 }
	return false
}

func drawAndPause(attempt *Attempt, status string, screen tcell.Screen, pause time.Duration, stop <-chan struct{}) (ok bool) {
	// Draw the attempt.
	screen.Clear()
	r := attempt.region
	drawCanvas(screen, r.width, r.height, defStyle)
	for i, sp := range attempt.placedShapes {
		drawShape(screen, sp.x, sp.y, sp.shape, shapeMarkers[i%len(shapeMarkers)], shapeStyles[i%len(shapeStyles)])
	}
	drawText(screen, 0, r.height+3, 80, fmt.Sprintf("% d", attempt.shapeCounts), defStyle)
	drawText(screen, 0, r.height+4, 80, status, defStyle)
	screen.Show()
	// Loop handling.
	select {
	case <-stop:
		return false
	case <-time.After(pause):
		return true
	}
}
func drawCanvas(screen tcell.Screen, width, height int, style tcell.Style) {
	screen.Put(0, 0, string(tcell.RuneULCorner), style)
	screen.Put(width+1, 0, string(tcell.RuneURCorner), style)
	screen.Put(0, height+1, string(tcell.RuneLLCorner), style)
	screen.Put(width+1, height+1, string(tcell.RuneLRCorner), style)
	for x := 1; x <= width; x++ {
		screen.Put(x, 0, string(tcell.RuneHLine), style)
		screen.Put(x, height+1, string(tcell.RuneHLine), style)
	}
	for y := 1; y <= height; y++ {
		screen.Put(0, y, string(tcell.RuneVLine), style)
		screen.Put(width+1, y, string(tcell.RuneVLine), style)
	}
}

func drawShape(screen tcell.Screen, ox, oy int, shape *Shape, r rune, style tcell.Style) {
	for y, row := range shape.mask {
		for x, set := range row {
			if set {
				screen.Put(1+ox+x, 1+oy+y, string(r), style)
			}
		}
	}
}

func drawText(screen tcell.Screen, ox, oy, width int, text string, style tcell.Style) {
	for i, r := range text {
		screen.Put(ox+i%width, oy+i/width, string(r), style)
	}
}

type Attempt struct {
	placedShapes []*ShapePlacement
	shapeCounts  []int
	field        [][]int  // counts number of shapes on cell, negative if one shape has a actual tile there
	region       *Region
}

func NewAttempt(r *Region) *Attempt {
	attempt := Attempt{
		region:      r,
		shapeCounts: slices.Clone(r.shapeCounts),
	}
	attempt.field = make([][]int, r.height)
	for y := 0; y < r.height; y++ {
		attempt.field[y] = make([]int, r.width)
	}
	return &attempt
}

// Place attempts to place the given shape at the given position.
func (a *Attempt) Place(s *Shape, x, y int) *ShapePlacement {
	// Check shape supply.
	if a.shapeCounts[s.index] == 0 {
		return nil
	}
	// Check bounds.
	if x < 0 || y < 0 || x+s.width > a.region.width || y+s.height > a.region.height {
		return nil
	}
	// Check overlap.
	for sy, row := range s.mask {
		for sx, isSet := range row {
			if isSet && a.field[y+sy][x+sx] < 0{
				return nil
			}
		}
	}
	// Commit the shape.
	sp := ShapePlacement{
		shape: s,
		x:     x,
		y:     y,
	}
	for sy, row := range s.mask {
		for sx, isSet := range row {
			cell := a.field[y+sy][x+sx]
			if cell != 0 {
				sp.overlap++
			}
			a.field[y+sy][x+sx] = cell + 1
			if isSet {
				a.field[y+sy][x+sx] *= -1
			}
		}
	}
	a.placedShapes = append(a.placedShapes, &sp)
	a.shapeCounts[s.index]--
	return &sp
}

// Remove removes the most recently placed shape.
func (a *Attempt) Remove() bool {
	if len(a.placedShapes) == 0 {
		return false
	}
	sp := a.placedShapes[len(a.placedShapes)-1]
	a.placedShapes = a.placedShapes[:len(a.placedShapes)-1]
	for sy, row := range sp.shape.mask {
		for sx, isSet := range row {
			if isSet {
				a.field[sp.y+sy][sp.x+sx] *= -1
			}
			a.field[sp.y+sy][sp.x+sx] -= 1
		}
	}
	a.shapeCounts[sp.shape.index]++
	return true
}

type ShapePlacement struct {
	shape *Shape
	x, y  int
	overlap int
}

type Shape struct {
	mask          [][]bool
	width, height int
	index         int
}

func (s *Shape) Equals(t *Shape) bool {
	if s.index!= t.index || s.width != t.width || s.height != t.height {
		return false
	}
	for y, row := range s.mask {
		for x, cell := range row {
			if cell != t.mask[y][x] {
				return false
			}
		}
	}
	return true
}

func (s *Shape) FlipLR() *Shape {
	t := &Shape{
		width:  s.width,
		height: s.height,
		index:  s.index,
	}
	t.mask = make([][]bool, t.height)
	for y, row := range s.mask {
		t.mask[y] = make([]bool, t.width)
		for x, cell := range row {
			t.mask[y][t.width-x-1] = cell
		}
	}
	return t
}

func (s *Shape) RotateCW() *Shape {
	t := &Shape{
		width:  s.height,
		height: s.width,
		index:  s.index,
	}
	t.mask = make([][]bool, t.height)
	for x := 0; x < s.width; x++ {
		t.mask[x] = make([]bool, t.width)
		for y := 0; y < s.height; y++ {
			t.mask[x][s.height-y-1] = s.mask[y][x]
		}
	}
	return t
}

func (s *Shape) Variants() []*Shape {
	res := []*Shape{s}
	// Rotate shape. A rotational symmetry of 90 degrees implies a rotational symmetry of 180 degrees.
	t := s
	for i := 0; i < 3; i++ {
		t = t.RotateCW()
		if t.Equals(s) {
			break
		}
		res = append(res, t)
	}
	// Every flip (horizontally, diagonally, vertically) can be created from a horizontal flip and a rotation.
	t = s.FlipLR()
	for i := 0; i < 4; i++ {
		if !t.Equals(s) {
			res = append(res, t)
		}
		t = t.RotateCW()
	}
	return res
}

type Region struct {
	width, height int
	shapeCounts   []int
}

var (
	regionRE = regexp.MustCompile(`^(\d+)x(\d+):\s+((?:\d+ ?)+)$`)
)

func LoadData(path string) (shapes []*Shape, regions []*Region, err error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// Parse shapes.
		if line == fmt.Sprintf("%d:", len(shapes)) {
			shape := Shape{index: len(shapes)}
			for scanner.Scan() {
				line := scanner.Text()
				if line == "" {
					break
				}
				if shape.width == 0 {
					shape.width = len(line)
				} else if len(line) != shape.width {
					return nil, nil, fmt.Errorf("expected shape width %d, got %q", shape.width, line)
				}
				row := make([]bool, shape.width)
				for i, r := range line {
					row[i] = r == '#'
				}
				shape.mask = append(shape.mask, row)
			}
			shape.height = len(shape.mask)
			shapes = append(shapes, &shape)
			continue
		}
		// Parse requirements.
		var reg Region
		m := regionRE.FindStringSubmatch(line)
		if m == nil {
			return nil, nil, fmt.Errorf("unexpected line: %q", line)
		}
		reg.width, err = strconv.Atoi(m[1])
		if err != nil {
			return nil, nil, err
		}
		reg.height, err = strconv.Atoi(m[2])
		if err != nil {
			return nil, nil, err
		}
		for _, token := range strings.Split(m[3], " ") {
			sc, err := strconv.Atoi(token)
			if err != nil {
				return nil, nil, err
			}
			reg.shapeCounts = append(reg.shapeCounts, sc)
		}
		regions = append(regions, &reg)
	}
	if err := scanner.Err(); err != nil {
		return nil, nil, err
	}
	return shapes, regions, nil
}
