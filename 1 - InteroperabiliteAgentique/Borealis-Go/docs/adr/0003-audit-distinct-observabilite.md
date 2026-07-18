# ADR-0003 — Journal d'audit WORM distinct de l'observabilité

- **Statut :** accepté
- **Date :** 2026-07-07
- **Phase :** M0.5 / M0.7 / M0.6

## Contexte

Deux besoins de traçabilité coexistent : (1) l'**audit** réglementaire — qui a
fait quoi, vérifiable, inaltérable (WORM) ; (2) l'**observabilité** SRE — traces
et logs pour le diagnostic. Les confondre exposerait l'audit à la rétention, à
l'échantillonnage et à la mutabilité des puits d'observabilité (invariant 7 de
l'Annexe B).

## Décision

Deux puits **séparés** :

- `internal/audit` — journal **append-only** (M0), portant les 7 champs
  (KYA/KYH/KYC/ts/action/résultat/version). Horloge injectable (déterminisme).
  La chaîne de hachage vérifiable (`entry_hash`/`prev_hash`, PRD §9.4) et son
  vérificateur arrivent en M3.1.
- `internal/observability` — traçage OTel (exportateur injectable) + logger
  slog. Jamais utilisé comme source de vérité d'audit.

Aucune écriture d'audit ne transite par le puits OTel, et réciproquement.

## Conséquences

- L'audit reste vérifiable indépendamment de l'infrastructure d'observabilité.
- Léger coût de duplication (deux chemins d'écriture) — assumé.
- La corrélation reste possible via le `trace_id`/`task id` **copié** dans
  l'entrée d'audit (M2.7), sans fusionner les puits.

## Reversal condition

Une exigence réglementaire imposant un format d'audit qui coïnciderait
exactement avec le puits d'observabilité (improbable) — à réévaluer alors.

## Alternatives écartées

- **Un seul journal pour tout** : viole l'invariant 7 ; l'audit hériterait de la
  mutabilité/rétention des logs — rejeté.
