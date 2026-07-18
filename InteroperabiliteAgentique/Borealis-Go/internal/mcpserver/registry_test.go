package mcpserver

import "testing"

func TestDefaultRegistry(t *testing.T) {
	r, err := DefaultRegistry()
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"identity-hub", "credit-bureau-sim", "capacity-calculator", "policy-engine"} {
		s, ok := r.Server(name)
		if !ok {
			t.Errorf("serveur %s absent du registre", name)
			continue
		}
		if s.ContractVersion == "" || len(s.Tools) == 0 || len(s.Errors) == 0 || s.SLOP95Ms <= 0 {
			t.Errorf("%s incomplet : %+v", name, s)
		}
	}
	if _, ok := r.Server("inexistant"); ok {
		t.Error("serveur inexistant trouvé")
	}
}

func TestRegistryValidate(t *testing.T) {
	bad := []struct {
		name string
		yaml string
	}{
		{"vide", "servers: []"},
		{"nom vide", "servers:\n  - contractVersion: '1'\n    tools: [a]\n    errors: [X]"},
		{"sans version", "servers:\n  - name: x\n    tools: [a]\n    errors: [X]"},
		{"sans outil", "servers:\n  - name: x\n    contractVersion: '1'\n    errors: [X]"},
		{"sans erreurs", "servers:\n  - name: x\n    contractVersion: '1'\n    tools: [a]"},
	}
	for _, c := range bad {
		if _, err := LoadRegistry([]byte(c.yaml)); err == nil {
			t.Errorf("%s : erreur de validation attendue", c.name)
		}
	}
	// YAML malformé → erreur de parsing.
	if _, err := LoadRegistry([]byte("servers: [unterminated")); err == nil {
		t.Error("YAML malformé : erreur attendue")
	}
}
