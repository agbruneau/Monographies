package mcpserver

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/agbruneau/borealis/pkg/mcpcontract"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// toolSchemas connecte un serveur, liste ses outils et retourne les schémas
// (entrée, sortie) de l'outil nommé, parsés pour pkg/mcpcontract.
func toolSchemas(t *testing.T, s *mcp.Server, tool string) (in, out mcpcontract.Schema) {
	t.Helper()
	cs := connectServer(t, s)
	lt, err := cs.ListTools(context.Background(), nil)
	if err != nil {
		t.Fatalf("ListTools : %v", err)
	}
	for _, to := range lt.Tools {
		if to.Name != tool {
			continue
		}
		ib, _ := json.Marshal(to.InputSchema)
		if in, err = mcpcontract.Parse(ib); err != nil {
			t.Fatalf("parse input : %v", err)
		}
		if to.OutputSchema != nil {
			ob, _ := json.Marshal(to.OutputSchema)
			if out, err = mcpcontract.Parse(ob); err != nil {
				t.Fatalf("parse output : %v", err)
			}
		}
		return in, out
	}
	t.Fatalf("outil %q absent", tool)
	return
}

// TestCapacitySubstitutionV1toV2 : le Capacity Calculator v2 se substitue à v1
// **sans modifier le consommateur** — qui lit toujours la sortie dans la struct
// v1 (CapacityOut). Pour toute entrée v1, la décision (mensualité, DTI, capacité)
// est **identique** ; l'e2e du consommateur reste vert (M4.5, DoD 5, G-4, US-9).
func TestCapacitySubstitutionV1toV2(t *testing.T) {
	v1 := connectServer(t, NewCapacityServer())
	v2 := connectServer(t, NewCapacityServerV2())

	cases := []map[string]any{
		{"income": 120000.0, "debts": 1200.0, "requestedAmount": 30000.0, "termMonths": 60},
		{"income": 60000.0, "debts": 2000.0, "requestedAmount": 40000.0, "termMonths": 48},
		{"income": 240000.0, "debts": 500.0, "requestedAmount": 15000.0, "termMonths": 36},
		{"income": 0.0, "debts": 0.0, "requestedAmount": 10000.0, "termMonths": 0}, // dégénéré → hors capacité
	}
	for _, args := range cases {
		// Consommateur INCHANGÉ : décode dans la struct v1, avec les mêmes args v1.
		var d1, d2 CapacityOut
		callStructured(t, v1, "compute_capacity", args, &d1)
		callStructured(t, v2, "compute_capacity", args, &d2)
		if d1 != d2 {
			t.Errorf("args=%v : la substitution change la décision v1\n v1=%+v\n v2=%+v", args, d1, d2)
		}
		// La sortie v2 porte bien le champ ajouté (présent, mais ignoré ci-dessus).
		var enriched CapacityOutV2
		callStructured(t, v2, "compute_capacity", args, &enriched)
		if enriched.RiskBand == "" {
			t.Errorf("args=%v : v2 devrait exposer riskBand", args)
		}
	}
}

// TestRiskBandConsistency : la bande de risque v2 est cohérente avec capacityOk
// v1 — jamais « high » sur un dossier que v1 approuve. Cas limite du réfutateur :
// demande valide au DTI arrondi à 0,00.
func TestRiskBandConsistency(t *testing.T) {
	// income énorme, montant minuscule → DTI ≈ 0 arrondi à 0,00, mais capacityOk.
	safe, err := computeCapacityV2(CapacityInV2{Income: 1_000_000, Debts: 0, RequestedAmount: 100, TermMonths: 360})
	if err != nil {
		t.Fatalf("prérequis : entrée valide, erreur inattendue : %v", err)
	}
	if !safe.CapacityOk {
		t.Fatalf("prérequis : la demande devrait être dans les capacités : %+v", safe)
	}
	if safe.RiskBand == "high" {
		t.Errorf("demande valide (capacityOk) classée high : %+v — incohérence de contrat", safe)
	}
	// Demande dégénérée (durée nulle) → hors capacité → high.
	degen, err := computeCapacityV2(CapacityInV2{Income: 60000, Debts: 0, RequestedAmount: 10000, TermMonths: 0})
	if err != nil {
		t.Fatalf("prérequis : entrée valide (dégénérée mais stressRate nul), erreur inattendue : %v", err)
	}
	if degen.CapacityOk || degen.RiskBand != "high" {
		t.Errorf("demande dégénérée : %+v, want capacityOk=false riskBand=high", degen)
	}
}

// TestCapacityContractBackwardV1toV2 confirme, via pkg/mcpcontract, que le
// contrat v2 est BACKWARD-compatible avec v1 : v2 accepte toute entrée v1 (le
// champ ajouté stressRate est optionnel) et sa sortie reste lisible par un
// consommateur v1 (riskBand est un ajout). L'assertion négative garantit que le
// contrôle n'est pas vide : un v2 qui rendrait stressRate requis casserait.
func TestCapacityContractBackwardV1toV2(t *testing.T) {
	v1in, v1out := toolSchemas(t, NewCapacityServer(), "compute_capacity")
	v2in, v2out := toolSchemas(t, NewCapacityServerV2(), "compute_capacity")

	if b := mcpcontract.InputBackwardCompatible(v1in, v2in); len(b) != 0 {
		t.Errorf("entrée v2 non BACKWARD-compatible : %v", b)
	}
	if b := mcpcontract.OutputBackwardCompatible(v1out, v2out); len(b) != 0 {
		t.Errorf("sortie v2 non BACKWARD-compatible : %v", b)
	}

	// Oracle d'entrée : un v2 rendant stressRate REQUIS romprait la compat.
	broken := v2in
	broken.Required = map[string]bool{}
	for k := range v2in.Required {
		broken.Required[k] = true
	}
	broken.Required["stressRate"] = true
	if b := mcpcontract.InputBackwardCompatible(v1in, broken); len(b) == 0 {
		t.Error("un champ ajouté rendu requis devrait rompre la compat (contrôle vide ?)")
	}

	// Oracle de sortie : un v2 qui SUPPRIME un champ v1 romprait la lecture v1.
	brokenOut := mcpcontract.Schema{
		Properties: map[string]string{}, Nullable: map[string]bool{},
		Required: map[string]bool{}, AdditionalProperties: v2out.AdditionalProperties,
	}
	for k, typ := range v2out.Properties {
		if k != "capacityOk" { // régression simulée : le champ v1 disparaît
			brokenOut.Properties[k] = typ
		}
	}
	if b := mcpcontract.OutputBackwardCompatible(v1out, brokenOut); len(b) == 0 {
		t.Error("suppression d'un champ de sortie v1 devrait rompre la compat (contrôle vide ?)")
	}
}
