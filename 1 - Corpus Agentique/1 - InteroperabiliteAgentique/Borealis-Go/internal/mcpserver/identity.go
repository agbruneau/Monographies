// Package mcpserver assemble les serveurs MCP du démonstrateur. Le schéma
// d'entrée est inféré par le SDK à partir du type Go ; une entrée invalide est
// reportée en résultat IsError (ADR-0001), non en -32602.
package mcpserver

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"math"
	"strings"

	"github.com/agbruneau/borealis/internal/fixtures"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// VerifyIdentityInput est l'entrée de verify_identity (PRD §6.2).
type VerifyIdentityInput struct {
	ApplicantID string `json:"applicantId" jsonschema:"identifiant du demandeur"`
	Name        string `json:"name,omitempty" jsonschema:"nom complet (verifie si fourni)"`
	SIN         string `json:"sin,omitempty" jsonschema:"NAS fictif (verifie si fourni)"`
}

// VerifyIdentityOutput est la sortie de verify_identity.
type VerifyIdentityOutput struct {
	Match  bool   `json:"match"`
	Reason string `json:"reason"`
}

// CheckWatchlistsInput est l'entrée de check_watchlists.
type CheckWatchlistsInput struct {
	ApplicantID string `json:"applicantId" jsonschema:"identifiant du demandeur"`
	Name        string `json:"name,omitempty" jsonschema:"nom complet"`
}

// CheckWatchlistsOutput est la sortie de check_watchlists.
type CheckWatchlistsOutput struct {
	WatchlistHit bool    `json:"watchlistHit"`
	AMLScore     float64 `json:"amlScore"`
	Reason       string  `json:"reason"`
}

// ofacSimNames est une liste OFAC/PPE **simulée** (fictive, noms de la fixture).
var ofacSimNames = map[string]bool{"rossi": true, "diallo": true}

// verifyIdentity vérifie l'identité contre le jeu synthétique (logique pure).
func verifyIdentity(ds *fixtures.Dataset, in VerifyIdentityInput) VerifyIdentityOutput {
	if in.ApplicantID == "" {
		return VerifyIdentityOutput{Match: false, Reason: "missing_applicant_id"}
	}
	b, ok := ds.Borrower(in.ApplicantID)
	if !ok {
		return VerifyIdentityOutput{Match: false, Reason: "not_found"}
	}
	if in.Name != "" && !strings.EqualFold(in.Name, b.FirstName+" "+b.LastName) {
		return VerifyIdentityOutput{Match: false, Reason: "name_mismatch"}
	}
	if in.SIN != "" && in.SIN != b.SIN {
		return VerifyIdentityOutput{Match: false, Reason: "sin_mismatch"}
	}
	return VerifyIdentityOutput{Match: true, Reason: "verified"}
}

// amlScore dérive un score AML déterministe dans [0,1) (NFR-10) du couple id+nom.
func amlScore(applicantID, name string) float64 {
	sum := sha256.Sum256([]byte(applicantID + "|" + name))
	v := binary.BigEndian.Uint32(sum[:4])
	return float64(v) / float64(math.MaxUint32)
}

// checkWatchlists évalue les listes de surveillance simulées (logique pure). Le
// score AML doit être **discriminant** : un demandeur propre reste sous les
// seuils FINTRAC (0,7/0,8) — sinon le happy path échouerait pour ~30 % des
// dossiers propres ; un hit force un score élevé (refus + escalade HITL).
//
// Le nom criblé est **résolu côté serveur** depuis le dossier (jamais celui,
// falsifiable, fourni par le sujet) : le criblage AML porte ainsi sur l'identité
// autoritative. Le nom d'entrée n'est retenu qu'en repli (dossier introuvable).
func checkWatchlists(ds *fixtures.Dataset, in CheckWatchlistsInput) CheckWatchlistsOutput {
	name := in.Name
	if b, ok := ds.Borrower(in.ApplicantID); ok {
		name = b.FirstName + " " + b.LastName
	}
	raw := amlScore(in.ApplicantID, name) // [0,1) déterministe
	hit := false
	for _, part := range strings.Fields(strings.ToLower(name)) {
		if ofacSimNames[part] {
			hit = true
			break
		}
	}
	var score float64
	reason := "clear"
	if hit {
		score = 0.85 + raw*0.15 // hit → [0.85, 1.0) : au-dessus des seuils
		reason = "watchlist_match_simulated"
	} else {
		score = raw * 0.5 // propre → [0, 0.5) : sous les seuils FINTRAC
	}
	return CheckWatchlistsOutput{WatchlistHit: hit, AMLScore: round2(score), Reason: reason}
}

// NewIdentityServer construit le serveur MCP Identity Hub (verify_identity +
// check_watchlists) sur le jeu synthétique fourni.
func NewIdentityServer(ds *fixtures.Dataset) *mcp.Server {
	s := mcp.NewServer(&mcp.Implementation{Name: "identity-hub", Version: "v1.0.0"}, nil)
	mcp.AddTool(s, &mcp.Tool{
		Name:        "verify_identity",
		Description: "Verifie l'identite synthetique d'un demandeur (donnees fictives).",
	}, func(_ context.Context, _ *mcp.CallToolRequest, in VerifyIdentityInput) (*mcp.CallToolResult, VerifyIdentityOutput, error) {
		return nil, verifyIdentity(ds, in), nil
	})
	mcp.AddTool(s, &mcp.Tool{
		Name:        "check_watchlists",
		Description: "Verifie les listes PPE/OFAC simulees et calcule un score AML.",
	}, func(_ context.Context, _ *mcp.CallToolRequest, in CheckWatchlistsInput) (*mcp.CallToolResult, CheckWatchlistsOutput, error) {
		return nil, checkWatchlists(ds, in), nil
	})
	return s
}
