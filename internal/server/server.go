// Package server contains the local HTTP API and embedded asset routing.
package server

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"net/url"
	"strconv"

	"quietatlas/internal/terrain"
)

func New(assets fs.FS) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/island", island)
	mux.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) { http.NotFound(w, r) })
	static := http.FileServer(http.FS(assets))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			w.Header().Set("Allow", "GET, HEAD")
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		path := r.URL.Path
		if path == "/" {
			path = "/index.html"
		}
		info, err := fs.Stat(assets, path[1:])
		if err != nil || info.IsDir() {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Content-Security-Policy", "default-src 'self'; img-src 'self' blob: data:; style-src 'self' 'unsafe-inline'; object-src 'none'; base-uri 'none'; frame-ancestors 'none'")
		static.ServeHTTP(w, r)
	})
	return mux
}

func island(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		writeError(w, http.StatusMethodNotAllowed, "use GET to create an island")
		return
	}
	p, err := parseParameters(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	result, err := terrain.Build(p)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if r.URL.Query().Get("format") == "svg" {
		w.Header().Set("Content-Type", "image/svg+xml; charset=utf-8")
		w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="quiet-atlas-%d.svg"`, p.Seed))
		_, _ = w.Write([]byte(result.SVG))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(result)
}

func parseParameters(r *http.Request) (terrain.Parameters, error) {
	p := terrain.Parameters{Seed: 2408, SeaLevel: 0.45, Roughness: 0.55, Contours: true}
	values, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		return p, err
	}
	for key, v := range values {
		if len(v) != 1 {
			return p, fmt.Errorf("%s must occur only once", key)
		}
		switch key {
		case "seed":
			p.Seed, err = strconv.ParseInt(v[0], 10, 64)
		case "seaLevel":
			p.SeaLevel, err = strconv.ParseFloat(v[0], 64)
		case "roughness":
			p.Roughness, err = strconv.ParseFloat(v[0], 64)
		case "contours":
			if v[0] != "true" && v[0] != "false" {
				err = fmt.Errorf("expected true or false")
			} else {
				p.Contours = v[0] == "true"
			}
		case "format":
			if v[0] != "svg" && v[0] != "json" {
				err = fmt.Errorf("expected svg or json")
			}
		default:
			return p, fmt.Errorf("unknown parameter: %s", key)
		}
		if err != nil {
			return p, fmt.Errorf("invalid %s", key)
		}
	}
	return p, p.Validate()
}

func writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": message})
}
