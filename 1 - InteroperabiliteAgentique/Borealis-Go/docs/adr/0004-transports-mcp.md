# ADR-0004 — Choix des transports MCP

- **Statut :** accepté
- **Date :** 2026-07-07
- **Phase :** M1.9

## Contexte

Le go-sdk v1.6.1 offre plusieurs transports MCP : `StdioTransport`,
`CommandTransport` (client → sous-processus stdio), transports en mémoire
(`NewInMemoryTransports`), SSE et *streamable HTTP*. P0.1 a noté que la révision
MCP 2026-07-28 (RC) déprécie Roots/Sampling/Logging — primitives **non utilisées**
ici, donc sans impact.

## Décision

- **Binaires serveurs** (`cmd/mcp-*`) : servis sur **stdio** (`StdioTransport`),
  lancés comme sous-processus par l'orchestrateur / la composition.
- **Tests unitaires et d'intégration** : transports **en mémoire**
  (`mcp.NewInMemoryTransports`) — pas de processus, déterministe, rapide.
- **Tests e2e** (`test/e2e`, tag `e2e`) : **`CommandTransport`** lançant le vrai
  binaire, pour exercer la voie de déploiement stdio réelle.
- SSE / streamable HTTP : non retenus (aucun besoin de transport réseau MCP dans
  le périmètre ; A2A porte les échanges réseau, ch. M2).

## Conséquences

- Un seul mode de transport MCP à opérer (stdio) → surface réduite.
- Le calibre du gate reste local et déterministe (in-memory).
- La visualisation Jaeger et `docker compose` supposent Docker (absent sur l'hôte
  courant — voir docs/verification-p01.md).

## Reversal condition

Un besoin d'exposer un serveur MCP sur le réseau (multi-hôtes, hors sous-processus)
justifierait d'ajouter le transport *streamable HTTP*.

## Alternatives écartées

- **SSE d'emblée** : complexité réseau non requise au niveau MCP.
- **Tout en mémoire** (y compris e2e) : ne prouverait pas la voie stdio réelle.
