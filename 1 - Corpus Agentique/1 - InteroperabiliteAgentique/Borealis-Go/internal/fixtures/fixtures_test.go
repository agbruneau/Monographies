package fixtures

import "testing"

func TestGenerateDeterministic(t *testing.T) {
	a := Generate(DefaultSeed, 100)
	b := Generate(DefaultSeed, 100)
	if len(a.Borrowers) != 100 || len(b.Borrowers) != 100 {
		t.Fatalf("comptes %d/%d, want 100", len(a.Borrowers), len(b.Borrowers))
	}
	for i := range a.Borrowers {
		if a.Borrowers[i] != b.Borrowers[i] {
			t.Fatalf("non déterministe à %d : %+v != %+v", i, a.Borrowers[i], b.Borrowers[i])
		}
	}
	if len(a.Policies) != 50 {
		t.Errorf("policies = %d, want 50", len(a.Policies))
	}
}

func TestBorrowersSortedAndLookup(t *testing.T) {
	ds := Generate(DefaultSeed, 50)
	for i := 1; i < len(ds.Borrowers); i++ {
		if ds.Borrowers[i-1].ID >= ds.Borrowers[i].ID {
			t.Fatalf("emprunteurs non triés à %d", i)
		}
	}
	if _, ok := ds.Borrower("A00001"); !ok {
		t.Error("A00001 absent")
	}
	if _, ok := ds.Borrower("ZZZ"); ok {
		t.Error("ZZZ devrait être absent")
	}
}

func TestPolicyFor(t *testing.T) {
	ds := Generate(DefaultSeed, 10)
	if _, ok := ds.PolicyFor("retail", 30); !ok {
		t.Error("politique retail/30 absente")
	}
	if _, ok := ds.PolicyFor("retail", 200); ok {
		t.Error("âge 200 ne devrait matcher aucune politique")
	}
	if _, ok := ds.PolicyFor("nonexistent", 30); ok {
		t.Error("segment inexistant ne devrait pas matcher")
	}
}
