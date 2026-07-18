# Makefile — démonstrateur Borealis (interopérabilité agentique MCP + A2A).
# Cibles alignées sur le PRD §14 et le gate du PRDPlan §2.
MODULE  = github.com/agbruneau/borealis
GO      = go
BUILD   = ./build
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
COMMIT  ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo unknown)
DATE    ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS = -ldflags "-s -w \
	-X $(MODULE)/internal/meta.Version=$(VERSION) \
	-X $(MODULE)/internal/meta.Commit=$(COMMIT) \
	-X $(MODULE)/internal/meta.BuildDate=$(DATE)"

.PHONY: all build test race coverage bench vuln lint e2e audit-export docker-build docker-verify check tidy clean help

all: check

## build: compile tous les binaires cmd/ dans ./build
build:
	@mkdir -p $(BUILD)
	$(GO) build $(LDFLAGS) -trimpath -o $(BUILD)/ ./cmd/...

## test: tests unitaires + intégration
test:
	$(GO) test ./...

## race: tests avec détecteur de course (requiert CGO/gcc)
race:
	$(GO) test -race ./...

## coverage: profil + total
coverage:
	$(GO) test -coverprofile=coverage.out ./...
	$(GO) tool cover -func=coverage.out | tail -1

## bench: benchmarks
bench:
	$(GO) test -bench=. -benchmem -run='^$$' ./...

## vuln: analyse de vulnérabilités
vuln:
	govulncheck ./...

## lint: golangci-lint (config .golangci.yml)
lint:
	golangci-lint run ./...

## e2e: scénarios de bout en bout (tag e2e)
e2e:
	$(GO) test -tags e2e ./test/e2e/...

## audit-export: export JSONL rejouable du journal WORM (FR-23, M5.2b)
audit-export:
	$(GO) run ./cmd/orchestrator --audit-export

## docker-build: image multi-stage (requiert Docker — voir docs/REPRODUCIBILITY.md, M4.4)
docker-build:
	@command -v docker >/dev/null 2>&1 || { echo "Docker absent sur cet hôte."; exit 1; }
	docker build -f deploy/Dockerfile -t borealis:$(VERSION) .

## docker-verify: double build, comparaison SHA256 (NFR-09, M4.4)
docker-verify:
	@command -v docker >/dev/null 2>&1 || { echo "Docker absent sur cet hôte."; exit 1; }
	bash deploy/docker-build-verify.sh

## check: gate local complet (build, vet, test, vuln, lint, couverture)
check:
	bash scripts/check.sh

## tidy: go mod tidy
tidy:
	$(GO) mod tidy

## clean: artefacts
clean:
	rm -rf $(BUILD) coverage.out

## help: liste des cibles
help:
	@grep -E '^## ' $(MAKEFILE_LIST) | sed 's/## //'
