package terrain

import (
	"math"
	"reflect"
	"testing"
)

func TestContourTopologies(t *testing.T) {
	tests := []struct {
		name  string
		field Field
		level float64
		want  int
	}{
		{"empty ocean", Field{3, []float64{0, 0, 0, 0, 0, 0, 0, 0, 0}}, .5, 0},
		{"single island", Field{3, []float64{0, 0, 0, 0, 1, 0, 0, 0, 0}}, .5, 1},
		{"closed boundary", Field{3, []float64{1, 1, 1, 1, 1, 1, 1, 1, 1}}, .5, 1},
		{"island with lake", Field{3, []float64{1, 1, 1, 1, 0, 1, 1, 1, 1}}, .5, 2},
		{"separate islets", Field{3, []float64{1, 0, 0, 0, 0, 0, 0, 0, 1}}, .5, 2},
		{"exact level corners", Field{3, []float64{.5, .5, .5, .5, 1, .5, .5, .5, .5}}, .5, 1},
		{"exact level field", Field{2, []float64{.5, .5, .5, .5}}, .5, 0},
		{"tied saddle", Field{2, []float64{1, 0, 0, 1}}, .5, 2},
		{"connected saddle", Field{2, []float64{1, .4, .4, 1}}, .5, 1},
		{"disconnected saddle", Field{2, []float64{.6, 0, 0, .6}}, .5, 2},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rings := Contour(tc.field, tc.level)
			if len(rings) != tc.want {
				t.Fatalf("got %d rings, want %d", len(rings), tc.want)
			}
			for _, ring := range rings {
				if len(ring) < 3 {
					t.Fatal("not a polygon")
				}
				area := 0.0
				for i, p := range ring {
					next := ring[(i+1)%len(ring)]
					area += p.X*next.Y - next.X*p.Y
					if math.IsNaN(p.X) || math.IsNaN(p.Y) || p.X < 0 || p.Y < 0 || p.X > float64(tc.field.Size-1) || p.Y > float64(tc.field.Size-1) {
						t.Fatal("invalid coordinate")
					}
				}
				if math.Abs(area) < 1e-12 {
					t.Fatal("zero-area polygon")
				}
			}
			if !reflect.DeepEqual(rings, Contour(tc.field, tc.level)) {
				t.Fatal("unstable ring order")
			}
		})
	}
}

func TestInvalidField(t *testing.T) {
	for _, f := range []Field{{}, {1, []float64{1}}, {3, []float64{1}}} {
		if len(Contour(f, .5)) != 0 {
			t.Fatal("invalid field produced a ring")
		}
	}
}

func TestAsymptoticDeciderTies(t *testing.T) {
	if got := cellEdges(5, [4]float64{1, 0, 1, 0}, .5); !reflect.DeepEqual(got, [][2]int{{0, 3}, {1, 2}}) {
		t.Fatal(got)
	}
	if got := cellEdges(10, [4]float64{0, 1, 0, 1}, .5); !reflect.DeepEqual(got, [][2]int{{0, 1}, {2, 3}}) {
		t.Fatal(got)
	}
}
