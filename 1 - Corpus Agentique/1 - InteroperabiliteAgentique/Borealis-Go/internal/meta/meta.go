// Package meta expose les informations de version du démonstrateur Borealis.
// Les valeurs sont injectées au link par les ldflags du Makefile (voir build/).
package meta

import "fmt"

// Valeurs par défaut ; surchargées au link (-X ...).
var (
	Version   = "dev"
	Commit    = "unknown"
	BuildDate = "unknown"
)

// Info retourne une chaîne compacte version/commit/date, utilisable par le
// drapeau --version de chaque binaire.
func Info() string {
	return fmt.Sprintf("borealis %s (commit %s, built %s)", Version, Commit, BuildDate)
}
