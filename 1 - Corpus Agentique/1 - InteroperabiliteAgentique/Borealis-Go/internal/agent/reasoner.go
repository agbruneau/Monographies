// Package agent porte les abstractions communes aux agents Borealis :
// l'interface Reasoner (couture du moteur de raisonnement) et ses types de
// requête/résultat. Un LLM réel n'entre jamais dans le gate (déterminisme,
// G-8) ; le défaut d'exécution et de test est un fake (internal/testfakes).
package agent

import "context"

// Recommendation est l'orientation produite par un Reasoner. Aucune valeur
// n'exprime un octroi ferme (FR-25) : au mieux « pré-qualifier ».
type Recommendation string

const (
	// RecommendPreQualify : orientation vers une pré-qualification.
	RecommendPreQualify Recommendation = "pre_qualify"
	// RecommendDecline : orientation vers un refus.
	RecommendDecline Recommendation = "decline"
	// RecommendEscalate : renvoi à un humain (HITL).
	RecommendEscalate Recommendation = "escalate"
)

// ReasonRequest agrège les signaux collectés pour une demande.
type ReasonRequest struct {
	ApplicantID string
	Signals     map[string]string
}

// ReasonResult est la sortie déterminisable du raisonnement.
type ReasonResult struct {
	Recommendation Recommendation
	Rationale      string
	Confidence     float64
}

// Reasoner est la couture du moteur de raisonnement (règles ou LLM).
type Reasoner interface {
	Reason(ctx context.Context, req ReasonRequest) (ReasonResult, error)
}
