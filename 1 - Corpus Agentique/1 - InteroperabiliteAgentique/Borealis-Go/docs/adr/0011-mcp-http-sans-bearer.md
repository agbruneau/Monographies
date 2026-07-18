# ADR-0011 — MCP HTTP sans bearer : frontière réseau assumée, pas d'auth theater

- **Statut :** accepté
- **Date :** 2026-07-08
- **Phase :** post-M4.3 (audit de refactorisation)
- **Nature :** ADR d'**écart** (option alternative de l'audit, documentée plutôt que codée)

## Contexte

L'audit du 2026-07-08 (m8) relève que le mode MCP HTTP streamable
(`MCP_HTTP_ADDR`, `internal/servercmd`) ne porte **aucune authentification**,
contrairement aux agents A2A qui exigent un bearer JWT à audience restreinte
(FR-26, `internal/a2aserver.SecureMux`). L'option par défaut envisagée était
d'appliquer la même garde bearer.

Vérification avant d'agir (inspection du dépôt) :

1. **Aucun appelant en dépôt n'atteint les serveurs MCP HTTP du compose.** Les
   `cmd/agent-*` construisent leur handler via `agentcmd.handlerFor`, qui sert
   un `a2aserver.ReplyExecutor` — un **stub** de réponse, jamais un exécuteur
   qui appelle réellement `verify_identity`/`get_credit_report`/etc. La
   logique métier réelle (`internal/agent`, `internal/orchestrator`) appelle
   les serveurs MCP via des **transports en mémoire** (tests, e2e), jamais via
   `MCP_HTTP_ADDR`.
2. **La frontière réseau est déjà un périmètre de confiance explicite, pas une
   hypothèse implicite** : `deploy/docker-compose.yml` place les 4 serveurs
   MCP sur le réseau `mesh`, déclaré `internal: true` (0 egress, NFR-06), et
   **ne publie aucun port hôte** pour ces services — seul Jaeger (UI
   d'observabilité) est exposé via le réseau `edge` séparé.

Ajouter un bearer maintenant reviendrait à poser une garde qu'**aucun appelant
du dépôt ne peut satisfaire** (pas de couture d'émission/attache de jeton pour
un client MCP-sur-HTTP) : de l'authentification décorative, testable
seulement par un client fictif inventé pour l'occasion, sans rapport avec un
flux réel du démonstrateur.

## Décision (repli)

Ne pas coder de garde bearer sur `internal/servercmd` pour l'instant.
S'appuyer sur la frontière réseau **déjà posée et documentée** du
docker-compose (`mesh: internal: true`, aucun port publié) comme périmètre de
confiance assumé pour `MCP_HTTP_ADDR`. Un commentaire dans `servercmd.go`
renvoie à cette ADR.

## Conséquences

- Aucun changement de code fonctionnel ; la limite de taille de corps (m7,
  `limitBody`) reste la seule garde active sur ce chemin, et suffit contre le
  DoS mémoire trivial qui motivait la remontée en premier lieu.
- L'écart est **assumé et signalé** (règle « jamais de déviation silencieuse »),
  pas un oubli : l'option bearer reste faisable (le motif `limitBody` /
  `a2aserver.BearerAuth` est directement réutilisable) le jour où elle a un
  appelant réel à protéger.
- Un lecteur qui déploierait ces serveurs **hors** du compose fourni (port
  publié manuellement, réseau non `internal`) perdrait cette garantie sans
  avertissement au runtime — limite assumée d'un démonstrateur, pas d'un
  système de production.

## Reversal condition

Dès qu'un exécuteur A2A réel (remplaçant `ReplyExecutor`) appelle effectivement
les serveurs MCP **via `MCP_HTTP_ADDR`** (et non plus in-memory), ou que le
compose publie un port hôte pour l'un de ces services : appliquer la même
garde bearer que FR-26 (`a2aserver.BearerAuth`/`idpmock`), à audience dérivée
du nom du serveur MCP.

## Alternatives écartées

- **Bearer générique à audience partagée « mcp » sans appelant réel** : théâtre
  de sécurité — testable uniquement contre un client jouet, aucune preuve que
  le mécanisme correspond à un flux réel ; risque de fausse assurance.
- **Publier les ports en compose pour permettre un test end-to-end du bearer** :
  changerait la topologie réseau documentée (NFR-06) pour les besoins d'un
  test, contraire à la chirurgie minimale attendue par l'audit.
