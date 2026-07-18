// Package mcpcontract vérifie la compatibilité **BACKWARD** des contrats
// d'outils MCP (§12.2, G-4) : une v2 doit accepter toute entrée valide selon
// v1, et ses sorties rester lisibles par un consommateur v1.
//
// Le modèle de schéma suit propriétés typées + nullabilité + champs requis +
// additionalProperties : suffisant pour les schémas inférés par le SDK.
package mcpcontract

import (
	"encoding/json"
	"sort"
)

// Schema est une vue simplifiée d'un JSON Schema objet.
type Schema struct {
	Properties           map[string]string // nom → type de base
	Nullable             map[string]bool   // nom → accepte null (union avec "null")
	Required             map[string]bool
	AdditionalProperties bool // false = strict (aucune propriété additionnelle)
}

// Parse extrait un Schema d'un JSON Schema objet (issu de l'inférence du SDK).
func Parse(raw []byte) (Schema, error) {
	var doc struct {
		Properties map[string]struct {
			Type json.RawMessage `json:"type"`
		} `json:"properties"`
		Required             []string `json:"required"`
		AdditionalProperties *bool    `json:"additionalProperties"`
	}
	if err := json.Unmarshal(raw, &doc); err != nil {
		return Schema{}, err
	}
	s := Schema{
		Properties:           make(map[string]string, len(doc.Properties)),
		Nullable:             make(map[string]bool),
		Required:             make(map[string]bool, len(doc.Required)),
		AdditionalProperties: true, // défaut JSON Schema si absent
	}
	for name, p := range doc.Properties {
		base, nullable := parseType(p.Type)
		s.Properties[name] = base
		if nullable {
			s.Nullable[name] = true
		}
	}
	for _, r := range doc.Required {
		s.Required[r] = true
	}
	if doc.AdditionalProperties != nil {
		s.AdditionalProperties = *doc.AdditionalProperties
	}
	return s, nil
}

// parseType canonicalise le champ "type", qui peut être une chaîne ("string")
// ou une union (["string","null"]) : retourne le premier type non-"null" et si
// l'union admet null. Absent ou non reconnu → ("", false).
func parseType(raw json.RawMessage) (base string, nullable bool) {
	if len(raw) == 0 {
		return "", false
	}
	var single string
	if json.Unmarshal(raw, &single) == nil {
		return single, false
	}
	var arr []string
	if json.Unmarshal(raw, &arr) == nil {
		for _, tp := range arr {
			if tp == "null" {
				nullable = true
				continue
			}
			if base == "" {
				base = tp
			}
		}
		return base, nullable
	}
	return "", false
}

// InputBackwardCompatible vérifie que v2 accepte toute entrée valide selon v1.
// Retourne la liste triée des ruptures (vide = compatible).
func InputBackwardCompatible(v1, v2 Schema) []string {
	var breaks []string
	// v2 ne doit pas exiger un champ qui n'était pas requis en v1.
	for f := range v2.Required {
		if !v1.Required[f] {
			breaks = append(breaks, "v2 exige un nouveau champ requis: "+f)
		}
	}
	for f, t1 := range v1.Properties {
		t2, ok := v2.Properties[f]
		switch {
		case !ok:
			// Une propriété v1 supprimée n'est un problème que si v2 interdit
			// les propriétés additionnelles (sinon v2 l'accepte encore).
			if !v2.AdditionalProperties {
				breaks = append(breaks, "v2 supprime la propriété d'entrée "+f+" et interdit les additionnelles")
			}
		case t2 != t1:
			breaks = append(breaks, "type d'entrée changé pour "+f+": "+t1+" -> "+t2)
		case v1.Nullable[f] && !v2.Nullable[f]:
			breaks = append(breaks, "v2 n'accepte plus null pour l'entrée "+f)
		}
	}
	// v2 ne doit pas interdire les propriétés additionnelles que v1 acceptait.
	if v1.AdditionalProperties && !v2.AdditionalProperties {
		breaks = append(breaks, "v2 interdit les proprietes additionnelles acceptees par v1")
	}
	sort.Strings(breaks)
	return breaks
}

// OutputBackwardCompatible vérifie qu'un consommateur v1 peut lire une sortie v2 :
// tout champ attendu par v1 doit exister en v2, du même type, et v2 ne doit pas
// introduire de nullabilité inattendue.
func OutputBackwardCompatible(v1, v2 Schema) []string {
	var breaks []string
	for f, t1 := range v1.Properties {
		t2, ok := v2.Properties[f]
		switch {
		case !ok:
			breaks = append(breaks, "champ de sortie retiré en v2: "+f)
		case t2 != t1:
			breaks = append(breaks, "type de sortie changé pour "+f+": "+t1+" -> "+t2)
		case v2.Nullable[f] && !v1.Nullable[f]:
			breaks = append(breaks, "v2 rend nullable la sortie "+f+" (inattendu pour un consommateur v1)")
		}
	}
	sort.Strings(breaks)
	return breaks
}
