# PRDPlan — Borealis : plan d'exécution du PRD v2.0 (mode ultracode)

| | |
|---|---|
| **Titre** | Plan d'exécution Borealis — du PRD au dépôt exécutable |
| **Source de vérité** | [`PRD-Borealis.md`](PRD-Borealis.md) v2.0 (2026-07-07). En cas de conflit entre ce plan et le PRD, **le PRD prime** ; signaler le conflit dans le rapport de phase. |
| **Exécutant cible** | Claude Opus 4.8 (Claude Code), **mode ultracode** — orchestration multi-agents par phases, vérification adverse avant clôture de chaque phase |
| **Version** | 1.0 |
| **Date** | 2026-07-07 |
| **Livrable** | Dépôt Go `borealis/` conforme à la Definition of Done du PRD §16 |
| **Langue** | Français (commits, ADR, documentation) ; code et identifiants en anglais |

> **Avertissement.** Reprend l'avertissement de fiction du PRD : données 100 % synthétiques, aucun système réel, jamais de PII. Le dépôt `borealis/` est **distinct** du dépôt d'écriture *InteroperabiliteAgentique* (qui ne contient aucun code applicatif — ne jamais y créer le code).

---

## 0. Suivi des phases

Dépôt de code : initialement `../borealis` (frère du dépôt d'écriture), rapatrié depuis dans `Borealis-Go/` au sein du dépôt d'écriture (commit `613add3`) — les mentions « répertoire frère » ci-dessous (P0.2) sont conservées comme trace du plan d'origine. Statut mis à jour à la clôture *prouvée* de chaque phase.

| Phase | Portée | Statut | Preuve |
|---|---|---|---|
| **P0** | Amorçage + levée des ⚠ SDK | ✅ **Terminé** (2026-07-07) | Gate local vert ; `docs/verification-p01.md` ; commit `13cdfea` |
| **M0** | Walking skeleton | ✅ **Terminé** (2026-07-07) | Gate vert (couv. **99,3 %**) ; e2e stdio vert ; vérif. adverse (9 constats corrigés) ; commits `c5d3ae2`, `9d73050`. ADR-0001..0003. |
| **M1** | MCP production-grade | ✅ **Terminé** (2026-07-07) | 4 serveurs MCP (§6.2) + fixtures seedées + ressource/prompt/progress + résilience + `pkg/mcpcontract` (compat BACKWARD) + registre tiers TIC ; gate **96,2 %**, e2e vert ; vérif. adverse (11 constats corrigés) ; commits `dd7eebc`…`99f1288`. ADR-0004/0009. |
| **M2** | A2A : découverte, délégation, streaming | ✅ **Terminé** (2026-07-07) | `pkg/a2a` (JWS ES256+JCS) + 5 agents serveurs A2A + orchestrateur (découverte par config, routage **anti-SSRF/anti-collision**) + cycle de vie de Task + SSE + `referenceTaskIds` + bearer JWT à audience restreinte + `traceparent` W3C ; gate **94,4 %** ; vérif. adverse (**14 constats corrigés**, dont 2 vulns sécurité) ; commits `9f4d35e`…`04867fa`. ADR-0004/0005. |
| **M3** | Gouvernance, audit, explicabilité | ✅ **Terminé** (2026-07-07) | Audit chaîne de hachage + ancrage d'export tamper-evident (FR-19/20, G-7, NFR-03) + PEP fail-closed L1→L0 + HITL deux files + explicabilité fr/en golden (sélection `lang`, invariant 1 verrouillé) + biais Gini + AML→HITL (escalade avant scellement) + maker-checker topologique (AST durci) + suite STRIDE ; gate **94,3 %** ; vérif. adverse en 2 passes (1 MAJEUR audit + 2 MAJEURS de tautologie de tests, corrigés et **prouvés par mutation**) ; commits `b3f5862`…`c86c919`. ADR-0008. |
| **M4** | Robustesse, reproductibilité, substitution | ✅ **Terminé** (2026-07-07) | Idempotence des délégations (FR-18) + injection de panne/disjoncteur (FR-17, +test HANG) + substitution Capacity v1→v2 au contrat BACKWARD (G-4, US-9) + webhook signé anti-rejeu double garde (FR-09) + 9 binaires (transport HTTP MCP) cross-compilés 3 OS+arm64 (NFR-08) + artefacts docker/build reproductible (⚠ non exécutés — Docker absent) ; ADR-0006/0007 ; gate **90,7 %** ; vérif. adverse 3 réfutateurs (**3 MAJEURS** corrigés + mutation-prouvés : deliver état final, anti-rejeu horodatage, sur-revendication OTLP reformulée) ; commits `b551054`…`7104ed8`. **Reliquats honnêtes** : `docker compose up`/SHA256 (DoD 1/10) non vérifiables localement ; export OTLP→Jaeger non câblé dans les binaires (⛔ reliquat). |
| **M5** | Finition, documentation, traçabilité | ✅ **Terminé** (2026-07-07) | 6 scénarios e2e §12.4 + golden canonique déterministe (×3) ; bench boucle P99 ~336 µs ≪ 2 s (NFR-01) ; audit-export JSONL rejouable (FR-23) ; ARCHITECTURE.md + README pont + matrice de traçabilité **sans orphelin** ; placeholder PQ (ES256 actif / ML-DSA-65 réservé, NG-9) ; 10 ADR (9 + écart 0010) ; **vérification finale adverse** (2 réfutateurs frais → 1 MAJEUR honnêteté FR-26 câblé réellement + mineurs, corrigés) ; `docs/RAPPORT-FINAL.md` ; gate **90,3 %** ; commits `841a1fc`…`a22920f`. **DoD : 8/11 pleinement prouvées** ; reliquats honnêtes : `docker compose`/SHA256 (Docker absent) et export OTLP (⛔ non câblé). |

**Contraintes d'hôte constatées en P0.1** (à lever sur hôte CGO/Linux + Docker) : **Docker absent** → DoD 1 (`docker compose up`) et 10 (build reproductible SHA256) non vérifiables localement ; **CGO/gcc absent** → `-race` sauté avec avertissement (DoD 9 partiel). Tout le reste du gate s'exécute.

---

## 1. Mission et critère de fin

**Mission.** Construire le démonstrateur Borealis décrit par le PRD : 5 agents A2A + 4 serveurs MCP en Go, cas de pré-qualification de crédit, au calibre d'ingénierie FibGo.

**Critère de fin (seuil de vérification global).** Les 11 cases de la DoD projet (PRD §16) sont cochées et **prouvées** :

1. `docker compose up` démarre tout ; `make e2e` déroule le scénario de crédit de bout en bout, de façon déterministe.
2. Découverte A2A par Agent Card (≥ 4 pairs) sans registre codé en dur.
3. ≥ 3 délégations A2A chaînées ; états de Task couverts ; ≥ 1 chemin `input-required`.
4. ≥ 4 serveurs MCP avec outils validés ; ≥ 1 ressource + ≥ 1 prompt.
5. Substitution d'un agent au contrat compatible → e2e reste vert.
6. Explicabilité fr/en, couverture 100 %, < 500 ms.
7. Journal WORM vérifiable (chaîne de hachage) distinct des logs SRE.
8. Trace OTel ininterrompue A2A + MCP corrélée par `task id`.
9. Gate local vert : lint, race, vuln, couverture ≥ 90 %, bench sans régression.
10. Build Docker reproductible (SHA256 identiques).
11. 9 ADR rédigées ; README pont-monographie ; matrice de traçabilité sans orphelin.

« Prouvé » = la commande de vérification a été exécutée et sa sortie constatée ; jamais d'affirmation sans exécution.

---

## 2. Protocole d'exécution ultracode

Chaque phase suit le même cycle — aucune phase ne se clôt sans son étape 5 verte :

1. **Planifier.** Relire la section du PRD couverte par la phase ; décomposer en tâches ; identifier les tâches indépendantes (marquées ‖ ci-dessous).
2. **Fan-out.** Répartir les tâches ‖ sur des sous-agents parallèles. Les tâches ‖ de ce plan touchent des répertoires **ou fichiers** disjoints (`cmd/mcp-*`, `cmd/agent-*` ; pour M3.4 ‖ M3.5, fichiers distincts de `internal/conformance` précisés dans les tâches) ; si deux tâches partagent un même fichier, les sérialiser ou isoler par worktree et intégrer ensuite.
3. **Intégrer.** Fusionner, résoudre, faire compiler l'ensemble.
4. **Vérifier (adverse).** (a) Gate local : `scripts/check.{sh,ps1}` (lint, `-race -short`, vuln, couverture au plancher de la phase, bench). (b) Agents réfutateurs : ≥ 2 sous-agents à contexte frais tentent de **casser** les critères de sortie de la phase (entrées invalides, card non signée, réécriture d'audit, arrêt d'un composant…) et vérifient la conformité au PRD (IDs d'exigences couverts, MoSCoW respecté). Un critère réfuté → retour à l'étape 1.
5. **Clore.** (a) Pour chaque FR livrée dans la phase, dérouler la **DoD par exigence** du PRD §16 (10 items : mappée à une user story + test d'acceptation ; `go test -v -race -cover` local ; couverture ; golden ; bench si perf-sensible ; ADR si décision ; doc à jour ; piste d'audit de l'action ; renvois PRD/Annexe B validés ; commit conventionnel en français). (b) Les commits sont **par tâche/exigence** au fil de la phase ; le commit de clôture (ex. `M1 — serveurs MCP et contrats`) est le dernier, après gate vert. (c) Rapport de phase (critères → preuve) et mise à jour de la matrice de traçabilité (§5).

**Règles non négociables** (violations = arrêt et correction immédiate) :

- **Jamais d'octroi ferme** : aucun chemin de code d'octroi ; le e2e s'arrête à la pré-évaluation (FR-25, invariant 1).
- **Déterminisme** : le fake `Reasoner` est le défaut ; `make e2e` × 3 → **sorties canoniques** identiques (G-8, NFR-10). La sortie canonique exclut les champs volatils (timestamps, `trace_id`, `task id`, hachages qui en dépendent) : horloge et générateur d'IDs injectables et figés en mode e2e, ou normalisation avant comparaison au golden `decision_matrix.json` ; `make e2e` émet cette sortie canonique. Un LLM réel n'entre jamais dans le gate.
- **Golden immuables — à partir de la clôture de la phase qui couvre la fonctionnalité concernée** : avant cette clôture, le golden évolue librement ; après, toute modification exige une ADR justificative (esprit FibGo, PRD §12.3).
- **Could décidables** : FR-04, FR-09, FR-23 s'implémentent **par défaut** ; si un blocage technique persiste après une itération complète de l'étape 4, abandonner et marquer « non couvert (Could) » dans la matrice et le rapport de phase.
- **Aucun PII réel** ; générateur de données seedé et déterministe (PRD §9.5).
- **Pas de CI distante** : la rigueur est locale (`scripts/check`), comme FibGo.
- **TDD** : pour chaque FR, écrire d'abord le test d'acceptation qui échoue (le critère d'acceptation du PRD §8 en donne l'oracle), puis le faire passer. Répartition cible de la pyramide du PRD §12.1 (~70 % unitaire avec `mcp.NewInMemoryTransports()`, ~20 % intégration avec fakes LLM + serveurs MCP en mémoire, ~10 % e2e), mesurée au rapport M5.7.
- **Couverture > 90 %** (directive utilisateur du 2026-07-07 ; **prime sur le PRD**, relève NFR-11 de 80 → **90 % global**) : plancher du gate à **90 % dès M0** ; unitaire ≥ 95 % maintenu ; `main()` gardés minces (logique en `run(ctx)` testable) pour l'atteindre. Toute phase se clôt à ≥ 90 % de couverture globale.
- **Écarts** : si une exigence s'avère inapplicable telle quelle (API de SDK différente du PRD, primitive dépréciée…), ne pas dévier en silence — rédiger l'ADR de l'écart, ajuster la matrice, le signaler au rapport de phase.

---

## 3. Phase P0 — Amorçage et levée des ⚠ (préalable à M0)

### P0.1 Vérification adverse des ⚠ du PRD (bloquant)

Le PRD marque ⚠ tous les énoncés SDK non couverts par les sources vérifiées de la monographie (PRD §2.2, §7.1, §7.3, §20). Avant toute ligne de code, vérifier sur le web (dépôts officiels) :

| À vérifier | Attendu (PRD) | Si écart |
|---|---|---|
| `github.com/modelcontextprotocol/go-sdk` : versions, gel d'API, licence | v1.0.0 (gel) puis v1.6.x ; MIT + Apache 2.0 | Épingler la version stable réelle ; noter au rapport P0.1 (alimentera ADR-0001 à sa rédaction, M0.6) |
| `a2a-go` : dépôt canonique (a2aproject vs a2aserver), version | v1.x (épingler la révision courante) | Choisir le dépôt actif ; noter au rapport P0.1 (alimentera ADR-0005, M2.8) |
| Statut OAuth client du SDK MCP (`auth`/`oauthex`) | Non stabilisé côté client | S'il est stabilisé : le noter, sans élargir le périmètre (FR-26 reste la couture) |
| Révision MCP 2026-07-28 (RC) : publiée en finale ? Dépréciations ? | RC du 21 mai 2026 ; Roots/Sampling/Logging *Deprecated* | Noter au rapport P0.1 (alimentera ADR-0004, M1.9) |
| `google/adk-go` v2.x | Optionnel, jeune | Sans impact structurel (derrière `Reasoner`) |

**Critère de sortie.** P0.1 produit un **rapport de vérification** (versions, sources datées) conservé dans le dépôt de code (`docs/verification-p01.md`) et repris par les ADR concernées à leur rédaction. **Ne jamais modifier le PRD** (il vit dans le dépôt d'écriture) : tout écart au PRD se consigne dans l'ADR correspondante et au rapport de phase. **Écart majeur** (ex. gel d'API démenti, dépôt A2A abandonné) → s'arrêter et rapporter avant M0.

### P0.2 Squelette du dépôt

Créer `borealis/` en **répertoire frère du dépôt d'écriture** (ex. `../borealis`), `git init` avec branche `main`, premier commit après le squelette — jamais dans le dépôt d'écriture. Arborescence exacte du PRD §14 : `cmd/` (9 binaires), `internal/` (agent, a2aserver, mcpserver, audit, observability, creditdomain, conformance, hitl, testfakes), `pkg/` (a2a, mcpcontract), `api/`, `deploy/`, `docs/adr/`, `test/{e2e,fixtures,golden}`, `scripts/check.{sh,ps1}`, `Makefile` (build, test, race, vuln, coverage, bench, lint, e2e, audit-export, docker-*), `.golangci.yml`, `go.mod` (Go 1.22+, dépendances épinglées P0.1), `LICENSE` (Apache 2.0 code ; CC BY 4.0 docs/données), `CLAUDE.md` (règles §2 de ce plan), `README.md` (squelette).

**Critères de sortie.** `make lint` et `go test ./...` verts sur le squelette ; `scripts/check.sh` et `scripts/check.ps1` s'exécutent (plancher de couverture **90 % dès M0** — directive utilisateur, défaut du gate ; le gate tolère l'absence de base benchstat — *skip* avec avertissement jusqu'à M1 — et l'absence de paquets couvrables au squelette) ; premier commit.

---

## 4. Phases M0 → M5

Fenêtres : effort nominal ~2 sem/phase, borne calendaire 12 mois (PRD §15) — sans objet pour une exécution agentique ; **l'ordre et les critères de sortie, eux, sont normatifs**.

### M0 — Walking skeleton

Un chemin de bout en bout minimal avant toute richesse.

| # | Tâche | Exigences |
|---|---|---|
| M0.1 | `internal/creditdomain` : modèle métier pur (Applicant, Decision, ratios), sans I/O ; tests unitaires | — |
| M0.2 | `internal/agent` : interface `Reasoner` (PRD §7.2) + `internal/testfakes/reasoner_stub.go` (fake déterministe par scénario `approved`/`denied`/`escalate`) | G-8, ADR-0002 |
| M0.3 | Serveur MCP Identity Hub (`cmd/mcp-identity`) minimal : `verify_identity` via `AddTool[In, Out]`, transport `stdio` | FR-10 (amorce) |
| M0.4 | Agent KYC minimal (client MCP via SDK officiel) + orchestrateur minimal (boucle L0, appel direct — A2A vient en M2). L'oracle FR-11 porte sur les **agents consommateurs MCP** (KYC, Scoring, Politique — PRD §6.1/§6.2) ; l'Orchestrateur et l'Explicabilité n'en consomment pas par conception. | FR-11 (amorce) |
| M0.5 | `internal/observability` : `slog` + OTel (OTLP HTTP) + `otelslog` ; compose minimal (otel-collector + Jaeger) | NFR-07 (amorce) |
| M0.6 | ADR-0001 (abstraction protocole), ADR-0002 (Reasoner), ADR-0003 (audit ≠ observabilité) | §19 |
| M0.7 | `internal/audit` : interface et journal **minimal** append-only (champs KYA/KYH/KYC/ts/action/résultat/version) — la chaîne de hachage et le vérificateur arrivent en M3.1 ; permet aux oracles « + entrée d'audit » de M2 de passer | FR-19/FR-20 (amorce) |

**Critères de sortie (PRD §15/M0).** `go test` vert ; 1 appel MCP réel de bout en bout ; trace visible dans Jaeger ; couverture ≥ 90 % ; ADR-0001..0003 rédigées (format PRD §14, avec *Reversal condition*).

### M1 — MCP production-grade

| # | Tâche | Exigences |
|---|---|---|
| M1.1 | Générateur de fixtures seedé : `borrowers.json` (10 000), `policies.csv` (50), `history.jsonl` (~100 000 ; Git LFS si > 100 Mo), `decision_matrix.json` (~30 cas golden) | §9.5, NFR-10 |
| M1.2 ‖ | `cmd/mcp-bureau` (`get_credit_report`) | FR-10, §6.2 |
| M1.3 ‖ | `cmd/mcp-capacity` (`compute_capacity`, PMT/DTI — PRD §9.2) | FR-10, §6.2 |
| M1.4 ‖ | `cmd/mcp-policy` (`evaluate_policy`) | FR-10, §6.2 |
| M1.5 ‖ | `cmd/mcp-identity` complété (`check_watchlists`, listes PPE/OFAC simulées) | FR-10, §6.2 |
| M1.6 | ≥ 1 ressource MCP (`credit:///application/{id}/assessment`) + ≥ 1 prompt (gabarit d'avis) + `NotifyProgress` | FR-12, FR-13, FR-14 |
| M1.7 | Erreurs : protocole (JSON-RPC -32602) vs métier (`isError`) ; aucun `panic` ; timeouts + retries + circuit-breaker | FR-15, FR-17, NFR-04 |
| M1.8 | `pkg/mcpcontract` : validation de schéma + **contrôle de compatibilité BACKWARD** (v2 accepte toute entrée valide v1 **et** ses sorties restent lisibles par un consommateur v1), branché sur les 4 serveurs et exécuté par le gate | §12.2, G-4 (amorce) |
| M1.9 | Registre YAML des serveurs MCP (tiers TIC : SLO, énum d'erreurs, version de contrat) ; ADR-0004 (transports), ADR-0009 (tiers TIC) | §11.4, §19 |

M1.2-M1.5 en fan-out parallèle (répertoires disjoints) ; chaque sous-agent livre serveur + tests + benchmark. Périmètre des tests du fan-out : rejet d'input invalide (-32602) via les schémas inférés du SDK ; la validation `pkg/mcpcontract` et le contrôle BACKWARD sont branchés **ensuite**, en M1.8, sur les 4 serveurs.

**Critères de sortie (PRD §15/M1).** Chaque outil validé par test de contrat (rejet d'input invalide → -32602) ; P95 appel outil < 500 ms (NFR-02, métrique `mcp_tool_latency`) ; couverture ≥ 90 % ; `make bench` établit la base de référence.

### M2 — A2A : découverte, délégation, streaming

| # | Tâche | Exigences |
|---|---|---|
| M2.1 | `pkg/a2a` : `AgentCard` + `Validate()` contre le schéma officiel A2A v1.0 ; signature/vérification JWS (RS256 ou ES256) + canonicalisation JCS (RFC 8785), stdlib | FR-01, FR-03 |
| M2.2 ‖ | Les 5 agents comme serveurs A2A (`/.well-known/agent-card.json`, `skills` du PRD §6.1) sur `a2a-go` | FR-01, §6.1 |
| M2.3 | Orchestrateur : découverte par card (URL de base configurées, **aucun endpoint en dur**), vérification de signature avant routage, sélection de `skills` | FR-02, FR-03, G-1 |
| M2.4 | Délégation `message/send` + cycle de vie de Task (8 états ; terminaux irréversibles — test de machine à états) ; chaînage `referenceTaskIds` porté par le Message (PRD §9.3) ; `contextId` par dossier | FR-05, FR-06, FR-07, G-2 |
| M2.5 | Streaming SSE `message/stream` (≥ 2 `TaskStatusUpdateEvent`) | FR-08, US-8 |
| M2.6 | FR-26 : bearer JWT court via IdP mock **in-process** (`internal/idpmock`, clés dev — ni binaire ni service compose supplémentaire ; ajout à l'arborescence consigné au rapport de phase), **audience restreinte** (claim `aud` par destinataire, RFC 8707) ; rejet 401 + entrée dans le journal minimal M0.7 sur jeton absent/invalide/`aud` non conforme | FR-26, §11.1 |
| M2.7 | NHI placeholder : champs `kya`/`owner_kyh` propagés jusqu'au journal (minimal M0.7 ; chaîne de hachage en M3.1) ; propagation `traceparent` W3C dans tous les appels A2A et MCP (test de contrat d'en-têtes) | §11.2, NFR-07, §12.2 |
| M2.8 | ADR-0005 (a2a-go) ; FR-04 (capabilities déclarées — Could, règle de décidabilité du §2) | §19, FR-04 |

**Critères de sortie (PRD §15/M2).** Découverte de ≥ 4 pairs ; scénario happy path e2e vert (fake LLM ; en M2, l'agent Explicabilité retourne un **avis stub** à texte fixe — les golden fr/en et la couverture 100 % arrivent en M3.4) ; états de Task couverts : `submitted → working → completed` observés, plus un test d'intégration qui force et **observe** la transition `working → input-required` (assertion sur l'état ; la release et les deux files arrivent en M3.3) ; card à signature invalide rejetée (test négatif) ; couverture ≥ 90 %.

### M3 — Gouvernance, audit, explicabilité

| # | Tâche | Exigences |
|---|---|---|
| M3.1 | `internal/audit` complété (sur la base M0.7) : chaîne de hachage (`entry_hash` = sha256 des champs, `prev_hash` chaîné — formule PRD §9.4) sur les 7 champs ; vérificateur de chaîne ; **distinct** du puits OTel ; micro-benchmark : P99 < 200 ms **et** ≥ 1000 écritures/s | FR-19, FR-20, G-7, NFR-03 |
| M3.2 | PEP minimal : garde non contournable avant tout appel d'outil (identifie, autorise, journalise) + **retombée fail-closed L1→L0 journalisée** quand le PEP/orchestrateur est indisponible (testée par le scénario M5.1) | invariant 3, §6.6, §12.4 |
| M3.3 | HITL : **deux files distinctes** proposition/commande (`internal/hitl`) ; `input-required` ; aucun message en file de commande sans release journalisée | FR-16, US-7 |
| M3.4 ‖ | Explicabilité fr/en (`internal/conformance/explain*.go`) : avis hybride (règles + importance + contexte), **sélection de la langue par paramètre `lang`** (US-11), golden fr et en, couverture 100 %, < 500 ms | FR-21, FR-22, G-5, NFR-12, US-11 |
| M3.5 ‖ | Détecteur de biais build-time (`internal/conformance/bias*.go` + `factsheet*.go` — fichiers disjoints de M3.4) : divergence **Gini** par groupe, rapport CSV, alerte > 15 % ; factsheet minimale du reasoner (E-23 illustratif) | §11.4 |
| M3.6 | Seuils AML en cumul : `amlScore > 0.7` → refus motivé ; `> 0.8` → refus **et** escalade HITL avant scellement | §11.4, FR-16 |
| M3.7 | Maker-checker : test topologique — maker = agents préparateurs (KYC/Scoring/Politique), checker = point de release HITL ; aucun agent n'appelle directement un pair ni la file de commande (tout flux passe par l'orchestrateur ; la release ne se déclenche que par le canal HITL) | FR-24 |
| M3.8 | Tests STRIDE (PRD §11.3) : spoofing (card falsifiée), tampering (réécriture d'audit → chaîne rompue), elevation (`aud` étranger → 401) ; ADR-0008 (**objet : report du module d'identité complet — NHI JIT/HSM/WORM matériel/OAuth multi-hop — au candidat #2**, PRD §19) | §11.3, §19 |

**Critères de sortie (PRD §15/M3).** Vérification de la chaîne de hachage verte (réécriture détectée) ; 100 % des décisions expliquées (fr et en, golden, sélection par `lang`) ; scénario AML → HITL vert ; maker-checker testé ; écriture d'audit < 200 ms P99 **et** ≥ 1000 écritures/s (micro-benchmark).

### M4 — Robustesse, reproductibilité, substitution

| # | Tâche | Exigences |
|---|---|---|
| M4.1 | Idempotence des délégations (rejeu même `contextId`+params → un seul effet, doublon signalé à l'audit) | FR-18, NFR-05 |
| M4.2 | Injection de panne : serveur MCP arrêté → retry borné, circuit-breaker ouvert, échec propre ; détection < 30 s | FR-17, NFR-04 |
| M4.3 | `docker compose` complet (orchestrateur + 4 agents + 4 MCP + otel + jaeger) ; **0 egress hors hôtes déclarés** (réseau fermé + test) ; portabilité NFR-08 : cross-compilation `GOOS=linux|darwin|windows go build ./...` verte + images `docker buildx` linux/amd64 + arm64 (l'exécution sur les 2 autres OS est rapportée « non vérifiée localement ») | NFR-06, NFR-08 |
| M4.4 | Build reproductible : `Dockerfile.multi` (base épinglée, `SOURCE_DATE_EPOCH`), `docker-build-verify.sh` → double build, SHA256 identiques ; rédaction de `docs/REPRODUCIBILITY.md` (arborescence PRD §14) | NFR-09 |
| M4.5 | **Substitution Scoring v1→v2** au contrat compatible (BACKWARD, M1.8) sans modifier l'orchestrateur → e2e vert ; ADR-0006 (crypto-agilité : suite d'algorithmes dans le contrat, placeholder PQ), ADR-0007 (stratégie de test) | G-4, US-9, §19 |
| M4.6 | FR-09 (webhook push signé, anti-rejeu) — Could : implémenter par défaut (règle de décidabilité du §2) | FR-09, US-10 |

**Critères de sortie (PRD §15/M4).** `docker compose up` + `make e2e` < 15 min depuis un clone ; double build SHA256 identiques ; substitution v1→v2 sans casse ; `govulncheck` vert ; gate local vert.

### M5 — Finition, documentation, traçabilité

| # | Tâche | Exigences |
|---|---|---|
| M5.1 | Les **6 scénarios e2e** (PRD §12.4) : happy path ; refus AML ; timeout modèle ; signaux conflictuels ; échec Explicabilité ; **fail-closed** (PEP/orchestrateur arrêté → aucun agent n'agit, retombée L1→L0 journalisée — mécanisme livré en M3.2) — 1 test Go par scénario, fixtures isolées ; golden `test/golden/trace-approved.json` (trace canonique du happy path, PRD §12.3) | §12.4, G-2, G-8 |
| M5.2a | Latence boucle P99 < 2 s constatée (`make bench`) — obligatoire | NFR-01 |
| M5.2b | `make audit-export` (JSONL rejouable) — Could, règle de décidabilité du §2 | FR-23 |
| M5.3 ‖ | `docs/ARCHITECTURE.md` : mapping ArchiMate avec les stéréotypes **verbatim du PRD §6.6** (registre ch.6 §6.1.9 cité par le PRD). **Repli sans accès à l'Annexe B** (cas nominal de l'exécutant) : dériver le schéma maître des diagrammes du PRD §6.3-6.5 (C4, séquence, flux) et remplacer la table « vue A-L → paquet Go » par une table « composant du PRD §6 → paquet Go » ; consigner ce repli en ADR d'écart | §14, KPI §17 |
| M5.4 ‖ | README pont-monographie : démarrage < 15 min, renvois PRD/Annexe B/ch.6 par section, avertissement de fiction | §14, persona Marc |
| M5.5 | Matrice de traçabilité **sans orphelin** : ce plan §5, complétée (exigence → code → test → preuve) et versionnée dans `docs/` | DoD 11 |
| M5.6 | Placeholder post-quantique (suite d'algorithmes encodée dans le contrat — aucune crypto PQ réelle, NG-9) ; revue finale des 9 ADR ; audit final style FibGo | §19, NG-9 |
| M5.7 | **Vérification finale adverse** : rejouer intégralement la DoD §1 de ce plan, chaque case avec sa commande et sa sortie constatée ; `make e2e` × 3 (sorties canoniques identiques) ; relecture par ≥ 2 agents réfutateurs frais (personae David et Carmen du PRD §4 comme grilles de lecture) ; rapporter la répartition de la pyramide de tests et les sous-cibles de couverture (global ≥ 90 %, unitaire ≥ 95 %, intégration 60-70 % — NFR-11/KPI §17) | §16, §17 |

**Critères de sortie.** DoD projet 11/11 prouvée ; KPI §17 : les KPI **techniques** mesurés par commande ; les KPI **pédagogiques** nécessitant un tiers ou l'Annexe B (« extractible < 1 h », « concepts retrouvables ≥ 90 % ») rapportés « évaluation humaine différée », avec leurs éléments préparatoires (guide d'extraction pas-à-pas dans le README ; mapping §6.6 complété) ; rapport final : tableau critère → commande → résultat.

---

## 5. Matrice de traçabilité (exigence → phase → vérification)

À compléter au fil des phases (colonnes code/test/preuve) ; les orphelins sont interdits à M5.

| Exigence | Priorité | Phase | Vérification (oracle du PRD §8/§10) |
|---|---|---|---|
| FR-01 card conforme schéma | Must | M2.1-M2.2 | `GET` sans auth → JSON valide contre le schéma officiel |
| FR-02 découverte sans endpoint en dur | Must | M2.3 | Retrait d'un agent de la config → plus routé, sans modif de code |
| FR-03 signatures JWS+JCS | Must | M2.1, M2.3 | Card à signature invalide → rejet + audit |
| FR-04 capabilities déclarées | Could | M2.8 | Transport choisi selon la capability |
| FR-05 message/send + cycle de vie | Must | M2.4 | États `submitted → working → completed` observés |
| FR-06 referenceTaskIds | Must | M2.4 | Assertion sur les IDs chaînés (KYC → Scoring) |
| FR-07 états terminaux irréversibles | Must | M2.4 | `failed` ↛ `completed` (test machine à états) |
| FR-08 streaming SSE | Must | M2.5 | ≥ 2 `TaskStatusUpdateEvent` reçus |
| FR-09 webhook push | Could | M4.6 | `POST` signé, anti-rejeu (token + timestamp) |
| FR-10 AddTool schémas inférés | Must | M0.3, M1.2-M1.5 | Input invalide rejeté — le SDK v1.6.1 reporte via **IsError** (schéma inféré : requis/type/additionalProperties), **pas -32602** (ADR-0001) ; -32602 réservé au protocole (M1.7) |
| FR-11 clients MCP SDK officiel | Must | M0.4, M2.2 | Chaque agent **consommateur MCP** (KYC, Scoring, Politique — PRD §6.1) appelle ≥ 1 outil. M0 : KYC en mémoire (unit) **+ e2e stdio sur le binaire réel** (`test/e2e`, tag e2e) |
| FR-12 ressource MCP | Must | M1.6 | `ReadResource` OK ; `ResourceNotFoundError` sur URI inconnue |
| FR-13 prompt MCP | Must | M1.6 | `GetPrompt` → `GetPromptResult` |
| FR-14 NotifyProgress | Must | M1.6 | ≥ 1 `ProgressNotification` |
| FR-15 erreurs protocole vs métier | Must | M1.7 | -32602 vs `isError` ; aucun `panic` |
| FR-16 HITL deux files | Must | M3.3, M3.6 | AML élevé → `input-required` ; file de commande vide sans release |
| FR-17 timeouts/retries/CB | Should | M1.7, M4.2 | Injection de panne → échec propre |
| FR-18 idempotence | Should | M4.1 | Rejeu → un seul effet, doublon audité |
| FR-19 journal chaîné par hachage | Must | M3.1 | Réécriture → chaîne rompue détectée |
| FR-20 7 champs d'audit | Must | M3.1 | Assertion sur KYA/KYH/KYC/ts/action/résultat/version |
| FR-21 explicabilité fr/en | Must | M3.4 | 100 % des décisions ; < 500 ms ; golden fr+en |
| FR-22 explication hybride | Should | M3.4 | Refus `score<300` → motif codé déterministe (golden) |
| FR-23 audit-export | Could | M5.2b | `make audit-export` → JSONL rejouable |
| FR-24 maker-checker | Should | M3.7 | Test topologique |
| FR-25 aucun octroi ferme | Won't | M5.1 (test négatif) | Aucun chemin de code d'octroi ; e2e s'arrête à la pré-évaluation |
| FR-26 bearer à audience restreinte | Should | M2.6 | Jeton absent/invalide/`aud` étranger → 401 + audit |
| NFR-01 boucle P99 < 2 s | — | M5.2a | `make bench` |
| NFR-02 MCP P95 < 500 ms | — | M1 (sortie) | `TestToolRoundtripSmoke` (aller-retour in-memory des 5 outils, smoke ; SLO mesuré non falsifiable en mémoire, renommé T4 2026-07-08) ; l'histogramme `mcp_tool_latency` et son fournisseur (`observability.Provider`, jamais câblé en production) ont été retirés (B2, 2026-07-08) |
| NFR-03 audit < 200 ms P99 **et** ≥ 1000 écritures/s | — | M3 (sortie) | Micro-benchmark (les deux cibles) |
| NFR-04 résilience | — | M1.7, M4.2 | Injection de panne |
| NFR-05 idempotence | — | M4.1 | = FR-18 |
| NFR-06 0 egress | — | M4.3 | Réseau compose fermé + test d'hôtes |
| NFR-07 trace ininterrompue | — | M0.5, M2.7 | 1 `trace_id` par dossier dans Jaeger |
| NFR-08 portabilité 3 OS | — | M4.3 | Cross-compilation `GOOS=linux\|darwin\|windows` verte + `docker buildx` amd64+arm64 (exécution hors hôte : « non vérifiée localement ») |
| NFR-09 build reproductible | — | M4.4 | 2 builds → SHA256 égaux |
| NFR-10 déterminisme e2e | — | M5.7 | `make e2e` × 3 identiques |
| NFR-11 couverture ≥ 90 % global (unitaire ≥ 95 %) | — | M0+ (gate, plancher 90 %) ; sous-cibles rapportées en M5.7 | `make coverage` |
| NFR-12 i18n fr/en | — | M3.4 | Golden fr et en ; sélection par `lang` (US-11) |
| NFR-13 licences | — | P0.2 | `go-licenses` vert |
| NFR-14 0 vuln haute/critique | — | gate | `govulncheck` |
| G-1..G-8 | — | M2, M2, M1, M4.5, M3.4, gate, M3.1, M5.7 | Critères du PRD §3.1 |
| ADR-0001..0009 | — | M0.6, M0.6, M0.6, M1.9, M2.8, M4.5, M4.5, M3.8, M1.9 | Rédigées au format §14 du PRD (avec *Reversal condition*) |

---

## 6. Risques d'exécution (extraits du PRD §18 pertinents pour l'exécutant)

| Risque | Conduite à tenir |
|---|---|
| API réelle des SDK ≠ PRD (a2a-go jeune, ADK v2.x) | P0.1 d'abord ; épingler ; écart → ADR, jamais de déviation silencieuse |
| Non-déterminisme (LLM, ordre de goroutines, itération de maps) | Fake par défaut ; seeds fixes ; tri stable avant sérialisation ; `-race` systématique |
| Latence de composition à la limite (4 délégations ≈ 2 s vs NFR-01) | Mesurer d'abord ; paralléliser KYC ‖ Scoring seulement si le bench l'exige (Scoring ne consomme pas la sortie de KYC) |
| Sur-ingénierie | Non-buts NG-1..NG-9 du PRD ; question du pair « surdimensionné ? » ; raccourcis marqués `// ponytail:` |
| SSE : fuite de ressources | `context` + timeout + test d'annulation dès M2.5 |

---

## 7. Rapport final attendu

À la clôture de M5, produire `docs/RAPPORT-FINAL.md` : (1) DoD 11/11 avec commande + sortie constatée par case ; (2) KPI du PRD §17 mesurés ; (3) écarts au PRD (chacun adossé à une ADR) ; (4) matrice de traçabilité complète sans orphelin ; (5) reliquats et candidats (#2 identité, #3 échelle/sémantique — PRD §21).

---

## 8. Plan d'exécution de l'audit (2026-07-08)

Planifie la mise en œuvre des constats de [PRD §23 — Audit du code livré](PRD-Borealis.md#23-audit-du-code-livré-2026-07-08) (numérotation M/m = §23.2, B# = §23.3, P# = §23.4, T# = §23.5). **Ordre : correction d'abord, simplification ensuite** — supprimer du code non couvert (phase 4) *augmente* la couverture globale, donc la faire avant le resserrement du gate serait trompeur ; on corrige, on resserre le gate, puis on taille.

**Règles transverses (CLAUDE.md du dépôt) :**
- TDD : chaque correctif de bogue commence par un test qui reproduit le défaut et échoue.
- `bash scripts/check.sh` vert après **chaque** phase (plancher 90 %).
- Golden immuables : si un correctif change une sortie canonique e2e ⇒ ADR (`docs/adr/`) avant de régénérer. Les correctifs M1/m4/m12 ne touchent que des chemins d'entrées invalides ou absents du scénario nominal — vérifier néanmoins `make e2e` ×3 à chaque phase.
- Commits conventionnels en français, un par tâche.
- `-race` indisponible sur cet hôte (CGO absent) : les correctifs de course (M3, P2) se vérifient par test unitaire concurrent + revue ; noter dans le rapport de phase que la preuve `-race` reste à faire sur hôte CGO.

---

### 8.1 Phase 1 — Défauts majeurs (M1-M3)

| # | Tâche | Test d'abord | Vérification | Commit |
|---|---|---|---|---|
| 1.1 | **M1** — rejeter (toolError) `stressRate` tel que `0.09+stressRate <= 0` dans `capacity_v2.go` ; borner aussi le taux dans `capacityAt` (garde défensive). | Test appelant l'outil v2 avec `stressRate=-0.09` et `-0.19` → attend `IsError`, pas de NaN. | `check.sh` vert ; le test passe ; sortie v1 inchangée à `stressRate=0`. | `fix(mcp): borner stressRate de capacity_v2 (NaN et sous-évaluation)` |
| 1.2 | **M2** — propager `ctx` jusqu'à `Run(ctx, transport)` dans la branche stdio (`servercmd.go`) ; en profiter pour replier `Main` dans `RunMCP` (B7, même fichier, branche `--version` dupliquée inatteignable). | Test : `RunMCP` en stdio avec ctx annulé → retour sans blocage. | Test passe ; `go vet` propre. | `fix(mcp): propager le contexte de signal au transport stdio` |
| 1.3 | **M3** — `Discover` reconstruit `bySkill`/`collisions` à neuf, collisions comptées sur l'index complet ; `sync.RWMutex` partagé avec `Route`/`Delegate`/`Skills`/`Collisions`. Couvre aussi P2 par symétrie : `atomic.Bool` (ou mutex) sur `pep.available`/`level`. | Tests : (a) re-`Discover` sans un agent → skill non routé ; (b) 2ᵉ `Discover` introduisant une collision → skill retiré du routage ; (c) test concurrent `Discover`+`Route` (défendable sans `-race` : assertions d'état final). | Tests passent ; comportement mono-`Discover` inchangé (golden e2e ×3 identiques). | `fix(a2a): reconstruire l'index à chaque Discover et verrouiller le routage` |

### 8.2 Phase 2 — Mineurs sécurité et robustesse (m4-m12)

Un commit par ligne ; chaque tâche = test qui échoue → correctif.

| # | Tâche | Vérification |
|---|---|---|
| 2.1 | m12 — `amlGate` : refus si `scr.WatchlistHit`, indépendamment du score (invariant « hit ⇒ refus » porté par l'agent, plus seulement par le serveur simulé). | Test : `watchlistHit=true, score=0.1` → `OutcomeDeclined`. |
| 2.2 | m9 — refuser de démarrer en mode réseau sans `IDP_SECRET` explicite ; le secret en dur ne reste que pour le mode démo local documenté (sinon suppression). | Test : `run()` réseau sans env → erreur explicite. |
| 2.3 | m7 — `http.MaxBytesHandler` (1 MiB) sur les mux A2A (`agentcmd.go`) et MCP HTTP (`servercmd.go`). | Test : corps > 1 MiB → 4xx, pas d'allocation du corps. |
| 2.4 | m10 — `interfacesMatchBase` compare `Scheme` en plus de `Host`. | Test : card `http://` sous base `https://` → rejetée. |
| 2.5 | m8 — **résolu par ADR-0011** plutôt que par code : inspection avant action (règle « réfléchir avant d'agir ») a montré qu'aucun appelant du dépôt n'atteint les serveurs MCP HTTP (les `cmd/agent-*` servent un `ReplyExecutor` stub ; la logique métier réelle appelle les serveurs MCP en mémoire, jamais via `MCP_HTTP_ADDR`) et que la frontière réseau est déjà un périmètre de confiance explicite (`deploy/docker-compose.yml` : réseau `mesh` `internal: true`, aucun port publié pour ces services). Un bearer sans appelant réel aurait été du théâtre de sécurité non testable. Voir `docs/adr/0011-mcp-http-sans-bearer.md` (condition de réversion : dès qu'un exécuteur réel appelle ces serveurs via HTTP). | ADR relu ; commentaire de renvoi posé dans `servercmd.go`. |
| 2.6 | m4 — `deliver` : `last == nil` en fin de flux → erreur, non mémorisée par le cache d'idempotence. | Test : stream vide → erreur ; rejeu suivant ré-exécute. |
| 2.7 | m5 — disjoncteur : `context.Canceled`/`DeadlineExceeded` n'incrémentent pas les échecs. | Test : 5 annulations → circuit fermé. |
| 2.8 | m6 — deadline par défaut sur le chemin de délégation non-streaming. | Test : pair muet → erreur d'échéance, pas de blocage. |
| 2.9 | m11 — `hitl.Release` valide `approve|reject`. | Test : « aprove » → erreur, rien au journal. |

**Sortie de phase :** `check.sh` vert, `make e2e` ×3 identiques, rapport de phase.

### 8.3 Phase 3 — Gate et tests (T1-T7)

| # | Tâche | Vérification |
|---|---|---|
| 3.1 | **T2** — ajouter au gate `go test -tags e2e ./test/e2e/...` exécuté 2× avec comparaison des sorties (matérialise la règle « ×3 » du PRDPlan à coût raisonnable ; noter l'écart si on reste à 2×). | `make check` échoue si le golden dérive. |
| 3.2 | **T5** — tests des seams à 0 % : `run()` de chaque `cmd/*` (motif `TestRunVersionFlag` existant dans servercmd), `InjectTraceparent`, `SpecByName`, les deux `Cancel`. | Couverture des paquets concernés remonte (préalable à 3.3). |
| 3.3 | **T1** — appliquer « ≥ 95 % par paquet `internal/*` » dans `check.sh` ; si un paquet reste structurellement sous 95 % après 3.2 et la phase 4, **amender la directive par ADR** plutôt que tricher avec des tests de remplissage. | `check.sh` échoue sur paquet < 95 % ; liste des paquets amendés vide ou justifiée. |
| 3.4 | **T3** — mutualiser dans `internal/testfakes` : helper de connexion MCP in-memory, `fixedClock`, `goodApplicant`, `testPolicy` ; remplacer les 8 copies. | Diff net négatif ; tous les tests passent inchangés. |
| 3.5 | **T4** — ramener les `Test*P9x` SLO à un appel de fumée chacun (garder les benchmarks). | Durée du paquet `agent` < 1 s ; gate vert. |
| 3.6 | **T6/T7** — fondre `zz_repro_test.go` dans `bias_test.go` ; renommer `zz_probe_test.go` → `dataset_invariants_test.go` ; commentaire sur le golden expliquant que l'égalité byte-à-byte est l'exigence. | Gate vert. |

### 8.4 Phase 4 — Simplifications (B1-B14, ~−440 lignes, −4 dépendances)

Ordre : suppressions franches d'abord (zéro appelant), réductions ensuite. Chaque suppression : vérifier `grep` zéro appelant hors tests, supprimer code + tests dédiés, `check.sh`.

| # | Tâche | Note |
|---|---|---|
| 4.1 | B2 — supprimer `observability.Provider` + plombage `WithLatency` ; retirer `otel/sdk`, `otel/sdk/metric`, `otel/metric` de go.mod (`go mod tidy`). | Résout aussi la couverture 79,2 % du paquet (T1). Garder `tracing.go` (câblé). |
| 4.2 | B3 — **non appliqué, conservé tel quel** : ADR-0009 décide explicitement « Tenir un registre **YAML**… » comme « source unique de vérité des métadonnées de service » — un artefact de gouvernance délibéré, validé par son propre test (`LoadRegistry`+`Validate`), pas du code mort. Passer en JSON contredirait la décision acceptée sans renversement documenté ; supprimer l'artefact irait à l'encontre de l'ADR. `gopkg.in/yaml.v3` reste. | ADR-0009 relu (vérification faite, comme prévu par cette ligne) ; aucun changement de code. |
| 4.3 | B4, B5, B9, B11, B12 — supprimer `fixtures.WriteAll`, `Factsheet`/`WriteMarkdown` (remplacé par un `.md` committé si l'info doit survivre), `BiasReport.WriteCSV`, `SupportedSigAlgs()`/`ReservedSigAlgs()`, `creditdomain.DTI()`. | Un commit groupé `refactor: retirer le code sans appelant de production`. |
| 4.4 | B1 — remplacer la canonicalisation manuelle par `json.Marshal` post-`UseNumber` ; **test de non-régression d'abord** : golden des octets canoniques d'une card de référence avant/après. | Sensible (signatures JWS) — si le moindre octet diffère, abandonner et documenter pourquoi. |
| 4.5 | B6 — `decode[T any]` partagé dans `mcpcall.go`. | Tests existants inchangés. |
| 4.6 | B8 — **appliqué à `resilience.WithClock` et `webhook.WithSignerClock`/`WithVerifierClock`** (0 appelant hors leur propre test → champ `now` réglé directement depuis le paquet). **Non appliqué à `audit.WithClock`** : vérification avant action a trouvé 18 appelants dans 5 paquets externes, dont un appelant de **production** (`internal/auditexport/auditexport.go:38`) — ce n'est pas un « config nobody sets », c'est une couture réellement adoptée ; la boilerplate `Option`/variadique est un idiome Go légitime pour une valeur externally-settable, pas du sur-dimensionnement ici. Retirer l'option forcerait une réécriture mécanique de 18 sites pour un gain cosmétique. | Tests resilience/webhook adaptés mécaniquement ; audit.WithClock conservé tel quel (décision documentée, pas un oubli). |
| 4.7 | B13, B14 — `cmp.Or` pour `env()` ; `math.Round` pour `fixtures.round2`. | ⚠ `round2` : vérifier que les valeurs générées ne bougent pas (fixtures seedées → golden e2e ×3). Si elles bougent : ADR ou abandon. |
| 4.8 | B10 — trancher `Skills()`/`Collisions()` **après** 1.3 : s'ils restent test-only une fois le mutex posé, les déclasser en méthodes non exportées utilisées par les tests du paquet. | — |

### 8.5 Phase 5 — Plafonds de performance (P1, P3-P5)

Pas de correctif immédiat justifié à l'échelle du démonstrateur ; documenter les plafonds pour ne pas les redécouvrir :

- P1 : commentaire `// ponytail:` sur le journal WORM en mémoire (plafond : croissance non bornée en service long ; chemin d'évolution : sink JSONL derrière la même interface). Si les binaires sont censés tourner en continu (compose), envisager le sink dès maintenant — décision à prendre au moment de la phase.
- P3/P4 : compléter les commentaires `ponytail:` existants (croissance du cache d'idempotence, purge O(n) des nonces).
- P5 : hisser `cardPayload` hors de la boucle de signatures (3 lignes, sans risque) ou laisser avec note.

---

### 8.6 Critère de sortie global

1. ✅ `bash scripts/check.sh` **et** `scripts/check.ps1` verts avec le gate resserré (95 % par paquet, exemptions documentées : `internal/auditexport` 90 %, `internal/testfakes` 87,5 % — branches défensives de connexion SDK non déclenchables).
2. ✅ `go test -tags e2e ./test/e2e/...` ×2 (réduit de ×3, écart documenté T2) → sorties canoniques identiques à chaque phase ; golden intouchés.
3. ⚠️ `go mod tidy` : **2 dépendances retirées** (`otel/sdk/metric` direct, `gopkg.in/yaml.v3` reste — B3 non appliqué, ADR-0009 exige explicitement le YAML) sur les 4 visées initialement ; écart justifié (4.2).
4. ✅ Aucun constat majeur (M1-M3) ni mineur (m4-m11) ouvert. m8 résolu par ADR-0011 (frontière réseau, pas de code).
5. ✅ Rapport de phase : 27/27 tâches complétées, un commit par tâche (voir `git log`). Écarts documentés inline : B3 (§4.2), audit.WithClock (§4.6), P5 (§Phase 5).

**Reporté sciemment :** preuve `-race` (hôte CGO requis — à exécuter dès qu'un hôte gcc/Docker est disponible, cf. docs/verification-p01.md) ; sink d'audit persistant si les processus restent éphémères (documenté en commentaire, P1).

**Couverture totale finale : 96,2 %** (départ : 90,3 %).

---

*Fin du plan v1.0. Exécution : commencer par P0.1 (levée des ⚠) ; ne clore aucune phase sans son gate et sa vérification adverse ; le PRD-Borealis.md v2.0 reste la source de vérité.*
