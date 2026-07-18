# ADR-0007 — Stratégie de test

- **Statut :** accepté
- **Date :** 2026-07-07
- **Phase :** M4.5

## Contexte

Le démonstrateur touche des invariants sensibles (jamais d'octroi ferme, journal
non répudiable, séparation maker-checker, criblage AML) et doit rester **vérifiable
localement**, sans CI distante (calibre FibGo). Une couverture chiffrée seule ne
prouve rien : un test peut couvrir une ligne sans en verrouiller le comportement.
Il faut une stratégie explicite, sinon chaque phase invente la sienne.

## Décision

**Pyramide + déterminisme + gate local + vérification adverse**, dans cet ordre de
défense :

1. **Pyramide (PRD §12.1)** : ~70 % unitaire (`mcp.NewInMemoryTransports()`,
   logique pure), ~20 % intégration (fakes LLM + serveurs MCP en mémoire), ~10 %
   e2e (binaire réel, tag `e2e`). Répartition mesurée au rapport M5.7.
2. **Déterminisme (G-8, NFR-10)** : `Reasoner` fake par défaut (jamais de LLM réel
   dans le gate) ; fixtures **seedées** ; horloges et générateurs d'ID **injectables**
   (`WithClock`) ; tri stable avant sérialisation ; golden octet-à-octet.
3. **Gate local (`scripts/check.{sh,ps1}`)** : build, vet, test (`-race` si CGO),
   `govulncheck`, `golangci-lint`, **plancher de couverture 90 %** (directive
   utilisateur, > NFR-11), bench sans régression. Aucune phase ne se clôt gate rouge.
4. **Vérification adverse (protocole §2.4 du plan)** : à chaque clôture de phase,
   ≥ 2 agents réfutateurs à contexte frais tentent de **casser** les critères de
   sortie ; les constats confirmés sont corrigés avant clôture.
5. **Test de mutation ciblé** : pour les tests qui gardent un invariant critique
   (ancrage d'audit, topologie maker-checker, substitution BACKWARD), on **neutralise
   la défense** et on vérifie que le test **échoue** — preuve qu'il n'est pas
   tautologique. Appliqué en M3 (2 tests MAJEURS re-prouvés).

**Substitution de contrat** (M4.5) : la compatibilité BACKWARD (`pkg/mcpcontract`)
est vérifiée *par test* (Capacity v1→v2), avec un **oracle négatif** (un champ
ajouté rendu requis doit rompre) pour que l'assertion positive ne soit pas vide.

## Conséquences

- Un défaut d'invariant échoue **localement**, tôt, avec une commande reproductible.
- La couverture chiffrée est un plancher, pas la preuve : la preuve est la mutation
  et l'adversaire.
- Coût : chaque phase paie sa vérification adverse (jetons, temps) — assumé pour le
  calibre visé.

## Reversal condition

Si une CI distante de confiance était introduite, une partie du gate local pourrait
y migrer (garder le déterminisme et la mutation ; déléguer vuln/lint/couverture).

## Alternatives écartées

- **Couverture chiffrée comme seul critère** : verrouille des lignes, pas des
  comportements ; laisse passer des tests tautologiques (constaté et corrigé en M3).
- **e2e lourd comme socle** : lent, non déterministe sans fakes ; inverse la pyramide.
