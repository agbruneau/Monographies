# Reproductibilité et déploiement (M4.3 / M4.4)

Ce document décrit le déploiement conteneurisé du démonstrateur Borealis et la
stratégie de **build reproductible** (NFR-09), ainsi que ce qui a été — ou non —
vérifié.

> **Données 100 % synthétiques, aucun système réel.** Les extraits de
> configuration sont illustratifs ; aucun secret réel.

## État de vérification

| Élément | Vérifié localement ? | Détail |
|---|---|---|
| Cross-compilation 3 OS + arm64 (NFR-08) | ✅ **oui** | `GOOS=linux\|darwin\|windows GOARCH=amd64 go build ./...` + `linux/arm64` verts (9 binaires) |
| Transport HTTP streamable des serveurs MCP | ✅ **oui** | round-trip en process (`TestServeHTTPRoundtrip`) |
| `docker compose up` (DoD 1) | ⚠ **non** | **Docker absent de l'hôte** — artefacts livrés, exécution non lancée |
| Build reproductible SHA256 (DoD 10, NFR-09) | ⚠ **non** | idem — `deploy/docker-build-verify.sh` prêt, non exécuté |
| 0 egress (NFR-06) | ⚠ **partiel** | réseau `internal: true` déclaré ; test d'hôtes en M5 |
| Export OTLP → Jaeger (NFR-07) | ⛔ **non livré** | spans + `traceparent` instrumentés ; **exportateur non câblé dans les binaires** (reliquat, non imputable à Docker) |

Les cases ⚠ relèvent d'une **contrainte d'hôte** (Docker non installé), pas d'un
défaut de conception. La case ⛔ est un **reliquat de code** (câblage de
l'exportateur), distinct de la contrainte Docker — signalé sans détour.

## Topologie (`deploy/docker-compose.yml`)

Onze services : orchestrateur + 4 agents A2A (HTTP) + 4 serveurs MCP (HTTP
streamable, `MCP_HTTP_ADDR`) + collecteur OTel + Jaeger. Deux réseaux :

- **`mesh`** (`internal: true`) : porte tout le trafic applicatif ; **aucune
  sortie Internet** (0 egress, NFR-06).
- **`edge`** : réservé à l'inbound de l'UI Jaeger (port 16686).

Le code **instrumente les spans** et **propage `traceparent` W3C** de bout en bout
(A2A + MCP, NFR-07 ; `TraceRoundTripper`, `InjectTraceparent`). **Limite honnête :
le câblage de l'exportateur OTLP dans les binaires autonomes n'est pas encore livré.**
Le fournisseur de traces applicatif (`observability.Provider`, jamais instancié par
les binaires) a été retiré (audit de refactorisation, 2026-07-08, B2) faute
d'appelant de production ; l'instrumentation au niveau bibliothèque (propagation)
demeure, mais son câblage à un fournisseur OTel réel (tracer/meter provider +
exportateur `otlptrace`) reste à construire au moment où l'export sera
effectivement branché dans `servercmd`/`agentcmd`. Le collecteur et Jaeger restent
**provisionnés** dans le compose pour l'accueillir, mais **ne reçoivent pas de
traces en l'état**.

## Build reproductible

`deploy/Dockerfile` (multi-étages) produit un binaire par service (argument
`BIN`) avec les garanties de reproductibilité :

- `CGO_ENABLED=0` (binaire statique, pas de dépendance à la libc de l'hôte) ;
- `-trimpath` (chemins de compilation absolus retirés) ;
- `-ldflags "-s -w -buildid="` (pas d'identifiant de build variable) ;
- `SOURCE_DATE_EPOCH` partagé (horodatages figés) ;
- bases **à épingler par digest** (`golang:1.26`, `distroless/static:nonroot`).

Vérification (sur hôte Docker) :

```bash
SOURCE_DATE_EPOCH=1720310400 bash deploy/docker-build-verify.sh mcp-identity
# → deux builds, empreintes SHA256 comparées, doivent être identiques.
```

## Reste à faire (sur hôte doté de Docker)

1. Épingler les images de base par digest (`@sha256:…`).
2. `docker compose -f deploy/docker-compose.yml up` puis `make e2e` (< 15 min).
3. `docker-build-verify.sh` sur les 9 binaires → SHA256 identiques.
4. Test d'egress : aucune connexion sortante hors hôtes déclarés (M5).
