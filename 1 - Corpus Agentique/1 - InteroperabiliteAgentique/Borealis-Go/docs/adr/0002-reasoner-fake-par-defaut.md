# ADR-0002 — Le moteur de raisonnement derrière l'interface Reasoner ; fake déterministe par défaut

- **Statut :** accepté
- **Date :** 2026-07-07
- **Phase :** M0.2 / M0.6

## Contexte

La décision assistée fait intervenir un « raisonnement » (règles ou LLM). Un LLM
réel est non déterministe (G-8) et rendrait le gate irreproductible (NFR-10). Le
PRD §7.2 prescrit une couture d'abstraction.

## Décision

Définir l'interface **`agent.Reasoner`** (`Reason(ctx, ReasonRequest) →
ReasonResult`). Le **défaut d'exécution et de test est un fake déterministe**
(`internal/testfakes.StubReasoner`, piloté par scénario : approved / denied /
escalate ; scénario inconnu → escalade prudente). Aucun LLM réel n'entre dans le
gate. Un moteur réel (ex. `adk-go`) pourra être branché derrière la même
interface sans toucher au domaine ni à l'orchestrateur.

## Conséquences

- Reproductibilité : `make e2e` × 3 → sorties canoniques identiques (le fake ne
  varie pas).
- Les recommandations du fake sont cohérentes avec l'issue métier (mapping
  outcome → scénario dans l'orchestrateur), suffisant pour le walking skeleton.
- L'invariant *jamais d'octroi ferme* (FR-25) est porté par le domaine, pas par
  le reasoner : aucune recommandation ne vaut « octroyé ».

## Reversal condition

Un besoin avéré d'évaluer un LLM réel dans un test d'intégration ciblé (hors
gate) justifierait un mode « reasoner réel » optionnel, toujours derrière
l'interface et jamais dans le gate déterministe.

## Alternatives écartées

- **Appeler un LLM directement dans l'orchestrateur** : couplage fort, non
  déterministe, non testable — rejeté.
- **Pas d'abstraction (règles en dur)** : empêcherait de démontrer la couture
  agentique du PRD.
