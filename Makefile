export GOCACHE := $(CURDIR)/.cache/go-build
export GOMODCACHE := $(CURDIR)/.cache/go-mod

.PHONY: install build test run clean
install:
	npm ci --prefix web --cache .cache/npm --ignore-scripts --no-audit --no-fund

build:
	npm run build --prefix web
	mkdir -p bin
	go build -trimpath -o bin/quiet-atlas ./cmd/quiet-atlas

test:
	go test -race ./...
	go vet ./...
	npm run typecheck --prefix web
	npm test --prefix web

run:
	./bin/quiet-atlas

clean:
	rm -rf bin web/dist
	mkdir -p web/dist
	touch web/dist/.gitkeep
