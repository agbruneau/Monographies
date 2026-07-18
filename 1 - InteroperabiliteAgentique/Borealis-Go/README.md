# Borealis — démonstrateur d'interopérabilité agentique (MCP + A2A)

> ⚠ **Fiction pédagogique.** *Coopérative financière Boréalis* est une entreprise fictive ;
> toutes les données sont **100 % synthétiques**, aucun système réel, **jamais de PII**.
> Ce dépôt matérialise l'**Annexe B** et le **chapitre 6** de la monographie
> *Interopérabilité agentique* (au sein du dépôt d'écriture *InteroperabiliteAgentique*).

Démonstrateur Go : **5 agents A2A** + **4 serveurs MCP** orchestrant une **pré-qualification de crédit**
(jamais un octroi ferme). Illustre l'invariant transversal de la monographie — *découplage, contrat, évolution* —
et la triade d'identité **KYA / KYH / KYC** avec journal d'audit **WORM** vérifiable par chaîne de hachage.

## Démarrage rapide (< 15 min depuis un clone)

**Prérequis :** Go ≥ 1.22 (dév. sous 1.26) et **GNU Make** (pour les cibles ci-dessous).
Pour la voie conteneurisée : Docker + compose.

> **Windows sans `make`** : le gate passe par `scripts/check.ps1` ; les autres
> cibles se lancent directement en Go —
> `go test -tags e2e ./test/e2e/...` (e2e),
> `go run ./cmd/orchestrator --audit-export` (export),
> `go test -bench=. -run='^$' ./...` (bench).

```bash
# 1. Gate local (build, vet, test, govulncheck, lint, couverture ≥ 90 %)
make check

# 2. Scénarios e2e déterministes (6 scénarios PRD §12.4) — sans Docker
make e2e

# 3. Export d'audit rejouable (JSONL, chaîne de hachage vérifiable)
make audit-export

# 4. Benchmarks (latence boucle P99 < 2 s, débit d'audit)
make bench
```

Voie conteneurisée (topologie complète — **nécessite Docker**, cf.
[`docs/REPRODUCIBILITY.md`](docs/REPRODUCIBILITY.md)) :

```bash
docker compose -f deploy/docker-compose.yml up   # orchestrateur + 4 agents + 4 MCP + otel + jaeger
```

> **Extraire un patron en < 1 h** (persona David) : un serveur MCP autonome vit
> dans `internal/mcpserver/*.go` + `cmd/mcp-*/main.go` ; un agent A2A dans
> `internal/a2aserver` + `cmd/agent-*`. Les deux se lancent en stdio (défaut) ou
> HTTP (`MCP_HTTP_ADDR` / `ADDR`).

## Ce que le démonstrateur prouve

- **Découverte A2A** par Agent Card signée (JWS ES256 + JCS), sans registre en dur.
- **Délégation** MCP/A2A avec cycle de vie de Task (états terminaux irréversibles).
- **Substitution** d'un agent au contrat compatible (Capacity v1→v2) sans casse.
- **Gouvernance** : journal WORM, PEP fail-closed, HITL deux files, explicabilité fr/en, biais Gini.
- **Sécurité** : STRIDE (spoofing/tampering/elevation), bearer à audience restreinte, webhook anti-rejeu.

## Pont vers la monographie (par section)

| Concept démonstrateur | Paquet Go | Renvoi monographie |
|---|---|---|
| Agents A2A, serveurs MCP, orchestrateur | `internal/{a2aserver,mcpserver,orchestrator}` | PRD §6 ; ch. 6 ; Annexe B |
| Contrat : Agent Card + compat BACKWARD | `pkg/{a2a,mcpcontract}` | ch. 3 (MCP/A2A) ; Annexe B §6 |
| Triade d'identité NHI/KYA/KYH/KYC | `internal/audit` (champs) | Annexe B §3 |
| Journal WORM (chaîne de hachage + ancrage) | `internal/audit` | PRD §9.4 ; Annexe B §3.7 |
| Calibrage d'autonomie L0–L3, PEP, HITL | `internal/{pep,hitl}` | Annexe B §1.3, §12 |
| Explicabilité, biais, factsheet | `internal/conformance` | Annexe B §11 |
| Crypto-agilité (placeholder PQ) | `pkg/a2a/cryptosuite.go` | ch. 7 (horizon PQ) ; ADR-0006 |

Mapping complet composant → paquet : [`docs/ARCHITECTURE.md`](docs/ARCHITECTURE.md).
Traçabilité exigence → code → test → preuve : [`docs/TRACEABILITY.md`](docs/TRACEABILITY.md).
Présentation synthèse (diapositives) : [`Présentation Borealis-Go.pdf`](Présentation%20Borealis-Go.pdf).

## Structure

```
cmd/            9 binaires (4 mcp-*, 4 agent-*, orchestrator)
internal/       cœur : agent, a2aserver, mcpserver, orchestrator, audit,
                conformance, hitl, pep, resilience, observability, webhook,
                idpmock, security, creditdomain, fixtures, auditexport, ...
pkg/            a2a (cards JWS/JCS), mcpcontract (compat BACKWARD)
deploy/         docker-compose, Dockerfile reproductible, otel
docs/           ARCHITECTURE, TRACEABILITY, REPRODUCIBILITY, adr/ (11 ADR)
test/e2e/       6 scénarios + golden canonique
scripts/        check.{sh,ps1} (gate local)
```

## État et reliquats honnêtes

Construction pilotée par `PRDPlan-Borealis.md` (phases P0 → M5, vérification adverse
à chaque clôture). Reliquats signalés sans détour (cf. `docs/TRACEABILITY.md`) :

- **⚠ Docker absent de l'hôte** → `docker compose up` et build reproductible SHA256 non vérifiés localement.
- **⛔ Export OTLP** non câblé dans les binaires (instrumentation présente ; reliquat candidat #2).
- Module d'identité complet (NHI JIT/HSM/WORM matériel) reporté (ADR-0008).

## Licences

Code sous **Apache 2.0** (`LICENSE`) ; documentation et données synthétiques sous **CC BY 4.0**.
