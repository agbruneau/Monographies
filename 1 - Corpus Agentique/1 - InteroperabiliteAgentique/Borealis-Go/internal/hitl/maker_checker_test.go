package hitl

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"strings"
	"testing"
)

// selectsField déballe les enveloppes d'accès (index, déréférence, parenthèses)
// et dit si l'expression désigne le champ `field` d'un récepteur. Couvre
// `b.commands`, `b.commands[i]`, `(*b).commands`.
func selectsField(expr ast.Expr, field string) bool {
	for {
		switch x := expr.(type) {
		case *ast.ParenExpr:
			expr = x.X
		case *ast.IndexExpr:
			expr = x.X
		case *ast.StarExpr:
			expr = x.X
		case *ast.SelectorExpr:
			return x.Sel.Name == field
		default:
			return false
		}
	}
}

// commandQueueMutators retourne l'ensemble des fonctions qui **peuvent** muter le
// champ `field` : affectation directe/indexée, ou **prise d'adresse** (`&b.field`)
// — car une adresse échappée permet une mutation hors du site (helper, closure).
// C'est le cœur du test topologique maker-checker : durci contre l'indirection
// par pointeur, l'angle mort d'une simple détection d'affectation.
func commandQueueMutators(files []*ast.File, field string) map[string]bool {
	mutators := map[string]bool{}
	for _, file := range files {
		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Body == nil { // corps nil (déclaration externe) → rien à inspecter
				continue
			}
			ast.Inspect(fn.Body, func(n ast.Node) bool {
				switch x := n.(type) {
				case *ast.AssignStmt:
					for _, lhs := range x.Lhs {
						if selectsField(lhs, field) {
							mutators[fn.Name.Name] = true
						}
					}
				case *ast.UnaryExpr:
					if x.Op == token.AND && selectsField(x.X, field) {
						mutators[fn.Name.Name] = true // &b.field : l'adresse échappe
					}
				}
				return true
			})
		}
	}
	return mutators
}

// TestMakerCheckerTopology est un test **topologique** (FR-24, M3.7) : il prouve,
// par analyse du code source du paquet, que la file de commande n'est mutée que
// par Release — le point *checker*. Un contournement futur (une méthode qui
// écrirait `commands`, l'indexerait, ou en prendrait l'adresse sans passer par
// une release humaine) casse la séparation maker-checker et doit échouer ici,
// pas silencieusement à l'exécution.
//
// Complément du test comportemental TestProposeThenReleaseFlow : celui-là montre
// que le flux courant respecte l'invariant ; celui-ci interdit d'en écrire un
// autre.
func TestMakerCheckerTopology(t *testing.T) {
	const queueField = "commands"

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, ".", func(fi fs.FileInfo) bool {
		return !strings.HasSuffix(fi.Name(), "_test.go") // code de production seul
	}, 0)
	if err != nil {
		t.Fatalf("analyse du paquet : %v", err)
	}

	var files []*ast.File
	for _, pkg := range pkgs {
		for _, f := range pkg.Files {
			files = append(files, f)
		}
	}
	mutators := commandQueueMutators(files, queueField)

	if len(mutators) == 0 {
		t.Fatalf("aucun mutateur de %q trouvé : le champ a-t-il été renommé ? (test à réviser)", queueField)
	}
	for name := range mutators {
		if name != "Release" {
			t.Errorf("la file de commande est mutée par %q : seul Release (checker) doit y écrire (FR-24)", name)
		}
	}
}

// TestMakerCheckerHeuristicCatchesEscape prouve que l'heuristique n'est PAS
// aveugle à l'indirection par pointeur : un contournement furtif (helper mutant
// `*dst` alimenté par `&b.commands`) est bien détecté au site d'appel. Sans ce
// durcissement, `sneak` passerait inaperçu — le test topologique deviendrait
// tautologique face à ce vecteur.
func TestMakerCheckerHeuristicCatchesEscape(t *testing.T) {
	const src = `package fake
type Broker struct{ commands []int }
func enqueue(dst *[]int, c int) { *dst = append(*dst, c) }
func (b *Broker) Release()      { b.commands = append(b.commands, 1) }
func (b *Broker) sneak()        { enqueue(&b.commands, 9) }
func (b *Broker) touch()        { b.commands[0] = 7 }
func (b *Broker) read() int     { return len(b.commands) }
`
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "fake.go", src, 0)
	if err != nil {
		t.Fatalf("analyse de l'extrait : %v", err)
	}
	got := commandQueueMutators([]*ast.File{f}, "commands")

	for _, want := range []string{"Release", "sneak", "touch"} {
		if !got[want] {
			t.Errorf("mutateur %q non détecté (angle mort de l'heuristique)", want)
		}
	}
	if got["read"] {
		t.Error("read() ne mute pas commands : faux positif")
	}
}
