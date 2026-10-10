# Quiet Atlas

Tiny worlds, happy discoveries. Quiet Atlas is a local, account-free island generator with a playful pastel interface. Adjust a seed, sea level, terrain roughness, and contour lines, then keep your island as a standalone SVG.

## Build and run

Requirements: Go 1.24 or newer, Node.js 22.18 or newer, npm, and make. The verified toolchain was Go 1.26.0, Node.js 26.8.2, and npm 11.19.1. Dependency installation uses the public npm registry; the Go backend uses only the standard library.

```sh
make install
make test
make build
make run
```

Open **http://127.0.0.1:8080**. Keep the terminal open while exploring; press **Ctrl+C** to stop. If the port is occupied:

```sh
./bin/quiet-atlas -addr 127.0.0.1:8090
```

`make build` type-checks and builds the React application into `web/dist`, then embeds it in `bin/quiet-atlas` using Go's `embed` package. It builds with `-trimpath`. Build and package caches stay inside the project. Node.js is needed to develop or rebuild the frontend. The resulting executable includes the frontend and needs no Node.js, asset folder, account, or external API at runtime. It listens on loopback by default.

This repository contains source code, tests, configuration, and the npm lockfile. Build output, release binaries, screenshots, and generated example maps are excluded. The executable is not signed or notarized.

## Development

After `make install` and `make build`, run the Go server in one terminal:

```sh
make run
```

In a second terminal:

```sh
npm run dev --prefix web
```

Open Vite's local address, usually http://127.0.0.1:5173. Its `/api` proxy targets the Go server on port 8080. Restart or rebuild the Go server after changing Go code. The production executable serves both the UI and API from the same origin.

## Make a tiny world

- **Seed:** a whole number from `0` through `9007199254740991`. The same seed and settings reproduce the same island.
- **Sea level:** `25%` to `70%`. Higher values leave less land.
- **Terrain roughness:** `0%` to `100%`. Higher values add coves and ridges.
- **Contour lines:** show or hide elevation strokes without changing the terrain.
- **Surprise me** or the dice button: choose a new seed using the browser's cryptographic random generator.
- **Reset settings:** return to seed `2408`, sea level `45%`, roughness `55%`, and visible contours.
- **Save SVG:** download the exact SVG displayed in the preview. The button waits for the latest valid map to finish. Invalid seed input keeps the last valid map visible.

Settings update the map automatically after a brief debounce. Old requests are canceled and ignored, so rapid changes cannot replace the newest map with an obsolete response. Maps are generated in memory; the app has no database, accounts, analytics, or remote runtime dependencies. Keep a seed and its settings to recreate a favorite island later.

## API

`GET /api/island` returns a JSON object containing `parameters`, `svg`, `landPercent`, and a deterministic island `name`.

| Query parameter | Default | Valid values |
| --- | --- | --- |
| `seed` | `2408` | Integer, `0`–`9007199254740991` |
| `seaLevel` | `0.45` | Number, `0.25`–`0.70` |
| `roughness` | `0.55` | Number, `0`–`1` |
| `contours` | `true` | `true` or `false` |
| `format` | `json` | `json` or `svg` |

```sh
curl 'http://127.0.0.1:8080/api/island?seed=42&seaLevel=0.45&roughness=0.55&contours=true'
curl 'http://127.0.0.1:8080/api/island?seed=42&format=svg' -o island.svg
```

The SVG response includes an attachment filename. Invalid, nonfinite, duplicate, malformed, or unknown parameters receive `400` with a JSON error. Other API methods receive `405`; missing routes receive `404`. API responses use `Cache-Control: no-store`.

## Architecture

```text
cmd/quiet-atlas/       Local HTTP startup, timeouts, shutdown
internal/terrain/     Parameter validation, seeded field, coastlines, contours, SVG
internal/server/      HTTP parsing, JSON/SVG endpoints, static routing
web/src/lib/          Seed validation, request serialization, SVG downloads
web/src/hooks/        Debounced fetches and stale-response protection
web/src/components/   Controls, map preview, lightweight vector icons
web/dist/             Compiled frontend, embedded in the executable
docs/                 Verification report and release policy
```

Go owns all terrain, shoreline, contour, SVG, name, and land-percentage calculations. React and TypeScript own only inputs, request state, display, and download actions. The preview uses a Blob URL for the server's SVG rather than inserting markup into the page.

The terrain is a `161 × 161` height field built from request-local seeded value noise and a shaped island falloff. An explicit ocean margin keeps the boundary underwater. Marching squares uses canonical shared edges, an asymptotic decider for saddle cells, and closed compound rings with an even-odd fill rule to preserve lakes and separate islets. Display and downloads use the same 800 × 800 vector document.

## Verification

See [the verification report](docs/verification.md) for passed checks and limits. Tests cover reproducibility, independent concurrent requests, parameter boundaries, coastline topology, saddle ties, exact-level corners, valid SVG XML, HTTP errors, seed parsing, and random-seed bounds. Browser checks cover controls, rapid updates, download byte equality, errors, retries, and responsive widths from 320 to 1440 pixels.

## Release policy

This repository remains private during preparation. Public visibility is planned for the associated video publication. That visibility change is a separate release action; it is not part of the source push. See [the release policy](docs/release-policy.md).

## References

- [React effect cleanup](https://react.dev/learn/synchronizing-with-effects)
- [Vite documentation](https://vite.dev/guide/)
- [Go embed documentation](https://pkg.go.dev/embed)
