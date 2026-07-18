# CLAUDE.md

Guide pour Claude Code (claude.ai/code) sur ce dépôt.

## Nature du dépôt

Projet d'**écriture** en **français canadien** autour de l'**interopérabilité agentique** : **quatre livrables rédigés**, plus un **démonstrateur de code** qui les accompagne ([`Borealis-Go/`](Borealis-Go/), régi par son propre [`CLAUDE.md`](Borealis-Go/CLAUDE.md)) :

1. **Monographie** de science et génie informatique, à double public — **recherche** (modèles, formalismes, état de l'art) et **praticien-architecte** (normes, protocoles, mises en œuvre). Sept chapitres en spirale du général au spécifique ; **socle des sources arrêté à juin 2026** ; invariant transversal rappelé à chaque couche : *découplage, contrat, évolution*.
2. **Architecture détaillée de solution (ADS)**, intégrée comme **Annexe B** de la monographie — source [`Chapitres/Annexe B - Architecture de Solutions.md`](Chapitres/Annexe%20B%20-%20Architecture%20de%20Solutions.md). Elle projette la monographie (surtout ch. 5-6) sur une entreprise fictive (*Coopérative financière Boréalis*) et une pile **IBM** consolidée : livrable d'ingénierie *prêt au déploiement* (diagrammes, contrats, NFR, topologie, runbooks), de nature distincte du corps doctrinal mais rendu dans le même PDF.
3. **Article de synthèse** autonome — revue condensée de la monographie au format académique (`Synthese Monographie.md`).
4. **Veille technologique** autonome — état de l'art au 7 juillet 2026 sur les protocoles ouverts, la normalisation, l'identité et la sécurité agentiques, chaque énoncé adossé à une source primaire (`Veille Technologique.md`).
5. **Démonstrateur Boréalis** — code Go (**MCP + A2A**, pré-qualification de crédit) matérialisant le PRD, dans [`Borealis-Go/`](Borealis-Go/). Livrable de nature distincte : conventions, commandes (`make check`, gate ≥ 90 %) et règles non négociables dans [`Borealis-Go/CLAUDE.md`](Borealis-Go/CLAUDE.md), qui prime dans ce répertoire.

**État : rédaction terminée, démonstrateur clos** (audit 27/27, gate vert à 96,2 %). Le travail courant est la finalisation et la maintenance — corrections, vérification adverse des citations, régénération des PDF. Outillage : `git`, le pipeline de rendu (voir *Commandes*) et la chaîne Go du démonstrateur (voir `Borealis-Go/CLAUDE.md`).

## Organisation

Contenu rédigé dans [`Chapitres/`](Chapitres/), deux fichiers par chapitre, plus l'ADS :

| Fichier | Rôle |
|---|---|
| `Chapitre N - {Sujet}.md` | chapitre rédigé (corps) |
| `Chapitre N - Bibliographie - {Sujet}.md` | références vérifiées du chapitre |
| `Annexe B - Architecture de Solutions.md` | corps de l'ADS, repris tel quel dans `Monographie.md` |

Assemblages et livrables à la racine :

- [`Monographie.md`](Monographie.md) → `Monographie.pdf` (**569 p.**) — **source unique** du PDF principal : liminaires (page titre Typst, résumé, table des matières, abréviations) + 7 chapitres + Bibliographie générale + **Annexe A** (documents d'accompagnement) + **Annexe B** (l'ADS : 18 sections + vues *blueprint* §0.1-0.2, 6 sous-annexes, 28 diagrammes Mermaid). Rendu FESP (Pandoc → Typst), diagrammes pré-rendus.
- [`Synthese Monographie.md`](Synthese%20Monographie.md) → `Synthese Monographie.pdf` (**69 p.**) — article de synthèse autonome (12 sections + bibliographie propre), **même pipeline FESP**, sans diagramme.
- [`Veille Technologique.md`](Veille%20Technologique.md) → `Veille Technologique.pdf` (**48 p.**) — veille autonome (13 sections numérotées + 116 références), **pipeline distinct** (gabarit Pandoc/Typst par défaut, police New Computer Modern ; voir *Commandes*).
- [`Borealis-Go/`](Borealis-Go/) — démonstrateur Go (voir livrable 5 ci-dessus) : `cmd/` (9 binaires), `internal/`, `pkg/`, `docs/adr/` (11 ADR), PRD et plan d'exécution inclus.
- `index.html` — page de présentation (GitHub Pages) ; lie les trois PDF.

Les trois PDF (monographie, synthèse, veille) sont publiés sur GitHub Pages.

Progression de la monographie (chaque chapitre suppose les précédents) :

1. Interopérabilité des SI (fondements, intégration d'entreprise)
2. IA agentique
3. Interopérabilité agentique (MCP, A2A…)
4. … en entreprise
5. … dans le domaine financier
6. Blueprint d'architecture d'entreprise ArchiMate (formalise les ch. 1-5)
7. Horizon 2027-2032 (chapitre prospectif / *capstone* projetant les ch. 1-6)

## Commandes

Deux chaînes de rendu **distinctes**.

**Monographie et synthèse** — pipeline FESP (`build/build-pdf.sh`) :

```bash
bash build/build-pdf.sh                            # Monographie.pdf (défaut ; pré-rend les 28 diagrammes Mermaid)
bash build/build-pdf.sh "Synthese Monographie.md"  # Synthese Monographie.pdf (sans diagramme)
```

Le script prend un `.md` source en argument (défaut `Monographie.md`) et en déduit le `.pdf`. Les blocs ` ```mermaid ` sont pré-rendus en SVG (mermaid-cli) puis injectés avant Pandoc→Typst ; si `mmdc` est absent, ils restent en bloc de code (mode dégradé, sans échec). La pagination liminaire romaine→arabe ne s'applique qu'à `Monographie.md`. Gabarit : `build/fesp.template` (police Liberation Sans, repli Arial).

**Veille technologique** — invocation Pandoc directe (gabarit Typst *par défaut*, table des matières générée, police New Computer Modern) :

```bash
pandoc "Veille Technologique.md" --pdf-engine=typst --toc -o "Veille Technologique.pdf"
```

La veille n'utilise **pas** `build/build-pdf.sh` : son identité visuelle (police, mise en page) provient du gabarit Pandoc/Typst par défaut. Les sauts de page (avant la table des matières et l'introduction) sont pilotés par des blocs Typst bruts ` ```{=typst} #pagebreak(weak: true) ``` ` dans la source.

**Prérequis :** Pandoc ≥ 3.1.7, Typst ≥ 0.12, `python3` + `pypdf`, polices Liberation Sans + DejaVu Sans (repli déclaré dans le gabarit FESP : Arial + Segoe UI Symbol) et New Computer Modern (veille). Diagrammes : Node ≥ 18 + `@mermaid-js/mermaid-cli` (`mmdc`) + un Chromium (surcharges `MMDC=…`, `PUPPETEER_CONFIG=…`).

**Notes d'environnement :** sous Windows, exporter `PYTHONUTF8=1` avant le build (sinon le pré-rendu plante sur les caractères hors cp1252). Les SVG Mermaid sont mis en cache (`DIAGCACHE`) et réutilisés d'un build à l'autre.

## Convention des chapitres rédigés

- **Sections numérotées, préfixées du numéro de chapitre** : sections `N.0, N.1, N.2…` (`##`) ; sous-sections `N.1.1` (`###`) ; items feuilles `N.1.1.1` (`####`). Le premier nombre est **toujours** le numéro du chapitre (ch. 1 → `1.6.3.1`, ch. 3 → `3.2.1`). *Exception documentée :* les ch. 2 et 4 ouvrent à `N.1` (pas de section `N.0`) — ne pas renuméroter.
- **En-tête de chapitre** : titre `#`, bloc `> **Résumé.**`, puis ligne *Public visé / Fil conducteur / dates*.
- **Renvois inter-sections** en clair : intra-chapitre `(voir §1.1.1)`, `(détaillé en 1.6.2)` ; inter-chapitres avec le numéro du **chapitre cible** `(cf. ch.1 §1.6.3.1)`, `(renvoi ch. 7)`. Tout renvoi doit pointer vers une cible existante.
- **Invariant transversal** rappelé à chaque couche : *découplage, contrat, évolution*.

## Convention de l'ADS (Annexe B)

- **Numérotation native conservée, décalée d'un niveau sous `# Annexe B`** : sections `## 0…17` (plus `## 0.1`-`0.2`, vues *blueprint* d'ouverture) + `## Annexe A…F` internes ; sous-sections `### N.x`. Ne pas confondre : le `§6` de l'ADS (sécurité) ≠ ch. 6 de la monographie ; l'`## Annexe B` interne (matrice de traçabilité) ≠ l'**Annexe B** de la monographie (= l'ADS entière). Les renvois à la monographie restent explicites (`cf. Monographie ch.6 §6.1.3`).
- **Diagrammes Mermaid *inline*** (aucun fichier image séparé), légendés `**Figure N — …**`. Toute syntaxe doit passer `mmdc` (séquences : pas de `;` ni `{}` dans les messages ; libellés d'arête sans parenthèses nues).
- **Vues d'entreprise en ArchiMate 4** (stéréotypes du registre Monographie ch.6 §6.1.9), équivalent 3.2 noté si requis.
- **Extraits de configuration *illustratifs*** (jamais de secret réel) ; statuts vivants marqués ⚠ ; valeurs chiffrées (SLO, RPO/RTO, débits) = hypothèses à calibrer.

## Convention de la veille technologique

- **Document autonome**, distinct de la monographie (non repris dans `Monographie.md`). En-tête YAML complet (titre, sous-titre, auteur, résumé, `mainfont: "New Computer Modern"`, `section-numbering: "1.1.1"`).
- **Sections `#` numérotées automatiquement** par Pandoc (Introduction = 1 … Conclusion = 13) ; sous-sections `##` → `N.1`. La section 12 (*Horizon prospectif 2027-2030*) trie ses sous-sections en **Programmé / Projeté / Spéculatif** (même logique épistémique que le ch. 7). Les sections liminaires ou finales sans numéro portent `{-}` (`# Divulgation {-}`, `# Références {-}`).
- **Références manuelles** : liste numérotée sous `# Références {-}` dans un bloc `::: {#refs} … :::` ; le corps cite par crochets **littéraux** `[N]` (pas de champ `bibliography` Pandoc, pas de clés `@…`). Toute modification du compte met à jour le bilan de vérification.
- **Sauts de page** via blocs Typst bruts ` ```{=typst} #pagebreak(weak: true) ``` ` ; le saut avant la table des matières passe par `header-includes` (`#show outline: it => [#pagebreak(weak: true) #it]`).

## Convention des fichiers Bibliographie

- **Format de citation** : `Auteurs/Organisme (Année). *Titre*. Support/éditeur. Identifiant stable (DOI / RFC / ISO/IEC / W3C / URL).`
- **Catégories** thématiques ; ordre alphabétique par auteur/organisme (normes ISO/IEC et RFC par numéro croissant). Exceptions déclarées en tête de fichier : le ch. 1 classe les RFC au rang alphabétique de leur auteur/éditeur ; les textes réglementaires des ch. 4-5 suivent l'ordre chronologique.
- **Marqueur ⚠** = *ressource vivante* (doc produit, spec versionnée par date, travail en cours) dont la version/ancre doit être fixée au moment de citer. Ne pas le retirer sans avoir figé la version.
- `cf. §X` relie une entrée aux sections qui s'en servent ; `à re-vérifier à la rédaction` = contrôle à refaire avant publication (conserver tel quel).
- **Accents des chaînes de citation** : préservés *verbatim* même fautifs (`fevrier`), pour ne pas altérer les métadonnées sources. Ne pas « corriger » au fil de l'eau.

## Règles de travail

- **Vérification adverse des citations.** Toute référence ajoutée ou modifiée est vérifiée sur le web (existence, auteurs, année, numéro de norme/RFC/recommandation) avant d'être affirmée exacte. Chaque bibliographie de chapitre porte un bilan (`N références vérifiées → M entrées…`) ; le mettre à jour si on touche au compte. Les **116 références de la veille** ont fait l'objet d'une validation adverse en deux passes (repérage parallèle + confirmation indépendante contre sources primaires) — elles existent et sont conformes.
- **Dates et statut prospectif.** Signaler toute nouveauté postérieure à juin 2026. Le **ch. 7** de la monographie (et la **§12** de la veille) projettent l'avenir et trient chaque énoncé en **PROGRAMMÉ / PROJETÉ / SPÉCULATIF** (faits datés / prévisions millésimées / paris) — ne jamais présenter du spéculatif comme acquis. Dans l'ADS, ne jamais présenter un statut *preview* (⚠) comme GA.
- **Synchronisation et publication.** Après toute édition d'un chapitre : reporter le changement dans `Monographie.md` (et dans `Synthese Monographie.md` si le fait y figure), régénérer les PDF **et les pousser avec la source** — jamais la source seule. La veille est autonome : l'éditer implique de régénérer son PDF. Si la pagination de la monographie change, reporter le nouveau compte aux **deux** endroits d'`index.html` (`"numberOfPages"` du JSON-LD et le compteur animé `data-count`) ainsi que dans `CLAUDE.md` et `README.md` ; mettre à jour le compte de références s'il change.
- **Registre et ton.** Qualité rédactionnelle soutenue ; **ton professionnel et neutre** (pas de marketing, pas de 1ʳᵉ personne). Gloses, titres et notes en français.
- **Longueur.** Chaque **chapitre de la monographie** fait au moins **10 000 mots** ; préserver ce seuil lors des révisions. La synthèse et la veille, condensées par nature, n'y sont pas soumises.
- **Messages de commit** : courts, en français, nommés par livrable touché — `Chapitre 5`, `Chapitre 3 et 4`, `Annexe B`, `Synthèse`, `Veille technologique`.
