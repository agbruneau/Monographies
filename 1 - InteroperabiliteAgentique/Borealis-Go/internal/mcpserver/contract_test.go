package mcpserver

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/agbruneau/borealis/internal/fixtures"
	"github.com/agbruneau/borealis/pkg/mcpcontract"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// TestServerContractsSelfBackwardCompatible branche pkg/mcpcontract sur les 4
// serveurs : chaque schéma d'outil inféré doit être BACKWARD-compatible avec
// lui-même (invariant minimal exécuté par le gate ; la substitution v1→v2
// réelle est vérifiée en M4.5, G-4).
func TestServerContractsSelfBackwardCompatible(t *testing.T) {
	ds := fixtures.Generate(fixtures.DefaultSeed, 10)
	servers := map[string]*mcp.Server{
		"identity": NewIdentityServer(ds),
		"bureau":   NewBureauServer(ds),
		"capacity": NewCapacityServer(),
		"policy":   NewPolicyServer(ds),
	}
	for name, s := range servers {
		cs := connectServer(t, s)
		lt, err := cs.ListTools(context.Background(), nil)
		if err != nil {
			t.Fatalf("%s ListTools : %v", name, err)
		}
		if len(lt.Tools) == 0 {
			t.Errorf("%s : aucun outil exposé", name)
		}
		for _, tool := range lt.Tools {
			in, _ := json.Marshal(tool.InputSchema)
			v1, err := mcpcontract.Parse(in)
			if err != nil {
				t.Fatalf("%s/%s parse input : %v", name, tool.Name, err)
			}
			if b := mcpcontract.InputBackwardCompatible(v1, v1); len(b) != 0 {
				t.Errorf("%s/%s entrée auto-incompatible : %v", name, tool.Name, b)
			}
			if tool.OutputSchema != nil {
				out, _ := json.Marshal(tool.OutputSchema)
				vo, err := mcpcontract.Parse(out)
				if err != nil {
					t.Fatalf("%s/%s parse output : %v", name, tool.Name, err)
				}
				if b := mcpcontract.OutputBackwardCompatible(vo, vo); len(b) != 0 {
					t.Errorf("%s/%s sortie auto-incompatible : %v", name, tool.Name, b)
				}
			}
		}
	}
}

// TestContractDetectsRealBreak donne au contrôle un ORACLE réel (au-delà de
// l'auto-comparaison) : un v2 dérivé d'un schéma serveur qui exige un nouveau
// champ doit être détecté comme incompatible BACKWARD.
func TestContractDetectsRealBreak(t *testing.T) {
	ds := fixtures.Generate(fixtures.DefaultSeed, 10)
	cs := connectServer(t, NewIdentityServer(ds))
	lt, err := cs.ListTools(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	var v1 mcpcontract.Schema
	for _, tool := range lt.Tools {
		if tool.Name == "verify_identity" {
			in, _ := json.Marshal(tool.InputSchema)
			if v1, err = mcpcontract.Parse(in); err != nil {
				t.Fatal(err)
			}
		}
	}
	// v2 diverge : exige un nouveau champ (mfaToken) absent en v1.
	v2 := mcpcontract.Schema{
		Properties:           map[string]string{"applicantId": "string", "mfaToken": "string"},
		Nullable:             map[string]bool{},
		Required:             map[string]bool{"applicantId": true, "mfaToken": true},
		AdditionalProperties: false,
	}
	if b := mcpcontract.InputBackwardCompatible(v1, v2); len(b) == 0 {
		t.Error("ajout d'un champ requis divergent non détecté (oracle réel)")
	}
}
