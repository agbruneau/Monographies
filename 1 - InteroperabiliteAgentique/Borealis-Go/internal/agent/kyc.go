package agent

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// IdentityCheck est le résultat de la vérification d'identité (via MCP).
type IdentityCheck struct {
	Match  bool
	Reason string
}

// VerifyIdentity appelle l'outil MCP verify_identity via un caller résilient
// (disjoncteur + retry + timeout). C'est l'« appel MCP réel » du KYC (FR-11).
func VerifyIdentity(ctx context.Context, caller *ResilientCaller, applicantID string) (IdentityCheck, error) {
	res, err := caller.CallTool(ctx, "verify_identity", map[string]any{"applicantId": applicantID})
	if err != nil {
		return IdentityCheck{}, fmt.Errorf("verify_identity: %w", err)
	}
	return parseIdentityResult(res)
}

// parseIdentityResult extrait l'IdentityCheck d'un résultat d'outil. Une
// violation de schéma est reportée par le SDK en IsError (ADR-0001).
func parseIdentityResult(res *mcp.CallToolResult) (IdentityCheck, error) {
	out, err := decode[struct {
		Match  bool   `json:"match"`
		Reason string `json:"reason"`
	}](res, "verify_identity")
	if err != nil {
		return IdentityCheck{}, err
	}
	return IdentityCheck{Match: out.Match, Reason: out.Reason}, nil
}

// decode extrait T d'un résultat d'outil MCP : échec métier (IsError,
// ADR-0001), erreur de sérialisation ou de désérialisation — factorisé entre
// VerifyIdentity et ScreenWatchlists (B6), auparavant dupliqué.
func decode[T any](res *mcp.CallToolResult, tool string) (T, error) {
	var out T
	if res.IsError {
		return out, fmt.Errorf("%s rejected: %s", tool, toolErrorText(res))
	}
	b, err := json.Marshal(res.StructuredContent)
	if err != nil {
		return out, fmt.Errorf("%s: marshal: %w", tool, err)
	}
	if err := json.Unmarshal(b, &out); err != nil {
		return out, fmt.Errorf("%s: unmarshal: %w", tool, err)
	}
	return out, nil
}

// Screening est le résultat du criblage AML/listes de surveillance (via MCP
// check_watchlists). Le nom est résolu côté serveur depuis le dossier : l'agent
// n'a que l'identifiant à fournir.
type Screening struct {
	WatchlistHit bool
	AMLScore     float64
	Reason       string
}

// ScreenWatchlists appelle l'outil MCP check_watchlists via le caller résilient.
func ScreenWatchlists(ctx context.Context, caller *ResilientCaller, applicantID string) (Screening, error) {
	res, err := caller.CallTool(ctx, "check_watchlists", map[string]any{"applicantId": applicantID})
	if err != nil {
		return Screening{}, fmt.Errorf("check_watchlists: %w", err)
	}
	out, err := decode[struct {
		WatchlistHit bool    `json:"watchlistHit"`
		AMLScore     float64 `json:"amlScore"`
		Reason       string  `json:"reason"`
	}](res, "check_watchlists")
	if err != nil {
		return Screening{}, err
	}
	return Screening{WatchlistHit: out.WatchlistHit, AMLScore: out.AMLScore, Reason: out.Reason}, nil
}

// toolErrorText extrait le texte d'erreur d'un résultat d'outil en échec.
func toolErrorText(res *mcp.CallToolResult) string {
	for _, c := range res.Content {
		if tc, ok := c.(*mcp.TextContent); ok {
			return tc.Text
		}
	}
	return "unknown tool error"
}
