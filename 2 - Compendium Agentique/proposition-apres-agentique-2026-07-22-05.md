# Proposition de passe — après-agentique, itération 05

| Champ | Valeur |
|---|---|
| Date | 22 juillet 2026 |
| Itération | 05 — relève v0.10 nᵒ 5 instruite à ses sources primaires |
| Objet | Fiche d'instruction : l'extension déclarative et les campagnes de skills malveillants — la source est identifiée, deux écarts du journal consignés, l'échelle réévaluée |
| Sièges | Vol. IV `TOC.md` v0.10 — ch. 52 (Livre IX, provenance), ch. 21 (Livre III, rug-pull/révocation), ch. 20 (Livre III, taxonomie) ; journal v0.10, relève 5 |
| Statut | **Proposition** — rien n'est intégré, rien n'est poussé ; l'intégration est une décision d'auteur |
| `TOC.md` touché | **Non par cette passe** — entrée de journal prête à intégrer en annexe |

## L'incrément

**La relève v0.10 nᵒ 5 — la seule dont aucune source primaire n'était nommée au dépôt (« troisième degré d'absence ») — est instruite : l'éditeur de sécurité est identifié (Koi Security), sa publication primaire est extraite, et le constat en ressort transformé.** Le libellé du journal parlait d'« une extension tierce » ; les pièces documentent une campagne de chaîne d'approvisionnement à l'échelle du magasin d'extensions entier. Deux éléments du libellé ne survivent pas tels quels et sont consignés.

### Sources (consultées le 22 juillet 2026)

1. **[B]** O. Yomtov (Koi Security), *ClawHavoc: 341 malicious ClawHub skills…*, billet de recherche, **1ᵉʳ février 2026**, mis à jour le **16 février 2026**. https://www.koi.ai/blog/clawhavoc-341-malicious-clawedbot-skills-found-by-the-bot-they-were-targeting — éditeur de sécurité vendant un outil d'analyse (Clawdex) : chiffres auto-déclarés, attribués.
2. **[B]** S. Bellary Seetharam, N. Mohamed, B. Melicher, O. Starov (Unit 42, Palo Alto Networks), *OpenClaw's Skill Marketplace and the Emerging AI Supply Chain Threat*, **23 juin 2026**. https://unit42.paloaltonetworks.com/openclaw-ai-supply-chain-risk/
3. **[B]** The Hacker News, *Researchers Find 341 Malicious ClawHub Skills…*, **2 février 2026** — presse spécialisée, utile pour l'attribution croisée (Koi Security ; OpenSourceMalware, P. McCarty). https://thehackernews.com/2026/02/researchers-find-341-malicious-clawhub.html

### Constat 1 — [B] : le mécanisme d'installation déclarative, confirmé et précisé

La pièce 1 documente l'installation d'extensions (*skills*) par simple déclaration, sans compilation ni déploiement : documentation d'apparence légitime portant une section « *Prerequisites* » qui fait exécuter des scripts malveillants, typosquattage de paquets populaires, **absence de vérification du magasin avant installation**. La pièce 2 porte le verbatim le plus fort pour le ch. 52 : « *The lack of isolation between skill logic and agent authority means that installation results in complete control over the agent's identity* » — l'installation d'une extension équivaut à la prise de contrôle de l'identité de l'agent. Le constat du journal (« une nomenclature logicielle d'artefact ne les capte pas ») est cohérent avec les deux pièces ; aucune des deux ne mentionne de signature d'artefact ni de nomenclature.

### Constat 2 — [B] : l'échelle — le libellé du journal sous-estimait d'un facteur de plusieurs centaines

Le journal évoquait « une extension tierce ». Les pièces documentent : **341 skills malveillants sur 2 857 audités** (Koi, 1ᵉʳ février), portés à **824 sur ≈ 10 700** à la mise à jour du 16 février (+142 %) ; charges utiles : voleur d'information AMOS (mots de passe de trousseau, navigateurs, plus de 60 familles de portefeuilles, sessions Telegram, clés SSH), interpréteur inversé (*reverse shell*) caché « *among hundreds of lines of working code* », exfiltration de `~/.clawdbot/.env` vers webhook.site. La pièce 2 ajoute cinq skills passés au travers des défenses de février à mai 2026, **après** l'intégration de VirusTotal et de ClawScan. L'écart d'échelle **renforce** les trois atterrissages — c'est une classe d'attaque systémique, pas un incident isolé.

### Constat 3 — deux écarts du journal, consignés

1. **Les dates « 28-29 janvier 2026 » ne s'épinglent pas à la source primaire.** La publication primaire de Koi est datée du 1ᵉʳ février 2026 et, telle qu'extraite, ne donne pas de dates de détection fin janvier ; des reprises de presse évoquent un premier skill malveillant le 27 janvier et une vague le 31, mais aucune pièce primaire consultée ne porte « 28-29 janvier ». L'ancrage daté défendable est : **campagne détectée fin janvier, publiée le 1ᵉʳ février 2026**. Piste : si l'auteur tient aux 28-29 janvier, la source qui les porte reste à produire.
2. **Le composé « exfiltration de données et injection d'invite » se scinde.** La vague de fin janvier (pièce 1) documente l'exfiltration et des portes dérobées dans du code — **pas d'injection d'invite explicite**. C'est la pièce 2 (juin 2026) qui documente les attaques de la couche agentique — « *semantic instruction hijacking* », détournement d'affiliation à l'exécution — sur d'autres skills, découverts de février à mai. L'énoncé instruit est donc : *deux découvertes datées distinctes*, janvier (exfiltration par code) et février-mai (détournement par instructions), non un incident unique cumulant les deux.

### Constat 4 — [B] : la révocation, absente — le cas documente l'asymétrie du ch. 21

Aucune des pièces ne relève de mécanisme de révocation automatique : le retrait dépend de l'équipe du magasin (« *We've reported our findings to ClawHub's security team* » ; Unit 42 : comptes bannis et skills supprimés après signalement). L'éditeur de sécurité a construit son propre outil d'analyse préventive et rétroactive (Clawdex) — un tiers comblant l'absence de mécanisme natif. C'est, au grain de l'extension, exactement l'asymétrie émission soignée / révocation négligée que le ch. 21 prend pour objet.

### Conséquences sur les trois atterrissages du journal

- **Ch. 52 — tient, renforcé.** Le front « provenance » gagne une campagne documentée et un verbatim d'architecture (installation = contrôle de l'identité de l'agent). La réserve du journal demeure : cela ne prouve pas le jugement « front le plus mûr ».
- **Ch. 21 — tient, renforcé.** Rug-pull au grain de l'extension, masqué dans du code fonctionnel, sans révocation native ; l'histoire des PKI se rejoue au magasin d'extensions.
- **Ch. 20 — tient, avec l'échelle en plus.** Une classe d'attaque dont le vecteur est le harnais (mécanisme d'extension), ni le protocole ni le mécanisme d'identité — elle entre au dénombrement exigé par la décision 8 **contre** la thèse d'absorption identitaire, et son poids numérique (centaines de skills) rend le dénombrement d'autant moins contournable.

### Tri épistémique

- Mécanisme d'installation, absence de vérification, absence de révocation native, verbatims : **[B]** (recherche d'éditeurs de sécurité, datée, extraite).
- Chiffres (341/2 857 ; 824/10 700 ; +142 % ; 5 skills post-défenses) : **[B] auto-déclarés par les éditeurs de sécurité, attribués** — deux éditeurs concordants sur l'ordre de grandeur, mais pas de décompte indépendant.
- Dates fines de fin janvier (27, 31) : **[C]** — presse seulement, non retrouvées dans les pièces primaires extraites.
- La scission du composé exfiltration/injection : résultat d'instruction, opposable aux deux pièces citées.

## Annexe — entrée de journal prête à intégrer (si l'auteur consomme la fiche)

> **Instruction de la relève 5 (date d'intégration à poser).** Sources primaires identifiées et extraites le 22 juillet 2026 : Koi Security (ClawHavoc, 1ᵉʳ février 2026, mise à jour 16 février) et Unit 42 (23 juin 2026). Le mécanisme d'installation déclarative, l'absence de vérification et l'absence de révocation native sont confirmés ; l'échelle passe d'« une extension » à une campagne (341 puis 824 skills malveillants, chiffres d'éditeurs attribués). **Deux écarts consignés** : les dates « 28-29 janvier » ne s'épinglent pas à la source primaire (ancrage : 1ᵉʳ février) ; le composé « exfiltration + injection d'invite » se scinde en deux découvertes datées (janvier : exfiltration par code ; février-mai : *semantic instruction hijacking*, Unit 42). Les trois atterrissages (ch. 52, 21, 20) tiennent et se renforcent. La relève 5 passe de « à instruire, troisième degré d'absence » à « instruite, [B] ». Fiche : `proposition-apres-agentique-2026-07-22-05.md`.

Note de réalignement `Conspectus.md` : cardinal « huit relèves » re-mesuré, inchangé ; statut ventilé — quatre instruites [B], quatre [C] à instruire.

## Vérification adverse (exécutée avant livraison)

- **Réfutation tentée** : le libellé du journal confronté à trois pièces indépendantes entre elles (deux éditeurs de sécurité concurrents, une presse spécialisée) — deux éléments réfutés et consignés (dates, composé), deux confirmés et renforcés (mécanisme, absence de révocation).
- **Chiffres** : relevés sur les publications des éditeurs eux-mêmes, pas sur la presse ; mise à jour du 16 février distinguée du rapport initial.
- **Cardinaux** : « huit relèves » re-mesuré, inchangé ; aucun décompte du dépôt modifié.
- **Renvois** : ch. 20, 21, 52 cités — tous dans 1-57.
- **Zones gelées** : aucune touchée ; Livre X terminal ; journaux non réécrits.
- **Décision 8** : la thèse du ch. 20 n'est pas réécrite — le dénombrement qu'elle exige gagne une pièce, contre elle, comme le journal v0.10 l'avait déjà déclaré.

## Prochaine itération envisagée

Instruire la relève v0.10 nᵒ 3 (la cinquième horloge — le harnais versionné par son éditeur, ch. 53) : les pièces Mastra déjà extraites portent la version d'éditeur (`@mastra/core@1.45.0`) ; il reste à documenter que le changement de version du harnais modifie le comportement observable à modèle, invites, outils et politique constants — le journal des versions (*releases*) du dépôt public `mastra-ai/mastra` est la pièce candidate.
