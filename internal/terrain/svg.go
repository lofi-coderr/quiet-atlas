package terrain

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

type Island struct {
	Parameters Parameters `json:"parameters"`
	SVG        string     `json:"svg"`
	Land       int        `json:"landPercent"`
	Name       string     `json:"name"`
}

// Build creates a complete SVG in Go. The browser displays and saves these
// same bytes, so downloads always match the island currently on screen.
func Build(p Parameters) (Island, error) {
	if err := p.Validate(); err != nil {
		return Island{}, err
	}
	f := Generate(p)
	land := 0
	for _, h := range f.Heights {
		if h > p.SeaLevel {
			land++
		}
	}
	first := []string{"Peach", "Clover", "Mallow", "Pebble", "Honey", "Mint", "Sunny", "Coral"}
	last := []string{"Cove", "Isle", "Bay", "Haven", "Lagoon", "Key", "Shore", "Atoll"}
	name := first[p.Seed%8] + " " + last[(p.Seed/8)%8]
	var svg strings.Builder
	fmt.Fprintf(&svg, `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 800 800" width="800" height="800" role="img" aria-labelledby="title desc"><title id="title">%s</title><desc id="desc">Quiet Atlas island. Seed %d, sea level %.2f, roughness %.2f.</desc>`, name, p.Seed, p.SeaLevel, p.Roughness)
	svg.WriteString(`<defs><pattern id="waves" width="70" height="56" patternUnits="userSpaceOnUse"><path d="M10 28q5-5 10 0t10 0" fill="none" stroke="#70bbbd" stroke-width="1.4" opacity=".32"/></pattern></defs><rect width="800" height="800" fill="#b4e3df"/><rect width="800" height="800" fill="url(#waves)"/>`)
	shore := compoundPath(Contour(f, p.SeaLevel), f.Size)
	// A wide outer stroke gives the coast a soft shallow-water halo.
	fmt.Fprintf(&svg, `<path d="%s" fill="#92d2cf" stroke="#92d2cf" stroke-width="19" stroke-linejoin="round" fill-rule="evenodd"/>`, shore)
	colors := []string{"#ffe2a8", "#d3e6a3", "#a9d59c", "#80bd99", "#67a893", "#578f85", "#c2cfba", "#eff0d7"}
	for i, color := range colors {
		level := p.SeaLevel + float64(i)*0.082
		if level >= 1 {
			break
		}
		path := shore
		if i != 0 {
			path = compoundPath(Contour(f, level), f.Size)
		}
		fmt.Fprintf(&svg, `<path d="%s" fill="%s" fill-rule="evenodd"/>`, path, color)
		if p.Contours && i > 0 {
			fmt.Fprintf(&svg, `<path d="%s" fill="none" stroke="#456f68" stroke-width="1" opacity=".22" stroke-linejoin="round"/>`, path)
		}
	}
	fmt.Fprintf(&svg, `<path d="%s" fill="none" stroke="#ddbb83" stroke-width="1.5" stroke-linejoin="round"/>`, shore)
	// A small illustrated compass lives in the export too.
	svg.WriteString(`<g transform="translate(714 82)" fill="#4e8180"><path d="M0-20 7 0 0-5-7 0Z"/><path d="M0 20 7 0 0 5-7 0Z" opacity=".35"/><text y="-29" text-anchor="middle" font-family="sans-serif" font-size="12" font-weight="700">N</text></g></svg>`)
	return Island{p, svg.String(), int(math.Round(float64(land) * 100 / float64(len(f.Heights)))), name}, nil
}

func compoundPath(rings []Ring, size int) string {
	var b strings.Builder
	appendNumber := func(v float64) { b.WriteString(strconv.FormatFloat(v*680/float64(size-1)+60, 'f', 2, 64)) }
	for _, ring := range rings {
		for i, p := range ring {
			if i == 0 {
				b.WriteByte('M')
			} else {
				b.WriteByte('L')
			}
			appendNumber(p.X)
			b.WriteByte(' ')
			appendNumber(p.Y)
		}
		b.WriteByte('Z')
	}
	return b.String()
}
