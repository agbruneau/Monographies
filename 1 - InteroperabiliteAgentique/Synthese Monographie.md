```{=typst}
#set page(numbering: "1")
#v(1.0cm)
#align(center)[
  #text(size: 18pt, weight: "bold")[L'interopérabilité agentique : une revue de synthèse]

  #v(0.35cm)
  #text(size: 13pt, weight: "bold")[Des fondements de l'intégration d'entreprise aux protocoles, à la finance régulée et à l'horizon 2027-2032]

  #v(0.7cm)
  #text(size: 12pt, weight: "bold")[André-Guy Bruneau, M.Sc. IT]

  #v(0.15cm)
  #text(size: 10pt)[agbruneau\@gmail.com]

  #v(0.15cm)
  #text(size: 10pt, style: "italic")[Article de synthèse · juin 2026]
]
#v(0.55cm)
#line(length: 100%, stroke: 0.4pt + luma(150))
#v(0.35cm)
```

**Résumé.** Cet article propose une revue de synthèse de l'*interopérabilité agentique*, c'est-à-dire la capacité d'agents logiciels fondés sur de grands modèles de langage à coopérer de manière fiable, gouvernée et évolutive. Il consolide un corpus traité en spirale, du général au spécifique, sur sept strates : la théorie de l'interopérabilité des systèmes d'information, l'ingénierie des systèmes agentiques, la couche protocolaire émergente (Model Context Protocol, Agent2Agent, et initiatives concurrentes de découverte, d'identité et de paiement), son déploiement à l'échelle de l'entreprise, sa spécialisation au domaine financier régulé, sa formalisation en architecture d'entreprise ArchiMate, puis sa projection sur l'horizon 2027-2032. La thèse défendue est qu'un invariant unique — *découplage, contrat, évolution* — structure toutes ces strates : la couche agentique n'invente pas une nouvelle théorie de l'interopérabilité, elle réinstancie l'ancienne à un niveau d'abstraction supérieur, où les tensions irréductibles d'hétérogénéité, d'autonomie et de distribution resurgissent sous des mécanismes nouveaux (découverte, sémantique partagée, identité non humaine, gouvernance des frontières). L'article dégage cinq apports : une lecture unifiée raccordant des littératures disjointes ; une grille d'analyse des protocoles par les axes classiques de l'interopérabilité ; un patron de conception sectoriel — l'*autonomie graduée sous contrôle de finalité* ; une formalisation architecturale opposable ; et une validation appliquée qui éprouve l'ensemble par construction. Le socle documentaire est arrêté à juin 2026 ; les énoncés sur l'avenir sont triés entre **PROGRAMMÉ**, **PROJETÉ** et **SPÉCULATIF**. La revue s'adresse simultanément à un public de recherche et de praticiens-architectes.

**Mots-clés —** interopérabilité ; IA agentique ; Model Context Protocol (MCP) ; Agent2Agent (A2A) ; intégration d'entreprise ; services financiers ; ArchiMate ; identité non humaine ; gouvernance de l'IA.

## 1. Introduction

### 1.1 Contexte : d'une question ancienne sur un terrain neuf

L'essor des systèmes d'intelligence artificielle dits *agentiques* — des logiciels capables de poursuivre un but en planifiant, en invoquant des outils et en agissant sur leur environnement, plutôt que de seulement répondre à une requête — déplace une question ancienne du génie informatique vers un terrain neuf. Cette question est celle de l'**interopérabilité** : comment des systèmes conçus séparément, hétérogènes, autonomes et distribués, parviennent-ils non seulement à échanger de l'information, mais à l'*utiliser* mutuellement de manière fiable et durable ? Longtemps cantonnée à l'intégration d'applications d'entreprise — bus de services, files de messages, interfaces de programmation, schémas de données —, l'interopérabilité se rejoue aujourd'hui entre agents logiciels qui négocient, délèguent et transigent en langage quasi naturel.

La proposition cardinale de cette revue est que cette transposition n'abolit pas les acquis de l'intégration d'entreprise : elle les *réinstancie* à un niveau d'abstraction supérieur, où les mêmes tensions fondamentales — hétérogénéité, autonomie, distribution —, identifiées de longue date dans la littérature des bases de données fédérées (Sheth et Larson, 1990), resurgissent sous des mécanismes nouveaux. La distinction théorique fondatrice, due à Wegner (1996), entre *échanger* de l'information et savoir *l'utiliser*, demeure le pivot conceptuel : un agent peut recevoir un message syntaxiquement correct d'un autre agent sans en partager le sens pragmatique, exactement comme deux applications d'entreprise peuvent dialoguer sur le fil sans accord sémantique. Ce qui change n'est pas la nature du problème, mais l'identité et le comportement des parties : des acteurs probabilistes, non déterministes et auto-descriptifs, qui se découvrent et se coordonnent à l'exécution.

Cette difficulté a un coût économique mesurable. L'intégration point-à-point de N systèmes croît en nombre de connexions selon une loi quadratique — de l'ordre de N(N−1)/2 liens —, de sorte qu'au-delà de quelques dizaines de systèmes, la *dette d'intégration* devient le principal poste de complexité et d'immobilisation des systèmes d'information. C'est cette combinatoire que les architectures d'intégration (section 3.6) puis, à un niveau d'abstraction supérieur, les protocoles agentiques (section 5) cherchent à briser, en substituant à l'appariement direct une médiation contractuelle qui ramène le coût marginal d'un système supplémentaire d'un facteur quadratique à un facteur additif. La proposition N×M → N+M, centrale au passage à l'échelle de l'entreprise (section 6.2), n'est que la transposition agentique de ce même principe d'économie combinatoire ; et l'arbitrage *faire ou acheter*, comme le coût total de possession, s'y rejouent à l'identique. Reconnaître cette continuité évite de présenter comme une découverte ce qui est, en réalité, une nouvelle instance d'un arbitrage que l'ingénierie de l'intégration formule depuis les premiers bus de messages : c'est précisément le rôle de l'invariant que de rendre cette filiation visible.

### 1.2 Problématique

À partir de 2024-2025, une famille de protocoles d'interopérabilité agentique a émergé — le *Model Context Protocol* (MCP) pour relier un agent à ses outils et à ses données, le protocole *Agent2Agent* (A2A) pour la collaboration inter-agents, et plusieurs initiatives concurrentes de découverte, d'identité et de paiement. Ce foisonnement de spécifications a suscité une situation paradoxale : la connectivité technique progresse rapidement, mais la cohérence d'ensemble, la gouvernance, la sécurité et la conformité à l'échelle de l'entreprise restent mal cernées. Le praticien-architecte fait face à une pile encore instable, faiblement gouvernée. Le chercheur dispose de cadres théoriques éprouvés pour l'interopérabilité classique — niveaux conceptuels (dont le LCIM, modèle de diagnostic plutôt que de mesure), formalismes de compatibilité, théorie du couplage — mais peine à les raccorder aux objets agentiques, où le contrat comportemental cesse d'être déterministe.

Il manque une lecture unifiée qui relie la théorie de l'interopérabilité, l'ingénierie des systèmes agentiques, la couche protocolaire émergente, son déploiement en entreprise, sa spécialisation sectorielle et sa formalisation architecturale, puis qui projette cet ensemble sur un horizon prospectif maîtrisé. Combler ce manque est l'objet de la présente revue.

### 1.3 Question de recherche et thèse

La question directrice se formule ainsi : *quelles propriétés architecturales permettent à des agents d'intelligence artificielle d'interopérer de manière sûre, gouvernée et évolutive au sein de l'entreprise, en particulier dans le domaine exigeant des services financiers ?* La réponse défendue prend la forme d'un **invariant transversal** : l'interopérabilité durable, à toutes les couches et quelle que soit la nature des parties en présence, repose sur trois notions solidaires — auxquelles la validation appliquée à l'échelle du parc d'agents (section 10) ajoutera un quatrième terme, l'*exploitation* (voir §11.1).

Le **découplage** réduit les dépendances mutuelles entre parties ; il hérite directement de la décomposition modulaire et du masquage de l'information (Parnas, 1972). Le **contrat** est l'interface explicite qui rend le découplage opérationnel : il fixe ce sur quoi les parties s'engagent — dans son ordre syntaxique, sémantique ou comportemental (cf. §3.3) — tout en libérant tout le reste à l'évolution indépendante. L'**évolution** est la capacité du contrat à changer sans rompre la coopération. Cet invariant, éprouvé du transport d'octets jusqu'à la négociation entre agents, prend au niveau agentique la forme de la *découverte*, de la *sémantique partagée*, de l'*identité* et de la *gouvernance des frontières*. Le résultat principal de la revue est la continuité de cet invariant : ce qui change d'une strate à l'autre n'est pas le principe, mais ses réalisations — le contrat prend successivement la forme d'une description d'interface, d'un schéma de message, d'un schéma d'outil, d'une carte d'agent, d'une politique d'autorisation, puis d'un élément d'architecture formalisé.

### 1.4 Contributions

La revue dégage et organise cinq contributions issues du corpus analysé.

D'abord, une **lecture unifiée** qui raccorde, sous un invariant unique, des littératures jusqu'ici disjointes : l'intégration d'entreprise, l'ingénierie des systèmes agentiques et la couche protocolaire émergente. Ensuite, une **grille d'analyse** des protocoles agentiques par les quatre axes éprouvés de l'interopérabilité — découverte, sémantique, identité et confiance, gouvernance des frontières (cf. §2.3) —, qui permet de comparer des spécifications hétérogènes et d'anticiper leurs zones de recouvrement et de manque. Troisièmement, un **patron de conception sectoriel** — l'*autonomie graduée sous contrôle de finalité* — qui articule degré d'autonomie de l'agent, criticité de l'opération et exigence réglementaire. Quatrièmement, une **formalisation architecturale** en ArchiMate qui rend l'ensemble opposable, communicable et révisable par les équipes d'architecture d'entreprise. Enfin, une **validation appliquée** — l'architecture détaillée de solution présentée à la section 10 — qui éprouve l'ensemble par construction et étend l'invariant à l'exploitation.

### 1.5 Organisation de l'article

L'article suit une progression en spirale qui élargit puis resserre la focale. La section 2 expose la méthode et le cadre d'analyse. La section 3 établit les fondements théoriques de l'interopérabilité des systèmes d'information. La section 4 résume l'ingénierie des systèmes agentiques. La section 5, cœur de la revue, analyse la couche d'interopérabilité agentique et ses protocoles. La section 6 traite le passage à l'échelle de l'entreprise ; la section 7, la spécialisation au domaine financier ; la section 8, la formalisation ArchiMate. La section 9 projette l'ensemble sur l'horizon 2027-2032. La section 10 présente une validation appliquée sous la forme d'une architecture détaillée de solution. La section 11 discute le résultat central, les limites et les fronts de recherche ouverts ; la section 12 conclut.

## 2. Méthode et cadre d'analyse

### 2.1 Nature de la revue et double public

Cette revue est une synthèse critique organisée en spirale plutôt qu'une revue systématique au sens des protocoles PRISMA : elle consolide un corpus doctrinal et normatif déjà constitué et vérifié, et son ambition est l'intégration conceptuelle, non le dénombrement exhaustif. Elle s'adresse simultanément à deux lectorats dont les attentes diffèrent : un public **recherche** (modèles, taxonomies, formalismes, questions ouvertes) et un public **praticien-architecte** (normes, protocoles, technologies, décisions de mise en œuvre datées). Pour les servir sans cloisonner le propos, deux dispositifs parcourent le texte. Les encadrés *Perspective recherche* isolent les apports théoriques ; les encadrés *Mise en œuvre* isolent les normes datées et les considérations de déploiement.

> **Perspective recherche.** L'unité d'analyse de la revue n'est pas le modèle de langage isolé, mais le *système* qui l'enveloppe et les *relations* entre systèmes. Conformément à Wegner (1996), l'interopérabilité est une relation entre systèmes, non une propriété intrinsèque de l'un d'eux ; elle ne se mesure donc qu'au point de contact, et se maintient sous évolution gouvernée plutôt qu'elle ne s'acquiert une fois pour toutes.

### 2.2 L'invariant transversal comme grille de lecture

L'invariant *découplage, contrat, évolution* fournit la grille de lecture unique de la revue. Il est rappelé à chaque strate, et chaque mécanisme rencontré y est rapporté. Le tableau 1 anticipe la manière dont les trois termes se déclinent à chacune des couches canoniques de l'interopérabilité, depuis la couche technique jusqu'à la couche organisationnelle et pragmatique propre aux agents — une progression que le *Levels of Conceptual Interoperability Model* (LCIM, §3.2) formalise au-delà du syntaxique et que la section 5 transpose aux agents. Cette grille n'est pas une simple commodité d'exposition : elle est l'hypothèse de travail testée tout au long de l'article, et le critère qui sépare ce qui relève de la nouveauté réelle de ce qui n'est qu'une réinstanciation d'un acquis ancien.

Table: Tableau 1 — Déclinaison de l'invariant par couche d'interopérabilité.

| Couche | Enjeu d'interopérabilité | Mécanisme de découplage | Porteur du contrat | Stratégie d'évolution |
| --- | --- | --- | --- | --- |
| Technique | Transporter des octets entre points distants | Indirection, asynchronisme, courtier | Protocole de transport (HTTP, files) | Négociation de version, compatibilité de transport |
| Syntaxique | Structurer le message de façon analysable | Schémas, formats auto-descriptifs | OpenAPI, AsyncAPI, JSON Schema | Compatibilité ascendante/descendante, lecteur tolérant |
| Sémantique | Partager le sens des termes échangés | Modèle canonique, ontologies | RDF/OWL, SHACL, vocabulaires | Médiation, alignement, versionnement de vocabulaire |
| Organisationnelle / pragmatique | Aligner buts, autorité et responsabilité | Découverte, identité, politiques | Carte d'agent, mandat, politique d'autorisation | Gouvernance, registres, régime de responsabilité |

### 2.3 Les axes d'analyse des protocoles

Pour comparer des spécifications hétérogènes, la revue mobilise quatre axes hérités de la théorie de l'interopérabilité et transposés aux agents : la **découverte** (comment une partie trouve-t-elle l'autre et sait-elle ce qu'elle sait faire ?), la **sémantique** (comment s'accordent-elles sur le sens d'une intention ou d'une capacité ?), l'**identité et la confiance** (comment s'authentifient-elles et délèguent-elles l'autorité ?) et la **gouvernance des frontières** (comment la sécurité et la conformité tiennent-elles au point de contact ?). Ces axes structurent en particulier la section 5. Un cinquième découpage, orthogonal, distingue quatre *plans d'interaction* — agent-outil, agent-agent, agent-humain et agent-données — que la revue emploie pour situer chaque protocole dans son emplacement fonctionnel.

### 2.4 Discipline épistémique et socle documentaire

La revue applique une discipline de sourçage stricte : tout énoncé factuel porte, dans le corpus dont il est tiré, une source nominative et un millésime ; les ressources vivantes — spécifications versionnées par date, textes réglementaires en cours d'adoption — sont signalées par le marqueur ⚠. Le socle documentaire est arrêté à **juin 2026**. La section 9, prospective, déroge seule à cette borne, en triant explicitement chaque énoncé sur l'avenir entre **PROGRAMMÉ** (engagement daté), **PROJETÉ** (prévision millésimée d'analyste) et **SPÉCULATIF** (pari de recherche) ; un statut *preview* ou non encore finalisé n'y est jamais présenté comme acquis. Conformément à cette discipline, les chiffres rapportés ici le sont avec leur statut épistémique : une mesure académique, une étude tierce, une projection d'analyste et un chiffre de fournisseur non audité n'ont pas la même valeur probante, et la revue ne les présente jamais sur le même plan.

### 2.5 Positionnement par rapport aux revues existantes

Plusieurs revues récentes recensent les protocoles d'interopérabilité agentique : des panoramas comparatifs des spécifications MCP, ACP, A2A et ANP (Ehtesham et coll., 2025 ; Yang et coll., 2025a), une taxonomie technique des protocoles de communication d'agents (Sander et coll., 2026) et une lecture sémantique d'un large corpus de protocoles (Yuan et coll., 2026). Ces travaux cartographient l'espace des protocoles avec rigueur, mais ils l'isolent de son amont — la théorie de l'interopérabilité d'entreprise — et de son aval — le déploiement gouverné, la spécialisation sectorielle et la formalisation architecturale. L'apport propre de la présente revue n'est pas de recenser un protocole de plus, mais de *raccorder* ces strates sous un invariant unique, et de soumettre chaque mécanisme au même test : relève-t-il d'une nouveauté réelle, ou d'une réinstanciation d'un acquis de l'intégration d'entreprise ? C'est ce raccordement, plutôt que l'exhaustivité du recensement, qui distingue la contribution. La portée est, en contrepartie, délibérément moins large sur chaque protocole pris isolément : la revue renvoie aux panoramas spécialisés pour le détail d'une spécification donnée, et concentre son effort sur les jointures.

## 3. Fondements : la théorie de l'interopérabilité des systèmes d'information

### 3.1 Définitions : échanger n'est pas utiliser

La première strate fixe le vocabulaire. L'interopérabilité se définit comme la capacité de deux systèmes ou plus à échanger de l'information *et à utiliser* l'information échangée (Wegner, 1996) — une définition qui dissocie l'échange (le transfert réussi) de l'usage (l'interprétation effective), et la *capacité* d'un système à coopérer de l'*acte* ponctuel d'interopérer. Cette double exigence la distingue de notions voisines : l'*intégration* désigne le couplage durable de composants en un tout ; la *compatibilité*, l'aptitude à coexister sans interférence ; la *portabilité*, la capacité à transférer un composant d'un environnement à un autre. La norme ISO/IEC 25010:2023, modèle de qualité du produit logiciel — dont la révision introduit la flexibilité comme caractéristique de premier rang —, range l'interopérabilité parmi les sous-caractéristiques de la compatibilité, confirmant son statut d'exigence qualité plutôt que de simple détail d'implémentation.

Cette définition emporte une conséquence méthodologique forte : l'interopérabilité n'est pas un état acquis mais une propriété à instaurer par contrat et à maintenir sous évolution gouvernée. Elle relève donc de l'ingénierie continue, non de la configuration ponctuelle.

### 3.2 Niveaux d'interopérabilité et modèle LCIM

La pile canonique distingue quatre niveaux cumulatifs : **technique** (transport), **syntaxique** (format), **sémantique** (sens) et **organisationnel** (processus et finalité). Chaque palier présuppose les précédents. Le modèle de référence le plus influent pour raisonner au-delà du syntaxique est le *Levels of Conceptual Interoperability Model* (LCIM), introduit par Tolk et Muguira (2003) puis raffiné à sept niveaux par Wang, Tolk et Wang (2009). Le LCIM superpose aux niveaux techniques des niveaux *pragmatique*, *dynamique* et *conceptuel*, articulés autour de trois propriétés : l'*integratability* (fondations techniques), l'*interoperability* (données et services) et la *composability* (modèles conceptuels). Son enseignement central est qu'un échange syntaxiquement et sémantiquement correct ne suffit pas tant que les modèles conceptuels sous-jacents ne sont pas explicités. Le LCIM est un outil de diagnostic, non une métrique : il situe une difficulté, il ne la mesure pas.

> **Perspective recherche.** Le LCIM jouera un rôle pivot à la section 5 : ses niveaux *pragmatique* et *dynamique*, longtemps théoriques, deviennent les verrous opérationnels de l'interopérabilité entre agents, car ce qui transite n'est plus une donnée stable mais une intention négociée à l'exécution.

### 3.3 Couplage, contrats et les trois tensions irréductibles

Le couplage se décline en facettes indépendantes — technologie, format, localisation, temporalité, interaction —, chacune réductible par un mécanisme propre : indirection pour la localisation, asynchronisme pour la temporalité, schémas tolérants pour le format (tableau 2). Cette décomposition guide la conception, car réduire une facette sans toucher aux autres permet un découplage *ciblé* plutôt qu'une refonte coûteuse. Le contrat, support formel du découplage, se décline lui-même en trois ordres : **syntaxique** (la structure du message), **sémantique** (le sens des termes) et **comportemental** (l'ordre admissible des interactions). Cette tripartition annonce directement la difficulté agentique, où le contrat comportemental devient probabiliste.

Trois tensions demeurent irréductibles, identifiées de longue date dans la littérature des bases de données fédérées (Sheth et Larson, 1990) : l'**hétérogénéité** (technologie, structure et sens divergents — la sémantique étant la plus difficile à réduire), l'**autonomie** (chaque partie reste maîtresse de ses décisions, ce qui interdit toute solution par autorité centrale) et la **distribution** (séparation dans l'espace et le temps, latence, pannes partielles). Ces trois axes ne disparaissent jamais ; ils appellent les mêmes arbitrages à chaque couche. La revue soutient qu'ils resurgissent, intacts, au niveau agentique : un agent peut être couplé à un pair par la localisation, le format, la temporalité ou l'interaction, et les mêmes mécanismes de découplage — registres, schémas tolérants, tâches asynchrones, contrats comportementaux — sont l'objet même de la couche d'interopérabilité agentique (section 5).

Table: Tableau 2 — Facettes du couplage et mécanismes de découplage (d'après Hohpe et Woolf, 2003).

| Facette | Ce qu'elle lie | Mécanisme de découplage |
| --- | --- | --- |
| Technologie | Pile d'implémentation des parties | Indirection, adaptateurs, standards ouverts |
| Format | Structure des données échangées | Schémas, lecteur tolérant, modèle canonique |
| Localisation | Adresse du partenaire | Annuaire, courtier, résolution dynamique |
| Temporalité | Disponibilité simultanée | Asynchronisme, file de messages |
| Interaction | Ordre admissible des échanges | Contrat comportemental, types de session |

### 3.4 Évolution et versionnement

L'évolution gouvernée distingue compatibilité ascendante (un ancien client tolère un nouveau service) et descendante (un nouveau client tolère un ancien service), et s'appuie sur la discipline du *lecteur tolérant* : n'exiger que ce qui est nécessaire, ignorer le reste. Côté recherche, des formalismes garantissent la compatibilité comportementale de façon décidable : les *types de session* (Honda, Yoshida et Carbone, 2008) et les *automates d'interface* (de Alfaro et Henzinger, 2001) modélisent les protocoles d'interaction et vérifient qu'une composition ne se bloque pas. Côté pratique, les en-têtes normalisés `Deprecation` et `Sunset`, ainsi que les registres de schémas à vérification automatique, opérationnalisent l'évolution sans rupture. L'écart entre la rigueur des formalismes et le caractère structurel des contrats d'API industriels demeure un front ouvert, sur lequel la revue reviendra.

### 3.5 Cadres de référence et modèles de maturité

Pour rendre l'interopérabilité gouvernable, des cadres en fixent les dimensions et des modèles en jaugent l'état — l'équivalent organisationnel du rôle que le contrat joue à l'échelle technique. Le cadre européen d'interopérabilité (NEIF, 2017) distingue quatre couches — juridique, organisationnelle, sémantique, technique — traversées par une gouvernance transversale ; la norme ISO 11354 (*Framework for Enterprise Interoperability*) organise un cube croisant préoccupations (données, services, processus, affaires), barrières (conceptuelles, technologiques, organisationnelles) et approches — intégrée (format commun imposé), unifiée (méta-modèle pivot) et fédérée (adaptation dynamique à l'exécution, sans format imposé), cette dernière incarnant l'idéal du couplage faible.

Les modèles de maturité par paliers — de LISI et l'OIM, d'origine militaire, à l'EIMM et à ses dérivés orientés services publics (IMM/IMAPS) — fournissent une échelle de progression commune, mais souffrent d'un statisme mal adapté aux architectures dynamiques pilotées par l'IA : un système peut satisfaire les critères formels d'un palier sans interopérer effectivement (conformité de façade), et la maturité s'y mesure mieux en continu (test de contrat, observabilité) qu'en paliers figés. Cette limite prépare la transposition du LCIM aux agents employée aux sections 6 et 8.

Le cadre européen illustre le passage d'un référentiel de principes à une architecture outillée : l'*European Interoperability Reference Architecture* (EIRA v6.1.0, 31 mai 2024) décompose l'interopérabilité en blocs architecturaux réutilisables sur les quatre couches du NEIF, avec son outillage de conformité (CarTool). Le règlement pour une Europe interopérable (Règlement (UE) 2024/903, 13 mars 2024) fait basculer cette logique en obligation juridique : il impose, depuis le 12 janvier 2025, une évaluation formelle avant toute décision touchant des systèmes de services publics transfrontaliers, et institue des solutions labellisées partagées. Cette bascule consacre, au sommet de la pile des cadres, l'invariant de la revue en exigence de conformité opposable, non plus seulement en bonne pratique recommandée.

### 3.6 Architectures d'intégration : de la SOA aux maillages

L'histoire des architectures d'intégration est celle d'un déplacement progressif de l'intelligence vers les extrémités. L'architecture orientée services (SOA) et les services web (SOAP, WSDL, UDDI) ont posé la médiation centralisée par bus de services d'entreprise (ESB), substituant à la prolifération point-à-point une médiation par contrats de service réutilisables. Les microservices ont renversé le principe — *smart endpoints, dumb pipes* —, déportant la logique d'intégration dans les services eux-mêmes et décentralisant la donnée et le déploiement. Le maillage de services (*service mesh*), avec son plan de contrôle distinct du plan de données, a externalisé les préoccupations transversales (sécurité, routage, observabilité) hors du code applicatif — distinction structurante qui fournira, transposée, le patron du *plan de contrôle agentique* (section 6.3), où la décision est séparée de l'autorité. La leçon historique de l'échec d'UDDI — un registre central découvre mal ce qui évolue vite et appartient à des organisations distinctes — informe directement la conception des registres d'agents.

Entre la SOA et le maillage, la pratique d'entreprise a sédimenté un outillage que la couche agentique réutilise plus qu'elle ne le remplace. Les plateformes d'intégration en tant que service (iPaaS) industrialisent la création de connecteurs ; l'échange de données informatisé (EDI, EDIFACT, X12) régit depuis des décennies les flux interentreprises sous contrat strict ; le maillage d'événements prolonge l'architecture événementielle au-delà d'un domaine unique. La modernisation de l'héritage suit le patron du *strangler fig* : encapsuler progressivement un système ancien derrière des interfaces contractuelles jusqu'à pouvoir le remplacer sans rupture. Ces dispositifs — connecteurs gérés, contrats interentreprises, encapsulation progressive — réapparaissent intacts à la section 6, où l'agent se greffe sur le tissu existant par des serveurs d'outils, connecteurs gérés d'un genre nouveau.

### 3.7 Styles d'API, messagerie et patrons d'intégration

Les mécanismes d'échange se déploient sur trois strates : les styles d'interface synchrone, le middleware asynchrone et le vocabulaire de patrons qui les orchestre. Côté synchrone, REST (Fielding, 2000) fait de la ressource adressable par HTTP le pivot du découplage client-serveur ; le modèle de maturité de Richardson, popularisé par Fowler, en gradue l'adoption sans se confondre avec la définition de Fielding — l'hypermédia (HATEOAS) y est une condition nécessaire, non un degré facultatif, et la pratique industrielle s'arrête souvent à un RPC discipliné sur HTTP. gRPC (HTTP/2, Protocol Buffers) privilégie la latence interservices ; GraphQL déplace le cadrage de la requête vers le client et pose la question de la composition d'un graphe unifié entre équipes autonomes — fédération ou schémas composites, encore en cours de normalisation ouverte.

Le contrat se rend machine-lisible par OpenAPI (3.1.0 sur JSON Schema 2020-12 ; 3.2.0 depuis septembre 2025, étendue au streaming) pour le synchrone et par AsyncAPI (3.0, puis 3.1) pour l'événementiel, deux formalismes cousins partageant schémas et discipline d'outillage. L'approche *contract-first*, où le contrat précède le code, sert directement l'invariant : elle permet simulation, génération parallèle producteur/consommateur et test de contrat piloté par le consommateur (Pact), détectant une rupture de compatibilité avant déploiement. La gestion de ces contrats se déplace vers la passerelle d'API (routage, authentification, limitation de débit) puis vers le plan de gestion (*API Management*), qui en gouverne le cycle de vie — catalogue, portail développeur, dépréciation — et traite l'API comme un produit.

Le découplage asynchrone repose sur le middleware orienté message (MOM) : producteur et consommateur, jamais tenus de coexister dans le temps, communiquent par un canal intermédiaire (file, publication-abonnement ou flux persistant rejouable). Trois sémantiques de livraison encadrent l'échange — au plus une fois, au moins une fois, et une livraison *exactement une fois* qui relève de l'abus de langage, la livraison fiable sous pannes étant théoriquement impossible (problème des deux généraux, Akkoyunlu et coll., 1975) ; seul un *traitement* exactement-une-fois, par idempotence côté consommateur, est atteignable. L'architecture événementielle généralise ce principe en érigeant le fait métier en primitive d'intégration ; la boîte d'envoi transactionnelle et la capture de changement (CDC) réconcilient base de données et flux d'événements sans transaction distribuée, sous un régime de cohérence éventuelle.

Le langage de patrons d'intégration (Hohpe et Woolf, 2003) — routeur fondé sur le contenu, traducteur de messages, modèle canonique, séquenceur — codifie un demi-siècle de solutions éprouvées ; la coordination de processus distribués en règle la tenue sous défaillance partielle : orchestration centralisée (visibilité, couplage de contrôle) contre chorégraphie répartie (autonomie, observabilité difficile), saga compensatoire contre validation à deux phases, sous la contrainte du théorème CAP. Ces patrons préfigurent le basculement de 2025-2026 traité à la section 5 : l'orchestration classique, déterministe et auditable par construction, cède devant l'orchestration agentique, où un agent planifie dynamiquement l'enchaînement des appels sans plan préétabli — gain d'adaptabilité payé d'une perte de reproductibilité, qu'outillent tout juste les protocoles d'invocation d'outils, au premier rang MCP.

### 3.8 Données, sémantique et formats de table

L'interopérabilité des données mobilise des formats — textuels ou binaires, orientés ligne ou colonne — et des schémas à compatibilité gouvernée (ascendante, descendante, complète, transitive), vérifiée automatiquement par des registres de schémas dans l'architecture événementielle. Le modèle canonique de données ramène un appariement quadratique entre *N* systèmes à un appariement linéaire — *N×N* transformations directes réduites à *2N* via un format pivot. Au plan analytique, les formats colonnaires (Parquet, ORC, Arrow) et les formats de table qui les surmontent — Apache Iceberg, Delta Lake — ajoutent aux fichiers une couche de métadonnées apportant transactions ACID, voyage dans le temps et évolution de schéma ; le contrat s'y déplace vers le catalogue, que l'Iceberg REST Catalog (Apache Polaris, Unity Catalog OSS) standardise en interface ouverte entre moteurs hétérogènes.

Le mouvement *data mesh* transpose au domaine des données les principes de découplage et de propriété décentralisée : le *contrat de données* exécutable — formalisé par l'Open Data Contract Standard — fixe schéma et garanties de qualité entre fournisseur et consommateurs, tandis que le *produit de données* en expose les ports de sortie. L'échange inter-organisationnel de données souveraines s'appuie en outre sur le Dataspace Protocol, qui normalise négociation de contrat et transfert entre espaces de données mutuellement méfiants.

L'interopérabilité sémantique mobilise la pile du Web sémantique (RDF, RDFS/OWL, JSON-LD ; SPARQL pour l'interrogation fédérée ; SHACL pour la validation sous hypothèse de monde clos), la médiation ontologique par appariement (campagne OAEI) et le graphe de connaissances d'entreprise. Les grands modèles de langage y introduisent une capacité nouvelle d'alignement et d'extraction, mais d'une fiabilité insuffisante pour une médiation non supervisée : la question centrale devient celle de la *confiance*, soit les conditions auxquelles un appariement généré peut être promu au rang de contrat opposable. Le mouvement inverse, où le graphe discipline le modèle, se manifeste dans GraphRAG : ancrer la génération sur un graphe structuré réduit les hallucinations et restaure la traçabilité, préfigurant la couche d'interopérabilité agentique où le sens partagé devient condition de coopération entre agents autonomes (section 5).

### 3.9 Sécurité, identité et gouvernance de l'interopérabilité

La sécurité n'est pas une couche superposée à l'interopérabilité : elle en est une propriété émergente, vérifiée à chaque échange plutôt que garantie par le périmètre réseau. La frontière de confiance se déplace vers le message, et le contrat de sécurité — schéma de jeton, profil d'autorisation, identité de charge — devient aussi explicite et versionné que le contrat d'interface. L'intégration multiplie les frontières de confiance et les classes d'attaque propres : le *confused deputy* (un service privilégié induit à agir pour le compte d'un appelant non autorisé), le vol et le rejeu de jetons porteurs, et — en tête du référentiel OWASP API Security — les défauts d'autorisation au niveau de l'objet (BOLA) ou de ses propriétés (BOPLA), fautes de contrat où l'interface expose une surface plus large que la politique ne la garde.

L'identité fédérée répond à trois questions : qui parle, à quoi a-t-il droit, comment le prouver sans secret durable. L'authentification migre de l'assertion SAML 2.0 vers OpenID Connect, adossé à l'autorisation déléguée OAuth 2.0 — sa consolidation en OAuth 2.1 (PKCE obligatoire, flux implicite retiré) demeure un projet à la mi-2026, codifiée entre-temps par la BCP 240. Le jeton (JWT) reste porteur et rejouable sauf liaison cryptographique à son détenteur — mTLS (RFC 8705) ou DPoP (RFC 9449), mécanismes combinés en profils de haute sécurité tels que FAPI 2.0 (22 février 2025), pertinent notamment en finance ouverte. Le provisionnement (SCIM 2.0) et la propagation de contexte (échange de jetons, jetons de transaction) évitent qu'une réémission de jeton n'efface l'identité du mandant d'origine.

L'architecture *zero-trust* (NIST SP 800-207) généralise le principe « never trust, always verify » : aucune confiance dérivée de la localisation réseau, chaque requête vérifiée. Pour les charges de travail, l'identité SPIFFE/SPIRE substitue une attestation cryptographique éphémère aux secrets statiques ; à l'échelle inter-organisationnelle, l'identité décentralisée — identifiants et justificatifs vérifiables du W3C, ancrés par eIDAS 2.0 en droit européen, dont la généralisation des portefeuilles est attendue dès 2026 — déplace la racine de confiance hors d'un fournisseur unique. Ce socle, repris et durci par la couche agentique (sections 5.7 et 6.4), y devient multi-saut et composite : un parc d'agents multiplie les porteurs d'identité, donc les contrats de confiance à émettre, faire tourner et révoquer. La confiance, comme l'interopérabilité, est une relation établie par contrat et maintenue sous révocation, non une propriété acquise une fois pour toutes.

Gouverner ce dispositif suppose un cycle de vie de contrat outillé : gouvernance fédérée et *policy-as-code* (linting de contrats en intégration continue), triplets SLA-SLO-SLI et budget d'erreur, test d'interopérabilité distinct de la simple conformité (Wegner, 1996) et vérifié en continu, certification tierce (OpenID Foundation, FIDO Alliance) pour l'attester à l'échelle d'un écosystème. L'observabilité — traçage distribué normalisé par W3C Trace Context, outillé par OpenTelemetry — en est la condition opérationnelle et la base de l'auditabilité, désormais exigée par la réglementation : directive PSD2 en banque ouverte, Data Act (Règlement (UE) 2023/2854, applicable depuis le 12 septembre 2025) pour la portabilité inter-nuage. Le tableau 1 (section 2) résume la déclinaison de l'invariant qui sous-tend cette strate : à chaque couche, un mécanisme découple, un contrat engage, une stratégie d'évolution absorbe le changement.

## 4. L'ingénierie des systèmes agentiques

### 4.1 Définitions et héritages

La deuxième strate établit que l'agent fondé sur un grand modèle de langage est un *système d'intégration à part entière* : il consomme des outils, mobilise une mémoire, se coordonne avec d'autres agents. La bascule décisive de 2024-2026 n'est donc pas une montée en capacité du modèle, mais un déplacement de la question d'ingénierie : de la qualité de la réponse vers la fiabilité, la contrôlabilité en coût et la défendabilité en sécurité du *système* qui enveloppe le modèle. L'agentique est une propriété du système, non un attribut du modèle.

La définition canonique de l'agent rationnel — une entité qui perçoit par des capteurs et agit par des effecteurs pour atteindre un but (Russell et Norvig, 2021) — encadre des définitions plus récentes et concurrentes, faute de consensus formel sur le terme. Pour Anthropic (2024), un agent est un système où le modèle dirige dynamiquement ses propres processus et son usage des outils — une formulation centrée sur la *direction dynamique* du processus ; Willison (2025) en propose une lecture plus opérationnelle, un agent étant un LLM qui exécute des outils en boucle jusqu'à l'atteinte d'un but, qui met l'accent sur la mécanique exécutable plutôt que sur l'intention. La distinction proposée par Sapkota et coll. (2025) entre *AI agents* (entités outillées individuelles) et *Agentic AI* (systèmes coordonnés de plusieurs agents) clarifie ce vocabulaire encore flottant. Trois héritages irriguent le domaine : l'architecture *croyances-désirs-intentions* (BDI) de Bratman (1987) et de Rao et Georgeff (1995), qui fonde l'engagement vers une intention pour éviter la reconsidération permanente ; les systèmes multi-agents classiques, avec le *Contract Net Protocol* (Smith, 1980) et les langages de communication d'agents (FIPA-ACL) ; et la lignée des processus de décision markoviens et de l'apprentissage par renforcement (Sutton et Barto, 2018).

### 4.2 La boucle agentique et les patrons d'architecture

Le squelette de tout agent est la boucle *perception-raisonnement-action-observation* : le modèle augmenté itère — il observe l'état, raisonne, choisit une action (souvent un appel d'outil), observe le résultat — jusqu'à un critère d'arrêt (but atteint, budget épuisé, échec, intervention humaine). Sur ce squelette se greffent des patrons mono-agent éprouvés : **ReAct** entrelace raisonnement et action (Yao et coll., 2022) ; **Plan-and-Execute** sépare la planification de l'exécution, au prix d'une rigidité accrue face à un environnement changeant ; **Reflexion** ajoute une auto-critique par renforcement verbal (Shinn et coll., 2023), limitée par l'effet de chambre d'écho d'un agent qui se juge lui-même ; **CodeAct** exprime l'action sous forme de code exécutable, ce qui exige en retour un bac à sable non négociable. Un principe directeur traverse cette littérature d'ingénierie : *minimiser la surface agentique*. Agentifier est un choix qui doit se justifier, non un point de départ ; un flux stable se traite mieux par un enchaînement déterministe que par un agent. Le cadre intégrateur CoALA (Sumers et coll., 2023) organise ces éléments — modules de mémoire, espaces d'action interne et externe, procédure de décision — en une architecture cognitive de référence.

La revue esquisse ce même principe sous la forme d'un continuum d'autonomie à six niveaux (0 à 5), par analogie avec les échelles d'autonomie de conduite, du simple outil sous contrôle humain total jusqu'à l'agent pleinement autonome ; il gouverne le choix d'architecture, le placement des points d'approbation humaine et le calibrage des garde-fous. Le tableau suivant en propose une lecture opérationnelle :

Table: Tableau 3 — Continuum d'autonomie et déplacement de la tenue de la boucle.

| Niveau | Désignation | Qui tient la boucle de contrôle | Rôle de l'humain |
| --- | --- | --- | --- |
| 0-1 | Outil / assistant | Le code ou l'utilisateur, intégralement | Décide et exécute |
| 2 | Copilote | L'humain, assisté par le modèle | Valide chaque étape |
| 3 | Orchestration sous revue | Le modèle propose, le code contraint | Approuve les actions sensibles |
| 4 | Autonomie bornée | Le modèle, dans un périmètre fixé | Supervise et intervient sur exception |
| 5 | Autonomie étendue | Le modèle, de bout en bout | Audite a posteriori |

### 4.3 Raisonnement et calcul à l'inférence

La capacité de raisonnement des agents s'est déplacée du simple enchaînement linéaire (*chain-of-thought*, Wei et coll., 2022) vers la recherche structurée (arbres et graphes de pensées) et le *calcul à l'inférence* (test-time compute), où l'allocation de davantage de calcul au moment de la requête améliore la qualité — une nouvelle loi d'échelle qui déplace une part du budget de l'entraînement vers le service. Les vérificateurs — récompense de processus contre récompense de résultat — encadrent cette recherche, la supervision étape par étape s'avérant souvent supérieure à la seule supervision du résultat. Une classe distincte de **modèles de raisonnement** — la série o d'OpenAI, DeepSeek-R1, le raisonnement étendu de Claude — délibère longuement avant de répondre ; ce budget de réflexion se règle explicitement (*reasoning_effort*, *thinking_budget*) et se facture comme des jetons de sortie, sous peine de sur-réflexion (*overthinking*) coûteuse et sans gain. Une réserve théorique majeure subsiste par ailleurs : les grands modèles de langage échouent de manière fiable comme planificateurs autonomes sur des bancs formels (PlanBench), ce qui motive l'approche **LLM-Modulo** (Kambhampati et coll., 2024), où le modèle génère des candidats qu'un vérificateur externe valide — le modèle est un traducteur et un générateur d'heuristiques, non un solveur garanti. Enfin, la trace de raisonnement verbalisée n'est pas toujours fidèle au calcul réel qui produit la réponse : elle demeure un signal utile pour la surveillance de l'agent, mais non probant (cf. §4.6).

### 4.4 Outils, mémoire et ancrage informationnel

L'usage d'outils — l'*appel de fonction* — est le mécanisme par lequel l'agent agit sur le monde. Le schéma de l'outil, exprimé en JSON Schema, est un *contrat de premier ordre* : sa stabilité conditionne la non-régression de l'agent. Le Model Context Protocol (MCP), analysé en détail à la section 5, standardise ce contrat entre l'agent et le fournisseur d'outils. Un second espace d'action, plus englobant, prolonge l'appel d'outils typé : l'exécution de code en bac à sable et le *computer use* — piloter une interface graphique générale par capture d'écran, clic et saisie, ouvert en bêta publique fin 2024 — étendent l'agent au-delà d'un catalogue d'API nommées, au prix d'un rayon d'impact et d'une surface d'attaque élargis d'autant (cf. §4.6) ; leur fiabilité, mesurée sur des bancs dédiés (OSWorld), progresse sans encore atteindre le plafond de performance humaine. La mémoire se décline en mémoire de travail et mémoire à long terme (épisodique, sémantique, procédurale), et l'*ingénierie du contexte* — écrire, sélectionner, compacter, isoler l'information dans la fenêtre du modèle — devient une discipline à part entière, confrontée aux phénomènes de perte d'information en milieu de contexte. L'ancrage informationnel mûrit du RAG statique vers le **RAG agentique** : une boucle planifier-récupérer-critiquer-itérer (Self-RAG, Corrective RAG) où l'agent décide quand et quoi récupérer, et évalue lui-même la pertinence de ce qu'il a obtenu.

### 4.5 Systèmes multi-agents et orchestration

Au-delà de l'agent unique, les systèmes multi-agents organisent des topologies (étoile, chaîne, arbre, graphe) et des patrons (orchestrateur-travailleurs, débat). Leur justification doit être pesée : Anthropic rapporte, pour son système de recherche orchestrateur-travailleurs, un gain d'environ +90 % sur une évaluation interne par rapport à un agent unique, mais au prix d'une consommation de jetons de l'ordre de quinze fois celle d'une conversation simple (Anthropic, 2025) — une mesure ponctuelle, non une loi générale, dont le gain doit être mis en balance avec ce coût à chaque décision d'architecture. La coordination inter-agents prépare la couche d'interopérabilité proprement dite, traitée à la section 5. La taxonomie MAST (Cemri et coll., 2025) classe les modes de défaillance propres aux systèmes multi-agents en trois catégories — spécification insuffisante, désalignement inter-agents, vérification défaillante — et établit que ces défaillances sont *émergentes*, non réductibles à l'agent isolé.

### 4.6 Évaluation, sécurité et opérationnalisation

L'évaluation des agents est intrinsèquement difficile : on évalue un système, non un modèle, dans un contexte non déterministe. La métrique de potentiel (*pass@k*) masque l'effondrement de la consistance que révèle la métrique de fiabilité (*pass^k*), et la probabilité d'échec par pas se compose géométriquement sur les tâches à horizon long. La sécurité, elle, est une propriété structurelle et non un retard de capacité. La *triade létale* (Willison, 2025) — la réunion, dans un même agent, de données privées, de contenu non fiable et d'un canal de sortie — caractérise la vulnérabilité fondamentale ; l'injection d'invite indirecte (Greshake et coll., 2023) n'est pas résoluble au niveau du modèle, car instructions et données partagent le même flux de jetons. Les défenses connues — patrons *Dual LLM*, *CaMeL* (Debenedetti et coll., 2025), *Règle de deux* — réduisent le risque sans l'éliminer ; CaMeL, par exemple, neutralise par conception environ les deux tiers des scénarios d'un banc adversarial de référence (AgentDojo), la fraction restante exigeant une politique que le mécanisme ne sait pas exprimer — proportion indicative qui interdit toute lecture de « solution ». Au-delà des attaques externes, l'agent peut dévier de lui-même : le *reward hacking* — contourner la lettre d'une métrique sans en accomplir l'intention — et le *scheming* — dissimuler stratégiquement ses objectifs réels à ses superviseurs — ont été documentés sur des modèles de frontière par des équipes de recherche indépendantes (METR, 2025 ; Meinke et coll., 2024 ; Greenblatt et coll., 2024), étendant l'enjeu d'alignement au-delà de la seule robustesse aux attaques externes. La défense se conçoit donc en profondeur, jamais comme une élimination. L'opérationnalisation (AgentOps) traite l'agent en production comme un contrat de service révisable, dont les budgets, les habilitations et les modes d'approbation se versionnent comme du code.

### 4.7 Entraînement, adaptation et bancs d'essai

Deux régimes d'amélioration, qu'il importe de distinguer, coexistent. L'*adaptation en contexte* est réversible : elle procède par édition de l'invite, des outils ou de la mémoire, sans toucher aux poids du modèle ; elle relève du contrat et de l'évolution gouvernée d'un système déployé. La *mise à jour paramétrique* est persistante : elle procède par affinage supervisé, apprentissage par renforcement à récompense vérifiable ou distillation, et relève du cycle de vie du modèle lui-même. La frontière entre les deux est aussi une frontière de gouvernance — la première se journalise et se rétablit, la seconde exige une revalidation. Des familles de bancs d'essai distinctes mesurent le codage (par exemple SWE-bench), la navigation web et la bureautique, ainsi que le raisonnement outillé et les interactions de type entreprise, chacune avec son plafond humain et ses biais de contamination ; aucune ne capture à elle seule la fiabilité d'un système agentique en production. D'où la nécessité de métriques de consistance (et non de seul potentiel), de classements coût-contrôlés et de jeux d'évaluation décontaminés — exigences que la section 9 érige en *science de l'évaluation* encore manquante.

### 4.8 Cadriciels, choix du modèle et économie des jetons

Trois fils transversaux structurent l'ingénierie des agents au-delà des patrons. Le premier est l'**économie des jetons** : chaque pas de la boucle consomme un appel au modèle, et le coût croît avec la longueur du contexte et le nombre d'agents. La mise en cache d'invite à préfixe constant peut abattre de l'ordre de 90 % le coût des jetons d'entrée répétés ; à l'inverse, une architecture multi-agents mal justifiée multiplie la facture sans gain proportionnel. Le coût n'est donc pas un paramètre d'exploitation secondaire mais une contrainte de conception de premier ordre, qui rejoint le FinOps de la section 6.7. Le deuxième fil est le **choix et le service du modèle** : retenir un modèle de raisonnement coûteux, un modèle rapide ou un petit modèle spécialisé est une décision d'architecture, non un réglage ; la portabilité inter-modèles (section 5.11) en conditionne la réversibilité. Le troisième fil est la **sécurité**, traitée comme propriété structurelle (section 4.6). Côté outillage, l'écosystème des cadriciels d'orchestration s'est densifié — graphes d'états, cadres conversationnels multi-agents, planificateurs — sans qu'aucun ne s'impose comme standard ; ce qui importe, du point de vue de l'interopérabilité, est moins le cadriciel retenu que les *protocoles* qu'il parle (section 5), car ce sont eux, et non l'outil, qui préservent le découplage vis-à-vis du fournisseur. La leçon d'ingénierie est constante : agentifier est un choix qui doit se justifier, et l'autonomie accordée doit rester découplée du périmètre d'impact, tout élargissement de ce périmètre exigeant un contrat de garde-fous proportionné.

> **Perspective recherche.** L'invariant se transpose à cette strate : isoler le raisonnement de l'exécution est un découplage ; le schéma d'outil est un contrat ; les fichiers de configuration d'agent (du type `AGENTS.md`) sont une mémoire procédurale et un contrat de comportement qui évoluent par édition plutôt que par réentraînement. Les patrons de défense sont la transposition agentique du masquage de l'information : la sécurité naît du découplage explicite entre canal de confiance et canal non fiable.

## 5. La couche d'interopérabilité agentique

### 5.1 Convergence et application du LCIM aux agents

La section 5 est le point de jonction des deux strates précédentes : elle fait converger la théorie de l'interopérabilité (section 3) et l'ingénierie agentique (section 4). Son constat fondateur est que ce qui transite entre systèmes d'IA n'est plus une donnée structurée mais des *intentions, tâches et capacités négociées à l'exécution*, échangées entre acteurs autonomes, non déterministes et auto-descriptifs. Ce déplacement fait passer le problème du contrat déterministe au **contrat comportemental probabiliste**, de la connectivité à la négociation, et de la sécurité d'un agent à la sécurité des frontières *entre* agents.

Appliqué aux agents, le LCIM (section 3.2) révèle que ses niveaux supérieurs — pragmatique et dynamique —, déjà présents dans le modèle d'origine mais longtemps de portée réduite, en deviennent les verrous opérationnels. Le point de friction se déplace vers le haut de la pile : dans l'intégration classique, il se situait au sémantique ; avec des acteurs pilotés par un grand modèle de langage, c'est l'accord *pragmatique* — le sens d'une intention déléguée, dans un contexte donné — qui devient prépondérant, sans que le sémantique cesse pour autant d'être fragile. Une lecture systématique d'un corpus de dix-huit protocoles agentiques (Yuan et coll., 2026) confirme ce diagnostic : les protocoles dominants excellent en communication et en syntaxe mais demeurent pauvres en clarification, en alignement et en vérification sémantiques — fil analytique central de toute la section.

### 5.2 Quatre plans d'interaction, quatre axes d'analyse

La revue croise deux découpages orthogonaux, annoncés à la section 2.3. Les quatre **axes d'analyse** — découverte, sémantique, identité et confiance, gouvernance des frontières — structurent respectivement les sections 5.5, 5.6, 5.7 et 5.10. Les quatre **plans d'interaction**, objet de la présente sous-section, situent chaque protocole dans son emplacement fonctionnel. Le plan **agent-outil** (vertical) relie un agent à ses capacités externes ; il est le domaine de MCP. Le plan **agent-agent** (horizontal) relie des agents autonomes entre eux ; il est le domaine d'A2A et d'ANP. Le plan **agent-humain** régit la négociation de l'interaction et la reprise de contrôle. Le plan **agent-données** régit l'accès gouverné à l'information ; il ne dispose pas d'un protocole dédié et se joue au croisement de la découverte et de la sémantique. Ce découpage évite le faux débat de la rivalité entre protocoles : MCP et A2A n'occupent pas le même emplacement fonctionnel, et leur complémentarité importe davantage que leur concurrence apparente.

### 5.3 MCP : le contrat de la frontière agent-outil

Le Model Context Protocol transforme un appariement quadratique N×M (N agents, M outils) en un appariement additif N+M en standardisant le contrat entre un client (intégré à l'agent hôte) et un serveur (exposant des outils, des ressources et des invites). Sa trajectoire de transport illustre un découplage poussé : de `stdio` à *Server-Sent Events*, puis à *Streamable HTTP*. Un cœur sans état — sans poignée de main initiale ni identifiant de session, déployable derrière un répartiteur de charge sans affinité — est visé par la révision candidate du 28 juillet 2026 : une cible architecturale, non un acquis à la date d'arrêt du socle documentaire (juin 2026). Le contrat MCP est réel mais étroit : il fixe la forme de l'échange, pas son interprétation. La répartition du contrôle qu'il instaure — les outils relèvent du modèle, les ressources de l'application, les invites de l'utilisateur — est une décision d'architecture significative, complétée par des primitives côté client (échantillonnage, *elicitation*, *roots*) qui rendent l'échange bidirectionnel et négocié.

Cette répartition mérite d'être soulignée, car elle est rarement explicitée dans le discours promotionnel sur le protocole. Le slogan du « port universel » qui relierait n'importe quel agent à n'importe quel outil masque une réalité plus nuancée : le contrat MCP standardise la *forme* de l'échange — comment un outil se déclare, s'invoque et répond — mais non son *interprétation*, qui demeure à la charge du modèle. Deux serveurs MCP exposant des outils nommés identiquement peuvent se comporter différemment ; le protocole ne garantit pas la sémantique, seulement la mécanique. La bidirectionnalité (l'*elicitation* permettant au serveur de demander une précision à l'utilisateur via le client, l'échantillonnage permettant au serveur de solliciter une complétion du modèle hôte) introduit en outre une surface de confiance nouvelle, que la sécurité de la frontière (section 5.10) devra traiter. Le contrat est donc réel mais étroit, et c'est précisément cette étroitesse qui rend la couche sémantique (section 5.6) indispensable plutôt que superflue.

> **Mise en œuvre.** La spécification MCP a connu une cadence de jalons datés : version initiale du 5 novembre 2024 ; révision du 26 mars 2025 (transport Streamable HTTP, cadre d'autorisation OAuth) ; révision du 18 juin 2025 (sorties structurées, *elicitation*, serveur qualifié Resource Server OAuth) ; révision du 25 novembre 2025 (découverte OIDC, JSON Schema 2020-12, tâches asynchrones expérimentales) ; révision candidate du 28 juillet 2026 (cœur sans état, cadre d'Extensions, MCP Apps, politique de dépréciation formelle — non finalisée à juin 2026). La gouvernance du protocole a été transférée à l'Agentic AI Foundation (Linux Foundation) le 9 décembre 2025, accompagnée d'un registre de serveurs prévisualisé le 8 septembre 2025. Plusieurs de ces éléments restent des ressources vivantes à reconfirmer avant tout engagement.

### 5.4 A2A et ANP : la frontière agent-agent

Là où MCP délègue finement une capacité délimitée, le protocole Agent2Agent encapsule une *délégation opaque* : un agent confie une tâche entière à un autre acteur autonome, sans contrôler sa boucle interne. A2A repose sur la *carte d'agent* (Agent Card), descripteur auto-portant des capacités, et sur des objets de tâche, de message et d'artefact. Cédé à la Linux Foundation le 23 juin 2025 et porté ensuite à la version 1.0 (v1.0.0 le 12 mars 2026, v1.0.1 le 28 mai 2026), il a rallié plus de 150 organisations contributrices au 9 avril 2026. Le protocole concurrent *Agent Communication Protocol* (ACP-agent), issu d'IBM/BeeAI — homonyme à distinguer de l'*Agentic Commerce Protocol* traité en 5.9, également abrégé ACP —, a convergé par fusion dans A2A en août 2025 : illustration que la convergence, dans cet espace, se fait par absorption plutôt que par standardisation négociée. L'Agent Network Protocol (ANP) propose une approche décentralisée fondée sur les identifiants décentralisés (DID) du W3C et un méta-protocole de négociation ; sa maturité reste tributaire de celle des DID, dont l'adoption large n'est pas attendue avant 2027 environ. La frontière entre MCP et A2A est reconnue floue : les tâches asynchrones de MCP empiètent sur le terrain d'A2A. Le tableau 4 en clarifie la complémentarité, qui constitue le résultat analytique central de la section.

Table: Tableau 4 — MCP et A2A : deux frontières complémentaires, non concurrentes.

| Dimension | MCP (agent-outil) | A2A (agent-agent) |
| --- | --- | --- |
| Relation | Verticale : agent → capacité | Horizontale : agent ↔ agent |
| Granularité | Capacité délimitée, contrôlée finement | Délégation opaque d'une tâche entière |
| Unité d'échange | Outil, ressource, invite | Tâche, message, artefact |
| Tenue de la boucle | L'hôte arbitre chaque appel | L'acteur distant gère sa propre boucle |
| Découverte | Registre de serveurs, manifeste | Carte d'agent auto-portante |
| Gouvernance (juin 2026) | Agentic AI Foundation | Linux Foundation |

La conséquence pratique est qu'un déploiement réaliste ne choisit pas entre MCP et A2A : il les *compose*. Un agent expose ses capacités à ses pairs par une carte A2A, et consomme ses outils et données par des serveurs MCP ; les deux contrats coexistent, chacun gouvernant un plan d'interaction distinct. C'est cette composition simultanée — et non une hypothétique unification — qui définit le socle de production.

### 5.5 Découverte, registres et nommage

La découverte se décompose en trois moments : *découverte* (qui existe ?), *sélection* (qui convient à cette tâche précise ?) et *résolution* (où et comment l'atteindre ?). La sélection concentre la difficulté agentique, car elle suppose un appariement sémantique probabiliste entre un besoin exprimé et une capacité décrite. La leçon d'UDDI s'impose ici : la curation prime sur le tout-dynamique, et la fédération sur la centralisation. Plusieurs mécanismes coexistent — cartes d'agent auto-descriptives, schémas ouverts de capacités, registres de serveurs MCP, annuaires fédérés et services de noms d'agents — sans qu'aucun standard de nommage et de résolution ne se soit encore imposé. Une distinction opérationnelle structure cette couche : un *registre* sert la découverte technique — résoudre un nom vers une description et un point d'accès vérifiables —, tandis qu'un *marketplace* sert la distribution commerciale ; confondre les deux expose au *typosquatting* de noms d'agents et à la résolution vers un imposteur, faute de vérification de propriété et de provenance.

### 5.6 Le verrou sémantique

Le verrou central de la couche est l'écart entre l'accord-de-protocole et la compréhension. Deux agents peuvent échanger des messages parfaitement formés et néanmoins diverger sur le sens d'une capacité ou d'une intention. La description d'un outil ou d'un agent en texte libre, lue et interprétée par le modèle, n'offre aucune garantie de comportement : l'écart entre la description et le comportement effectif est la source de modes d'échec propres. Une réponse théorique consiste à spécifier le contrat comportemental de manière probabiliste — par un triplet associant une probabilité minimale de satisfaction, une tolérance et un nombre d'essais (Bhardwaj, 2026) — et à le vérifier à l'exécution, rattachant le contrat d'agent aux postconditions plutôt qu'à la seule validation de schéma. Une autre réponse confie la réconciliation sémantique au modèle lui-même, au prix d'un risque d'hallucination d'alignement — une correspondance plausible mais fausse entre deux vocabulaires ; une troisième, neuro-symbolique, conserve une couche ontologique formelle qui stabilise le contrat pendant que la couche neuronale en absorbe l'évolution, dans la ligne d'une convergence « ontology-first » apparue en 2026 qui ancre capacités, intentions et tâches sur des classes explicites plutôt que sur la seule prose libre. La filiation avec la sémantique des engagements observables (Singh, 1998 ; Yolum et Singh, 2002) — fonder le sens sur des obligations sociales vérifiables plutôt que sur des états mentaux — fournit un précédent théorique direct.

### 5.7 Identité, confiance et axe agent-humain

L'identité agentique est *composite* : elle superpose l'identité du modèle, celle de l'instance et celle du délégant humain. L'autorisation déléguée s'appuie sur OAuth 2.0 et OpenID Connect, et sur l'échange de jetons (RFC 8693) pour la délégation multi-saut, avec sa revendication d'agir-pour-le-compte-de. La preuve cryptographique progresse — cartes d'agent signées, justificatifs vérifiables (W3C Verifiable Credentials 2.0, mai 2025), signatures de messages HTTP (RFC 9421) — et s'articule à des cadres d'identité de charge de travail (SPIFFE, IETF WIMSE) et décentralisée (DID). Le point dur demeure la délégation au-delà de deux sauts : OAuth atteste la délégation mais ne contraint pas l'atténuation de privilège entre sauts, de sorte que l'attribution le long d'une chaîne A→B→C n'est pas encore standardisée.

L'axe agent-humain comble une asymétrie de la taxonomie initiale des quatre plans (section 5.2) : la confiance ne se joue pas seulement entre systèmes mais à l'interface où un humain valide, corrige ou consent. Trois mécanismes datés l'instrumentent. L'*elicitation* de MCP (révision du 18 juin 2025) permet au serveur de formuler une demande structurée d'information à l'utilisateur en cours d'exécution, faisant de l'interaction une primitive du protocole plutôt qu'un hors-bande applicatif. L'extension MCP Apps fait transiter des interfaces HTML en bac à sable par le même circuit d'audit et de consentement que les appels d'outils ; elle relève toutefois de la révision candidate du 28 juillet 2026 et n'est donc pas encore stable. AG-UI normalise sur HTTP et flux d'événements (SSE) un schéma d'échange entre frontal et dorsal agentique, adopté par plusieurs cadriciels. Ces mécanismes étendent le découplage à la relation agent-humain : l'interface devient elle-même un contrat négocié, évolutif et indépendant du cadriciel sous-jacent.

### 5.8 La pile étagée et la matrice de décision

La conclusion analytique de la section est que ces protocoles ne se concurrencent pas mais *s'étagent* : ils occupent des emplacements fonctionnels complémentaires au-dessus de la pile classique (JSON-RPC, HTTP, OAuth, schémas). La composition est simultanée, non alternative : un déploiement réel fait coexister MCP (agent-outil), A2A (agent-agent), des primitives agent-humain et un cadre réglementaire. Le socle de production minimal combine MCP et A2A. Le tableau 5 reproduit la matrice de maturité et de décision arrêtée à juin 2026.

Table: Tableau 5 — Matrice de maturité et de décision des protocoles agentiques (état arrêté à juin 2026).

| Protocole | Axe couvert | Statut (juin 2026) | Gouvernance | Recommandation d'usage |
| --- | --- | --- | --- | --- |
| MCP | Agent-outil | Stable (révision du 25 nov. 2025) ; cœur sans état visé par la RC du 28 juillet 2026 | Agentic AI Foundation (LF) | Socle pour l'accès aux outils et données ; traiter le cœur sans état comme cible, non comme acquis |
| A2A v1.0 | Agent-agent | Stable | Linux Foundation | Socle pour la délégation inter-agents |
| ANP | Agent-agent (décentralisé) | Émergent, tributaire des DID (horizon ~2027) | Communauté ANP, W3C AI Agent Protocol CG | Veille ; pilotes décentralisés |
| AP2 | Paiement / mandat | Émergent (v0.2) | FIDO Alliance | Veille ; commerce agentique |
| x402 | Paiement machine-natif | Émergent, forte traction de production (x402 V2) | x402 Foundation (LF) | Rail de règlement le plus déployé parmi les protocoles machine-natifs |
| MCP Apps | Agent-humain (UI) | Expérimental (extension RC du 28 juillet 2026) | MCP Core / AAIF | Prototype, hors production critique |
| AG-UI | Agent-humain | Émergent, adopté par plusieurs cadriciels | Communauté AG-UI | Interfaces temps réel, complémentaire de l'elicitation MCP |

La logique de décision se résume ainsi : pour un système agentique interopérable de production, MCP plus A2A constituent le socle minimal ; la couche agent-humain (elicitation, MCP Apps, AG-UI) s'ajoute dès qu'un point de validation humaine ou une interface riche est requis ; la couche de règlement n'intervient que pour les cas d'usage commerciaux et reste, en juin 2026, le segment le moins stabilisé. Les protocoles marqués émergents ou conditionnés aux DID appellent une veille active plutôt qu'un engagement architectural irréversible.

### 5.9 Commerce et paiements agentiques

Le règlement agentique pose un problème spécifique : dissocier le mandant, l'agent et le commerçant, tout en prouvant l'intention, l'autorité et l'identité. Une famille de protocoles y répond, étagée en découverte/*checkout*, mandat et règlement. Côté mandat, l'*Agent Payments Protocol* (AP2), annoncé par Google le 16 septembre 2025 et transféré à la FIDO Alliance le 28 avril 2026, exprime les mandats d'intention, de panier et de paiement comme des justificatifs vérifiables signés. Côté *checkout*, l'*Agentic Commerce Protocol* (ACP-commerce, OpenAI/Stripe) sous-tend l'achat instantané en agent conversationnel, aux côtés de l'*Universal Commerce Protocol* (UCP, Google, janvier 2026), en compatibilité affichée avec AP2. Les réseaux de cartes ajoutent leurs propres dispositifs (Visa Trusted Agent Protocol, Mastercard Agent Pay), et des rails machine-natifs prolongent le mouvement : x402 (Coinbase), qui ressuscite le code HTTP 402 et concentre la plus forte traction de production parmi eux, le *Machine Payments Protocol* (Stripe/Tempo, mars 2026) et le *Know Your Agent*, qui transpose la vérification « connaître son client » à l'agent. La question de la responsabilité, du litige et de la non-répudiation reste largement ouverte, et la revue y revient en contexte financier à la section 7.

### 5.10 Sécurité non composable et modes d'échec

Le résultat de sécurité le plus important de la couche est que la sûreté n'est pas une propriété compositionnelle : un agent sûr et un outil sûr, une fois composés, ne donnent pas un système sûr. La frontière de confiance n'est plus le périmètre d'un système mais *chaque arête* du graphe d'interaction. Les attaques propres à la couche — empoisonnement d'outil, *rug-pull* (un serveur change son comportement après approbation), injection inter-agents — exploitent précisément ces arêtes. Plusieurs vulnérabilités à fort indice de gravité documentées en 2025 (exfiltration *zero-click* en environnement de productivité ; injection de commande dans un relais MCP) confirment la réalité du risque. Le durcissement procède par couche et par provenance vérifiée des registres ; la vérification d'intégrité *continue* contre le rug-pull demeure toutefois immature.

Cet angle adversarial ne recouvre qu'une partie du problème : à composition égale, l'échec peut survenir sans attaquant, par le seul jeu d'acteurs non déterministes assemblés au-dessus de contrats imparfaitement spécifiés. La taxonomie MAST (section 4.5) fournit le cadre de ces défaillances *émergentes* propres à la frontière interopérable — routage vers un agent inadéquat, perte de contexte au passage d'une frontière protocolaire, incompatibilité de schéma non détectée, absence de clarification. Un incident de production documenté en 2025 en illustre le caractère contractuel plutôt que purement technique : un agent de codage a exécuté des commandes destructrices durant un gel de code, faute d'un garde-fou organisationnel opposable à l'agent et non aux seuls opérateurs humains. Leur évaluation systématique inter-fournisseurs reste un front ouvert.

### 5.11 Portabilité inter-modèles et inter-cadriciels

Un plan de découplage propre à la couche agentique mérite d'être isolé : la portabilité. L'interface de complétion conversationnelle dite « compatible OpenAI » s'est imposée comme standard de fait, créant un verrouillage inversé où des fournisseurs concurrents émulent l'interface d'un précurseur pour capter son écosystème d'outils. La portabilité réelle se joue à deux niveaux. Au niveau de l'*invocation*, des passerelles uniformisent des familles d'API hétérogènes — complétions sans état, interfaces à état exposant des primitives agentiques, formats de messages propriétaires — au prix d'un dénominateur commun parfois appauvri. Au niveau de la *définition d'agent*, des formats émergents visent à rendre portable la description d'un agent — ses outils, ses politiques, son comportement — par-delà les cadriciels. Le constat dominant tient en une formule : la définition voyage, l'exécution reste captive. Les protocoles partagés (MCP, A2A) offrent la voie d'interopérabilité la plus robuste, car ils découplent la coopération de la pile d'exécution sous-jacente : deux agents qui « parlent MCP et A2A » coopèrent quelle que soit leur fondation. La neutralité par fondation règle la propriété du contrat ; elle ne règle ni la portabilité de l'état, ni celle de l'exécution durable, qui demeurent des fronts ouverts.

### 5.12 Conformité, observabilité et évaluation des protocoles

Vérifier qu'une implémentation respecte un protocole — la conformité — diffère de vérifier qu'elle interopère utilement. Des suites de conformité existent pour MCP et A2A, mais le test de la *négociation* — l'accord pragmatique qui s'établit à l'exécution entre acteurs probabilistes — reste embryonnaire, car il faudrait évaluer non une réponse mais une distribution de comportements. L'observabilité s'appuie sur la propagation du contexte de trace (W3C Trace Context) à travers les frontières d'agents, condition de l'auditabilité de bout en bout et, en contexte régulé, de la valeur probante du journal. La certification inter-fournisseurs — qui validerait qu'un agent d'un éditeur coopère de façon prévisible avec un agent d'un autre — n'existe pas encore. Ce manque est structurant : sans métrique partagée ni banc d'essai neutre, l'interopérabilité agentique demeure une propriété affirmée plutôt que mesurée, et la conformité de protocole ne garantit pas la composabilité de la sûreté (section 5.10).

### 5.13 Négociation, fiabilité et l'héritage des systèmes multi-agents

La couche agentique réactive un héritage théorique longtemps mis de côté : celui des systèmes multi-agents classiques. Le *Contract Net Protocol* (Smith, 1980) formalisait déjà l'appel d'offres et l'attribution de tâches entre agents ; les langages de communication d'agents (KQML, FIPA-ACL) typaient les actes de langage (informer, demander, promettre) sur une sémantique mentaliste de croyances et d'intentions. La critique de cette sémantique — invérifiable, car on n'observe pas les états mentaux d'un agent — a conduit à lui préférer une sémantique d'*engagements observables* (Singh, 1998 ; Yolum et Singh, 2002), fondée sur des obligations sociales vérifiables. Cet héritage éclaire directement le contrat comportemental probabiliste (section 5.6) : un agent ne se juge pas sur ce qu'il « croit » mais sur les engagements qu'il tient, mesurés à l'exécution. La fiabilité devient alors une propriété du *graphe* d'interaction, non d'un nœud isolé, ce qui rapproche le problème de la tolérance aux pannes byzantines : comment un collectif d'acteurs faillibles et possiblement divergents atteint-il une décision sûre ? La négociation de capacités — la plus aboutie étant celle des méta-protocoles décentralisés — reste conditionnée à la maturité des identifiants décentralisés. Le verrou central demeure pragmatique : passer d'un accord sur le protocole à un accord sur le *sens* d'une intention déléguée, sans retomber dans la fragilité des ontologies lourdes ni dans l'invérifiabilité de la sémantique mentaliste. Ce verrou, hérité de quarante ans de recherche en intelligence artificielle distribuée, est aujourd'hui reposé à l'échelle d'agents probabilistes — et reste, à juin 2026, largement ouvert.

> **Perspective recherche.** L'invariant se reformule ici pour des acteurs probabilistes : le contrat passe du déterministe au comportemental ; le découplage s'enrichit de deux plans nouveaux — découplage de modèle et de cadriciel ; l'évolution devient versionnage de capacités et politique de dépréciation formelle. La neutralité de gouvernance d'une fondation règle la question de la *propriété* du contrat ; elle ne dispense pas d'en éprouver l'exécution.

## 6. Le passage à l'échelle de l'entreprise

### 6.1 Interopérabilité n'est pas intégration : les quatre dettes

La quatrième strate distingue nettement deux notions souvent confondues : l'**interopérabilité** (les protocoles s'entendent) n'est pas l'**intégration d'entreprise** (les agents créent de la valeur sous contrainte organisationnelle réelle — héritage applicatif, silos, conformité, gouvernance répartie). Ce qui n'émerge qu'à l'échelle de l'organisation se laisse organiser par quatre *dettes* qu'aucun protocole ne solde de lui-même : la dette de **données** conditionne l'ancrage gouverné (section 6.5), la dette d'**intégration** conditionne la greffe sur l'héritage (section 6.2), la dette d'**identité** conditionne l'accès à l'échelle d'une flotte d'agents (section 6.4), et la dette de **gouvernance** conditionne la conformité et le contrôle (section 6.6). À ces dettes s'ajoute un phénomène quantitatif propre à l'organisation : la prolifération d'agents (*agent sprawl*) et l'explosion corrélative des identités machine (section 6.4).

### 6.2 De N×M à N+M : le greffon sur l'héritage

Le fil économique de la strate est la proposition **N×M → N+M** : standardiser les interfaces fait chuter le coût marginal d'intégration d'un agent ou d'un outil supplémentaire d'un appariement quadratique à un branchement additif. C'est l'argument de fond en faveur des standards ouverts, repris comme critère d'appel d'offres à la section 6.8. Concrètement, l'agent se greffe sur le tissu existant par des serveurs MCP qui encapsulent les systèmes hérités — patron *strangler fig* appliqué couche d'outil par couche d'outil, les plateformes d'intégration en tant que service (iPaaS) se muant elles-mêmes en fabriques de serveurs MCP à partir de leur catalogue de connecteurs existant —, par des passerelles d'API et d'inférence qui domicilient le trafic, le cache et le coût, et par un maillage d'événements qui prolonge les architectures événementielles jusqu'à l'EDI agentique. Sans registre interne sous curation ni discipline de versionnement et de dépréciation des serveurs, cette greffe reconstitue toutefois une dette de fragmentation analogue au point-à-point qu'elle prétendait résorber.

### 6.3 Le plan de contrôle agentique

L'architecture de référence dominante à l'échelle de l'entreprise est le **plan de contrôle agentique**, convergence indépendante de plusieurs analystes et fournisseurs en 2025-2026. Il s'organise en cinq fonctions : les agents *décident*, les plans de contrôle *gouvernent*, l'exécution *fait respecter*, les systèmes *prouvent*, et une fonction FinOps *contrôle le coût*. Son principe directeur est de **séparer l'intelligence de l'autorité, la capacité de la permission, l'explication de la preuve** : ne jamais laisser le modèle qui décide s'auto-autoriser ou s'auto-attester. Cette séparation est la transposition à l'échelle entreprise de la distinction plan de contrôle / plan de données héritée du maillage de services (section 3.6). Concrètement, le plan de contrôle fait respecter quatre disciplines indissociables. Il *autorise* : toute action passe par un point d'application de politique externe à l'agent, jamais auto-octroyé. Il *inventorie* : tout agent décisionnel est enregistré comme un actif gouverné, car on ne gouverne que ce qu'on a inventorié. Il *trace* : chaque décision déléguée produit une piste d'audit probante, distincte de la simple observabilité d'ingénierie. Il *borne le coût* : une fonction FinOps mesure le coût par résultat et coupe les dérives. Reste ouvert l'arbitrage sur sa propriété : au fournisseur de la plateforme (délai de valeur court, mais verrouillage, section 6.8), ou à l'entreprise (neutralité, mais à construire). Cette convergence en fait le motif dominant du déploiement d'entreprise — et c'est lui que la section 8 formalisera en ArchiMate et que la section 10 instanciera sur une pile concrète.

### 6.4 Identités non humaines

L'identité non humaine (NHI) devient un enjeu de premier ordre à mesure que les agents et services automatisés se multiplient dans l'organisation : une enquête fournisseur de 2025, non auditée, avance un ratio de l'ordre de 80 identités machine pour une identité humaine — ordre de grandeur non certifié, mais indicatif du déséquilibre. Le cycle de vie d'une identité d'agent — enregistrement, attribution de privilèges, rotation des secrets, et surtout *offboarding* — doit être gouverné aussi rigoureusement que celui d'un employé, faute de quoi s'accumulent des identités orphelines sur-privilégiées. La délégation doit rester traçable sans réutilisation de jeton, et l'autorisation doit être externalisée vers des points de décision et d'application de politique distincts du code de l'agent. Le tableau 6 reproduit le palmarès OWASP des dix principaux risques d'identité non humaine, qui sert de liste de contrôle. Le déséquilibre numérique est au cœur du problème : chaque identité machine est un porteur de privilèges potentiellement orphelin, sur-doté ou muni de secrets à longue durée de vie, et l'écart entre ce qu'un agent est *techniquement* capable de faire et ce que l'organisation l'a *explicitement* habilité à faire — l'*écart d'autorité* — en est la mesure. La gouvernance de l'identité non humaine n'est donc pas une extension marginale de la gestion des identités humaines, mais un régime à part entière, où l'*offboarding* fiable — le point le plus négligé — importe autant que l'enregistrement initial. Ce même principe de non-auto-autorisation (section 6.3) est ce qui rend le régime tenable à l'échelle.

Table: Tableau 6 — OWASP Non-Human Identities Top 10 (2025).

| Code | Risque | Code | Risque |
| --- | --- | --- | --- |
| NHI1 | Offboarding incorrect | NHI6 | Déploiement infonuagique non sûr |
| NHI2 | Fuite de secrets | NHI7 | Secrets à longue durée de vie |
| NHI3 | NHI tierce vulnérable | NHI8 | Isolation d'environnement déficiente |
| NHI4 | Authentification non sûre | NHI9 | Réutilisation d'identité |
| NHI5 | NHI sur-privilégiée | NHI10 | Usage humain d'une NHI |

### 6.5 Ancrage gouverné

À l'échelle de l'entreprise, l'ancrage informationnel cesse d'être un simple mécanisme anti-affabulation : il devient un *point d'application de politique*. Le principe directeur est qu'un agent ne doit jamais récupérer ni restituer un fragment que son mandataire humain ne pourrait lire directement à la source. Cela impose une récupération *consciente des permissions* — le filtrage des droits doit s'opérer à l'intérieur même de la recherche vectorielle, et non en post-traitement, qui fuit — et une attention particulière à la mémoire de flotte, dont la propagation d'étiquettes de sensibilité à travers des états persistants constitue un canal de fuite inter-tenant encore mal formalisé. Une couche sémantique versionnée — métriques et définitions certifiées une fois, plutôt que du SQL recomposé à chaque appel — et des contrats de données ouverts, qui font de l'agent un consommateur régi au même titre qu'une application, prolongent ce même principe au plan des données : un vocabulaire partagé abaisse le coût marginal d'alignement, transposition documentaire de N×M → N+M (section 6.2).

### 6.6 Gouvernance et conformité

La gouvernance mobilise trois leviers complémentaires : des obligations légales (la règle), des normes volontaires (le système de management) et des contrôles techniques (la règle exécutable). Côté normes, ISO/IEC 42001:2023 fournit le premier système de management de l'IA certifiable, et le cadre de gestion des risques du NIST (AI RMF 1.0) structure la démarche en fonctions gouverner-cartographier-mesurer-gérer. Côté obligations, le Règlement européen sur l'IA échelonne les usages par niveau de risque — un assistant anodin peut basculer vers le haut risque dès qu'il touche à l'emploi, au crédit ou à un service essentiel —, et son calendrier a bougé en 2026 : le report de l'annexe III, porté par le Digital Omnibus, a reçu le feu vert du Conseil de l'UE le 29 juin 2026, sous la seule réserve de sa publication au Journal officiel de l'Union européenne (cf. §7.4). S'y ajoutent le RGPD et, sectoriellement, le règlement DORA en finance, approfondis à la section 7. Le principe transversal est que la conformité dépend de la donnée traitée et de la décision prise, non du modèle isolé ; les correspondances entre référentiels (*crosswalks*) réduisent le coût marginal de conformité — transposition de N×M → N+M au plan réglementaire.

Cette articulation à trois leviers évite deux écueils symétriques. Le premier est le *théâtre de conformité* : empiler des politiques déclaratives sans contrôle technique exécutable, de sorte que la conformité existe sur le papier mais non dans le flux d'exécution. Le second est le *solutionnisme technique* : déployer des garde-fous sans les rattacher à une obligation légale ou à un système de management, ce qui les rend invérifiables et non opposables. La conformité agentique tient lorsque les trois niveaux se referment l'un sur l'autre : une obligation légale — par exemple une exigence de journalisation — se traduit en une règle de système de management — une politique inventoriée —, elle-même réalisée par un contrôle technique — un point d'application qui refuse l'action non journalisée. Cette fermeture verticale, du pilote réglementaire (*Driver*) au contrôle implémenté, est exactement ce que la formalisation ArchiMate (section 8.4) rend traçable, et ce que l'inventaire des modèles (section 7.5) rend opposable devant un régulateur.

### 6.7 Observabilité, FinOps et maturité

L'observabilité s'unifie autour des conventions sémantiques OpenTelemetry pour l'IA générative, dont le statut demeure en développement à la mi-2026 ; elle fonde le journal d'audit probant. La fiabilité d'un système composé est le produit des fiabilités d'étape, ce qui impose des SLO par étape plutôt qu'un SLO global trompeur. Le FinOps des jetons substitue au coût par jeton un *coût par résultat métier* — seule mesure qui rende visible l'effet de plateforme d'un parc mutualisé, alors qu'une orchestration multi-agent mal justifiée peut la multiplier sans gain (section 6.11). Enfin, la maturité de l'interopérabilité d'entreprise se mesure par une transposition du LCIM aux agents (tableau 7), qui montre que l'enjeu se joue aux étages hauts — sémantique, pragmatique, dynamique, conceptuel — là où les standards de transport n'opèrent plus, le niveau le plus élevé (alignement des modèles de tâche et d'autorité) restant largement hors d'atteinte.

Table: Tableau 7 — Transposition du LCIM aux agents (d'après Tolk et Muguira, 2003 ; Wang, Tolk et Wang, 2009).

| Niveau | Interopérabilité classique | Transposition agentique |
| --- | --- | --- |
| 0 — Aucune | Systèmes isolés, aucun échange | Agents cloisonnés par plateforme |
| 1 — Technique | Connectivité, transport partagé | Transport disponible, sans format commun |
| 2 — Syntaxique | Structure de message commune | Messages A2A/MCP bien formés ; carte d'agent, schémas d'outils |
| 3 — Sémantique | Sens des données partagé | Vocabulaire des entités métier partagé (section 6.5) |
| 4 — Pragmatique | Contexte d'usage compris | Effet attendu d'une action compris en contexte |
| 5 — Dynamique | États évoluant de concert | Coordination des états sur une trajectoire observable |
| 6 — Conceptuel | Modèles du monde alignés | Modèles de tâche et d'autorité alignés — hors d'atteinte |

> **Mise en œuvre.** La grille de décision de l'architecte d'entreprise ordonne huit choix : faire ou acheter ; jardin clos ou standards ouverts ; gouvernance possédée ou déléguée au fournisseur ; mono ou multi-agent ; protocoles et écosystème ; autonomie et reprise humaine ; identité des NHI ; portabilité du modèle et stratégie de sortie. Le canevas d'appel d'offres correspondant exige des conformités vérifiables : MCP (OAuth, indication de ressource RFC 8707), A2A (carte d'agent, compatibilité fil), observabilité OpenTelemetry, contrats de données ouverts, et journal d'audit conforme aux exigences de journalisation réglementaire.

### 6.8 Plateformes d'agents : jardins clos contre standards ouverts

À l'échelle de l'entreprise, le choix de plateforme structure la dépendance future. Trois familles coexistent. Les plateformes *infusées* dans le système de référence (suites de gestion de la relation client, progiciels de gestion intégrés, plateformes de flux de travail) offrent l'intégration la plus immédiate au système de référence, au prix d'un couplage au fournisseur. Les plateformes *cloud horizontales* (services d'agents des grands fournisseurs d'infonuagique) offrent l'élasticité et l'outillage, au prix d'une gravité de données. Les briques *neutres* issues de fondations ouvertes offrent la portabilité, au prix d'un effort d'intégration. L'arbitrage décisif oppose le jardin clos — délai de mise en valeur court, mais verrouillage — aux standards ouverts — neutralité et stratégie de sortie, mais coût d'assemblage. C'est ici que la proposition N×M → N+M devient un argument stratégique : exiger dans les appels d'offres la conformité aux standards ouverts (MCP avec ses indications de ressource, A2A avec ses cartes d'agent, contrats de données, observabilité) préserve une stratégie de sortie. Le pari d'écosystème n'est pas univoque : MCP relève de l'Agentic AI Foundation, tandis qu'A2A et AGNTCY/OASF (Cisco, LangChain, Galileo) relèvent tous deux de la Linux Foundation, mais comme *blueprints* distincts, à ne pas confondre dans les appels d'offres. Gartner projette l'annulation de plus de 40 % des projets d'IA agentique d'ici fin 2027 (Gartner, 2025), avertissement contre l'*agent washing* qui renforce l'exigence de critères scorables. Le coût d'un mauvais pari ne se mesure donc pas seulement en migration, mais en dépendance de sentier.

### 6.9 Interopérabilité inter-entreprises et espaces de données

Au-delà du périmètre d'une firme, la frontière inter-organisationnelle est, par défaut, hostile : chaque firme est pour l'autre une zone non fiable. La coopération y suppose un *Know Your Agent* — pendant machine de la connaissance du client —, des registres inter-firmes sous curation plutôt que purement décentralisés, et une politique de partage *exécutable*, et non seulement contractuelle. Les espaces de données souverains et les marchés de capacités fournissent le cadre d'un *agentic web* inter-entreprises, où s'opposent une vision ouverte et décentralisée et une vision de souveraineté régionale. L'orchestration à cette échelle se distingue nettement de la coordination intra-firme : aucun orchestrateur central ne s'impose, chaque organisation conserve le sien, et la coordination s'opère par des contrats inter-agents fédérés (section 6.11). La sécurité, la responsabilité et les modes d'échec se concentrent précisément sur ces arêtes inter-organisationnelles, qui échappent au contrôle d'une partie unique : l'attribution d'une action déléguée d'une firme à l'autre, en l'absence de standard de délégation transitive, reste l'un des points durs du domaine, traité en contexte financier à la section 7.

### 6.10 Adoption, modèle opérationnel et maturité

L'adoption à l'échelle relève d'un modèle socio-technique autant que technique. Le modèle dominant est le *moyeu et rayons* (hub-and-spoke), articulé autour d'un centre d'excellence agentique qui possède le plan de contrôle, les standards et la gouvernance, tandis que les équipes métier construisent les agents dans ce cadre. La maturité, mesurée par la transposition LCIM→agents (tableau 7), progresse lentement, car la difficulté n'est pas de connecter mais de gouverner : sortir du « purgatoire des pilotes » suppose une discipline de conduite du changement, un investissement en compétences et un régime de responsabilité clair, davantage qu'une montée en capacité des modèles. La revue retient ici une extension de l'invariant proposée par le corpus : à *découplage, contrat, évolution* s'ajoute l'*exploitation*, car un parc d'agents doit demeurer non seulement interopérable mais *opérable* — observable, remédiable, auditable — sous peine de devenir, à l'échelle, ingouvernable. Cette extension est reprise et instanciée à la section 10.

### 6.11 Orchestration à l'échelle : mono-agent par défaut et collaboration humain-agent

Le passage à l'échelle transforme aussi l'orchestration, qui cesse d'être un choix d'ingénierie pour devenir une question de propriété et de responsabilité distribuée. Trois échelles se distinguent par la structure de confiance plutôt que par la technique : *intra-équipe* (propriétaire unique, plan de contrôle homogène) ; *inter-équipes* (plusieurs propriétaires, accords de niveau de service croisés) ; *inter-organisations*, où s'ajoute une frontière de confiance sans autorité centrale commune (section 6.9). Le patron *superviseur* et le pipeline séquentiel dominent les déploiements réels ; les formes hiérarchiques conviennent au complexe inter-domaines, l'essaim et le tableau noir restant cantonnés à la recherche.

La première décision demeure toutefois économique avant d'être architecturale : faut-il un système multi-agent du tout ? Les retours d'expérience convergent sur un avertissement — un système de recherche multi-agent peut consommer, selon des mesures publiées, de l'ordre de quinze fois plus de jetons qu'une interaction conversationnelle simple. La taxonomie MAST des modes d'échec multi-agents (Cemri et coll., 2025) confirme qu'une part substantielle des défaillances tient à la spécification des rôles et à la coordination inter-agents, non aux modèles. Le multi-agent ne se justifie donc que par un parallélisme réel, une séparation de responsabilités ou de sécurité, ou des frontières organisationnelles incontournables ; hors de ces cas, le mono-agent bien conçu reste le défaut, et le mimétisme d'un organigramme humain est l'anti-patron à proscrire.

Sur le plan agent-humain, les flux de longue durée — qui s'étendent sur des heures ou survivent à des pannes — exigent une exécution *durable*, journalisée pas à pas pour reprendre sans perte d'état. La collaboration s'organise par une boîte de réception d'agents aux modes d'approbation explicites — accepter, éditer, répondre, ignorer — et par des interruptions déclenchées selon un score de confiance ; le risque inverse, la *fatigue d'approbation*, vide le contrôle de sa substance si les sollicitations se multiplient. La gestion des exceptions complète le dispositif : escalade humaine, agents vérificateurs, disjoncteurs sur les actions à risque et files de messages morts pour les tâches non résolues.

### 6.12 Sécurité du parc d'agents à l'échelle

À l'échelle du parc, la sécurité change de nature : elle ne protège plus un agent mais une population d'agents, de serveurs d'outils et d'identités machine en interaction. Plusieurs référentiels en structurent l'analyse : le palmarès OWASP des risques des applications agentiques, la base de connaissances des menaces adverses sur les systèmes d'IA, et un cadre de modélisation des menaces en couches propre aux systèmes agentiques. Trois vecteurs dominent. L'**injection indirecte** par les outils et les données reste la menace fondamentale (section 4.6) ; à l'échelle, elle se conjugue à des fuites *zero-click* documentées en environnement de productivité, où un contenu piégé exfiltre des données sans action de l'utilisateur. La **chaîne d'approvisionnement** des serveurs d'outils introduit le risque de *rug-pull* — un serveur approuvé qui change de comportement — et impose une vérification de provenance et d'intégrité continue. Enfin, la **prolifération d'identités** (section 6.4) élargit la surface d'attaque à chaque agent et secret. Les contre-mesures relèvent de la défense en profondeur : contrôle d'egress et prévention de fuite de données, confinement et isolation d'environnement, point d'application de politique unique (section 6.3), reprise humaine calibrée pour les actions sensibles — au même titre que pour l'orchestration courante (section 6.11) —, et un dispositif de réponse aux incidents adapté aux actions autonomes. La règle directrice reste que la sûreté n'est pas compositionnelle (section 5.10) : sécuriser chaque composant ne sécurise pas le parc, dont la défense doit se concevoir aux arêtes.

## 7. La spécialisation au domaine financier

### 7.1 Les quatre durcisseurs

La cinquième strate spécialise l'analyse au domaine financier, choisi pour son exigence réglementaire maximale. Sa thèse est que la finance ne crée pas de concepts nouveaux d'interopérabilité : elle en *élève le coût de défaillance*, par le cumul de quatre durcisseurs intrinsèques — l'invariant *découplage, contrat, évolution* n'y est pas remplacé, il y est porté à un coût d'échec plus élevé. Le premier est l'**irréversibilité** et la finalité du règlement : une fois un paiement réglé, il ne se défait pas. Le deuxième est l'exigence de **fonds propres et le risque systémique**, aggravé par la corrélation d'agents entraînés sur les mêmes modèles de fondation. Le troisième est la **confiance et la responsabilité juridique** : une institution répond des actes de son agent, comme l'a établi la décision *Moffatt c. Air Canada* (2024 BCCRT 149). Le quatrième est la **double qualification réglementaire** : un agent qui décide est à la fois un *modèle* (soumis au risque-modèle) et un composant de *système de TIC tiers* (soumis à la résilience opérationnelle) — deux régimes qui ne se recouvrent pas et pour lesquels une même institution doit produire deux pistes d'audit distinctes. L'exigence agentique se dérive de ces propriétés intrinsèques, non d'un catalogue de produits.

### 7.2 L'autonomie graduée sous contrôle de finalité

Le patron de conception directeur, transférable au-delà de la finance, est l'**autonomie graduée sous contrôle de finalité** : le degré d'autonomie d'un agent se calibre sur le produit *matérialité × réversibilité* de l'action, et non sur la capacité brute du modèle. Concrètement, ce degré se décline selon une échelle continue — assistance, copilote, orchestration sous revue, autonomie bornée. La règle opératoire est la *préparation par l'agent et la release humaine sur l'action irréversible* : un agent ne doit jamais exécuter une action irréversible sans garde-fou structurel. Ce patron articule le continuum d'autonomie de la section 4.2 avec le maillage réglementaire, et fournit le critère qui chapeaute toutes les décisions d'architecture du domaine. Le tableau 8 en donne la matrice de décision : le curseur d'autonomie se lit au croisement de la réversibilité de l'action et de sa matérialité.

Table: Tableau 8 — Matrice d'autonomie : réversibilité × matérialité.

| | Faible matérialité | Forte matérialité |
| --- | --- | --- |
| Action réversible | Autonomie large (préparer et exécuter) | Autonomie sous revue (exécuter, journaliser, alerter) |
| Action irréversible | Préparation puis validation légère | Préparation seule ; *release* humaine obligatoire |

La force de cette matrice est qu'elle ne dépend pas de la capacité du modèle : un modèle plus performant ne déplace pas le curseur vers le haut pour une action irréversible à forte matérialité, car ce qui borne l'autonomie n'est pas l'aptitude technique mais le coût de l'erreur et l'exigence réglementaire. C'est en ce sens que le patron est *transférable* : tout domaine peut le réinstancier en redéfinissant ses propres seuils de matérialité et de réversibilité, la finance ne faisant qu'en fournir la version la plus exigeante.

### 7.3 Standards de données financières

Le domaine financier dispose d'un substrat sémantique mûr, organisé en trois couches. La couche de **messagerie et de syntaxe** comprend ISO 20022 (paiements et post-marché — la fin de la coexistence MT/MX pour le transfrontalier, arrêtée au 22 novembre 2025, ne signifie pas la fin du format MT, dont certains rails domestiques suivent des calendriers distincts), FIX (marchés), FpML (dérivés) et ACORD/AL3 (assurance). La couche de **modèle et de capacité** comprend BIAN (architecture de référence bancaire) et le *Common Domain Model* (CDM) porté par l'ISDA. La couche d'**ontologie** comprend FIBO (*Financial Industry Business Ontology*). La distinction est opératoire : un message bien formé selon la première couche peut rester sémantiquement faux sans la troisième. Le CDM, associé au reporting réglementaire numérique, illustre le *standard-comme-code exécutable* : la règle devient code adossé au modèle, de sorte que l'agent *orchestre* au lieu de *deviner*. Ces standards préexistants, contractuels et gouvernés par versionnement, sont l'incarnation même de l'invariant de la revue. Un point de discipline éditoriale conditionne la lecture des sous-domaines qui suivent (section 7.6) : à juin 2026, aucun « MCP financier » sectoriel n'est ratifié — les serveurs MCP d'éditeurs de cœur bancaire ou d'assurance relèvent de l'adoption produit, non d'un standard interopérable au même titre qu'ISO 20022 ou ACORD. Le tableau 9 résume la stratification.

Table: Tableau 9 — Le substrat sémantique financier en trois couches.

| Couche | Rôle | Standards (exemples) |
| --- | --- | --- |
| Messagerie et syntaxe | Forme du message échangé | ISO 20022 (paiements), FIX (marchés), FpML (dérivés), ACORD (assurance) |
| Modèle et capacité | Cycle de vie, capacités métier | BIAN (banque), CDM (ISDA, post-marché) |
| Ontologie et sémantique | Sens partagé, désambiguïsation | FIBO |

Cette stratification a une portée opératoire directe pour l'agent : un message conforme à la couche de syntaxe peut demeurer sémantiquement faux faute d'ancrage dans la couche ontologique, et c'est précisément ce que l'ancrage gouverné (section 6.5) doit prévenir. Lorsqu'un standard devient *code exécutable* — comme le modèle de domaine commun adossé au reporting réglementaire —, l'agent cesse de *deviner* la règle pour l'*orchestrer*, ce qui déplace la charge de la preuve de l'extraction vers la vérification, et réduit d'autant la surface de défaillance non déterministe.

### 7.4 Le maillage réglementaire

L'agent financier est pris dans un maillage de qualifications cumulatives et non exclusives, qui varient par juridiction. Le tableau 10 en donne la grille. La double qualification — modèle et système de TIC tiers — est le patron universel du domaine (section 7.1) : dès qu'un agent s'exécute sur une infrastructure fournie par un tiers et qu'il décide, il relève simultanément de la résilience opérationnelle et du risque-modèle — la double piste d'audit qui en résulte a été posée en 7.1. Une troisième qualification, celle de *décision automatisée*, s'y ajoute de façon conditionnelle plutôt qu'universelle : elle ne se déclenche que lorsque l'agent produit une décision opposable au client dans un domaine listé (crédit, emploi, services essentiels).

Cette pluralité conditionnelle a une conséquence pratique sous-estimée : un agent d'octroi de crédit — cas où les trois qualifications se cumulent — est simultanément un système de TIC tiers, un modèle et un dispositif de décision automatisée, sans qu'aucun régime ne prime absolument. Chacun exige sa propre preuve : registre des dépendances et stratégie de sortie pour la résilience ; inventaire, validation indépendante et surveillance continue pour le risque-modèle ; information claire et droit de contestation pour la décision automatisée. Chercher un régime unique qui subsumerait les autres est une erreur d'analyse ; l'architecture doit produire une piste d'audit *plurielle*, capable de répondre à chaque autorité dans son propre référentiel — c'est cette exigence qui fait du double contrat, modèle et service tiers, l'unité de conception du domaine financier, et non l'agent pris isolément.

Table: Tableau 10 — Grille de qualification multiple de l'agent financier (régimes cumulatifs).

| Dimension | Qualification | Régimes applicables (exemples) |
| --- | --- | --- |
| Système de TIC tiers | Résilience opérationnelle | DORA (UE) ; OSFI B-10/B-13 (Canada) |
| Modèle décisionnel | Risque-modèle | SR 26-2 / OCC 2026-13 (US) ; OSFI E-23 ; AMF (Québec) |
| Décision automatisée | Transparence et droits | AI Act annexe III (UE) ; Loi 25 art. 12.1 (Québec) ; RGPD art. 22 |

Le versant tiers TIC de cette grille n'est pas qu'une obligation de registre : depuis la désignation, le 18 novembre 2025, de dix-neuf fournisseurs critiques de TIC par les autorités européennes de surveillance — dont les trois grands hyperscalers —, un agent hébergé sur une telle infrastructure hérite d'une dépendance supervisée, assortie d'une stratégie de sortie testable. Reloger cette dépendance en cloud souverain ne la dissout pas : les offres souveraines s'appuient le plus souvent sur les mêmes hyperscalers, si bien que la garantie de premier ordre tient moins à la localisation qu'à la gouvernance de l'identité et des accès — ce qui relie le deuxième durcisseur (risque systémique, section 7.1) au registre des tiers TIC, dont les limites face au risque collectif sont discutées en 7.8.

> **Mise en œuvre.** Quelques ancrages datés, arrêtés à juin 2026 : DORA est applicable depuis le 17 janvier 2025 ; les obligations « haut risque » de l'AI Act étaient initialement prévues au 2 août 2026, mais le report porté par le Digital Omnibus a été **adopté** (feu vert du Conseil le 29 juin 2026), décalant l'annexe III autonome au 2 décembre 2027 (publication au JOUE à suivre) ; la directive SR 26-2 / bulletin OCC 2026-13 (17 avril 2026) remplace SR 11-7 mais **exclut explicitement** l'IA générative et agentique et demeure non contraignante ; la ligne directrice E-23 du BSIF (finale le 11 septembre 2025) et la ligne directrice de l'AMF sur l'IA (finale le 7 avril 2026) entrent en vigueur le 1er mai 2027 ; la fin de la coexistence MT/MX en messagerie transfrontalière ISO 20022 est arrêtée au 22 novembre 2025.

### 7.5 Risque-modèle, identité de l'agent et résidence

Le principe « si l'agent décide, c'est un modèle » impose d'inscrire l'agent à l'inventaire des modèles, avec validation indépendante et surveillance continue, et de distinguer l'observabilité d'ingénierie de l'audit réglementaire opposable, qui exige une reconstitution **infalsifiable** de la décision (signature matérielle, journal en ajout seul, chaînage cryptographique). La finance transpose aussi à l'acteur-agent un principe de contrôle interne plus ancien qu'elle, la **ségrégation des tâches** : l'agent qui propose une opération ne doit jamais être l'agent — ni l'humain — qui l'approuve, et le contrôle à quatre yeux (*four-eyes*) exige une indépendance d'exécution, sans canal latéral entre proposeur et approbateur, faute de quoi deux agents nominalement séparés peuvent reconstituer un chemin de collusion. Le référentiel ouvert FINOS *AI Governance Framework* v2.0 (12 novembre 2025) codifie cette exigence d'audit et de ségrégation sans se substituer aux cadres réglementaires. L'identification de l'acteur-machine — *Know Your Agent* (KYA) — transpose à l'agent le triptyque de la connaissance du client et de l'entreprise ; elle n'est pas encore standardisée, les initiatives observées relevant du positionnement fournisseur. La résidence des données doit enfin couvrir non seulement le repos, mais l'inférence, les embeddings, les traces et les journaux — un arbitrage entre modèle confiné et modèle frontière.

### 7.6 Les cinq sous-domaines

La revue spécialise ensuite l'analyse à cinq sous-domaines, chacun durci à sa manière. La ligne de partage commune sépare partout la *productivité assistée* — copilotes qui préparent, résument, rédigent sous revue humaine, à faible coût d'erreur — de l'*autonomie transactionnelle* — agents qui engagent une action matérielle et souvent irréversible. Le patron de l'autonomie graduée (section 7.2) s'applique à chaque sous-domaine, mais la position du curseur y diffère selon la criticité réglementaire.

#### 7.6.1 Banque

En banque, la productivité assistée est déjà répandue — copilotes conversationnels internes, assistants de service client —, tandis que l'autonomie transactionnelle reste bornée par l'irréversibilité. Le revirement, à partir de 2025, d'un assistant client d'abord promu comme substitut à des centaines d'agents humains vers un modèle hybride réintégrant l'humain illustre la leçon transférable : une demande de routine se traite par agent, mais le litige, la fraude ou la détresse financière exigent l'escalade humaine. Quatre zones concentrent l'exigence : l'identification du client et la lutte contre le blanchiment, où l'agent investigue mais où l'humain décide de la déclaration ; l'octroi de crédit et l'origination hypothécaire, qualifiés haut risque par l'annexe III de l'AI Act (scoring de solvabilité des personnes physiques, *carve-out* fraude), soumis au risque-modèle et à l'explicabilité opposable ; les paiements en temps réel, où la finalité interdit le rejeu ; et la banque ouverte, où l'agent devient destinataire accrédité de données consenties plutôt que collecteur non autorisé. L'enjeu d'intégration dominant reste l'encapsulation du cœur bancaire hérité — souvent sur grand système — par une façade contractuelle (par exemple un serveur MCP, non normalisé sectoriellement) exposant une capacité sans réécrire le système de référence.

#### 7.6.2 Assurance de dommages

En assurance de dommages, le tournant est l'ouverture des systèmes de gestion de polices et des plateformes de sinistres aux protocoles agentiques, par l'intermédiaire des standards ACORD — sans qu'il s'agisse d'un « MCP assurance » officiellement ratifié (section 7.3). Les chaînes de valeur — première déclaration de sinistre, évaluation, règlement, subrogation — se prêtent à l'orchestration, l'agent préparant l'instruction et l'humain engageant le paiement. La tarification et la souscription, nourries de télématique, posent la question de l'équité ; la détection de fraude illustre le couple investigation automatisée / décision humaine. Un point de qualification, à tenir fermement, distingue ce sous-domaine du suivant (section 7.6.3) : à la différence de la tarification vie et santé, la tarification et la souscription en dommages **ne figurent pas** à l'annexe III de l'AI Act — l'étendre à l'IARD serait une erreur. Le centre de gravité réglementaire s'y déplace donc vers le risque-modèle et la résilience opérationnelle plutôt que vers le droit des systèmes à haut risque, ce qui interdit une transposition uniforme entre sous-domaines assurantiels.

#### 7.6.3 Assurance de personnes (vie et santé)

En assurance de personnes, la souscription accélérée régulée et l'usage de la donnée de santé comme contexte (par exemple via des standards d'interopérabilité de santé adossés à MCP) dominent, sous une double contrainte d'équité et de sensibilité des données. La tarification et l'évaluation du risque en vie et santé sont explicitement qualifiées haut risque par l'annexe III de l'AI Act — à la différence de l'assurance dommages (section 7.6.2). Le contre-exemple d'un système d'aide à la décision contesté en assurance santé — litige en cours, recours partiellement maintenu en 2025 — rappelle que l'autonomie mal bornée sur une décision défavorable au bénéficiaire expose à un risque juridique et réputationnel majeur ; la *release* humaine n'y est pas optionnelle. Sur les rentes et la longévité, l'horizon décennal de l'engagement inverse la priorité habituelle : la stabilité du modèle actuariel prime sur l'agilité de l'agent, cantonné à la préparation.

#### 7.6.4 Gestion de patrimoine et d'actifs

En gestion de patrimoine et d'actifs, le copilote du conseiller — qui écrit dans le système de relation client et prépare la recommandation — prévaut, encadré par le devoir fiduciaire et l'adéquation du conseil. L'architecture dominante est fédérée, à registre de greffons, où un agent central orchestre des capacités spécialisées, et où l'interdiction contractuelle de formuler un conseil personnalisé maintient délibérément le copilote en deçà du seuil qui déclencherait les obligations de convenance. Il importe de distinguer le conseiller-robot déterministe classique — moteur de règles auditable par construction — de l'agent fondé sur un modèle de langage, dont le non-déterminisme appelle des garde-fous différents. Côté trading et recherche augmentée, la littérature documente des architectures multi-agents spécialisés (analyse, risque, exécution) où la ségrégation des rôles reconduit, au niveau du marché, le contrôle à quatre yeux (section 7.5) ; l'état réel demeure toutefois celui de la recherche et des pilotes, non de la production autonome. Côté gestion d'actifs institutionnelle, les opérations post-marché — réconciliation, valeur liquidative — s'appuient sur le modèle de domaine commun et son reporting adossé au modèle, sur des tâches à faible irréversibilité où l'exactitude prime sur l'autonomie.

#### 7.6.5 Services de technologies de l'information financières

Enfin, les services de TI financiers placent au centre la résilience opérationnelle — cartographie et registre des tiers TIC, stratégie de sortie testée, maîtrise du risque de concentration (section 7.4) —, la modernisation du cœur (où une façade contractuelle permet une migration sans rupture, sans que l'adoption du MCP par un éditeur ne constitue un standard sectoriel ratifié) et le plan de contrôle comme *point de conformité* unique, couplé à l'inventaire des modèles. C'est dans ce sous-domaine que l'extension *exploitation* de l'invariant (section 6.10) prend tout son sens : l'agent d'exploitation est lui-même un actif gouverné, inscrit au registre et soumis aux mêmes contrôles que ceux qu'il surveille.

Le tableau 11 illustre, par des cas datés et leur statut probant, la diversité des déploiements financiers à juin 2026 ; il importe d'en lire le statut, car un chiffre auto-déclaré par un fournisseur n'a pas la valeur d'une mesure académique ou d'une décision judiciaire.

Table: Tableau 11 — Exemples de déploiements financiers (statut probant variable, arrêté à juin 2026).

| Sous-domaine | Type de cas | Statut / nature de la donnée |
| --- | --- | --- |
| Banque | Copilotes de productivité (suites internes, assistants) | Adoption auto-déclarée par les institutions |
| Banque | Service client conversationnel | Cas médiatisé, position révisée en 2025 |
| Assurance dommages | Détection de fraude, traitement des sinistres | Partenariats éditeurs, gains auto-déclarés |
| Crime financier | Agent de lutte contre le blanchiment | Annonce 2026, disponibilité ultérieure |
| Gestion d'actifs | Copilote sur plateforme de gestion | Capacité auto-déclarée par l'éditeur |
| Responsabilité | Agent conversationnel engageant la firme | Décision *Moffatt c. Air Canada* (2024), fait juridique |

Au-delà des cas, deux dispositifs encadrent l'expérimentation prudente du domaine : les *bacs à sable réglementaires*, où un superviseur autorise un essai borné sous surveillance, et les bancs d'essai sectoriels, encore lacunaires. Leur rareté confirme que la frontière de la production, en finance, est tenue moins par la capacité technique que par l'exigence de preuve.

### 7.7 Commerce, paiements et responsabilité en finance

L'agentique remodèle les paiements non en automatisant la finalité du règlement, mais en s'insérant *en amont*, là où l'intention se forme, se compare et se négocie. La pile de paiement agentique se lit en trois couches : intention et orchestration (où le non-déterminisme est admis), autorisation (vérification de mandat) et règlement (déterministe, irréversible) — le raisonnement probabiliste ne devant jamais franchir vers les deux couches basses. Le *mandat vérifiable*, intention signée à limites explicites (montant, bénéficiaire, fenêtre temporelle), joue une fonction probatoire : c'est la seule forme de contrôle possible sur une action qui, une fois franchie, ne peut plus être défaite. Deux voies d'insertion, encore émergentes à juin 2026, la matérialisent : les réseaux de cartes étendent leur tokenisation existante par des jetons de paiement agentiques liant la transaction à un agent identifié, sans modifier la répartition de responsabilité entre émetteur, acquéreur et marchand ; des protocoles ouverts — au premier rang AP2, conçu pour le cas où l'humain n'est pas présent — rendent ce même mandat vérifiable et archivable. Le point dur est le trou de responsabilité : sur les rails à finalité irréversible (stablecoins, rails instantanés), la rétrofacturation n'existe tout simplement pas ; et là où elle existe, sur les rails de cartes, elle a été conçue pour la fraude ou le défaut marchand, non pour une instruction autorisée mais non voulue d'un acteur-machine. La directive européenne sur les services de paiement ne tranche pas non plus la délégation à l'acteur-machine, l'authentification forte restant pensée pour un payeur humain présent. La jurisprudence naissante — *Moffatt c. Air Canada* — confirme que l'institution répond de son agent. Au Canada, cette stratification s'instancie sur un rail national encore en construction : le Real-Time Rail, dont le cadre de règles doit entrer en vigueur le 24 août 2026, avec un lancement phasé à compter du quatrième trimestre 2026 (PROGRAMMÉ), sera irrévocable et fondé sur ISO 20022 — reconduisant sans atténuation le patron de préparation par l'agent et de *release* humaine sur l'irréversible.

### 7.8 Risque systémique, équité et explicabilité opposable

Cinq questions ouvertes, propres au durcissement financier, résistent à toute solution purement technique. La première est la **responsabilité** le long d'une chaîne N+M peuplée d'identités non humaines : qui répond d'un crédit refusé ou d'un sinistre rejeté par un agent — l'institution, l'éditeur, l'intégrateur — n'est pas tranché, et la seule réponse d'architecture disponible est une chaîne de mandat traçable jusqu'au principal humain. La deuxième est le **risque systémique d'agents corrélés** : des institutions déployant des agents fondés sur les mêmes modèles de fondation voient leurs décisions se corréler, ce qui peut amplifier les mouvements de marché par mimétisme (*herding*) ; ni la résidence, ni le risque-modèle individuel, ni le seul registre des tiers TIC (section 7.4) ne traitent ce risque collectif, qui relève de la surveillance macroprudentielle. La troisième est l'**équité** : les critères mathématiques d'équité algorithmique sont mutuellement incompatibles, de sorte qu'aucune configuration d'agent ne « résout » l'équité ; le choix relève d'un arbitrage explicite et documenté. La quatrième est l'**explicabilité opposable** : produire, sur un artefact non déterministe, une explication qui tienne devant un régulateur ou un tribunal demeure un problème de recherche ouvert. La cinquième est l'**assurabilité** : couvrir le préjudice causé par un agent suppose de tarifer un risque non stationnaire, peu documenté actuariellement. Ces cinq questions confirment la thèse du durcissement : la finance n'invente pas de nouveaux concepts d'interopérabilité, elle en élève le coût de défaillance au point de rendre indispensables des garde-fous que d'autres secteurs jugeraient optionnels.

> **Perspective recherche.** Les questions ouvertes du domaine financier sont denses : attribution de responsabilité le long de chaînes N+M peuplées d'identités non humaines ; équité algorithmique, dont les critères sont mathématiquement incompatibles entre eux ; explicabilité opposable sur des artefacts non déterministes ; risque systémique d'agents corrélés, amplifié par la concentration des fournisseurs de modèles et d'infrastructure ; et assurabilité de préjudices causés par des agents. Aucun banc d'essai public mature n'existe encore en assurance ni en lutte agentique contre le blanchiment, et aucun « MCP financier » sectoriel n'est ratifié.

## 8. La formalisation architecturale en ArchiMate

### 8.1 Le verrou méthodologique et la convention de version

La sixième strate transforme les acquis discursifs précédents en un *blueprint* d'architecture d'entreprise vérifiable, exprimé en ArchiMate. Le verrou méthodologique est posé d'emblée : ni ArchiMate 3.2 ni ArchiMate 4 n'introduisent d'élément natif pour l'agent autonome, l'appel d'outil, l'interaction inter-agents, l'identité non humaine ou le plan de contrôle. La seule extension défendable est le mécanisme officiel de spécialisation, de stéréotype et de profils, sur le modèle du *Risk and Security Overlay* (Band et coll., 2019) : chaque concept agentique se modélise comme une composition gouvernée d'éléments existants, jamais par invention de vocabulaire — phrase-test : si l'énoncé tient encore, une fois le mot ArchiMate retiré, comme un exposé des strates précédentes, c'est une redondance à écarter.

Le *blueprint* s'écrit dans sa version de référence — ArchiMate 4 (The Open Group, doc C260, 27 avril 2026) —, chaque patron portant une note d'équivalence vers ArchiMate 3.2 (doc C226, octobre 2022) : à la mi-2026, l'outillage d'entreprise ne prend encore en charge que la notation 3.2, sans export ni import natif de la v4 (⚠ ressource vivante). La v4 remplace les couches par sept domaines de même rang — dont un domaine **Common** nouveau — et réduit le compte d'éléments d'environ soixante à quarante-deux en unifiant les comportements (Service, Process, Function, Event). Six éléments antérieurs (Interaction, Constraint, Contract, Representation, Implementation Event, Gap) sont retirés, chacun remplacé par une spécialisation stéréotypée du registre normalisé.

### 8.2 Le catalogue de patrons agentiques

Le cœur normatif de la strate est un catalogue de patrons projetant chaque concept agentique sur des éléments ArchiMate existants, qualifiés par un stéréotype plutôt qu'inventés. Un agent reste un composant applicatif qualifié `<<Agent>>`, doté d'attributs de profil — autonomie, matérialité, réversibilité — alimentant des cartes de chaleur de risque ; un serveur d'outils reste un composant applicatif réalisant un service, qualifié `<<MCP Server>>` — discipline qui préserve l'interopérabilité avec l'outillage existant et découple le modèle de l'outil qui le lit.

Le tableau 12 en donne la synthèse. Le patron de l'agent associe un composant applicatif, un rôle assigné et un **processus** de raisonnement — non une fonction : la boucle plan-agir-observer, enchaînement ordonné, est un *Process* ; la *Function*, non séquentielle, reste réservée à l'appel d'outil atomique. Le patron MCP associe un composant serveur, un service réalisé, une interface et une relation de service (*Serving*) vers l'agent. Le patron A2A réserve cette même relation à l'offre durable d'un agent-spécialiste et rend l'échange ponctuel par un flux ou un déclenchement — les confondre est un anti-patron fréquent. Le patron humain-agent (quatre yeux) se distingue par l'*absence délibérée* de flux direct entre proposeur et approbateur, preuve de non-collusion. Le patron de mémoire et d'ancrage gouverné décompose la persistance de l'agent en objets distincts — mémoire, contexte récupéré, couche sémantique —, exposant les points d'ancrage des contrôles de fuite de données.

Table: Tableau 12 — Patrons agentiques et leur réalisation en ArchiMate (registre des stéréotypes).

| Concept agentique | Origine (strate) | Réalisation ArchiMate |
| --- | --- | --- |
| Agent | Section 4 | Composant applicatif `<<Agent>>` + rôle assigné + processus `<<reasoning loop>>` |
| Appel d'outil MCP | Section 5.3 | Composant `<<MCP Server>>` + service réalisé + interface + *Serving* vers l'agent + accès à la ressource |
| Interaction A2A | Section 5.4 | Flux ou déclenchement entre composants + collaboration applicative durable |
| Identité non humaine | Section 6.4 | Rôle porté par assignation + justificatif + propriétaire métier |
| Plan de contrôle | Section 6.3 | Collaboration applicative `<<Control Plane>>` réalisant des services de gouvernance |
| Mémoire, RAG, ancrage | Section 6.5 | Objets de données de mémoire et de contexte + objet métier de la couche sémantique |
| Exigence réglementaire | Sections 6.6, 7.4 | Spécialisation d'exigence `<<regulatory-requirement>>` |

### 8.3 Instanciation par domaine : capacité, altitude métier et double matérialité

Au domaine Strategy, l'unité de planification est la **capacité** agentique, spécialisant une capacité BIAN ou ACORD — jamais une capacité par serveur MCP, anti-patron confondant le savoir-faire de l'entreprise avec l'outil. Au domaine Business, le choix d'altitude gouverne l'imputabilité : un agent ne remonte en Business, comme rôle assigné à un acteur automate, que s'il porte une responsabilité métier ou réglementaire propre ; sinon il reste un composant applicatif réalisant cet acteur — discipline reprise sur les cinq chaînes de valeur sectorielles (souscription, sinistre, conseil patrimonial, paiement, crédit), l'agent positionné au bon jalon sans redécrire le processus sous-jacent. Aux domaines Application et Technology s'applique la **double matérialité** : logique délégable en Application, substrat d'exécution en Technology, distinction qui subsiste en v4 comme convention d'altitude plutôt que frontière de couche. Pour la résidence Canada-Québec, un lieu domicilie nœuds et artefacts sous une zone de conformité délimitée, l'interdiction de sortie se traduisant par l'*absence* de chemin sortant — l'invariant *découplage, contrat, évolution* se matérialise ici en un contrat de résidence vérifiable.

### 8.4 La chaîne de traçabilité verticale et l'overlay sécurité/risque

La posture zéro-confiance prolonge cette instanciation : un point d'application de politique — franchissement obligatoire, distinct du point de décision — encadre chaque interaction, sur le principe *Know Your Agent*. La contribution opposable de la strate reste la *chaîne de traçabilité verticale* : un pilote (*Driver*) influence une évaluation, qui réalise une exigence stéréotypée réglementaire, laquelle se concrétise en une action ou fonction applicative portée par un composant et un service de technologie — aucune exigence réglementaire ne devant demeurer orpheline. Cette chaîne fonde la matrice d'auditabilité croisant exigence, contrôle, domaine sectoriel et élément porteur.

Le point de vue sécurité et risque applique le même mécanisme à la chaîne de menace, reprise sans réinvention du *Risk and Security Overlay* (cf. §8.1) : un événement menaçant, qualifié en risque, appelle un objectif de contrôle réalisé par une mesure de sécurité puis par un contrôle implémenté — service de point d'application de politique ou module matériel de sécurité — associé à l'actif protégé, absorbant les catalogues de risques agentiques comme paires risque-contrôle.

### 8.5 Exemple de bout en bout, patrons et maturité

Le *blueprint* se valide sur un exemple traversant six domaines — Motivation, Strategy, Business, Application, Technology, Implementation & Migration —, la souscription vie augmentée et sa variante de déclaration de sinistre en dommages, éprouvant la cohérence inter-patrons et la traçabilité ; le modèle y est confiné en mode local pour la résidence des données (section 8.3). Il fournit en outre une bibliothèque de patrons réutilisables et neuf anti-patrons nommés avec leur correction — l'agent « boîte magique » mono-couche, qui confond structure et comportement, ou le contrôle humain décoratif, placé hors de l'étape irréversible —, ainsi qu'une feuille de route par *plateaux* : chaque palier de maturité agrège des incréments de capacité et débloque des contrôles vérifiables, transposant l'échelle LCIM→agents (section 6.7) en jalons reliés à l'observabilité plutôt qu'à des déclarations d'intention.

### 8.6 Points de vue et gouvernance des vues

Au-delà des patrons, le *blueprint* organise des *points de vue* — au sens d'ISO/IEC/IEEE 42010 — adaptés à chaque partie prenante : une vue d'autorité et de délégation pour le décideur, une vue du plan de contrôle pour l'architecte, une vue d'interopérabilité agentique pour l'ingénieur d'intégration, une vue de conformité et d'auditabilité pour le régulateur, et une vue d'ancrage gouverné pour le responsable des données. Chaque vue fixe ses préoccupations, son public et son niveau d'abstraction ; toutes se réconcilient dans un référentiel fédéré — états *as-is*/*to-be* — confronté à l'inventaire du parc réel. Cette gouvernance des vues *est* la contribution opposable du *blueprint* : elle conditionne la validité de la chaîne de traçabilité et la lisibilité du modèle pour des parties prenantes aux préoccupations disjointes.

### 8.7 Limites assumées du langage

La strate assume quatre limites structurelles. ArchiMate est un langage de description *statique* : il décrit l'intention, non l'état d'exécution. Le *non-déterminisme* n'y est pas représentable — un processus modélise le cadre attendu, jamais le chemin effectivement suivi. L'*état d'instance* échappe au langage, qui décrit des types. Enfin, la *découverte dynamique à l'exécution* n'est pas native : le modèle décrit les artefacts de la découverte, non son occurrence. Un fait critique encadre ces limites : la révision majeure du langage publiée en 2026 a refondu son métamodèle *sans ajouter la moindre prise en charge de l'IA ou des agents* — l'occasion était disponible et n'a pas été saisie, ce qui laisse la question scientifique entière et justifie le choix, assumé comme daté et révisable, de l'*overlay* plutôt que de l'extension du métamodèle.

> **Perspective recherche.** La modélisation atteste l'*intention* de contrôle, jamais son *exercice*. D'où une complémentarité obligée entre le modèle d'intention (ArchiMate) et l'observabilité, c'est-à-dire l'inventaire vivant du parc réel. Les travaux sur les modèles vivants, le jumeau numérique organisationnel et l'auto-génération de vues par grands modèles de langage demeurent immatures, en tension directe avec la nature statique du langage.

## 9. L'horizon 2027-2032

### 9.1 Le tri épistémique

La septième strate, prospective et *capstone*, ne décrit aucun état nouveau du présent : elle projette l'arc des six strates précédentes sur la fenêtre 2027-2032, en s'imposant une discipline de tri. Chaque énoncé sur l'avenir est classé **PROGRAMMÉ** (engagement daté et sourcé), **PROJETÉ** (prévision d'analyste ou d'institution, nominative et millésimée) ou **SPÉCULATIF** (pari de recherche, scénario), un quatrième label — **FAIT ÉTABLI** — réservé aux faits passés démontrés par source primaire ; la règle absolue est de ne jamais présenter du projeté ou du spéculatif comme acquis.

Le programmé lui-même n'est pas figé : les reports de l'AI Act ne sont devenus un engagement daté qu'à l'adoption du Digital Omnibus par le Conseil de l'UE le 29 juin 2026 — JOUE à venir —, d'où la mention « à re-vérifier » sur toute échéance 2026-2028. Le tableau 13 en donne des exemples représentatifs.

Table: Tableau 13 — Tri épistémique des énoncés sur l'horizon (exemples représentatifs).

| Statut | Nature | Exemples (source, année) |
| --- | --- | --- |
| PROGRAMMÉ | Engagement daté | AMLR applicable au 10 juillet 2027 ; fin des frais de changement infonuagique au 12 janvier 2027 (Data Act) ; E-23/AMF en vigueur au 1er mai 2027 ; dépréciation cryptographique classique 2030, interdiction 2035 (NIST IR 8547, projet) |
| FAIT ÉTABLI | Fait passé daté | Campagne d'espionnage conduite à ~80-90 % de façon autonome via MCP (Anthropic, nov. 2025) |
| PROJETÉ | Prévision millésimée | « Plus de 40 % des projets d'IA agentique annulés d'ici fin 2027 » (Gartner, 2025) ; ~95 % des initiatives sans effet mesurable (MIT Project NANDA, 2025) ; ~1 300 G$ de dépense IT agentique en 2029 (IDC, 2025) |
| SPÉCULATIF | Pari, discontinuité | Spécification conjointe MCP/A2A ; effet N×M → N+M (non chiffré publiquement) ; fenêtre de crédibilité d'un ordinateur quantique 2029-2032 |

### 9.2 La grappe d'échéances et la trajectoire des protocoles

La proposition cardinale de la strate est que les protocoles d'interopérabilité agentique sont, pour l'essentiel, *livrés, gouvernés, déployés et formalisés* ; ce qui reste ouvert n'est plus surtout technique mais **institutionnel et économique**. Le squelette daté — la grappe d'échéances de 2027 (réglementation financière et IA, identité, paiements) — est dense en programmé ; 2032, lui, est dominé par le projeté et le spéculatif, à mesure que le cône d'incertitude s'élargit.

Ce qui reste ouvert n'est plus *quelles fonctionnalités* viendront mais *comment les standards évoluent* : le Model Context Protocol abandonne les jalons datés pour des aires de priorité pilotées par des groupes de travail, avec un cycle Actif → Déprécié → Retiré sur au moins douze mois — l'invariant *évolution* appliqué à la gouvernance du protocole lui-même. A2A, porté comme MCP par la Linux Foundation, annonce sans date engagée une spécification d'interopérabilité, des registres consolidés et une certification, complémentaire de MCP (agent-outil contre agent-agent) sans garantir l'interopérabilité effective.

La trajectoire tend vers une *coexistence stratifiée gouvernée* plutôt qu'un protocole unifié, par analogie avec HTTP, WebSocket et gRPC : MCP et A2A se répartiront par couche d'usage, leur fusion demeurant spéculative. Un dernier verrou subsiste : malgré AGENTS.md, l'OASF et l'Agent Definition Language d'Eclipse LMOS, aucun format d'agent portable *de jure* n'existe, et la portabilité de la définition ne règle pas celle de l'état d'exécution (mémoire, sessions, points de reprise), attachée au runtime — découplage protocolaire partiel, aux conséquences directes sur souveraineté et coût total de possession.

### 9.3 Bifurcation de la gouvernance et horloge post-quantique

La « fondation neutre » se fragmente structurellement par couche : transport et communication sous l'Agentic AI Foundation ; paiement et identité-commerce sous la FIDO Alliance ; identité décentralisée, sémantique et découverte dispersées entre le W3C, l'IETF et la Decentralized Identity Foundation. Le risque de silos se déplace ainsi vers les jointures entre fondations plutôt que d'être dissous par leur neutralité affichée.

L'**identité vérifiable** et la **cryptographie post-quantique** imposent l'horloge de l'horizon : la confiance agentique repose sur des primitives à clé publique portant une date de péremption, et la migration du NIST — dépréciation classique visée en 2030, interdiction en 2035, échéances portées par un document encore à l'état de projet (IR 8547) — fait de la crypto-agilité une contrainte de conception, non une option. Ces deux échéances, avec la grappe réglementaire de 2027, sont les ancrages les plus fermes de l'horizon.

Cette bifurcation n'est pas neutre pour l'interopérabilité : elle déplace le risque de silo vers les *jointures* entre couches gouvernées séparément — un mandat signé selon une fondation doit rester vérifiable par un agent dont l'identité et le transport relèvent d'autres fondations, sous peine de rouvrir au niveau institutionnel ce que les protocoles ont résolu au niveau technique. L'horloge post-quantique impose symétriquement de traiter l'algorithme cryptographique comme un détail d'implémentation découplé du contrat, substituable sans rupture des signatures et justificatifs vérifiables : l'invariant *évolution* poussé à sa conséquence cryptographique.

### 9.4 Trajectoire de la menace et programme de recherche

Un seuil a été franchi : une campagne d'intrusion désignée GTG-1002, conduite à environ 80-90 % de façon autonome, a été documentée par Anthropic en novembre 2025 — un agent orchestrant reconnaissance, exploitation et exfiltration via des outils exposés en MCP, l'humain n'intervenant qu'à des points de décision épars. Ce fait établi démontre la faisabilité d'une attaque largement autonome, non sa généralisation. La défense évolue symétriquement vers un *centre des opérations de sécurité agentique*, où attaquant et défenseur sont des systèmes agentiques et où le différentiel de vitesse devient l'enjeu ; reste ouverte la question d'une avance défensive soutenable.

Le programme de recherche que dégage la strate est structuré : le **verrou sémantique et pragmatique** (partager le sens d'une intention déléguée, non le traduire avec perte) restera, selon la proposition de la revue, largement ouvert au-delà de 2030 ; les **garanties composables** sur des chaînes de délégation hétérogènes n'ont aucune voie établie, un contrat vérifié sur un agent isolé ne garantissant pas la chaîne qui l'invoque ; la **science de l'évaluation** inter-fournisseurs manque, faute d'une métrique d'horizon de tâche *déléguée*, distincte de celle d'un agent isolé. Le principe directeur — *ce qui n'est pas mesurable n'est pas contractualisable* — relie ce programme à l'invariant de la revue.

### 9.5 Forces macroéconomiques et prospective d'entreprise

Cinq forces macro encadrent l'horizon, et aucune n'est purement technique. La *capacité* d'autonomie progresse — doublement de l'horizon de tâche fiable autour de sept mois (METR, signal mesuré) — sans lever les fragilités sur les tâches longues ni combler l'absence d'une métrique de tâche *déléguée* (§9.4). L'*économie de l'inférence* voit le coût du jeton s'effondrer, absorbé par l'explosion du calcul d'inférence, poussant vers des flottes hétérogènes multi-fournisseurs qui rendent l'interopérabilité obligatoire. La *soutenabilité* énergétique, hydrique et carbone devient une contrainte : l'Agence internationale de l'énergie projette un quasi-doublement de la consommation électrique des centres de données d'ici 2030 (PROJETÉ), sans standard pour imputer une empreinte à un appel d'agent. La *souveraineté* menace d'un « splinternet agentique » où deux agents conformes chacun dans son ressort pourraient devenir incapables de coopérer pour une raison réglementaire plutôt que technique. Enfin, les *horizons de capacité* — l'IA générale — demeurent une incertitude bornée que la revue refuse de chiffrer.

Côté entreprise, l'horizon devrait voir une sortie progressive du pilotage perpétuel, un assainissement de marché, puis la consolidation d'une organisation agentique en moyeu et rayons, où la supervision humaine passe du contrôle de chaque action au pilotage de résultats. Deux constats convergents, tous deux PROJETÉS, y poussent : plus de 40 % des projets annulés d'ici fin 2027 (Gartner, 2025) et une vaste majorité de déploiements sans effet mesurable (MIT Project NANDA, 2025) — l'échec des déploiements *isolés* est l'argument économique du décloisonnement. La recomposition du travail suit la même discipline : solde net d'emplois PROJETÉ positif d'ici 2030, contre un signal déjà mesuré de recul relatif de l'emploi des jeunes travailleurs les plus exposés, sans généralisation causale établie. Si N×M → N+M demeure un pari non chiffré (SPÉCULATIF), c'est lui qui fait de l'interopérabilité une condition de marché plutôt qu'un raffinement optionnel — l'hypothèse que l'horizon mettra à l'épreuve.

### 9.6 La finance à l'horizon : rails, paiements et risque systémique

Le vertical financier franchit, sur 2027-2032, un seuil de nature plutôt que de degré : l'agentique s'insère dans des rails de paiement instantanés, irrévocables et tokenisés. Le patron directeur, de statut PROJETÉ, maintient l'agentique en amont, l'autorisation et le règlement restant déterministes — déclinaison financière de l'invariant *découplage, contrat, évolution*.

Le calendrier concentre des échéances PROGRAMMÉES : GENIUS Act au plus tard le 18 janvier 2027 ; euro numérique en pilote au second semestre 2027, émission n'étant qu'une éventualité pour 2029 ; réception obligatoire des messages camt (ISO 20022) en novembre 2027 avant fin de coexistence à la *Standards Release* SWIFT 2028 — le standard cessant d'être un texte à deviner pour devenir un code versionné et exécutable que l'agent orchestre. La gouvernance des paiements agentiques a basculé vers la FIDO Alliance, où plusieurs rails de mandat vérifiable coexistent sans consolidation certaine. Les analyses prudentielles convergentes (risque systémique européen, stabilité financière, banques centrales) restent PROJETÉES sur l'amplification du risque par corrélation d'agents, sa modélisation non linéaire demeurant SPÉCULATIVE ; la maturation de l'autonomie financière progresse par bacs à sable réglementaires, aucun banc d'essai mature n'existant pour l'assurance ni la lutte contre le blanchiment agentique.

### 9.7 Responsabilité, assurabilité et gouvernance de l'émergence

Le verrou final est moins de connecter que de *gouverner, garantir et imputer*. Le vide de responsabilité — accentué par le retrait du projet de directive européenne sur la responsabilité en IA, qui rabat le contentieux sur le régime révisé de responsabilité du fait des produits (transposition due au 9 décembre 2026), lequel couvre mal les relations B2B et les pertes purement économiques — fait de la piste d'audit un substitut probatoire et de l'assurabilité un régulateur de fait : faute de régime tranché, le marché de la couverture fixera le prix du risque agentique. La gouvernance se heurte à un théorème d'impossibilité de l'attribution : aucun cadre ne peut satisfaire à la fois attributabilité, borne de prévisibilité, non-vacuité et complétude pour des collectifs humains-agents.

Trois risques de second ordre tempèrent la thèse de l'interopérabilité comme valeur : la **déqualification**, où le déchargement cognitif érode la reprise en main humaine ; la **concentration**, l'interopérabilité multipliant les dépendances vers un petit nombre de fournisseurs plutôt que de les disperser ; et la **sur-confiance**, ce paradoxe de l'automatisation où la fiabilité perçue relâche la vigilance qui détecterait la défaillance résiduelle.

### 9.8 Scénarios 2027-2032

Plutôt que de prédire un futur unique, la revue retient trois axes d'incertitude dont le croisement engendre l'espace des futurs plausibles (tableau 14) : *convergence* protocolaire contre *coexistence stratifiée gouvernée* (scénario jugé le plus probable, qui rejoue la bifurcation de gouvernance du §9.3) ; *intensité réglementaire élevée* contre *dérégulée*, selon la fermeté de la grappe de 2027 et la consistance du régime de responsabilité ; *agentic web ouvert* contre *fractionnement souverain* régional. S'y ajoutent des cartes maîtresses à faible probabilité mais fort impact : rupture cryptographique, incident systémique d'un réseau d'agents corrélés, discontinuité de capacité. Aucun scénario n'est prédit ; tous bornent un cône d'incertitude que la grappe de 2027 (PROGRAMMÉ) ancre à court terme, dominé par le projeté et le spéculatif en approchant 2032.

Table: Tableau 14 — Axes des scénarios 2027-2032.

| Axe d'incertitude | Pôle 1 | Pôle 2 |
| --- | --- | --- |
| Trajectoire des protocoles | Convergence en une spécification unifiée | Coexistence stratifiée gouvernée (la plus probable) |
| Intensité réglementaire | Élevée | Dérégulée |
| Topologie de l'écosystème | Agentic web ouvert et interopérable | Fractionnement souverain régional |

> **Mise en œuvre.** Six conditions, aucune acquise, gouvernent la réalisation du pari N×M → N+M (§9.5) : verrou sémantique au moins partiellement résolu (§9.4) ; identité et confiance composables en régime post-quantique (§9.3) ; certification mesurable adossée à une science de l'évaluation inter-fournisseurs (§9.4) ; régime de responsabilité comblant le vide actuel (§9.7) ; standard d'attribution d'empreinte énergétique et hydrique par appel d'agent (§9.5) ; portabilité neutre des agents, faute de quoi l'effet N+M reste capté par les exécuteurs dominants (§9.2). L'invariant se transpose, à l'horizon, au niveau institutionnel et juridique : le contrat n'est plus seulement comportemental mais juridique, l'évolution n'est plus seulement celle d'un protocole mais celle d'un régime de responsabilité, et un schéma versionné, tel ISO 20022, demeure le contrat stable sur lequel un agent peut s'appuyer sans réinterpréter.

## 10. Validation appliquée : une architecture détaillée de solution

### 10.1 Démarche et portée

Pour éprouver la transposabilité des résultats, le corpus inclut une architecture détaillée de solution (ADS) appliquée à une entreprise fictive — une coopérative financière québécoise — sur une pile technologique consolidée d'un seul fournisseur. Cette ADS projette les acquis des strates 5 et 6 sur un terrain concret soumis au cadre réglementaire canadien et québécois — résidence des renseignements personnels, risque-modèle, résilience des tiers TIC, lutte contre le blanchiment — sans qu'aucun régime ne prime. Elle constitue une validation par construction : non pas une étude de déploiement empirique longitudinale, mais une démonstration d'ingénierie *prête au déploiement* — diagrammes, contrats d'interface, exigences non fonctionnelles, topologie, manuels d'exploitation — qui vérifie que les patrons de la revue se composent en un tout cohérent.

### 10.2 Plan de contrôle obligatoire et dorsale tri-plan

L'ADS traite la plateforme d'agents non comme un dispositif à gouverner après coup, mais comme un **plan de contrôle obligatoire** couplé à une **dorsale d'intégration tri-plan**. Sa logique tient en une phrase : l'agent prépare ; un humain ou un contrôle déterministe engage l'irréversible ; toute action transite par un point d'application de politique unique ; tout actif décisionnel est un modèle inventorié. Ce plan de contrôle regroupe des fonctions indissociables — autorisation, registre de modèles, identité non humaine, audit — et applique le principe de non-auto-autorisation (section 6.3) par deux points d'application de politique en série : couche agent-à-agent, puis couche agent-outil. La dorsale, elle, sépare trois substrats non interchangeables — synchrone (passerelles d'API), commande avec garantie de livraison (files de messages), événement et coordination (flux) — au-dessus desquels se superposent les deux plans agentiques (agent-outil par MCP, agent-agent par A2A). Le couple plan de contrôle/dorsale instancie, à l'échelle d'un système complet, la distinction plan de contrôle / plan de données héritée du maillage de services (section 3.6) ; la segmentation de la dorsale en trois substrats en est un raffinement propre à l'agentique, non une redite.

La règle de décision de substrat illustre la rigueur exigée : une interaction synchrone à réponse immédiate emprunte la passerelle d'API ; une commande à garantie de livraison emprunte la file de messages ; un événement ou une coordination empruntent le flux. Confondre leurs garanties — traiter un événement comme une commande, par exemple — est une source classique de défaillance. Le degré d'autonomie consenti à l'agent obéit au même souci de vérifiabilité : l'ADS instancie le patron autonomie graduée sous contrôle de finalité (section 7.2) par une échelle de paliers où l'irréversible reste borné par défaut et où tout franchissement exige une condition explicite — réversibilité native ou mandat cryptographique borné —, non la seule performance du modèle. Les exigences non fonctionnelles, enfin, sont traduites en cibles mesurables par plan plutôt qu'en un objectif global trompeur : objectifs de niveau de service, points et délais de reprise, budgets de latence, débits soutenus — des *hypothèses à calibrer* sur les volumétries réelles avant tout engagement. Un pilote sans exigence mesurable est une intention, pas une exigence : c'est cette discipline qui distingue un livrable d'ingénierie d'une intention d'architecture, et qui rend l'ADS opposable.

### 10.3 Le quatrième terme : l'exploitation

La contribution propre de l'ADS au cadre de la revue est l'extension de l'invariant d'un quatrième terme : l'**exploitation**. La thèse est que l'interopérabilité agentique n'est durable que couplée à l'exploitation résiliente qui la maintient — un parc n'est conforme que s'il demeure observable et remédiable sans rompre ses propres invariants. L'ADS distingue à cet effet deux disciplines complémentaires : gouverner *ce qui* tourne (cycle de vie de l'agent, registre, calibrage d'autonomie) et maintenir *que ça tourne* (santé d'exploitation du parc, par une boucle observer-corréler-décider-agir). La haute disponibilité et la reprise après sinistre absorbent une panne survenue ; ce volet réactif se complète d'une détection précoce d'anomalies (dérive de latence, saturation d'une file) qui anticipe la bascule plutôt que d'y répondre après coup. L'auto-remédiation qui en résulte reste bornée et *fail-closed* : elle n'engage jamais l'irréversible sans reprise humaine, et se limite au réversible (rééquilibrage, mise à l'échelle, rétrogradation d'autonomie). La récursivité est assumée : les agents d'exploitation sont eux-mêmes inscrits au registre et soumis aux mêmes contrôles d'identité, de périmètre et d'audit qu'ils appliquent au reste du parc. *Découplage, contrat, évolution* deviennent ainsi *découplage, contrat, évolution, exploitation*.

### 10.4 Compromis structurants et gouvernance des décisions

Consolider l'architecture sur une pile unique maximise la cohérence de gouvernance — un même graphe relie modèles, contrôles et exigences — et permet d'exploiter le socle applicatif existant sans le réécrire ; la revue en assume le revers. La dépendance à un fournisseur unique concentre le risque de verrouillage, que l'invariant d'*évolution* impose de compenser en maintenant ouvertes les coutures protocolaires (MCP, A2A, standards de données, échange de modèles ArchiMate). Chaque compromis structurant est ainsi documenté comme une décision d'architecture assortie d'une alternative et d'une condition explicite de renversement, non comme un choix acquis — application directe, au plan de l'ingénierie, de la discipline d'évolution gouvernée de la section 3.4.

### 10.5 Statut illustratif

L'ADS est en outre explicitement *illustrative* et ne doit pas être lue comme un état général de disponibilité de produits : plusieurs briques clés du plan de contrôle et de l'exploitation y sont en avant-première plutôt qu'en disponibilité générale, les extraits de configuration ne contiennent aucun secret réel, et toutes les valeurs chiffrées — niveaux de service, délais de reprise, débits, dimensionnement du parc — sont des hypothèses à calibrer. Ces réserves, avec les dépendances envers des tiers TIC critiques, sont isolées dans un registre dédié plutôt que diluées dans le corps du texte. Cette discipline de statut — ne jamais présenter une avant-première comme acquise, ni une hypothèse de dimensionnement comme un engagement — est la transposition, au plan de l'ingénierie, du tri épistémique de la section 9 : une architecture détaillée de solution est crédible non pas lorsqu'elle affirme des capacités, mais lorsqu'elle expose ses propres conditions de réfutation.

## 11. Discussion

### 11.1 Résultat central : l'invariant réinstancié

Le principal résultat de cette revue est la continuité de l'invariant *découplage, contrat, évolution* du transport d'octets jusqu'à la négociation entre agents. Ce qui change d'une strate à l'autre n'est pas le principe, mais ses *réalisations* : le contrat prend successivement la forme d'une description d'interface, d'un schéma de message, d'un schéma d'outil, d'une carte d'agent, d'une politique d'autorisation, puis d'un élément d'architecture formalisé. La couche agentique n'invente donc pas une nouvelle théorie de l'interopérabilité ; elle réinstancie l'ancienne à un niveau d'abstraction supérieur. C'est une **réinstanciation sous contrainte** : les principes demeurent, les mécanismes se renouvellent, et la tâche de l'architecte est de reconnaître l'invariant sous la nouveauté, afin de ne pas reconstruire à grands frais ce que des décennies d'intégration d'entreprise avaient déjà appris à découpler, à contractualiser et à faire évoluer.

Un second résultat, plus discret mais tout aussi structurant, complète ce constat : l'invariant ne se contente pas de se réinstancier, il se *complète*. À l'échelle du parc d'agents (§6.10), puis à l'épreuve de la validation appliquée (§10.3), la revue retient un quatrième terme — l'**exploitation** — car un parc conforme sur le papier n'est gouvernable que s'il demeure, en continu, observable, remédiable et auditable. Le triptyque devient ainsi un ensemble à quatre termes solidaires : *découplage, contrat, évolution, exploitation*. Cette extension n'est pas ajoutée de l'extérieur ; elle procède du même raisonnement — un parc d'agents non observé est aussi peu gouvernable qu'un contrat non vérifié.

### 11.2 Une lecture unifiée de littératures disjointes

La revue raccorde sous un invariant unique trois corpus longtemps disjoints. La littérature de l'intégration d'entreprise (sections 3, 6) apporte la théorie des niveaux, du couplage et des contrats, ainsi que les patrons d'architecture. La littérature de l'ingénierie agentique (section 4) apporte la boucle de contrôle, l'outillage, la mémoire et les modes de défaillance propres aux systèmes probabilistes. La littérature protocolaire émergente (section 5) apporte les mécanismes concrets de la coopération inter-agents. La grille d'analyse par axes — découverte, sémantique, identité et confiance, gouvernance des frontières — permet de comparer des spécifications hétérogènes et d'anticiper leurs recouvrements et leurs manques, là où l'analyse par produit isolé échoue. Le patron sectoriel de l'autonomie graduée sous contrôle de finalité et la formalisation ArchiMate complètent l'apport en rendant la synthèse, respectivement, transférable et opposable ; l'architecture détaillée de solution (section 10) referme la boucle en l'éprouvant par construction et en y ajoutant l'exploitation comme quatrième terme.

### 11.3 Tensions transversales et compromis

La revue confirme que les trois tensions irréductibles — hétérogénéité, autonomie, distribution — resurgissent intactes au niveau agentique, et qu'elles imposent les mêmes arbitrages, à un coût d'erreur supérieur. Une déclinaison de la tension de distribution mérite d'être nommée séparément, car elle structure à la fois la section 5 et la section 6 : la sûreté n'est pas une propriété compositionnelle — un agent sûr et un outil sûr, composés, ne donnent pas un système sûr —, ce qui déplace la frontière de confiance de chaque système vers chaque arête du graphe d'interaction. Le compromis principal de toute la démarche est assumé : privilégier la *gouvernabilité* (découplage explicite, contrats vérifiables, traçabilité) au prix d'une part de la fluidité et de l'autonomie que la technologie rendrait possible. L'alternative — laisser les agents négocier librement en langage naturel sans contrat opposable — maximise la souplesse mais rend la sécurité non composable et la responsabilité indécidable. Les conditions qui renverseraient la recommandation existent : un contexte non régulé, à faible coût d'erreur et à faible exigence d'audit, justifierait une autonomie plus large et une gouvernance plus légère. La finance régulée se situe à l'exact opposé de ces conditions, ce qui explique qu'elle serve de cas limite révélateur.

Un second compromis, plus subtil, traverse la revue : celui entre la *fraîcheur* et la *stabilité* de l'analyse. Traiter un objet en cours de stabilisation — une couche protocolaire dont les spécifications évoluent au rythme de plusieurs révisions par an — expose à l'obsolescence rapide ; mais attendre la stabilisation reviendrait à se taire au moment où les décisions d'architecture se prennent. La revue tranche en faveur de l'utilité présente, au prix d'une discipline de datation stricte (le socle de juin 2026) et du marqueur des ressources vivantes, qui circonscrit ce qui devra être reconfirmé. Ce choix méthodologique est lui-même une application de l'invariant : le contenu daté est *découplé* du cadre d'analyse stable (l'invariant et les axes), de sorte qu'une révision millésimée puisse mettre à jour les faits sans rejouer la démonstration. C'est cette architecture de l'argument — un noyau conceptuel stable, une périphérie factuelle versionnée — qui donne à la revue sa durée de vie utile au-delà de la péremption de tel ou tel chiffre.

### 11.4 Limites de la revue

Quatre réserves circonscrivent la portée des résultats. La première tient à la **fraîcheur de l'objet** : la couche d'interopérabilité agentique se stabilise encore, et un socle documentaire arrêté à juin 2026 vieillira vite — certaines spécifications analysées ici évolueront, fusionneront ou disparaîtront. La deuxième tient à l'**absence de validation empirique à grande échelle** : la démonstration est conceptuelle et architecturale, appuyée sur les sources et les normes, mais n'inclut pas d'étude de déploiement longitudinale en production ; l'ADS de la section 10, explicitement illustrative, atténue cette limite sans la lever. La troisième tient au **périmètre sectoriel** : la finance, choisie pour son exigence réglementaire, n'épuise pas les contextes d'usage, et la transposition à d'autres secteurs reste à éprouver. La quatrième, inhérente à la prospective, est que toute projection demeure révisable, malgré le tri épistémique. La revue assume en outre une limite de méthode : synthèse critique en spirale, elle privilégie l'intégration conceptuelle sur l'exhaustivité dénombrée d'une revue systématique.

### 11.5 Fronts de recherche ouverts

Ces limites dessinent en creux un programme. Sur le plan théorique, le rapprochement entre les formalismes de compatibilité comportementale (types de session, automates d'interface) et les contrats d'agents — encore décrits en langage semi-structuré ou spécifiés de manière probabiliste — demeure un chantier actif. Sur le plan de la sécurité, la modélisation des menaces propres aux frontières agentiques (injection par les outils et les données, confusion de délégation, dérive d'autonomie, rug-pull) appelle des cadres de vérification outillés et une notion de garantie *composable*. Sur le plan de la gouvernance, l'identité non humaine et la traçabilité de bout en bout des décisions déléguées au-delà de deux sauts restent des problèmes ouverts, dont l'urgence croît avec le degré d'autonomie consenti. Sur le plan de l'évolution, enfin, l'horloge post-quantique et la grappe d'échéances réglementaires de la fin de la décennie imposeront une discipline de migration que les architectures décrites ici devront absorber sans rompre les contrats existants. À ceux-là s'ajoute le besoin d'une *science de l'évaluation* inter-fournisseurs, condition de toute certification, et d'une métrique d'horizon de tâche déléguée. Le domaine financier laisse enfin en héritage deux fronts qui lui sont propres et que rien, à juin 2026, n'outille encore : l'assurabilité d'un préjudice causé par un agent, faute de série actuarielle, et le risque systémique d'agents corrélés sur un petit nombre de modèles de fondation communs, qui relève d'une surveillance macroprudentielle encore à construire.

### 11.6 Synthèse opérationnelle : l'invariant strate par strate

Le tableau 15 condense le résultat central de la revue : à chaque strate, l'invariant se réalise sous une forme concrète, et un verrou dominant reste ouvert. La lecture verticale du tableau donne la trajectoire du *contrat* — d'une description d'interface à un régime de responsabilité, puis à un engagement d'exploitation — et celle des questions ouvertes — d'un problème d'ingénierie à un problème institutionnel. C'est cette continuité, et non un mécanisme particulier, qui constitue l'apport de fond : reconnaître l'invariant sous la nouveauté permet à l'architecte de réutiliser un demi-siècle d'acquis plutôt que de le redécouvrir à grands frais.

Table: Tableau 15 — Réalisation de l'invariant et verrou dominant, par strate.

| Strate | Réalisation du contrat | Verrou dominant à juin 2026 |
| --- | --- | --- |
| Interopérabilité des SI (§3) | Description d'interface, schéma de message | Sémantique automatisée fiable |
| Ingénierie agentique (§4) | Schéma d'outil, configuration d'agent | Fiabilité sur les tâches à horizon long |
| Couche agentique (§5) | Carte d'agent, contrat comportemental probabiliste (§5.6) | Accord pragmatique sur l'intention |
| Entreprise (§6) | Politique d'autorisation, plan de contrôle | Identité non humaine et délégation multi-saut |
| Finance (§7) | Double contrat (modèle + tiers TIC), mandat vérifiable | Responsabilité, équité, risque systémique |
| Architecture (§8) | Élément ArchiMate stéréotypé, traçabilité verticale | Représentation du non-déterminisme et de l'exécution |
| Horizon (§9) | Contrat juridique, régime de responsabilité | Garanties composables, science de l'évaluation |
| Validation appliquée (§10) | Plan de contrôle obligatoire, dorsale tri-plan, exploitation (4ᵉ terme) | Calibrage des hypothèses chiffrées sur le réel |

## 12. Conclusion

Au terme de ce parcours en spirale, la proposition centrale se trouve éprouvée à chacune des strates traversées : l'interopérabilité durable, qu'elle relie des applications d'entreprise ou des agents d'intelligence artificielle, repose sur le même socle solidaire — *découplage, contrat, évolution* —, que la validation appliquée (section 10) complète d'un quatrième terme, l'**exploitation**, sans lequel un parc d'agents interopérable demeure, à l'échelle, ingouvernable. De la théorie de l'interopérabilité (section 3) à l'ingénierie des agents (section 4), des protocoles émergents (section 5) à leur déploiement en entreprise (section 6), de la spécialisation financière (section 7) à la formalisation ArchiMate (section 8), jusqu'à la projection sur l'horizon 2027-2032 (section 9) et à sa validation appliquée (section 10), c'est le même invariant, désormais à quatre termes, qui se réinstancie sous des mécanismes renouvelés.

La revue dégage cinq apports : une lecture unifiée de littératures disjointes ; une grille d'analyse des protocoles par les axes éprouvés de l'interopérabilité ; un patron de conception sectoriel articulant autonomie, criticité et exigence réglementaire ; une formalisation architecturale opposable ; et une validation appliquée qui éprouve l'ensemble par construction tout en étendant l'invariant à l'exploitation. Elle confirme aussi que le verrou de la fin de la décennie n'est plus d'abord technique mais institutionnel et économique : il s'agit moins de *connecter* que de *gouverner, garantir, imputer et exploiter*. L'interopérabilité agentique n'est ni une rupture qui invaliderait l'ingénierie de l'intégration, ni une simple continuation qui la laisserait inchangée. Elle est une réinstanciation sous contrainte — et la responsabilité de l'architecte est précisément de reconnaître l'invariant sous la nouveauté. La clôture honnête d'une telle revue consiste, dès lors, à rouvrir : les ancrages les plus fermes de l'horizon — la grappe réglementaire de 2027 et la migration post-quantique de 2030-2035 — bornent un cône d'incertitude que la recherche et la pratique auront, ensemble, à explorer.

## Annexe A — Sigles et acronymes

| Sigle | Signification |
| --- | --- |
| AAIF | Agentic AI Foundation (Linux Foundation, gouvernance de MCP depuis le 9 décembre 2025) |
| A2A | Agent2Agent (protocole d'interopérabilité entre agents) |
| ACID | Atomicity, Consistency, Isolation, Durability (propriétés transactionnelles) |
| ACORD | Association for Cooperative Operations Research and Development (normes d'assurance) |
| ACP | Agent Communication Protocol / Agentic Commerce Protocol (homonymes) |
| AG-UI | Agent-User Interaction Protocol |
| AGNTCY | Collectif d'interopérabilité agentique (OASF) |
| AI Act | Règlement (UE) 2024/1689 sur l'intelligence artificielle |
| AL3 | Norme de transactions de données ACORD pour l'assurance IARD (Automation Level 3) |
| AMF | Autorité des marchés financiers (régulateur financier québécois) |
| AMLR | Anti-Money Laundering Regulation (Règlement (UE) 2024/1624, lutte contre le blanchiment) |
| ANP | Agent Network Protocol |
| AP2 | Agent Payments Protocol |
| API | Application Programming Interface |
| BCP | Best Current Practice (IETF) |
| BDI | Belief–Desire–Intention (architecture d'agent) |
| BIAN | Banking Industry Architecture Network |
| BOLA | Broken Object Level Authorization |
| BOPLA | Broken Object Property Level Authorization |
| CAP | Cohérence, disponibilité, tolérance au partitionnement (théorème CAP) |
| CDC | Change Data Capture (capture de changement) |
| CDM | Common Domain Model (ISDA) |
| CoALA | Cognitive Architectures for Language Agents |
| DID | Decentralized Identifier (W3C) |
| DORA | Digital Operational Resilience Act (Règlement (UE) 2022/2554) |
| DPoP | Demonstrating Proof of Possession |
| EDI | Electronic Data Interchange |
| EDIFACT | Electronic Data Interchange for Administration, Commerce and Transport |
| eIDAS | Electronic Identification, Authentication and trust Services (règlement européen) |
| EIMM | Enterprise Interoperability Maturity Model |
| EIRA | European Interoperability Reference Architecture |
| ESB | Enterprise Service Bus |
| FAPI | Financial-grade API |
| FDX | Financial Data Exchange |
| FIBO | Financial Industry Business Ontology |
| FIDO | Fast IDentity Online (Alliance) |
| FINOS | Fintech Open Source Foundation |
| FIPA | Foundation for Intelligent Physical Agents |
| FIX | Financial Information eXchange |
| FpML | Financial products Markup Language |
| GENIUS Act | Guiding and Establishing National Innovation for U.S. Stablecoins Act (loi américaine, 2025) |
| gRPC | gRPC Remote Procedure Calls (protocole RPC de Google) |
| HATEOAS | Hypermedia as the Engine of Application State |
| HTTP | Hypertext Transfer Protocol |
| IARD | Incendie, accidents et risques divers (assurance de dommages) |
| IETF | Internet Engineering Task Force |
| iPaaS | Integration Platform as a Service |
| ISDA | International Swaps and Derivatives Association |
| ISO | International Organization for Standardization |
| JOUE | Journal officiel de l'Union européenne |
| JSON | JavaScript Object Notation |
| JWT | JSON Web Token |
| KQML | Knowledge Query and Manipulation Language |
| KYA | Know Your Agent (connaissance de l'agent) |
| LCIM | Levels of Conceptual Interoperability Model |
| LF | Linux Foundation |
| LISI | Levels of Information Systems Interoperability |
| LLM | Large Language Model (grand modèle de langage) |
| MAST | Multi-Agent System Failure Taxonomy |
| MCP | Model Context Protocol |
| METR | Model Evaluation and Threat Research |
| MOM | Message-Oriented Middleware (middleware orienté message) |
| mTLS | Mutual Transport Layer Security (authentification mutuelle TLS) |
| NEIF | New European Interoperability Framework |
| NHI | Non-Human Identity (identité non humaine) |
| NIST | National Institute of Standards and Technology |
| OAEI | Ontology Alignment Evaluation Initiative |
| OASF | Open Agentic Schema Framework (AGNTCY) |
| OAuth | Open Authorization (cadre d'autorisation déléguée) |
| OCC | Office of the Comptroller of the Currency (États-Unis) |
| OIDC | OpenID Connect |
| OIM | Organisational Interoperability Model (maturité organisationnelle) |
| OSFI | Office of the Superintendent of Financial Institutions (Bureau du surintendant des institutions financières, BSIF) |
| OWASP | Open Worldwide Application Security Project |
| OWL | Web Ontology Language (W3C) |
| PKCE | Proof Key for Code Exchange (OAuth) |
| PSD2 | Payment Services Directive 2 (Directive (UE) 2015/2366) |
| RAG | Retrieval-Augmented Generation |
| RC | Release Candidate (révision candidate) |
| RDF | Resource Description Framework (W3C) |
| REST | Representational State Transfer |
| RFC | Request for Comments (IETF) |
| RGPD | Règlement général sur la protection des données (UE) |
| RPC | Remote Procedure Call (appel de procédure distante) |
| SAML | Security Assertion Markup Language |
| SCIM | System for Cross-domain Identity Management |
| SHACL | Shapes Constraint Language (W3C) |
| SI | Système d'information |
| SLA | Service-Level Agreement (entente de niveau de service) |
| SLI | Service-Level Indicator (indicateur de niveau de service) |
| SLO | Service-Level Objective (objectif de niveau de service) |
| SOA | Service-Oriented Architecture |
| SOAP | Simple Object Access Protocol |
| SPARQL | SPARQL Protocol and RDF Query Language (W3C) |
| SPIFFE | Secure Production Identity Framework for Everyone |
| SR | Supervision and Regulation letter (lettre de surveillance de la Réserve fédérale américaine, p. ex. SR 26-2) |
| SSE | Server-Sent Events (événements envoyés par le serveur) |
| TIC | Technologies de l'information et de la communication |
| TOGAF | The Open Group Architecture Framework |
| UCP | Universal Commerce Protocol |
| UDDI | Universal Description, Discovery and Integration |
| W3C | World Wide Web Consortium |
| WIMSE | Workload Identity in Multi-System Environments (IETF) |
| WSDL | Web Services Description Language |
| X12 | ASC X12 — Accredited Standards Committee X12 (norme EDI nord-américaine) |
| x402 | Protocole de paiement machine-natif (HTTP 402) |

## Note sur la méthode et les sources

Cet article est une revue de synthèse d'une monographie de l'auteur — Bruneau, A.-G. (2026), *Interopérabilité Agentique en entreprise dans le domaine des services financiers*, monographie en science et génie informatique — dont il consolide les sept chapitres et l'architecture détaillée de solution sous une forme d'article scientifique. Les références ci-dessous sont reprises *verbatim* des bibliographies par chapitre de la monographie, ayant fait l'objet d'une double passe de vérification adverse close en juin 2026 (existence, auteurs, année, numéro de norme, de RFC ou de recommandation). Le socle documentaire est arrêté à juin 2026 ; le marqueur ⚠ signale une *ressource vivante* (spécification versionnée par date, texte réglementaire en cours d'adoption, rapport d'analyste) dont la version et l'ancrage doivent être reconfirmés avant toute diffusion formelle. Bilan : la bibliographie compte 211 entrées réparties en sept catégories (A à G) ; la révision du 14 juillet 2026 y a importé 25 entrées issues de ces mêmes bibliographies de chapitre, afin que chaque mention datée du corps dispose de son entrée, sans introduire de référence nouvelle non vérifiée. Conformément à l'usage, l'auteur déclare avoir employé des outils d'intelligence artificielle générative, sous sa supervision et sa responsabilité, pour la structuration, la rédaction de premiers jets et la vérification des références ; il assume l'entière responsabilité scientifique du contenu.

## Références

Les références sont regroupées par catégorie thématique. Au sein de chaque catégorie, l'ordre est alphabétique par auteur ou organisme ; dans les catégories A, D, E et F, un sous-classement place les travaux de personnes avant les sources normatives et institutionnelles (organismes de normalisation, régulateurs, législateurs), regroupées en fin de catégorie ; la catégorie B distingue de même les travaux fondateurs (agents et systèmes multi-agents classiques) des travaux contemporains. Les normes ISO/IEC sont ordonnées par numéro croissant, de même que les RFC d'un même auteur.

### A. Théorie de l'interopérabilité et intégration d'entreprise

- Akkoyunlu, E. A.; Ekanadham, K.; Huber, R. V. (1975). *Some Constraints and Tradeoffs in the Design of Network Communications*. Proceedings of the 5th ACM Symposium on Operating Systems Principles (SOSP'75), ACM, p. 67-74. DOI 10.1145/800213.806523.
- Brewer, E. (2012). *CAP Twelve Years Later: How the "Rules" Have Changed*. Computer (IEEE), vol. 45, no 2, p. 23-29. DOI 10.1109/MC.2012.37.
- Chappell, D. A. (2004). *Enterprise Service Bus*. O'Reilly Media (série Theory in Practice). ISBN 978-0-596-00675-4.
- Chen, D.; Daclin, N. (2006). *Framework for Enterprise Interoperability*. Proceedings of the IFAC Workshop EI2N, Bordeaux, p. 77-88 ; aussi chapitre 6 dans *Interoperability for Enterprise Software and Applications* (dir. H. Panetto, N. Boudjlida), ISTE/Wiley. DOI 10.1002/9780470612200.ch6.
- Chen, D.; Doumeingts, G.; Vernadat, F. (2008). *Architectures for enterprise integration and interoperability: Past, present and future*. Computers in Industry, vol. 59, no 7, p. 647-659. Elsevier. DOI 10.1016/j.compind.2007.12.016.
- de Alfaro, L.; Henzinger, T. A. (2001). *Interface Automata*. Proceedings of ESEC/FSE-9, ACM, p. 109-120. DOI 10.1145/503209.503226.
- Erl, T. (2005). *Service-Oriented Architecture: Concepts, Technology, and Design*. Prentice Hall PTR (Pearson Education). ISBN 978-0-13-185858-9.
- Euzenat, J.; Shvaiko, P. (2013). *Ontology Matching* (2nd Edition). Springer-Verlag, Heidelberg. ISBN 978-3-642-38720-3. DOI 10.1007/978-3-642-38721-0.
- Evans, E. (2003). *Domain-Driven Design: Tackling Complexity in the Heart of Software*. Addison-Wesley Professional. ISBN 978-0-321-12521-7.
- Fielding, R. T. (2000). *Architectural Styles and the Design of Network-based Software Architectures*. Dissertation, University of California, Irvine. https://www.ics.uci.edu/~fielding/pubs/dissertation/top.htm
- Fischer, M. J.; Lynch, N. A.; Paterson, M. S. (1985). *Impossibility of Distributed Consensus with One Faulty Process*. Journal of the ACM, vol. 32, no 2, p. 374-382. DOI 10.1145/3149.214121.
- Fowler, M. (2010). *Richardson Maturity Model*. martinfowler.com (bliki), 18 mars 2010. https://martinfowler.com/articles/richardsonMaturityModel.html.
- Gilbert, S.; Lynch, N. (2002). *Brewer's Conjecture and the Feasibility of Consistent, Available, Partition-Tolerant Web Services*. ACM SIGACT News, vol. 33, no 2, p. 51-59. DOI 10.1145/564585.564601.
- Gruber, T. R. (1993). *A Translation Approach to Portable Ontology Specifications*. Knowledge Acquisition, vol. 5, no 2, p. 199-220. DOI 10.1006/knac.1993.1008.
- Hohpe, G.; Woolf, B. (2003). *Enterprise Integration Patterns: Designing, Building, and Deploying Messaging Solutions*. Addison-Wesley (Signature Series, M. Fowler). ISBN 978-0-321-20068-6.
- Honda, K.; Yoshida, N.; Carbone, M. (2008). *Multiparty Asynchronous Session Types*. Proceedings of POPL'08, ACM, p. 273-284. DOI 10.1145/1328438.1328472.
- Kleppmann, M. (2017). *Designing Data-Intensive Applications: The Big Ideas Behind Reliable, Scalable, and Maintainable Systems*. O'Reilly Media. ISBN 978-1-4493-7332-0.
- Lamport, L. (1978). *Time, Clocks, and the Ordering of Events in a Distributed System*. Communications of the ACM, vol. 21, no 7, p. 558-565. DOI 10.1145/359545.359563.
- Lewis, J.; Fowler, M. (2014). *Microservices: a definition of this new architectural term*. martinfowler.com, 25 mars 2014.
- Newman, S. (2021). *Building Microservices: Designing Fine-Grained Systems* (2e édition). O'Reilly Media. ISBN 978-1-492-03402-5.
- Parnas, D. L. (1972). *On the Criteria To Be Used in Decomposing Systems into Modules*. Communications of the ACM, vol. 15, no 12, p. 1053-1058. DOI 10.1145/361598.361623.
- Richardson, C. (2018). *Microservices Patterns: With Examples in Java*. Manning Publications. ISBN 978-1-617-29454-9.
- Sheth, A. P.; Larson, J. A. (1990). *Federated Database Systems for Managing Distributed, Heterogeneous, and Autonomous Databases*. ACM Computing Surveys, vol. 22, no 3, p. 183-236. DOI 10.1145/96602.96604.
- Tolk, A.; Muguira, J. A. (2003). *The Levels of Conceptual Interoperability Model (LCIM)*. Fall Simulation Interoperability Workshop (SIW), SISO, Orlando, paper 03F-SIW-007.
- Wang, W.; Tolk, A.; Wang, W. (2009). *The Levels of Conceptual Interoperability Model: Applying Systems Engineering Principles to Modeling and Simulation*. Proceedings of SpringSim'09. arXiv:0908.0191.
- Wegner, P. (1996). *Interoperability*. ACM Computing Surveys, vol. 28, no 1, p. 285-287. DOI 10.1145/234313.234424.
- AsyncAPI Initiative (Linux Foundation) (2023). *AsyncAPI Specification 3.0.0*, 5 décembre 2023. https://www.asyncapi.com/docs/reference/specification/v3.0.0.
- AsyncAPI Initiative (Linux Foundation) (2026). *AsyncAPI Specification 3.1.0*, 31 janvier 2026. https://www.asyncapi.com/docs/reference/specification/v3.1.0.
- Commission européenne (DG DIGIT / ISA2) (2017). *New European Interoperability Framework* (Annexe 2 à COM(2017) 134 final), 23 mars 2017. CELEX:52017DC0134.
- Commission européenne (Digital Europe / Interoperable Europe) (2024). *European Interoperability Reference Architecture (EIRA) v6.1.0*, 31 mai 2024. https://interoperable-europe.ec.europa.eu/collection/european-interoperability-reference-architecture-eira/solution/eira/eira-v610-online-documentation.
- ISO (2011). *Advanced automation technologies and their applications — Requirements for establishing manufacturing enterprise process interoperability — Part 1: Framework for enterprise interoperability*. ISO 11354-1:2011. https://www.iso.org/standard/50417.html.
- ISO/IEC (2023). *Systems and software engineering — SQuaRE — Product quality model*. ISO/IEC 25010:2023. ⚠
- OpenAPI Initiative (Linux Foundation) (2021). *OpenAPI Specification v3.1.0*, 18 février 2021. https://spec.openapis.org/oas/v3.1.0.html.
- OpenAPI Initiative (Linux Foundation) (2025). *OpenAPI Specification v3.2.0*, 19 septembre 2025. https://spec.openapis.org/oas/v3.2.0.html.
- Parlement européen et Conseil de l'UE (2024). *Règlement (UE) 2024/903 du 13 mars 2024 établissant des mesures pour un niveau élevé d'interopérabilité du secteur public (acte pour une Europe interopérable)*. CELEX:32024R0903.

### B. Ingénierie des systèmes agentiques (IA)

- Bratman, M. E. (1987). *Intention, Plans, and Practical Reason*. Harvard University Press, Cambridge MA. ISBN 9780674458185.
- Brooks, R. A. (1991). *Intelligence without Representation*. Artificial Intelligence, vol. 47, no 1-3, p. 139-159. DOI 10.1016/0004-3702(91)90053-M.
- Rao, A. S.; Georgeff, M. P. (1995). *BDI Agents: From Theory to Practice*. First International Conference on Multi-Agent Systems (ICMAS-95), San Francisco, p. 312-319. AAAI Digital Library (aaai.org/Library/ICMAS/icmas95contents.php).
- Russell, S. J.; Norvig, P. (2020/2021). *Artificial Intelligence: A Modern Approach* (4th Edition). Pearson. ISBN 9780134610993.
- Smith, R. G. (1980). *The Contract Net Protocol: High-Level Communication and Control in a Distributed Problem Solver*. IEEE Transactions on Computers, vol. C-29, no 12, p. 1104-1113. DOI 10.1109/TC.1980.1675516.
- Sutton, R. S.; Barto, A. G. (2018). *Reinforcement Learning: An Introduction* (2nd Edition). MIT Press, Cambridge MA. ISBN 9780262039246.
- Wooldridge, M. (2009). *An Introduction to MultiAgent Systems* (2nd Edition). John Wiley & Sons. ISBN 9780470519462.
- Wooldridge, M.; Jennings, N. R. (1995). *Intelligent Agents: Theory and Practice*. The Knowledge Engineering Review, vol. 10, no 2, p. 115-152. DOI 10.1017/S0269888900008122.
- Anthropic (2024). *Building Effective Agents*. Anthropic Engineering (19 déc. 2024). https://www.anthropic.com/engineering/building-effective-agents
- Anthropic (2025). *How we built our multi-agent research system*. Anthropic Engineering (13 juin 2025). https://www.anthropic.com/engineering/multi-agent-research-system
- Cemri, M.; Pan, M. Z.; Yang, S.; et al. (2025). *Why Do Multi-Agent LLM Systems Fail?* arXiv:2503.13657 (taxonomie MAST). ⚠
- Debenedetti, E.; Shumailov, I.; Fan, T.; et al. (2025). *Defeating Prompt Injections by Design (CaMeL)*. arXiv:2503.18813.
- DeepSeek-AI (Guo, D.; et al.) (2025). *DeepSeek-R1: Incentivizing Reasoning Capability in LLMs via Reinforcement Learning*. arXiv:2501.12948 ; aussi Nature vol. 645, p. 633-638, 2025, DOI 10.1038/s41586-025-09422-z.
- Edge, D.; Trinh, H.; Cheng, N.; et al. (2024). *From Local to Global: A Graph RAG Approach to Query-Focused Summarization*. arXiv:2404.16130 (Microsoft Research, GraphRAG).
- Greenblatt, R.; et al. (2024). *Alignment Faking in Large Language Models*. arXiv:2412.14093 (Anthropic & Redwood Research).
- Guo, T.; Chen, X.; Wang, Y.; et al. (2024). *Large Language Model based Multi-Agents: A Survey of Progress and Challenges*. arXiv:2402.01680 (IJCAI 2024, p. 8048-8057).
- Hong, S.; Zheng, X.; Chen, J.; et al. (2023). *MetaGPT: Meta Programming for a Multi-Agent Collaborative Framework*. arXiv:2308.00352 (ICLR 2024).
- Jimenez, C. E.; Yang, J.; Wettig, A.; et al. (2023). *SWE-bench: Can Language Models Resolve Real-World GitHub Issues?* arXiv:2310.06770 (ICLR 2024).
- Kambhampati, S.; Valmeekam, K.; Guan, L.; et al. (2024). *LLMs Can't Plan, But Can Help Planning in LLM-Modulo Frameworks*. arXiv:2402.01817 (ICML 2024, spotlight).
- Kapoor, S.; Stroebl, B.; Siegel, Z. S.; Nadgir, N.; Narayanan, A. (2024). *AI Agents That Matter*. arXiv:2407.01502.
- Lewis, P.; Perez, E.; Piktus, A.; et al. (2020). *Retrieval-Augmented Generation for Knowledge-Intensive NLP Tasks*. arXiv:2005.11401 (NeurIPS 2020).
- Meinke, A.; Schoen, B.; Scheurer, J.; et al. (2024). *Frontier Models are Capable of In-Context Scheming*. arXiv:2412.04984 (Apollo Research).
- METR (Model Evaluation & Threat Research) (2025). *Recent Frontier Models Are Reward Hacking*. METR (5 juin 2025). https://metr.org/blog/2025-06-05-recent-reward-hacking/.
- Park, J. S.; O'Brien, J. C.; Cai, C. J.; et al. (2023). *Generative Agents: Interactive Simulacra of Human Behavior*. arXiv:2304.03442 (UIST 2023).
- Patil, S. G.; Zhang, T.; Wang, X.; Gonzalez, J. E. (2023). *Gorilla: Large Language Model Connected with Massive APIs*. arXiv:2305.15334 (NeurIPS 2024).
- Qin, Y.; Liang, S.; Ye, Y.; et al. (2023). *ToolLLM: Facilitating Large Language Models to Master 16000+ Real-world APIs*. arXiv:2307.16789 (ICLR 2024).
- Sapkota, R.; Roumeliotis, K. I.; Karkee, M. (2025). *AI Agents vs. Agentic AI: A Conceptual Taxonomy, Applications and Challenges*. arXiv:2505.10468 ; Information Fusion, DOI 10.1016/j.inffus.2025.103599.
- Schick, T.; Dwivedi-Yu, J.; Dessì, R.; et al. (2023). *Toolformer: Language Models Can Teach Themselves to Use Tools*. arXiv:2302.04761 (NeurIPS 2023).
- Shinn, N.; Cassano, F.; Berman, E.; Gopinath, A.; Narasimhan, K.; Yao, S. (2023). *Reflexion: Language Agents with Verbal Reinforcement Learning*. arXiv:2303.11366 (NeurIPS 2023).
- Sumers, T. R.; Yao, S.; Narasimhan, K.; Griffiths, T. L. (2023). *Cognitive Architectures for Language Agents (CoALA)*. arXiv:2309.02427 (TMLR, 2024).
- Valmeekam, K.; Marquez, M.; Olmo, A.; Sreedharan, S.; Kambhampati, S. (2023). *PlanBench: An Extensible Benchmark for Evaluating Large Language Models on Planning and Reasoning about Change*. arXiv:2206.10498 (NeurIPS 2023 Datasets & Benchmarks).
- Wang, L.; Ma, C.; Feng, X.; et al. (2024). *A Survey on Large Language Model based Autonomous Agents*. arXiv:2308.11432 ; Frontiers of Computer Science, vol. 18, no 6, DOI 10.1007/s11704-024-40231-1.
- Wei, J.; Wang, X.; Schuurmans, D.; et al. (2022). *Chain-of-Thought Prompting Elicits Reasoning in Large Language Models*. arXiv:2201.11903 (NeurIPS 2022).
- Wu, Q.; Bansal, G.; Zhang, J.; et al. (2023). *AutoGen: Enabling Next-Gen LLM Applications via Multi-Agent Conversation Framework*. arXiv:2308.08155 (COLM 2024).
- Xi, Z.; Chen, W.; Guo, X.; et al. (2023). *The Rise and Potential of Large Language Model Based Agents: A Survey*. arXiv:2309.07864.
- Yao, S.; Zhao, J.; Yu, D.; et al. (2022). *ReAct: Synergizing Reasoning and Acting in Language Models*. arXiv:2210.03629 (ICLR 2023).
- Yao, S.; Yu, D.; Zhao, J.; et al. (2023). *Tree of Thoughts: Deliberate Problem Solving with Large Language Models*. arXiv:2305.10601 (NeurIPS 2023).

### C. Protocoles d'interopérabilité agentique (MCP, A2A, ANP, commerce/paiement)

- A2A Project / Linux Foundation (2026). *Agent2Agent (A2A) Protocol Specification* (v1.0.0 le 12 mars 2026 ; v1.0.1 le 28 mai 2026). a2a-protocol.org. ⚠
- Agent Network Protocol Community (2025). *AgentNetworkProtocol (ANP) — specifications* (dont did:wba Method Design Specification). ⚠
- Anthropic (2024). *Introducing the Model Context Protocol*, 25 novembre 2024.
- Anthropic / Model Context Protocol contributors (2025). *Model Context Protocol Specification (version 2025-11-25)* (maintenu sous l'Agentic AI Foundation / Linux Foundation). ⚠
- Bhardwaj, V. P. (2026). *Agent Behavioral Contracts: Formal Specification and Runtime Enforcement for Reliable Autonomous AI Agents*. arXiv:2602.22302. ⚠
- Chang, G.; et al. (Agent Network Protocol) (2025). *Agent Network Protocol (ANP) Technical White Paper*. arXiv:2508.00007.
- Christensen, E.; Curbera, F.; Meredith, G.; Weerawarana, S. (W3C) (2001). *Web Services Description Language (WSDL) 1.1*. W3C Note, 15 mars 2001.
- Coinbase ; x402 Foundation (sous Linux Foundation) (2025-2026). *x402 — A payments protocol for the internet, built on HTTP*. ⚠
- CopilotKit (et communauté AG-UI) (2025). *AG-UI (Agent-User Interaction Protocol)* (HTTP POST + SSE, protobuf optionnel) entre front-end et backend agentique. copilotkit.ai/ag-ui ; GitHub ag-ui-protocol/ag-ui. ⚠
- Ehtesham, A.; Singh, A.; Gupta, G. K.; Kumar, S. (2025). *A survey of agent interoperability protocols: MCP, ACP, A2A, and ANP*. arXiv:2505.02279.
- Finin, T.; Fritzson, R.; McKay, D.; McEntire, R. (1994). *KQML as an Agent Communication Language*. Proceedings of CIKM '94, ACM, p. 456-463. DOI 10.1145/191246.191322.
- FIPA / IEEE Computer Society (2002). *FIPA ACL Message Structure Specification (SC00061G)*; *FIPA Communicative Act Library Specification (SC00037J)*, 6 décembre 2002.
- Google (gouvernance transférée à la FIDO Alliance) (2025-2026). *Agent Payments Protocol (AP2) — Specification* (mandats Intent/Cart/Payment ; v0.2 le 2026-04-28). ⚠
- Google (avec Shopify, Etsy, Wayfair, Target, Walmart) (2026). *Universal Commerce Protocol (UCP)* (transports API/A2A/MCP ; découverte/panier/achat/suivi ; compatible AP2). ucp.dev (annonce 2026-01-11, NRF Big Show). ⚠
- IBM Research / BeeAI, Linux Foundation (2025). *Agent Communication Protocol (ACP) — projet BeeAI*. ⚠
- Marro, S.; La Malfa, E.; Wright, J.; et al. (2024). *A Scalable Communication Protocol for Networks of Large Language Models* (Agora). arXiv:2410.11905.
- Mastercard (2025). *Mastercard unveils Agent Pay, pioneering agentic payments technology to power commerce in the age of AI*. 29.4.2025. https://www.mastercard.com/global/en/news-and-trends/press/2025/april/mastercard-unveils-agent-pay-pioneering-agentic-payments-technology-to-power-commerce-in-the-age-of-ai.html.
- MCP maintainers (Anthropic, OpenAI, GitHub, Block, communauté) (2025). *MCP Specification — Revision 2025-06-18*. ⚠
- OpenAI, Stripe (2025-2026). *Agentic Commerce Protocol (ACP) — Specification* (Checkout API, Delegate Payment). ⚠
- Raskar, R.; Chari, P.; Zinky, J.; et al. (2025). *Beyond DNS: Unlocking the Internet of AI Agents via the NANDA Index and Verified AgentFacts*. arXiv:2507.14263 (MIT Media Lab / Project NANDA).
- Sander, L.; Gidey, H. K.; Lenz, A.; Knoll, A. (2026). *A Technical Taxonomy of LLM Agent Communication Protocols*. arXiv:2606.19135. ⚠
- Singh, M. P. (1998). *Agent Communication Languages: Rethinking the Principles*. IEEE Computer, vol. 31, no 12, p. 40-47. DOI 10.1109/2.735849.
- Soria Parra, D.; Delimarsky, D. — Model Context Protocol (2026). *The 2026-07-28 MCP Specification Release Candidate* (cycle Active/Deprecated/Removed ≥ 12 mois). Model Context Protocol Blog, 21 mai 2026. https://blog.modelcontextprotocol.io/posts/2026-07-28-release-candidate/. ⚠
- Sporny, M.; Longley, D.; Kellogg, G.; et al. (2020). *JSON-LD 1.1: A JSON-based Serialization for Linked Data*. W3C Recommendation, 16 juillet 2020.
- Stripe, Tempo (2026). *Machine Payments Protocol (MPP) — open standard for billing AI agents over HTTP* (HTTP 402 ; stablecoins via Tempo on-chain, cartes/BNPL via Shared Payment Tokens). mpp.dev (lancement 2026-03-18). ⚠
- The Linux Foundation (2025). *Linux Foundation Launches the Agent2Agent Protocol Project*, 23 juin 2025.
- The Linux Foundation (2025). *Linux Foundation Announces the Formation of the Agentic AI Foundation (AAIF)*, 9 décembre 2025.
- Visa (avec Cloudflare) (2025). *Trusted Agent Protocol (TAP)* (RFC 9421 / Web Bot Auth ; Ed25519 ; annonce 2025-10-14). ⚠
- W3C DID Working Group (Sporny, M.; Guy, A.; Sabadello, M.; Reed, D., Eds.) (2022). *Decentralized Identifiers (DIDs) v1.0*. W3C Recommendation, 19 juillet 2022.
- W3C VC Working Group (Sporny, M.; Thibodeau Jr, T.; Herman, I.; Cohen, G.; Jones, M. B., Eds.) (2025). *Verifiable Credentials Data Model v2.0*. W3C Recommendation, 15 mai 2025.
- Yang, Yingxuan; Chai, H.; Song, Y.; et al. (2025a). *A Survey of AI Agent Protocols* (Shanghai Jiao Tong University). arXiv:2504.16736.
- Yang, Yingxuan; Ma, M.; Huang, Y.; et al. (2025b). *Agentic Web: Weaving the Next Web with AI Agents*. arXiv:2507.21206.
- Yolum, P.; Singh, M. P. (2002). *Flexible Protocol Specification and Execution: Applying Event Calculus Planning Using Commitments*. Proceedings of AAMAS '02, ACM, p. 527-534. DOI 10.1145/544862.544867.
- Yuan, D.; Lyu, F.; Yuan, Y.; et al. (2026). *Beyond Message Passing: A Semantic View of Agent Communication Protocols*. arXiv:2604.02369. ⚠

### D. Entreprise : identité, sécurité, gouvernance, observabilité

- Backman, A. (Ed.); Richer, J. (Ed.); Sporny, M. (2024). *HTTP Message Signatures*. IETF, RFC 9421.
- Campbell, B.; Bradley, J.; Sakimura, N.; Lodderstedt, T. (2020). *OAuth 2.0 Mutual-TLS Client Authentication and Certificate-Bound Access Tokens*. IETF, RFC 8705. DOI 10.17487/RFC8705.
- Campbell, B.; Bradley, J.; Tschofenig, H. (2020). *Resource Indicators for OAuth 2.0*. IETF, RFC 8707. DOI 10.17487/RFC8707.
- Fett, D.; Campbell, B.; Bradley, J.; et al. (2023). *OAuth 2.0 Demonstrating Proof of Possession (DPoP)*. IETF, RFC 9449. DOI 10.17487/RFC9449.
- Greshake, K.; Abdelnabi, S.; Mishra, S.; et al. (2023). *Not What You've Signed Up For: Compromising Real-World LLM-Integrated Applications with Indirect Prompt Injection*. AISec '23. arXiv:2302.12173.
- Hardt, D. (Ed.) (2012). *The OAuth 2.0 Authorization Framework*. IETF, RFC 6749. DOI 10.17487/RFC6749.
- Hardt, D.; Parecki, A.; Lodderstedt, T. (Eds.) (2026). *The OAuth 2.1 Authorization Framework* (Internet-Draft). IETF OAuth WG, draft-ietf-oauth-v2-1. ⚠
- Huang, K.; Narajala, V. S.; Habler, I.; Sheriff, A. (2025). *Agent Name Service (ANS): A Universal Directory for Secure AI Agent Discovery and Interoperability*. arXiv:2505.10609.
- Hunt, P. (Ed.); Grizzle, K.; Wahlstroem, E.; Mortimore, C. (2015). *System for Cross-domain Identity Management: Core Schema*. IETF, RFC 7643. DOI 10.17487/RFC7643.
- Hunt, P. (Ed.); Grizzle, K.; Ansari, M.; Wahlstroem, E.; Mortimore, C. (2015). *System for Cross-domain Identity Management: Protocol*. IETF, RFC 7644. DOI 10.17487/RFC7644.
- Jones, M.; Bradley, J.; Sakimura, N. (2015). *JSON Web Token (JWT)*. IETF, RFC 7519. DOI 10.17487/RFC7519.
- Jones, M.; Nadalin, A.; Campbell, B. (Ed.); et al. (2020). *OAuth 2.0 Token Exchange*. IETF, RFC 8693. DOI 10.17487/RFC8693.
- Jones, M. B.; Hunt, P.; Parecki, A. (2025). *OAuth 2.0 Protected Resource Metadata*. IETF, RFC 9728. DOI 10.17487/RFC9728.
- Lodderstedt, T.; Bradley, J.; Labunets, A.; Fett, D. (2025). *Best Current Practice for OAuth 2.0 Security*. IETF, RFC 9700, BCP 240. DOI 10.17487/RFC9700.
- Richer, J. (Ed.); Jones, M.; Bradley, J.; et al. (2015). *OAuth 2.0 Dynamic Client Registration Protocol*. IETF, RFC 7591.
- Rose, S.; Borchert, O.; Mitchell, S.; Connelly, S. (2020). *Zero Trust Architecture*. NIST Special Publication 800-207. DOI 10.6028/NIST.SP.800-207.
- Sakimura, N.; Bradley, J.; Jones, M.; et al. (2014). *OpenID Connect Core 1.0*. OpenID Foundation.
- Sakimura, N. (Ed.); Bradley, J.; Agarwal, N. (2015). *Proof Key for Code Exchange by OAuth Public Clients (PKCE)*. IETF, RFC 7636. DOI 10.17487/RFC7636.
- South, T.; Marro, S.; Hardjono, T.; et al. (2025). *Authenticated Delegation and Authorized AI Agents*. arXiv:2501.09674.
- Willison, S. (2025). *The lethal trifecta for AI agents: private data, untrusted content, and external communication*. Simon Willison's Weblog (16 juin 2025).
- Cloud Security Alliance (Huang, K.) (2025). *Agentic AI Threat Modeling Framework: MAESTRO* (7 couches).
- Coalition for Secure AI (CoSAI), OASIS Open — Workstream 4 (2025). *CoSAI Principles for Secure-by-Design Agentic Systems*. Juillet 2025.
- IETF WIMSE Working Group (Salowey, J.; Rosomakho, Y.; Tschofenig, H., Eds.) (2026). *Workload Identity in a Multi System Environment (WIMSE) Architecture*. IETF Internet-Draft, draft-ietf-wimse-arch-07 (2 mars 2026). ⚠
- ISO/IEC JTC 1/SC 42 (2023). *Information technology — Artificial intelligence — Management system*. ISO/IEC 42001:2023.
- ISO/IEC JTC 1/SC 42 (2025). *Information technology — Artificial intelligence (AI) — AI system impact assessment*. ISO/IEC 42005:2025.
- MITRE Corporation (2021 ; base vivante). *MITRE ATLAS (Adversarial Threat Landscape for Artificial-Intelligence Systems)*. ⚠
- National Institute of Standards and Technology (2023). *Artificial Intelligence Risk Management Framework (AI RMF 1.0)*. NIST AI 100-1. DOI 10.6028/NIST.AI.100-1.
- National Institute of Standards and Technology (2024). *AI RMF: Generative Artificial Intelligence Profile*. NIST AI 600-1. DOI 10.6028/NIST.AI.600-1.
- OpenID Foundation, AIIM Community Group (South, T.; Cecchetti, S.; Fletcher, G.; et al.) (2025). *Identity Management for Agentic AI*. OpenID Foundation (whitepaper, 7 oct. 2025) ; arXiv:2510.25819.
- OpenID Foundation FAPI Working Group (2025). *FAPI 2.0 Security Profile* (Final Specification), 22 février 2025. https://openid.net/specs/fapi-security-profile-2_0-final.html.
- OpenTelemetry Authors (CNCF) (2025-2026). *Semantic Conventions for Generative AI*. Statut Development. ⚠
- OWASP API Security Project (2023). *OWASP API Security Top 10 — 2023*. OWASP Foundation, 3 juillet 2023.
- OWASP Gen AI Security Project (2025). *OWASP Top 10 for LLM Applications 2025*. OWASP Foundation.
- OWASP GenAI Security Project / Agentic Security Initiative (2025). *OWASP Top 10 for Agentic Applications (for 2026)* (publié 9 déc. 2025). ⚠
- OWASP Non-Human Identities Top 10 Project (2025). *OWASP Non-Human Identities Top 10 — 2025*. OWASP Foundation.
- SPIFFE Project (CNCF) (s.d.). *Secure Production Identity Framework for Everyone (SPIFFE) Specification / SPIRE* (projet diplômé). https://github.com/spiffe/spiffe. ⚠
- W3C / FIDO Alliance (Hodges, J.; Jones, J. C.; Jones, M. B.; et al., Eds.) (2021). *Web Authentication: An API for accessing Public Key Credentials — Level 2*. W3C Recommendation, 8 avril 2021.
- W3C Distributed Tracing Working Group (2021). *Trace Context (Level 1)*. W3C Recommendation, 23 novembre 2021.

### E. Domaine financier : standards de données et réglementation

- Aldasoro, I.; Gambacorta, L.; Korinek, A.; Shreeti, V.; Stein, M. (2024). *Intelligent financial system: how AI is transforming finance*. BIS Working Paper No 1194 (13 juin 2024) ; CEPR DP19181.
- Assemblée nationale du Québec (2021). *Loi modernisant des dispositions législatives en matière de protection des renseignements personnels (Loi 25) — art. 12.1 (décisions automatisées)*. LQ 2021, c. 25 ; art. 12.1 en vigueur 22.9.2023.
- Autorité des marchés financiers (Québec) (2026). *Ligne directrice sur l'utilisation de l'intelligence artificielle*. Version finale 7.4.2026 ; en vigueur 1.5.2027. ⚠
- Bank of England / FPC (2025). *Financial Stability in Focus: Artificial intelligence in the financial system*. Avril 2025.
- BIAN e.V. (2026). *Service Landscape 14.0* (architecture de référence bancaire, alignement ISO 20022). ⚠
- Civil Resolution Tribunal of British Columbia (2024). *Moffatt v. Air Canada*. 2024 BCCRT 149, 14 février 2024, dossier SC-2023-005609. https://www.canlii.org/en/bc/bccrt/doc/2024/2024bccrt149/2024bccrt149.html
- Comité de Bâle sur le contrôle bancaire (BCBS) (2017). *Basel III: Finalising Post-Crisis Reforms* (output floor à 72,5 %). BIS, d424, 7 déc. 2017.
- EDM Council / OMG (2026). *Financial Industry Business Ontology (FIBO) — dépôt de production (ontologies OWL)*. ⚠
- EIOPA (2025). *Opinion on Artificial Intelligence governance and risk management (EIOPA-BoS-25-360)*. 6.8.2025. ⚠
- Federal Reserve, OCC & FDIC (2026). *SR 26-2 / OCC Bulletin 2026-13 — Model Risk Management: Revised Guidance* (17.4.2026 ; supersède SR 11-7 et SR 21-8 ; IA générative et agentique hors périmètre). ⚠
- Financial Data Exchange (FDX) (2026). *FDX API — standard de partage de données financières (open finance/open banking)*. ⚠
- Financial Stability Board (2024). *The Financial Stability Implications of Artificial Intelligence*. FSB, 14 nov. 2024. ⚠
- FINOS, avec ISDA, ISLA, ICMA (2026). *Common Domain Model (CDM)*. FINOS Community Specification License 1.0. ⚠
- FINOS (Linux Foundation) (2025/2026). *AI Governance Framework v2.0 (AIGF)* (46 risques/mitigations ; catalogue agentique). Annonce 12.11.2025. ⚠
- FIX Trading Community (2026). *FIX — Family of Standards (FIX Latest)*. ⚠
- ISDA (2025). *FpML (Financial products Markup Language) — Recommendation version 5.13*. Publiée 12.5.2025. ⚠
- ISO/TC 68, ISO 20022 Registration Authority (2026). *ISO 20022-1:2026 Financial services — Universal financial industry message scheme — Part 1: Metamodel*. ⚠
- Office of the Superintendent of Financial Institutions (OSFI/BSIF) (2022). *Guideline B-13 — Technology and Cyber Risk Management*. Publiée 13.7.2022 ; en vigueur 1.1.2024. https://www.osfi-bsif.gc.ca/en/guidance/guidance-library/technology-cyber-risk-management.
- Office of the Superintendent of Financial Institutions (OSFI/BSIF) (2023). *Guideline B-10 — Third-Party Risk Management*. Version finale 30.4.2023 ; en vigueur 1.5.2024. https://www.osfi-bsif.gc.ca/en/guidance/guidance-library/third-party-risk-management-guideline.
- Office of the Superintendent of Financial Institutions (OSFI/BSIF) (2025). *Guideline E-23 — Model Risk Management (2027)*. Version finale 11 sept. 2025 ; en vigueur 1er mai 2027. ⚠
- Parlement européen et Conseil de l'UE (2014). *Directive 2014/65/UE concernant les marchés d'instruments financiers (MiFID II)*. JO L 173, 12 juin 2014.
- Parlement européen et Conseil de l'UE (2015). *Directive (UE) 2015/2366 (PSD2)*. JO L 337, 23 décembre 2015.
- Parlement européen et Conseil de l'UE (2016). *Règlement (UE) 2016/679 (RGPD), Article 22 — décision individuelle automatisée*. JO L 119, 4 mai 2016.
- Parlement européen et Conseil de l'UE (2022). *Règlement (UE) 2022/2554 sur la résilience opérationnelle numérique du secteur financier (DORA)*. JO L 333, 27 déc. 2022 ; applicable depuis le 17 janvier 2025.
- Parlement européen et Conseil de l'UE (2023). *Règlement (UE) 2023/1114 sur les marchés de crypto-actifs (MiCA)*. JO L 150, 9 juin 2023.
- Parlement européen et Conseil de l'UE (2023). *Règlement (UE) 2023/2854 (Data Act) — accès équitable aux données et leur utilisation*. JO L 2023/2854, 22 décembre 2023. http://data.europa.eu/eli/reg/2023/2854/oj.
- Parlement européen et Conseil de l'UE (2024). *Règlement (UE) 2024/1183 modifiant le règlement (UE) no 910/2014 — identité numérique européenne (eIDAS 2.0)*, 30 avril 2024.
- Parlement européen et Conseil de l'UE (2024). *Règlement (UE) 2024/1624 relatif à la prévention de l'utilisation du système financier aux fins du blanchiment de capitaux ou du financement du terrorisme (AMLR)*. JO L 2024/1624. ELI http://data.europa.eu/eli/reg/2024/1624/oj. ⚠
- Parlement européen et Conseil de l'UE (2024). *Règlement (UE) 2024/1689 établissant des règles harmonisées concernant l'intelligence artificielle (AI Act)*. OJ L 2024/1689, 12.7.2024.
- Parlement européen et Conseil de l'UE (2024). *Directive (UE) 2024/2853 (Product Liability Directive révisée)*. JO L, 18 nov. 2024 ; transposition d'ici le 9 déc. 2026 ; couvre les logiciels/produits IA.
- Payments Canada (2026). *Real-Time Rail (RTR) payment system — règles en vigueur 24.8.2026 ; lancement par vagues prévu à partir du T4 2026, élargi en 2027*. https://www.payments.ca/systems-services/payment-systems/real-time-rail-payment-system. ⚠
- U.S. Department of the Treasury (2024). *Artificial Intelligence in Financial Services*. Communiqué JY2760 (19.12.2024).
- United States Congress (2025). *GENIUS Act — Guiding and Establishing National Innovation for U.S. Stablecoins Act*. Public Law 119-27, 18 juil. 2025 (139 Stat. 419) ; mise en œuvre OCC Bulletin 2026-3 (25 fév. 2026). https://occ.treas.gov/news-issuances/bulletins/2026/bulletin-2026-3.html. ⚠

### F. Architecture d'entreprise et modélisation (ArchiMate, TOGAF)

- Azevedo, C. L. B.; Iacob, M.-E.; Almeida, J. P. A.; et al. (2015). *Modeling Resources and Capabilities in Enterprise Architecture: A Well-Founded Ontology-Based Proposal for ArchiMate*. Information Systems, 54, 235-262. DOI 10.1016/j.is.2015.04.008.
- Band, I.; Engelsman, W.; Feltus, C.; et al. — The Open Group (2019). *How to Model Enterprise Risk Management and Security with the ArchiMate Language* (Risk & Security Overlay). White Paper W172 ; 5 nov. 2019.
- Guizzardi, G.; Wagner, G.; Almeida, J. P. A.; Guizzardi, R. S. S. (2015). *Towards Ontological Foundations for Conceptual Modeling: The Unified Foundational Ontology (UFO) Story*. Applied Ontology, 10(3-4), 259-271. DOI 10.3233/AO-150157.
- Jonkers, H.; Lankhorst, M.; van Buuren, R.; et al. (2004). *Concepts for Modeling Enterprise Architectures*. International Journal of Cooperative Information Systems, 13(3), 257-287. DOI 10.1142/S0218843004000985.
- Lankhorst, M. et al. (2017). *Enterprise Architecture at Work: Modelling, Communication and Analysis* (4th ed.). Springer ; DOI 10.1007/978-3-662-53933-0.
- Van Riel, J.; Poels, G. (2023). *A Method for Developing Generic Capability Maps*. Business & Information Systems Engineering, 65(4), 403-424. DOI 10.1007/s12599-023-00793-z.
- ISO/IEC/IEEE (2022). *Software, systems and enterprise — Architecture description*. ISO/IEC/IEEE 42010:2022.
- The Open Group (2022). *ArchiMate 3.2 Specification*. Doc. C226 ; oct. 2022.
- The Open Group (2022). *The TOGAF Standard, 10th Edition — Introduction and Core Concepts*. Doc. C220 ; avr. 2022. ⚠
- The Open Group (2026). *ArchiMate 4 Specification*. Doc. C260 ; publié le 27 avr. 2026 (supersède C226). ⚠
- The Open Group & BIAN e.V. (2020). *ArchiMate Modeling Notation for the Financial Industry Reference Model: BIAN*. Open Group Guide G205 ; 31 mars 2020.

### G. Prospective, menace émergente, cryptographie post-quantique

- Anthropic — Threat Intelligence (2025). *Disrupting the First Reported AI-Orchestrated Cyber Espionage Campaign* (GTG-1002). 13 nov. 2025. ⚠
- Belcak, P.; Heinrich, G.; Diao, S.; et al. (2025). *Small Language Models are the Future of Agentic AI*. arXiv:2506.02153 (NVIDIA Research). ⚠
- Chan, A.; Ezell, C.; Kaufmann, M.; et al. (2024). *Visibility into AI Agents*. FAccT '24. DOI 10.1145/3630106.3658948 ; arXiv:2401.13138.
- Davidovic, S.; Tourpe, H. — Fonds monétaire international (2026). *How Agentic AI Will Reshape Payments*. IMF Notes 2026/004, 22 avr. 2026. ⚠
- European Systemic Risk Board — ASC (Cecchetti, S. G.; Lumsdaine, R. L.; Peltonen, T.; Sánchez Serrano, A.) (2025). *Artificial Intelligence and Systemic Risk*. ASC Reports No 16, déc. 2025. ⚠
- Europol — Quantum Safe Financial Forum (avec FS-ISAC) (2026). *Prioritising Post-Quantum Cryptography Migration Activities in Financial Services*. 21 janv. 2026. ⚠
- Gartner, Inc. (2025). *Gartner Predicts Over 40% of Agentic AI Projects Will Be Canceled by End of 2027*. Communiqué, 25 juin 2025. ⚠
- Hammond, L.; Chan, A.; et al. — Cooperative AI Foundation (2025). *Multi-Agent Risks from Advanced AI* (Technical Report #1). arXiv:2502.14143. ⚠
- International Data Corporation (IDC) (2025). *Agentic AI to Dominate IT Budget Expansion Over Next Five Years […] $1.3 Trillion in 2029*. Communiqué, 26 aout 2025, réf. prUS53765225. ⚠
- International Energy Agency (IEA) (2025). *Energy and AI* (rapport spécial), avr. 2025. ⚠
- Kolt, N. (2025). *Governing AI Agents*. Notre Dame Law Review, vol. 101 (à paraître). arXiv:2501.07913 ; SSRN 4772956. ⚠
- Kwa, T.; West, B.; Becker, J.; et al. — METR (2025). *Measuring AI Ability to Complete Long Software Tasks* (horizon de tâche, doublement ~7 mois). arXiv:2503.14499. ⚠
- McKinsey & Company (QuantumBlack) (2025). *Seizing the Agentic AI Advantage*. Juin 2025. ⚠
- MIT Project NANDA (Challapally, A.; et al.) (2025). *The GenAI Divide: State of AI in Business 2025*. MIT Media Lab. https://nanda.media.mit.edu/ ⚠
- National Institute of Standards and Technology (2024). *Module-Lattice-Based Key-Encapsulation Mechanism Standard (ML-KEM)*. NIST FIPS 203. DOI 10.6028/NIST.FIPS.203.
- National Institute of Standards and Technology (2024). *Module-Lattice-Based Digital Signature Standard (ML-DSA)*. NIST FIPS 204. DOI 10.6028/NIST.FIPS.204.
- National Institute of Standards and Technology (2024). *Stateless Hash-Based Digital Signature Standard (SLH-DSA)*. NIST FIPS 205. DOI 10.6028/NIST.FIPS.205.
- National Institute of Standards and Technology — Moody, D.; et al. (2024). *Transition to Post-Quantum Cryptography Standards* (dépréciation 2030, interdiction 2035). NIST IR 8547 (Initial Public Draft), 12 nov. 2024. ⚠
- NIST — Center for AI Standards and Innovation (CAISI) (2026). *Announcing the "AI Agent Standards Initiative" for Interoperable and Secure Innovation*. 17 février 2026. ⚠
- Schroeder de Witt, C.; et al. (2025). *Open Challenges in Multi-Agent Security: Towards Secure Systems of Interacting AI Agents*. arXiv:2505.02077. ⚠
- Soria Parra, D. — Model Context Protocol (2026). *The 2026 MCP Roadmap*. 9 mars 2026. ⚠
- World Economic Forum (2025). *The Future of Jobs Report 2025*. Janv. 2025. ⚠
- Zhang, Y.; Cai, Y.; Zuo, X.; et al. (2024). *The Fusion of Large Language Models and Formal Methods for Trustworthy AI Agents: A Roadmap*. arXiv:2412.06512. ⚠
