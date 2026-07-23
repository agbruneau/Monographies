# Proposition de passe — après-agentique, itération 03

| Champ | Valeur |
|---|---|
| Date | 22 juillet 2026 |
| Itération | 03 — relève v0.10 nᵒ 2 instruite à ses sources primaires |
| Objet | Fiche d'instruction : chaîne d'approbation d'outils (pièces Mastra) — un maillon confirmé verbatim, un énoncé du journal non confirmé |
| Sièges | Vol. IV `TOC.md` v0.10 — ch. 25-26 (Livre IV), ch. 18 (Livre III), ch. 29 (Livre V), ch. 45 (Livre VII) ; journal v0.10, relève 2 |
| Statut | **Proposition** — rien n'est intégré, rien n'est poussé ; l'intégration est une décision d'auteur |
| `TOC.md` touché | **Non par cette passe** — entrée de journal prête à intégrer en annexe ; elle devra consigner l'écart relevé ci-dessous |

## L'incrément

**La relève v0.10 nᵒ 2 est instruite à ses deux sources primaires écrites. Résultat en deux temps : la chaîne de résolution est confirmée mot pour mot — l'ordre du journal est exact — mais l'énoncé « les octrois persistent d'une session à l'autre » n'est confirmé par aucune des deux pièces citées.** Un incrément qui ne survivrait pas tel quel n'est pas édulcoré : l'écart est consigné, avec sa piste de vérification.

### Sources (consultées le 22 juillet 2026)

1. S. Bhagwat (Mastra), *Anatomy of a harness: building a coding agent that can run for hours*, blogue Mastra, 5 juin 2026. https://mastra.ai/blog/anatomy-of-a-coding-agent
2. S. Bhagwat (Mastra), *Announcing Mastra Harness*, blogue Mastra, 18 juin 2026 — version d'éditeur `@mastra/core@1.45.0` confirmée dans la pièce. https://mastra.ai/blog/announcing-agent-harness

Les deux pièces sont un éditeur décrivant son propre produit : niveau **[B]** au mieux (document public daté, extraction citée), toute métrique auto-déclarée et attribuée à chaque occurrence (règle héritée du PRD Vol. II §7.5).

### Constat 1 — [B] : la chaîne de résolution, confirmée verbatim

La pièce du 5 juin 2026 donne la chaîne en code, « *first match wins* » :

> `if (toolPolicy === 'deny') return 'deny';` — 1. *per-tool deny, a hard block*
> `if (state.yolo === true) return 'allow';` — 2. *YOLO auto-approves anything not denied above*
> `if (toolPolicy) return toolPolicy;` — 3. *explicit per-tool allow/ask*
> `if (this.sessionGrantedTools.has(toolName)) return 'allow';` — 4. *session tool grant*
> `if (this.sessionGrantedCategories.has(category)) return 'allow';` — 5. *session category grant*
> `if (categoryPolicy) return categoryPolicy;` — 6. *explicit per-category policy*
> `return 'ask';` — 7. *default to ask*

**L'ordre à sept maillons du journal v0.10 est exact**, y compris sur le point le plus contre-intuitif (l'octroi de session par catégorie **précède** la politique par catégorie). Le refus par outil est décrit comme « *an absolute block that nothing downstream can reopen* ». Le mode d'auto-approbation globale (« YOLO ») existe bien en deuxième maillon.

**Note de méthode consignée** : une première extraction automatisée de cette passe avait inversé les maillons 5 et 6 — c'est la contre-vérification au verbatim qui l'a réfutée, non le journal. L'incident illustre la règle du dépôt : une extraction est du contenu, elle se vérifie comme le reste.

### Constat 2 — non confirmé : la persistance des octrois entre sessions

Le journal v0.10 écrit : « Les octrois **persistent d'une session à l'autre**. » Aucune des deux pièces ne le soutient au 22 juillet 2026 :

- la pièce du 5 juin nomme les octrois `sessionGrantedTools` / `sessionGrantedCategories` — octrois *de session* — et ne mentionne aucune survie au-delà ;
- la pièce du 18 juin écrit « *Approvals carry across the session* » — à travers **la** session (persistance intra-session), non à travers **les** sessions — et « *Grant permission for a single tool or a whole category, and the agent won't ask again* », formulation compatible avec la seule session courante.

Ce qui persiste aux redémarrages, selon les pièces, ce sont les **règles configurées** (politiques par outil et par catégorie) — un acte de configuration explicite, non un octroi accordé en cours de session. Fait négatif **établi** sur les deux pièces citées par le journal, non *vérifié* au sens fort. **Piste de vérification avant de trancher** : la documentation officielle des permissions du Harness et le code source `@mastra/core` (dépôt public `mastra-ai/mastra`) — si une promotion d'octroi en règle persistante y existe, l'énoncé du journal se re-source ; sinon, il se corrige dans la passe suivante (un journal publié ne se réécrit pas).

### Conséquences sur les trois atterrissages du journal

- **(a) Ch. 25-26 — tient.** La chaîne confirmée est bien une réalisation concrète et datée d'un *frame* opérationnel — à instruire comme **cas**, jamais comme fondement de la taxonomie OO1-OO4 (lacune PRD Vol. II §10.10 inchangée).
- **(b) Ch. 18 — à reformuler.** L'argument « octroi par catégorie survivant à la session = élargissement de mandat sans acte de délégation » ne survit pas tel quel : ce qui survit aux sessions est une **règle**, posée par un acte explicite. La forme instruite de l'angle mort est plus fine : une règle par **catégorie** (« *allow every read without prompting* ») est un acte de délégation explicite mais **grossier** — un seul acte couvre une classe ouverte d'outils futurs, sans énumération de ce qu'il autorisera demain. L'angle mort du ch. 18 demeure, déplacé de « sans acte » vers « acte sans granularité ».
- **(c) Ch. 29 / ch. 45 — tient et se renforce.** Le mode YOLO auto-approuve tout ce qui n'est pas refusé par outil (verbatim au maillon 2) : le distinguer d'une approbation effective demeure une condition de validité des indicateurs de supervision du ch. 45, et il n'est pas un contrôle au sens de la supervision **attendue** par E-23.

### Tri épistémique

- Chaîne à sept maillons, mode YOLO, version `1.45.0`, dates et auteur : **[B]** (pièces d'éditeur datées, extraction citée).
- Persistance des octrois entre sessions : **non confirmée** — ne pas intégrer en l'état.
- Reformulation de l'atterrissage ch. 18 : **construction d'auteur**, déclarée telle.

## Annexe — entrée de journal prête à intégrer (si l'auteur consomme la fiche)

> **Instruction de la relève 2 (date d'intégration à poser).** Les deux pièces Mastra (5 et 18 juin 2026) ont été extraites le 22 juillet 2026. La chaîne de résolution à sept maillons est confirmée verbatim (code de la pièce du 5 juin) — l'ordre du journal v0.10 était exact. **Écart consigné** : « les octrois persistent d'une session à l'autre » n'est soutenu par aucune des deux pièces — ce sont les *règles* configurées qui survivent, les octrois sont de session ; l'atterrissage ch. 18 se reformule (« acte de délégation sans granularité » plutôt que « élargissement sans acte »), sous réserve d'une vérification du code source `@mastra/core`. Les atterrissages ch. 25-26 et ch. 29/45 tiennent. La relève 2 passe de « à instruire » à « instruite, [B], avec un énoncé corrigé ». Fiche : `proposition-apres-agentique-2026-07-22-03.md`.

Note de réalignement `Conspectus.md` : même nuance qu'à l'itération 02 — le cardinal « huit relèves » ne change pas (re-mesuré) ; leur statut se ventile (deux instruites [B], six [C] à instruire).

## Vérification adverse (exécutée avant livraison)

- **Réfutation tentée et productive dans les deux sens** : la première extraction (ordre 5-6 inversé) réfutée par le verbatim ; l'énoncé de persistance du journal réfuté par les deux pièces — consigné, non lissé.
- **Cardinaux** : « huit relèves » re-mesuré, inchangé ; aucun décompte du dépôt modifié.
- **Renvois** : ch. 18, 25, 26, 29, 45 cités — tous dans 1-57.
- **Zones gelées** : aucune touchée ; Livre X terminal ; les journaux v0.8-v0.10 ne sont pas réécrits — l'écart du journal v0.10 est consigné pour la passe suivante, conformément au protocole du dépôt.
- **Décision 8** : aucune thèse réécrite.

## Prochaine itération envisagée

Instruire la relève v0.10 nᵒ 6 (mémoire observationnelle — ch. 5, ch. 44) : les mêmes deux pièces la documentent, et l'extraction de cette passe a déjà relevé l'observateur, le réflecteur, le seuil de déclenchement (40 000 jetons, intervalles de 20 %), les déclencheurs d'inactivité par fournisseur et la structure à trois couches — il restera à confronter ces éléments au libellé du journal et à qualifier la réserve « aucune propriété de conservation publiquement établie ».
