# Proposition de passe — après-agentique, itération 04

| Champ | Valeur |
|---|---|
| Date | 22 juillet 2026 |
| Itération | 04 — relève v0.10 nᵒ 6 instruite à ses sources primaires |
| Objet | Fiche d'instruction : mémoire observationnelle — dispositif confirmé verbatim, réserve du journal précisée par une pièce antérieure à la passe |
| Sièges | Vol. IV `TOC.md` v0.10 — ch. 5 (Livre I, ancrage informationnel) et ch. 44 (Livre VII, dérive) ; journal v0.10, relève 6 |
| Statut | **Proposition** — rien n'est intégré, rien n'est poussé ; l'intégration est une décision d'auteur |
| `TOC.md` touché | **Non par cette passe** — entrée de journal prête à intégrer en annexe |

## L'incrément

**La relève v0.10 nᵒ 6 est instruite. Le dispositif est confirmé verbatim ; deux libellés se précisent ; et la réfutation de la réserve du journal exhume une pièce que la passe v0.10 n'avait pas vue — une évaluation publique de l'éditeur datée du 9 février 2026, antérieure de cinq mois à la passe.** La réserve ne tombe pas : elle se reformule.

### Sources (consultées le 22 juillet 2026)

1. S. Bhagwat (Mastra), *Anatomy of a harness*, blogue Mastra, 5 juin 2026. https://mastra.ai/blog/anatomy-of-a-coding-agent
2. Mastra (page de recherche), *Observational Memory: 95% on LongMemEval*, 9 février 2026. https://mastra.ai/research/observational-memory
3. Référence du banc d'essai cité par la pièce 2 : LongMemEval, Wu et al., arXiv 2410.10813 (2024) — banc d'essai public tiers ; identité relevée, non extraite à cette passe.

Pièces 1 et 2 : éditeur décrivant et mesurant son propre produit — **[B]** au mieux, toute métrique auto-déclarée et attribuée à chaque occurrence.

### Constat 1 — [B] : le dispositif, confirmé verbatim

- L'observateur consigne des « *structured observations: the decisions, facts, and state changes* » — la triade du journal v0.10 (« décisions, faits, changements d'état ») est exacte au mot près.
- Le réflecteur compresse le journal d'observations une fois accumulé ; l'agent reçoit trois couches : messages récents bruts, journal d'observations, réflexions comprimées.
- Déclenchement : seuils personnalisables (« *40k tokens by default* »), exécution anticipée « *at 20% intervals* », période d'inactivité à délai propre au fournisseur (5 min Anthropic, 1 h OpenAI/DeepSeek), et changement de fournisseur.
- Le mal visé est nommé « *context rot* » — le terme du journal v0.10 est verbatim.

### Constat 2 — précisions de libellé (à porter à l'intégration)

1. **« Compaction destructrice » n'est pas le terme de la pièce** : le texte écrit « *lossy compaction* » — compaction *avec perte* (elle « *drops the exact wording of a requirement* ») ; « destructrice » surtraduit.
2. **La revendication de préservation n'est pas, dans la pièce du 5 juin, « préserve l'information plutôt que de résumer »** : le verbatim est « *The decisions survive as decisions, not as a paraphrase buried in a summary* » — une revendication de *forme* (les décisions restent des décisions), pas de non-perte. La pièce ne revendique aucune conservation sans perte.

### Constat 3 — la réserve du journal, précisée par une pièce antérieure à la passe

Le journal v0.10 écrit : « aucune propriété de conservation n'est publiquement établie, et "préserve l'information plutôt que de résumer" est une affirmation d'éditeur, **pas une mesure** ». L'instruction nuance les deux moitiés :

- **La première moitié tient.** Aucune propriété de *conservation* (non-perte d'information) n'est établie — ni par l'éditeur, ni par un tiers.
- **La seconde moitié est périmée depuis avant la passe.** Une mesure publique existe : la page de recherche du 9 février 2026 rapporte, sur LongMemEval (banc d'essai public tiers — 500 questions, ≈ 57 M jetons de conversations), des scores auto-déclarés de **94,87 %** (gpt-5-mini), **84,23 %** (gpt-4o), **93,27 %** (gemini-3-pro-preview) et **89,20 %** (gemini-3-flash-preview), avec code et méthode publiés (« *completely open source end to end* ») et des limites déclarées (catégorie multi-session la plus dure ; catégorie préférences volatile, 30 questions). **Trois bornes d'attribution** : (a) évaluation exécutée par l'éditeur lui-même — reproduction indépendante non documentée au 22 juillet 2026 ; (b) LongMemEval mesure le **rappel en question-réponse**, pas une propriété de conservation — les deux moitiés de la réserve ne portent donc pas sur le même objet ; (c) deux des quatre scores sont obtenus sur des modèles en **préversion** (*preview*) — jamais à présenter comme une disponibilité générale. Accessoirement, le titre de la page arrondit 94,87 % en « 95% » : citer le chiffre, pas le titre.

### Conséquences sur les deux atterrissages du journal

- **Ch. 5 — tient et se précise.** L'ancrage informationnel devient bien un artefact **dérivé et daté, produit par d'autres modèles que l'agent** (observateur et réflecteur sont des modèles auxiliaires). L'évaluation LongMemEval, attribuée, donne à la rédaction un point d'appui chiffré — auto-déclaré — que la passe v0.10 ne connaissait pas.
- **Ch. 44 — tient.** La quatrième source de dérive (« la mémoire réécrite ») demeure : la mémoire observationnelle est elle-même un mécanisme de réécriture par modèles auxiliaires, dont la qualité est auto-mesurée sur un banc de rappel — la dérive de la mémoire n'y est pas mesurée en production.

### Tri épistémique

- Dispositif, triade, seuils, « *context rot* », « *lossy compaction* » : **[B]** (pièce d'éditeur datée, verbatim cité).
- Scores LongMemEval : **[B] auto-déclarés, attribués** — jamais cités sans « selon l'éditeur » ; préversions marquées.
- Propriété de conservation : **toujours non établie** — la réserve subsiste sous sa forme reformulée.
- Rapprochement d'auteur, déclaré tel : l'écart entre *mesure de rappel* et *propriété de conservation* est une distinction d'auteur de cette fiche.

## Annexe — entrée de journal prête à intégrer (si l'auteur consomme la fiche)

> **Instruction de la relève 6 (date d'intégration à poser).** Pièces extraites le 22 juillet 2026. Dispositif confirmé verbatim (triade des observations, seuil 40 k, intervalles de 20 %, inactivité par fournisseur, « *context rot* »). Deux libellés précisés : « *lossy compaction* » (avec perte, non « destructrice ») ; revendication réelle « *The decisions survive as decisions…* ». **Écart consigné** : la moitié « pas une mesure » de la réserve était périmée dès l'écriture de la passe — une évaluation publique auto-exécutée sur LongMemEval (page de recherche Mastra, 9 février 2026 ; 94,87 % gpt-5-mini, scores attribués, deux modèles en préversion) antédate la v0.10 de cinq mois. La moitié « aucune propriété de conservation établie » tient : LongMemEval mesure le rappel, pas la conservation. La relève 6 passe de « à instruire » à « instruite, [B], réserve reformulée ». Fiche : `proposition-apres-agentique-2026-07-22-04.md`.

Note de réalignement `Conspectus.md` : cardinal « huit relèves » re-mesuré, inchangé ; statut ventilé — trois instruites [B], cinq [C] à instruire.

## Vérification adverse (exécutée avant livraison)

- **Réfutation tentée et productive** : la recherche d'une mesure publique — censée réfuter la réserve — en a trouvé une, datée d'avant la passe ; l'écart est consigné contre le journal plutôt que lissé, et la distinction rappel/conservation empêche de sur-corriger dans l'autre sens.
- **Chiffres** : tous relevés sur la page de recherche elle-même (pas sur la couverture de presse) ; le titre arrondi (« 95% ») n'est pas repris.
- **Cardinaux** : « huit relèves » re-mesuré, inchangé ; aucun décompte du dépôt modifié.
- **Renvois** : ch. 5, 44 cités — tous dans 1-57.
- **Zones gelées** : aucune touchée ; Livre X terminal ; journaux non réécrits.
- **Décision 8** : aucune thèse réécrite.

## Prochaine itération envisagée

Instruire la relève v0.10 nᵒ 5 (extension déclarative, ch. 52/21/20) : identifier et extraire la publication des chercheurs de l'éditeur de sécurité sur l'extension malveillante de fin janvier 2026 (exfiltration et injection d'invite) — la seule des huit relèves dont la source primaire n'est pas encore nommée dans le dépôt (« aucune source primaire n'a été extraite » au journal v0.10).
