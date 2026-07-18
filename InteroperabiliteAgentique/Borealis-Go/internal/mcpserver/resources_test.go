package mcpserver

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/agbruneau/borealis/internal/fixtures"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestApplicantIDFromURI(t *testing.T) {
	if id, ok := applicantIDFromURI("credit:///application/A00001/assessment"); !ok || id != "A00001" {
		t.Errorf("URI valide : id=%q ok=%v", id, ok)
	}
	for _, u := range []string{
		"file:///x",
		"credit:///application//assessment", // id vide
		"credit:///application/A1",          // pas de suffixe
		"credit:///other/A1/assessment",     // pas le bon préfixe
	} {
		if _, ok := applicantIDFromURI(u); ok {
			t.Errorf("%q devrait être invalide", u)
		}
	}
}

func TestAssessmentResource(t *testing.T) {
	ds := fixtures.Generate(fixtures.DefaultSeed, 100)
	cs := connectServer(t, NewBureauServer(ds))
	ctx := context.Background()

	res, err := cs.ReadResource(ctx, &mcp.ReadResourceParams{URI: "credit:///application/A00001/assessment"})
	if err != nil {
		t.Fatalf("ReadResource : %v", err)
	}
	if len(res.Contents) == 0 {
		t.Fatal("aucun contenu")
	}
	var report BureauOut
	if err := json.Unmarshal([]byte(res.Contents[0].Text), &report); err != nil {
		t.Fatalf("contenu illisible : %v", err)
	}
	if report.Score < 300 {
		t.Errorf("score hors plage : %d", report.Score)
	}
	// URI dont l'id est absent → ResourceNotFoundError.
	if _, err := cs.ReadResource(ctx, &mcp.ReadResourceParams{URI: "credit:///application/ZZZ/assessment"}); err == nil {
		t.Error("URI inconnue : erreur attendue")
	}
}

func TestAdvicePrompt(t *testing.T) {
	ds := fixtures.Generate(fixtures.DefaultSeed, 10)
	cs := connectServer(t, NewBureauServer(ds))
	ctx := context.Background()

	fr, err := cs.GetPrompt(ctx, &mcp.GetPromptParams{Name: "credit_assessment_advice", Arguments: map[string]string{"applicantId": "A00001", "lang": "fr"}})
	if err != nil {
		t.Fatalf("GetPrompt fr : %v", err)
	}
	if len(fr.Messages) == 0 || !strings.Contains(fr.Messages[0].Content.(*mcp.TextContent).Text, "pré-qualification") {
		t.Errorf("prompt fr inattendu : %+v", fr.Messages)
	}
	en, err := cs.GetPrompt(ctx, &mcp.GetPromptParams{Name: "credit_assessment_advice", Arguments: map[string]string{"applicantId": "A00001", "lang": "en"}})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(en.Messages[0].Content.(*mcp.TextContent).Text, "pre-qualification") {
		t.Errorf("prompt en inattendu : %q", en.Messages[0].Content.(*mcp.TextContent).Text)
	}
}

func TestProgressNotification(t *testing.T) {
	ds := fixtures.Generate(fixtures.DefaultSeed, 10)
	got := make(chan struct{}, 8)
	client := mcp.NewClient(&mcp.Implementation{Name: "t", Version: "v0"}, &mcp.ClientOptions{
		ProgressNotificationHandler: func(_ context.Context, _ *mcp.ProgressNotificationClientRequest) {
			select {
			case got <- struct{}{}:
			default:
			}
		},
	})
	ctx := context.Background()
	t1, t2 := mcp.NewInMemoryTransports()
	if _, err := NewBureauServer(ds).Connect(ctx, t1, nil); err != nil {
		t.Fatal(err)
	}
	cs, err := client.Connect(ctx, t2, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = cs.Close() }()

	params := &mcp.CallToolParams{Name: "get_credit_report", Arguments: map[string]any{"applicantId": "A00001"}}
	params.SetProgressToken("tok-1")
	if _, err := cs.CallTool(ctx, params); err != nil {
		t.Fatalf("CallTool : %v", err)
	}
	select {
	case <-got:
	case <-time.After(3 * time.Second):
		t.Fatal("aucune ProgressNotification reçue (FR-14)")
	}
}
