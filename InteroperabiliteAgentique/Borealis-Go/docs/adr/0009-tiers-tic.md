# ADR-0009 — Serveurs MCP traités comme tiers TIC (registre de services)

- **Statut :** accepté
- **Date :** 2026-07-07
- **Phase :** M1.9

## Contexte

Dans le domaine financier, un service externe consommé par le système entre dans
le périmètre des **tiers TIC** (fournisseurs de services numériques ; posture
DORA / Loi 25). Les serveurs MCP de Borealis (Identity Hub, Bureau, Capacity,
Policy) sont, du point de vue des agents, de tels fournisseurs : ils exposent un
contrat de service (E/S, SLO, taxonomie d'erreurs). Sans registre, ces contrats
restent implicites.

## Décision

Tenir un **registre YAML** (`internal/mcpserver/registry.yaml`, embarqué via
`go:embed`) décrivant chaque serveur MCP comme tiers TIC : `name`, `tier`,
`contractVersion`, `sloP95Ms`, `tools`, `errors` (énum). Un loader
(`LoadRegistry`/`DefaultRegistry`) le parse et le **valide** (nom, version de
contrat, ≥ 1 outil, énum d'erreurs non vide). Le registre est la source unique
de vérité des métadonnées de service.

## Conséquences

- Contrats de service explicites et versionnés → base pour la gouvernance
  d'évolution (compat BACKWARD, ADR à venir M4.5) et la surveillance des SLO.
- Les valeurs de SLO sont des **hypothèses à calibrer** (PRD §17) — jamais
  présentées comme mesurées.
- Léger couplage : ajouter un serveur MCP impose de mettre à jour le registre
  (validé par le test, donc non oublié).

## Reversal condition

Un registre de services centralisé externe (catalogue d'entreprise) rendrait ce
fichier redondant — à migrer alors vers la source d'autorité.

## Alternatives écartées

- **Métadonnées codées en dur dans chaque binaire** : dispersées, non
  inspectables d'un coup, difficiles à auditer.
- **Pas de registre** : contrats implicites, contraire à l'invariant *contrat*.
