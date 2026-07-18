package mcpserver

import (
	_ "embed"
	"errors"
	"fmt"

	"gopkg.in/yaml.v3"
)

//go:embed registry.yaml
var registryYAML []byte

// ServiceEntry décrit un serveur MCP comme tiers TIC (ADR-0009) : SLO, version
// de contrat et énum d'erreurs exposée.
type ServiceEntry struct {
	Name            string   `yaml:"name"`
	Tier            string   `yaml:"tier"`
	ContractVersion string   `yaml:"contractVersion"`
	SLOP95Ms        int      `yaml:"sloP95Ms"`
	Tools           []string `yaml:"tools"`
	Errors          []string `yaml:"errors"`
}

// Registry est le registre des serveurs MCP.
type Registry struct {
	Servers []ServiceEntry `yaml:"servers"`
}

// LoadRegistry parse et valide un registre YAML.
func LoadRegistry(data []byte) (*Registry, error) {
	var r Registry
	if err := yaml.Unmarshal(data, &r); err != nil {
		return nil, err
	}
	if err := r.Validate(); err != nil {
		return nil, err
	}
	return &r, nil
}

// DefaultRegistry retourne le registre embarqué (registry.yaml).
func DefaultRegistry() (*Registry, error) { return LoadRegistry(registryYAML) }

// Validate vérifie que chaque entrée porte un nom, une version de contrat, au
// moins un outil et une énum d'erreurs non vide.
func (r *Registry) Validate() error {
	if len(r.Servers) == 0 {
		return errors.New("registry: aucun serveur")
	}
	for _, s := range r.Servers {
		switch {
		case s.Name == "":
			return errors.New("registry: nom de serveur vide")
		case s.ContractVersion == "":
			return fmt.Errorf("registry: %s sans version de contrat", s.Name)
		case len(s.Tools) == 0:
			return fmt.Errorf("registry: %s sans outil", s.Name)
		case len(s.Errors) == 0:
			return fmt.Errorf("registry: %s sans énum d'erreurs", s.Name)
		}
	}
	return nil
}

// Server retourne l'entrée du serveur nommé (ok=false si absent).
func (r *Registry) Server(name string) (ServiceEntry, bool) {
	for _, s := range r.Servers {
		if s.Name == name {
			return s, true
		}
	}
	return ServiceEntry{}, false
}
