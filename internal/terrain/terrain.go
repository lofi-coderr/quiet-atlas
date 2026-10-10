// Package terrain generates deterministic height fields and vector islands.
package terrain

import (
	"errors"
	"math"
	"math/rand"
)

const (
	MaxSeed    int64 = 9007199254740991 // JavaScript's largest exact integer.
	Resolution       = 161
	CanvasSize       = 800
)

// Parameters are shared by the JSON API, SVG export, and generator.
type Parameters struct {
	Seed      int64   `json:"seed"`
	SeaLevel  float64 `json:"seaLevel"`
	Roughness float64 `json:"roughness"`
	Contours  bool    `json:"contours"`
}

func (p Parameters) Validate() error {
	if p.Seed < 0 || p.Seed > MaxSeed {
		return errors.New("seed must be an integer between 0 and 9007199254740991")
	}
	if math.IsNaN(p.SeaLevel) || math.IsInf(p.SeaLevel, 0) || p.SeaLevel < 0.25 || p.SeaLevel > 0.7 {
		return errors.New("seaLevel must be between 0.25 and 0.70")
	}
	if math.IsNaN(p.Roughness) || math.IsInf(p.Roughness, 0) || p.Roughness < 0 || p.Roughness > 1 {
		return errors.New("roughness must be between 0 and 1")
	}
	return nil
}

// Field stores row-major heights in [0, 1]. Its boundary is always underwater.
type Field struct {
	Size    int
	Heights []float64
}

func (f Field) At(x, y int) float64 { return f.Heights[y*f.Size+x] }

// Generate owns its RNG: concurrent requests never share mutable random state.
func Generate(p Parameters) Field {
	rng := rand.New(rand.NewSource(p.Seed))
	noise := makeNoise(rng)
	phase := rng.Float64() * 2 * math.Pi
	stretch := 0.85 + rng.Float64()*0.3
	f := Field{Size: Resolution, Heights: make([]float64, Resolution*Resolution)}
	for y := 0; y < f.Size; y++ {
		for x := 0; x < f.Size; x++ {
			u := float64(x)/float64(f.Size-1)*2 - 1
			v := float64(y)/float64(f.Size-1)*2 - 1
			angle := math.Atan2(v, u)
			radius := math.Hypot(u*stretch, v/stretch)
			coast := 1 + 0.12*math.Sin(3*angle+phase) + 0.06*math.Cos(5*angle-phase)
			// Broad hills keep the island coherent; fine octaves add coves and ridges.
			n := 0.56*noise.at(u*2.4+13, v*2.4+17) + 0.27*noise.at(u*5.5+31, v*5.5+9)
			n += 0.12*noise.at(u*11+7, v*11+41) + 0.05*noise.at(u*23+19, v*23+23)
			h := 1.02 - math.Pow(radius/coast, 1.7) + n*(0.24+p.Roughness*0.62)
			// Explicit ocean margin guarantees closed coastlines at every valid setting.
			margin := math.Min(1, math.Max(0, (0.95-radius)/0.16))
			margin = math.Min(margin, math.Min(1, (1-math.Max(math.Abs(u), math.Abs(v)))/0.12))
			f.Heights[y*f.Size+x] = math.Max(0, math.Min(1, h*margin))
		}
	}
	return f
}

type valueNoise [64 * 64]float64

func makeNoise(rng *rand.Rand) valueNoise {
	var n valueNoise
	for i := range n {
		n[i] = rng.Float64()*2 - 1
	}
	return n
}

func (n valueNoise) at(x, y float64) float64 {
	x0, y0 := int(math.Floor(x)), int(math.Floor(y))
	xf, yf := x-float64(x0), y-float64(y0)
	// Quintic interpolation makes octave boundaries smooth.
	smooth := func(t float64) float64 { return t * t * t * (t*(t*6-15) + 10) }
	sample := func(x, y int) float64 { return n[((y%64+64)%64)*64+(x%64+64)%64] }
	lerp := func(a, b, t float64) float64 { return a + (b-a)*t }
	a := lerp(sample(x0, y0), sample(x0+1, y0), smooth(xf))
	b := lerp(sample(x0, y0+1), sample(x0+1, y0+1), smooth(xf))
	return lerp(a, b, smooth(yf))
}
