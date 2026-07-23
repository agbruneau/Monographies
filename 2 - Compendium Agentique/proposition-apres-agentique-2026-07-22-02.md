# Proposition de passe — après-agentique, itération 02

| Champ | Valeur |
|---|---|
| Date | 22 juillet 2026 |
| Itération | 02 — relève v0.10 nᵒ 4 instruite à sa source primaire |
| Objet | Fiche d'instruction : préimpression A. Metere, taxonomie de la divergence action/journal |
| Sièges | Vol. IV `TOC.md` v0.10 — ch. 54 (Livre IX, sémantique d'effet, chapitre d'accueil principal) et ch. 43 (Livre VII, observabilité) ; journal v0.10, relève 4 |
| Statut | **Proposition** — rien n'est intégré, rien n'est poussé ; l'intégration est une décision d'auteur |
| `TOC.md` touché | **Non par cette passe** — une entrée de journal prête à intégrer est fournie en annexe, si l'auteur consomme la fiche |

## L'incrément

**La relève v0.10 nᵒ 4 est instruite à sa source primaire : la préimpression existe, est datée et dit ce que le journal v0.10 lui prête — avec deux précisions de libellé et un recoupement nouveau.** L'extraction effective ci-dessous satisfait la condition d'entrée posée par la gouvernance (« une relève n'entre au socle qu'après extraction effective de cette source »).

### Fiche d'instruction

**Source.** A. Metere, *Architectural Obsolescence of Unhardened Agentic-AI Runtimes*, arXiv 2605.01740**v1**, soumise le 3 mai 2026 (06:38 UTC), catégories cs.CR / cs.AI / cs.MA. Consultée le 22 juillet 2026 (page de résumé et texte intégral HTML). Aucune version ultérieure, aucune mention de publication ni de révision par les pairs.

**Énoncé instruit (niveau [B] — document public daté, extraction citée ; élévation depuis [C] au titre du risque 11).** La préimpression pose comme propriété de sûreté d'un runtime agentique la correspondance biunivoque entre mutations réelles et enregistrements d'audit — verbatim : « *every real mutation has exactly one successful audit record, and vice versa* » — et fait de la **détection de l'écart** (non de la richesse de la trace) le critère porteur. Elle énumère quatre formes de divergence :

| Forme | Terme d'origine | Définition (extraite le 22 juillet 2026) |
|---|---|---|
| F1 | *gate bypass* | mutation sans enregistrement correspondant — l'agent a changé l'état sans passer la garde d'admission |
| F2 | *audit forgery* | enregistrement sans mutation correspondante — le journal fabrique la preuve d'un travail jamais advenu |
| F3 | *approved-but-failed silent* | mutation partiellement advenue avec enregistrement marqué en échec — le résidu se projette comme F1 |
| F4 | *wrong target* | garde approuvée pour (capacité, cible A), mutation advenue sur (capacité, cible B) |

: Taxonomie F1-F4 de la préimpression Metere, termes d'origine et définitions extraites.

La détection exige, selon la pièce, trois primitives : un journal d'audit consignant chaque paire (capacité, cible), une source de delta d'état, et un vérificateur comparant les deux multiensembles.

**Précisions de libellé sur le journal v0.10 (à porter à l'intégration).** (1) La paraphrase « échec silencieux de l'hôte » (F3) est proche mais non verbatim : le terme d'origine est *approved-but-failed silent* — échec silencieux de l'exécution approuvée, l'« hôte » n'étant pas dans la définition. (2) Les quatre noms d'origine (*gate bypass*, *audit forgery*, *approved-but-failed silent*, *wrong target*) doivent accompagner la taxonomie à sa première occurrence, conformément à la règle de terminologie du dépôt.

**Réserves (reconduites et confirmées).** La pièce est non révisée par les pairs ; son auteur propose une implémentation concurrente (« enclawed-oss ») de l'objet qu'il mesure, et le résumé revendique pour elle des scores parfaits — **les résultats chiffrés comparatifs demeurent exclus de la relève**, comme le journal v0.10 le prescrivait : une mesure produite par l'auteur de l'implémentation concurrente se qualifie avant de se citer. Le protocole de mesure (1 600 cas de base, extension à 80 000, contenu adversarial généré par dix modèles) est rapporté ici comme *dispositif déclaré*, non comme résultat.

**Recoupement nouveau (constat de convergence, non une réplication).** L'objet mesuré, « OpenClaw », est décrit par la pièce comme « *the single-user agentic-AI gateway by Steinberger* » — le même auteur que la « loi de Steinberger » de la relève 8, et un candidat naturel au « runtime largement déployé » de la chronologie de la relève 5 (sources ouvertes, non extraites à ce jour). Trois relèves v0.10 (4, 5, 8) convergent donc vers le même écosystème d'exécution ; convergence à noter, **jamais** à compter comme sources indépendantes — les pièces partagent un objet, pas une méthode.

### Tri épistémique

- Existence, date, contenu et taxonomie de la préimpression : **[B]** (document public daté, extraction citée, non révisé par les pairs).
- La taxonomie comme grille du ch. 54 : **candidate** — décision d'auteur, comme le journal v0.10 le posait.
- Les scores d'enclawed-oss : **exclus** (auto-mesure en conflit d'intérêt).
- Rapprochement d'auteur, déclaré tel : le « vérificateur comparant deux multiensembles » est, dans le vocabulaire du ch. 54, une **réconciliation** — la taxonomie F1-F4 atterrit donc au cœur de la thèse du chapitre (« tracer l'effet, pas seulement l'appel »), et F2 nomme le versant que le ch. 43 appelle le chaînon manquant (trace ↔ chaîne de mandat).

## Atterrissages proposés

- **Ch. 54 (principal).** La taxonomie F1-F4 comme grille candidate de la section « réconciliation » ; les termes d'origine en première occurrence ; la réserve de conflit d'intérêt attachée à toute citation.
- **Ch. 43.** F2 (*audit forgery*) comme cas limite de l'observabilité : une trace riche peut être fausse — l'argument renforce la thèse du chapitre sans nouvelle structure.
- **Ch. 20.** Sans objet direct par cette fiche (la relève 4 n'y atterrit pas) ; le recoupement OpenClaw intéresse en revanche l'instruction future de la relève 5 (ch. 52, 21, 20).

## Annexe — entrée de journal prête à intégrer (si l'auteur consomme la fiche seule)

> **Actualisation v0.11 — instruction de la relève 4 (date d'intégration à poser).** La préimpression arXiv 2605.01740v1 (Metere, 3 mai 2026) a été extraite le 22 juillet 2026 : la taxonomie F1-F4 et la propriété de correspondance biunivoque sont confirmées verbatim ; deux précisions de libellé (F3 = *approved-but-failed silent* ; noms d'origine requis) ; résultats chiffrés toujours exclus (conflit d'intérêt confirmé — implémentation concurrente du même auteur). La relève 4 passe de « à instruire » à « instruite, [B] » ; son entrée au socle du ch. 54 demeure une décision de rédaction. Fiche : `proposition-apres-agentique-2026-07-22-02.md`.

Note de réalignement `Conspectus.md` : la phrase « huit relèves datées du 21 juillet 2026, toutes **[C] et à instruire à la source primaire** » (section « L'angle mort déclaré ») devrait alors se nuancer — sept restent [C] à instruire, une est instruite [B]. **Re-mesure faite : le cardinal « huit relèves » ne change pas** ; seul leur statut se ventile.

## Vérification adverse (exécutée avant livraison)

- **Réfutation tentée** : la description du journal v0.10 confrontée au document réel — existence, auteur, date, version, taxonomie : conformes ; une paraphrase non verbatim relevée (F3) et consignée plutôt que lissée.
- **Conflit d'intérêt** : confirmé à la source (enclawed-oss, scores parfaits auto-déclarés) — l'exclusion des chiffres est maintenue.
- **Cardinaux** : « huit relèves » re-mesuré, inchangé ; aucun décompte du dépôt modifié par cette livraison.
- **Renvois** : ch. 20, 21, 43, 52, 54 cités — tous dans 1-57.
- **Zones gelées** : aucune touchée ; Livre X terminal ; aucune renumérotation citée.
- **Décision 8** : la thèse du ch. 53 et celle du ch. 20 ne sont pas réécrites ; tout rapprochement est marqué construction d'auteur.

## Prochaine itération envisagée

Instruire la relève v0.10 nᵒ 2 (chaîne d'approbation d'outils) à ses sources primaires écrites : *Anatomy of a harness* (5 juin 2026) et *Announcing Mastra Harness* (18 juin 2026, `@mastra/core@1.45.0`) — pièces d'éditeur, à traiter avec attribution systématique des métriques auto-déclarées.
