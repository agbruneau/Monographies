package audit

import (
	"bytes"
	"encoding/json"
	"sort"
	"sync"
	"testing"
	"time"
)

func fixedClock() func() time.Time {
	ts := time.Date(2026, 7, 7, 12, 0, 0, 0, time.UTC)
	return func() time.Time { return ts }
}

func TestAppendAssignsSeqAndTime(t *testing.T) {
	j := New(WithClock(fixedClock()))
	e1 := j.Append(Record{KYA: "agent:kyc", OwnerKYH: "human:lead", KYC: "cust:hash", Action: "verify", Result: "ok", Version: "v1"})
	e2 := j.Append(Record{KYA: "agent:kyc", Action: "score", Result: "ok", Version: "v1"})
	if e1.Seq != 1 || e2.Seq != 2 {
		t.Errorf("séquences = %d,%d ; want 1,2", e1.Seq, e2.Seq)
	}
	if !e1.Time.Equal(time.Date(2026, 7, 7, 12, 0, 0, 0, time.UTC)) {
		t.Errorf("horodatage non injecté : %v", e1.Time)
	}
	if e1.KYA != "agent:kyc" || e1.OwnerKYH != "human:lead" {
		t.Errorf("champs perdus : %+v", e1)
	}
	if j.Len() != 2 {
		t.Errorf("Len = %d, want 2", j.Len())
	}
}

func TestEntriesReturnsCopy(t *testing.T) {
	j := New()
	j.Append(Record{Action: "a"})
	got := j.Entries()
	got[0].Action = "muté"
	if j.Entries()[0].Action != "a" {
		t.Error("Entries() ne retourne pas une copie défensive")
	}
}

func TestVerifyChainIntegrity(t *testing.T) {
	j := New(WithClock(fixedClock()))
	for i := 0; i < 5; i++ {
		j.Append(Record{KYA: "agent:kyc", OwnerKYH: "human:lead", KYC: "cust:h", Action: "verify", Result: "ok", Version: "v1"})
	}
	if err := j.Verify(); err != nil {
		t.Fatalf("chaîne intègre rejetée : %v", err)
	}
	// Les hachages sont chaînés et non vides.
	es := j.Entries()
	if es[0].PrevHash != "genesis" {
		t.Errorf("première entrée prev_hash = %q, want genesis", es[0].PrevHash)
	}
	for i := 1; i < len(es); i++ {
		if es[i].PrevHash != es[i-1].EntryHash {
			t.Errorf("entrée %d : prev_hash ne chaîne pas la précédente", i)
		}
	}
}

// TestHashEntryInjectiveOnNULBoundary : l'encodage haché doit être injectif
// même quand un champ contient l'octet séparateur — sinon deux Records
// distincts obtenus en déplaçant la frontière entre champs adjacents partagent
// le même entry_hash, et une telle altération d'un export passerait
// VerifyEntries/VerifyExport (FR-19, G-7). Même défense que webhook.mac
// (préfixe de longueur).
func TestHashEntryInjectiveOnNULBoundary(t *testing.T) {
	ts := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	e1 := Entry{Seq: 1, Time: ts, Record: Record{Action: "a", Result: "b"}, PrevHash: genesisHash}
	e2 := Entry{Seq: 1, Time: ts, Record: Record{Action: "a\x00b", Result: ""}, PrevHash: genesisHash}
	if hashEntry(e1) == hashEntry(e2) {
		t.Error("collision de frontière : deux Records distincts partagent le même entry_hash (encodage non injectif)")
	}
}

func TestVerifyDetectsTampering(t *testing.T) {
	// Réécriture d'un champ → chaîne rompue.
	j := New(WithClock(fixedClock()))
	for i := 0; i < 3; i++ {
		j.Append(Record{KYA: "a", Action: "x", Result: "ok", Version: "v1"})
	}
	j.entries[1].Result = "octroi_falsifie" // réécriture directe (attaque)
	if err := j.Verify(); err == nil {
		t.Error("réécriture d'un champ : chaîne devrait être détectée rompue (FR-19)")
	}

	// Réécriture du prev_hash → chaîne rompue.
	j2 := New(WithClock(fixedClock()))
	for i := 0; i < 3; i++ {
		j2.Append(Record{KYA: "a", Action: "x", Result: "ok", Version: "v1"})
	}
	j2.entries[2].PrevHash = "genesis"
	if err := j2.Verify(); err == nil {
		t.Error("prev_hash réécrit : chaîne devrait être détectée rompue")
	}
}

// TestVerifyExportAnchor : l'ancrage (compte + tête) referme l'angle mort de la
// relecture nue — troncature de fin et effacement total deviennent détectables.
func TestVerifyExportAnchor(t *testing.T) {
	j := New(WithClock(fixedClock()))
	for i := 0; i < 3; i++ {
		j.Append(Record{KYA: "a", Action: "x", Result: "ok", Version: "v1"})
	}
	entries, wantLen, wantTip := j.Entries(), j.Len(), j.Tip()

	// Chaîne intacte + ancrage correct → OK.
	if err := VerifyExport(entries, wantLen, wantTip); err != nil {
		t.Fatalf("chaîne intacte rejetée : %v", err)
	}
	// Troncature de fin : relecture nue passe (chaîne valide plus courte)...
	if err := VerifyEntries(entries[:2]); err != nil {
		t.Errorf("relecture nue d'une chaîne tronquée devrait passer (limite connue) : %v", err)
	}
	// ...mais l'ancrage la détecte (compte).
	if err := VerifyExport(entries[:2], wantLen, wantTip); err == nil {
		t.Error("troncature de fin non détectée par l'ancrage")
	}
	// Effacement total → détecté (compte 0).
	if err := VerifyExport(nil, wantLen, wantTip); err == nil {
		t.Error("effacement total non détecté par l'ancrage")
	}
	// Chaîne ENTIÈREMENT RE-FORGÉE : cohérente en interne (VerifyEntries passe),
	// **même longueur**, mais contenu — donc tête — différents. C'est la seule
	// attaque que ni le compte ni la relecture nue n'attrapent : elle isole le
	// contrôle de tête (retirer ce contrôle rendrait ce cas vert à tort).
	reforged := New(WithClock(fixedClock()))
	for i := 0; i < 3; i++ {
		reforged.Append(Record{KYA: "a", Action: "x", Result: "octroi_falsifie", Version: "v1"})
	}
	ref := reforged.Entries()
	if err := VerifyEntries(ref); err != nil {
		t.Fatalf("la chaîne re-forgée doit être interne-valide (sinon le cas ne teste pas la tête) : %v", err)
	}
	if len(ref) != wantLen || reforged.Tip() == wantTip {
		t.Fatalf("cas invalide : re-forge de même longueur et de tête différente attendu")
	}
	if err := VerifyExport(ref, wantLen, wantTip); err == nil {
		t.Error("re-forge cohérent de même longueur non détecté — le contrôle de tête ne sert à rien")
	}
	// Journal vide : Tip == genesis, compte 0, cohérent.
	empty := New(WithClock(fixedClock()))
	if err := VerifyExport(empty.Entries(), empty.Len(), empty.Tip()); err != nil {
		t.Errorf("journal vide cohérent rejeté : %v", err)
	}
}

// TestExportJSONLReplayable : le journal exporté en JSONL est **rejouable** —
// relu ligne à ligne, il reconstruit une chaîne intègre (VerifyEntries) et son
// ancrage (VerifyExport) tient (FR-23).
func TestExportJSONLReplayable(t *testing.T) {
	j := New(WithClock(fixedClock()))
	for i := 0; i < 3; i++ {
		j.Append(Record{KYA: "agent:kyc", KYC: "cust:h", Action: "decision", Result: "declined", Version: "m5"})
	}
	var buf bytes.Buffer
	if err := j.ExportJSONL(&buf); err != nil {
		t.Fatal(err)
	}
	if lines := bytes.Count(buf.Bytes(), []byte("\n")); lines != 3 {
		t.Errorf("JSONL = %d lignes, want 3", lines)
	}
	var replayed []Entry
	dec := json.NewDecoder(bytes.NewReader(buf.Bytes()))
	for dec.More() {
		var e Entry
		if err := dec.Decode(&e); err != nil {
			t.Fatal(err)
		}
		replayed = append(replayed, e)
	}
	if err := VerifyEntries(replayed); err != nil {
		t.Errorf("chaîne rejouée non intègre : %v", err)
	}
	if err := VerifyExport(replayed, j.Len(), j.Tip()); err != nil {
		t.Errorf("ancrage rejoué : %v", err)
	}
}

// failWriter échoue à la première écriture (couvre le chemin d'erreur d'export).
type failWriter struct{}

func (failWriter) Write([]byte) (int, error) { return 0, errWrite }

var errWrite = bytes.ErrTooLarge // erreur sentinelle réutilisée

func TestExportJSONLWriteError(t *testing.T) {
	j := New(WithClock(fixedClock()))
	j.Append(Record{KYA: "a", Action: "x", Result: "ok", Version: "v1"})
	if err := j.ExportJSONL(failWriter{}); err == nil {
		t.Error("un writer défaillant devrait propager l'erreur")
	}
	// Journal vide : rien à écrire → pas d'erreur.
	empty := New(WithClock(fixedClock()))
	if err := empty.ExportJSONL(failWriter{}); err != nil {
		t.Errorf("journal vide : %v, want nil", err)
	}
}

func TestAuditThroughput(t *testing.T) {
	// NFR-03 : P99 < 200 ms par écriture ET débit >= 1000 écritures/s.
	j := New()
	const n = 2000
	lat := make([]time.Duration, 0, n)
	start := time.Now()
	for i := 0; i < n; i++ {
		t0 := time.Now()
		j.Append(Record{KYA: "a", Action: "x", Result: "ok", Version: "v1"})
		lat = append(lat, time.Since(t0))
	}
	elapsed := time.Since(start)

	rate := float64(n) / elapsed.Seconds()
	if rate < 1000 {
		t.Errorf("débit = %.0f écritures/s, want >= 1000 (NFR-03)", rate)
	}
	sort.Slice(lat, func(i, j int) bool { return lat[i] < lat[j] })
	if p99 := lat[int(0.99*float64(len(lat)))]; p99 > 200*time.Millisecond {
		t.Errorf("P99 = %v, want < 200ms (NFR-03)", p99)
	}
}

func BenchmarkAppend(b *testing.B) {
	j := New()
	r := Record{KYA: "a", Action: "x", Result: "ok", Version: "v1"}
	for i := 0; i < b.N; i++ {
		j.Append(r)
	}
}

func TestAppendConcurrent(t *testing.T) {
	j := New(WithClock(fixedClock()))
	const n = 100
	var wg sync.WaitGroup
	wg.Add(n)
	for range n {
		go func() {
			defer wg.Done()
			j.Append(Record{Action: "x"})
		}()
	}
	wg.Wait()
	if j.Len() != n {
		t.Fatalf("Len = %d, want %d", j.Len(), n)
	}
	seen := make(map[uint64]bool, n)
	for _, e := range j.Entries() {
		if seen[e.Seq] {
			t.Fatalf("séquence dupliquée : %d", e.Seq)
		}
		seen[e.Seq] = true
	}
	if len(seen) != n {
		t.Errorf("séquences uniques = %d, want %d", len(seen), n)
	}
}
