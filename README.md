# Interopérabilité, orchestration et entreprise agentiques — un triptyque et sa veille

Travaux d'André-Guy Bruneau sur les agents d'IA en écosystème d'entreprise, et plus
particulièrement en services financiers. Le dépôt réunit **trois monographies** conçues en
progression — les protocoles, puis les cadres réglementaires, puis l'organisation qui doit les
faire tenir ensemble — et **une veille technologique autonome** qui les traverse et les tient à jour.

> **Où entrer.** Le lecteur pressé lit la [veille technologique](Veille%20Technologique.md) : c'est
> l'état de l'art le plus récent (18 juillet 2026), et le seul document qui cite les trois autres.
> Le lecteur méthodique suit l'ordre des volumes, du général au spécifique.

## Les quatre livrables

| | **Veille technologique** | **Vol. I — Interopérabilité** | **Vol. II — Orchestration** | **Vol. III — Entreprise** |
|---|---|---|---|---|
| **Dossier** | racine du dépôt | [`1 - InteroperabiliteAgentique/`](1%20-%20InteroperabiliteAgentique/) | [`2 - OrchestrationAgentique/`](2%20-%20OrchestrationAgentique/) | [`3 - EntrepriseAgentique/`](3%20-%20EntrepriseAgentique/) |
| **Titre** | Interopérabilité agentique et orchestration des processus d'affaires en entreprise | Interopérabilité agentique en entreprise dans le domaine des services financiers | L'autonomie encadrée | L'entreprise agentique — la fabrique de confiance |
| **Rôle** | État de l'art vérifié, mis à jour par éditions | Cadre général, mondial et théorique | Cas canadien réglementé, instruit au grain du droit | Le verrou commun : identité, maillage, exploitation |
| **Portée** | Mondiale | Mondiale (UE / É.-U. / R.-U. / Asie) | Canada-Québec (E-23, AMF, Loi 25, ACVM, Lynx/RTR) | Organisation et cycle de vie (NHI, *agent mesh*, AgentOps) |
| **Thèse** | « L'agent d'entreprise fiable de 2026 est un agent *enveloppé* » | « Autonomie graduée sous contrôle de finalité » | « Autonomie encadrée » (*framed autonomy*) | « La confiance ne se décrète pas, elle se fabrique » |
| **Méthode** | Revue structurée, vérification adverse à trois votants | Formalisme d'ingénierie (ArchiMate 4, ADS « Boréalis ») | Socle factuel F-01…F-48, niveaux de preuve [A]/[B]/[C] | Reconduction des deux méthodes (annoncée) |
| **Gel de l'information** | 18 juillet 2026 | Juin 2026 | 16-17 juillet 2026 | — (proposition) |
| **État** | Publiée (107 p., 218 références) | Rédaction terminée (569 p.) | Publiée, `mono-v1.0` (387 p.) | **Cadrage seul** — table des matières v0.4 |

## Veille technologique — le document transversal

[`Veille Technologique.md`](Veille%20Technologique.md) → `Veille Technologique.pdf` (**107 p.**,
13 sections numérotées, **218 références**). Revue vérifiée où chaque énoncé factuel est adossé à
une source primaire consultée et soumis à contradiction — vérificateurs indépendants chargés de
*réfuter*, contre-vérification directe sinon. Elle couvre les trois protocoles structurants (MCP,
A2A, ANP), leur gouvernance, l'adoption documentée, la sécurité, et **six couches** que la pile
protocolaire laisse implicites : événementielle, de contrôle, transactionnelle, sémantique, de
confiance et d'orchestration des processus d'affaires.

**Elle est aussi le point d'articulation du triptyque.** L'édition du 18 juillet 2026 intègre deux
corpus compagnons du dépôt et en fait des sections à part entière :

- **§4.12 — « De la spécification au code »** confronte le corpus documentaire à l'épreuve d'une
  implémentation de référence : le démonstrateur `Borealis-Go` du Vol. I (référence [217]) ;
- **§8.4 — « L'instruction sectorielle canadienne »** reprend le croisement systématique entre
  trajectoire protocolaire et textes canadiens établi par le Vol. II (référence [218]).

L'auto-citation est assumée et divulguée ; ses limites (circularité possible, implémentation
unique, chiffres institutionnels auto-déclarés) sont exposées en section 10 de la veille.

*Historique des éditions : 2, 4, 7, 12, 13, 15 puis 18 juillet 2026. Chaque édition ajoute une
couche ou un corpus et revérifie les faits périssables.*

## Vol. I — Interopérabilité agentique

Monographie de science et génie informatique, construite **en spirale du général au spécifique**,
pour un double public (recherche et praticien-architecte). Invariant transversal : *découplage,
contrat, évolution*.

- **Monographie** (`Monographie.pdf`, **569 p.**) — 7 chapitres : interopérabilité des SI, IA
  agentique, interopérabilité agentique, en entreprise, dans le domaine financier, blueprint
  ArchiMate, horizon 2027-2032.
- **Architecture détaillée de solution** (Annexe B) — la monographie projetée sur une entreprise
  fictive, la *Coopérative financière Boréalis*, consolidée sur la pile IBM ; 18 sections,
  6 sous-annexes, 28 diagrammes Mermaid, rendus dans le PDF principal.
- **Article de synthèse** (`Synthese Monographie.pdf`, **69 p.**) — la monographie condensée au
  format académique, avec bibliographie propre.
- **Démonstrateur `Borealis-Go/`** — code Go exécutable matérialisant l'ADS : **5 agents A2A** et
  **4 serveurs MCP** orchestrant une pré-qualification de crédit (jamais un octroi ferme), sur les
  SDK officiels des deux protocoles. **12 ADR**, journal d'audit à chaîne de hachage, vérification
  adverse à chaque phase, invariants critiques prouvés par mutation ; couverture déclarée 96,2 %
  au rapport final. C'est ce démonstrateur qui fournit la §4.12 de la veille.

## Vol. II — L'autonomie encadrée

Monographie sur l'interopérabilité et l'orchestration agentique en services financiers canadiens,
publiée et étiquetée `mono-v1.0`. **92 059 mots** en 29 pièces (24 chapitres, avant-propos,
annexes A-D) selon son README ; `Monographie.pdf` **387 p.** ; article de synthèse **66 p.**

Sa contribution la plus citable est un **résultat négatif** : en croisant trois protocoles
(MCP, A2A, AP2) et cinq corpus de textes canadiens, **aucun lien documenté par source primaire** —
quinze croisements, zéro lien. D'où sa thèse probatoire : sous exigence réglementaire stricte, le
cadre déterministe invoque les agents, jamais l'inverse, parce que le cadre est la seule pièce
dont l'exploitant puisse démontrer la teneur devant un tiers.

Sa méthode est son autre apport : socle factuel de **46 entrées** (F-01 à F-48) cotées par niveau
de preuve — **[A]** vote adversarial 3-0 > **[B]** source primaire extraite > **[C]** repérage —,
huit garde-fous de formulation, onze lacunes exposées plutôt que comblées.

## Vol. III — L'entreprise agentique

**Au stade du cadrage seul.** Le dossier ne contient à ce jour qu'une table des matières commentée
([`TOC.md`](3%20-%20EntrepriseAgentique/TOC.md), v0.4 du 18 juillet 2026, statut *proposition*) :
aucun chapitre rédigé, **aucun socle factuel constitué** — les « socles candidats » y sont des
repérages [C] à instruire. Cadrage annoncé : 28 chapitres en 9 parties, ≈ 100 000 mots, autour de
trois capacités — *émettre* une identité opposable (le passeport d'agent), l'*appliquer* au
maillage d'agents, l'*exploiter* dans la durée (AgentOps) — sous l'horloge post-quantique.

Le volume naît des lacunes des deux précédents : identité non humaine et délégation multi-saut
(verrou identifié au Vol. I), mécanique des attaques et valeur cryptographique des Agent Cards
(questions ouvertes du Vol. II).

## Ordre de lecture et renvois

**Vol. I → Vol. II → Vol. III**, la veille servant d'entrée rapide ou de mise à jour.

- **Vol. II présuppose Vol. I** pour la théorie du découplage, l'ingénierie des agents LLM,
  l'anatomie des protocoles, la sécurité de la couche agentique et la cryptographie post-quantique.
- **Vol. I illustre mondialement** ce que **Vol. II instruit au grain du droit canadien**.
- **Vol. III prolonge les deux** sur leur verrou commun, l'identité et son exploitation.
- **La veille recoupe les trois** et les cite explicitement (§4.12 et §8.4).
- Un lecteur pressé côté canadien peut entrer directement par le **chapitre 13** du Vol. II
  (« le pont : des contraintes réglementaires aux frames déterministes »), son pivot.

## Divergences factuelles entre volumes

Deux faits datés divergent d'un corpus à l'autre. Ils sont **signalés, non arbitrés** — la veille
les expose en §8.4, et le lecteur doit les trancher à sa date de citation :

| Objet | Vol. II (gel 16-17 juill.) | Veille (18 juill.) |
|---|---|---|
| Ligne directrice IA de l'AMF — version finale | 30 mars 2026, avec dette de vérification déclarée (`lautorite.qc.ca` renvoie 403 aux outils) | 7 avril 2026 |
| Gouvernance d'AP2 | aucun transfert documenté au socle | don à la FIDO Alliance, 28 avril 2026 |

L'entrée en vigueur du 1er mai 2027 est, elle, concordante entre les corpus. La divergence tient à
une frontière de socle plus qu'à un désaccord : deux corpus vérifiés de bonne foi peuvent diverger
sur un fait daté selon leurs périmètres de sources — argument pour le millésimage systématique.

> ⚠ Le fichier `commun/faits-partages.md`, évoqué par le cadrage du Vol. III comme source unique
> de vérité pour les faits partagés, **n'existe pas encore** (marqué « à créer » dans sa TOC).
> En son absence, chaque volume porte ses propres faits datés.

## Structure du dépôt

```
.
├── README.md                          ← ce fichier (avant-propos croisé)
├── Veille Technologique.md / .pdf     ← veille autonome, 18 juillet 2026 (107 p., 218 réf.)
├── 1 - InteroperabiliteAgentique/     ← Vol. I
│   ├── Chapitres/                       7 chapitres + bibliographies + Annexe B (ADS)
│   ├── Monographie.md / .pdf            assemblage (569 p.)
│   ├── Synthese Monographie.md / .pdf   article de synthèse (69 p.)
│   ├── Borealis-Go/                     démonstrateur Go (MCP + A2A), 12 ADR
│   ├── build/                           pipeline FESP (Mermaid → Pandoc → Typst)
│   └── index.html                       page de présentation (GitHub Pages)
├── 2 - OrchestrationAgentique/        ← Vol. II
│   ├── monographie/                     29 pièces (parties I-VII, annexes, registre des gels)
│   ├── doc/                             PRD, PRDPlan, TOC, audit — gouvernance
│   ├── verification/                    revalidations et grille de conformité CA-1..CA-8
│   ├── build/                           pipeline Pandoc → Typst
│   ├── Monographie.md / .pdf            assemblage (387 p.)
│   └── Synthese Monographie.md / .pdf   article de synthèse (66 p.)
└── 3 - EntrepriseAgentique/           ← Vol. III
    └── TOC.md                           table des matières commentée (v0.4) — seul livrable
```

## Construire les PDF

Trois chaînes distinctes, à lancer depuis le dossier concerné.

**Veille technologique** (racine) — invocation Pandoc directe, gabarit Typst par défaut :

```bash
pandoc "Veille Technologique.md" --pdf-engine=typst --toc -o "Veille Technologique.pdf"
```

**Vol. I** — pipeline FESP, avec pré-rendu des 28 diagrammes Mermaid :

```bash
bash build/build-pdf.sh                              # Monographie.pdf
bash build/build-pdf.sh "Synthese Monographie.md"    # Synthese Monographie.pdf
```

**Vol. II** — assemblage des 29 pièces, puis même pipeline :

```bash
python build/assemble.py                    # monographie/ → Monographie.md
bash   build/build-pdf.sh Monographie.md    # → Monographie.pdf
```

**Prérequis :** Pandoc ≥ 3.1.7, Typst ≥ 0.12, `python3` + `pypdf` ; polices Liberation Sans et
DejaVu Sans (pipeline FESP), New Computer Modern (veille) ; pour les diagrammes, Node ≥ 18 +
[`@mermaid-js/mermaid-cli`](https://github.com/mermaid-js/mermaid-cli) et un Chromium. Sous
Windows, exporter `PYTHONUTF8=1` avant le build. **Règle permanente :** régénérer et versionner le
PDF avec sa source — jamais la source seule.

## Ce qui reste vivant

Le domaine se périme par trimestres, et ces corpus par morceaux. Échéances datées à revalider
avant toute réutilisation ou publication :

| Échéance | Objet | Documents touchés |
|---|---|---|
| 28 juillet 2026 | Révision de la spécification MCP (protocole sans état) | Veille §4.1 ; Vol. I ch. 3 ; Vol. II ch. 1, 2, 7 |
| après le 26 août 2026 | Texte final du règlement du cadre bancaire canadien ; arrêté désignant l'organisme de normalisation | Veille §8.4 ; Vol. II ch. 14, 15, 24 |
| cible T4 2026 | Lancement effectif du RTR — cible précédée de quatre cibles abandonnées depuis 2019 | Veille §8.4 ; Vol. II ch. 15, 24 |
| 2 décembre 2026 | Marquage des contenus générés (règlement européen sur l'IA) | Veille §8.1, §12 |
| **1er mai 2027** | Entrée en vigueur simultanée d'E-23 (BSIF) et de la ligne directrice IA de l'AMF | Veille §4.11.5, §8.4 ; Vol. I ch. 5 à 7 ; Vol. II ch. 9, 11, 20 |
| continue | Trajectoire du projet de loi C-36 | Veille §8.4 ; Vol. II ch. 10 |

## Avertissements

- **Aucun avis juridique ni conseil d'investissement.** Ces ouvrages rapportent des textes et en
  proposent des lectures d'architecture qui engagent leur auteur seul.
- **Aucune recommandation de fournisseur.** Les instanciations sur une pile d'éditeur (IBM
  notamment) sont des cas documentés, pas des verdicts comparatifs.
- **Statuts et chiffres.** Les métriques d'adoption sont, sauf mention contraire, auto-déclarées
  par les acteurs et attribuées comme telles ; les statuts *preview* ne sont jamais présentés comme
  *disponibilité générale* ; les projections d'analystes portent leur millésime.
- **Lacunes exposées, non comblées.** Le Vol. II en recense onze ; la veille, seize questions
  ouvertes. Aucune n'est comblée par une source de moindre qualité.
- **Assistance par agents.** Ces travaux ont été produits avec l'assistance de pipelines de
  recherche multi-agents, selon les méthodes de vérification décrites dans chaque document ; la
  responsabilité éditoriale est celle de l'auteur.

## Notes de maintenance

Deux désynchronisations connues, **signalées et non corrigées** dans cette révision :

- Le `README.md` et le `CLAUDE.md` du **Vol. I** décrivent encore `Veille Technologique.md` comme
  un livrable interne au dossier (« état de l'art au 7 juillet 2026, 48 p., 116 références »),
  alors que la veille a été déplacée à la racine du dépôt et porte désormais l'édition du
  18 juillet 2026 (107 p., 218 références). Ils annoncent aussi 11 ADR pour `Borealis-Go`, qui en
  compte 12.
- Le `CLAUDE.md` du **Vol. II** signale que le déplacement de ses documents de gouvernance vers
  `doc/` a cassé des liens relatifs et que `build/assemble.py` lit encore `TOC.md` à la racine.

Conventions de rédaction et règles de travail : voir le `CLAUDE.md` de chaque volume.
