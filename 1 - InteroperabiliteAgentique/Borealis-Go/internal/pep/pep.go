// Package pep fournit un Policy Enforcement Point minimal : une garde **non
// contournable** avant tout appel d'outil, qui identifie (KYA), autorise selon
// le niveau d'autonomie, et journalise la décision. En cas d'indisponibilité,
// il **retombe fail-closed** du niveau préparateur (L1) au copilote (L0) —
// aucune action autonome n'est alors permise (invariant 3, §6.6, §12.4).
package pep

import (
	"errors"
	"sync/atomic"

	"github.com/agbruneau/borealis/internal/audit"
)

// Level est le niveau d'autonomie (échelle L0-L3 de l'Annexe B §1.3).
type Level int

const (
	// L0 : copilote — toute action requiert un humain.
	L0 Level = iota
	// L1 : préparateur — la préparation (appels d'outils) est autorisée ; la
	// release finale reste humaine (HITL).
	L1
)

// Erreurs de la garde.
var (
	ErrFailClosed = errors.New("pep: indisponible — retombée fail-closed L1->L0")
	ErrDenied     = errors.New("pep: action non autorisée à ce niveau (humain requis)")
)

// PEP est le point d'application de la politique. level est fixé à la
// construction (immuable) ; available est mutable et lu/écrit concurremment
// par des handlers HTTP — atomic.Bool évite la course de données (P2).
type PEP struct {
	journal   *audit.Journal
	level     Level
	available atomic.Bool
}

// New construit un PEP disponible au niveau donné.
func New(journal *audit.Journal, level Level) *PEP {
	p := &PEP{journal: journal, level: level}
	p.available.Store(true)
	return p
}

// SetAvailable simule la disponibilité du PEP/orchestrateur (test fail-closed).
func (p *PEP) SetAvailable(v bool) { p.available.Store(v) }

// Guard autorise (ou non) une action avant un appel d'outil, et journalise la
// décision. C'est la voie **non contournable** : la logique d'appel d'outil
// passe par Guard, jamais autour.
func (p *PEP) Guard(kya, kyc, action string) error {
	switch {
	case !p.available.Load():
		p.audit(kya, kyc, "pep.fail_closed", "L1->L0 deny")
		return ErrFailClosed
	case p.level < L1:
		p.audit(kya, kyc, "pep.deny", "L0 requires human")
		return ErrDenied
	default:
		p.audit(kya, kyc, "pep.allow:"+action, "allow")
		return nil
	}
}

func (p *PEP) audit(kya, kyc, action, result string) {
	if p.journal != nil {
		p.journal.Append(audit.Record{
			KYA: kya, OwnerKYH: "human:credit-lead@borealis.example", KYC: kyc,
			Action: action, Result: result, Version: "m3",
		})
	}
}
