# P0.1 — Vérification adverse des ⚠ du PRD (levée des incertitudes SDK)

**Date de vérification :** 2026-07-07. **Sources :** dépôts officiels GitHub (pages releases/README), consultées ce jour.
**Verdict global :** aucun écart majeur → feu vert pour M0. Un seul correctif structurel : `a2a-go` est en **v2.x**, pas v1.x.

| Énoncé PRD (⚠) | Attendu (PRD) | Constaté | Action |
|---|---|---|---|
| `modelcontextprotocol/go-sdk` versions / licence | v1.0.0 (gel) puis v1.6.x ; MIT + Apache 2.0 | Actif ; latest **v1.6.1** (v1.7.0-pre) ; **Apache 2.0 (contributions neuves) + MIT (code existant)** ; v1.0.0 a existé (compat spec 2025-06-18) | Épingler **v1.6.1**. Conforme au PRD. → ADR-0001 |
| OAuth **client** du SDK MCP | Non stabilisé côté client | **Stabilisé en v1.5.0** (build tag `mcp_go_client_oauth` retiré ; garanties de compat standard) | Écart *favorable* : FR-26 reste la couture applicative, mais le SDK la porte nativement. Ne pas élargir le périmètre. → noter en ADR-0004/0009 |
| `a2a-go` : dépôt canonique + version | `a2aproject` ; v1.x | Dépôt **`a2aproject/a2a-go`** actif ; latest **v2.3.1** ; Go ≥ 1.24.4 ; Apache 2.0 ; module `/v2` | **Écart de version** : épingler **v2.3.1** (chemin `github.com/a2aproject/a2a-go/v2`). → ADR-0005 |
| MCP révision 2026-07-28 (RC) | RC 21 mai 2026 ; Roots/Sampling/Logging *Deprecated* | Non re-confirmé sur le web ce jour (à revalider avant M1.9) | Sans impact structurel (le PRD ne dépend pas de ces primitives). → à noter en ADR-0004 |
| `google/adk-go` v2.x | Optionnel, jeune | **v2.0.0** (30 juin 2026), code-first | Optionnel, derrière `Reasoner`. Sans impact structurel |

**Résolution effective des modules (probe `go get`, 2026-07-07) :**
`github.com/modelcontextprotocol/go-sdk v1.6.1` et `github.com/a2aproject/a2a-go/v2 v2.3.1` se téléchargent via `proxy.golang.org`. Le build est viable.

**Contraintes d'hôte** (hors PRD ; à lever sur un hôte CGO/Linux + Docker) :
- **Docker absent** → `docker compose` (DoD 1), build reproductible SHA256 (DoD 10), M4.3/M4.4 non vérifiables sur cet hôte ; cross-compilation `GOOS=… go build` reste possible.
- **CGO/gcc absent** → `-race` non exécutable sur cet hôte (DoD 9 partiel) ; s'exécute sur un hôte CGO/Linux ou en conteneur.
