# Rapport final — Borealis (clôture M5.7)

> Données 100 % synthétiques, aucun système réel, jamais de PII.
> Vérification finale : gate rejoué, `make e2e` ×3 déterministe, relecture par
> deux réfutateurs frais (grilles David / Carmen). Rapport **fidèle** : les cases
> non vérifiables localement sont marquées ⚠/⛔, jamais présentées comme acquises.

> **Mise à jour post-clôture (2026-07-08).** Un audit du dépôt (voir `PRDPlan-Borealis.md` §8) a, après cette clôture M5.7, porté la couverture de 90,3 % à **96,2 %**, ajouté **ADR-0011** (11 ADR hors gabarit) et étendu la suite de tests (~210 fonctions). Les chiffres ci-dessous reflètent l'état figé à la clôture de M5.7 (2026-07-07).

## 1. Definition of Done (§16 / §1 du plan) — commande → sortie constatée

| # | Case | Commande | Résultat constaté | Statut |
|---|---|---|---|---|
| 1 | `docker compose up` + `make e2e` déterministe | `make e2e` (×3) ; `docker compose up` | e2e **PASS ×3, golden inchangé** ; compose non lancé | ✅ e2e / ⚠ compose (Docker absent) |
| 2 | Découverte A2A ≥ 4 pairs sans registre en dur | `go test ./internal/orchestrator` | Discover vérifie signatures + épingle l'hôte + rejette collisions | ✅ |
| 3 | ≥ 3 délégations + états de Task + input-required | `go test ./internal/a2aserver ./pkg/a2a` | machine à états (terminaux irréversibles), input-required observé | ✅ |
| 4 | ≥ 4 serveurs MCP + ressource + prompt | `go test ./internal/mcpserver` | 4 serveurs, `resources_test.go` (ressource + prompt + progress) | ✅ |
| 5 | Substitution au contrat compatible → e2e vert | `go test ./internal/mcpserver -run Substitution` | Capacity v1→v2 BACKWARD, décision v1 identique | ✅ |
| 6 | Explicabilité fr/en, 100 %, < 500 ms | `go test ./internal/conformance` | golden fr/en (refusé + pré-qualifié), 1000 appels < 500  ms | ✅ |
| 7 | Journal WORM vérifiable, distinct des logs SRE | `go test ./internal/audit` | chaîne de hachage + ancrage d'export ; ADR-0003 (puits distinct) | ✅ |
| 8 | Trace OTel A2A+MCP corrélée par task id | — | instrumentation (spans, traceparent) présente ; **export non câblé** | ⛔ reliquat |
| 9 | Gate local vert (lint, race, vuln, couverture, bench) | `bash scripts/check.sh` | **PASS, couverture 90,3 %**, 0 vuln, lint propre ; `-race` **sauté** (CGO absent) | ✅ (race sauté, signalé) |
| 10 | Build Docker reproductible (SHA256 identiques) | `deploy/docker-build-verify.sh` | script prêt, non exécuté (Docker absent) | ⚠ (Docker absent) |
| 11 | ADR + README pont + matrice sans orphelin | `ls docs/adr` ; `docs/TRACEABILITY.md` | **10 ADR** (9 décisions + 1 écart), README pont, matrice sans orphelin | ✅ |

**Bilan DoD : 8/11 pleinement prouvées ; 2 non vérifiables localement (Docker absent,
cases 1-compose et 10) ; 1 reliquat de code (case 8, export OTLP).** Aucune case
n'est déclarée acquise sans exécution.

## 2. KPI (§17)

| KPI | Type | Mesure |
|---|---|---|
| Latence boucle P99 < 2 s | technique | ✅ **~336 µs/op** (`make bench`, ×6000 de marge) |
| MCP P95 < 500 ms | technique | ✅ `TestToolRoundtripSmoke` (smoke ; renommé lors de l'audit de refactorisation du 2026-07-08, T4) |
| Audit < 200 ms P99 & ≥ 1000/s | technique | ✅ `TestAuditThroughput` |
| Couverture globale ≥ 90 % | technique | ✅ **90,3 %** |
| Sous-cible unitaire ≥ 95 % | technique | cœur métier ≥ 95 % (agent 95,2 %, audit 98,3 %, mcpserver 97,9 %, pep/hitl/resilience/creditdomain 100 %) ; launchers/observability plus bas (glue + reliquat export) tirent le global à 90,3 % |
| « Extractible < 1 h » (David) | pédagogique | **évaluation humaine différée** — patron MCP/A2A isolable (README) |
| « Concepts retrouvables ≥ 90 % » (Aline) | pédagogique | **évaluation humaine différée** — mapping `docs/ARCHITECTURE.md` |

## 3. Répartition de la pyramide de tests

- **180** fonctions de test/bench ; **7** e2e (`test/e2e`, tag e2e) ≈ **4 %**.
- Le gros est unitaire (transports en mémoire) + intégration (fakes + serveurs MCP
  en mémoire). Le ratio e2e (4 %) est en deçà de la cible ~10 % : les 6 scénarios
  §12.4 concentrent l'e2e ; l'essentiel de la vérification est unitaire/intégration
  (plus rapide, déterministe). Signalé sans enjolivure.

## 4. Écarts au PRD (chacun adossé à un ADR)

- **ADR-0001** — validation d'entrée MCP reportée en `IsError`, pas `-32602` (SDK v1.6.1).
- **ADR-0005** — `a2a-go` v2 (le PRD visait v1) ; épinglé v2.3.1.
- **ADR-0008** — module d'identité complet (NHI JIT/HSM/WORM matériel/OAuth multi-hop) reporté au candidat #2.
- **ADR-0010** — mapping ArchiMate dérivé du PRD §6 (sans accès à l'Annexe B / registre §6.6).

## 5. Reliquats et candidats

- **⛔ Export OTLP** non câblé dans les binaires (instrumentation présente ; case DoD 8).
- **⚠ Docker absent de l'hôte** → cases 1 (compose) et 10 (SHA256) non vérifiées localement.
- **`-race` sauté** (CGO/gcc absent) → DoD 9 partielle sur ce point.
- Effet partiel + erreur SSE → double effet possible au rejeu (plafond à-moins-une-fois, documenté).
- **Candidat #2** : module d'identité matériel, crypto PQ réelle, export OTLP.
- **Candidat #3** : échelle, sémantique, agents à logique métier complète (exécuteurs actuels : stubs de réponse).

## 5 bis. Relecture adverse finale (2 réfutateurs frais, M5.7)

Deux réfutateurs à contexte frais (grilles **David** — réutilisation/correction e2e ;
**Carmen** — conformité/honnêteté) ont exécuté les chaînes réelles. Verdict :
*« le démonstrateur tient largement ses promesses »* / *« massivement honnête »*.
Constats **corrigés avant clôture** :

- **MAJEUR (Carmen)** — FR-26 revendiquée déployée alors que le binaire agent servait
  le Mux **nu** (commentaire `auth.go` inversé). **Corrigé** : `SecureMux` câblé dans
  `internal/agentcmd` (card publique, /a2a → 401), vérifié par `TestSecureInvokeRequiresAuth` ;
  couture appelant (orchestrateur → jeton) consignée en reliquat.
- **MINEUR (David/Carmen)** — `make docker-build` pointait `deploy/Dockerfile.multi`
  inexistant. **Corrigé** → `deploy/Dockerfile`.
- **MINEUR (David)** — scénario e2e 5 vacus (n'exerçait pas la dégradation). **Corrigé** :
  étage d'explicabilité injectable, chemin d'échec réellement pris (avis absent,
  décision intacte), contraste explicateur sain.
- **NIT (David)** — README `make` non listé / Windows ; `make audit-export` masquait
  l'échec. **Corrigés** (prérequis + équivalents `go` ; masquage retiré).

Invariants **confirmés tenus** par les réfutateurs (échantillon exécuté) : FR-25
(par le typage), WORM + ancrage (re-forge démasqué), maker-checker AST (durci
pointeur), fail-closed non contournable, anti-rejeu webhook, pseudonymisation KYC,
placeholder PQ actif-vs-réservé.

## 6. Verdict

Six phases (P0 → M5) construites, **chacune gate-verte et vérifiée adversairement**
avant clôture (M0 : 9 constats ; M1 : 11 ; M2 : 14 dont 2 vulns ; M3 : 3 majeurs ;
M4 : 3 majeurs ; tous corrigés, plusieurs **prouvés par mutation**). Les invariants
de la monographie — *découplage, contrat, évolution* — et « jamais un octroi ferme »
tiennent sur tous les chemins testés. Les limites sont **signalées sans détour**.
