package day08

import (
	"bufio"
	"cmp"
	"fmt"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"
)

func Round1(path string, numPairs int, numClusters int, verbose bool) (int64, error) {
	points, err := LoadPoints(path)
	if err != nil {
		return 0, err
	}
	// Compute pairwise distances of points.
	var pairs []*pointPair
	for i := 0; i < len(points)-1; i++ {
		for j := i + 1; j < len(points); j++ {
			pair := &pointPair{
				p1:   points[i],
				p2:   points[j],
				i:    i,
				j:    j,
				dist: points[i].SLDist(points[j]),
			}
			pairs = append(pairs, pair)
		}
	}
	slices.SortFunc(pairs, func(a, b *pointPair) int {
		return cmp.Compare(a.dist, b.dist)
	})
	// Create clusters for the closest pairs.
	type cluster struct {
		indexes []int
	}
	clusters := make([]*cluster, len(points))
	for i, pair := range pairs {
		if i >= numPairs {
			break
		}
		if verbose {
			fmt.Printf("dist(%s, %s) = %.1f\n", pair.p1, pair.p2, pair.dist)
		}
		c1 := clusters[pair.i]
		c2 := clusters[pair.j]
		if c1 == nil && c2 == nil {
			c := cluster{[]int{pair.i, pair.j}}
			clusters[pair.i] = &c
			clusters[pair.j] = &c
			continue
		}
		if c2 == nil {
			c1.indexes = append(c1.indexes, pair.j)
			clusters[pair.j] = c1
			continue
		}
		if c1 == nil {
			c2.indexes = append(c2.indexes, pair.i)
			clusters[pair.i] = c2
			continue
		}
		if c1 == c2 {
			continue
		}
		c1.indexes = append(c1.indexes, c2.indexes...)
		for _, k := range c2.indexes {
			clusters[k] = c1
		}
	}
	// Find the distinct clusters.
	var sortedClusters []*cluster
	distinctClusters := make(map[int]struct{})
	for _, c := range clusters {
		if c == nil {
			continue
		}
		_, ok := distinctClusters[c.indexes[0]]
		if ok {
			continue
		}
		distinctClusters[c.indexes[0]] = struct{}{}
		sortedClusters = append(sortedClusters, c)
	}
	slices.SortFunc(sortedClusters, func(a, b *cluster) int {
		return len(b.indexes) - len(a.indexes)
	})
	var product int64 = 1
	for i, c := range sortedClusters {
		if i >= numClusters {
			break
		}
		product *= int64(len(c.indexes))
	}
	return product, nil
}

func Round2(path string, verbose bool) (int, error) {
	points, err := LoadPoints(path)
	if err != nil {
		return 0, err
	}
	// Compute pairwise distances of points.
	var pairs []*pointPair
	for i := 0; i < len(points)-1; i++ {
		for j := i + 1; j < len(points); j++ {
			pair := &pointPair{
				p1:   points[i],
				p2:   points[j],
				i:    i,
				j:    j,
				dist: points[i].SLDist(points[j]),
			}
			pairs = append(pairs, pair)
		}
	}
	slices.SortFunc(pairs, func(a, b *pointPair) int {
		return cmp.Compare(a.dist, b.dist)
	})
	// Create clusters for the closest pairs.
	type cluster struct {
		indexes []int
	}
	clusters := make([]*cluster, len(points))
	for _, pair := range pairs {
		if verbose {
			fmt.Printf("dist(%s, %s) = %.1f\n", pair.p1, pair.p2, pair.dist)
		}
		var c *cluster
		c1 := clusters[pair.i]
		c2 := clusters[pair.j]
		if c1 == nil && c2 == nil {
			c = &cluster{[]int{pair.i, pair.j}}
			clusters[pair.i] = c
			clusters[pair.j] = c
		} else if c2 == nil {
			c1.indexes = append(c1.indexes, pair.j)
			clusters[pair.j] = c1
			c = c1
		} else if c1 == nil {
			c2.indexes = append(c2.indexes, pair.i)
			clusters[pair.i] = c2
			c = c2
		} else if c1 == c2 {
			continue
		} else {
			c1.indexes = append(c1.indexes, c2.indexes...)
			for _, k := range c2.indexes {
				clusters[k] = c1
			}
			c = c1
		}
		if len(c.indexes) == len(points) {
			return pair.p1.X * pair.p2.X, nil
		}
	}
	return 0, fmt.Errorf("did not find a single circuit")
}

type Point3D struct {
	X, Y, Z int
}

func ParsePoint3D(enc string) (*Point3D, error) {
	vals := strings.Split(enc, ",")
	if len(vals) != 3 {
		return nil, fmt.Errorf("expected 3 comma-separated ints, got %q", enc)
	}
	var p Point3D
	var err error
	p.X, err = strconv.Atoi(vals[0])
	if err != nil {
		return nil, err
	}
	p.Y, err = strconv.Atoi(vals[1])
	if err != nil {
		return nil, err
	}
	p.Z, err = strconv.Atoi(vals[2])
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (p *Point3D) String() string {
	return fmt.Sprintf("%d,%d,%d", p.X, p.Y, p.Z)
}

func (p *Point3D) SLDist(o *Point3D) float64 {
	dx := p.X - o.X
	dy := p.Y - o.Y
	dz := p.Z - o.Z
	return math.Pow(float64(dx*dx+dy*dy+dz*dz), 0.5)
}

func LoadPoints(path string) ([]*Point3D, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	var points []*Point3D
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		p, err := ParsePoint3D(line)
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

type pointPair struct {
	p1, p2 *Point3D
	i, j   int
	dist   float64
}
