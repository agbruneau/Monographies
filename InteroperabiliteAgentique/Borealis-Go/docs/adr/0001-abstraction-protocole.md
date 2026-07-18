# ADR-0001 — Dépendre des SDK officiels, isoler le domaine derrière des adaptateurs

- **Statut :** accepté
- **Date :** 2026-07-07
- **Phase :** M0.3 / M0.6

## Contexte

Le démonstrateur doit parler MCP (et A2A en M2). Deux voies : réimplémenter le
protocole, ou dépendre des SDK officiels. La monographie insiste sur l'invariant
*découplage, contrat, évolution* : le code métier ne doit pas être couplé aux
détails d'un SDK. P0.1 a épinglé `modelcontextprotocol/go-sdk` v1.6.1
(dual-licence Apache 2.0 + MIT) ; les modules se résolvent via GOPROXY.

## Décision

Dépendre du **go-sdk officiel v1.6.1** pour MCP. Le domaine pur
(`internal/creditdomain`) et l'orchestration (`internal/agent`) ignorent le SDK ;
les seuls points de contact sont des adaptateurs minces : `internal/mcpserver`
(construction de serveurs + outils) et `agent.VerifyIdentity` (appel d'outil +
`parseIdentityResult`). Le lancement des serveurs est factorisé dans
`internal/servercmd` (main() d'une ligne, logique testable).

## Conséquences

- **Écart au PRD acté (FR-10/FR-15) :** le SDK valide l'entrée contre le schéma
  inféré, mais reporte toute violation (champ requis manquant, type invalide,
  propriété additionnelle) en **`CallToolResult.IsError = true`**, et **non** en
  erreur protocole `-32602`. Ce dernier reste réservé aux erreurs de protocole
  (méthode inconnue, params malformés). La distinction protocole/métier de FR-15
  sera traitée en M1.7 en tenant compte de ce comportement ; la matrice est
  ajustée en conséquence.
- Le schéma d'entrée est strict par défaut (`additionalProperties: false`,
  champs requis) — bon pour la robustesse.

## Reversal condition

Un changement d'API majeur du SDK (rupture de compat au-delà de v1.x) ou une
divergence du SDK avec la spec MCP qui empêcherait un cas d'usage du PRD.

## Alternatives écartées

- **Réimplémenter le protocole** : coûteux, source de dérive avec la spec,
  contraire à *ne pas réinventer*.
- **Envelopper le SDK derrière notre propre interface `Protocol`** : abstraction
  spéculative à une seule implémentation (anti-pattern), écartée tant qu'un
  second protocole n'impose pas la couture.
