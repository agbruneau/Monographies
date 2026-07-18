# ADR-0012 — Encodage injectif de l'entry_hash (préfixe de longueur)

- **Statut :** accepté
- **Date :** 2026-07-15
- **Phase :** maintenance post-M5 (audit de code du 2026-07-15)

## Contexte

Le PRD §9.4 définit `entry_hash = sha256(seq || ts || kya || owner_kyh || kyc ||
action || input_hash || result || prev_hash)` sans spécifier l'encodage de la
concaténation. L'implantation (`internal/audit.hashEntry`) séparait les champs
d'un simple octet `\x00`. Cet encodage n'est **pas injectif** : si un champ
contient un octet NUL (possible via un JSON portant la séquence d'échappement
Unicode du caractère nul, reprise dans un `Result`), deux Records distincts
obtenus en déplaçant la frontière entre champs adjacents produisent le même
hash — une telle altération d'un export passerait
`VerifyEntries`/`VerifyExport`, à rebours de la garantie tamper-evident
(FR-19, G-7). `webhook.mac` corrigeait déjà cette ambiguïté par préfixe de
longueur ; le journal WORM, dont c'est la garantie centrale, ne l'appliquait
pas.

## Décision

Chaque champ variable de `hashEntry` est préfixé de sa longueur
(`\x00<len>:<champ>`). Les champs hachés et leur ordre restent ceux de la
formule PRD §9.4 (restreinte aux 7 champs du journal, restriction déjà
documentée) ; seul l'encodage change.

## Conséquences

- Positif : l'encodage devient injectif — toute altération de frontière de
  champ change le hash ; cohérence avec `webhook.mac`.
- Négatif : les hash produits avant ce changement ne sont pas recomputables
  avec le nouveau code. Aucun export persistant n'existe (journal en mémoire,
  processus courts, aucun golden ne fige de hash) : pas de migration.
- Test de non-régression : `TestHashEntryInjectiveOnNULBoundary`.

## Reversal condition

Si un vérificateur externe tiers implémente la formule PRD §9.4 au pied de la
lettre (séparateur nu) et doit re-vérifier des exports produits ici, revenir à
un encodage commun documenté dans le contrat d'export (ou publier l'encodage
préfixé comme partie du contrat FR-23).

## Alternatives écartées

- **Statu quo (séparateur `\x00` nu)** : laisse une classe d'altérations
  indétectables — contraire à l'objet même du journal (G-7).
- **JCS/JSON canonique de l'entrée** : plus lourd, dépendance de plus, sans
  gain d'injectivité par rapport au préfixe de longueur.
- **Échappement des NUL dans les champs** : déplace le problème (l'échappement
  lui-même doit être injectif) et touche les données au lieu de l'encodage.
