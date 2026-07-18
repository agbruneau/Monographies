package mcpcontract

import "testing"

func TestParse(t *testing.T) {
	s, err := Parse([]byte(`{"type":"object","properties":{"id":{"type":"string"},"n":{"type":"integer"}},"required":["id"],"additionalProperties":false}`))
	if err != nil {
		t.Fatal(err)
	}
	if s.Properties["id"] != "string" || s.Properties["n"] != "integer" {
		t.Errorf("propriétés : %+v", s.Properties)
	}
	if !s.Required["id"] || s.Required["n"] {
		t.Errorf("required : %+v", s.Required)
	}
	if s.AdditionalProperties {
		t.Error("additionalProperties devrait être false")
	}
	// additionalProperties absent → défaut true.
	s2, _ := Parse([]byte(`{"properties":{"a":{"type":"string"}}}`))
	if !s2.AdditionalProperties {
		t.Error("additionalProperties par défaut devrait être true")
	}
	// type en union (["array","null"]) → type de base + nullable.
	if s3, _ := Parse([]byte(`{"properties":{"r":{"type":["array","null"]}}}`)); s3.Properties["r"] != "array" || !s3.Nullable["r"] {
		t.Errorf("type union : base=%q nullable=%v, want array/true", s3.Properties["r"], s3.Nullable["r"])
	}
	// type absent → chaîne vide.
	if s4, _ := Parse([]byte(`{"properties":{"x":{}}}`)); s4.Properties["x"] != "" {
		t.Errorf("type absent = %q, want vide", s4.Properties["x"])
	}
	// type non reconnu (nombre) → chaîne vide.
	if s5, _ := Parse([]byte(`{"properties":{"y":{"type":123}}}`)); s5.Properties["y"] != "" {
		t.Errorf("type non reconnu = %q, want vide", s5.Properties["y"])
	}
	// union entièrement null → base vide mais nullable.
	if s6, _ := Parse([]byte(`{"properties":{"z":{"type":["null"]}}}`)); s6.Properties["z"] != "" || !s6.Nullable["z"] {
		t.Errorf("union null : base=%q nullable=%v, want vide/true", s6.Properties["z"], s6.Nullable["z"])
	}
	// JSON invalide.
	if _, err := Parse([]byte(`{not json`)); err == nil {
		t.Error("erreur JSON attendue")
	}
}

func schema(props map[string]string, required []string, addl bool) Schema {
	s := Schema{Properties: props, Nullable: map[string]bool{}, Required: map[string]bool{}, AdditionalProperties: addl}
	for _, r := range required {
		s.Required[r] = true
	}
	return s
}

func TestInputBackwardCompatible(t *testing.T) {
	v1 := schema(map[string]string{"id": "string", "amt": "number"}, []string{"id"}, true)
	if b := InputBackwardCompatible(v1, v1); len(b) != 0 {
		t.Errorf("auto-compat : %v", b)
	}
	// v2 ajoute un champ requis → rupture.
	if b := InputBackwardCompatible(v1, schema(map[string]string{"id": "string", "amt": "number"}, []string{"id", "amt"}, true)); len(b) != 1 {
		t.Errorf("champ requis ajouté : %v", b)
	}
	// v2 change un type → rupture.
	if b := InputBackwardCompatible(v1, schema(map[string]string{"id": "string", "amt": "string"}, []string{"id"}, true)); len(b) != 1 {
		t.Errorf("type changé : %v", b)
	}
	// v2 durcit additionalProperties → rupture.
	if b := InputBackwardCompatible(v1, schema(map[string]string{"id": "string", "amt": "number"}, []string{"id"}, false)); len(b) != 1 {
		t.Errorf("additionalProperties durci : %v", b)
	}
	// v2 supprime une propriété ET interdit les additionnelles → rupture.
	if b := InputBackwardCompatible(schema(map[string]string{"id": "string", "opt": "string"}, []string{"id"}, false), schema(map[string]string{"id": "string"}, []string{"id"}, false)); len(b) != 1 {
		t.Errorf("propriété supprimée sous strict : %v", b)
	}
	// v2 supprime une propriété mais accepte les additionnelles → pas de rupture.
	if b := InputBackwardCompatible(schema(map[string]string{"id": "string", "opt": "string"}, []string{"id"}, true), schema(map[string]string{"id": "string"}, []string{"id"}, true)); len(b) != 0 {
		t.Errorf("suppression sous additionalProperties=true ne doit pas casser : %v", b)
	}
	// v1 acceptait null, v2 non → rupture d'entrée.
	v1null := schema(map[string]string{"id": "string"}, []string{"id"}, true)
	v1null.Nullable["id"] = true
	if b := InputBackwardCompatible(v1null, schema(map[string]string{"id": "string"}, []string{"id"}, true)); len(b) != 1 {
		t.Errorf("nullabilité d'entrée retirée : %v", b)
	}
}

func TestOutputBackwardCompatible(t *testing.T) {
	v1 := schema(map[string]string{"score": "integer", "ratio": "number"}, nil, true)
	if b := OutputBackwardCompatible(v1, v1); len(b) != 0 {
		t.Errorf("auto-compat : %v", b)
	}
	// v2 retire un champ → rupture.
	if b := OutputBackwardCompatible(v1, schema(map[string]string{"score": "integer"}, nil, true)); len(b) != 1 {
		t.Errorf("champ retiré : %v", b)
	}
	// v2 change un type → rupture.
	if b := OutputBackwardCompatible(v1, schema(map[string]string{"score": "string", "ratio": "number"}, nil, true)); len(b) != 1 {
		t.Errorf("type changé : %v", b)
	}
	// v2 ajoute un champ (surensemble) → OK pour la sortie.
	if b := OutputBackwardCompatible(v1, schema(map[string]string{"score": "integer", "ratio": "number", "extra": "string"}, nil, true)); len(b) != 0 {
		t.Errorf("surensemble devrait être compatible : %v", b)
	}
	// v2 rend un champ nullable que v1 ne l'était pas → rupture de sortie.
	v2null := schema(map[string]string{"score": "integer", "ratio": "number"}, nil, true)
	v2null.Nullable["score"] = true
	if b := OutputBackwardCompatible(v1, v2null); len(b) != 1 {
		t.Errorf("nullabilité de sortie ajoutée : %v", b)
	}
}
