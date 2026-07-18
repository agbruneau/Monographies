package a2a

import (
	"testing"

	a2asdk "github.com/a2aproject/a2a-go/v2/a2a"
)

func TestTaskMachineHappyPath(t *testing.T) {
	m := NewTaskMachine()
	if m.State() != a2asdk.TaskStateSubmitted {
		t.Fatalf("état initial = %s, want submitted", m.State())
	}
	for _, s := range []a2asdk.TaskState{a2asdk.TaskStateWorking, a2asdk.TaskStateCompleted} {
		if err := m.To(s); err != nil {
			t.Fatalf("transition vers %s : %v", s, err)
		}
	}
	if m.State() != a2asdk.TaskStateCompleted {
		t.Errorf("état final = %s, want completed", m.State())
	}
}

func TestTerminalIrreversible(t *testing.T) {
	// completed → working refusé.
	m := NewTaskMachine()
	_ = m.To(a2asdk.TaskStateWorking)
	_ = m.To(a2asdk.TaskStateCompleted)
	if err := m.To(a2asdk.TaskStateWorking); err == nil {
		t.Error("completed → working devrait être refusé (terminal irréversible)")
	}
	// failed ↛ completed.
	m2 := NewTaskMachine()
	_ = m2.To(a2asdk.TaskStateFailed)
	if err := m2.To(a2asdk.TaskStateCompleted); err == nil {
		t.Error("failed → completed devrait être refusé")
	}
	// input-required (non terminal) → working permis.
	m3 := NewTaskMachine()
	_ = m3.To(a2asdk.TaskStateWorking)
	if err := m3.To(a2asdk.TaskStateInputRequired); err != nil {
		t.Fatal(err)
	}
	if err := m3.To(a2asdk.TaskStateWorking); err != nil {
		t.Errorf("input-required → working devrait être permis : %v", err)
	}
}

func TestCanTransitionAndIsTerminal(t *testing.T) {
	if CanTransition(a2asdk.TaskStateWorking, a2asdk.TaskStateUnspecified) {
		t.Error("transition vers unspecified devrait être refusée")
	}
	for _, s := range []a2asdk.TaskState{a2asdk.TaskStateCompleted, a2asdk.TaskStateFailed, a2asdk.TaskStateCanceled, a2asdk.TaskStateRejected} {
		if !IsTerminal(s) {
			t.Errorf("%s devrait être terminal", s)
		}
	}
	for _, s := range []a2asdk.TaskState{a2asdk.TaskStateSubmitted, a2asdk.TaskStateWorking, a2asdk.TaskStateInputRequired} {
		if IsTerminal(s) {
			t.Errorf("%s ne devrait pas être terminal", s)
		}
	}
}
