# Interopérabilité et autonomie agentiques en services financiers — un diptyque

Deux monographies d'André-Guy Bruneau sur l'orchestration et l'interopérabilité des
agents en services financiers. Conçues séparément, elles se complètent : ce dépôt les
articule explicitement en **diptyque** — un cadre général mondial et théorique, puis son
instruction approfondie sur le cas canadien réglementé.

> **Verdict de complémentarité :** franchement complémentaires, mais recouvrement concentré
> sur le matériau le plus saillant et complémentarité longtemps *latente* (les deux ouvrages
> ne se citaient pas). Le présent dépôt convertit cette latence en articulation explicite.
> Évaluation détaillée : [`eval.md`](eval.md) (complémentarité nette **64/100** ;
> qualité Vol. I **80/100**, Vol. II **77/100**).

## Les deux volumes

| | **Vol. I — Interopérabilité agentique** | **Vol. II — L'autonomie encadrée** |
|---|---|---|
| **Fichier** | [`vol-1-interoperabilite-agentique/`](vol-1-interoperabilite-agentique/interoperabilite-agentique-services-financiers.md) | [`vol-2-autonomie-encadree/`](vol-2-autonomie-encadree/autonomie-encadree.md) |
| **Titre** | Interopérabilité agentique en entreprise dans le domaine des services financiers | L'autonomie encadrée |
| **Rôle** | Cadre général mondial et théorique | Cas canadien réglementé, approfondi et durci |
| **Portée** | Mondiale (UE / US / UK / Asie) | Canada-Québec (E-23, AMF, Loi 25, ACVM, Lynx/RTR) |
| **Thèse** | « Autonomie graduée sous contrôle de finalité » | « Autonomie encadrée » (*framed autonomy*, paradigme APM) |
| **Méthode-signature** | Formalisme d'ingénierie (ArchiMate 4, ADS « Boréalis ») | Rigueur évidentiaire (socle F-01…F-48, niveaux [A]/[B]/[C]) |
| **Gel de l'information** | Juin 2026 | 16-17 juillet 2026 |

## Ordre de lecture

Lire **Vol. I d'abord** (il pose la théorie de l'interopérabilité, l'ingénierie des agents,
l'anatomie des protocoles et la sécurité de la couche que Vol. II présuppose sans les
reconstruire), puis **Vol. II** (qui restreint le périmètre au Canada réglementé et instruit
le droit applicable au grain de l'architecture d'exécution). Un lecteur pressé côté canadien
peut entrer directement par le ch. 13 de Vol. II (« le pont : des contraintes réglementaires
aux frames déterministes »), en revenant à Vol. I au besoin.

## Renvois entre les deux volumes

- **Vol. II présuppose Vol. I** pour : théorie du découplage/contrat/évolution, ingénierie
  des agents LLM, anatomie complète des protocoles (MCP/A2A/AP2/AGNTCY/x402), sécurité de la
  couche agentique, sémantique/ontologies, cryptographie post-quantique.
- **Vol. I illustre mondialement** ce que **Vol. II instruit au grain du droit canadien** :
  double-qualification agent = modèle + tiers TIC, standards de données, adoption
  institutionnelle.

## Faits partagés et divergences

Les faits datés cités par les deux volumes (protocoles, réglementation) sont regroupés dans
[`commun/faits-partages.md`](commun/faits-partages.md) — **source unique de vérité** : chaque
volume devrait y renvoyer plutôt que redupliquer la donnée. Ce fichier recense aussi **deux
divergences factuelles à trancher** (date de finalisation de la ligne directrice IA de l'AMF ;
gouvernance d'AP2) — voir [`eval.md`](eval.md) §3.5.

## Structure du dépôt

```
.
├── README.md                                 ← avant-propos croisé (ce fichier)
├── eval.md                                   ← évaluation académique de la complémentarité
├── commun/
│   └── faits-partages.md                     ← faits datés partagés (source unique) + divergences
├── vol-1-interoperabilite-agentique/
│   ├── interoperabilite-agentique-services-financiers.md
│   └── interoperabilite-agentique-services-financiers.pdf
└── vol-2-autonomie-encadree/
    ├── autonomie-encadree.md
    └── autonomie-encadree.pdf
```
