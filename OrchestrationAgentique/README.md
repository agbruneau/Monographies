# AgentMesh — « L'autonomie encadrée »

Dépôt de rédaction d'une **monographie exhaustive** sur l'interopérabilité et l'orchestration agentique en écosystème d'entreprise de services financiers au Canada — protocoles ouverts (MCP, A2A, AP2, AGNTCY), cadre réglementaire canadien (E-23, AMF, ACVM, cadre bancaire) et blueprint d'intégration d'entreprise. État des lieux 2024-2026.

| Champ | Valeur |
|---|---|
| Livrable | **mono-v1.0** — publiée le 17 juillet 2026 |
| Volumétrie | **92 059 mots**, 29 pièces (24 chapitres, avant-propos, annexes A-D) |
| Socle factuel | **46 entrées** F-01 à F-48 (F-12 à F-14 non attribués ; F-23b) |
| Conformité | CA-1 à CA-8 |

## Par où commencer

- **Lire la monographie** → [`monographie/README.md`](monographie/README.md) : index de lecture ordonné, à commencer par l'[avant-propos](monographie/00-avant-propos.md). Le [chapitre 13](monographie/03-partie-III/ch-13-pont-frames.md) en est le pivot.
- **PDF assemblé** → [`Monographie.pdf`](Monographie.pdf) (les 29 pièces reliées en un volume).
- **Contribuer / reprendre** → [`CLAUDE.md`](CLAUDE.md) : conventions, garde-fous et procédure de reprise.

## Structure du dépôt

| Chemin | Contenu |
|---|---|
| `monographie/` | Les 29 pièces, un fichier par chapitre (parties I-VII, annexes `90-annexes/`, registre des gels `99-registre-gel.md`) |
| `doc/` | Documents de gouvernance et sources — [`PRD.md`](doc/PRD.md), [`PRDPlan.md`](doc/PRDPlan.md), [`TOC.md`](doc/TOC.md), [`audit.md`](doc/audit.md), et les deux PDF académiques du socle (F-36, F-37) |
| `verification/` | Rapports de revalidation et grille de conformité CA-1..CA-8 |
| `build/` | Pipeline de rendu PDF (Pandoc → Typst) |
| `Monographie.md` / `Monographie.pdf` | Assemblage versionné des 29 pièces et son rendu |

## Gouvernance

Documents par ordre d'autorité — **en cas de conflit, le PRD prime** :

1. [`PRD.md`](doc/PRD.md) — contenu, socle factuel, garde-fous (R-1..R-8), critères d'acceptation ;
2. [`PRDPlan.md`](doc/PRDPlan.md) — plan d'exécution et boucle qualité par chapitre (§4.2) ;
3. [`TOC.md`](doc/TOC.md) — titre, abstract, table des matières commentée.

Toute affirmation factuelle centrale est tracée à une entrée du socle F-xx avec son niveau de preuve — **[A]** vote adversarial 3-0 > **[B]** source primaire extraite > **[C]** repérage à confirmer.

## Régénérer le PDF

Après toute modification des chapitres :

```bash
python build/assemble.py            # concatène monographie/ → Monographie.md
bash   build/build-pdf.sh Monographie.md   # → Monographie.pdf (US-letter, gabarit build/fesp.template)
```

## Avertissements

- **Aucun avis juridique ni conseil d'investissement** : l'ouvrage rapporte des textes et en propose des lectures d'architecture qui engagent l'auteur seul.
- **Aucune recommandation de fournisseur** : la Partie VII instancie le blueprint sur le portefeuille d'IBM à titre de cas documenté, pas de verdict comparatif.
- **L'ouvrage se périme par morceaux** : chaque pièce porte sa date de gel ; les échéances de revalidation sont suivies dans [`CLAUDE.md`](CLAUDE.md) et au [chapitre 24](monographie/07-partie-VII/ch-24-lacunes-revalidation.md).
- **Onze lacunes ouvertes** sont exposées plutôt que comblées ([chapitre 21](monographie/06-partie-VI/ch-21-frontiere.md)).
