# ADR-0010 — Écart : mapping ArchiMate dérivé du PRD (sans Annexe B)

- **Statut :** accepté
- **Date :** 2026-07-07
- **Phase :** M5.3
- **Nature :** ADR d'**écart** (déviation documentée, non une décision de conception)

## Contexte

Le PRD §14 / KPI §17 attend un `docs/ARCHITECTURE.md` mappant les **vues
ArchiMate A-L** de l'Annexe B aux paquets Go, avec les **stéréotypes verbatim du
registre §6.6** (ch.6 §6.1.9). Or l'exécutant n'a, en cas nominal, **pas accès à
l'Annexe B ni au registre §6.6** : ces artefacts vivent dans le dépôt d'écriture
*InteroperabiliteAgentique*, hors du dépôt de code.

## Décision (repli)

Conformément à la règle de repli du plan (M5.3), `docs/ARCHITECTURE.md` :

1. **dérive le schéma maître des diagrammes du PRD §6.3-6.5** (C4, séquence, flux)
   au lieu des vues A-L de l'Annexe B ;
2. remplace la table « **vue A-L → paquet Go** » par une table
   « **composant du PRD §6 → paquet Go** » ;
3. emploie les stéréotypes **ArchiMate 3.2/4 standard de la couche application**
   (`«application component»`, `«application service»`, `«application interface»`,
   `«data object»`) faute du registre §6.6 verbatim.

## Conséquences

- Le mapping reste **fidèle au PRD** (source accessible) et vérifiable contre le
  code ; il ne prétend pas reproduire les vues A-L qu'il n'a pas lues.
- Un lecteur disposant de l'Annexe B pourra **rapprocher** la table §6 → paquet
  des vues A-L sans contradiction (les composants §6 sont un sur-ensemble nommé).
- Écart **assumé et signalé** (règle « jamais de déviation silencieuse ») ; il
  n'affecte aucun comportement de code.

## Reversal condition

Si l'Annexe B et le registre §6.6 deviennent accessibles à l'exécutant, remplacer
la table « composant §6 → paquet » par la table « vue A-L → paquet » attendue et
substituer les stéréotypes verbatim — cet écart devient alors sans objet.

## Alternatives écartées

- **Inventer des vues A-L** sans les avoir lues : fabriquerait une fausse
  conformité, contraire à l'honnêteté du rapport.
- **Omettre ARCHITECTURE.md** : laisserait DoD 11 (README pont + mapping)
  incomplet.
