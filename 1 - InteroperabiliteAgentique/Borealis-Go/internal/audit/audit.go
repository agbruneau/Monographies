// Package audit fournit le journal d'audit WORM du démonstrateur : journal
// **append-only** (aucune mise à jour ni suppression) portant les 7 champs de
// la piste (KYA/KYH/KYC/ts/action/résultat/version), scellés dans une **chaîne
// de hachage vérifiable** (entry_hash/prev_hash, formule PRD §9.4). Verify()
// détecte toute réécriture (FR-19, FR-20, G-7).
//
// Ce journal est **distinct** du puits d'observabilité OTel (invariant 7 /
// ADR-0003) : ne jamais confondre les deux.
package audit

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"sync"
	"time"
)

// genesisHash est le prev_hash de la première entrée (ancre de la chaîne).
const genesisHash = "genesis"

// Record porte les champs sémantiques fournis par l'appelant (triade
// d'identité KYA/KYH/KYC + action, résultat, version de contrat/politique).
type Record struct {
	KYA      string // Know-Your-Agent : identité de l'agent non-humain
	OwnerKYH string // Know-Your-Human : humain responsable de l'agent
	KYC      string // Know-Your-Customer : sujet (pseudonymisé, jamais de PII)
	Action   string
	Result   string
	Version  string
}

// Entry est une entrée scellée : le Record plus un numéro de séquence, un
// horodatage, et la chaîne de hachage (prev_hash chaîné, entry_hash = sha256
// des champs — FR-19, formule PRD §9.4) attribués par le journal.
type Entry struct {
	Seq  uint64
	Time time.Time
	Record
	PrevHash  string
	EntryHash string
}

// Journal est un journal append-only sûr en concurrence (en mémoire pour M0).
//
// ponytail: entries croît sans borne pour la durée de vie du processus — un
// binaire tournant en service long (agents/orchestrateur en compose) accumule
// indéfiniment. Chemin d'évolution si le débit/la durée l'exigent : un puits
// JSONL (append-only sur disque) derrière la même interface Append/Entries,
// ou une purge/rotation bornée par taille — pas avant d'en avoir la preuve
// (le démonstrateur tourne en processus courts, e2e/tests).
type Journal struct {
	mu      sync.Mutex
	clock   func() time.Time
	entries []Entry
}

// Option configure un Journal.
type Option func(*Journal)

// WithClock injecte une horloge (déterminisme des tests et du mode e2e).
func WithClock(c func() time.Time) Option {
	return func(j *Journal) { j.clock = c }
}

// New construit un journal (horloge par défaut : time.Now).
func New(opts ...Option) *Journal {
	j := &Journal{clock: time.Now}
	for _, o := range opts {
		o(j)
	}
	return j
}

// Append scelle un enregistrement dans la chaîne de hachage et retourne l'entrée
// créée. Seule opération de mutation : jamais de mise à jour ni de suppression
// (WORM). entry_hash chaîne l'entrée à la précédente (FR-19).
func (j *Journal) Append(r Record) Entry {
	j.mu.Lock()
	defer j.mu.Unlock()
	prevHash := genesisHash
	if n := len(j.entries); n > 0 {
		prevHash = j.entries[n-1].EntryHash
	}
	e := Entry{Seq: uint64(len(j.entries)) + 1, Time: j.clock(), Record: r, PrevHash: prevHash}
	e.EntryHash = hashEntry(e)
	j.entries = append(j.entries, e)
	return e
}

// hashEntry calcule l'entry_hash = sha256 des 7 champs (KYA/KYH/KYC/ts/action/
// résultat/version) + seq + prev_hash (champs de la formule PRD §9.4, restreinte
// aux 7 champs du journal). Encodage **injectif** : chaque champ variable est
// préfixé de sa longueur — un simple séparateur laisserait deux Records
// distincts (frontière déplacée autour d'un octet NUL) partager le même hash,
// altération indétectable par Verify* (FR-19, G-7 ; même défense que
// webhook.mac ; ADR-0012). Horodatage canonique UTC.
func hashEntry(e Entry) string {
	h := sha256.New()
	fmt.Fprintf(h, "%d\x00%s", e.Seq, e.Time.UTC().Format(time.RFC3339Nano))
	for _, f := range []string{e.KYA, e.OwnerKYH, e.KYC, e.Action, e.Result, e.Version, e.PrevHash} {
		fmt.Fprintf(h, "\x00%d:%s", len(f), f)
	}
	return hex.EncodeToString(h.Sum(nil))
}

// Verify rejoue la chaîne de hachage interne. Voir VerifyEntries pour la portée
// exacte (et ses limites) de la relecture (FR-19, G-7).
func (j *Journal) Verify() error {
	j.mu.Lock()
	defer j.mu.Unlock()
	return VerifyEntries(j.entries)
}

// Tip retourne le hash de la dernière entrée (genesisHash si la chaîne est vide).
// C'est l'**ancrage** à engager hors du journal (reçu signé, registre externe) :
// il fixe la queue de la chaîne, que la relecture seule ne peut pas sceller.
func (j *Journal) Tip() string {
	j.mu.Lock()
	defer j.mu.Unlock()
	if n := len(j.entries); n > 0 {
		return j.entries[n-1].EntryHash
	}
	return genesisHash
}

// VerifyEntries rejoue une chaîne d'entrées **détenue à l'extérieur** du journal
// (log exporté, re-vérifié par un auditeur tiers). Elle détecte toute mutation de
// champ, réordonnancement, insertion, ou retrait **en milieu ou en tête** de
// chaîne (prev_hash rompu).
//
// Limite intrinsèque d'une chaîne de hachage nue : la relecture seule NE détecte
// PAS une **troncature de fin** (retrait des N dernières entrées) ni un log
// **entièrement effacé** — la queue tronquée reste une chaîne valide plus courte.
// Sceller la queue exige un état externe : voir VerifyExport, à préférer pour
// re-vérifier un log exporté (FR-19, G-7).
func VerifyEntries(entries []Entry) error {
	prevHash := genesisHash
	for _, e := range entries {
		if e.PrevHash != prevHash {
			return fmt.Errorf("audit: chaîne rompue à seq %d (prev_hash incohérent)", e.Seq)
		}
		if want := hashEntry(e); e.EntryHash != want {
			return fmt.Errorf("audit: entrée %d altérée (entry_hash incohérent)", e.Seq)
		}
		prevHash = e.EntryHash
	}
	return nil
}

// VerifyExport re-vérifie un log exporté contre un **ancrage engagé à part** au
// moment de l'export : le nombre d'entrées attendu et le hash de tête (Tip). Elle
// referme l'angle mort de VerifyEntries — une troncature de fin ou un effacement
// total changent le compte et/ou la tête, donc sont détectés. C'est la garantie
// tamper-evident complète pour un auditeur tiers (FR-19, G-7).
func VerifyExport(entries []Entry, expectedLen int, expectedTip string) error {
	if len(entries) != expectedLen {
		return fmt.Errorf("audit: %d entrées, %d attendues (troncature ou insertion ?)", len(entries), expectedLen)
	}
	if err := VerifyEntries(entries); err != nil {
		return err
	}
	gotTip := genesisHash
	if n := len(entries); n > 0 {
		gotTip = entries[n-1].EntryHash
	}
	if gotTip != expectedTip {
		return fmt.Errorf("audit: tête %q != attendue %q (queue altérée)", gotTip, expectedTip)
	}
	return nil
}

// Entries retourne une copie des entrées (l'appelant ne peut pas muter l'état).
func (j *Journal) Entries() []Entry {
	j.mu.Lock()
	defer j.mu.Unlock()
	out := make([]Entry, len(j.entries))
	copy(out, j.entries)
	return out
}

// ExportJSONL écrit le journal en **JSONL** (une entrée JSON par ligne), format
// rejouable (FR-23) : un auditeur relit les lignes, reconstruit la chaîne et la
// re-vérifie (VerifyEntries / VerifyExport avec l'ancrage), hors du processus.
func (j *Journal) ExportJSONL(w io.Writer) error {
	j.mu.Lock()
	defer j.mu.Unlock()
	enc := json.NewEncoder(w)
	for _, e := range j.entries {
		if err := enc.Encode(e); err != nil {
			return err
		}
	}
	return nil
}

// Len retourne le nombre d'entrées.
func (j *Journal) Len() int {
	j.mu.Lock()
	defer j.mu.Unlock()
	return len(j.entries)
}
