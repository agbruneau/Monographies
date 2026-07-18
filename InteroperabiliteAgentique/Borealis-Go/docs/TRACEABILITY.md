# Matrice de traçabilité — Borealis (DoD 11)

Exigence → phase → **code** → **test** → **preuve**. Sans orphelin : chaque
exigence pointe vers un paquet et un test exécuté. Statut ✅ = vérifié par le gate
ou un test ; ⚠ = livré mais non vérifiable localement (contrainte d'hôte) ; ⛔ =
reliquat honnête (non livré, signalé).

> Données 100 % synthétiques, aucun système réel, jamais de PII.

## Exigences fonctionnelles (FR)

| Exigence | Prio | Code | Test | Statut |
|---|---|---|---|---|
| FR-01 card conforme schéma | Must | `pkg/a2a/card.go`, `internal/a2aserver` | `card_test.go`, `agent_test.go` | ✅ |
| FR-02 découverte sans endpoint en dur | Must | `internal/orchestrator` (`Discover`) | `orchestrator_test.go` | ✅ |
| FR-03 signatures JWS+JCS | Must | `pkg/a2a` (`Sign`/`Verify`) | `card_test.go` (rejets alg/clé/altération) | ✅ |
| FR-04 capabilities déclarées | Could | `internal/orchestrator` (`SupportsStreaming`) | `orchestrator_test.go` (choix transport) | ✅ |
| FR-05 message/send + cycle de vie | Must | `pkg/a2a/task.go`, `internal/a2aserver` | `task_test.go`, e2e | ✅ |
| FR-06 referenceTaskIds | Must | `internal/a2aserver` | `agent_test.go` | ✅ |
| FR-07 états terminaux irréversibles | Must | `pkg/a2a/task.go` (`TaskMachine`) | `task_test.go` | ✅ |
| FR-08 streaming SSE | Must | `internal/orchestrator` (`deliver`) | `orchestrator_test.go` (état final replié) | ✅ |
| FR-09 webhook push signé anti-rejeu | Could | `internal/webhook` | `webhook_test.go` (rejeu, NUL, périmé) | ✅ |
| FR-10 AddTool schémas inférés | Must | `internal/mcpserver` (4 serveurs) | `invalidinput_test.go` (rejet input → IsError, ADR-0001) | ✅ |
| FR-11 clients MCP SDK officiel | Must | `internal/agent` | `mcpcall_test.go` + e2e stdio | ✅ |
| FR-12 ressource MCP | Must | `internal/mcpserver/resources.go` | `resources_test.go` | ✅ |
| FR-13 prompt MCP | Must | `internal/mcpserver/resources.go` | `resources_test.go` | ✅ |
| FR-14 NotifyProgress | Must | `internal/mcpserver` | `*_test.go` | ✅ |
| FR-15 erreurs protocole vs métier | Must | `internal/mcpserver/errors.go` | `protocol_test.go`, `invalidinput_test.go` | ✅ |
| FR-16 HITL deux files | Must | `internal/hitl`, `internal/agent` (`amlGate`) | `hitl_test.go`, `orchestrator_test.go` | ✅ |
| FR-17 timeouts/retries/CB | Should | `internal/resilience`, `internal/agent` | `resilience_test.go`, `fault_test.go` (+HANG) | ✅ |
| FR-18 idempotence | Should | `internal/orchestrator/idempotency.go` | `idempotency_test.go` | ✅ |
| FR-19 journal chaîné par hachage | Must | `internal/audit` | `audit_test.go` (tampering, ancrage) | ✅ |
| FR-20 7 champs d'audit | Must | `internal/audit` | `audit_test.go` | ✅ |
| FR-21 explicabilité fr/en | Must | `internal/conformance/explain.go` | `explain_test.go` (golden fr/en, <500ms) | ✅ |
| FR-22 explication hybride | Should | `internal/conformance/explain.go` | `explain_test.go` (golden) | ✅ |
| FR-23 audit-export | Could | `internal/audit` (`ExportJSONL`), `internal/auditexport` | `auditexport_test.go` (rejouable) | ✅ |
| FR-24 maker-checker | Should | `internal/hitl` | `maker_checker_test.go` (topologique AST) | ✅ |
| FR-25 aucun octroi ferme | Won't | `internal/creditdomain` (aucun octroi) | e2e (invariant 1), `explain_test.go` | ✅ |
| FR-26 bearer à audience restreinte | Should | `internal/idpmock`, `internal/a2aserver/auth.go`, **câblé dans `internal/agentcmd`** | `auth_test.go`, `agentcmd_test.go` (`TestSecureInvokeRequiresAuth` : binaire déployé → card publique, /a2a → 401) | ✅ (agent) |

## Exigences non fonctionnelles (NFR)

| Exigence | Code | Test / preuve | Statut |
|---|---|---|---|
| NFR-01 boucle P99 < 2 s | `internal/agent` | `latency_test.go` (~336 µs/op), `make bench` | ✅ |
| NFR-02 MCP P95 < 500 ms | `internal/mcpserver` | `TestToolRoundtripSmoke` (smoke ; SLO mesuré non falsifiable en mémoire, T4) | ✅ |
| NFR-03 audit < 200 ms P99 & ≥1000/s | `internal/audit` | `TestAuditThroughput` | ✅ |
| NFR-04 résilience | `internal/resilience` | `fault_test.go` (arrêt + HANG) | ✅ |
| NFR-05 idempotence | `internal/orchestrator` | `idempotency_test.go` | ✅ |
| NFR-06 0 egress | `deploy/docker-compose.yml` (`internal: true`) | test d'hôtes (M5) | ⚠ (Docker absent) |
| NFR-07 trace ininterrompue | `internal/observability` (traceparent) | instrumentation ; **export non câblé** | ⛔ reliquat |
| NFR-08 portabilité 3 OS | `cmd/*` | cross-compile linux/darwin/windows + arm64 | ✅ |
| NFR-09 build reproductible | `deploy/Dockerfile`, `docker-build-verify.sh` | double build SHA256 | ⚠ (Docker absent) |
| NFR-10 déterminisme e2e | `test/e2e` | golden ×3 stable | ✅ |
| NFR-11 couverture ≥ 90 % | tout | `make coverage` (96,2 %) | ✅ |
| NFR-12 i18n fr/en | `internal/conformance` | `explain_test.go` (golden, `lang`) | ✅ |
| NFR-13 licences | `LICENSE` | Apache 2.0 / CC BY 4.0 | ✅ |
| NFR-14 0 vuln haute/critique | gate | `govulncheck` | ✅ |

## Buts (G) et décisions (ADR)

| Élément | Code / preuve | Statut |
|---|---|---|
| G-1 découverte A2A | `internal/orchestrator` ; ≥ 4 pairs | ✅ |
| G-2 6 scénarios e2e | `test/e2e/scenarios_test.go` | ✅ |
| G-3 4 serveurs MCP + ressource/prompt | `internal/mcpserver` | ✅ |
| G-4 substitution v1→v2 | `internal/mcpserver/capacity_v2.go` | ✅ |
| G-5 explicabilité 100 % | `internal/conformance` | ✅ |
| G-6 gate local vert | `scripts/check.sh` | ✅ |
| G-7 WORM vérifiable | `internal/audit` | ✅ |
| G-8 déterminisme | fakes, golden, horloges injectables | ✅ |
| ADR-0001..0011 | `docs/adr/` (11 fichiers : 9 décisions + 2 écarts M5.3/refactorisation) | ✅ |

## Reliquats honnêtes (candidat #2)

- **⛔ NFR-07** : instrumentation présente (spans, traceparent) mais exportateur
  OTLP non câblé dans les binaires. `observability.Provider` (fournisseur de
  traces/métriques applicatif, jamais instancié en production) a été retiré
  (audit de refactorisation, 2026-07-08, B2) : le câblage à un fournisseur réel
  reste entièrement à faire quand l'export sera branché.
- **⚠ DoD 1/10** : `docker compose up` et build reproductible SHA256 non
  vérifiables localement (Docker absent de l'hôte).
- Effet partiel + erreur → double effet possible au rejeu (plafond à-moins-une-fois).
- **FR-26** : l'agent déployé **applique** la garde bearer (401 sur /a2a non authentifié).
  La **délivrance/attache du jeton côté appelant** (orchestrateur → agent, IdP partagé)
  reste une couture non câblée pour les appels inter-agents vivants.
- Module d'identité complet (NHI JIT/HSM/WORM matériel/OAuth multi-hop) reporté
  ([ADR-0008](adr/0008-report-module-identite.md)).
