# ADR-0008 — Report du module d'identité complet au candidat #2

- **Statut :** accepté
- **Date :** 2026-07-07
- **Phase :** M3.8

## Contexte

Le PRD §19 décrit un module d'identité agentique de qualité production :
identités non humaines (NHI) provisionnées **just-in-time**, clés de signature
protégées par **HSM**, journal d'audit scellé sur support **WORM matériel**, et
propagation d'autorisation **OAuth multi-hop** (délégation de jeton d'agent en
agent avec réduction de portée). Le démonstrateur (candidat #1) doit prouver la
posture de sécurité STRIDE (§11.3) sans porter ce module au complet : le budget
et la portée du candidat #1 visent la démonstration d'interopérabilité, non le
durcissement infrastructurel.

Le risque : simuler à demi ces mécanismes (un « faux » HSM logiciel, un WORM en
fichier local présenté comme scellé) donnerait une **fausse assurance** de
sécurité — précisément l'anti-patron que la règle « ne jamais présenter un
*preview* comme GA » proscrit.

## Décision

**Reporter le module d'identité complet au candidat #2.** Le candidat #1
implémente des **substituts explicitement étiquetés** couvrant la *forme* du
contrat de sécurité, jamais sa robustesse matérielle :

- **NHI / identité de service** → `idpmock` : bearer JWT HS256 à **audience
  restreinte** (claim `aud`, RFC 8707). Pas de provisionnement JIT, pas de
  rotation, pas de protection anti-rejeu — documenté dans le paquet.
- **Clés de signature** → clés ECDSA P-256 **en mémoire** (`pkg/a2a`), signature
  ES256/JCS des cards. Pas de HSM ; les clés ne quittent pas le processus.
- **WORM** → journal à **chaîne de hachage logicielle** (`internal/audit`),
  vérifiable (`Verify`, `VerifyEntries`). Tamper-**évident**, non tamper-**proof** :
  aucun scellement matériel ni stockage immuable.
- **OAuth multi-hop** → un seul saut d'audience vérifié ; pas de délégation
  chaînée ni de réduction de portée transitive.

La suite STRIDE consolidée (`internal/security/stride_test.go`) prouve que ces
substituts défendent bien contre spoofing, tampering et elevation of privilege
**au niveau du contrat** ; leur montée en robustesse relève du candidat #2.

## Conséquences

- Posture de sécurité **honnête** : chaque substitut est étiqueté, aucune
  assurance matérielle n'est revendiquée. Les menaces §11.3 sont couvertes par
  des tests, pas par des allégations.
- Le candidat #2 hérite d'un périmètre net (NHI JIT, HSM, WORM matériel, OAuth
  multi-hop) et de contrats déjà éprouvés à substituer sans réécrire les appelants
  (l'abstraction protocolaire d'ADR-0001 isole ces changements).
- Limite assumée : le démonstrateur ne résiste pas à un attaquant disposant d'un
  accès au processus (extraction de clés en mémoire) — hors modèle de menace du
  candidat #1.

## Reversal condition

Si le candidat #1 devait être exposé hors d'un environnement de démonstration
cloisonné (données non synthétiques, accès tiers), ce report serait invalidé :
le module d'identité matériel deviendrait un prérequis de mise en service, pas un
report.

## Alternatives écartées

- **Implémenter un HSM/WORM logiciel présenté comme équivalent** : fausse
  assurance de sécurité, contraire à la règle *preview ≠ GA*.
- **Omettre toute défense d'identité dans le candidat #1** : rendrait la suite
  STRIDE (§11.3) invérifiable et l'invariant *contrat* non démontré.
