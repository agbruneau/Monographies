# ADR-0006 — Crypto-agilité : suite d'algorithmes portée par le contrat

- **Statut :** accepté
- **Date :** 2026-07-07
- **Phase :** M4.5 (placeholder post-quantique concret en M5.6)

## Contexte

Les signatures d'Agent Card (FR-03) utilisent aujourd'hui **ES256** (ECDSA
P-256), figé dans `pkg/a2a`. L'horizon post-quantique (PRD ch.7, veille §12)
impose de pouvoir **migrer d'algorithme** sans réécrire les consommateurs ni
rompre les cards déjà émises. Coder l'algorithme en dur rendrait cette migration
coûteuse et non testable ; l'inventer complètement (crypto PQ réelle) sort du
périmètre du démonstrateur (NG-9).

## Décision

**Porter l'identifiant d'algorithme dans le contrat**, pas dans le code appelant.
Deux applications :

1. **Signature de card** — l'en-tête protégé JWS transporte déjà `alg` (vérifié :
   `checkProtectedHeader` rejette `none`/`crit`, exige `alg == ES256`). C'est la
   **couture de crypto-agilité** : un vérifieur lit l'algorithme de la card, il ne
   le présume pas. Migrer vers un algorithme PQ (ex. ML-DSA) reviendra à étendre
   l'ensemble des `alg` acceptés et à publier des cards signées avec le nouveau —
   sans toucher les consommateurs qui lisent déjà `alg`.
2. **Évolution de contrat d'outil** — la compatibilité **BACKWARD** vérifiée par
   `pkg/mcpcontract` (M1.8) et démontrée par la substitution Capacity v1→v2
   (M4.5) est le même principe appliqué aux schémas : un champ **ajouté** (ex. une
   future `cryptoSuite`) n'est jamais requis, donc n'invalide pas les producteurs
   antérieurs.

Le **placeholder post-quantique concret** (une suite d'algorithmes encodée, sans
crypto PQ réelle) est livré en **M5.6** ; cette ADR fixe la décision d'architecture
qui l'autorise.

## Conséquences

- Migration d'algorithme = extension d'un ensemble accepté + republication de
  cards, jamais une modification des vérifieurs → surface de changement minimale.
- La négociation reste explicite et auditable (l'`alg` d'une card est inspectable) ;
  pas de rétrogradation silencieuse (l'agilité n'affaiblit pas : `none` reste rejeté).
- Coût : un ensemble d'algorithmes acceptés à maintenir ; un test de contrat le
  garde honnête (le durcissement de `alg` est déjà couvert par `TestVerifyRejectsAlgConfusion`).

## Reversal condition

Si une seule primitive devait être imposée par norme (interdiction réglementaire
de tout algorithme classique), l'agilité deviendrait un coût mort : figer alors
l'unique algorithme et retirer la négociation.

## Alternatives écartées

- **Algorithme codé en dur dans les vérifieurs** : migration = modification de
  chaque consommateur ; rupture des cards existantes.
- **Crypto PQ réelle dès le candidat #1** : hors périmètre (NG-9), fausse maturité ;
  reportée (placeholder M5.6, module complet au candidat #2 — [[0008-report-module-identite]]).
