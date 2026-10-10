package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/fstest"

	"quietatlas/internal/terrain"
)

func handler() http.Handler {
	return New(fstest.MapFS{"index.html": {Data: []byte("<h1>Quiet Atlas</h1>")}, "assets/app.js": {Data: []byte("const island = true;")}})
}
func call(method, path string) *httptest.ResponseRecorder {
	r := httptest.NewRecorder()
	handler().ServeHTTP(r, httptest.NewRequest(method, path, nil))
	return r
}

func TestDefaultAPIAndSVG(t *testing.T) {
	r := call("GET", "/api/island")
	if r.Code != 200 {
		t.Fatal(r.Code, r.Body.String())
	}
	var result terrain.Island
	if err := json.Unmarshal(r.Body.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if result.Parameters.Seed != 2408 || result.SVG == "" || result.Land <= 0 {
		t.Fatal("invalid response")
	}
	svg := call("GET", "/api/island?format=svg")
	if svg.Code != 200 || svg.Body.String() != result.SVG {
		t.Fatal("SVG export differs from JSON")
	}
	if !strings.HasPrefix(svg.Header().Get("Content-Type"), "image/svg+xml") || svg.Header().Get("Content-Disposition") != `attachment; filename="quiet-atlas-2408.svg"` {
		t.Fatal(svg.Header())
	}
	if r.Header().Get("Cache-Control") != "no-store" {
		t.Fatal("missing no-store")
	}
}

func TestInvalidInputs(t *testing.T) {
	for _, query := range []string{"seed=-1", "seed=9007199254740992", "seed=1.5", "seed=cat", "seed=", "seed=1&seed=2", "seed=%xx", "seaLevel=NaN", "seaLevel=Inf", "seaLevel=.24", "seaLevel=.71", "roughness=-.1", "roughness=1.1", "roughness=no", "contours=1", "format=png", "unknown=1", "seed=3;foo=4"} {
		t.Run(query, func(t *testing.T) {
			r := call("GET", "/api/island?"+query)
			if r.Code != 400 {
				t.Fatalf("got %d: %s", r.Code, r.Body.String())
			}
			var message map[string]string
			if err := json.Unmarshal(r.Body.Bytes(), &message); err != nil || message["error"] == "" {
				t.Fatal("missing JSON error")
			}
		})
	}
}

func TestSeedAndParameterBoundaries(t *testing.T) {
	for _, query := range []string{"seed=0&seaLevel=.25&roughness=0&contours=false", "seed=9007199254740991&seaLevel=.7&roughness=1&contours=true"} {
		if r := call("GET", "/api/island?"+query); r.Code != 200 {
			t.Fatal(r.Code, r.Body.String())
		}
	}
}

func TestHTTPMethodsAndMissingPaths(t *testing.T) {
	for _, path := range []string{"/api/missing", "/missing", "/assets/missing.js", "/assets/"} {
		if r := call("GET", path); r.Code != 404 {
			t.Fatal(path, r.Code)
		}
	}
	for _, path := range []string{"/", "/assets/app.js"} {
		if r := call("GET", path); r.Code != 200 {
			t.Fatal(path, r.Code)
		}
		r := call("HEAD", path)
		if r.Code != 200 || r.Body.Len() != 0 {
			t.Fatal(path, r.Code, r.Body.String())
		}
		if r := call("POST", path); r.Code != 405 || r.Header().Get("Allow") != "GET, HEAD" {
			t.Fatal(path, r.Code)
		}
	}
	for _, method := range []string{"POST", "PUT", "DELETE", "HEAD"} {
		r := call(method, "/api/island")
		if r.Code != 405 || r.Header().Get("Allow") != "GET" {
			t.Fatal(method, r.Code)
		}
	}
	r := call("GET", "/")
	if r.Header().Get("X-Content-Type-Options") != "nosniff" || !strings.Contains(r.Header().Get("Content-Security-Policy"), "default-src 'self'") {
		t.Fatal("missing static headers")
	}
}
