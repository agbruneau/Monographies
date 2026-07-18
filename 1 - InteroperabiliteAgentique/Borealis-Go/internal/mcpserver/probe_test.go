package mcpserver

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/agbruneau/borealis/internal/fixtures"
)

// TestVerifyIdentitySchemaContract vérifie que le SDK infère un schéma d'entrée
// strict : applicantId requis, de type string, et aucune propriété additionnelle.
// Test de régression du contrat d'entrée (le rejet effectif est couvert par
// TestVerifyIdentityInvalidInput).
func TestVerifyIdentitySchemaContract(t *testing.T) {
	ds := fixtures.Generate(fixtures.DefaultSeed, 10)
	cs := connectServer(t, NewIdentityServer(ds))
	lt, err := cs.ListTools(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	var found bool
	for _, tool := range lt.Tools {
		if tool.Name != "verify_identity" {
			continue
		}
		found = true
		b, _ := json.Marshal(tool.InputSchema)
		var schema struct {
			Required             []string `json:"required"`
			AdditionalProperties *bool    `json:"additionalProperties"`
			Properties           map[string]struct {
				Type string `json:"type"`
			} `json:"properties"`
		}
		if err := json.Unmarshal(b, &schema); err != nil {
			t.Fatalf("schéma illisible : %v", err)
		}
		if !contains(schema.Required, "applicantId") {
			t.Errorf("required = %v, doit contenir applicantId", schema.Required)
		}
		if schema.AdditionalProperties == nil || *schema.AdditionalProperties {
			t.Errorf("additionalProperties = %v, want false", schema.AdditionalProperties)
		}
		if schema.Properties["applicantId"].Type != "string" {
			t.Errorf("applicantId type = %q, want string", schema.Properties["applicantId"].Type)
		}
	}
	if !found {
		t.Fatal("outil verify_identity absent de ListTools")
	}
}

func contains(ss []string, s string) bool {
	for _, x := range ss {
		if x == s {
			return true
		}
	}
	return false
}
