package terrain

import "math"

type Point struct{ X, Y float64 }
type Ring []Point

type vertex struct {
	key   int
	point Point
}
type segment struct{ a, b vertex }

// Contour traces a superlevel set with marching squares. One below-level sample
// ring closes shapes that touch the field boundary. Rings are suitable for an
// SVG evenodd compound path, preserving lakes and disconnected islands.
func Contour(f Field, level float64) []Ring {
	if f.Size < 2 || len(f.Heights) != f.Size*f.Size {
		return nil
	}
	width := f.Size + 2
	value := func(x, y int) float64 {
		if x == 0 || y == 0 || x == width-1 || y == width-1 {
			return level - 1
		}
		return f.At(x-1, y-1)
	}
	var segments []segment
	for y := 0; y < width-1; y++ {
		for x := 0; x < width-1; x++ {
			v := [4]float64{value(x, y), value(x+1, y), value(x+1, y+1), value(x, y+1)}
			mask := 0
			for i, h := range v {
				if h > level {
					mask |= 1 << i
				}
			}
			edge := func(e int) vertex {
				// Canonical grid edge IDs join neighboring cells exactly, without
				// rounded coordinate lookups or dependence on map iteration order.
				ends := [4][2]int{{0, 1}, {1, 2}, {3, 2}, {0, 3}}
				coords := [4]Point{{float64(x - 1), float64(y - 1)}, {float64(x), float64(y - 1)}, {float64(x), float64(y)}, {float64(x - 1), float64(y)}}
				i, j := ends[e][0], ends[e][1]
				t := (level - v[i]) / (v[j] - v[i])
				// Keep crossings on distinct edges when a sample equals the level.
				t = math.Max(1e-9, math.Min(1-1e-9, t))
				p := Point{coords[i].X + (coords[j].X-coords[i].X)*t, coords[i].Y + (coords[j].Y-coords[i].Y)*t}
				p.X = math.Max(0, math.Min(float64(f.Size-1), p.X))
				p.Y = math.Max(0, math.Min(float64(f.Size-1), p.Y))
				key := 2 * (y*width + x)
				switch e {
				case 1:
					key = 2*(y*width+x+1) + 1
				case 2:
					key = 2 * ((y+1)*width + x)
				case 3:
					key++
				}
				return vertex{key, p}
			}
			for _, pair := range cellEdges(mask, v, level) {
				segments = append(segments, segment{edge(pair[0]), edge(pair[1])})
			}
		}
	}
	adjacent := make(map[int][]int, len(segments))
	for i, s := range segments {
		adjacent[s.a.key] = append(adjacent[s.a.key], i)
		adjacent[s.b.key] = append(adjacent[s.b.key], i)
	}
	used := make([]bool, len(segments))
	var rings []Ring
	for start, s := range segments {
		if used[start] {
			continue
		}
		ring := Ring{s.a.point}
		current, index := s.a, start
		closed := false
		for {
			used[index] = true
			next := segments[index].b
			if current.key == next.key {
				next = segments[index].a
			}
			if next.key == s.a.key {
				closed = true
				break
			}
			ring = append(ring, next.point)
			neighbors := adjacent[next.key]
			index = -1
			for _, candidate := range neighbors {
				if !used[candidate] {
					index = candidate
					break
				}
			}
			if index == -1 {
				break
			}
			current = next
		}
		if closed && len(ring) >= 3 {
			rings = append(rings, ring)
		}
	}
	return rings
}

// cellEdges uses the bilinear asymptotic decider for saddle cases 5 and 10.
// At an exact tie, the two above-level corners remain separate.
func cellEdges(mask int, v [4]float64, level float64) [][2]int {
	switch mask {
	case 0, 15:
		return nil
	case 1, 14:
		return [][2]int{{3, 0}}
	case 2, 13:
		return [][2]int{{0, 1}}
	case 3, 12:
		return [][2]int{{3, 1}}
	case 4, 11:
		return [][2]int{{1, 2}}
	case 6, 9:
		return [][2]int{{0, 2}}
	case 7, 8:
		return [][2]int{{3, 2}}
	case 5, 10:
		q := (v[0]-level)*(v[2]-level) - (v[1]-level)*(v[3]-level)
		if q > 0 || (q == 0 && mask == 10) {
			return [][2]int{{0, 1}, {2, 3}}
		}
		return [][2]int{{0, 3}, {1, 2}}
	}
	return nil
}
