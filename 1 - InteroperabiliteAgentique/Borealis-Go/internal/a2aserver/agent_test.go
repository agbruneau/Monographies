package a2aserver

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"

	a2asdk "github.com/a2aproject/a2a-go/v2/a2a"
	"github.com/a2aproject/a2a-go/v2/a2aclient"
	"github.com/a2aproject/a2a-go/v2/a2aclient/agentcard"
	"github.com/a2aproject/a2a-go/v2/a2asrv"
	borealisa2a "github.com/agbruneau/borealis/pkg/a2a"
)

func TestReplyExecutorCancel(t *testing.T) {
	e := &ReplyExecutor{Reply: "x"}
	count := 0
	e.Cancel(context.Background(), nil)(func(a2asdk.Event, error) bool { count++; return true })
	if count != 0 {
		t.Errorf("Cancel a émis %d événements, want 0", count)
	}
}

// TestLifecycleExecutorCancel : comme ReplyExecutor, l'annulation n'a rien à
// émettre (exécution synchrone immédiate, T5).
func TestLifecycleExecutorCancel(t *testing.T) {
	e := &LifecycleExecutor{Reply: "x"}
	count := 0
	e.Cancel(context.Background(), nil)(func(a2asdk.Event, error) bool { count++; return true })
	if count != 0 {
		t.Errorf("Cancel a émis %d événements, want 0", count)
	}
}

// TestLifecycleExecutorStopsEarly : un consommateur qui arrête l'itération
// avant la fin (yield → false) doit interrompre Execute au point d'arrêt —
// les trois gardes de retour anticipé (T5, après submitted/working/état final).
func TestLifecycleExecutorStopsEarly(t *testing.T) {
	for n := 1; n <= 3; n++ {
		e := &LifecycleExecutor{Reply: "x"}
		execCtx := &a2asrv.ExecutorContext{Message: a2asdk.NewMessage(a2asdk.MessageRoleUser, a2asdk.NewTextPart("go"))}
		count := 0
		e.Execute(context.Background(), execCtx)(func(a2asdk.Event, error) bool {
			count++
			return count < n
		})
		if count != n {
			t.Errorf("n=%d : événements reçus avant arrêt = %d", n, count)
		}
	}
}

// TestSpecByName : nom connu retourne la spec (par son skill), nom inconnu
// retourne ok=false (T5).
func TestSpecByName(t *testing.T) {
	spec, ok := SpecByName("borealis-kyc")
	if !ok || spec.SkillID != "verify-identity" {
		t.Errorf("SpecByName(borealis-kyc) = %+v, ok=%v", spec, ok)
	}
	if _, ok := SpecByName("agent-inconnu"); ok {
		t.Error("nom inconnu : ok=false attendu")
	}
}

func TestAllSpecsBuildValidCards(t *testing.T) {
	for _, s := range Specs {
		if err := borealisa2a.Validate(BuildCard(s, "http://x")); err != nil {
			t.Errorf("%s : card invalide : %v", s.Name, err)
		}
	}
	if len(Specs) != 5 {
		t.Errorf("Specs = %d, want 5 agents (§6.1)", len(Specs))
	}
}

// TestAgentA2ARoundtrip exerce la chaîne A2A complète sur un agent réel servi
// par httptest : découverte par card (FR-01/FR-02), vérification de signature
// avant routage (FR-03), puis délégation message/send (FR-05 amorce).
func TestAgentA2ARoundtrip(t *testing.T) {
	ctx := context.Background()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	srv := httptest.NewUnstartedServer(nil)
	baseURL := "http://" + srv.Listener.Addr().String()

	card := BuildCard(Specs[1], baseURL) // KYC
	sig, err := borealisa2a.Sign(card, key, "kid-kyc")
	if err != nil {
		t.Fatal(err)
	}
	card.Signatures = []a2asdk.AgentCardSignature{sig}

	srv.Config.Handler = Mux(card, &ReplyExecutor{Reply: "identity verified"})
	srv.Start()
	defer srv.Close()

	// Découverte.
	resolved, err := agentcard.DefaultResolver.Resolve(ctx, baseURL)
	if err != nil {
		t.Fatalf("resolve : %v", err)
	}
	// Vérification de signature avant routage.
	if len(resolved.Signatures) == 0 {
		t.Fatal("card résolue sans signature")
	}
	if verr := borealisa2a.Verify(resolved, resolved.Signatures[0], &key.PublicKey); verr != nil {
		t.Fatalf("vérif signature de la card résolue : %v", verr)
	}
	// Délégation.
	client, err := a2aclient.NewFromCard(ctx, resolved, a2aclient.WithJSONRPCTransport(srv.Client()))
	if err != nil {
		t.Fatalf("client : %v", err)
	}
	res, err := client.SendMessage(ctx, &a2asdk.SendMessageRequest{
		Message: a2asdk.NewMessage(a2asdk.MessageRoleUser, a2asdk.NewTextPart("verify A00001")),
	})
	if err != nil {
		t.Fatalf("SendMessage : %v", err)
	}
	b, _ := json.Marshal(res)
	if !strings.Contains(string(b), "identity verified") {
		t.Errorf("réponse inattendue : %s", b)
	}
}
