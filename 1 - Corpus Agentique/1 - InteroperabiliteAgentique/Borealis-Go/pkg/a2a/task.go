package a2a

import (
	"fmt"

	a2asdk "github.com/a2aproject/a2a-go/v2/a2a"
)

var terminalStates = map[a2asdk.TaskState]bool{
	a2asdk.TaskStateCompleted: true,
	a2asdk.TaskStateFailed:    true,
	a2asdk.TaskStateCanceled:  true,
	a2asdk.TaskStateRejected:  true,
}

// IsTerminal indique si un état de Task est terminal (irréversible).
func IsTerminal(s a2asdk.TaskState) bool { return terminalStates[s] }

// CanTransition indique si la transition from→to est permise : les états
// terminaux sont **irréversibles** (FR-07, G-2) et l'état non spécifié n'est
// jamais une cible valide.
func CanTransition(from, to a2asdk.TaskState) bool {
	if IsTerminal(from) {
		return false
	}
	return to != a2asdk.TaskStateUnspecified
}

// TaskMachine suit l'état d'une Task et refuse les transitions interdites.
type TaskMachine struct {
	state a2asdk.TaskState
}

// NewTaskMachine crée une machine à l'état initial submitted.
func NewTaskMachine() *TaskMachine {
	return &TaskMachine{state: a2asdk.TaskStateSubmitted}
}

// State retourne l'état courant.
func (m *TaskMachine) State() a2asdk.TaskState { return m.state }

// To effectue une transition ; retourne une erreur si elle est interdite (état
// terminal déjà atteint, ou cible non spécifiée).
func (m *TaskMachine) To(s a2asdk.TaskState) error {
	if !CanTransition(m.state, s) {
		return fmt.Errorf("a2a: transition interdite %s -> %s", m.state, s)
	}
	m.state = s
	return nil
}
