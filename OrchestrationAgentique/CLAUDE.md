# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Nature du dépôt

Projet documentaire. Livrable : une monographie exhaustive — titre de travail « L'autonomie encadrée » — sur l'interopérabilité et l'orchestration agentique en écosystème d'entreprise de services financiers au Canada. Les deux PDF académiques dans le socle (F-36, F-37, dans `doc/`) sont des sources analysées et intégrées au PRD ; ne pas les supprimer.

**Emplacement des documents de gouvernance** *(depuis le 17 juillet 2026)* : PRD.md, PRDPlan.md, TOC.md, audit.md et les deux PDF académiques vivent dans `doc/`. ⚠ Le déplacement a cassé des références : `build/assemble.py` lit encore `TOC.md` à la racine (pipeline PDF hors service), et les liens relatifs de `monographie/README.md`, `verification/relecture-CA.md` et des fichiers de `doc/` eux-mêmes pointent vers les anciens chemins.

**Pipeline de rendu PDF** *(ajouté le 17 juillet 2026)* : `build/` contient un pipeline Pandoc → Typst calqué sur le dépôt gabarit `InteroperabiliteAgentique`. `build/assemble.py` concatène les 29 pièces de `monographie/` en `Monographie.md` (retrait de l'appareil de rédaction interne — tableaux d'en-tête, bandeaux de correctif, blocs de gouvernance HTML ; préfixage des notes ; pages de garde de partie), puis `bash build/build-pdf.sh Monographie.md` produit `Monographie.pdf` (US-letter, Liberation Sans/Arial, gabarit `build/fesp.template`). **Régénérer après toute modification des chapitres.** Le PDF et le `Monographie.md` assemblé sont versionnés ; les artefacts intermédiaires du build sont ignorés (`.gitignore`).

Documents de gouvernance, par ordre d'autorité :
1. [PRD.md](doc/PRD.md) — contenu, socle factuel, garde-fous, critères d'acceptation (prime en cas de conflit) ;
2. [PRDPlan.md](doc/PRDPlan.md) — plan d'exécution (phases P0–P4, boucle qualité par chapitre §4.2) ;
3. [TOC.md](doc/TOC.md) — titre, abstract et table des matières commentée (livrable J-2) : découpage en 24 chapitres tracés au socle ; toute modification du découpage passe par ce fichier (version++).

Arborescence de rédaction : `monographie/` (gabarits de chapitres en place — un fichier par chapitre, en-tête avec date de gel, socle F-xx et garde-fous R-x ; annexes A–D dans `90-annexes/` ; registre des gels dans `99-registre-gel.md`) et `verification/` (rapports de revalidation, grille CA). Ne pas rédiger un chapitre sans passer par la boucle qualité de PRDPlan §4.2, et renseigner le registre de gel à chaque fusion.

## PRD.md fait autorité

Toute contribution (rédaction de la monographie, mise à jour du PRD, nouvelle recherche) doit se conformer au PRD :

- **Socle factuel (§7, 46 entrées — identifiants F-01 à F-48, F-12 à F-14 non attribués, plus F-23b)** : toute affirmation factuelle centrale doit être tracée vers une entrée F-xx ou vers une nouvelle source primaire de qualité équivalente, avec son niveau de preuve — **[A] > [B] > [C]** : [A] vote adversarial 3-0 (le plus haut), [B] source primaire lue/extraite avec citation mais sans vote, [C] repérage à confirmer. Une entrée [C] ne porte jamais un fait central sans élévation préalable. ⚠ **Une extraction de source primaire n'élève pas une entrée déjà votée 3-0 : elle l'enrichit d'un contenu de rang inférieur** (faute commise sur F-09 le 16 juill. 2026, rectifiée au PRD v1.7 — l'entrée porte le marquage [A/B mixte]). ⚠ **Ne jamais dériver le cardinal du socle d'une borne d'identifiants** : le total (46) et le plus haut identifiant ont coïncidé jusqu'en P0, et la coïncidence rendait l'erreur indétectable. Commande de référence : PRD §7 chapeau.
- **Garde-fous (§8.1, R-1 à R-8)** : affirmations réfutées ou non vérifiées, interdites dans tout texte (ex. : ACP comme standard vivant, FDX comme standard retenu du cadre bancaire canadien, conformité E-23/B-13 revendiquée par IBM, confusion des trois emplois de « (agentic) control plane »).
- **Attribution (§8.2)** : toutes les métriques auto-déclarées (institutions, éditeurs, Linux Foundation, études d'analystes commandées) sont attribuées à leur source à chaque occurrence.
- **Correspondances réglementaires** : tout lien produit↔réglementation (E-23, AMF, art. 12.1, B-13) est marqué « documenté » ou « inférence d'auteur » (tableau B.3). Ne jamais présenter une inférence comme un fait.
- **Neutralité fournisseur (§8.4)** : IBM est un cas d'instanciation documenté (Partie VII, Annexe B), pas une recommandation.

## Sensibilité temporelle

Le domaine évolue par trimestres ; les faits sensibles au temps s'échelonnent pour l'essentiel de déc. 2023 à juill. 2026 (contexte historique antérieur inclus). Avant toute réutilisation ou publication : revalider les faits « chauds » listés en §8.3 du PRD (A2A/MCP, RTR, règlements du cadre bancaire, C-36). La clôture de l'acquisition Confluent, fait chaud jusqu'en P0, est résolue (17 mars 2026 — F-41). Chaque chapitre rédigé porte une date de gel de l'information.

## Conventions

- **Langue** : livrables et prose en français soutenu ; terminologie technique anglaise entre parenthèses à la première occurrence ; citations verbatim en langue originale.
- **Commits** : messages en anglais, format Conventional Commits (`docs(prd): ...` pour le PRD ; `docs(mono): draft chapter N — <slug anglais court>` pour les chapitres), sujet ≤ 50 caractères quand possible (plafond 72), mode impératif.
- **Versionnage du PRD** : incrémenter la version et la table des jalons (§12) à chaque enrichissement substantiel ; consigner la méthode de toute nouvelle recherche en Annexe A.
- **Chirurgie** : ne modifier que les sections concernées du PRD ; ne pas renuméroter les entrées F-xx ni les garde-fous R-x existants (ils sont cités en références croisées).

## Conventions additionnelles

- **Commits et push** : règle permanente du projet — committer tous les changements et pousser sur `main` à la fin de chaque tâche, sans attendre de demande.

## État du projet (17 juillet 2026)

**Monographie publiée et étiquetée `mono-v1.0`.** Tous les jalons sont atteints (J-0 à J-5, PRD §12) et toutes les phases exécutées (P0 à P4, PRDPlan §1.4 — 47 activités closes, aucune en attente). Documents à jour : **PRD v1.10, PRDPlan v1.4, TOC v1.5** (versions relevées le 17 juill. 2026 par les suites de l'audit global — [`audit.md`](doc/audit.md)).

Livrable : **29 pièces, 92 059 mots** (90 362 à la publication ; re-mesuré après la passe corrective de l'audit) — 24 chapitres, avant-propos, annexes A à D. Index de lecture : [`monographie/README.md`](monographie/README.md). Socle : **46 entrées** (F-01 à F-48 ; F-12 à F-14 non attribués ; F-23b). Conformité : **CA-1 à CA-8** — le « 8/8 » de la relecture de publication a été **partiellement démenti** par l'audit du 17 juill. 2026 (CA-1 et CA-7 portaient des écarts, depuis corrigés) ; voir l'addendum de [`verification/relecture-CA.md`](verification/relecture-CA.md). Revalidation temporelle du 17 juillet 2026 : aucun amendement matériel des faits chauds. **Audit global du 17 juillet 2026 : deux élévations du socle (F-33, F-43), 15 constats majeurs et 44 mineurs traités, relus adversarialement.**

**Ce qui reste vivant après publication** — l'ouvrage se périme par morceaux, et ces points sont datés :

| Échéance | Objet | Pièces touchées |
|---|---|---|
| **28 juill. 2026** | Révision de la spécification MCP (protocole sans état, retrait de `Mcp-Session-Id`) — **confirmée à sa source** le 17 juill. | ch. 1, 2, 7 |
| après le 26 août 2026 | Texte final du règlement du cadre bancaire ; arrêté désignant l'organisme de normalisation technique (lève ou maintient **R-5**) | ch. 14, 15, 24 |
| cible T4 2026 | Lancement effectif du RTR (lève ou maintient la **réserve F-29**) | ch. 15, 24 |
| 1er mai 2027 | Entrées en vigueur d'E-23 et de la ligne directrice IA de l'AMF — le « futur proche » de l'ouvrage devient du présent | ch. 9, 11, 20 |
| continue | Trajectoire de C-36 | ch. 10 |

**Onze lacunes ouvertes** (PRD §10), toutes exposées et aucune comblée. La plus coûteuse : le contenu de la ligne directrice IA de l'AMF n'est pas au socle (§10.4 ; `lautorite.qc.ca` renvoie 403 aux outils) — le chapitre pivot n'en dérive aucune contrainte, sur l'un des deux textes qui fixent l'échéance du 1er mai 2027. **Ne jamais combler une lacune §10 avec du non-vérifié.**

**Toute reprise** (mise à jour, republication, dérivation) repasse par : revalidation des faits chauds (PRD §8.3) → amendement du **socle d'abord**, jamais des chapitres seuls (PRDPlan §7) → boucle qualité §4.2, **dont la relecture adversariale par un relecteur distinct** → date de gel au registre → version++ des documents touchés.
