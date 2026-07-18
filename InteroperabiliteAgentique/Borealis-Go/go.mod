module github.com/agbruneau/borealis

go 1.26

// Dépendances externes ajoutées à leur premier usage (règle YAGNI) :
//   - github.com/modelcontextprotocol/go-sdk (mcp)      → M0.3+
//   - github.com/a2aproject/a2a-go/v2                   → M2
//   - go.opentelemetry.io/otel + contrib/otelslog       → M0.5
// Versions vérifiées en P0.1 (docs/verification-p01.md).

require (
	github.com/a2aproject/a2a-go/v2 v2.3.1
	github.com/modelcontextprotocol/go-sdk v1.6.1
	go.opentelemetry.io/otel v1.44.0
	go.opentelemetry.io/otel/sdk v1.44.0
	go.opentelemetry.io/otel/trace v1.44.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/google/jsonschema-go v0.4.3 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/segmentio/asm v1.1.3 // indirect
	github.com/segmentio/encoding v0.5.4 // indirect
	github.com/yosida95/uritemplate/v3 v3.0.2 // indirect
	go.opentelemetry.io/auto/sdk v1.2.1 // indirect
	go.opentelemetry.io/otel/metric v1.44.0 // indirect
	golang.org/x/mod v0.33.0 // indirect
	golang.org/x/oauth2 v0.35.0 // indirect
	golang.org/x/sync v0.19.0 // indirect
	golang.org/x/sys v0.45.0 // indirect
)
