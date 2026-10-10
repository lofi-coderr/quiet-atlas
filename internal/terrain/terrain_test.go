package terrain

import (
	"bytes"
	"encoding/xml"
	"io"
	"math"
	"reflect"
	"strings"
	"sync"
	"testing"
)

func defaults() Parameters { return Parameters{2408, .45, .55, true} }

func TestValidation(t *testing.T) {
	p := defaults()
	cases := []Parameters{{-1, .45, .55, true}, {MaxSeed + 1, .45, .55, true}, {0, .24, .55, true}, {0, .71, .55, true}, {0, math.NaN(), .55, true}, {0, math.Inf(1), .55, true}, {0, .45, -.01, true}, {0, .45, 1.01, true}, {0, .45, math.NaN(), true}, {0, .45, math.Inf(-1), true}}
	for _, bad := range cases {
		if bad.Validate() == nil {
			t.Errorf("accepted invalid parameters: %+v", bad)
		}
		if _, err := Build(bad); err == nil {
			t.Errorf("Build accepted invalid parameters: %+v", bad)
		}
	}
	for _, seed := range []int64{0, MaxSeed} {
		p.Seed = seed
		if err := p.Validate(); err != nil {
			t.Fatal(err)
		}
	}
}

func TestDeterminismAndIndependentRequests(t *testing.T) {
	p := defaults()
	expected, err := Build(p)
	if err != nil {
		t.Fatal(err)
	}
	if actual, _ := Build(p); !reflect.DeepEqual(actual, expected) {
		t.Fatal("same seed changed output")
	}
	var wg sync.WaitGroup
	for i := 0; i < 12; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			other := p
			other.Seed += int64(i + 1)
			_, _ = Build(other)
			actual, _ := Build(p)
			if actual.SVG != expected.SVG {
				t.Error("concurrent requests changed seeded output")
			}
		}(i)
	}
	wg.Wait()
	other := p
	other.Seed++
	actual, _ := Build(other)
	if actual.SVG == expected.SVG {
		t.Fatal("different seeds produced the same SVG")
	}
}

func TestHeightBoundsAndOceanBoundary(t *testing.T) {
	for seed := int64(0); seed < 40; seed++ {
		for _, roughness := range []float64{0, 1} {
			f := Generate(Parameters{seed, .25, roughness, true})
			land := 0
			for y := 0; y < f.Size; y++ {
				for x := 0; x < f.Size; x++ {
					h := f.At(x, y)
					if math.IsNaN(h) || math.IsInf(h, 0) || h < 0 || h > 1 {
						t.Fatalf("invalid height: %v", h)
					}
					if (x == 0 || y == 0 || x == f.Size-1 || y == f.Size-1) && h != 0 {
						t.Fatalf("seed %d: boundary is not ocean", seed)
					}
					if h > .7 {
						land++
					}
				}
			}
			if land == 0 {
				t.Fatalf("seed %d disappeared at the highest sea level", seed)
			}
		}
	}
}

func TestSeaLevelAndContourSettings(t *testing.T) {
	p := defaults()
	p.SeaLevel = .25
	low, _ := Build(p)
	p.SeaLevel = .7
	high, _ := Build(p)
	if low.Land <= high.Land {
		t.Fatal("higher sea did not reduce land")
	}
	p = defaults()
	on, _ := Build(p)
	p.Contours = false
	off, _ := Build(p)
	if on.SVG == off.SVG || on.Land != off.Land || on.Name != off.Name {
		t.Fatal("contour setting changed terrain or failed to change SVG")
	}
	if !strings.Contains(on.SVG, `opacity=".22"`) || strings.Contains(off.SVG, `opacity=".22"`) {
		t.Fatal("contour strokes do not match setting")
	}
}

func TestSVGIsStandaloneXML(t *testing.T) {
	for _, seed := range []int64{0, 2408, MaxSeed} {
		p := defaults()
		p.Seed = seed
		island, err := Build(p)
		if err != nil {
			t.Fatal(err)
		}
		decoder := xml.NewDecoder(bytes.NewBufferString(island.SVG))
		roots := 0
		paths := 0
		for {
			token, err := decoder.Token()
			if err == io.EOF {
				break
			}
			if err != nil {
				t.Fatal(err)
			}
			if element, ok := token.(xml.StartElement); ok {
				if element.Name.Local == "svg" {
					roots++
					if element.Name.Space != "http://www.w3.org/2000/svg" {
						t.Fatal("missing SVG namespace")
					}
				}
				if element.Name.Local == "path" {
					paths++
				}
			}
		}
		if roots != 1 || paths < 4 || strings.Contains(island.SVG, "NaN") || strings.Contains(island.SVG, "Infinity") {
			t.Fatal("invalid standalone SVG")
		}
	}
}

func BenchmarkBuild(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = Build(defaults())
	}
}
