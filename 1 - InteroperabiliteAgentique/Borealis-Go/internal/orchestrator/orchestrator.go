// Package orchestrator implémente l'orchestrateur A2A : découverte des agents
// pairs par Agent Card (URL de base **configurées**, aucun endpoint en dur —
// FR-02), vérification de signature avant routage (FR-03), et sélection par
// compétence (skill, G-1). Retirer un agent de la configuration suffit à ne
// plus le router.
//
// Défenses (durcissement M2) : (1) validation de schéma de la card ;
// (2) **épinglage de l'interface à l'hôte de base** (anti-SSRF : un pair ne
// peut pas rediriger la délégation vers un autre hôte) ; (3) **rejet des
// collisions de compétence** (un skill revendiqué par ≥ 2 agents est ambigu →
// non routé, anti-détournement).
package orchestrator

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"iter"
	"net/http"
	"net/url"
	"sort"
	"sync"
	"time"

	a2asdk "github.com/a2aproject/a2a-go/v2/a2a"
	"github.com/a2aproject/a2a-go/v2/a2aclient"
	"github.com/a2aproject/a2a-go/v2/a2aclient/agentcard"
	"github.com/agbruneau/borealis/internal/audit"
	"github.com/agbruneau/borealis/internal/observability"
	borealisa2a "github.com/agbruneau/borealis/pkg/a2a"
)

// AgentRef est un pair découvert et vérifié (signature de confiance valide).
type AgentRef struct {
	Card   *a2asdk.AgentCard
	Client *a2aclient.Client
	Skills map[string]bool
}

// SupportsStreaming indique si le pair déclare la capability streaming (FR-04).
func (r *AgentRef) SupportsStreaming() bool {
	return r.Card != nil && r.Card.Capabilities.Streaming
}

// Orchestrator découvre les pairs et route les délégations par compétence.
type Orchestrator struct {
	baseURLs   []string
	trust      map[string]*ecdsa.PublicKey
	httpClient *http.Client
	mu         sync.RWMutex // protège bySkill/collisions (écrits par Discover, lus par Route/Delegate)
	bySkill    map[string]*AgentRef
	collisions map[string]bool
	cache      *delegationCache // déduplication des délégations (FR-18)
	journal    *audit.Journal   // optionnel : trace des délégations et doublons
}

// New construit l'orchestrateur avec les URL de base configurées et les clés
// publiques de confiance (indexées par kid). Le client HTTP propage le
// traceparent W3C sur les appels A2A sortants (NFR-07).
func New(baseURLs []string, trust map[string]*ecdsa.PublicKey) *Orchestrator {
	return &Orchestrator{
		baseURLs:   baseURLs,
		trust:      trust,
		httpClient: &http.Client{Transport: &observability.TraceRoundTripper{}},
		bySkill:    make(map[string]*AgentRef),
		collisions: make(map[string]bool),
		cache:      &delegationCache{seen: make(map[string]a2asdk.SendMessageResult)},
	}
}

// Discover résout chaque card configurée, vérifie sa signature, la valide, et
// épingle son interface à l'hôte de base avant d'indexer ses compétences. Un
// pair non fiable, malformé, ou pointant hors de son hôte de base est **rejeté**
// (non routé) et listé dans rejected (FR-03, anti-SSRF).
func (o *Orchestrator) Discover(ctx context.Context) (rejected []string, err error) {
	cands := make([]*AgentRef, 0, len(o.baseURLs))
	for _, base := range o.baseURLs {
		card, e := agentcard.DefaultResolver.Resolve(ctx, base)
		if e != nil {
			return rejected, fmt.Errorf("orchestrator: résolution %s: %w", base, e)
		}
		if o.verify(card) != nil || borealisa2a.Validate(card) != nil || !interfacesMatchBase(card, base) {
			rejected = append(rejected, base)
			continue
		}
		client, e := a2aclient.NewFromCard(ctx, card, a2aclient.WithJSONRPCTransport(o.httpClient))
		if e != nil {
			return rejected, fmt.Errorf("orchestrator: client %s: %w", base, e)
		}
		cands = append(cands, &AgentRef{Card: card, Client: client, Skills: skillSet(card)})
	}
	o.index(cands)
	return rejected, nil
}

// index route les compétences à **fournisseur unique** ; une compétence
// revendiquée par plusieurs agents est ambiguë → non routée (anti-détournement :
// un agent malveillant ne peut pas capter un skill légitime en le revendiquant).
// Reconstruit bySkill/collisions **à neuf** à chaque appel, à partir de
// l'ensemble complet des candidats résolus par CET appel : un agent retiré de
// la configuration disparaît du routage, et un agent introduit à un Discover
// ultérieur qui revendique un skill déjà indexé produit une collision — jamais
// un écrasement silencieux de l'entrée précédente (M3).
func (o *Orchestrator) index(cands []*AgentRef) {
	count := map[string]int{}
	for _, c := range cands {
		for s := range c.Skills {
			count[s]++
		}
	}
	bySkill := make(map[string]*AgentRef)
	collisions := make(map[string]bool)
	for _, c := range cands {
		for s := range c.Skills {
			if count[s] == 1 {
				bySkill[s] = c
			} else {
				collisions[s] = true
			}
		}
	}
	o.mu.Lock()
	o.bySkill, o.collisions = bySkill, collisions
	o.mu.Unlock()
}

// ponytail: borealisa2a.Verify recanonicalise la card (Marshal→Decode→Marshal)
// à CHAQUE itération, alors que le résultat est identique pour toutes les
// signatures d'une même card (les signatures en sont exclues). Borné par le
// nombre de signatures par card — au plus 1-2 en pratique (un seul algorithme
// actif, ADR-0006) — donc laissé tel quel ; si une fenêtre de migration
// crypto-agile fait cohabiter plusieurs signatures par card, hisser la
// canonicalisation hors de la boucle (exigerait d'exposer une variante de
// Verify acceptant le payload précalculé, pkg/a2a/card.go).
func (o *Orchestrator) verify(card *a2asdk.AgentCard) error {
	for _, sig := range card.Signatures {
		kid, e := borealisa2a.KID(sig)
		if e != nil {
			continue
		}
		if pub, ok := o.trust[kid]; ok {
			if borealisa2a.Verify(card, sig, pub) == nil {
				return nil
			}
		}
	}
	return errors.New("aucune signature de confiance valide")
}

// interfacesMatchBase vérifie que toutes les interfaces de la card pointent vers
// le **même hôte ET le même schéma** que l'URL de base configurée (anti-SSRF).
// Le schéma est vérifié en plus de l'hôte (m10) : sans ce contrôle, un pair
// pourrait déclarer une interface http:// sous une base https:// — même hôte,
// donc acceptée par le seul contrôle d'hôte, mais délégation rétrogradée en
// clair sous une base configurée en TLS.
func interfacesMatchBase(card *a2asdk.AgentCard, base string) bool {
	bu, err := url.Parse(base)
	if err != nil || len(card.SupportedInterfaces) == 0 {
		return false
	}
	for _, iface := range card.SupportedInterfaces {
		iu, err := url.Parse(iface.URL)
		if err != nil || iu.Host != bu.Host || iu.Scheme != bu.Scheme {
			return false
		}
	}
	return true
}

func skillSet(card *a2asdk.AgentCard) map[string]bool {
	s := make(map[string]bool, len(card.Skills))
	for _, sk := range card.Skills {
		s[sk.ID] = true
	}
	return s
}

// Route retourne le pair exposant la compétence demandée (G-1).
func (o *Orchestrator) Route(skill string) (*AgentRef, bool) {
	o.mu.RLock()
	defer o.mu.RUnlock()
	ref, ok := o.bySkill[skill]
	return ref, ok
}

// skills liste, triées, les compétences routables. Non exportée (B10) : lue
// par les seuls tests du paquet, aucun appelant de production.
func (o *Orchestrator) skills() []string {
	o.mu.RLock()
	defer o.mu.RUnlock()
	out := make([]string, 0, len(o.bySkill))
	for s := range o.bySkill {
		out = append(out, s)
	}
	sort.Strings(out)
	return out
}

// collisionSkills liste, triées, les compétences ambiguës (non routées). Non
// exportée (B10) : lue par les seuls tests du paquet, aucun appelant de
// production.
func (o *Orchestrator) collisionSkills() []string {
	o.mu.RLock()
	defer o.mu.RUnlock()
	out := make([]string, 0, len(o.collisions))
	for s := range o.collisions {
		out = append(out, s)
	}
	sort.Strings(out)
	return out
}

// Delegate délègue un message à l'agent d'une compétence. La délégation est
// **idempotente par (contextId, compétence, params)** : un rejeu ne produit
// qu'un seul effet et le doublon est journalisé (FR-18, NFR-05). Le transport
// est choisi selon la capability déclarée (FR-04) : streaming si l'agent le
// supporte (flux SSE drainé, dernier résultat retourné), sinon message/send.
func (o *Orchestrator) Delegate(ctx context.Context, skill string, msg *a2asdk.Message) (a2asdk.SendMessageResult, error) {
	ref, ok := o.Route(skill)
	if !ok {
		return nil, fmt.Errorf("orchestrator: aucun agent pour la compétence %q", skill)
	}
	return o.idempotent(msg.ContextID, skill, msg, func() (a2asdk.SendMessageResult, error) {
		return deliver(ctx, ref, msg)
	})
}

// defaultDelegateTimeout borne les DEUX transports (streaming et message/send)
// quand l'appelant n'a fixé aucune échéance (m6) : sans cette garde, un pair
// muet bloque Delegate indéfiniment — et les cards de production déclarant
// toutes Streaming: true, le chemin streaming est celui réellement pris.
// Variable (non const) pour être raccourcie en test.
var defaultDelegateTimeout = 30 * time.Second

// deliver exécute l'envoi effectif au pair selon son transport (le seul « effet »
// que l'idempotence protège d'un rejeu). Le flux SSE est entièrement drainé ici
// (reduceStream), donc l'échéance couvre bien toute la délégation.
func deliver(ctx context.Context, ref *AgentRef, msg *a2asdk.Message) (a2asdk.SendMessageResult, error) {
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, defaultDelegateTimeout)
		defer cancel()
	}
	req := &a2asdk.SendMessageRequest{Message: msg}
	if ref.SupportsStreaming() {
		return reduceStream(ref.Client.SendStreamingMessage(ctx, req))
	}
	return ref.Client.SendMessage(ctx, req)
}

// reduceStream réduit un flux SSE à son **dernier résultat observable**. Le
// flux émet un Task (submitted) puis des TaskStatusUpdateEvent (working,
// completed) : ces derniers ne satisfont PAS SendMessageResult, il faut donc
// **replier leur statut** sur le Task suivi pour retourner son état final, pas
// le Task submitted initial (sinon le résultat mis en cache serait périmé —
// FR-05/07/08). Un flux qui se termine sans **aucun** événement retourne une
// erreur, jamais (nil, nil) : sinon ce nil deviendrait un résultat de
// délégation valide, mis en cache d'idempotence (m4).
func reduceStream(events iter.Seq2[a2asdk.Event, error]) (a2asdk.SendMessageResult, error) {
	var task *a2asdk.Task
	var last a2asdk.SendMessageResult
	for ev, streamErr := range events {
		if streamErr != nil {
			return nil, streamErr
		}
		switch e := ev.(type) {
		case *a2asdk.Task:
			task, last = e, e
		case *a2asdk.Message:
			last = e
		case *a2asdk.TaskStatusUpdateEvent:
			if task != nil {
				task.Status = e.Status // état final (completed) + message de réponse
				last = task
			}
		}
	}
	if last == nil {
		return nil, errors.New("orchestrator: flux SSE vide (aucun événement reçu)")
	}
	return last, nil
}
