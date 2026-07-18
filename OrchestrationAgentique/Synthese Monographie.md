```{=typst}
#set page(numbering: "1")
#v(1.0cm)
#align(center)[
  #text(size: 18pt, weight: "bold")[L'autonomie encadrée : une revue de synthèse]

  #v(0.35cm)
  #text(size: 13pt, weight: "bold")[Interopérabilité et orchestration agentique dans les services financiers canadiens — protocoles ouverts, cadre réglementaire et blueprint d'intégration d'entreprise]

  #v(0.7cm)
  #text(size: 12pt, weight: "bold")[André-Guy Bruneau, M.Sc. IT]

  #v(0.15cm)
  #text(size: 10pt)[agbruneau\@gmail.com]

  #v(0.15cm)
  #text(size: 10pt, style: "italic")[Article de synthèse · juillet 2026]
]
#v(0.55cm)
#line(length: 100%, stroke: 0.4pt + luma(150))
#v(0.35cm)
```

**Résumé.** Cet article propose une revue de synthèse de l'*orchestration agentique* en contexte financier canadien, c'est-à-dire des conditions auxquelles des systèmes composés d'agents logiciels fondés sur de grands modèles de langage peuvent exécuter des processus soumis à une exigence réglementaire. Il consolide un corpus traité en sept parties : la couche protocolaire émergente (Model Context Protocol, Agent2Agent, Agent Payments Protocol, AGNTCY), l'orchestration multi-agents en entreprise, le cadre réglementaire canadien (ligne directrice E-23 du Bureau du surintendant des institutions financières, vide fédéral consécutif à la mort de la Loi sur l'intelligence artificielle et les données, ligne directrice sur l'intelligence artificielle de l'Autorité des marchés financiers, article 12.1 de la Loi 25, avis 11-348 des Autorités canadiennes en valeurs mobilières), l'interopérabilité financière canadienne (cadre des services bancaires axés sur le consommateur, ISO 20022 sur Lynx et le *Real-Time Rail*), l'adoption par dix institutions, une architecture de référence neutre et un blueprint instancié sur un portefeuille d'éditeur documenté. La thèse défendue est celle de l'**autonomie encadrée** (*framed autonomy*) : sous exigence réglementaire stricte, le processus est imposé de façon déterministe par le cadre, qui invoque les agents ; les agents n'orchestrent pas le processus. Cette thèse est indissociable d'une seconde proposition, méthodologique : ce qu'une institution peut architecturer est borné par ce qu'elle peut prouver, et le cadre est précisément le seul objet dont l'exploitant puisse démontrer la teneur devant un tiers. L'article dégage cinq apports : le raccordement de deux trajectoires jusqu'ici disjointes ; une taxonomie de l'encadrement transposée, avec ses conditions, au contexte canadien ; un résultat négatif établi — quinze croisements entre trois protocoles et cinq textes canadiens, aucun lien documenté — et cinq zones de compensation architecturale ; une architecture de référence à quatre couches et cinq points de contrôle obligatoires, séquencée sur l'échéance du 1er mai 2027 ; et une discipline probatoire opposable, qui distingue le fait éprouvé, le fait extrait, le repérage et l'inférence d'auteur. Le socle documentaire est arrêté au 17 juillet 2026 ; onze lacunes sont exposées plutôt que comblées.

**Mots-clés —** orchestration agentique ; autonomie encadrée ; Model Context Protocol (MCP) ; Agent2Agent (A2A) ; services financiers canadiens ; ligne directrice E-23 ; Loi 25 (article 12.1) ; ISO 20022 ; risque de modélisation ; gouvernance de l'intelligence artificielle.

## 1. Introduction

### 1.1 Contexte : deux trajectoires qui n'avaient aucune raison de se rencontrer

Entre la fin de 2024 et le milieu de 2026, deux histoires distinctes se sont déroulées en parallèle, chacune dans son vocabulaire, ses publications et ses institutions. Leur rencontre est le sujet de cette revue.

La première est **technique et protocolaire**. En novembre 2024, un éditeur publie une interface permettant à un modèle de langage d'appeler des outils et d'accéder à des données typées : le *Model Context Protocol* (MCP), fondé sur JSON-RPC 2.0 et assorti d'un cadre d'autorisation reposant sur OAuth. Dix-sept mois plus tard — de novembre 2024 au communiqué de la Linux Foundation du 9 avril 2026 —, la couche compte quatre projets ouverts, plusieurs fondations et des comités de pilotage technique pluriorganisationnels. Google lance le protocole *Agent2Agent* (A2A) en avril 2025 et le transfère à la Linux Foundation sous licence Apache 2.0 dès juin ; Cisco ouvre AGNTCY en mars 2025 avec LangChain et Galileo, puis le place sous la même fondation le 29 juillet 2025 ; la gouvernance de MCP passe à l'Agentic AI Foundation en décembre 2025. Aucun de ces protocoles n'est resté sous l'intendance exclusive de son créateur au-delà de treize mois.

Cette consolidation n'est toutefois pas uniforme, et la nuance décide de la portée du constat. Le protocole de la transaction — l'*Agent Payments Protocol* (AP2), annoncé par Google le 16 septembre 2025 comme protocole compagnon d'A2A, celui-là même qui intéresserait le plus directement une institution financière — ne fait l'objet d'**aucun transfert de gouvernance documenté**. Et deux versions stables seulement sont établies : A2A v1.0 et la révision de MCP du 25 novembre 2025. Le régime de la couche protocolaire n'est donc pas homogène : il est ouvert là où l'ouverture était la plus facile, et il reste indéterminé là où l'enjeu financier est le plus direct.

La seconde histoire est **sectorielle et canadienne**. Pendant les mêmes mois, les institutions financières du pays ont mis en service des systèmes qui décident, ou qui s'en approchent : pré-adjudication du prêt garanti par l'immobilier chez TD, traitement autonome de courriels commerciaux chez la Banque Scotia, plateformes agentiques d'entreprise chez la Banque Royale du Canada et chez Manuvie. Simultanément, le cadre réglementaire s'est réorganisé autour d'une date unique. La ligne directrice E-23 du Bureau du surintendant des institutions financières (BSIF), publiée en version finale le 11 septembre 2025, et la ligne directrice sur l'intelligence artificielle de l'Autorité des marchés financiers (AMF), finale depuis le 30 mars 2026, entrent en vigueur **le même jour : le 1er mai 2027**. Entre-temps, la mort au feuilleton du projet de loi C-27 à la prorogation du 6 janvier 2025 a emporté la Loi sur l'intelligence artificielle et les données ; le cadre des services bancaires axés sur le consommateur a été légiféré par le projet de loi C-15, sanctionné le 26 mars 2026, puis confié à la Banque du Canada ; et la migration ISO 20022 s'est achevée sur Lynx le 22 novembre 2025, tandis que le *Real-Time Rail* (RTR) entrait en tests industriels.

Table: Tableau 1 — Les deux trajectoires, jalons vérifiés (2023-2026). P = trajectoire protocolaire ; R = trajectoire réglementaire et financière canadienne.

| Date | Trajectoire | Événement |
| --- | --- | --- |
| 22 sept. 2023 | R | Entrée en vigueur de l'article 12.1 de la Loi sur la protection des renseignements personnels dans le secteur privé (Québec) |
| nov. 2024 | P | Publication de MCP par Anthropic |
| 5 déc. 2024 | R | Avis 11-348 des Autorités canadiennes en valeurs mobilières |
| 6 janv. 2025 | R | Prorogation ; mort au feuilleton de C-27 (dont la Loi sur l'intelligence artificielle et les données) |
| 17 mars 2025 | P | Lancement de BeeAI et de l'ACP protocolaire (*Agent Communication Protocol*, IBM Research/BeeAI) |
| avr. → juin 2025 | P | Lancement d'A2A par Google, puis transfert à la Linux Foundation (Apache 2.0) |
| 29 juill. 2025 | P | AGNTCY placé sous la Linux Foundation |
| 29 août 2025 | P | Fusion de l'ACP protocolaire dans A2A ; développement actif arrêté |
| 11 sept. 2025 | R | Version finale de la ligne directrice E-23 du BSIF (en vigueur le 1er mai 2027) |
| 16 sept. 2025 | P | Annonce d'AP2 par Google |
| 22 nov. 2025 | R | Lynx : fin de la coexistence MT/MX, alignée sur l'échéance SWIFT CBPR+ |
| 25 nov. 2025 | P | Révision de la spécification MCP en vigueur à la date de gel |
| déc. 2025 | P | Gouvernance de MCP transférée à l'Agentic AI Foundation |
| 26 mars 2026 | R | Sanction royale de C-15 (cadre des services bancaires axés sur le consommateur) |
| 30 mars 2026 | R | Version finale de la ligne directrice sur l'IA de l'AMF (en vigueur le 1er mai 2027) |
| 9 avr. 2026 | P | Bilan de première année d'A2A par la Linux Foundation |
| 15 juin 2026 | R | Première lecture du projet de loi C-36 (*Protecting Privacy and Consumer Data Act*) |
| 27 juin 2026 | R | Prépublication du règlement du cadre bancaire (Gazette du Canada, partie I) |
| 1er juill. 2026 | R | Règlement administratif no 10 du RTR (Gazette du Canada, partie II ; en vigueur le 24 août 2026) |

Entre ces deux histoires, il n'existe pas de pont documenté. Il existe un vide — et ce vide est l'objet de la revue.

### 1.2 Problématique : un vide établi, non supposé

Le constat mérite d'être posé avec précision, parce qu'il fonde tout le reste. La ligne directrice E-23, dont le texte intégral a été extrait et balayé mécaniquement dans ses deux versions linguistiques, **n'emploie ni « agentique », ni « agent(s) », ni « orchestration »** — zéro occurrence dans chaque cas ; la racine « autonom\* » y apparaît huit fois. Symétriquement, les spécifications protocolaires décrivent des appels, des formats et des habilitations : elles ne mentionnent aucun texte canadien. Croiser les trois protocoles majeurs avec les cinq textes canadiens pertinents produit quinze croisements dont **aucun n'est documenté par une source primaire**.

Ce vide n'est pas une curiosité bibliographique. Il place le praticien dans une position inconfortable et parfaitement identifiable : il doit démontrer, devant une seconde ligne de défense ou devant un régulateur, la maîtrise d'un système dont ni les spécifications techniques ni les textes réglementaires ne disent comment il devrait être construit. Deux échappatoires se présentent, également fautives. La première consiste à supposer que les protocoles apportent la conformité — la conformité n'est pas une propriété qu'un format d'échange puisse porter. La seconde consiste à supposer que le silence des textes vaut permission — ce serait confondre l'absence de prescription avec l'absence d'obligation, alors que les obligations existent déjà, en particulier celle de l'article 12.1 de la Loi 25, opposable depuis le 22 septembre 2023.

Il manque donc une lecture qui relie la couche protocolaire, l'ingénierie de l'orchestration, le maillage réglementaire canadien, l'infrastructure de paiement du pays et la pratique effective des institutions — puis qui expose, plutôt qu'elle ne les comble, les endroits où le raccordement n'est pas démontrable. Combler ce manque est l'objet de la présente revue.

### 1.3 Question de recherche et double thèse

La question directrice se formule ainsi : *à quelles conditions un système multi-agents peut-il exécuter un processus soumis à une exigence réglementaire canadienne, de manière défendable devant un tiers ?*

La réponse défendue est double, et ses deux versants sont solidaires.

Le premier versant est un **principe d'architecture** : sous exigence réglementaire stricte, le processus est imposé de façon déterministe par le cadre, qui invoque les agents ; les agents n'orchestrent pas le processus. C'est l'**autonomie encadrée** (*framed autonomy*) — l'agent conserve sa latitude à l'intérieur d'une tâche, il ne la conserve pas sur l'enchaînement. Dans la taxonomie des options d'orchestration exposée à la section 4, ce sont les positions OO3 et OO4.

Le second versant est une **discipline probatoire** : une affirmation ne vaut, dans un dossier de conformité comme dans un ouvrage, que par l'épreuve qu'elle a subie. D'où une hiérarchie explicite entre le fait ayant survécu à une tentative de réfutation, le fait extrait d'une source primaire sans avoir été contesté, et le simple repérage à confirmer ; d'où également le marquage systématique de ce qui relève de l'inférence de l'auteur, et la distinction stricte entre l'absence de documentation et le fait négatif vérifié.

Ces deux versants ne sont pas juxtaposés : le second explique pourquoi le premier s'impose. Le raisonnement, développé à la section 6, tient en trois pas. L'article 12.1 fait peser sur l'entreprise qui rend une décision fondée exclusivement sur un traitement automatisé l'obligation d'en communiquer, sur demande, les raisons ainsi que les principaux facteurs et paramètres, **sans lui ouvrir la ressource de désigner un tiers** — ni le développeur, ni le fournisseur du modèle, ni l'émergence elle-même. Or les travaux disponibles écartent successivement les producteurs possibles de cette explication : la journalisation confiée aux agents n'est pas recommandée ; les relations causales entre les entrées et les sorties d'un modèle d'apprentissage sont souvent indéterminables ; et les protocoles établissent qui parle, non ce qui est dit ni si cela est fondé. Reste un seul objet dont l'exploitant puisse démontrer la teneur, parce qu'il l'a écrit, qu'il précède l'exécution et qu'il ne dépend de la coopération d'aucun agent : le cadre qu'il a posé. **L'encadrement n'est pas, dans cette lecture, une bonne pratique d'ingénierie ; c'est le seul objet que l'assujetti puisse produire pour répondre de ce dont il répond de toute façon.**

Cette lecture explique pourquoi l'encadrement est nécessaire. Elle n'établit pas qu'il soit suffisant : aucune source ne dit qu'un cadre correctement posé libère l'assujetti, ni que la démonstrabilité d'un *frame* vaille preuve de conformité.

### 1.4 Contributions

La revue dégage et organise cinq contributions.

D'abord, un **raccordement** : la mise en rapport, sous une même grille, de deux littératures et de deux calendriers qui s'ignoraient — la couche protocolaire agentique et le maillage réglementaire canadien des services financiers. L'apport n'est pas de recenser un protocole ou un texte de plus, mais de tenir les deux à la fois et de dire exactement où ils ne se touchent pas.

Ensuite, une **taxonomie de l'encadrement transposée** : les quatre options d'orchestration OO1–OO4, reprises d'un préprint de la mi-2026, appliquées au contexte canadien avec l'énoncé explicite des quatre conditions qui bornent cette transposition — inférence d'auteur, absence de définition de l'« exigence réglementaire stricte », coût de l'encadrement, et silence du principe sur le contenu du cadre.

Troisièmement, un **résultat négatif établi et exploité** : la matrice de quinze croisements dont aucun n'est documenté, et les cinq zones de compensation architecturale qu'elle fait apparaître — production de la trace, point d'arrêt humain, détectabilité de l'auto-modification, interface d'accès du cadre bancaire, contenu et état. Leur propriété commune est qu'aucune ne se referme en choisissant un meilleur protocole.

Quatrièmement, une **architecture de référence** neutre à quatre couches, assortie de cinq points de contrôle obligatoires et d'une feuille de route en trois mouvements séquencés sur l'entrée en vigueur commune du 1er mai 2027 — puis éprouvée par construction sur un portefeuille d'éditeur documenté par sources primaires, retenu comme cas d'instanciation et non comme recommandation.

Enfin, une **discipline probatoire opposable**, qui n'est pas un appareil décoratif mais un instrument dont la revue rapporte les prises : deux fautes de méthode commises par le projet lui-même — une élévation de niveau de preuve doublement erronée, et un adjectif d'indépendance réfuté par le socle qui le portait — détectées, corrigées et racontées. La méthode n'a pas empêché les fautes ; elle les a rendues détectables.

### 1.5 Organisation de l'article

La progression suit l'ordre d'un raisonnement, non celui d'un catalogue. La section 2 expose la méthode et le cadre d'analyse. La section 3 traite la couche protocolaire : ce qu'elle est, ce qu'elle gouverne et ce qu'elle ne couvre pas. La section 4 porte sur l'orchestration multi-agents en entreprise — taxonomie, paradigme de l'autonomie encadrée, cadriciels, identité et registres. La section 5 expose le maillage réglementaire canadien. La section 6, pivot de l'article, traduit les exigences en contraintes de cadre et énonce le principe directeur avec ses conditions. La section 7 traite l'interopérabilité financière canadienne ; la section 8, l'adoption effective par dix institutions. La section 9 présente la synthèse architecturale — matrice, architecture de référence, points de contrôle, instrumentation. La section 10 présente la validation appliquée sous forme de blueprint instancié. La section 11 discute le résultat central, les limites et les onze lacunes ouvertes ; la section 12 conclut.

## 2. Méthode et cadre d'analyse

### 2.1 Nature de la revue et publics visés

Cette revue consolide un corpus doctrinal, normatif et documentaire constitué et vérifié pour l'ouvrage dont elle est tirée — vingt-quatre chapitres, quatre annexes et un avant-propos, soit vingt-neuf pièces et environ quatre-vingt-douze mille mots. Ce n'est pas une revue systématique au sens des protocoles PRISMA : son ambition est l'intégration conceptuelle et la traçabilité, non le dénombrement exhaustif d'une littérature.

Elle s'adresse à deux familles de lecteurs. Les **praticiens** — architectes d'entreprise et directions technologiques d'institutions financières canadiennes, responsables du risque, de la conformité et de la gouvernance des modèles, directions stratégie — y chercheront des décisions d'architecture datées et le périmètre exact de ce qui est démontrable. Les **chercheurs et analystes** y trouveront un corpus cité, une taxonomie, et une liste de questions ouvertes assortie, pour chacune, de ce qui la trancherait.

Deux dispositifs parcourent le texte pour les servir sans cloisonner le propos. Les encadrés *Perspective recherche* isolent les apports conceptuels et les questions ouvertes ; les encadrés *Mise en œuvre* isolent les décisions datées et les considérations de déploiement. Un troisième dispositif, propre à cet ouvrage, signale les **états de la connaissance vérifiable** : les points où la recherche n'a pas été conduite, ou a été conduite et a échoué.

> **Perspective recherche.** L'unité d'analyse n'est ni le modèle de langage ni le protocole, mais le *processus* et la place qu'y occupe la main qui commande l'enchaînement. Ce déplacement est décisif : il rend la question d'architecture indépendante du choix de produit, et il rend le choix de produit second — un cadriciel livre des matériaux qui autorisent plusieurs positions d'orchestration, il n'en arbitre aucune.

### 2.2 Le socle factuel et les trois niveaux de preuve

Toute affirmation factuelle centrale de l'ouvrage est tracée vers une entrée d'un **socle factuel de quarante-six entrées**, numérotées de F-01 à F-48 — trois identifiants n'ont jamais été attribués, à la suite de fusions de doublons lors de la synthèse, et une entrée porte un suffixe. Le cardinal ne se dérive jamais de la borne des identifiants : les deux ont coïncidé un temps, et cette coïncidence a rendu une erreur de décompte indétectable.

Le socle a été constitué en trois passes d'un harnais de recherche multi-agents, chacune en cinq temps : décomposition du sujet en angles, recherche parallèle par angle, extraction d'affirmations falsifiables avec déduplication des adresses, vérification adversariale par trois juges indépendants — la réfutation étant acquise à deux voix sur trois —, puis synthèse. Le bilan est arithmétiquement précis et il est exposé plutôt que résumé : **trois cent vingt-quatre lancements d'agents, soixante-dix-sept sources récupérées, trois cent quatre-vingt-quatre affirmations extraites, soixante-quinze soumises au vote adversarial, soixante-neuf confirmées à l'unanimité, six réfutées.**

Le rapport de soixante-quinze sur trois cent quatre-vingt-quatre est le chiffre à retenir, et il faut l'énoncer sans le maquiller : le vote à trois juges est borné par un plafond budgétaire aux vingt-cinq affirmations les plus importantes de chaque passe. **Les trois cent neuf autres n'ont pas été rejetées ; elles n'ont pas été soumises.** L'asymétrie ne disqualifie rien — elle indique où porterait le meilleur rendement d'une passe supplémentaire.

Table: Tableau 2 — Les trois niveaux de preuve, dans l'ordre strict [A] > [B] > [C].

| Niveau | Ce que l'affirmation a subi | Portée admise |
| --- | --- | --- |
| **[A]** | Vote adversarial 3-0 : trois vérificateurs indépendants ont cherché à la réfuter et ont échoué | Porte un fait central sans réserve de niveau |
| **[B]** | Source primaire lue et extraite, avec citation ou consultation directe, **sans vote** | Fait bien tenu, non éprouvé ; porte un fait central en le disant |
| **[C]** | Repérage documentaire, contenu non extrait | **Ne porte jamais un fait central** sans élévation préalable |

L'ordre est contre-intuitif et il commande tout le reste : **le niveau ne mesure pas la qualité de la source, il mesure ce que l'affirmation a subi.** Une entrée [A] repose elle aussi sur des sources primaires ; une entrée [C] en identifie une. Ce qui les sépare est l'épreuve, non le document.

Le projet a commis cette confusion et la rapporte. Le 16 juillet 2026, le texte intégral d'E-23 a été extrait et versé au socle ; l'opération a été consignée comme une « élévation [C] → [B] » de l'entrée concernée. Les deux termes étaient faux : l'entrée n'avait jamais été [C], et [B] se situe **sous** [A]. Une extraction primaire n'élève pas une entrée déjà votée à l'unanimité — elle l'enrichit d'un contenu de rang inférieur. L'entrée porte depuis un marquage mixte, et la correction s'est répercutée sur dix chapitres. La même méprise s'est reformée dans quatre chapitres distincts, parce que [B] se lit spontanément comme « mieux vérifié » là où il signifie « moins éprouvé ».

> **Perspective recherche.** Cette hiérarchie a une conséquence directe pour la pratique institutionnelle, indépendante de l'ouvrage : dans un dossier de diligence raisonnable, la question utile n'est pas « la source est-elle sérieuse ? » mais « qu'a-t-on tenté contre cette affirmation ? ». Classer ses propres affirmations en fait vérifié, inférence et repérage non confirmé — et n'appuyer une décision que sur la première catégorie — est un geste transposable tel quel.

### 2.3 Les garde-fous : ce qu'on s'interdit d'écrire

Huit affirmations réfutées ou non vérifiées sont proscrites dans tout texte dérivé de l'ouvrage. Elles ne sont pas des recommandations de style : chacune corrige une erreur que le corpus rend plausible, et plusieurs ont été effectivement commises avant d'être détectées.

Table: Tableau 3 — Les huit garde-fous et leur formulation de remplacement.

| # | Affirmation proscrite | Ce qu'il faut écrire |
| --- | --- | --- |
| R-1 | « L'ACP protocolaire est un standard vivant sous gouvernance de fondation » | Il a **fusionné dans A2A** le 29 août 2025 ; ses propriétés s'écrivent au passé |
| R-2 | « Entra Agent ID inclut un registre d'agents centralisé » | Identités d'agents et *blueprints* ; jamais « registre centralisé » |
| R-3 | « La spécification de la Cloud Security Alliance impose une identité SPIFFE » | Elle **s'appuie sur** SPIFFE/SPIRE comme fondation |
| R-4 | « Paiements Canada n'annonce aucune date de lancement pour le RTR » | La **cible officielle est le quatrième trimestre de 2026**, lancement graduel |
| R-5 | « FDX est le standard technique retenu du cadre bancaire canadien » | **Aucun standard n'est nommé** ; organisme de normalisation à désigner par arrêté ministériel |
| R-6 | Toute position au *Magic Quadrant* iPaaS d'un analyste | Ne pas citer de position ; deux quadrants distincts ne se substituent pas |
| R-7 | « Ce composant est conforme à E-23 » ou « satisfait B-13 » | **Aucun lien documenté** ; toute correspondance est une inférence d'auteur |
| R-8 | Employer le sigle « ACP » seul, ou confondre les quatre emplois de « (agentic) control plane » | Quatre objets distincts ; chaque emploi porte son qualificatif complet |

À ces garde-fous s'ajoutent des règles d'attribution. Toutes les métriques d'adoption des protocoles — cent cinquante organisations et plus déclarant soutenir A2A, plus de soixante pour AP2, plus de soixante-cinq pour AGNTCY — sont **auto-déclarées** et attribuées à leur communiqué à chaque occurrence, sous la formule qui gouverne leur lecture : *soutien n'est pas production*. Toutes les métriques d'institutions financières et d'éditeurs de plateformes sont des déclarations d'entreprise non auditées indépendamment, et portent leur attribution à chaque mention. La couverture de l'IA agentique par E-23 est une **inférence d'analystes juridiques**, jamais une exigence du régulateur. Les projections sont présentées comme projections de leur auteur. Les études d'analystes commandées mentionnent leur commanditaire.

Une règle transversale mérite d'être isolée, parce qu'elle a résisté à une tentative d'exemption : la règle d'attribution des métriques institutionnelles ne connaît aucun usage illustratif. Un premier jet citait deux chiffres d'institutions sans les nommer, au motif qu'ils servaient d'exemple ; la relecture a refusé l'exemption. **Un chiffre auto-déclaré qu'on cesse d'attribuer devient, en trois citations, un fait.**

### 2.4 Trois degrés d'absence

L'ouvrage distingue, et cette revue reprend, trois énoncés que la prose courante confond.

Le premier est le **fait négatif vérifié** : une absence établie par un balayage documenté d'un texte identifié, dont on peut énoncer le périmètre exact. Trois cas seulement dans l'ouvrage. E-23 n'emploie aucun terme du vocabulaire agentique, chaînes cherchées et décomptes donnés. Les textes officiels du cadre bancaire ne nomment aucun standard technique : recherche exhaustive des chaînes « FDX », « Financial Data Exchange », « FAPI » et « OAuth » dans le règlement prépublié, le communiqué ministériel et la page budgétaire, zéro occurrence — constat voté à neuf voix contre zéro. Et l'éditeur retenu au blueprint ne publie aucune architecture agentique propre aux services financiers, vérifié introuvable sur son site d'architectures et dans sa collection technique.

Le deuxième est le **fait négatif établi** : une absence portée par les réserves d'une entrée du socle sans balayage de texte — typiquement, l'absence de revendication de conformité d'un éditeur à un texte réglementaire canadien.

Le troisième, de loin le plus fréquent, est l'**absence de documentation dans le corpus de cet ouvrage**. Elle ne dit rien du monde : elle dit ce qui n'a pas été lu. Formule employée sans variation : *le socle ne dit pas qu'AP2 est resté propriétaire ; il dit qu'il ne documente pas le contraire.* Cette distinction n'est pas un scrupule d'auteur — elle décide de ce qu'une institution peut inscrire dans un dossier de diligence raisonnable.

Corollaire opératoire, tenu tout au long de l'ouvrage : **l'argument du silence est interdit hors des trois faits négatifs vérifiés.** On ne conclut rien de ce qu'un texte ne dit pas, sauf à avoir cherché et à pouvoir dire ce qu'on a cherché.

### 2.5 Datation, péremption et revalidation

Chaque pièce de l'ouvrage porte une **date de gel de l'information** — la date à laquelle son contenu a été arrêté. Pour les vingt-quatre chapitres, c'est le 16 ou le 17 juillet 2026. La conséquence est à prendre au sérieux : **l'ouvrage ne se périme pas d'un bloc, il se périme par morceaux, à des rythmes différents.** Un chapitre sur l'article 12.1 vieillira lentement ; un chapitre sur les versions de protocoles vieillira vite.

L'exemple est immédiat et il est assumé dans le texte : une révision majeure de la spécification MCP — cœur sans état, retrait d'un en-tête de session, nouveaux en-têtes de méthode et de nom — dont la version candidate a été verrouillée le 21 mai 2026 est en finalisation le **28 juillet 2026**, soit douze jours après le gel des chapitres qui la décrivent. Ces chapitres décrivent donc en connaissance de cause un état déjà daté, et le disent.

Le protocole de revalidation qui en découle tient en quatre gestes : constater à la source primaire, consigner dans un rapport daté, **amender le socle d'abord** — jamais les chapitres seuls —, puis reprendre les pièces dépendantes avec une nouvelle date de gel. Sa portée est une liste limitative de faits sensibles au temps ; un fait résolu en sort, parce qu'*une liste de veille qui ne se réduit jamais est une liste que l'on cesse de lire*. Le seuil retenu est arbitraire et déclaré tel : moins de trente jours au moment de la publication, calibré sur la cadence trimestrielle du domaine. La revalidation, précise l'ouvrage, a été menée sans vote adversarial : c'est un filet, pas une garantie.

> **Mise en œuvre.** Trois gestes transposent ce dispositif dans une institution, sans rien emprunter au contenu de l'ouvrage. Reconstituer sa propre liste de veille, à sa propre date. Classer chaque affirmation de son dossier d'architecture en fait vérifié, inférence ou repérage non confirmé — et n'appuyer une décision que sur la première catégorie. Rouvrir auprès de son canal éditeur ce que les outils publics n'ont pas pu établir : ce qu'une recherche documentaire n'obtient pas, une institution cliente l'obtient souvent par contrat.

### 2.6 Positionnement : ce que cette revue n'est pas

Trois avertissements bornent la portée du propos, et ils sont constitutifs plutôt que protocolaires.

**Aucun avis juridique.** L'ouvrage rapporte des textes, cite leurs formules officielles et en propose des lectures d'architecture. Ces lectures engagent l'auteur, jamais le régulateur. Qu'un processus donné relève ou non de l'article 12.1, qu'un assemblage d'agents constitue ou non « un modèle » au sens d'E-23 : ces questions sont posées, elles ne sont pas tranchées.

**Aucun conseil d'investissement.** Les institutions nommées à la section 8 le sont comme cas documentés par sources primaires, jamais comme recommandations, et l'inégalité de leur documentation publique est un trait de la couverture médiatique, non un classement de leur maturité.

**Aucune recommandation de fournisseur.** Le portefeuille instancié à la section 10 est retenu parce que sa documentation publique permettait de tracer chaque composant à une source datée. Ce n'est pas un verdict comparatif. Les faits défavorables à cet éditeur — une dépréciation de produits, une position d'analyste non vérifiable, l'absence de tout lien documenté avec la réglementation canadienne — sont exposés au même titre que ses capacités.

## 3. La couche protocolaire : ce qu'elle gouverne, ce qu'elle laisse ouvert

### 3.1 Généalogie et gouvernance : dix-sept mois

La première strate établit un fait de gouvernance et en borne aussitôt la portée. En dix-sept mois — de la publication de MCP en novembre 2024 au bilan de première année d'A2A publié par la Linux Foundation le 9 avril 2026 —, la couche protocolaire agentique est sortie du régime propriétaire. AGNTCY a mis quatre mois entre son ouverture et sa mise sous fondation ; A2A, deux mois ; MCP, treize — le plus lent des trois. Aucun de ces protocoles n'est resté sous l'intendance exclusive de son créateur au-delà de treize mois.

Le comité de pilotage technique d'A2A réunit **huit organisations** : Amazon Web Services, Cisco, Google, IBM Research, Microsoft, Salesforce, SAP et ServiceNow. Les membres formateurs d'AGNTCY sont Cisco, Dell Technologies, Google Cloud, Oracle et Red Hat. Ce sont là des faits de **structure** ; le corpus documente la composition de ces instances, non leurs règles de vote — la différence compte pour qui voudrait apprécier ce qu'un siège emporte réellement.

Trois restrictions accompagnent le constat, et sans elles il serait trompeur.

D'abord, **la gouvernance neutre est une condition nécessaire, non suffisante** de la crédibilité en entreprise réglementée. Une fondation ne durcit pas une implémentation.

Ensuite, **l'acquis ne s'étend pas à AP2**. Le protocole de la transaction agentique ne fait l'objet d'aucun transfert de gouvernance documenté — et il faut énoncer ce point dans la forme exacte : le corpus ne dit pas qu'AP2 est resté propriétaire, il dit qu'il ne documente pas le contraire.

Enfin, **la gouvernance neutre ne protège pas de l'obsolescence**. L'exemple est net et il est daté. L'ACP protocolaire (*Agent Communication Protocol*, IBM Research/BeeAI) a été lancé le 17 mars 2025 en pré-alpha, d'abord comme extension de MCP, et donné à la Linux Foundation dès le mois de mars. Un billet du 28 mai 2025 le décrivait comme RESTful sur HTTP et lui prêtait l'ambition d'être « the HTTP of agent communication ». Le 29 août 2025, sa fusion dans A2A était officielle : développement actif arrêté, actifs versés, adaptateurs et guides de migration publiés, et son autrice principale entrée au comité de pilotage technique d'A2A. Cinq mois et douze jours entre le lancement et l'arrêt ; trois mois et un jour entre l'énoncé de l'ambition et sa fin. L'obsolescence a été **ordonnée** plutôt que subie — les utilisateurs ont eu des adaptateurs et un chemin de migration — mais rien n'établit que la gouvernance en soit la cause. La leçon d'architecture est ailleurs, et elle se formule en une distinction : **le protocole peut mourir, l'utilisateur peut être abandonné** ; ce sont deux événements différents, et le second seul est un risque de projet.

Sur les métriques d'adoption, la lecture critique se déploie en trois réserves de gravité croissante. L'unité comptée — l'« organisation de soutien » — n'est définie nulle part. Un soutien déclaré n'est pas un déploiement en production. Et le triplement rapporté par la Linux Foundation, de cinquante-et-quelques organisations au lancement à cent cinquante et plus en avril 2026, est le triplement d'une grandeur non définie. Il reste que ces chiffres établissent quelque chose de précis et de limité : un protocole soutenu par cent cinquante organisations a peu de chances d'être abandonné dans les dix-huit mois. C'est exactement ce qu'ils établissent, et rien de plus. L'arbitrage d'architecture, lui, regarde **qui décide, non combien applaudissent**.

### 3.2 MCP : le contrat de la frontière agent-outil

MCP se laisse décrire par quatre éléments d'anatomie : une architecture client-serveur ; le transport JSON-RPC 2.0 ; l'échange de données typées ; et un **cadre d'autorisation** fondé sur OAuth.

Le quatrième terme appelle une discipline lexicale que l'ouvrage tient sans exception : on écrit « cadre d'autorisation », **jamais « sécurisé »**. Le protocole fournit un cadre ; la sécurité dépend de l'implémentation. Parler de « transport sécurisé par MCP » dans un dossier interne est une erreur de catégorie, à faire corriger plutôt qu'à nuancer — parce que la phrase prête au protocole une propriété qu'aucune spécification ne peut porter.

La révision en vigueur à la date de gel est celle du **25 novembre 2025**. La version candidate suivante a été verrouillée le 21 mai 2026 pour une finalisation au 28 juillet 2026 ; son contenu annoncé — cœur sans état, retrait de l'en-tête de session, introduction d'en-têtes de méthode et de nom — est décrit **par les mainteneurs du projet** comme la plus importante révision depuis le lancement. La qualification leur est attribuée ; elle n'est pas reprise à son compte.

Un point de méthode mérite d'être relevé, parce qu'il illustre la discipline des trois degrés d'absence. Le corpus fournit un inventaire des intégrations infonuagiques **pour A2A seulement** — préversion chez l'un, disponibilité générale chez l'autre depuis avril 2026, exécution dédiée chez un troisième — et rien d'équivalent pour MCP. L'ouvrage écrit alors : l'absence d'un inventaire infonuagique de MCP dans le corpus **n'établit pas** l'absence d'un tel support. La distinction entre « ce qui n'existe pas » et « ce qui n'a pas été vérifié » est de celles qu'un texte technique franchit sans y penser ; l'ouvrage s'astreint à ne jamais la franchir.

### 3.3 A2A : la frontière agent-agent, et ce qu'une signature ne dit pas

A2A v1.0 est la première spécification stable du projet, que **le projet lui-même** qualifie de prête pour la production. Trois propriétés la caractérisent : le support multi-protocole (JSON sur HTTP, gRPC, JSON-RPC) ; la multi-location d'entreprise (*enterprise multi-tenancy*) ; et les **cartes d'agents signées** (*Signed Agent Cards*), portant vérification cryptographique d'identité. Une ligne de correctifs mineurs est datée du 28 mai 2026. Le corpus ne date pas la publication de la v1.0 elle-même ; l'ouvrage s'abstient donc de la dater.

La doctrine de complémentarité entre les deux protocoles — « MCP dans les agents, A2A entre les agents » — est citée verbatim de la page de lancement d'A2A : *Complementary to MCP, not a replacement*. Deux précisions en fixent la portée, et elles ont toutes deux été acquises par correction. C'est une doctrine **déclarée par le projet A2A**, non un accord entre les deux projets. Et elle fournit un **critère de découpage**, non une contrainte : rien n'empêche techniquement de faire transiter une délégation par des appels d'outils. **Un critère n'est pas une contrainte** — la formule vaut avertissement contre la tentation de lire une répartition fonctionnelle comme une règle d'architecture.

Sur les cartes signées, l'ouvrage énonce ce que le corpus ne documente pas, et la liste est décisive : ni l'**ancrage de confiance**, ni les mécanismes de **révocation**, ni la **gouvernance des clés**. Ces trois questions décident, en pratique, de la valeur d'une signature. Il en va de même de la multi-location d'entreprise, nommée sans que mécanisme, périmètre ou garanties d'isolation soient documentés : une propriété **annoncée** n'est pas une propriété **qualifiée**.

> **Mise en œuvre.** Trois découpages concurrents de l'espace de conception coexistent dans la littérature de la mi-2026 : la doctrine A2A raisonne par **axes** (vertical agent-outil, horizontal agent-agent) ; un panorama académique de mai 2025 raisonnait par **étapes** d'adoption séquentielle — lecture périmée moins de quatre mois après sa publication par la fusion qu'elle n'avait pas prévue ; AGNTCY raisonne par **couches**. Aucun n'a le statut d'une norme. L'architecte prudent lira le statut avant de lire la marque : préversion, disponibilité générale ou déprécié disent davantage, sur ce qu'un composant emporte comme engagement, que le nom de l'éditeur qui le publie.

### 3.4 AP2 et AGNTCY : la transaction et l'infrastructure

AP2 a été annoncé par Google le 16 septembre 2025 comme protocole compagnon d'A2A, dédié aux transactions pilotées par agents. En avril 2026, la Linux Foundation rapporte que **plus de soixante organisations des paiements et des services financiers** en déclarent le soutien ; sept sont nommées — Mastercard, PayPal, American Express, Adyen, Coinbase, Worldpay, Revolut —, la liste étant illustrative et non exhaustive. Aucun décompte n'existant au lancement, contrairement à A2A, **aucune dynamique de croissance ne peut en être tirée**.

Ces sept noms sont pourtant qualitativement plus informatifs qu'un décompte anonyme — non parce qu'ils prouveraient un déploiement, mais parce qu'ils **identifient qui aurait à bouger** pour qu'un tel déploiement existe. Les cinquante et quelques autres soutiens restent non identifiables.

Ce que le corpus ne porte pas d'AP2 est plus lourd que ce qu'il en porte : **aucune anatomie technique** — ni structure de message, ni mécanique de mandat, ni modèle de preuve d'intention —, aucun transfert de gouvernance, et aucune articulation avec les rails de paiement canadiens. Ce dernier point est traité à la section 7.3.

AGNTCY, ouvert par Cisco en mars 2025 avec LangChain et Galileo, placé sous fondation le 29 juillet 2025 avec plus de soixante-cinq entreprises déclarant leur soutien, apporte des annuaires de découverte et un transport dédié. Son positionnement — couche d'infrastructure et non concurrent des protocoles précédents — est un **positionnement déclaré par le projet**, non un fait vérifié ; des analyses tierces relèvent d'ailleurs un chevauchement avec A2A que le corpus n'endosse pas, et qui **nuance** ce positionnement sans le contester ni l'établir.

### 3.5 Un avertissement terminologique : quatre objets, un sigle

Un obstacle de vocabulaire traverse tout le domaine et mérite d'être isolé, car il produit des confusions coûteuses en dossier d'architecture. L'expression « (agentic) control plane » et le sigle « ACP » désignent **au moins quatre objets distincts**, dont trois relevaient déjà du corpus initial et dont le quatrième a été découvert en cours de rédaction.

Table: Tableau 4 — Les quatre emplois de « (agentic) control plane » et leurs formes imposées.

| Branche | Objet | Forme à employer |
| --- | --- | --- |
| (a) | Protocole d'IBM Research/BeeAI lancé le 17 mars 2025, **fusionné dans A2A le 29 août 2025**, développement arrêté | « l'ACP protocolaire (*Agent Communication Protocol*, IBM Research/BeeAI) » — toujours au passé |
| (b) | Programme d'un consortium canadien annoncé le 7 juillet 2026 (Lightworks, Banque Scotia, Sun Life, TELUS) | « l'*Agentic Control Plane* du consortium » — jamais abrégé sans qualificatif |
| (c) | Positionnement d'éditeur d'un produit d'orchestration depuis mai 2026 : ni nom de produit, ni protocole | « le positionnement d'IBM comme *agentic control plane* » — attribué à IBM à chaque occurrence |
| (d) | Composante homonyme au sein d'AGNTCY, dont l'intitulé complet et l'identité avec (a) ne sont pas établis | « la composante ACP d'AGNTCY » — jamais agrégée à (a), jamais distinguée comme un fait |

**Le sigle « ACP » employé seul est proscrit.** Les statuts respectifs de ces branches sont eux-mêmes différenciés, et cette différenciation est un exercice de discipline probatoire en miniature. Sur le couple (a)-(b), le corpus se prononce : pure homonymie, absence de lien établie. Sur (c) vis-à-vis de (a) et (b), il n'établit que la distinction des objets — écrire « sans aucun rapport » serait un fait négatif fabriqué, et précisément sur le seul couple où il ne va pas de soi, (a) et (c) étant tous deux des objets du même éditeur. Sur (d), l'identité avec (a) n'est ni établie ni exclue, et les deux hypothèses conduisent à des conclusions d'architecture **opposées**. L'ouvrage ne tranche pas : *un ouvrage qui trancherait ici gagnerait en fluidité ce qu'il perdrait en droit d'être cru.*

### 3.6 Risques protocolaires : trois surfaces nommées, une mécanique absente

Trois surfaces d'attaque sont nommées par les protocoles eux-mêmes ou par la littérature académique de référence : l'**empoisonnement d'outils** (*tool poisoning*), l'**injection d'invites** (*prompt injection*) et l'**empoisonnement de mémoire** (*memory poisoning*). Une typologie utile — et marquée comme lecture d'auteur — les sépare par ce qu'elles corrompent : l'outil corrompt la **capacité**, l'invite corrompt l'**instruction**, la mémoire corrompt l'**état**. Aucun contrôle ne les couvre ensemble, parce qu'aucune de ces trois choses ne circule par le même canal — c'est là une définition opératoire de ce que la littérature appelle une sécurité *holistique*.

Deux réponses protocolaires existent : le cadre d'autorisation fondé sur OAuth du côté de MCP, les cartes d'agents signées du côté d'A2A. Leur limite est commune et structurelle : elles établissent **qui** parle, ni ce qui est dit, ni si ce qui est dit est fondé. L'authentification est une condition nécessaire de la confiance, jamais une condition suffisante.

> **État de la connaissance vérifiable.** Le corpus **nomme** ces risques ; il ne verse aucune source primaire sur leur mécanique, aucun identifiant de vulnérabilité, aucun incident public daté, aucune date d'apparition de la documentation. Aucune passe de recherche n'a été conduite sur ce point. Un lecteur qui aurait besoin de la mécanique — pour construire un plan de tests d'intrusion, par exemple — doit la chercher hors de cet ouvrage. S'y ajoute une asymétrie qu'il faut lire correctement : **aucune attaque propre à A2A n'est documentée par le corpus**. Ce silence est une propriété du corpus, non une propriété du protocole, et en aucun cas un certificat de sûreté.

Une précaution de raisonnement s'ajoute, et elle vaut au-delà de ce point : les deux corpus qui nomment ces risques ne sont pas indépendants — la littérature académique cite les spécifications protocolaires. Le second commente le premier, et **une reprise n'est pas une corroboration croisée**.

De ces trois sections, l'ouvrage tire une proposition qu'il présente comme un rapprochement d'auteur et que la suite entière peut se lire comme vérifiant : **la sûreté d'un système agentique se décide, pour l'essentiel, au moment où l'on décide de son architecture, et non au moment où l'on choisit ses protocoles.**

## 4. L'orchestration multi-agents en entreprise

### 4.1 Les quatre options d'orchestration

La deuxième strate fournit l'instrument central de l'ouvrage : une taxonomie qui objective le choix d'architecture agentique comme une position sur un continuum d'encadrement. Elle provient d'un préprint de la Technische Universität München daté du 30 juin 2026 — et son statut épistémique doit être énoncé avant son contenu : **préprint version 1, non révisé par les pairs, dont les auteurs déclarent eux-mêmes des menaces à la validité** (expériences initiales, invites non comparées entre elles, facteurs confondants). Le cadre conceptuel est repris comme référence ; les résultats chiffrés le sont **à titre d'illustration seulement**.

Deux axes structurent la taxonomie : qui détient la connaissance du processus, et qui commande l'enchaînement.

Table: Tableau 5 — La taxonomie OO1–OO4 : qui commande, qui sait.

| Option | Commande l'enchaînement | Connaissance du processus par l'agent | Caractérisation |
| --- | --- | --- | --- |
| **OO1** | Les agents | Aucune | Orchestration agentique agnostique au processus : le processus n'existe ni en modèle, ni en moteur, ni en invite — seulement comme propriété émergente de la conversation, observable *a posteriori* dans les journaux |
| **OO2** | Les agents | Cadre transmis par l'invite ou le contexte | Orchestration agentique consciente d'un cadre ; le cadre est défini par son **canal de transmission**, la source ne se prononçant pas sur sa force exécutoire |
| **OO3** | Le processus | Aucune | Orchestration de processus invoquant des agents agnostiques : l'agent garde son autonomie *à l'intérieur* de la tâche, l'enchaînement n'est plus négociable |
| **OO4** | Le processus | Le processus est connu de l'agent | Orchestration explicite d'agents conscients du processus : les deux axes saturés |

Trois précisions gouvernent l'usage de cette grille. Les transitions entre options sont **fluides** : ce sont des positions sur un continuum, non des catégories étanches. Le partage correct est {OO1, OO2} — l'agent commande — contre {OO3, OO4} — le processus commande ; à l'intérieur de chaque paire, seule varie la connaissance. Et **aucun des deux axes n'est celui de la capacité** : monter dans la taxonomie n'est pas monter en puissance.

Une observation d'auteur mérite d'être isolée, car elle porte tout le reste : la seule discontinuité véritable se situe **entre OO2 et OO3**, là où la main qui commande change de côté. La source ne hiérarchise pas les transitions et ne désigne aucun seuil ; elle n'établit non plus aucun ordre de coût entre les options.

### 4.2 Cinq propriétés, sept critères, et ce que la source ne dit pas

Cinq propriétés servent à évaluer une orchestration : l'**autonomie** (latitude de décision), la **spécificité de tâche** (degré de taillage pour une tâche déterminée), la **réactivité** (capacité à répondre à un événement dans les délais), l'**assurance de correction** (garantie de conformité du résultat à l'attendu) et la **traçabilité** (capacité à reconstituer et à suivre l'exécution).

Table: Tableau 6 — Les cinq propriétés et leur instrumentation documentée.

| Propriété | Métriques quantitatives documentées |
| --- | --- |
| Autonomie | **Aucune** |
| Spécificité de tâche | Complexité cyclomatique ; métrique ABC |
| Réactivité | Taux de faux négatifs ; vitesse de réaction |
| Assurance de correction | Précision, rappel, F1 |
| Traçabilité | Correction du journal |

L'absence de métrique pour l'autonomie n'est pas un détail : c'est la seule des cinq propriétés qui ne se mesure pas, et c'est celle qui définit l'objet même du domaine. Le constat porte sur l'entrée du corpus, non nécessairement sur l'article intégral — trancher supposerait une relecture ciblée, qui n'a pas été conduite.

Deux avertissements encadrent l'emploi de ces propriétés. Premièrement, **le corpus ne porte aucune matrice croisant les quatre options et les cinq propriétés** : un architecte qui présenterait à son comité de risque une matrice complète en la créditant à cette source lui ferait dire ce qu'elle ne dit pas. Deuxièmement, la lecture qui range quatre de ces propriétés du côté de la démonstrabilité et isole l'autonomie du côté de la latitude est une **construction de l'ouvrage**, non de la source.

Sept critères qualitatifs guident la sélection : complexité du but, supervision humaine, contraintes, action humaine, espace de décision, effort initial, maintenance. Aucune fonction de score, aucune pondération, aucun arbre de décision ne les relie aux quatre options. Trois points saillants. Le mot « contraintes » est écrit **sans adjectif** : les contraintes réglementaires en sont une espèce, non le genre. Supervision humaine et action humaine sont deux critères **séparés**, sans définition différentielle — surveiller institue un point d'observation, accomplir une tâche institue un point d'arrêt. Enfin, effort initial et maintenance ne caractérisent pas le processus mais son **coût de possession** : leur présence interdit de traiter le choix d'orchestration comme une pure question de conformité fonctionnelle. **L'encadrement se paie**, et la source le dit.

Sur les résultats expérimentaux — un score F1 de 0,40 en orchestration non encadrée, de 0,97 en orchestration pleinement encadrée, de 1,00 en déterministe pur, sur un scénario d'éclairage prédictif —, la position de l'ouvrage est explicite : ces trois nombres sont des illustrations et rien d'autre ; ils n'entrent ni comme preuve, ni comme ordre de grandeur transposable, ni comme argument.

Deux enseignements de conception sont, eux, plus robustes que les mesures. Le premier : **la journalisation confiée aux agents n'est généralement pas recommandée** — la source déconseille un producteur de trace sans en désigner un autre, et ce vide sera comblé, à la section 9, par une décision d'architecture explicite. Le second : les contraintes temporelles exigent des *frames* ou des outils externes, **les temps de raisonnement des modèles de langage étant imprévisibles**.

Un verdict, enfin, est porté par la même source sur un cas régi : sur un processus de don de sang encadré par une directive européenne, l'orchestration non encadrée est jugée **inacceptable** lorsque des exigences strictes d'exécution et de documentation s'appliquent, et les tâches essentielles doivent être imposées de façon déterministe par le cadre. Ce verdict est européen, hors finance, et issu d'un préprint : sa transposition au Canada est traitée — et bornée — à la section 6.

> **État de la connaissance vérifiable.** La taxonomie OO1–OO4 structure trois des sept parties de l'ouvrage et repose sur **une source unique**, versée au corpus seize jours avant la date de gel. Aucune reproduction indépendante, aucune application documentée à un processus d'institution financière, aucune validation de ses métriques comme indicateurs de risque financier. S'y ajoute une non-indépendance qu'il faut nommer : l'autrice principale de ce préprint cosigne également le manifeste académique mobilisé à la section suivante.

### 4.3 Le paradigme de l'autonomie encadrée

Le second appui académique est un manifeste de recherche sur l'*Agentic Business Process Management*, signé par dix-huit auteurs d'universités et d'industrie, issu d'un séminaire spécialisé, publié en version de revue en 2026. Son statut est celui d'une position argumentée : la confiance du corpus porte sur l'**attribution** — ce qui est certain, c'est que ces auteurs soutiennent cette thèse. Un manifeste ne fait pas autorité ; il fournit un vocabulaire. **Il recommande, il ne juge pas.**

Ce vocabulaire est précieux, et il commence par une distinction. **Automatiser** fixe le comportement : l'ingénieur décide de la suite des opérations. **Rendre autonome** délègue la décision : l'ingénieur fixe ce qui **borne** le choix. C'est une différence de nature, non de degré, et elle emporte une conséquence documentaire décisive — un contrôle automatisé se prouve par sa spécification ; une décision autonome ne se documente que par **ce qui la bornait**.

L'**autonomie encadrée** (*framed autonomy*) est érigée en mécanisme premier de gouvernance des systèmes agentiques. Elle repose sur deux types de cadres. Le **frame normatif** est de nature déontique : obligations, permissions, interdictions. Le **frame opérationnel** en est déclaré distinct — et, point qu'il faut signaler parce qu'il touche au cœur de la thèse de l'ouvrage, **il n'est pas caractérisé par le corpus**. Une caractérisation avait été fabriquée dans un premier jet ; la relecture adversariale l'a retirée. L'articulation des deux types de cadres est la thèse de l'ouvrage, et l'un des deux termes n'est pas défini par la source qui les distingue.

Quatre capacités sont requises d'un système ainsi gouverné : l'**encadrement** ; l'**explicabilité**, que les auteurs relient explicitement à la conformité réglementaire en nommant les instruments européens et en rangeant la finance parmi les domaines à haut risque — le corpus n'établit **aucun lien** entre ces développements et les textes canadiens ; l'**actionnabilité conversationnelle**, nommée et rangée parmi les quatre requises mais **sans définition, sans terme anglais et sans critère de satisfaction** ; et l'**auto-modification**, scindée en **adaptation** — éphémère, portant sur une instance d'exécution — et **évolution** — persistante, modifiant le modèle de processus pour toutes les instances à venir.

Cette dernière distinction produira, à la section 9, l'un des cinq points de contrôle obligatoires. La source ne dit pas que les deux relèvent de régimes d'autorisation distincts : c'est une inférence de l'ouvrage, et elle repose sur une observation simple — un système qui traite l'adaptation et l'évolution par le même chemin technique **rend indétectable, dans ses journaux, le moment où une exception est devenue une règle**.

Deux défis transversaux du manifeste irriguent la suite. Le premier est la **sécurité holistique**, qui ajoute aux trois surfaces d'attaque de la section 3.6 un **paradoxe de confidentialité de l'explicabilité** : l'explicabilité expose, la confidentialité restreint ; deux capacités requises tirent en sens contraire, et le manifeste nomme la tension sans la résoudre. Le second est l'**écart de responsabilité** (*responsibility gap*) : l'indétermination de l'imputabilité juridique entre quatre porteurs — le développeur, **l'organisation qui impose le cadre**, le fournisseur du modèle, et le comportement émergent du système multi-agents. Le quatrième n'est pas une personne : il est le nom de ce dont aucun des trois autres n'est l'auteur au sens ordinaire. La section 6 y revient.

Un dernier apport, souvent négligé, retourne la notion de cadre : l'**opérationnalisation locale des frames** constitue une **frontière de sécurité et de confidentialité**. Restreindre le contexte et les capacités de chaque agent limite l'impact d'un agent compromis. Le cadre devient alors instrument de **confinement** autant que de gouvernance — logique de limitation des dégâts, non de prévention. Le corpus n'établit pas que l'encadrement soit une réponse *suffisante* aux risques de la section 3.6.

### 4.4 Les cadriciels : ce qu'ils livrent, ce qu'ils n'arbitrent pas

L'offre s'est industrialisée en 2025-2026, et la question utile n'est pas de la classer mais de mesurer ce que l'on sait d'elle. Une convention de lecture gouverne toute cette section : « le corpus ne documente pas » signifie *absence de documentation dans le corpus de cet ouvrage*, jamais un fait négatif vérifié.

Table: Tableau 7 — Cadriciels d'orchestration : statut probatoire du support protocolaire (état arrêté au 16 juillet 2026).

| Cadriciel | Jalon | Support MCP | Support A2A |
| --- | --- | --- | --- |
| Microsoft Agent Framework | Disponibilité générale 1.0 le 3 avr. 2026 ; successeur assumé de deux lignées antérieures | **Première main** : clients et serveurs MCP, outils locaux et hébergés | Non documenté pour le cadriciel lui-même |
| LangGraph Platform | Disponibilité générale le 14 mai 2025 | **Non documenté** : le billet de disponibilité générale ne mentionne ni MCP ni A2A | **Première main, pour la plateforme commerciale seulement** ; absent de la bibliothèque libre (demande ouverte, avr. 2026) |
| Confluent / Kafka | *Streaming Agents* en août 2025 ; mise à jour du 26 févr. 2026 | **Première main** : appel d'outils MCP natif ; serveur MCP officiel | Intégration en **préversion publique** ; aucun client nommé, aucun chiffre d'adoption |
| CrewAI | Journal des modifications depuis nov. 2025 | Repérage à confirmer, contenu non extrait | **Première main** : primitive de délégation de premier rang, extension et classes documentées |
| Temporal | — | Repérage à confirmer | Non documenté |

Ce décompte fonde une formulation prudente qu'il faut préférer à celle, plus commode, qui circule : le support de MCP est **répandu et inégalement établi** — deux offres sur cinq documentées de première main, deux au niveau du repérage, une non documentée. Écrire « généralisé » serait une thèse que le corpus réfute.

Deux points méritent d'être isolés parce qu'ils décident d'architectures.

Le premier est une **frontière commerciale prise pour une frontière technique**. Le support d'A2A par une plateforme d'orchestration très adoptée est confirmé de première main — point d'entrée dédié, méthodes documentées — **pour l'offre commerciale seulement** ; le cœur libre du même produit ne l'a pas, et la fonctionnalité y demeure une demande ouverte. Une institution qui bâtirait sur la bibliothèque libre en supposant la parité fonctionnelle découvrirait la différence tard.

Le second est un **risque d'exploitation documenté**. Le cadriciel le mieux caractérisé du corpus porte des limites connues de son magasin de points de contrôle en déploiement distribué multi-instances. Le corpus n'en établit ni la gravité, ni le contournement, ni le calendrier de correction — mais l'un des cinq points de contrôle obligatoires de la section 9 en dépend directement.

Les chiffres d'adoption, enfin, sont tous auto-déclarés et portent leur attribution. Un éditeur déclare, dans son billet de disponibilité générale du 14 mai 2025, que près de quatre cents entreprises avaient déployé des agents en production via sa plateforme depuis une bêta ouverte onze mois plus tôt — sans définition publiée de l'unité comptée. Un autre déclare une base de clientèle parmi les grandes entreprises et environ deux milliards d'exécutions sur douze mois, sans source datée ni définition de l'unité. L'écart entre les deux métriques joue **contre** la seconde : la première est datée et rattachée à un billet identifiable, donc situable et contestable.

L'enseignement principal de cette strate est une reformulation de la question d'architecture, et il est le plus transposable de l'ouvrage : **aucun des cadriciels examinés ne livre un positionnement d'orchestration ; ils livrent des matériaux qui en autorisent plusieurs. Le produit n'arbitre pas ; la configuration arbitre.** La question n'est donc pas « quel cadriciel choisir ? », mais « **quel cadre imposer, et quel produit sait le porter ?** »

### 4.5 Identité et registres : une moitié ratifiée, une moitié en brouillon

La dernière strate de cette partie traite ce qui nomme l'agent et borne ses droits hors de lui.

Du côté des produits, une offre d'identité d'agents est en disponibilité générale depuis le printemps 2026 : création et gestion d'identités d'agents distinctes de celles d'un utilisateur, et gestion de *blueprints* — terme de l'éditeur, dont le corpus ne livre pas la définition. Les fondations sont OAuth 2.0 pour l'autorisation et OpenID Connect pour l'authentification. Deux familles de scénarios sont documentées, l'une où l'agent agit pour son propre compte et l'autre où il agit par délégation — distinction qui détermine si la piste d'audit désigne un mandant humain identifié ou l'agent seul. Deux réserves accompagnent le fait : plusieurs sous-fonctionnalités demeurent en préversion et relèvent de licences additionnelles ; et les flux de délégation et de jeton d'agent **étendent** les normes de référence — le dispositif ne doit jamais être présenté comme du pur standard. Il ne doit pas davantage être présenté comme un registre centralisé : cette affirmation n'a pas été confirmée.

Du côté pré-normatif, une spécification de registre d'agents publiée le 27 mars 2026 par un laboratoire d'association sectorielle — **statut de brouillon**, à énoncer à chaque mention — assigne au registre d'entreprise faisant autorité quatre fonctions : **inventaire, découverte, cycle de vie, conformité**. Aucun de ces termes n'appartient au vocabulaire de l'infrastructure : ce sont des fonctions de gouvernance. Le profil d'agent s'ancre dans une extension d'une norme de provisionnement d'identités, portée par un brouillon individuel de l'organisme de normalisation de l'Internet. Deux champs obligatoires sont documentés : la **liste des outils et serveurs autorisés** pour l'agent, et le **bornage de ses privilèges** selon le principe du moindre privilège. La portée est remarquable : on ne peut enregistrer un agent sans déclarer, de façon explicite et lisible par machine, les outils qu'il peut atteindre et les bornes de ses droits. **La déclaration devient condition d'existence de l'agent dans l'annuaire.**

Cette structure est, en substance, un **frame opérationnel écrit hors de l'agent** — rapprochement d'auteur, mais qui donne à la section 9 sa deuxième couche.

L'édifice est cependant explicitement double : sa moitié inférieure est ratifiée et stable ; sa moitié supérieure — ce qui décrit spécifiquement les agents — est un brouillon individuel **expiré le 19 avril 2026**, repris par un brouillon de laboratoire, sans calendrier de consolidation. Vingt-trois jours séparent la publication de la spécification de l'expiration du brouillon qui l'ancre. Conséquence pratique : emprunter à ce corpus sa **structure** — déclarer les outils atteignables et les bornes de privilège — sans traiter sa **forme** — les identifiants de champ — comme un contrat d'interopérabilité.

> **État de la connaissance vérifiable.** Le corpus ne documente **aucune norme ratifiée de registre d'agents, à quelque niveau que ce soit**. Un architecte qui chercherait aujourd'hui la norme d'identité et de registre des agents ne la trouverait pas : elle n'existe pas. Et le corpus ne relie ni l'offre d'identité, ni la spécification de registre, aux textes canadiens qui entrent en vigueur le 1er mai 2027 — absence de documentation, non fait négatif vérifié.

## 5. Le cadre réglementaire canadien

### 5.1 E-23 : une couverture par inférence, à traiter comme acquise

La ligne directrice E-23 du Bureau du surintendant des institutions financières, consacrée à la gestion du risque de modélisation, a été publiée en version finale le **11 septembre 2025** et entre en vigueur le **1er mai 2027** — un préavis de dix-neuf mois et vingt jours. Elle s'applique à toutes les institutions financières fédérales.

Sa définition de « modèle » est le pivot de tout le raisonnement. Elle vise l'application de techniques statistiques ou d'hypothèses théoriques, empiriques ou fondées sur un jugement, **notamment des méthodes d'intelligence artificielle et d'apprentissage automatique**, qui traite des données en entrée pour générer des résultats. La lettre d'accompagnement établit que le régulateur a laissé cette définition **délibérément large**, en réponse à des commentaires qui demandaient de la restreindre. Le texte anticipe par ailleurs explicitement des modèles à apprentissage et à décision autonomes.

Sur cette base, une formulation obligatoire encadre tout énoncé relatif à la portée agentique du texte, et elle doit être reprise sans variation : **E-23 ne nomme ni l'agentique ni les agents ; sa définition de « modèle » englobe les méthodes d'intelligence artificielle et d'apprentissage automatique, d'où une couverture implicite que les analystes juridiques canadiens tiennent pour acquise.** Trois éléments sont ici superposés et ne se valent pas. Le fait négatif est vérifié — balayage mécanique du texte intégral, dans les deux langues, décomptes donnés. La couverture implicite est une **conclusion**, non une observation. Et ce qui a été vérifié est **un état de l'opinion juridique canadienne, non un état du droit**. Écrire « le régulateur exige, pour l'IA agentique, que… » est faux, et ce n'est pas discutable.

L'ouvrage recommande néanmoins de traiter la couverture comme acquise — en documentant que c'est par prudence. Trois éléments convergents l'y conduisent, dont aucun ne suffirait isolément : la convergence des analystes ; la construction fonctionnelle de la définition ; et l'anticipation explicite de la décision autonome par le régulateur lui-même. S'y ajoute une asymétrie de risque, marquée comme lecture d'auteur : se tromper en tenant la couverture pour acquise coûte une gouvernance surdimensionnée ; se tromper en la tenant pour inexistante coûte un parc décisionnel entier hors périmètre au 1er mai 2027, découvert par un tiers.

Le texte est **fondé sur des principes** et rédigé au conditionnel — douze principes formulés en « devrait ». La discipline lexicale qui en découle est stricte : on écrit « **attendu par** E-23 », jamais « exigé par ». Cinq familles d'attentes opératoires en ont été extraites : un **inventaire** des modèles au risque inhérent non négligeable, tenu à l'échelle de l'entreprise et qualifié d'exact, perpétuellement à jour et soumis à des contrôles robustes ; une **cote de risque** attribuée à chaque modèle, l'intensité du dispositif devant être proportionnée au risque ; des normes de **documentation de modèle** — jamais « fiche de modèle », terme absent du texte dans les deux langues ; un **cycle de vie** en cinq volets, explicitement non nécessairement séquentiels ; et une **surveillance continue** dont le libellé nomme la prise de décision autonome, la reparamétrisation autonome et un potentiel accru de dérive.

Un rapport conjoint du superviseur prudentiel et de l'agence de protection du consommateur en matière financière, daté du 24 septembre 2024, apporte deux éléments. Le premier est une trajectoire d'adoption : d'environ 30 % en 2019 à environ 50 % en 2023, **avec une projection d'environ 70 % d'ici 2026 — projection issue d'une enquête auto-déclarée, non un taux observé**, formule à reprendre à chaque occurrence. L'arithmétique mérite d'être exposée : vingt points en quatre ans déclarés, vingt points en trois ans projetés — la projection accélère d'un tiers la tendance qu'elle prolonge. Le corpus ne documente aucune mesure postérieure : l'horizon est atteint sans vérification, et aucune extrapolation n'est proposée.

Le second élément est plus lourd pour l'architecture : le rapport établit que le risque de modélisation propre à l'intelligence artificielle tient au très grand nombre de paramètres appris de façon autonome et à des **relations causales entre entrées et sorties souvent indéterminables**. Si la causalité est indéterminable, l'explication d'une décision **ne peut pas être obtenue en interrogeant le modèle après coup** : elle doit être produite par ce qui l'entoure. Cette inférence est exploitée à la section 6.

### 5.2 Le vide fédéral : de C-27 à C-36

Le projet de loi C-27 réunissait un volet de protection de la vie privée des consommateurs et la Loi sur l'intelligence artificielle et les données ; les deux sont morts au feuilleton à la prorogation du **6 janvier 2025**. Le corpus établit la concomitance, non la causalité, et ne documente **aucune obligation** qu'aurait imposée le volet consacré à l'intelligence artificielle : l'ouvrage ne compare donc cette loi défunte à rien.

Un ministre de l'intelligence artificielle et de l'innovation numérique — premier titulaire du portefeuille — a été nommé en mai 2025. Le 15 juin 2026, treize mois plus tard, le projet de loi **C-36** a été déposé et a franchi sa première lecture ; il était à l'étape de la deuxième lecture à la date d'observation du 16 juillet 2026.

Sa qualification appelle une formulation obligatoire, parce que l'erreur inverse est spontanée. C-36 est une **réforme de la protection des renseignements personnels** : il modifie la loi fédérale existante et prévoit la création d'une commission de la sécurité numérique et de la protection des données, laquelle **n'existe pas**. Trois volets y figurent — droit à la suppression, protection des données des enfants, **transparence algorithmique**, ce dernier étant seul qualifié d'intelligence artificielle par le corpus —, et ils sont portés **comme composantes d'une loi sur la vie privée, non comme les dispositifs d'une loi sur l'intelligence artificielle autonome**. L'expression « la loi canadienne sur l'IA » appliquée à C-36 est proscrite. Le texte est un projet de loi et non une loi ; l'ouvrage ne se prononce ni sur son adoption ni sur son échec.

Le vide qui subsiste est donc **spécifique** : le Canada n'est pas un espace de non-droit, mais il demeure sans cadre contraignant propre à l'intelligence artificielle au niveau fédéral. La charge se déplace vers quatre instruments datés — la ligne directrice prudentielle, la ligne directrice québécoise, l'avis en valeurs mobilières, et l'article 12.1 de la Loi 25 —, dont le dernier relève du droit général québécois plutôt que d'un régulateur sectoriel.

### 5.3 Québec : le cadre le plus explicite, et une friction directe

La ligne directrice sur l'intelligence artificielle de l'Autorité des marchés financiers a suivi un parcours daté : projet publié le 3 juillet 2025, consultation close le 7 novembre 2025, **version finale publiée le 30 mars 2026**, entrée en vigueur le **1er mai 2027**. Une réserve impérative s'attache à ce fait : il ne faut **jamais écrire « en attente » ni « en projet »** — cet état est périmé depuis le 30 mars 2026. Le corpus porte le **calendrier de ce texte, non son contenu** ; on y reviendra, car c'est la lacune la plus coûteuse de l'ouvrage.

L'article 12.1 de la loi québécoise sur la protection des renseignements personnels dans le secteur privé est, lui, opposable depuis le **22 septembre 2023**. Son déclencheur est une **décision fondée exclusivement sur un traitement automatisé** de renseignements personnels — deux conditions cumulatives : l'objet doit porter sur des renseignements personnels, et l'exclusivité doit être caractérisée.

Trois obligations en découlent, et leur structure temporelle est l'élément décisif.

**Informer** la personne qu'une telle décision a été rendue, **au plus tard au moment de la décision**. C'est la seule des trois qui soit proactive et datée : elle ne dépend d'aucune demande.

**Expliquer, sur demande** : communiquer les renseignements personnels utilisés, ainsi que les raisons et les principaux facteurs et paramètres ayant mené à la décision, et informer du droit de rectification. L'obligation porte sur **l'instance**, non sur le modèle — elle demande pourquoi *cette* décision a été rendue, pas comment le système fonctionne en général.

**Permettre la révision humaine** : donner l'occasion de présenter ses observations à un membre du personnel de l'entreprise **en mesure de réviser la décision**. Donc une personne rattachée à l'entreprise, dotée du pouvoir de **défaire** — non seulement d'expliquer ou de transmettre.

Le critère d'exclusivité n'est pas défini par le texte. Selon l'analyse d'un cabinet, une intervention humaine significative *avant* la décision écarterait l'application de l'article — nuance déterminante pour les systèmes à humain-dans-la-boucle. Son statut doit être énoncé : **interprétation de cabinet, retenue au niveau du repérage à confirmer, non confrontée aux positions de l'autorité de contrôle québécoise**. Elle ne peut donc porter aucun fait central.

Une distinction de vocabulaire s'impose ici, et elle a des conséquences que la section 10 rendra concrètes : **l'humain-dans-la-boucle et la révision humaine de l'article 12.1 ne sont pas la même chose.** Le premier est un point d'arrêt inséré dans l'exécution ; la seconde est un recours ouvert *après* la décision, à l'initiative de la personne concernée. Deux branches en découlent : si le dispositif d'intervention est jugé significatif au sens de l'analyse ci-dessus, il écarte l'article **en entier** ; sinon, l'article s'applique **intégralement** et la révision demeure due. Dans les deux cas, **outiller l'humain-dans-la-boucle n'outille pas la révision humaine**.

Trois contraintes d'architecture en découlent, marquées comme lecture d'auteur. La **traçabilité de l'instance** : restituer l'enchaînement, non la trace d'un appel. Le **synchronisme** : le moment de la décision doit être un événement identifiable, faute de quoi l'obligation d'information est mécaniquement inexécutable. Et la **réversibilité** : la révision suppose qu'une décision rendue, communiquée et peut-être exécutée en aval puisse être défaite.

### 5.4 Valeurs mobilières : les lois existantes s'appliquent

L'avis 11-348 des Autorités canadiennes en valeurs mobilières, publié le 5 décembre 2024 et dont le volet de consultation s'est clos le 31 mars 2025, vise sept catégories d'assujettis — personnes inscrites, émetteurs assujettis autres que les fonds d'investissement, marchés et leurs participants, chambres de compensation et services d'appariement, référentiels centraux, agences de notation désignées, administrateurs d'indices de référence désignés. Trois de ces sept catégories sont des **infrastructures** qui ne décident de rien au sens de l'opération : un système agentique déployé dans une infrastructure de marché relève donc du même champ que chez une personne inscrite.

La doctrine est double et symétrique. Les lois existantes **s'appliquent** aux systèmes d'intelligence artificielle, et les indications données **ne créent ni ne modifient aucune exigence**. Conséquence : aucune obligation nouvelle, aucun délai de transition, aucun programme daté — **et** aucune obligation existante écartée, tempérée ou suspendue au motif qu'une intelligence artificielle en serait l'instrument. Il n'y a **ni échéance, ni sursis**.

L'accroche pour l'agentique tient à la définition retenue, qui inclut explicitement des **niveaux variables d'autonomie et d'adaptativité après déploiement**. Le premier terme institue un continuum plutôt qu'une catégorie ; le second capte un système dont le comportement se modifie après sa mise en service — c'est-à-dire exactement l'objet que la distinction adaptation/évolution de la section 4.3 cherchait à nommer.

Deux réserves closent ce point. Le corpus **n'établit pas si l'avis formule des attentes** de gouvernance, d'explicabilité ou de supervision au-delà du constat d'applicabilité. Et les suites de la consultation ne sont documentées par aucune source : le corpus documente une **absence de source, non une absence de production** du régulateur — rien n'est présumé, ni instrument à venir, ni abandon, ni approbation tacite.

### 5.5 Le maillage, vu d'ensemble

Table: Tableau 8 — Les instruments canadiens mobilisés, avec leur nature et leur échéance.

| Instrument | Nature | Statut au 17 juillet 2026 | Échéance |
| --- | --- | --- | --- |
| Ligne directrice E-23 (BSIF) | Ligne directrice **fondée sur des principes**, au conditionnel — elle *attend* | Finale depuis le 11 sept. 2025 | **1er mai 2027** |
| Ligne directrice sur l'IA (AMF) | Ligne directrice sectorielle québécoise | **Finale depuis le 30 mars 2026** ; contenu hors corpus | **1er mai 2027** |
| Article 12.1, Loi 25 (Québec) | Disposition législative à l'impératif — elle *oblige* | **Opposable depuis le 22 sept. 2023** | En vigueur |
| Avis 11-348 (ACVM) | Avis du personnel : ne crée ni ne modifie aucune exigence | Publié le 5 déc. 2024 ; consultation close | Ni échéance, ni sursis |
| Projet de loi C-36 | Réforme de la protection des renseignements personnels, portant des volets IA | Deuxième lecture ; **adoption incertaine** | Indéterminée |
| Cadre bancaire (C-15 et règlement) | Loi sanctionnée ; règlement prépublié | Sanction le 26 mars 2026 ; commentaires clos le 26 août 2026 | Entrée en vigueur échelonnée |

Une différence de nature traverse ce tableau et commande la suite : les lignes directrices **attendent** un programme, l'article 12.1 **oblige** à un résultat. Les contraintes dérivées des unes sont des conditions de possibilité d'un programme attendu ; celles dérivées de l'autre sont des conditions de possibilité d'une obligation.

Une dernière observation, et elle n'est pas anodine : **le 1er mai 2027 n'est pas une coïncidence documentée.** Le corpus établit deux dates identiques et **aucune concertation** entre les deux régulateurs. Une simultanéité se constate ; une concertation s'établirait.

## 6. Le pont : des contraintes réglementaires aux cadres déterministes

Cette section est le pivot de l'ouvrage. Elle est aussi celle qu'il désigne lui-même comme la première à contester.

### 6.1 La table de traduction : onze entrées, neuf contraintes

L'opération consiste à demander, texte par texte, ce que l'architecture doit porter pour que ce que le texte attend ou oblige soit possible. Une précaution de méthode la borne : une exigence ne devient contrainte d'architecture que si **trois choses sont connues — ce que le texte impose, à qui, et sur quel objet**. L'une des trois manquant, aucune contrainte n'est dérivée.

Table: Tableau 9 — Traduction des exigences canadiennes en contraintes de cadre (onze entrées ; toute contrainte dérivée est une inférence d'auteur).

| # | Texte et attente | Ce que le cadre doit porter |
| --- | --- | --- |
| 1 | **E-23 — périmètre** (définition de « modèle ») | Un périmètre, **non un cadre** : cette entrée ouvre la porte que les cinq suivantes franchissent |
| 2 | **E-23 — inventaire** exact et perpétuellement à jour | L'ensemble de ce qui peut être invoqué doit être **énumérable avant l'exécution**, donc déclaré par le cadre |
| 3 | **E-23 — cotation de risque** par modèle, le niveau d'autonomie étant un facteur qualitatif | La cote doit être **connue au point d'invocation** |
| 4 | **E-23 — cycle de vie** en cinq volets, dont la mise hors service | La mise hors service doit **s'interposer entre la décision et l'invocation** : un point de passage, non une écriture au registre |
| 5 | **E-23 — documentation** : usages approuvés, limites, dépendances, sources | L'**usage approuvé** doit être opposable, donc porté par ce qui autorise l'invocation — et non seulement consigné |
| 6 | **E-23 — surveillance continue**, nommant décision et reparamétrisation autonomes | Le cadre est **producteur d'observations avant d'être producteur de décisions** |
| 7 | **Ligne directrice sur l'IA (AMF)** — seul le calendrier est au corpus | **Rien de dérivable.** Le contenu du texte n'est pas au corpus |
| 8 | **Avis 11-348** — adaptativité après déploiement | Distinguer l'**adaptation** d'instance de l'**évolution** du modèle de processus, et soumettre la seconde à un régime propre |
| 9 | **Art. 12.1 — informer** au plus tard au moment de la décision | Le **moment de la décision** doit être un événement émis par le cadre, identifiable et daté |
| 10 | **Art. 12.1 — expliquer** raisons, facteurs et paramètres de *cette* décision | Une **trace d'instance** produite par le cadre, **non par les agents** |
| 11 | **Art. 12.1 — réviser** par une personne en mesure de défaire | Un point de reprise, et des effets aval **bornés** — donc défaisables |

Le décompte est exact : **neuf des onze entrées produisent une contrainte** — cinq issues de la ligne directrice prudentielle, trois de l'article 12.1, une de l'avis en valeurs mobilières ; **une** produit un périmètre et non un cadre ; **une seule ne produit rien** — la ligne directrice québécoise, dont le corpus ne porte que le calendrier.

L'histoire de ce tableau est, plus que son contenu, l'argument méthodologique central de l'ouvrage. Dans son premier état, il comptait **six entrées dont deux ne produisaient rien**, et l'article 12.1 fournissait trois des quatre contraintes restantes — la ligne directrice prudentielle n'y produisait qu'un périmètre, parce que le corpus n'en portait alors que la définition de « modèle ». L'extraction du texte intégral a versé les attentes opératoires ; une passe de dérivation les a traduites ; le tableau est passé de six à onze entrées. Or le chapitre avait posé, dans son premier état, une proposition que sa propre correction vérifie : **la densité de contraintes dérivables d'un texte ne mesure pas son exigence propre, mais ce que le corpus en a extrait.** L'article 12.1 dominait parce qu'on l'avait lu ; la ligne directrice prudentielle domine pour la même raison, non parce qu'elle serait devenue plus contraignante.

Deux gardes explicites accompagnent cette table, et elles en bornent l'usage.

La première : **E-23 attend ces cinq choses d'un programme**, c'est-à-dire d'un dispositif d'organisation, de rôles et de contrôles. Elle n'attend rien d'une architecture, ne nomme ni les agents ni l'orchestration, et ne dit nulle part qu'il faille construire un cadre. Ce que l'ouvrage ajoute est une observation sur les **conditions de possibilité** : sur un système agentique, un inventaire exact, une cote connue au bon moment, une mise hors service opposable cessent d'être des propriétés d'un registre pour devenir des constats *a posteriori* — et le cadre est le seul objet qui les rende vraies par construction.

La seconde est une limite structurelle qui n'est pas comblée. L'objet de la ligne directrice est **le modèle**. Un système agentique compose plusieurs modèles, des outils et un enchaînement, et le texte **ne dit pas ce qu'est « le modèle » d'un tel assemblage** — ni s'il faut coter les composants, la composition, ou les deux. Le silence est entier.

### 6.2 Une convergence à trois sources non indépendantes

Le principe de l'encadrement déterministe des processus réglementés est formulé, dans la littérature de la mi-2026, par trois sources de registres différents.

La **première** est empirique : le préprint de la Technische Universität München déjà cité, dont le verdict sur un processus régi par une directive européenne juge l'orchestration non encadrée **inacceptable** lorsque des exigences strictes d'exécution et de documentation s'appliquent, et pose que les tâches essentielles doivent être imposées de façon déterministe par le cadre.

La **deuxième** est conceptuelle : le manifeste académique sur l'*Agentic Business Process Management*, qui érige l'autonomie encadrée en mécanisme premier de gouvernance et fournit le vocabulaire des cadres normatifs et opérationnels. Il **recommande ; il ne juge pas**.

La **troisième** est industrielle : un pattern d'architecture publié par un éditeur, mis à jour le 21 février 2025, qui recommande **explicitement les flux de travail statiques pour les processus sous surveillance réglementaire**, au nom de l'auditabilité, de la conformité et de transferts définis. Là encore, l'éditeur recommande — il ne juge pas. Fait remarquable par sa provenance : c'est le fournisseur qui borne lui-même l'emploi de ce qu'il vend.

Le point décisif est le suivant, et il a été acquis par correction : **ces trois sources ne sont pas indépendantes.** L'autrice principale du préprint figure parmi les auteurs du manifeste — les deux sources académiques **partagent une autrice**. L'organisation de recherche de l'éditeur figure parmi celles dont proviennent les dix-huit auteurs du manifeste — la source académique et la source industrielle **partagent une organisation**. L'adjectif « indépendantes » figurait dans le corpus, en gras ; il en a été retiré sur signalement de deux rédactions, et le corpus a été corrigé **en premier**, avant les chapitres qui s'en réclamaient.

La formulation exacte et suffisante est donc : *le principe est formulé, dans la littérature de la mi-2026, par un manifeste académique, une expérimentation et un pattern de fournisseur, dont deux partagent une autrice et deux une organisation*. Écrire « trois sources indépendantes recommandent l'encadrement déterministe » serait faux. Trois sources qui se recoupent deux à deux valent **moins qu'une corroboration indépendante, davantage qu'une opinion isolée** : un même milieu, partiellement chevauchant, l'a formulé dans trois registres.

> **État de la connaissance vérifiable.** Le corpus n'établit, pour aucune des trois sources, d'application au Canada, à la finance canadienne, ni à aucun des textes canadiens de la section 5. Mais cette absence est **une absence au corpus, non un fait négatif établi dans les sources** : aucune recherche exhaustive de chaînes n'a été conduite, à la différence des trois faits négatifs vérifiés de l'ouvrage. Le chapitre porte d'ailleurs son propre contre-exemple — le manifeste compte une affiliation canadienne. **La transposition au contexte canadien est entièrement une inférence d'auteur.**

### 6.3 L'écart de responsabilité et l'imputabilité

Le manifeste inscrit parmi ses défis transversaux l'indétermination de l'imputabilité juridique entre quatre porteurs candidats : le développeur, **l'organisation qui impose le cadre**, le fournisseur du modèle, et le comportement émergent du système multi-agents.

Face à cette énumération, le corpus porte trois objets de nature différente. Pour l'article 12.1, un **fait générateur** : l'obligation pèse sur toute entreprise qui rend une décision fondée exclusivement sur un traitement automatisé de renseignements personnels, sans condition tenant à l'auteur du système ni partage avec le fournisseur du modèle. Pour la ligne directrice prudentielle, un **périmètre d'application** — les institutions financières fédérales. Pour l'avis en valeurs mobilières, un périmètre également, le porteur qu'il désigne l'étant par renvoi au droit en vigueur et non par l'avis.

Une constante s'en dégage, et elle est marquée comme lecture d'auteur : **trois textes, trois régimes, un même porteur — celui qui exploite le système, non celui qui l'a construit.** L'expression « l'entité qui exploite » est de l'auteur : elle ne figure ni dans les textes, ni au corpus. Rapportée à l'énumération du manifeste, cette constante désigne le **deuxième** des quatre candidats : l'organisation qui impose le cadre.

Ce que cette section refuse d'affirmer est aussi important que ce qu'elle affirme. Elle ne dit pas que le droit canadien **résolve** l'écart de responsabilité : ce que le manifeste nomme est l'indétermination de l'**allocation** entre quatre porteurs, question que ces textes ne traitent pas. Elle ne dit pas que ces régimes se **cumulent** sur une même chaîne de décision — la portée du régime québécois à l'égard des institutions sous charte fédérale est une question sur laquelle le corpus est muet, et elle décide si les deux régimes s'appliquent à des périmètres disjoints ou se superposent. Elle ne dit rien de ce que la ligne directrice prudentielle et l'avis en valeurs mobilières prévoient quant à la désignation d'un tiers : le corpus ne portant pas ce contenu, l'ouvrage **se refuse explicitement l'argument du silence**.

Ce qu'elle affirme tient en une phrase, et c'est le cœur de l'ouvrage : **l'article 12.1 n'offre pas à l'entreprise assujettie la ressource de désigner un tiers — ni le développeur, ni le fournisseur de modèle, ni l'émergence elle-même. Sous l'article 12.1 du moins, l'écart de responsabilité n'est pas un moyen de défense devant le régulateur : il est le problème de l'assujetti.**

D'où le raisonnement final. L'entreprise assujettie répond d'un comportement qu'elle n'a pas spécifié et qu'elle ne peut pas reconstituer depuis le modèle. Le seul objet dont elle puisse démontrer la teneur — parce qu'elle l'a écrit, qu'il précède l'exécution et qu'il ne dépend de la coopération d'aucun agent — est le cadre qu'elle a posé. **L'encadrement n'est pas, dans cette lecture, une bonne pratique d'ingénierie : c'est le seul objet que l'assujetti puisse produire pour répondre de ce dont il répond de toute façon.**

Réserve immédiate, et elle n'est pas rhétorique : cette lecture explique **pourquoi l'encadrement est nécessaire, non qu'il soit suffisant**. Rien n'établit qu'un cadre correctement posé libère l'assujetti, ni que la démonstrabilité d'un cadre vaille preuve de conformité.

### 6.4 Le principe directeur et ses quatre conditions

> **Sous exigence réglementaire stricte, le processus est imposé de façon déterministe par le cadre, qui invoque les agents ; les agents n'orchestrent pas le processus.**

Dans la taxonomie de la section 4, ce sont les positions **OO3 et OO4** : la main qui commande l'enchaînement est le cadre, et la différence entre les deux tient à ce que l'agent **sait** du processus, non à ce qu'il **commande**.

Quatre conditions doivent être énoncées avec ce principe, sous peine de le falsifier.

**Première — c'est une inférence d'auteur.** Elle procède par transposition de trois sources dont le corpus n'établit l'application ni au Canada, ni à la finance canadienne, et dont deux partagent une autrice et deux une organisation.

**Deuxième — « exigence réglementaire stricte » n'est pas défini par le corpus.** Le verdict empirique borne son propre domaine à un processus européen, hors finance ; qualifier un processus canadien de strictement réglementé demeure un travail d'institution, non un acquis de la littérature.

**Troisième — l'encadrement se paie.** Effort initial et maintenance figurent explicitement parmi les sept critères de sélection. Un principe qui ignorerait son coût ne serait pas un principe d'architecture.

**Quatrième — le principe dit où placer la main qui commande ; il ne dit ni ce qu'un cadre doit contenir, ni comment on l'écrit.** Les cinq attentes de la ligne directrice prudentielle disent ce que le cadre doit **rendre possible**, non ce qu'il doit décider. Et de la ligne directrice québécoise — l'un des deux textes dont l'entrée en vigueur commune fixe l'échéance —, cette section ne dit rien du tout.

L'ouvrage referme ce chapitre par une phrase qui en donne le statut exact : *le pont n'est pas une déduction ; c'est un raisonnement, exposé pour être contesté — et sa partie la plus solide n'est ni la convergence de trois sources, ni la constante lue dans trois textes : c'est ce qu'un seul texte, l'article 12.1, fait peser sur l'entreprise qui rend la décision.*

## 7. L'interopérabilité financière canadienne

### 7.1 Le cadre bancaire : un fait négatif vérifié au cœur d'un dispositif

Le Canada s'est doté d'une première loi partielle sur les services bancaires axés sur le consommateur en 2024. Le projet de loi C-15, déposé le 18 novembre 2025 et **sanctionné le 26 mars 2026**, l'a **abrogée et remplacée** — non amendée. La loi complète couvre l'accréditation, la sécurité, la sécurité nationale, la responsabilité et le consentement. Le corpus documente la sanction royale, **non l'entrée en vigueur article par article** : l'ouvrage se refuse à dire que la loi est intégralement en vigueur. Second acquis, plus large : en modifiant la loi fédérale sur la protection des renseignements personnels, C-15 crée un **droit à la mobilité des données à l'échelle de l'économie**, dont le cadre bancaire est la première itération.

La supervision et la mise en œuvre ont été déplacées vers la **Banque du Canada** — l'agence de protection du consommateur en matière financière conservant son mandat général sous la loi sur les banques. Une institution fédérale relève donc de **trois autorités distinctes** sur ce périmètre : la banque centrale pour le cadre, le superviseur prudentiel pour le risque de modélisation, l'agence de protection du consommateur pour la relation client. L'accréditation est la porte d'entrée du dispositif : un acte administratif préalable, non une déclaration de conformité, dont les modalités restent à définir par règlement.

Le règlement a été **prépublié le 27 juin 2026** à la Gazette du Canada, partie I, au lendemain d'un communiqué ministériel ; la période de commentaires de soixante jours s'est close le **26 août 2026**. Le texte étant prépublié, il est **susceptible de changer**. L'entrée en vigueur sera échelonnée — l'accréditation d'abord, puis les règles communes et les frais d'évaluation dans l'année suivant la publication finale, le phasage par produits étant précisé à cette publication. Aucune de ces échéances n'est datée en valeur absolue : toutes sont ancrées sur une publication finale qui n'a pas de date.

Reste le point le plus important, et c'est un **fait négatif vérifié**. La loi impose un **standard technique unique**, fixé par un organisme de normalisation désigné par le ministre des Finances par arrêté ministériel, selon des critères énoncés — être significativement canadien, avoir une gouvernance équitable, ouverte et accessible, et décider de façon indépendante —, et assortis de quatre objectifs : sécurité, concurrence, innovation, interopérabilité mondiale.

La formulation à reprendre est la suivante : **au 16 juillet 2026, aucun organisme de normalisation technique n'a été désigné par arrêté ministériel et aucun standard n'est nommé dans les textes officiels ; la candidature d'un organisme d'industrie relève du commentaire d'industrie.** Le constat n'est pas une impression : recherche exhaustive de quatre chaînes de caractères — le sigle et le nom développé de l'organisme pressenti, ainsi que deux profils de sécurité couramment associés au domaine — dans le règlement prépublié, le communiqué ministériel et la page budgétaire ; **zéro occurrence**, constat voté à neuf voix contre zéro et reconduit lors d'une revalidation ultérieure sur les sites officiels. Écrire qu'un organisme donné est le standard retenu du cadre bancaire canadien est proscrit.

C'est là, dans tout l'ouvrage, l'exemple le plus net de ce qu'un fait négatif vérifié apporte de plus qu'une absence de documentation : il autorise une **décision d'architecture** — traiter le standard d'accès comme une variable derrière une frontière d'abstraction — là où une simple ignorance n'autoriserait rien.

### 7.2 ISO 20022 : un rail accompli, un rail visé

La couche sémantique commune des paiements canadiens est en place pour les paiements de grande valeur. Le **22 novembre 2025**, la coexistence entre l'ancien et le nouveau format de message a pris fin sur le système de paiement de grande valeur : les messages hérités ne sont plus pris en charge, au profit des messages ISO 20022 de virement client et de virement interbancaire. L'opérateur a annoncé l'achèvement le 26 novembre 2025. La bascule est **alignée sur l'échéance mondiale** de la messagerie transfrontalière — le corpus établit l'alignement, non la causalité. Elle avait été ouverte trente-deux mois plus tôt, et plus de 98 % des messages étaient déjà au nouveau format en octobre 2025, **selon l'opérateur**.

La valeur de cette date est précise et modeste : c'est la **disparition d'une alternative** à la frontière du rail, non une nouveauté de format. Le corpus ne dit rien des chaînes internes des institutions.

Table: Tableau 10 — Le *Real-Time Rail* : jalons vérifiés et statut au 17 juillet 2026.

| Jalon | Date | Statut |
| --- | --- | --- |
| Livraison de la composante d'échange | juin 2023 | Fait |
| Reprise du programme avec de nouveaux partenaires de livraison et d'exploitation | 16 avr. 2024 | Fait |
| Composante de compensation et de règlement complétée | T3 2025 | Fait |
| Tests d'intégration internes | T4 2025 | Fait |
| Tests d'acceptation utilisateur | T1 2026 | Fait |
| Tests industriels d'assurance de la solution | en cours à la mi-2026 | En cours |
| Règlement administratif no 10 (Gazette, partie II, 1er juill. 2026) | en vigueur le 24 août 2026 | Publié |
| Lancement, ISO 20022 dès l'origine, déploiement graduel se poursuivant en 2027 | **cible : T4 2026** | **Visé** |

Deux formulations symétriques doivent être tenues ensemble, et l'ouvrage insiste parce que chacune, prise seule, est fausse. **À la mi-juillet 2026, le système n'est pas en production** — on n'écrit ni « lancé », ni « en production ». **Et l'opérateur vise officiellement un lancement au quatrième trimestre de 2026** — on n'écrit pas non plus « aucune date annoncée », affirmation réfutée à l'unanimité. La cible a été successivement fixée à 2019, puis 2022, puis 2023, puis 2026 : ce sont les cibles successives, non les dates auxquelles les reports ont été décidés.

> **État de la connaissance vérifiable.** Sur ce que la couche sémantique commune change pour des agents, le corpus est plus pauvre qu'on ne l'attendrait. Il nomme les types de messages du rail de grande valeur et établit qu'ISO 20022 sera présent dès le lancement du rail temps réel, **sans nommer de type de message** pour ce dernier ; il **ne documente aucun attribut de richesse des données** — ni structure des données de remise, ni granularité des identifiants, ni comparaison avec le format retiré. Toute affirmation sur ce que les messages « portent de plus » excéderait le corpus. Ce qui est autorisé se réduit à un constat de surface, et à une asymétrie : **l'un des deux rails est un acquis, l'autre est une cible** — deux régimes de preuve qu'il ne faut pas additionner.

### 7.3 AP2 sur les rails canadiens : trois conditions, aucune remplie

La question se pose naturellement : le protocole de la transaction agentique s'articulera-t-il aux rails de paiement canadiens ? La réponse de l'ouvrage est un refus argumenté d'affirmer, dans les deux sens.

Le corpus ne documente **aucune articulation**. Cette lacune est un **défaut de documentation, non un fait négatif vérifié** : contrairement au cas du standard technique, aucune recherche exhaustive de chaînes n'a été menée dans les sources primaires du protocole ni du rail. Le corpus ne porte d'ailleurs, sur ce protocole, ni structure de message, ni mécanique de mandat, ni modèle de preuve d'intention. Le chapitre correspondant a été explicitement réorienté pendant la rédaction : il devait proposer des « scénarios d'articulation », il livre des **conditions de possibilité** — parce que des scénarios, dans un chapitre dont la thèse est qu'aucune source ne documente l'articulation, combleraient la lacune par de la fiction.

Trois conditions sont énumérées, dans leur ordre d'antériorité logique.

**Le rail doit exister.** Il ne l'est pas : le système est en tests industriels. L'asymétrie documentaire entre les deux objets est frappante — le rail a une chronologie jalonnée par des communications d'opérateur ; le protocole n'a que **deux points** : son annonce et un décompte de soutiens six mois et vingt-quatre jours plus tard, sans jalon de version, sans étape de maturité, sans déploiement. Le corpus n'établit pas que le protocole soit immature ; il établit que sa maturité n'est pas documentée.

**Les couches sémantiques doivent se rejoindre.** Le corpus n'établit **aucune correspondance** entre le protocole et la norme de message. Toute traduction d'un mandat de paiement agentique en message normalisé serait une invention.

**La gouvernance doit être documentable des deux côtés.** Côté rail, on a les partenaires et les dates du règlement administratif, mais **rien sur les règles d'accès d'un tiers** — et l'un des acteurs majeurs de l'écosystème canadien n'est documenté qu'en qualité de partenaire du programme, rien sur ses rails propres, de sorte que ce versant de la question ne peut même pas être instruit. Côté protocole, c'est celui-là même que la section 3.1 a dû **exclure** de son constat de consolidation, faute de tout transfert de gouvernance documenté.

Une quatrième famille de conditions — les conditions réglementaires — est signalée et renvoyée à la section 9.

Cinq questions de recherche closent ce point, et chacune est formulée avec ce qui la trancherait : l'un des soutiens déclarés du protocole opère-t-il ou participe-t-il à un rail canadien (la liste complète le dirait) ; la gouvernance du protocole fait-elle l'objet d'un transfert (un communiqué le dirait) ; une correspondance avec la norme de message est-elle documentée (la spécification du protocole, absente du corpus, le dirait) ; le règlement administratif du rail traite-t-il de l'initiation d'un paiement par un tiers ou par un agent (le texte publié le dirait) ; le lancement visé a-t-il lieu (l'opérateur le dira). Ces cinq questions sont **distinctes** de celle de la désignation du standard technique : les confondre reviendrait à fusionner deux ignorances que le corpus tient séparées — l'une est un fait négatif vérifié, l'autre une absence de documentation.

## 8. L'adoption : la production agentique canadienne

### 8.1 Dix institutions, une couverture inégale

Cette section rapporte ce que dix institutions financières canadiennes ont rendu public de leurs déploiements entre 2025 et le milieu de 2026. Une règle transversale la gouverne sans exception : **toutes les métriques ci-dessous sont auto-déclarées** — communiqués, présentations aux investisseurs, rapports annuels — et ne font l'objet d'aucune vérification indépendante. Leur attribution est portée à chaque occurrence.

Cinq appuis extérieurs seulement existent, sur cinq institutions distinctes : un prix sectoriel, une publication de gestion universitaire, un indice d'analyste, une étude de cas d'un institut de recherche partenaire, et une fondation d'écosystème. Les rapports annuels de deux institutions n'en sont pas : le support change le destinataire, pas l'auteur.

L'accès aux sources primaires a lui-même été inégal, et l'ouvrage le documente plutôt que de le lisser : cinq institutions ont livré un accès direct ; quatre un accès dégradé — page rendue par script, refus de chargement automatisé, document trop volumineux ; une n'a livré **aucune source primaire**.

Table: Tableau 11 — La production agentique canadienne, 2025-2026 (toutes les métriques sont auto-déclarées par l'institution nommée).

| Institution | Système | Métrique déclarée | Gouvernance documentée |
| --- | --- | --- | --- |
| TD | Premier modèle d'IA agentique (laboratoire interne), pré-adjudication du prêt garanti par l'immobilier et mémos de synthèse pour souscripteurs (communiqué du 21 mai 2026) | Traitement d'environ **quinze heures à moins de trois minutes**, résultats qualifiés de préliminaires | Équipe dédiée, **cinq critères** dont l'explicabilité ; prix sectoriel 2025 confirmé à sa source |
| Banque Scotia | Outil interne d'analyse documentaire en service depuis 2018, breveté, **capacités agentiques greffées en 2025** ; traitement autonome des courriels de clients en services aux entreprises | **≈ 90 % d'environ 1 500 courriels par jour** ; acheminement d'une à deux heures à quelques minutes ; escalade humaine des cas complexes | **Revue éthique préalable de chaque modèle, approbation graduée par niveau de risque** ; chef de l'IA nommé ; corroborée par une publication tierce |
| Banque Royale du Canada | Groupe IA rattaché au chef de la direction (18 févr. 2026) ; modèle de fondation interne, assistant employés, outil de marchés des capitaux, plateforme de calcul | Cible de **700 M$ à 1 G$** de valeur d'entreprise d'ici 2027 ; assistant à ~27 000 employés ; outil de marchés à 8 000 | Principes d'IA responsable dont la **transparence** ; seule institution engagée dans un effort inter-institutionnel documenté |
| Manuvie | Plateforme d'IA agentique d'entreprise ; environnement d'exécution retenu, alors en version bêta (communiqué conjoint du 10 mars 2026) | Cible de **plus de 1 G$** d'ici 2027, dont environ un cinquième en efficacité ; **≈ 300 M$ déclarés réalisés fin 2025** | **Chef de l'IA au comité de direction exécutive** depuis mai 2026 ; premier rang d'un indice d'analyste dans son secteur |
| Desjardins | Aucun système agentique en production documenté | **Aucune métrique de déploiement, aucune cible chiffrée** | IA générative et assistant bureautique en priorités ; l'ancrage dans les données et l'IA est l'une des six orientations du plan stratégique 2026-2029 |
| CIBC | Plateforme interne d'IA générative déployée à l'échelle de la banque (communiqué du 27 mai 2025), en pilote depuis juillet 2024 | **Plus de 200 000 heures** économisées durant le pilote ; accès conditionnel à une formation obligatoire | Se déclare première grande banque canadienne signataire du code de conduite volontaire fédéral (mars 2025) |
| Intact | Capacité analytique et prédictive à l'échelle (rapport annuel publié le 26 mars 2026) | **≈ 600 modèles en production** ; **plus de 200 M$** de bénéfices annuels attribués ; **plus de 90 %** du portefeuille personnel canadien tarifé par apprentissage machine | Aucune terminologie agentique attribuée — **fait établi par extraction intégrale** |
| BMO | Assistant d'IA générative interne (mai 2025) | **≈ 8 000 employés de première ligne** utilisateurs — le chiffre désigne des **employés**, non des politiques, et provient d'une soumission distincte du communiqué | Formations internes sur l'agentique attestées **par l'institut partenaire, non par la banque** |
| Sun Life | Membre d'un consortium annoncé le 7 juillet 2026 dont le programme phare est déclaré **en production en environnement réglementé** | **Plus de deux mille milliards de jetons par mois**, agrégés à l'échelle des organisations membres | Effort collectif d'infrastructure ; aucune imputation possible à une institution en particulier |
| Banque Nationale | — | — | **Aucune source primaire remontée** ; le seul élément identifié n'est attesté que par des sources secondaires |

### 8.2 Ce que ces cas établissent, et ce qu'ils n'établissent pas

Trois lectures se dégagent, et toutes trois sont bornées.

**Une chronologie.** L'année 2025 est celle de l'assistance — plateformes d'IA générative déployées à l'échelle, assistants internes. L'année 2026 est celle de la structure et du cœur de métier : du 18 février au 7 juillet 2026, quatre mois et dix-neuf jours séparent quatre annonces structurantes — création d'un groupe rattaché au chef de la direction, choix d'un environnement d'exécution agentique, premier modèle agentique en pré-adjudication de crédit, lancement d'un consortium d'infrastructure.

**Une gradation.** Une grille de lecture — de l'ouvrage, non du corpus — distingue la production déclarée, la structuration outillée et l'intention formalisée. Elle ne classe pas des maturités : elle classe des **régimes de documentation publique**.

**Une limite commune, et c'est le constat le plus instructif.** **Aucune des dix institutions ne publie de taux d'erreur, d'incident ou de reprise humaine hors escalade prévue, de périmètre de journalisation, ni de rattachement de ses contrôles à une exigence réglementaire nommée.** Cinq publient un gain chiffré ; les cinq autres publient une cible, une population d'utilisateurs, un volume de consommation, ou rien. Un seul mécanisme de gouvernance est décrit à un niveau opérationnel — une revue éthique préalable graduée par le risque —, et c'est le seul corroboré par un tiers.

Plusieurs précautions arithmétiques accompagnent ces chiffres et méritent d'être reprises, car elles illustrent la discipline d'attribution mieux qu'un principe général. Deux populations d'utilisateurs déclarées par une même institution **ne s'additionnent pas**, le recoupement n'étant pas documenté. Une fourchette de cible se cite **entière** : retenir la borne haute surestimerait la borne basse de plus de quarante pour cent. Un volume quotidien de courriels et un volume mensuel de demandes déclarés par une même source **ne se convertissent pas l'un en l'autre** — les unités diffèrent, et l'écart entre les deux ordres de grandeur l'établit. Un nombre d'experts et un nombre de modèles ne se croisent pas en ratio. Un volume de jetons mesure une **consommation**, non une valeur ni une correction. Et une cible n'est pas un gain.

Deux précautions qualitatives, enfin. La première : une plateforme d'IA générative assistive **ne doit pas être surqualifiée** en agentique quand son propre communiqué ne mentionne ni agents autonomes, ni protocole d'invocation d'outils. La seconde, symétrique et plus subtile : un classement bâti par recherche lexicale du mot « agentique » dans les communications d'entreprise serait un **artefact de méthode** — une institution dont les communications n'emploient pas le terme peut exploiter des centaines de modèles en production, et c'est le cas d'au moins une du tableau.

> **État de la connaissance vérifiable.** Trois questions restent ouvertes sur cette strate : le traitement de l'IA générative comme facteur de risque matériel dans le rapport annuel d'une institution (document non extractible, dépôts réglementaires refusés) ; l'existence d'agents en adhésion, en production de relevés fiscaux et en traitement des réclamations chez une autre (source restante secondaire et sous péage, à attribuer explicitement à cette source si elle est citée) ; et l'activité agentique d'une dixième institution, pour laquelle aucune source primaire n'a été remontée. **L'absence de sources primaires n'est pas l'absence d'activité** : cette section documente une couverture publique, non un état interne.

Un dernier avertissement terminologique s'impose ici, et il renvoie à la section 3.5 : le programme du consortium annoncé le 7 juillet 2026 porte un nom **homonyme** de celui d'un protocole fusionné dans A2A en août 2025. Le corpus établit, sur ce seul couple, l'absence de tout lien — pure homonymie. Le sigle employé seul reste proscrit.

## 9. Synthèse architecturale

### 9.1 La matrice : quinze croisements, aucun lien documenté

L'exercice consiste à croiser les trois protocoles majeurs avec les cinq textes canadiens de la section 5. Trois règles de remplissage le bornent : une cellule n'est renseignable que si le corpus porte **les deux termes** — ce que le texte attend *et* ce que le protocole fournit ; aucune cellule ne porte de verdict de conformité ; et **toute cellule non vide est une inférence d'auteur**.

Table: Tableau 12 — Matrice protocoles × textes canadiens : quinze croisements (état arrêté au 16 juillet 2026).

| Texte canadien | MCP | A2A v1.0 | AP2 |
| --- | --- | --- | --- |
| **E-23** (gestion du risque de modélisation) | Fait négatif vérifié : le texte ne nomme aucun agent | *idem* | *idem* |
| **Ligne directrice sur l'IA (AMF)** | Rien de croisable — contenu hors corpus | *idem* | *idem* |
| **Avis 11-348** (adaptativité après déploiement) | Rien sur l'auto-modification | *idem* | *idem* |
| **Art. 12.1, Loi 25** | Données typées : substrat de frontière, **non** trace d'instance | Cartes signées : identité de l'appelant, **non** explication | — |
| **Cadre bancaire** (standard technique) | Fait négatif vérifié sur le lexique d'autorisation | — | — |

Le résultat est net : **aucun des quinze croisements n'est documenté par une source primaire**. Ce n'est pas un tableau de correspondances, c'est un tableau de vides — et l'intérêt est dans la nature de ces vides, qui n'est pas uniforme.

Deux décomptes distincts, qu'il ne faut pas confondre, s'appliquent. Par **provenance de l'absence** : quatre absences ont été **cherchées dans la source** — les trois cellules de la ligne prudentielle, via le balayage mécanique du texte, et une cellule du cadre bancaire, via la recherche exhaustive de chaînes — et onze sont simplement **constatées au corpus**. Par **état de cellule** : six cellules portent un contenu, neuf sont vides. Le vide est celui du **lien**, jamais celui de la cellule.

Trois **espèces de vide** se répartissent sur les cinq lignes, et cette typologie est l'apport propre de l'exercice.

Le **vide de corpus** — une seule ligne, la ligne directrice québécoise : l'ouvrage ne sait pas ce que le texte attend. C'est la seule ligne dont le **premier terme** manque, et le seul vide qu'une extraction primaire comblerait.

Le **vide de texte** — une ligne, le cadre bancaire : le corpus porte l'exigence d'un standard technique unique, mais le texte officiel ne désigne rien. Ce n'est pas une lacune de l'ouvrage : c'est un état du droit, et il est vérifié.

Le **vide de protocole** — trois lignes : le corpus porte la demande, les protocoles ne portent rien.

D'où le renversement, qui est le résultat le plus utile de la matrice : **les deux textes dont le corpus porte le contenu le plus précis — la ligne directrice prudentielle, extraite intégralement, et l'article 12.1, cité dans son libellé officiel — sont ceux dont les six croisements ne portent aucun lien documenté. Et ce sont eux qui obligent le plus tôt.** Pour les deux protocoles d'interaction, le vide n'est pas documentaire : leurs spécifications décrivent des appels, pas des décisions. Pour le protocole de transaction, rien ne permet de trancher entre les deux espèces.

**Cinq zones de compensation architecturale** se déduisent de cette lecture. Une zone de compensation est un point où le corpus établit une demande, ou une contrainte de conception, que la couche protocolaire ne porte pas — et que quelque chose d'autre doit donc porter.

**Z1 — la production de la trace.** L'article 12.1 exige les raisons, facteurs et paramètres de *cette* décision ; le préprint écarte les agents comme producteurs de journal ; le rapport conjoint des régulateurs écarte le modèle, la causalité étant indéterminable ; la matrice écarte le protocole. Quatre candidats, trois éliminés.

**Z2 — le point d'arrêt humain.** Aucun humain n'est documenté dans les trois protocoles ; l'humain-dans-la-boucle et les points de reprise sont documentés au niveau d'**un** cadriciel d'orchestration — le singulier est exact.

**Z3 — la détectabilité de l'auto-modification.** L'avis en valeurs mobilières nomme l'adaptativité après déploiement ; le manifeste distingue adaptation et évolution ; la ligne directrice prudentielle vise la reparamétrisation autonome. Aucune ligne de la pile protocolaire ne l'exécute.

**Z4 — l'interface d'accès du cadre bancaire.** Le standard n'étant pas désigné, la discipline consiste à traiter le standard d'accès comme une variable derrière une frontière d'abstraction. Ce n'est pas un choix de protocole agentique.

**Z5 — le contenu et l'état.** Les réponses protocolaires ne couvrent aucune des trois surfaces d'attaque de la section 3.6. Le manifeste propose le confinement local comme frontière de sécurité : c'est une limitation des dégâts, non une prévention.

Leur propriété commune est décisive : **aucune ne se referme en choisissant un meilleur protocole.** La compensation est structurelle, non palliative.

Deux réserves closent cet exercice. La matrice **n'établit pas que les protocoles échouent** — la conformité n'est pas une propriété qu'un format d'échange puisse porter, et les protocoles répondent simplement à une autre question. Et les cinq zones **ne sont pas exhaustives** : une extraction primaire de la ligne directrice québécoise en ferait vraisemblablement apparaître d'autres.

### 9.2 L'architecture de référence : quatre couches

Une architecture de référence n'est pas une déduction du corpus : c'est un **ordonnancement**. Celle-ci est neutre — aucun composant d'éditeur n'y est prescrit.

Table: Tableau 13 — L'architecture de référence par couches.

| Couche | Responsabilité | Ce qu'elle ne porte pas |
| --- | --- | --- |
| **Protocoles** | Format et habilitation à la frontière : qui parle, par quel format | Le processus, la trace d'instance, la décision, l'humain |
| **Identité et registre** | Nommer l'agent ; borner ses outils et ses droits **hors de lui** | La conformité (statuts pré-normatifs) ; l'inventaire attendu par la ligne directrice prudentielle, qui est un objet distinct |
| **Orchestration** | Tenir l'enchaînement ; produire la trace ; arrêter pour l'humain | La qualification juridique du processus ; le contenu des cadres |
| **Gouvernance** | Inventaire, cotation, cycle de vie, documentation, surveillance | Ce que la ligne directrice québécoise attend |

Deux rapprochements structurent cette table, et l'un d'eux est un rapprochement **à écarter**.

Celui qu'il faut faire : les deux champs obligatoires du registre d'agents — outils atteignables et bornes de privilège — constituent un **cadre opérationnel écrit hors de l'agent**. C'est une inférence, mais elle donne à la couche d'identité une fonction qui excède l'authentification.

Celui qu'il faut écarter : **le registre d'agents n'est pas l'inventaire attendu par la ligne directrice prudentielle.** Le premier recense des agents ; le second recense des modèles au risque inhérent non négligeable. Ce sont deux objets, et les confondre produirait un inventaire qui ne satisferait ni l'un ni l'autre.

Une asymétrie mérite d'être relevée : la couche de gouvernance est spécifiable contre un texte — la ligne directrice prudentielle, dont les attentes ont été extraites — et pas contre l'autre — la ligne directrice québécoise, dont le corpus ne porte que le calendrier. Or les deux entrent en vigueur le même jour.

Le positionnement d'orchestration se décide ensuite, cas d'usage par cas d'usage, sous le principe directeur de la section 6 : sous exigence réglementaire stricte, positions OO3 ou OO4. Trois classes illustratives — et il faut insister sur le mot : **le corpus ne documente l'option d'orchestration d'aucun déploiement agentique canadien**, ces classes ne sont donc pas des architectures réelles.

La **classe « décider et expliquer »** — une décision de crédit produisant un mémo à un souscripteur — sature cinq des sept critères de sélection et appelle OO3 ou OO4.

La **classe « qualifier »** — un acheminement autonome de courriels — est la plus instructive par ce qu'elle rend indécidable : un acheminement est-il une décision fondée exclusivement sur un traitement automatisé ? La question n'est pas technique. **La qualification précède l'architecture.**

La **classe « exécuter »** — un paiement de grande valeur sur un rail normalisé — pose une contrainte d'exécution et non d'explication : l'agent observe, le rail reste déterministe.

Les trois classes se distinguent donc par la nature de ce qui les contraint : **expliquer, qualifier, exécuter**.

### 9.3 Cinq points de contrôle obligatoires

Quatre découlent des zones de compensation ; le cinquième vient de la table de traduction. La quatrième zone n'en produit aucun, faute d'objet — le standard n'étant pas désigné.

Table: Tableau 14 — Les cinq points de contrôle obligatoires et leur origine.

| # | Point de contrôle | Origine | Ce qu'il rend possible |
| --- | --- | --- | --- |
| **PC1** | L'**événement de décision** : un événement daté, émis par la couche d'orchestration | Art. 12.1 — obligation d'informer au plus tard au moment de la décision | La seule des trois obligations qui ne dépende d'aucune demande |
| **PC2** | La **trace d'instance produite par le cadre**, non par les agents | Z1 — trois producteurs candidats écartés | L'obligation d'expliquer, portant sur l'instance |
| **PC3** | Le **point d'arrêt humain**, à effets aval **bornés donc défaisables** | Z2 et art. 12.1 — révision par une personne en mesure de défaire | La réversibilité ; **ne remplace pas** la révision *a posteriori* |
| **PC4** | La **séparation de l'adaptation et de l'évolution** : deux chemins techniques, deux régimes d'autorisation | Z3, avis 11-348 et manifeste | Rendre détectable le moment où une exception devient une règle |
| **PC5** | Le **confinement local** : cadre opérationnel restreignant contexte et capacités de chaque agent | Z5 et manifeste | Limiter l'impact d'un agent compromis |

Deux observations closent ce dispositif.

La première est une tension qu'il faut documenter plutôt que déclarer résolue : **ce que PC2 exige d'exposer, PC5 exige de restreindre.** C'est le paradoxe de confidentialité de l'explicabilité, nommé par le manifeste et non résolu par lui. Un système qui déclarerait avoir tranché cette tension sans en exposer l'arbitrage aurait simplement omis de la voir.

La seconde : ces cinq points relèvent tous des propriétés que l'exploitant **démontre à un tiers** — spécificité de tâche, assurance de correction, réactivité, traçabilité. La cinquième propriété, l'autonomie, n'a au corpus **aucune métrique**. Ce n'est pas un hasard de présentation : ce qui se gouverne est ce qui se démontre.

### 9.4 Instrumentation et feuille de route vers le 1er mai 2027

Les métriques issues du préprint sont **des candidatures, pas des correspondances**. Trois rapprochements sont plausibles et aucun n'est documenté : la précision, le rappel et le score F1 face au suivi de la performance attendu par la ligne directrice — mais l'un est une mesure de laboratoire sur un scénario, l'autre une attente de dispositif sur un parc ; la vitesse de réaction et le taux de faux négatifs face aux seuils de dépassement — le texte ne dit pas de quelle grandeur ; la correction du journal face à la documentation de modèle et à l'inventaire — objets d'une tout autre nature qu'un journal d'exécution.

La feuille de route, en revanche, se séquence sans dépendre de ces rapprochements. Trois mouvements, dans un ordre qui n'est pas indifférent. À la date de gel, il restait **neuf mois et quinze jours** avant l'entrée en vigueur commune, et le texte prudentiel ne comporte **aucune disposition transitoire**.

**Inventorier.** L'attente est explicite et détaillée jusqu'aux champs : identifiant, nom, cote de risque, propriétaire, développeur, origine ; puis, pour les modèles au risque non négligeable, version, date de déploiement, réviseur, approbateur, dépendances, sources de données, usages approuvés, limites, date de revue, état de surveillance, prochaine revue. Ce mouvement ne dépend d'aucune réserve. Sa difficulté propre est le **périmètre** : traiter la couverture agentique comme acquise **par prudence**, et documenter qu'on le fait ainsi plutôt que par obligation établie.

**Coter, puis encadrer.** La cotation précède l'encadrement parce qu'elle le calibre : l'intensité du dispositif doit être proportionnée au risque. Une institution qui voudrait industrialiser ce calibrage devra construire sa propre pondération et l'assumer — le texte n'en fournit aucune.

**Surveiller.** C'est le seul mouvement qui hérite des réserves sur l'instrumentation. Il emporte néanmoins une décision d'architecture immédiate, indépendante de tout chiffre : la journalisation n'étant pas recommandée aux agents, **le lieu de production de la trace est une décision à prendre avant le 1er mai 2027, non après**.

La gradation de solidité doit être énoncée avec la feuille de route : inventorier et coter sont tracés au texte ; encadrer dépend d'une taxonomie à source unique ; surveiller hérite en outre des réserves sur les métriques. **Une institution qui commencerait par l'instrumentation commencerait par le maillon faible.**

> **Mise en œuvre.** Cinq jalons externes conditionnent cette trajectoire. La révision majeure de la spécification MCP, douze jours après le gel. L'entrée en vigueur du règlement administratif du rail temps réel, le 24 août 2026, sur fond de lancement visé au quatrième trimestre. La clôture des commentaires sur le règlement du cadre bancaire, le 26 août 2026 — le texte final peut changer. **L'arrêté ministériel désignant l'organisme de normalisation technique, sans date** : le jalon dont l'échéance est la moins prévisible et la portée la plus grande. Et l'entrée en vigueur commune du 1er mai 2027.

## 10. Validation appliquée : un blueprint instancié

### 10.1 Démarche, portée et réserve de neutralité

La dernière partie éprouve l'architecture de référence par construction, en l'instanciant sur le portefeuille d'un éditeur — IBM — documenté par sources primaires. Le choix tient à un critère unique et vérifiable : sa documentation publique permettait de tracer chaque composant à une source datée, avec son statut. **Ce n'est ni une recommandation d'achat, ni un verdict comparatif** ; l'éditeur est ci-après désigné comme tel, la démonstration ne dépendant en rien de son identité.

Deux réserves de lecture sont posées d'emblée. Les pages du site principal de l'éditeur refusent le chargement automatisé, de sorte qu'une partie des faits repose sur des extraits indexés citant l'adresse primaire — une relecture humaine est recommandée avant toute réutilisation. Et **aucune source du corpus ne positionne un produit sur l'échelle d'orchestration** : tout positionnement est, sans exception, une lecture de l'auteur.

Six principes directeurs gouvernent l'instanciation. **L'autonomie encadrée par construction** : tout processus sous exigence réglementaire stricte s'exécute en OO3 ou OO4, les contraintes réglementaires devenant des cadres opérationnels. **Aucune interaction non gouvernée** : chaque appel de modèle, chaque invocation d'outil, chaque échange inter-agents transite par une passerelle d'application de politiques. **Le contrat d'abord** : les actifs d'intégration existants sont republiés comme outils plutôt que réécrits. **Hybride et souveraineté** : charges sensibles auto-gérées, résidence canadienne des données. **La traçabilité de bout en bout, jamais déléguée aux agents** : le journal appartient à la plateforme. **La gouvernance du cycle de vie des modèles et des agents** : inventaire, évaluation graduée, catalogue gouverné.

Un mot sur le statut de ces six énoncés, car l'ouvrage le donne lui-même : trois énoncent une discipline d'architecture, deux imposent un point de contrôle, un énonce une contrainte de territoire. **Seul le premier dispose de trois sources convergentes** — avec la réserve de non-indépendance de la section 6.2. Les cinq autres sont des choix, appuyés sur des capacités documentées.

### 10.2 Huit couches, et un fait plus lourd que les autres

Table: Tableau 15 — Le blueprint en huit couches : nature, jalons datés et positionnement d'orchestration (tout positionnement est une lecture de l'auteur).

| # | Couche | Jalons datés au corpus | Positionnement |
| --- | --- | --- | --- |
| C1 | Exposition et gouvernance des interactions | Passerelle d'IA (2024, étendue 2025) ; agent d'API en disponibilité générale le 19 nov. 2025 ; passerelle compacte le même jour ; passerelle d'interaction **annoncée** en 2026 ; registre-passerelle libre | Aucun — s'applique à toutes les options |
| C2 | Intégration applicative et interentreprises | Cadence trimestrielle du moteur d'intégration ; nœud d'appel de modèle ; plateforme hybride en disponibilité générale le 16 juin 2025 | Flux déterministes = cadres opérationnels, **OO3** |
| C3 | Messagerie et événements | Messagerie en disponibilité générale juin 2026 (haute disponibilité, cryptographie post-quantique pour le transport) ; **deux produits d'événements dépréciés** ; trajectoire reportée sur une filiale acquise | Rail déterministe **OO3** ; agents événementiels sans cadre explicite **OO1** |
| C4 | Données | Plan de contrôle unifié d'intégration de données (11 juin 2025) | Aucun — couche d'alimentation |
| C5 | Couche agentique | Trousse de développement mai 2025 ; outils d'invocation dès mai 2025 ; support inter-agents depuis le 30 juin 2025 ; serveurs distants sept. 2025 ; garde-fous documentés | Flux de travail = cadre opérationnel, **OO4** — sous réserve que l'agent connaisse le processus, ce que le corpus n'établit pas |
| C6 | Gouvernance de l'IA et risque de modèle | Produit en disponibilité générale déc. 2023 ; gouvernance agentique livrée en 2025 ; module de conformité par accord d'origine externe (28 avr. 2025) | Aucun — couche transverse |
| C7 | Observabilité et opérations | Observabilité des agents et des modèles en **préversion publique** ; télémétrie ouverte native dans deux composants | Aucun — instrumentation des propriétés d'évaluation |
| C8 | Socle d'exécution et souveraineté | Plateforme en disponibilité générale le 30 juin 2026, support six ans ; deux régions canadiennes ; socle souverain en disponibilité générale le 5 mai 2026 ; cadre de contrôles infonuagiques sectoriel | Aucun — couche d'exécution |

Une convention de lecture gouverne ce tableau : le statut de disponibilité n'est indiqué que lorsque le corpus l'établit sous l'une des trois formes attendues — disponibilité générale, préversion, déprécié. Et **« annoncé » date une communication, non une disponibilité** — nuance qui touche précisément le composant sur lequel repose le deuxième principe.

Trois observations méritent d'être isolées.

**La dépréciation est le fait le plus lourd du blueprint.** Deux produits d'événements sont officiellement dépréciés, sans nouvelle fonctionnalité annoncée, la trajectoire étant reportée sur une filiale acquise — acquisition annoncée le 8 décembre 2025 et clôturée le 17 mars 2026. Pour une institution qui aurait bâti sur ces produits, cela vaut décision de migration. Or le corpus n'établit **ni calendrier de fin de support, ni chemin de migration, ni équivalence fonctionnelle**.

**Le point de contrôle unique du deuxième principe est en réalité une famille de points à coordonner.** Le corpus documente au moins quatre dispositifs de passerelle ou de registre ; il n'établit ni leur intégration mutuelle, ni qu'une politique écrite une fois s'y applique partout.

**Une capacité annoncée en disponibilité générale peut ne pas être comprise dans le droit d'usage de la plateforme** : le corpus documente au moins un cas où l'usage agentique de l'interface de plateforme exige un droit distinct. Le fait est mineur en apparence ; il ne l'est pas dans un dossier d'estimation.

La neutralité fournisseur se vérifie enfin à cinq obligations, dont la quatrième est la plus révélatrice : **exposer les faits négatifs au même titre que les capacités**. Cinq le sont — la dépréciation ci-dessus ; une position d'analyste **non vérifiée**, sur laquelle il faut se garder de substituer un quadrant à un autre ; l'absence de tout lien documenté entre le portefeuille et les textes prudentiels canadiens ; l'absence d'architecture agentique de l'éditeur propre aux services financiers, **fait négatif vérifié** ; et le niveau de repérage, non élevé, d'un lien entre la couche de messagerie et la norme de message financier.

Une position d'analyste, quand elle est citée, l'est avec son commanditaire : une étude d'impact économique portant un taux de rendement déclaré a été **commandée par l'éditeur**, et ce fait accompagne chaque citation.

### 10.3 La correspondance réglementaire : cinq inférences sur sept liens

Table: Tableau 16 — Correspondance blueprint ↔ réglementation : sept liens, quatre statuts.

| Exigence | Réponse d'architecture | Statut du lien |
| --- | --- | --- |
| Ligne directrice E-23 (1er mai 2027) | Couches de gouvernance et d'observabilité : inventaire, cotation graduée, cycle de vie, documentation | **Inférence d'auteur — fait négatif établi** |
| Ligne directrice sur l'IA (AMF) | Couche de gouvernance et paliers de risque | Inférence d'auteur |
| Article 12.1, Loi 25 | Couches agentique et de gouvernance : humain-dans-la-boucle outillé, traces | Inférence d'auteur |
| Avis 11-348 | Couches d'exposition et d'observabilité : passerelles d'audit, traçabilité par tâche | Inférence d'auteur |
| Cadre bancaire | Couche d'exposition, pour le standard **à venir** | **Attente réglementaire — ne rien anticiper** |
| Rails de paiement (ISO 20022, rail temps réel) | Couche de messagerie et continuité opérationnelle | **Documenté** pour le rôle de l'intégrateur ; solution de transformation au niveau du repérage |
| Ligne directrice B-13 | Couche d'exécution : souveraineté, résidence | **Inférence d'auteur — fait négatif établi** |

Le résultat est assumé et il est le point de la section : **cinq inférences d'auteur, un lien documenté, une attente réglementaire. Un blueprint honnête produit peu de liens documentés.**

La distinction entre les deux statuts d'inférence est verrouillée, et elle mérite explication car la faute inverse est tentante. Pour deux textes seulement — la ligne directrice prudentielle sur le risque de modélisation et celle sur les risques technologiques —, la formule impose de dire à la fois que le rapprochement est une inférence d'auteur **et** que l'éditeur ne revendique aucune conformité. Pour les trois autres textes canadiens, on écrit l'inférence **et rien de plus** : y adjoindre l'absence de revendication serait faux, puisque le corpus atteste des revendications de l'éditeur — mais pour des instruments européens et internationaux, jamais canadiens. **Étendre un fait négatif au-delà de son périmètre est une faute symétrique de celle qu'on cherche à éviter.**

Un point critique, remonté et non résorbé, doit être signalé au lecteur : la réponse d'architecture « humain-dans-la-boucle outillé » associée à l'article 12.1 **contredit l'acquis de la section 5.3**. Outiller un point d'arrêt dans l'exécution n'outille pas la révision *a posteriori* due à la personne concernée. Le blueprint le porte ; l'ouvrage le signale.

### 10.4 Trois flux de bout en bout

**Flux 1 — la décision de crédit : le processus commande.** Un flux de travail déterministe appelle des agents de synthèse documentaire ; il ne leur délègue ni la séquence, ni la décision de s'arrêter, ni la production du journal. Trois points de passage obligés : chaque appel de modèle transite par la passerelle ; un mémo est produit pour un souscripteur humain — humain-dans-la-boucle au sens strict, qui n'outille pas la révision de l'article 12.1 ; et la trace est produite par la plateforme, non par les agents. Ce que le flux démontre n'est pas la conformité d'un dossier de crédit, mais qu'**un dossier peut être instruit par des agents sans que quiconque perde la maîtrise de l'enchaînement ni la capacité de dire ce qui s'est produit**. Les trois obligations de l'article 12.1 restent des obligations de l'institution.

**Flux 2 — le paiement normalisé : l'agent observe, le rail exécute.** Sur le rail de grande valeur, la place de l'agent est délibérément étroite : **il observe et il alerte ; il n'émet pas**. Le positionnement d'orchestration s'en trouve scindé — rail d'émission déterministe sans agent en OO3, agent de surveillance branché sur un flux d'événements sans cadre de processus explicite en OO1, car **un bus ajoute le journal, non le cadre**. Ce que le flux démontre tient en une phrase : le confort de l'architecture ne tient pas à ce que l'agent soit encadré, **mais à ce qu'il soit sans effet sur une dorsale qu'il ne commande pas**.

> **État de la connaissance vérifiable.** Le lien entre la couche de messagerie de ce portefeuille et la norme de message financier n'est documenté qu'indirectement, par un ouvrage technique de 2016. Une tentative d'élévation a été conduite et a **échoué** — accès refusé aux pages de documentation ; une source actuelle localisée ne nomme pas les composants attendus. L'entrée reste au niveau du repérage et ne peut porter aucun fait central : **ce flux illustre une architecture de principe**. La valeur de ce constat est méthodologique : un échec de recherche s'écrit, avec sa date, ses moyens et son motif.

**Flux 3 — l'accès au cadre bancaire : concevoir contre une norme qui n'existe pas encore.** Trois propositions, dont la dernière est négative. La couche d'exposition existe et est documentée — gestion du cycle de vie, portail développeur, famille de passerelles. Elle sait republier des actifs existants sous forme d'outils et de serveurs à partir des définitions d'interfaces — c'est le troisième principe, le contrat d'abord. Et **le consentement de la personne ainsi que le registre public des participants sont hors du périmètre du portefeuille** : ce sont des objets d'institution et de régulateur. Ce que le flux démontre est une discipline : **concevoir la couche, pas le contrat** — et ne pas bâtir contre une spécification d'industrie au motif qu'elle serait probable.

### 10.5 Lacunes propres et péremption du blueprint

Six lacunes étaient énumérées ; une est résolue — la clôture de l'acquisition, établie par une vérification ciblée sur source primaire, et sortie de la liste de veille. Cinq subsistent : une position d'analyste non vérifiée ; l'absence d'architecture agentique de l'éditeur propre aux services financiers ; une chaîne de transformation de messages restée au niveau du repérage après une élévation tentée et échouée ; le fait que cinq correspondances sur sept soient des inférences — lacune **structurelle, non comblable par une recherche** ; et la datation du blueprint lui-même.

Cinq événements le périment, et un seul n'a pas de date : l'arrêté ministériel de désignation, qui rendrait le troisième flux spécifiable — ou faux. Les quatre autres sont datés : la révision de la spécification protocolaire, l'entrée en vigueur du règlement administratif du rail, la clôture des commentaires réglementaires, et l'entrée en vigueur commune du 1er mai 2027. Cette dernière est une péremption **de discours** : le futur proche bascule au présent, et la question passe de « comment se préparer » à « comment se conformer ».

L'ouvrage referme cette partie par la formule qui en résume l'ambition : **un blueprint honnête n'est pas un blueprint sans lacunes ; c'est un blueprint dont les lacunes sont écrites, datées et surveillées.**

## 11. Discussion

### 11.1 Résultat central : l'encadrement comme objet probatoire

Le résultat de cette revue n'est pas que l'encadrement déterministe soit techniquement supérieur — aucune source ne l'établit pour le Canada, et le seul verdict empirique disponible est européen, hors finance, issu d'un préprint dont les auteurs déclarent leurs propres limites. Le résultat est autre, et il est plus solide parce qu'il repose sur un texte plutôt que sur une convergence : **l'encadrement est le seul objet dont l'exploitant d'un système agentique puisse démontrer la teneur devant un tiers.**

L'argument tient en une élimination. L'article 12.1 fait peser l'obligation d'expliquer sur l'entreprise qui rend la décision, sans lui ouvrir la ressource de désigner un tiers. Or l'explication ne peut venir ni de l'agent — la journalisation ne lui est pas recommandée —, ni du modèle — la causalité entre entrées et sorties est souvent indéterminable —, ni du protocole — il établit qui parle, non ce qui est dit ni si cela est fondé. Reste le cadre : écrit avant l'exécution, indépendant de la coopération des agents, et par construction lisible.

Ce résultat déplace la question d'architecture. Elle n'est plus « quelle pile technique adopter », mais **« quelle propriété serai-je en mesure de démontrer, et par quel objet ? »**. C'est pourquoi les cinq propriétés d'évaluation de la section 4.2 se répartissent comme elles le font : quatre relèvent de ce qu'on démontre, et la cinquième — l'autonomie — est précisément celle qui n'a aucune métrique. **Ce qui se gouverne est ce qui se démontre.**

### 11.2 Une lecture unifiée de deux calendriers disjoints

L'apport propre de l'ouvrage n'est pas de recenser un protocole ou un texte de plus : c'est de tenir les deux calendriers ensemble et de dire exactement où ils ne se touchent pas.

Ce raccordement produit trois observations qu'aucune des deux littératures ne pouvait produire seule.

**La densité protocolaire et la densité réglementaire sont inversement corrélées.** Le protocole le plus proche de l'acte réglementé — le paiement — est celui dont le corpus documente le moins : aucune anatomie technique, aucun transfert de gouvernance, aucune articulation avec les rails du pays. Symétriquement, les deux textes les mieux documentés sont ceux dont aucun croisement protocolaire ne porte de lien.

**Les deux trajectoires ne progressent pas au même rythme et ce décalage est structurel.** Une spécification protocolaire se révise en huit mois ; une ligne directrice se donne dix-neuf mois de préavis. Un chapitre technique gelé douze jours avant une révision majeure est daté avant d'être publié ; un chapitre sur une disposition législative de 2023 ne l'est pas. Toute architecture qui absorberait le rythme protocolaire dans ses points de contrôle réglementaires hériterait de la volatilité du premier.

**La qualification précède l'architecture.** Le cas le plus instructif de la section 8 n'est pas le plus spectaculaire : c'est l'acheminement autonome de courriels, parce qu'il rend indécidable la question de savoir si l'on est en présence d'une décision fondée exclusivement sur un traitement automatisé. Aucune décision d'architecture n'est possible avant cette qualification, et cette qualification n'est pas technique.

### 11.3 Tensions transversales

Quatre tensions traversent l'ouvrage et ne sont résolues par aucune de ses sections.

**Encadrer coûte.** L'effort initial et la maintenance figurent explicitement parmi les critères de sélection d'une option d'orchestration. Un principe d'architecture qui ignorerait son coût de possession serait un slogan.

**Expliquer et confiner tirent en sens contraire.** Le paradoxe de confidentialité de l'explicabilité — l'explicabilité expose, la confidentialité restreint — oppose directement deux des cinq points de contrôle obligatoires. La littérature nomme la tension ; elle ne la résout pas. L'arbitrage se documente, il ne se déclare pas résolu.

**Une frontière commerciale peut se déguiser en frontière technique.** Le support d'un protocole par la plateforme commerciale d'un produit et son absence dans la bibliothèque libre du même produit ne se distinguent pas dans une architecture de principe. Ils se distinguent dans une facture.

**Le fait négatif est plus utile que l'ignorance, et plus dangereux.** Un fait négatif vérifié autorise une décision d'architecture — traiter un standard non désigné comme une variable derrière une abstraction. Une absence de documentation n'autorise rien. Confondre les deux dans un dossier de conformité expose à une découverte par un tiers ; les distinguer exige de documenter ce qu'on a cherché, et pas seulement ce qu'on a trouvé.

### 11.4 Limites : onze lacunes exposées

L'ouvrage recense onze lacunes ouvertes, aucune comblée, toutes assignées à un chapitre porteur.

Table: Tableau 17 — Les onze lacunes ouvertes, par famille.

| Famille | Lacune | Ce qui la comblerait |
| --- | --- | --- |
| **Acte réglementaire attendu** | Désignation de l'organisme de normalisation technique du cadre bancaire | Un arrêté ministériel |
| | Réglementaire fin : positions de l'autorité de contrôle québécoise sur l'art. 12.1 appliqué au multi-agents ; suites de la consultation en valeurs mobilières ; **contenu article par article de la ligne directrice sur l'IA de l'AMF** | Une décision, un avis subséquent — et, pour la troisième, **une extraction primaire non faite** |
| **Objets techniques** | Mécanique des risques protocolaires ; aucune attaque propre à A2A documentée | Source primaire dédiée, identifiant de vulnérabilité, incident public daté |
| | Anatomie non documentée : ancrage de confiance, révocation et gouvernance des clés des cartes signées ; date de la v1.0 ; multi-location ; inventaire infonuagique de MCP ; structure de message, mandat et preuve d'intention du protocole de transaction | Les spécifications elles-mêmes |
| | Sous-caractérisation du corpus académique : cadre **opérationnel** non caractérisé ; actionnabilité conversationnelle sans définition ni terme d'origine ; autonomie sans métrique ; taxonomie sans reproduction indépendante | Une relecture ciblée des articles, dont les documents sont au dépôt |
| **Adoption** | Une institution sans aucune source primaire ; deux élévations partiellement échouées (facteur de risque matériel, agents en production) | Sources primaires accessibles |
| | Chiffres d'adoption d'entreprise d'un cadriciel non vérifiés | Une passe de recherche non tentée |
| **Blueprint** | Position d'analyste non vérifiée ; chaîne de transformation de messages au niveau du repérage après échec ; date d'annonce d'un composant ; aucune annonce de l'éditeur avec une institution canadienne autre que deux | Documentation d'éditeur, canal client |
| **Prospective et divers** | Articulation du protocole de transaction avec les rails canadiens | Voir les cinq questions de la section 7.3 |
| | Identité de la composante homonyme au sein d'AGNTCY | Documentation du projet ; **aucune passe de recherche conduite** |
| | Une attribution sans date : la supervision du cadre bancaire est rattachée à un document budgétaire dont le corpus ne date pas l'annonce | Le document lui-même |

La plus coûteuse est la troisième de la deuxième ligne, et elle mérite d'être nommée sans atténuation : **le contenu de la ligne directrice sur l'IA de l'AMF n'est pas au corpus**, alors qu'il s'agit de l'un des deux textes qui fixent l'échéance du 1er mai 2027. Le chapitre pivot n'en dérive aucune contrainte ; la matrice porte, à cause d'elle, sa seule ligne entièrement vide. Cette lacune est d'une espèce distincte des autres : le texte est **final depuis le 30 mars 2026**, rien n'y est en attente. La lacune tient à une extraction primaire qui n'a pas été faite — et une revalidation ultérieure a constaté que le site de l'autorité refuse le chargement automatisé. *Une lacune qui attend une autorité et une lacune qui attend un lecteur ne se soldent pas de la même façon ; celle-ci est à la charge de l'ouvrage.*

La onzième lacune mérite également une mention, pour ce qu'elle enseigne sur la méthode plutôt que sur son objet. Elle a été découverte après le gel des vingt-quatre chapitres, à la construction d'une frise chronologique — parce qu'**une frise ne peut pas porter un événement sans date**, là où la prose absorbe l'imprécision sans broncher. Aucun fait central n'en dépend. L'enseignement, lui, est transposable : **un format de restitution nouveau — frise, matrice, tableau — est un révélateur de ce que la prose absorbe en silence.**

À ces onze lacunes s'ajoutent trois limites de méthode que l'ouvrage déclare.

La **charpente conceptuelle de trois parties sur sept repose sur une source unique** : un préprint version 1, non révisé par les pairs, sans reproduction indépendante ni application documentée à un processus d'institution financière.

Le **vote adversarial n'a couvert qu'une fraction du corpus** : soixante-quinze affirmations sur trois cent quatre-vingt-quatre, un peu moins d'un cinquième, par plafond budgétaire déclaré.

Et **« confirmé » ne signifie pas « vrai »** : la méthode vérifie qu'une source *dit* X à une date donnée. Trois erreurs sont documentées — une hallucination écartée en cours de passe, deux convergences fabriquées ; rien ne dit que ce soit tout.

### 11.5 Fronts de recherche ouverts

Six questions sont formulées avec, pour chacune, ce qui la trancherait — la formulation étant elle-même un exercice de discipline, puisqu'une question dont on ne sait pas ce qui la trancherait n'est pas une question de recherche.

La taxonomie de l'encadrement **résiste-t-elle à une reproduction indépendante ?** Une réplication par une équipe distincte, et une application documentée à un processus d'institution financière, le diraient.

Quelle est la **mécanique des risques protocolaires**, et existe-t-il des attaques propres au protocole inter-agents ? Une source primaire consacrée, un identifiant de vulnérabilité ou un incident public daté le diraient.

**Que vaut cryptographiquement une carte d'agent signée ?** La documentation de son ancrage de confiance, de sa révocation et de sa gouvernance des clés le dirait — ces trois questions décidant, en pratique, de la valeur d'une signature.

**L'article 12.1 s'applique-t-il à une décision rendue par un système multi-agents comportant un humain-dans-la-boucle ?** Une position de l'autorité de contrôle ou une décision judiciaire le dirait ; le corpus ne porte qu'une analyse de cabinet au niveau du repérage.

**Quelle organisation sera désignée organisme de normalisation technique ?** Un arrêté ministériel le dira.

**Comment le cadre opérationnel s'articule-t-il au cadre normatif ?** Une relecture ciblée du manifeste — dont le document est au dépôt — le dirait, et cette élévation est possible en une passe. C'est la question la plus embarrassante de la liste : l'articulation des deux types de cadres est la thèse de l'ouvrage, et l'un des deux termes n'est pas caractérisé par la source qui les distingue.

### 11.6 Synthèse opérationnelle : le principe directeur, strate par strate

Table: Tableau 18 — Réalisation du principe directeur et verrou dominant, par strate.

| Strate | Ce que le cadre doit y porter | Verrou dominant |
| --- | --- | --- |
| **Protocoles** | Rien : le protocole établit qui parle et par quel format | Aucune spécification ne porte de décision ni de trace d'instance |
| **Identité et registre** | Un cadre opérationnel écrit hors de l'agent : outils atteignables, bornes de privilège | Moitié supérieure de l'édifice en **brouillon expiré** ; aucune norme ratifiée de registre d'agents |
| **Orchestration** | L'enchaînement, la trace, le point d'arrêt humain | Taxonomie à **source unique** ; le produit n'arbitre pas, la configuration arbitre |
| **Gouvernance** | Inventaire, cotation, cycle de vie, documentation, surveillance | Spécifiable contre un texte, pas contre l'autre — les deux entrant en vigueur le même jour |
| **Rails financiers** | L'agent observe, le rail exécute | Un rail acquis, un rail **visé** : deux régimes de preuve à ne pas additionner |
| **Cadre bancaire** | Une frontière d'abstraction devant un standard non désigné | **Fait négatif vérifié** — le seul qui autorise une décision d'architecture |
| **Blueprint** | Six principes, dont un seul appuyé sur une convergence | Cinq correspondances sur sept sont des **inférences d'auteur** |

### 11.7 Ce qui périme l'ouvrage, et à quelle date

Table: Tableau 19 — Les échéances qui périment l'ouvrage, et ce qu'elles atteignent.

| Échéance | Événement | Ce qu'il périme |
| --- | --- | --- |
| 28 juill. 2026 | Finalisation de la révision majeure de la spécification MCP | L'anatomie protocolaire de la section 3.2 |
| 24 août 2026 | Entrée en vigueur du règlement administratif du rail temps réel | Le statut réglementaire de la section 7.2 |
| 26 août 2026 | Clôture des commentaires sur le règlement du cadre bancaire | La description du dispositif de la section 7.1 ; le texte final peut changer |
| **Sans date** | Arrêté ministériel désignant l'organisme de normalisation technique | Le fait négatif vérifié de la section 7.1, son garde-fou, et le troisième flux de la section 10.4 |
| Cible T4 2026 | Lancement effectif du rail temps réel | La réserve interdisant d'écrire « lancé » ou « en production » |
| Continue | Trajectoire du projet de loi C-36 | La caractérisation de la section 5.2 |
| 1er mai 2027 | Entrée en vigueur commune des deux lignes directrices | Le discours de « futur proche » des sections 5, 6 et 9, et la feuille de route |

Une nuance close ce tableau, et elle n'est pas de style : **périmer n'est pas falsifier.** Chaque échéance vise une formulation précise, bornée par sa date. L'événement la rend **non actuelle**, non fausse. C'est exactement ce que la date de gel achète : le droit de dire quelque chose de précis à une date, plutôt que quelque chose de vague indéfiniment.

## 12. Conclusion

Entre la fin de 2024 et le milieu de 2026, une couche protocolaire agentique s'est constituée et, pour l'essentiel, ouverte ; des institutions financières canadiennes ont mis en service des systèmes qui décident ; et un maillage réglementaire s'est réorganisé autour d'une date unique, le 1er mai 2027. Ces trois mouvements se sont produits sans se rencontrer : croiser trois protocoles majeurs avec cinq textes canadiens produit quinze croisements dont **aucun n'est documenté**.

Ce vide n'autorise ni l'attente, ni l'improvisation. Il déplace la charge. Puisqu'aucune spécification ne dit comment satisfaire un texte canadien, et puisqu'aucun texte canadien ne dit comment construire un système agentique, ce qui reste à l'institution est ce qu'elle écrit elle-même : **le cadre qui impose l'enchaînement, produit la trace, arrête pour l'humain, sépare l'adaptation de l'évolution et confine chaque agent à ce qu'il doit atteindre.** C'est le sens du principe directeur — sous exigence réglementaire stricte, le processus est imposé de façon déterministe par le cadre, qui invoque les agents ; les agents n'orchestrent pas le processus — et c'est le sens du titre : l'autonomie n'est pas supprimée, elle est encadrée.

Ce principe est une inférence d'auteur, et il est énoncé avec ses quatre conditions : la transposition au Canada n'est établie par aucune source ; l'« exigence réglementaire stricte » n'est pas définie ; l'encadrement se paie ; et le principe dit où placer la main qui commande, non ce que le cadre doit contenir. Sa partie la plus solide n'est ni la convergence de trois sources qui se recoupent deux à deux, ni la constante lue dans trois textes : c'est ce qu'un seul texte, l'article 12.1, fait peser sur l'entreprise qui rend la décision — une obligation d'expliquer, sans faculté de désigner un tiers.

Reste ce que l'ouvrage ne sait pas, et qu'il expose plutôt qu'il ne le comble : onze lacunes, dont le contenu de l'une des deux lignes directrices qui fixent l'échéance ; une charpente conceptuelle à source unique ; un vote adversarial qui n'a couvert qu'un cinquième des affirmations extraites ; et deux fautes de méthode commises, détectées et racontées. Cette exposition n'est pas une clause de style. Elle est la contrepartie d'une discipline dont la valeur se mesure à ce qu'elle refuse d'écrire — et dont l'ouvrage tire lui-même le bilan le plus utile qu'il puisse offrir à un lecteur : **tous les défauts lourds trouvés à la publication l'ont été par des relecteurs adversariaux, aucun par l'auto-contrôle de leur rédacteur.**

Ce que cette méthode achète n'est pas la certitude. C'est **la traçabilité de l'incertitude** — et, pour une institution qui devra démontrer sa maîtrise devant un tiers au printemps 2027, c'est très exactement ce qui est demandé.

## Annexe A — Sigles et acronymes

**A2A** — *Agent2Agent*, protocole de collaboration entre agents.
**ACFC** — Agence de la consommation en matière financière du Canada.
**ACP** — sigle **proscrit employé seul** : voir tableau 4 (quatre objets distincts — protocole d'IBM Research/BeeAI fusionné dans A2A ; *Agentic Control Plane* d'un consortium canadien ; positionnement d'éditeur ; composante d'AGNTCY).
**ACVM** — Autorités canadiennes en valeurs mobilières.
**ADK** — *Agent Development Kit*, trousse de développement d'agents.
**AGNTCY** — projet d'infrastructure agentique (annuaires de découverte, transport dédié).
**AMF** — Autorité des marchés financiers (Québec).
**AP2** — *Agent Payments Protocol*, protocole de transaction pilotée par agents.
**APM** — *Agentic Business Process Management*, gestion agentique des processus d'affaires.
**BPMN / BPEL** — langages de modélisation et d'exécution de processus d'affaires.
**BSIF** — Bureau du surintendant des institutions financières (Canada).
**CAI** — Commission d'accès à l'information (Québec).
**CBPR+** — profil de messagerie transfrontalière ISO 20022.
**CSA** — *Cloud Security Alliance* (à ne pas confondre avec les ACVM, dont le sigle anglais est identique).
**E-23** — ligne directrice du BSIF sur la gestion du risque de modélisation (en vigueur le 1er mai 2027).
**F1** — moyenne harmonique de la précision et du rappel.
**FDX** — *Financial Data Exchange* ; organisme d'industrie, **non désigné** comme organisme de normalisation technique du cadre bancaire.
**FINOS** — *Fintech Open Source Foundation*.
**FNR** — taux de faux négatifs (*false negative rate*).
**HITL** — humain-dans-la-boucle (*human-in-the-loop*).
**iPaaS** — plateforme d'intégration en tant que service.
**ISO 20022** — norme internationale de messages financiers.
**JSON-RPC** — protocole d'appel de procédure à distance sur JSON.
**LIAD** — Loi sur l'intelligence artificielle et les données (volet de C-27, mort au feuilleton le 6 janvier 2025).
**LPRPDE** — Loi sur la protection des renseignements personnels et les documents électroniques.
**Loi 25** — loi québécoise modernisant les dispositions législatives en matière de protection des renseignements personnels ; son article 110 introduit l'article 12.1 de la loi sur le secteur privé.
**Lynx** — système canadien de paiement de grande valeur.
**MCP** — *Model Context Protocol*, protocole d'accès des agents aux outils et aux données.
**MT / MX** — anciens et nouveaux formats de messages financiers (MX = ISO 20022).
**OAuth / OIDC** — cadre d'autorisation déléguée ; couche d'authentification associée.
**OO1–OO4** — les quatre options d'orchestration (voir tableau 5).
**pacs.008 / pacs.009** — types de messages ISO 20022 de virement client et interbancaire.
**RESL** — prêt garanti par l'immobilier résidentiel (*residential secured lending*).
**RTR** — *Real-Time Rail*, rail de paiement en temps réel canadien.
**SCIM** — *System for Cross-domain Identity Management*, norme de provisionnement d'identités.
**SPIFFE / SPIRE** — cadre et implémentation d'identité de charge de travail.
**TSC** — comité de pilotage technique (*technical steering committee*).

## Note sur la méthode et les sources

Cet article est une revue de synthèse tirée d'un ouvrage de vingt-neuf pièces et d'environ quatre-vingt-douze mille mots, dont le socle documentaire est arrêté au **17 juillet 2026** et repose sur quarante-six entrées vérifiées, portant chacune son niveau de preuve.

Les principes de sourçage exposés à la section 2 s'appliquent au présent article sans exception. Toute métrique d'adoption de protocole, d'institution financière ou d'éditeur de plateforme est **auto-déclarée** et rapportée comme telle. La couverture de l'intelligence artificielle agentique par la ligne directrice E-23 est une **inférence d'analystes juridiques**, jamais une exigence du régulateur. Aucune position de conformité n'est émise : les correspondances entre composants et textes réglementaires sont, sauf mention contraire, des **inférences d'auteur**. Trois faits négatifs seulement sont **vérifiés** au sens strict — absence de vocabulaire agentique dans le texte de la ligne directrice E-23, absence de standard technique désigné dans les textes du cadre bancaire, absence d'architecture agentique sectorielle publiée par l'éditeur instancié à la section 10 ; partout ailleurs, « le corpus ne documente pas » signifie une absence de lecture, non une absence de fait.

Les références ci-dessous sont regroupées par thème. Le marqueur **⚠** signale une **ressource vivante** — spécification versionnée par date, texte réglementaire en cours d'adoption, page produit susceptible d'évoluer — dont la version doit être fixée au moment de citer. Les accents et graphies des titres originaux sont préservés tels quels.

## Références

### A. Protocoles agentiques et gouvernance

Anthropic (2024). *Model Context Protocol* — annonce initiale, novembre 2024. Annonce d'éditeur.

Agentic AI Foundation / Linux Foundation (2025). *Transfert de la gouvernance du Model Context Protocol*, décembre 2025. Communiqué de fondation.

Model Context Protocol (2025). *Specification, revision 2025-11-25*. Spécification technique. modelcontextprotocol.io. ⚠

Model Context Protocol (2026). *Release candidate verrouillée le 21 mai 2026 ; finalisation annoncée pour le 28 juillet 2026 — cœur sans état, retrait de l'en-tête `Mcp-Session-Id`, introduction de `Mcp-Method` et `Mcp-Name`*. Spécification en cours. ⚠

Google (2025). *Lancement du protocole Agent2Agent*, avril 2025. Annonce d'éditeur.

Linux Foundation (2025). *Transfert d'Agent2Agent à la Linux Foundation sous licence Apache 2.0*, juin 2025. Communiqué.

A2A Project (2026). *Announcing A2A v1.0* — « Complementary to MCP, not a replacement ». Documentation de projet. a2a-protocol.org/latest/announcing-1.0. ⚠

A2A Project (2026). *A2A v1.0.1*, 28 mai 2026 ; gouvernance et publications du dépôt. github.com/a2aproject/A2A. ⚠

Linux Foundation (2026). *Bilan de première année d'Agent2Agent*, 9 avril 2026 — décomptes d'organisations de soutien pour A2A et pour AP2 (métriques auto-déclarées). Communiqué.

Google Cloud (2025). *Agent Payments Protocol (AP2)*, annonce du 16 septembre 2025. Annonce d'éditeur.

Cisco, LangChain et Galileo (2025). *Ouverture d'AGNTCY*, mars 2025. Annonce.

Linux Foundation (2025). *AGNTCY placé sous la gouvernance de la Linux Foundation*, 29 juillet 2025 — annuaires de découverte, transport dédié, membres formateurs. Communiqué.

Martineau, K. / IBM Research (2025). *BeeAI et l'Agent Communication Protocol*, billet du 28 mai 2025. research.ibm.com/blog/agent-communication-protocol-ai.

LF AI & Data (2025). *Fusion du protocole ACP d'IBM Research dans A2A*, 29 août 2025. Billet de fondation.

Microsoft, Amazon Web Services et Google Cloud (2026). *Intégrations infonuagiques d'A2A — préversion, disponibilité générale d'avril 2026, exécution dédiée*. Documentation d'éditeurs. ⚠

Ehtesham, A. et coll. (2025). *A survey of agent interoperability protocols*. Préprint non révisé par les pairs. arXiv:2505.02279. — *Jalon historiographique : sa séquence d'adoption a été périmée moins de quatre mois après publication.*

### B. Orchestration, cadriciels, identité et registres d'agents

Microsoft (2026). *Microsoft Agent Framework 1.0* — disponibilité générale le 3 avril 2026 ; successeur de deux lignées antérieures ; support MCP natif ; flux de travail à base de graphes, points de contrôle, humain-dans-la-boucle. Microsoft Learn ; DevBlogs. ⚠

Microsoft (2026). *Limites connues du magasin de points de contrôle en déploiement distribué multi-instances*. Discussion publique du dépôt (n° 2305). ⚠

LangChain (2025). *LangGraph Platform — disponibilité générale*, billet du 14 mai 2025 (métrique d'adoption auto-déclarée). langchain.com/blog.

LangChain (2026). *Serveur A2A de LangSmith Deployment* — point d'entrée et méthodes documentées. docs.langchain.com. ⚠

langchain-ai/langgraph (2026). *Demande de support A2A dans la bibliothèque libre*, ouverte le 3 avril 2026. Dépôt public. ⚠

Confluent (2025). *Streaming Agents*, août 2025 — agents événementiels, appel d'outils MCP natif déclaré. Billet d'éditeur.

Confluent (2026). *Mise à jour trimestrielle du 26 février 2026* — intégration A2A en préversion publique ; serveur MCP officiel pour l'offre infonuagique. Billet d'éditeur. ⚠

CrewAI (2026). *A2A Agent Delegation* — « CrewAI treats A2A protocol as a first-class delegation primitive » ; journal des modifications depuis novembre 2025. docs.crewai.com. ⚠

Temporal (2026). *Intégration d'une trousse d'agents en préversion publique ; recettes MCP*. docs.temporal.io. — *Repérage documentaire à confirmer.* ⚠

Microsoft (2026). *Entra Agent ID* — identités d'agents et *blueprints* ; fondations OAuth 2.0 et OpenID Connect ; scénarios d'action propre et de délégation. Microsoft Learn. ⚠

Cloud Security Alliance Labs (2026). *Agent Registry specification*, publiée le 27 mars 2026 — **brouillon** ; quatre fonctions du registre faisant autorité ; champs d'outils autorisés et de bornes de privilège. labs.cloudsecurityalliance.org. ⚠

IETF (2015). *System for Cross-domain Identity Management: Core Schema*. RFC 7643.

Abbey, R. / Okta (2026). *SCIM Agent Extension*, brouillon individuel `draft-abbey-scim-agent-extension-00`, **expiré le 19 avril 2026** ; consolidation en cours. datatracker.ietf.org. ⚠

### C. Littérature académique

Calvanese, D., De Giacomo, G., Dumas, M., Kampik, T., Montali, M., Rinderle-Ma, S., Weber, I. et coll. (2026). *Agentic Business Process Management: A Research Manifesto*. Dix-huit auteurs ; issu du séminaire Dagstuhl no 25192. *Information Systems*, 140, 102738 ; préprint arXiv:2603.18916v3 (12 avril 2026). — *Position argumentée de ses auteurs ; recommande, ne juge pas.*

Rinderle-Ma, S., Mangler, J. et coll. (2026). *Design and Implementation of Agentic Orchestrations and Orchestration of Agents*. Technische Universität München. Préprint version 1, **non révisé par les pairs**, 30 juin 2026. arXiv:2606.31518v1. — *Menaces à la validité déclarées par les auteurs ; résultats chiffrés repris à titre d'illustration seulement.*

### D. Réglementation canadienne

Bureau du surintendant des institutions financières (2025). *Ligne directrice E-23 — Gestion du risque de modélisation (2027)*, version finale du 11 septembre 2025, en vigueur le 1er mai 2027, et lettre d'accompagnement. Versions anglaise et française. osfi-bsif.gc.ca.

Bureau du surintendant des institutions financières et Agence de la consommation en matière financière du Canada (2024). *Rapport conjoint sur l'adoption de l'intelligence artificielle par les institutions financières fédérales*, 24 septembre 2024. Enquête auto-déclarée ; la trajectoire projetée est une projection de ses auteurs. osfi-bsif.gc.ca.

Parlement du Canada (2025). *Prorogation du 6 janvier 2025* — mort au feuilleton du projet de loi C-27, dont la Loi sur l'intelligence artificielle et les données.

Parlement du Canada (2026). *Projet de loi C-36 — Protecting Privacy and Consumer Data Act*, première lecture le 15 juin 2026 ; deuxième lecture à la date d'observation. LEGISinfo, 45-1/c-36. ⚠ — *Loi sur la vie privée portant des volets d'intelligence artificielle, non loi sur l'IA autonome ; adoption incertaine.*

Autorité des marchés financiers (2026). *Ligne directrice sur l'intelligence artificielle* — projet publié le 3 juillet 2025, consultation close le 7 novembre 2025, **version finale publiée le 30 mars 2026**, en vigueur le 1er mai 2027. lautorite.qc.ca. ⚠ — *Le contenu de la version finale n'est pas au corpus (lacune la plus coûteuse de l'ouvrage).*

Autorités canadiennes en valeurs mobilières (2024). *CSA Staff Notice and Consultation 11-348 — Applicability of Canadian Securities Laws and the use of Artificial Intelligence Systems in Capital Markets*, 5 décembre 2024 ; consultation close le 31 mars 2025. (2024) 47 OSCB 9307 ; osc.ca ; securities-administrators.ca.

Éditeur officiel du Québec (2023). *Loi sur la protection des renseignements personnels dans le secteur privé*, RLRQ c. P-39.1, **article 12.1**, en vigueur depuis le 22 septembre 2023 (introduit par 2021, c. 25, a. 110). legisquebec.gouv.qc.ca.

Commission des valeurs mobilières de l'Ontario et EY. *AI in Capital Markets: Exploring Use Cases in Ontario*. Rapport cité par l'avis 11-348.

Autorité des marchés financiers. *Document de réflexion — meilleures pratiques d'utilisation responsable de l'intelligence artificielle dans le secteur financier*. Cité par l'avis 11-348 ; antérieur au 5 décembre 2024, date exacte non établie par le corpus.

### E. Interopérabilité financière canadienne

Parlement du Canada (2026). *Projet de loi C-15 — cadre des services bancaires axés sur le consommateur*, déposé le 18 novembre 2025, **sanction royale le 26 mars 2026** ; abroge et remplace la loi partielle de 2024 ; crée un droit à la mobilité des données à l'échelle de l'économie. parl.ca.

Gouvernement du Canada (2026). *Règlement sur les services bancaires axés sur le consommateur* — **prépublication**, Gazette du Canada, partie I, vol. 160, no 26, 27 juin 2026 ; commentaires de soixante jours clos le 26 août 2026 ; entrée en vigueur échelonnée. gazette.gc.ca. ⚠

Ministère des Finances du Canada (2026). *Communiqué du 26 juin 2026* et pages budgétaires afférentes. canada.ca. ⚠

Banque du Canada (2026). *Supervision et mise en œuvre du cadre des services bancaires axés sur le consommateur ; registre public*. bankofcanada.ca. ⚠ — *Fait négatif vérifié associé : aucun organisme de normalisation technique désigné par arrêté ministériel, aucun standard nommé dans les textes officiels au 16 juillet 2026.*

Paiements Canada (2025). *Achèvement de la migration ISO 20022 sur Lynx* — fin de la coexistence MT/MX le 22 novembre 2025, annonce du 26 novembre 2025 ; alignement sur l'échéance mondiale de la messagerie transfrontalière ; proportion de messages au nouveau format déclarée par l'opérateur. payments.ca.

Paiements Canada (2026). *Real-Time Rail — état d'avancement du programme* : reprise le 16 avril 2024 avec ses partenaires de livraison et d'exploitation ; compensation et règlement complétés au troisième trimestre de 2025 ; tests d'intégration puis d'acceptation ; tests industriels en cours à la mi-2026 ; **cible de lancement au quatrième trimestre de 2026**, ISO 20022 dès l'origine. payments.ca. ⚠

Paiements Canada (2026). *Règlement administratif no 10* — Gazette du Canada, partie II, 1er juillet 2026 ; en vigueur le 24 août 2026. gazette.gc.ca. ⚠ — *Le contenu du règlement n'est pas au corpus.*

Mouvement Desjardins (2026). *Rapport annuel 2025* — rapport de gestion du 24 février 2026 ; participation au cadre bancaire conditionnelle à l'accréditation. desjardins.com.

### F. Institutions financières canadiennes (études de cas)

> *Toutes les métriques de cette section sont auto-déclarées par l'institution nommée et ne font l'objet d'aucune vérification indépendante.*

Groupe Banque TD (2026). *Premier modèle d'intelligence artificielle agentique de TD, développé par Layer 6* — pré-adjudication du prêt garanti par l'immobilier et mémos de synthèse ; communiqué du 21 mai 2026. stories.td.com.

Groupe Banque TD (2025). *Équipe Trustworthy AI — cinq critères* ; prix DataIQ « Best Responsible AI Program or Initiative in North America » 2025. stories.td.com ; dataiq.global.

Banque Scotia (2025). *Perspectives*, juillet 2025 — outil interne d'analyse documentaire, capacités agentiques ajoutées en 2025 ; traitement autonome des courriels de clients ; revue éthique préalable graduée. scotiabank.com. Corroboration : *MIT Sloan Management Review*.

Banque Royale du Canada (2026). *Création d'un groupe d'intelligence artificielle rattaché au chef de la direction*, communiqué du 18 février 2026 ; transcription officielle de la journée des investisseurs 2025. rbc.com/newsroom.

FINOS (2025-2026). *Travaux sur les contrôles communs des services d'intelligence artificielle ; présidence du conseil de gouvernance depuis août 2025 ; fonds sectoriel de juin 2026*. finos.org.

Manuvie et Akka (2026). *Communiqué conjoint du 10 mars 2026* — environnement d'exécution retenu pour la plateforme d'IA agentique d'entreprise, alors en version bêta ; partenariat antérieur d'ajustement fin (décembre 2025). Miroirs de fil de presse.

Evident AI Index (2026). *Classement de juin 2026* — méthodologie et périmètre non documentés par le corpus. ⚠

Banque CIBC (2025). *Déploiement à l'échelle de la banque de sa plateforme interne d'intelligence artificielle générative*, communiqué du 27 mai 2025 ; pilote depuis juillet 2024 ; formation obligatoire préalable à l'accès. newswire.ca. — *Gen-IA assistive : le communiqué ne mentionne ni agents autonomes, ni protocole d'invocation d'outils.*

Intact Corporation financière (2026). *Rapport annuel 2025*, publié le 26 mars 2026 — parc de modèles en production, bénéfices annuels attribués, part du portefeuille tarifée par apprentissage machine. intactfc.com. — *Fait établi par extraction intégrale : la terminologie agentique est absente de la lettre du chef de la direction.*

BMO (2025). *Lancement d'un assistant d'intelligence artificielle générative interne*, mai 2025 ; population d'utilisateurs déclarée dans une **soumission distincte** à un prix sectoriel (2025). newswire.ca ; qorusglobal.com. ⚠

Vector Institute (2026). *Étude de cas — formations internes sur l'intelligence artificielle agentique et multimodale*. vectorinstitute.ai. — *Attesté par le partenaire, non par la banque.*

Lightworks, Banque Scotia, Sun Life et TELUS (2026). *Lightworks, Scotiabank, Sun Life and TELUS Launch AI Consortium to Jointly Build Critical AI Control Infrastructure in Canada*, communiqué du 7 juillet 2026 — programme phare déclaré en production en environnement réglementé ; volume de jetons agrégé auto-déclaré. newswire.ca. ⚠ — *Homonymie sans lien avec le protocole d'IBM Research fusionné dans A2A (voir tableau 4).*

### G. Portefeuille d'instanciation et blueprint

> *Cas d'instanciation documenté par sources primaires, retenu pour la traçabilité de sa documentation publique. Ni recommandation d'achat, ni verdict comparatif.*

IBM (2023-2025). *Acquisition des plateformes webMethods et StreamSets* — annonce du 18 décembre 2023, clôture le 1er juillet 2024 ; **IBM webMethods Hybrid Integration**, disponibilité générale le 16 juin 2025. newsroom.ibm.com ; dépôts réglementaires.

IBM (2026). *Cloud Pak for Integration 16.2.0 LTS*, disponibilité générale le 30 juin 2026 ; *IBM MQ 10.0 LTS*, disponibilité générale en juin 2026 ; *App Connect Enterprise 13.0.7*, mars 2026 ; solution de messagerie ISO 20022. ibm.com ; community.ibm.com. ⚠

IBM (2025). *API Connect 12.1.0* (15 décembre 2025) ; *API Agent*, disponibilité générale le 19 novembre 2025 ; *DataPower Nano Gateway* (19 novembre 2025) ; *AI Gateway* (2024, étendue en 2025) ; *DataPower Interact Gateway*, **annoncée** en 2026. ibm.com/docs ; community.ibm.com. ⚠

IBM (2026). *ContextForge (`mcp-context-forge`)* — passerelle et registre unifiés devant serveurs d'outils, agents et interfaces. Dépôt public. github.com/IBM. ⚠

IBM (2025-2026). *Acquisition de Confluent* — annonce du 8 décembre 2025 ; expiration du délai réglementaire le 12 janvier 2026 ; approbation des actionnaires le 12 février 2026 ; **clôture le 17 mars 2026**. newsroom.ibm.com ; dépôts réglementaires.

IBM (2026). *Dépréciation d'Event Streams et d'Event Processing* — « Event Streams and Event Processing are deprecated and no enhancements will be provided » ; trajectoire reportée sur la filiale acquise ; Event Endpoint Management maintenu. ibm.github.io/event-automation. ⚠

IBM (2025-2026). *watsonx Orchestrate* — trousse de développement (12 mai 2025), outils d'invocation dès mai 2025, support inter-agents à compter du 30 juin 2025, serveurs distants (septembre 2025), garde-fous documentés ; positionnement d'éditeur comme *agentic control plane* depuis mai 2026 (attribué à IBM). github.com/IBM/ibm-watsonx-orchestrate-adk ; ibm.com. ⚠

IBM (2023-2025). *watsonx.governance* — disponibilité générale annoncée en décembre 2023 ; gouvernance agentique livrée en 2025 (studio d'évaluation, atlas des risques, catalogue gouverné). newsroom.ibm.com. ⚠

Credo AI et IBM (2025). *Accord d'origine externe du 28 avril 2025 — module d'accélérateurs de conformité* portant sur des instruments européens et internationaux, **aucun texte canadien**. credo.ai.

IBM (2025). *Agentic AI — pattern d'architecture*, mis à jour le 21 février 2025 : les flux de travail statiques (BPMN/BPEL) sont explicitement recommandés pour les processus sous surveillance réglementaire. ibm.com/think/architectures/patterns/agentic-ai. — *Pattern générique : aucune architecture agentique propre aux services financiers n'est publiée (fait négatif vérifié).*

IBM Canada (2019-2026). *Partenaire technologique principal de Lynx* (annonce du 2 mai 2019 ; lancement le 1er septembre 2021 ; centre de données additionnel en octobre 2023) ; *partenaire de livraison et d'exploitation du Real-Time Rail* depuis le 16 avril 2024. canada.newsroom.ibm.com ; payments.ca. — *Fait de contexte, jamais argument de conformité ; montants des contrats non publics.*

IBM (2025-2026). *IBM Cloud for Financial Services* — cadre de contrôles sectoriel fondé sur un référentiel public ; régions canadiennes (2021 et 3 avril 2025) ; *Sovereign Core*, disponibilité générale le 5 mai 2026. cloud.ibm.com/docs. ⚠ — *Aucun lien documenté avec les lignes directrices canadiennes ; toute correspondance est une inférence d'auteur.*

Forrester (2025). *Total Economic Impact of IBM webMethods* — **étude commandée par IBM** ; taux de rendement déclaré cité avec son commanditaire.

Forrester (2025). *Wave — plateformes d'intégration en tant que service, troisième trimestre 2025* ; Gartner (2025). *Magic Quadrant — gestion des interfaces de programmation*. — *Positions attestées ; elles ne disent rien de la position au quadrant portant sur les plateformes d'intégration, non vérifiée (garde-fou R-6).*
