package a2aserver

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"iter"
	"net/http/httptest"
	"testing"

	a2asdk "github.com/a2aproject/a2a-go/v2/a2a"
	"github.com/a2aproject/a2a-go/v2/a2aclient"
	"github.com/a2aproject/a2a-go/v2/a2aclient/agentcard"
	"github.com/a2aproject/a2a-go/v2/a2asrv"
	borealisa2a "github.com/agbruneau/borealis/pkg/a2a"
)

func startAndConnect(t *testing.T, exec a2asrv.AgentExecutor) (context.Context, *a2aclient.Client) {
	t.Helper()
	ctx := context.Background()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	srv := httptest.NewUnstartedServer(nil)
	baseURL := "http://" + srv.Listener.Addr().String()
	card := BuildCard(Specs[1], baseURL)
	sig, err := borealisa2a.Sign(card, key, "k")
	if err != nil {
		t.Fatal(err)
	}
	card.Signatures = []a2asdk.AgentCardSignature{sig}
	srv.Config.Handler = Mux(card, exec)
	srv.Start()
	t.Cleanup(srv.Close)

	resolved, err := agentcard.DefaultResolver.Resolve(ctx, baseURL)
	if err != nil {
		t.Fatal(err)
	}
	client, err := a2aclient.NewFromCard(ctx, resolved, a2aclient.WithJSONRPCTransport(nil))
	if err != nil {
		t.Fatal(err)
	}
	return ctx, client
}

func hasState(states []a2asdk.TaskState, want a2asdk.TaskState) bool {
	for _, s := range states {
		if s == want {
			return true
		}
	}
	return false
}

// collectStates draine le flux et retourne les états observés (Task initiale +
// mises à jour), l'ensemble des contextId, et le nombre de TaskStatusUpdateEvent.
func collectStates(t *testing.T, seq iter.Seq2[a2asdk.Event, error]) (states []a2asdk.TaskState, ctxIDs map[string]bool, suCount int) {
	t.Helper()
	ctxIDs = map[string]bool{}
	for ev, err := range seq {
		if err != nil {
			t.Fatalf("stream : %v", err)
		}
		switch e := ev.(type) {
		case *a2asdk.TaskStatusUpdateEvent:
			suCount++
			states = append(states, e.Status.State)
			ctxIDs[e.ContextID] = true
		case *a2asdk.Task:
			states = append(states, e.Status.State)
			ctxIDs[e.ContextID] = true
		}
	}
	return states, ctxIDs, suCount
}

// TestStreamingLifecycle : le streaming SSE observe ≥ 2 TaskStatusUpdateEvent
// (FR-08) couvrant working → completed (FR-05), tous sous un même contextId non
// vide (contextId par dossier).
func TestStreamingLifecycle(t *testing.T) {
	ctx, client := startAndConnect(t, &LifecycleExecutor{Reply: "done"})
	seq := client.SendStreamingMessage(ctx, &a2asdk.SendMessageRequest{
		Message: a2asdk.NewMessage(a2asdk.MessageRoleUser, a2asdk.NewTextPart("go")),
	})
	states, ctxIDs, suCount := collectStates(t, seq)

	if suCount < 2 {
		t.Errorf("TaskStatusUpdateEvent = %d, want >= 2 (FR-08)", suCount)
	}
	if !hasState(states, a2asdk.TaskStateSubmitted) || !hasState(states, a2asdk.TaskStateWorking) || !hasState(states, a2asdk.TaskStateCompleted) {
		t.Errorf("états = %v, want submitted + working + completed (FR-05)", states)
	}
	if len(ctxIDs) != 1 || ctxIDs[""] {
		t.Errorf("contextIDs = %v, want un seul, non vide (contextId par dossier)", ctxIDs)
	}
}

// TestStreamingInputRequired : force et observe la transition working →
// input-required (chemin HITL, M2 exit).
func TestStreamingInputRequired(t *testing.T) {
	ctx, client := startAndConnect(t, &LifecycleExecutor{Reply: "need input", Final: a2asdk.TaskStateInputRequired})
	seq := client.SendStreamingMessage(ctx, &a2asdk.SendMessageRequest{
		Message: a2asdk.NewMessage(a2asdk.MessageRoleUser, a2asdk.NewTextPart("go")),
	})
	states, _, _ := collectStates(t, seq)
	iw, ir := indexOf(states, a2asdk.TaskStateWorking), indexOf(states, a2asdk.TaskStateInputRequired)
	if iw < 0 || ir < 0 || iw >= ir {
		t.Errorf("états = %v, want working AVANT input-required (transition observée)", states)
	}
}

func indexOf(states []a2asdk.TaskState, want a2asdk.TaskState) int {
	for i, s := range states {
		if s == want {
			return i
		}
	}
	return -1
}

// TestStreamingTerminalGuard : la garde d'états branchée sur l'exécuteur bloque
// une transition post-terminale (working après completed) — FR-07 appliqué sur
// le chemin de streaming réel.
func TestStreamingTerminalGuard(t *testing.T) {
	ctx, client := startAndConnect(t, &LifecycleExecutor{Reply: "done", IllegalAfterFinal: true})
	seq := client.SendStreamingMessage(ctx, &a2asdk.SendMessageRequest{
		Message: a2asdk.NewMessage(a2asdk.MessageRoleUser, a2asdk.NewTextPart("go")),
	})
	states, _, _ := collectStates(t, seq)

	working := 0
	for _, s := range states {
		if s == a2asdk.TaskStateWorking {
			working++
		}
	}
	if working != 1 {
		t.Errorf("working observé %d fois, want 1 (transition post-terminale bloquée)", working)
	}
	if last := states[len(states)-1]; !borealisa2a.IsTerminal(last) {
		t.Errorf("dernier état = %s, want terminal", last)
	}
}

// refEchoExecutor renvoie les referenceTaskIds reçus (FR-06).
type refEchoExecutor struct{}

func (refEchoExecutor) Execute(_ context.Context, execCtx *a2asrv.ExecutorContext) iter.Seq2[a2asdk.Event, error] {
	return func(yield func(a2asdk.Event, error) bool) {
		refs := "refs:"
		if execCtx.Message != nil {
			for _, id := range execCtx.Message.ReferenceTasks {
				refs += string(id) + ";"
			}
		}
		yield(a2asdk.NewMessage(a2asdk.MessageRoleAgent, a2asdk.NewTextPart(refs)), nil)
	}
}

func (refEchoExecutor) Cancel(_ context.Context, _ *a2asrv.ExecutorContext) iter.Seq2[a2asdk.Event, error] {
	return func(_ func(a2asdk.Event, error) bool) {}
}

// TestReferenceTaskIds : les referenceTaskIds portés par le Message atteignent
// l'agent (chaînage KYC → Scoring, FR-06).
func TestReferenceTaskIds(t *testing.T) {
	ctx, client := startAndConnect(t, refEchoExecutor{})
	msg := a2asdk.NewMessage(a2asdk.MessageRoleUser, a2asdk.NewTextPart("score"))
	msg.ReferenceTasks = []a2asdk.TaskID{"task-kyc-1a2b"}
	res, err := client.SendMessage(ctx, &a2asdk.SendMessageRequest{Message: msg})
	if err != nil {
		t.Fatalf("SendMessage : %v", err)
	}
	m, ok := res.(*a2asdk.Message)
	if !ok {
		t.Fatalf("résultat de type %T, want *Message", res)
	}
	text := ""
	for _, p := range m.Parts {
		if p.Text() != "" {
			text = p.Text()
		}
	}
	if text != "refs:task-kyc-1a2b;" {
		t.Errorf("referenceTaskIds non transmis : %q", text)
	}
}
