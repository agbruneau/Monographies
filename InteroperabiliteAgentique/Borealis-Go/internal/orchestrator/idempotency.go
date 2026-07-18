package orchestrator

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sync"

	a2asdk "github.com/a2aproject/a2a-go/v2/a2a"
	"github.com/agbruneau/borealis/internal/audit"
)

// delegationCache mémorise le premier résultat de chaque délégation par clé
// d'idempotence, pour qu'un rejeu ne produise qu'un seul effet (FR-18, NFR-05).
//
// ponytail: seen croît d'une entrée par délégation unique, jamais purgée —
// borné par le nombre de délégations distinctes sur la durée de vie du
// processus, pas par le temps. Plafond acceptable pour le démonstrateur
// (processus courts) ; chemin d'évolution si un pilote tourne longtemps :
// LRU/TTL par clé, ou purge alignée sur la clôture du dossier (contextId).
type delegationCache struct {
	mu   sync.Mutex
	seen map[string]a2asdk.SendMessageResult
}

// WithAudit branche un journal : chaque délégation dédupliquée est journalisée
// (« delegation »), et tout rejeu signalé (« delegation.duplicate »). Sans
// journal, la déduplication reste active mais silencieuse.
func (o *Orchestrator) WithAudit(j *audit.Journal) *Orchestrator {
	o.journal = j
	return o
}

// idempotencyKey dérive la clé de déduplication d'une délégation à partir du
// **contextId (dossier)**, de la **compétence**, et du **contenu du message
// (params)** — jamais du messageId, aléatoire par appel (sinon aucun rejeu ne
// serait reconnu). Les Parts sont un tableau ordonné : leur sérialisation JSON
// est déterministe.
func idempotencyKey(skill string, msg *a2asdk.Message) (string, error) {
	parts, err := json.Marshal(msg.Parts)
	if err != nil {
		return "", fmt.Errorf("orchestrator: clé d'idempotence: %w", err)
	}
	h := sha256.New()
	fmt.Fprintf(h, "%s\x00%s\x00", msg.ContextID, skill)
	h.Write(parts)
	return hex.EncodeToString(h.Sum(nil)), nil
}

// idempotent exécute deliver **au plus une fois par clé** quand un contextId est
// présent : un rejeu retourne le résultat mémorisé et journalise le doublon.
// Sans contextId, aucune identité de dossier n'asseoit l'idempotence → exécution
// directe. Un échec n'est jamais mémorisé (un rejeu peut réussir).
//
// Plafond « à-moins-une-fois » : si un effet est appliqué côté pair PUIS que la
// réponse échoue en transit (flux SSE rompu après émission), deliver retourne une
// erreur, non mémorisée → un rejeu ré-exécute et double l'effet. La garantie « un
// seul effet » (FR-18) ne tient donc strictement que pour un rejeu APRÈS succès
// propre ; l'idempotence des effets partiels exigerait une clé d'idempotence
// honorée côté serveur (hors périmètre du démonstrateur, agents-stub sans effet).
//
// ponytail: dédup par map sérialisée ; deux délégations identiques *concurrentes*
// peuvent toutes deux s'exécuter (le rejeu visé est séquentiel). Verrou par clé
// (singleflight) si un pilote concurrent l'exige.
func (o *Orchestrator) idempotent(contextID, skill string, msg *a2asdk.Message, deliver func() (a2asdk.SendMessageResult, error)) (a2asdk.SendMessageResult, error) {
	if contextID == "" {
		return deliver()
	}
	key, err := idempotencyKey(skill, msg)
	if err != nil {
		return deliver() // clé indérivable → dégradé sûr : pas de dédup
	}

	o.cache.mu.Lock()
	if cached, ok := o.cache.seen[key]; ok {
		o.cache.mu.Unlock()
		o.auditDelegation(contextID, "delegation.duplicate", skill)
		return cached, nil
	}
	o.cache.mu.Unlock()

	res, err := deliver()
	if err != nil {
		return res, err // échec non mémorisé
	}
	o.cache.mu.Lock()
	o.cache.seen[key] = res
	o.cache.mu.Unlock()
	o.auditDelegation(contextID, "delegation", skill)
	return res, nil
}

// auditDelegation journalise une délégation (ou son doublon). No-op sans journal.
// Le contextId (identifiant de dossier synthétique) n'est pas une PII.
func (o *Orchestrator) auditDelegation(contextID, action, skill string) {
	if o.journal != nil {
		o.journal.Append(audit.Record{KYA: "agent:orchestrator", KYC: contextID, Action: action, Result: skill, Version: "m4"})
	}
}
