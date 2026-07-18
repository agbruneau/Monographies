package agent

import (
	"context"
	"fmt"
	"time"

	"github.com/agbruneau/borealis/internal/observability"
	"github.com/agbruneau/borealis/internal/pep"
	"github.com/agbruneau/borealis/internal/resilience"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// ResilientCaller enveloppe une session MCP d'un disjoncteur, d'un retry borné
// et d'un timeout par tentative (FR-17, NFR-04). C'est la voie standard par
// laquelle les agents appellent les serveurs MCP.
type ResilientCaller struct {
	sess    *mcp.ClientSession
	breaker *resilience.Breaker
	retry   resilience.RetryPolicy
	timeout time.Duration
	pep     *pep.PEP // optionnel : garde non contournable (M3.2)
	kya     string
}

// WithPEP branche un Policy Enforcement Point (garde non contournable) et
// l'identité KYA de l'agent appelant (M3.2, invariant 3).
func (c *ResilientCaller) WithPEP(p *pep.PEP, kya string) *ResilientCaller {
	c.pep = p
	c.kya = kya
	return c
}

// NewResilientCaller construit un caller aux réglages par défaut du démonstrateur
// (disjoncteur : 5 échecs / cooldown 30 s ; 3 tentatives ; timeout 500 ms).
func NewResilientCaller(sess *mcp.ClientSession) *ResilientCaller {
	return &ResilientCaller{
		sess:    sess,
		breaker: resilience.NewBreaker(5, 30*time.Second),
		retry:   resilience.RetryPolicy{MaxAttempts: 3, Backoff: 10 * time.Millisecond},
		timeout: 500 * time.Millisecond,
	}
}

// CallTool appelle un outil MCP sous disjoncteur + retry + timeout par tentative.
func (c *ResilientCaller) CallTool(ctx context.Context, name string, args map[string]any) (*mcp.CallToolResult, error) {
	// Garde non contournable : aucun appel d'outil sans autorisation du PEP.
	if c.pep != nil {
		if err := c.pep.Guard(c.kya, "", name); err != nil {
			return nil, err
		}
	}
	var res *mcp.CallToolResult
	err := c.breaker.Execute(ctx, c.retry, func() error {
		cctx, cancel := context.WithTimeout(ctx, c.timeout)
		defer cancel()
		params := &mcp.CallToolParams{Name: name, Arguments: args}
		// Propagation traceparent W3C via _meta (trace ininterrompue, NFR-07).
		if tp := observability.InjectTraceparent(ctx); len(tp) > 0 {
			meta := mcp.Meta{}
			for k, v := range tp {
				meta[k] = v
			}
			params.Meta = meta
		}
		r, err := c.sess.CallTool(cctx, params)
		if err != nil {
			// Timeout PAR TENTATIVE (cctx expiré, ctx appelant vivant) : c'est une
			// panne du service distant — pas une annulation de l'appelant — et elle
			// doit compter comme échec du disjoncteur (NFR-04). Ne pas envelopper
			// (%w) l'erreur de contexte sous-jacente : errors.Is(…, DeadlineExceeded)
			// la ferait re-neutraliser par onResult (m5).
			if cctx.Err() != nil && ctx.Err() == nil {
				return fmt.Errorf("agent: délai de tentative dépassé (%s) sur l'outil %s: %v", c.timeout, name, err)
			}
			return err
		}
		res = r
		return nil
	})
	return res, err
}
