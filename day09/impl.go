package day09

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/sekruse/adventofcode2025/day07"
)

type Point2D = day07.Point2D

func Round1(path string, verbose bool) (int, error) {
	points, err := LoadPoints(path)
	if err != nil {
		return 0, err
	}
	// Compute pairwise distances of points.
	var maxSquareSize int
	for i := 0; i < len(points)-1; i++ {
		p := points[i]
		for j := i + 1; j < len(points); j++ {
			q := points[j]
			size := (Abs(p.X-q.X) + 1) * (Abs(p.Y-q.Y) + 1)
			if size > maxSquareSize {
				maxSquareSize = size
			}
		}
	}
	return maxSquareSize, nil
}

func Round2(path string, verbose bool) (int, error) {
	points, err := LoadPoints(path)
	if err != nil {
		return 0, err
	}
	// Figure out the orientation of the cycle.
	orientation, err := Orientation(points, verbose)
	if err != nil {
		return 0, err
	}
	// Create vectors that point to where "outside" is around the perimeter.
	pv := PerimeterVectors(points, orientation, verbose)
	// Compute pairwise distances of points.
	var maxSquareSize int
	for i := 0; i < len(points)-1; i++ {
		p := points[i]
		for j := i + 1; j < len(points); j++ {
			q := points[j]
			if verbose {
				fmt.Printf("Testing %s and %s.\n", p, q)
			}
			size := (Abs(p.X-q.X) + 1) * (Abs(p.Y-q.Y) + 1)
			if size <= maxSquareSize {
				continue
			}
			if !IsBoxInPerimeter(p, q, pv, verbose) {
				continue
			}
			maxSquareSize = size
		}
	}
	return maxSquareSize, nil
}

func ParsePoint2D(enc string) (*Point2D, error) {
	vals := strings.Split(enc, ",")
	if len(vals) != 2 {
		return nil, fmt.Errorf("expected 2 comma-separated ints, got %q", enc)
	}
	var p Point2D
	var err error
	p.X, err = strconv.Atoi(vals[0])
	if err != nil {
		return nil, err
	}
	p.Y, err = strconv.Atoi(vals[1])
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func LoadPoints(path string) ([]*Point2D, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	var points []*Point2D
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		p, err := ParsePoint2D(line)
		if err != nil {
			return nil, err
		}
		points = append(points, p)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return points, nil
}

// Orientation determines if the cycle is oriented positively or negatively.
func Orientation(cycle []*Point2D, verbose bool) (int, error) {
	if len(cycle) < 4 {
		return 0, fmt.Errorf("expecting at least 4 points to form a box, got %d points", len(cycle))
	}
	dWrap := Point2D{
		X: Sign(cycle[0].X - cycle[len(cycle)-1].X),
		Y: Sign(cycle[0].Y - cycle[len(cycle)-1].Y),
	}
	dPrev := dWrap
	pPrev := cycle[0]
	var orientation int
	for i := 1; i < len(cycle); i++ {
		p := cycle[i]
		d := Point2D{
			X: Sign(p.X - pPrev.X),
			Y: Sign(p.Y - pPrev.Y),
		}
		det := Determinant(&dPrev, &d)
		orientation += det
		if verbose {
			fmt.Printf("%s -> %s: %+d %+d\n", &dPrev, &d, det, orientation)
		}
		pPrev = p
		dPrev = d
	}
	det := dPrev.X*dWrap.Y - dWrap.X*dPrev.Y
	orientation += det
	if verbose {
		fmt.Printf("%s -> %s: %+d %+d\n", &dPrev, &dWrap, det, orientation)
	}
	return Sign(orientation), nil
}

func PerimeterVectors(cycle []*Point2D, orientation int, verbose bool) map[Point2D][]*Point2D {
	res := make(map[Point2D][]*Point2D)
	pPrev := cycle[0]
	for i := 1; i < len(cycle)+1; i++ {
		p := cycle[i%len(cycle)]
		d := Point2D{
			X: Sign(p.X - pPrev.X),
			Y: Sign(p.Y - pPrev.Y),
		}
		// Turn against the cycles orientation to point outwards.
		dRot := Point2D{
			X: d.Y * orientation,
			Y: -d.X * orientation,
		}
		q := pPrev
		for {
			res[*q] = append(res[*q], &dRot)
			if verbose {
				fmt.Printf("%s -> %s\n", q, &dRot)
			}
			if *q == *p {
				break
			}
			q = &Point2D{
				X: q.X + d.X,
				Y: q.Y + d.Y,
			}
		}
		pPrev = p
	}
	// We need to remove the perimeter vectors for "concave" corners.
	pPrev = cycle[len(cycle)-1]
	pPrevPrev := cycle[len(cycle)-2]
	dPrev := Point2D{
		X: Sign(pPrev.X - pPrevPrev.X),
		Y: Sign(pPrev.Y - pPrevPrev.Y),
	}
	for i := 0; i < len(cycle); i++ {
		p := cycle[i]
		d := Point2D{
			X: Sign(p.X - pPrev.X),
			Y: Sign(p.Y - pPrev.Y),
		}
		do := Determinant(&dPrev, &d)
		if do != orientation {
			if verbose {
				fmt.Printf("Detected concave corner at %s.\n", pPrev)
			}
			delete(res, *pPrev)
		}
		pPrev = p
		dPrev = d
	}
	return res
}

func IsBoxInPerimeter(p, q *Point2D, pv map[Point2D][]*Point2D, verbose bool) bool {
	corners := []*Point2D{
		p,
		{X: p.X, Y: q.Y},
		q,
		{X: q.X, Y: p.Y},
	}
	src := corners[3]
	for i := 0; i < len(corners); i++ {
		// March along each edge and see if we're leaving the perimiter.
		dst := corners[i]
		d := Point2D{
			X: Sign(dst.X - src.X),
			Y: Sign(dst.Y - src.Y),
		}
		for r := src; *r != *dst; {
			// Check if we're at the perimeter and moving outside.
			vs, ok := pv[*r]
			if ok {
				for _, v := range vs {
					if *v == d {
						if verbose {
							fmt.Printf("Stepping out of bounds at %s in the direction %s.\n", r, &d)
						}
						// TODO: The perimeter might "touch" itself, so we're crossing the perimeter in and out simultaneously.
						return false
					}
				}
			}
			r = &Point2D{X: r.X + d.X, Y: r.Y + d.Y}
		}
		src = dst
	}
	return true
}

func Abs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}

func Sign(i int) int {
	if i < 0 {
		return -1
	}
	if i > 0 {
		return 1
	}
	return 0
}

func Determinant(a, b *Point2D) int {
	return a.X*b.Y - b.X*a.Y
}
