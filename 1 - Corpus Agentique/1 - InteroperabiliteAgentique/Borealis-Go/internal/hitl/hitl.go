// Package hitl gère l'intervention humaine (Human-In-The-Loop) par **deux files
// distinctes** : proposition (agent → humain, via input-required) et commande
// (humain → système, après release). Aucun message n'entre en file de commande
// sans une **release journalisée** (FR-16, US-7) — c'est le point checker du
// maker-checker (M3.7).
package hitl

import (
	"errors"
	"sync"

	"github.com/agbruneau/borealis/internal/audit"
)

// Proposal est une proposition soumise à un humain (file de proposition).
type Proposal struct {
	ID      string
	Subject string // sujet pseudonymisé (jamais de PII)
	Summary string // motif (ex. « AML élevé », « signaux conflictuels »)
}

// Command est une commande émise après release humaine (file de commande).
type Command struct {
	ProposalID string
	Decision   string // approve | reject
	Human      string // KYH du releaseur
}

// Décisions de release valides (m11 : toute autre valeur, ex. faute de frappe,
// est refusée par Release).
const (
	DecisionApprove = "approve"
	DecisionReject  = "reject"
)

// Erreurs de release.
var (
	ErrUnknownProposal = errors.New("hitl: proposition inconnue")
	ErrAlreadyReleased = errors.New("hitl: proposition déjà traitée")
	ErrInvalidDecision = errors.New("hitl: décision invalide (approve|reject attendu)")
)

// Broker relie les deux files. La file de commande n'est alimentée que par Release.
type Broker struct {
	mu        sync.Mutex
	journal   *audit.Journal
	proposals map[string]Proposal
	released  map[string]bool
	commands  []Command
}

// New construit un broker HITL.
func New(journal *audit.Journal) *Broker {
	return &Broker{
		journal:   journal,
		proposals: make(map[string]Proposal),
		released:  make(map[string]bool),
	}
}

// Propose place une proposition en file de proposition (input-required).
func (b *Broker) Propose(p Proposal) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.proposals[p.ID] = p
	b.audit("agent:orchestrator", "", p.Subject, "hitl.propose", p.Summary)
}

// Release traite une proposition par un humain : **seul chemin** vers la file de
// commande. La release est journalisée AVANT l'ajout de la commande.
func (b *Broker) Release(proposalID, decision, human string) error {
	if decision != DecisionApprove && decision != DecisionReject {
		return ErrInvalidDecision
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	p, ok := b.proposals[proposalID]
	if !ok {
		return ErrUnknownProposal
	}
	if b.released[proposalID] {
		return ErrAlreadyReleased
	}
	b.released[proposalID] = true
	b.audit("agent:orchestrator", human, p.Subject, "hitl.release", decision)
	b.commands = append(b.commands, Command{ProposalID: proposalID, Decision: decision, Human: human})
	return nil
}

// PendingProposals retourne le nombre de propositions non encore traitées.
func (b *Broker) PendingProposals() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	n := 0
	for id := range b.proposals {
		if !b.released[id] {
			n++
		}
	}
	return n
}

// Commands retourne une copie de la file de commande.
func (b *Broker) Commands() []Command {
	b.mu.Lock()
	defer b.mu.Unlock()
	out := make([]Command, len(b.commands))
	copy(out, b.commands)
	return out
}

func (b *Broker) audit(kya, human, kyc, action, result string) {
	if b.journal != nil {
		b.journal.Append(audit.Record{KYA: kya, OwnerKYH: human, KYC: kyc, Action: action, Result: result, Version: "m3"})
	}
}
