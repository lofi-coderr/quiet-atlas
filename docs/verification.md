# Verification report

Verified on 2026-10-10 using macOS Apple Silicon. The application was implemented and tested on an external APFS SSD in a dedicated project folder.

## Passed automated checks

| Check | Result |
| --- | --- |
| `go test ./...` | Passed |
| `go test -race -coverprofile=... ./...` | Passed; no data races |
| `go vet ./...` | Passed |
| `npm run typecheck --prefix web` | Passed |
| `npm test --prefix web` | 3 tests passed |
| `npm run build --prefix web` | Passed |
| Go executable with `-trimpath` and embedded frontend | Built successfully |

Coverage: `internal/terrain` **98.8%**, `internal/server` **96.9%**. Startup and the embed wrapper are exercised by the standalone executable checks rather than unit coverage. There are 12 top-level Go test groups, with additional input and topology subtests.

Terrain checks include safe seed boundaries, invalid/nonfinite inputs, equal-seed byte reproducibility, different-seed variation, concurrent RNG isolation, all-zero perimeter samples, finite normalized heights, land remaining at the highest permitted sea level, decreasing land as sea rises, and unchanged terrain when contour strokes are disabled. Contour fixtures cover closed boundaries, lakes, disconnected islets, tied and connected saddles, and exact-level samples. Generated SVGs are parsed as standalone XML.

HTTP checks include default and extreme parameters, matching JSON/SVG output, duplicate and unknown inputs, malformed query escapes, nonfinite values, invalid formats, unsupported methods, missing assets, and static headers.

## Passed browser checks

Tested in the Codex in-app browser on the local machine:

- Initial island rendered as a complete image; English labels and map metadata were visible.
- Seed `0`, seed `42`, and the maximum safe seed rendered successfully.
- An out-of-range seed displayed an accessible error and kept the previous valid map.
- Keyboard-operated sea-level and roughness sliders reached both limits.
- The contour checkbox updated the map.
- Six rapid seed changes settled on the final seed `1006`.
- Reset restored all four defaults; Surprise me produced and displayed a new safe seed.
- The browser saved `quiet-atlas-42.svg` with sea level `0.70`, roughness `0`, and contours disabled. Its **29,348 bytes** matched the Go response byte for byte. SHA-256: `8cb03531f64ceae01bb7a586d99712be71c740add98fbf6906c7bcae7052201b`.
- Stopping only the test application's server produced a readable error and Try again action. Restarting it and retrying rendered the requested seed `73`.
- Widths **320, 390, 768, 1024, and 1440 px** showed no horizontal document overflow. The map loaded at each width. The phone layout places the map above the controls.
- No unexpected browser console errors or warnings appeared in normal operation. The deliberate stopped-server check produced the expected failed local request.

Screenshots were saved separately as desktop (1440 px viewport) and mobile (390 px viewport) deliverables; screenshot files are excluded from this repository. These are browser viewport tests on the local Mac, not physical phone tests.

## Standalone executable proof

The macOS ARM64 executable was copied into an isolated directory with no adjacent frontend assets and launched with a PATH containing only system command directories. Node.js was unavailable in that PATH. It served the embedded index and API response identically to the project build on a separate local port. This verifies runtime independence from Node and `web/dist` on disk.

The source ZIP was also extracted into a fresh directory on the external SSD and rebuilt there. The restored source produced a byte-identical executable. Both ZIP archives passed CRC checks, and the source archive contains the compiled frontend while excluding caches and Git metadata.

## Scope and remaining limits

- macOS Apple Silicon is verified. Windows, Linux, Intel macOS, physical mobile devices, and other browser engines have not been tested.
- The binary is not signed or notarized.
- The authorized Library source handoff was successfully materialized and read on this machine. No new credentials or expanded authentication were created.
- This verification report covers implementation and local validation. Source publication is tracked separately in the release policy; no external application deployment was performed.
- Existing 2D services, data, and other projects were not changed.
- Source and deliverable packaging exclude dependency caches, Git metadata, machine paths, and private identifiers.
