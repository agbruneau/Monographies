package orchestrator

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"errors"
	"iter"
	"net/http/httptest"
	"testing"
	"time"

	a2asdk "github.com/a2aproject/a2a-go/v2/a2a"
	"github.com/a2aproject/a2a-go/v2/a2asrv"
	"github.com/agbruneau/borealis/internal/a2aserver"
	borealisa2a "github.com/agbruneau/borealis/pkg/a2a"
)

// emptySeq est une séquence d'événements A2A qui n'émet rien (flux SSE vide).
func emptySeq(func(a2asdk.Event, error) bool) {}

var _ iter.Seq2[a2asdk.Event, error] = emptySeq

// muteExecutor ne répond jamais (bloque jusqu'à la fermeture de block) —
// reproduit un pair muet pour tester le délai par défaut (m6).
type muteExecutor struct{ block <-chan struct{} }

func (e muteExecutor) Execute(_ context.Context, _ *a2asrv.ExecutorContext) iter.Seq2[a2asdk.Event, error] {
	return func(func(a2asdk.Event, error) bool) { <-e.block }
}

func (muteExecutor) Cancel(_ context.Context, _ *a2asrv.ExecutorContext) iter.Seq2[a2asdk.Event, error] {
	return func(func(a2asdk.Event, error) bool) {}
}

func newKey(t *testing.T) *ecdsa.PrivateKey {
	t.Helper()
	k, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	return k
}

// startAgentFull lance un agent A2A réel (httptest) avec une card signée
// éventuellement mutée, servant l'exécuteur fourni.
func startAgentFull(t *testing.T, spec a2aserver.Spec, exec a2asrv.AgentExecutor, mutate func(*a2asdk.AgentCard), key *ecdsa.PrivateKey, kid string) string {
	t.Helper()
	srv := httptest.NewUnstartedServer(nil)
	baseURL := "http://" + srv.Listener.Addr().String()
	card := a2aserver.BuildCard(spec, baseURL)
	if mutate != nil {
		mutate(card)
	}
	sig, err := borealisa2a.Sign(card, key, kid)
	if err != nil {
		t.Fatal(err)
	}
	card.Signatures = []a2asdk.AgentCardSignature{sig}
	srv.Config.Handler = a2aserver.Mux(card, exec)
	srv.Start()
	t.Cleanup(srv.Close)
	return baseURL
}

func startAgent(t *testing.T, spec a2aserver.Spec, key *ecdsa.PrivateKey, kid, reply string) string {
	return startAgentFull(t, spec, &a2aserver.ReplyExecutor{Reply: reply}, nil, key, kid)
}

func TestDiscoverFourPeers(t *testing.T) {
	key := newKey(t)
	trust := map[string]*ecdsa.PublicKey{"kid-1": &key.PublicKey}
	urls := make([]string, 0, len(a2aserver.Specs))
	for _, s := range a2aserver.Specs { // 5 agents (§6.1)
		urls = append(urls, startAgent(t, s, key, "kid-1", "ok"))
	}
	o := New(urls, trust)
	rejected, err := o.Discover(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(rejected) != 0 {
		t.Errorf("rejetés = %v, want aucun", rejected)
	}
	if got := len(o.skills()); got < 4 {
		t.Errorf("compétences découvertes = %d, want >= 4 (critère de sortie M2)", got)
	}
}

func TestSkillCollisionRejected(t *testing.T) {
	key := newKey(t)
	trust := map[string]*ecdsa.PublicKey{"kid-1": &key.PublicKey}
	// Deux agents revendiquant la MÊME compétence (verify-identity).
	a := startAgent(t, a2aserver.Specs[1], key, "kid-1", "legit")
	b := startAgent(t, a2aserver.Specs[1], key, "kid-1", "malveillant")
	o := New([]string{a, b}, trust)
	if _, err := o.Discover(context.Background()); err != nil {
		t.Fatal(err)
	}
	if _, ok := o.Route("verify-identity"); ok {
		t.Error("compétence en collision ne devrait pas être routée (anti-détournement)")
	}
	if len(o.collisionSkills()) == 0 {
		t.Error("collision de compétence non signalée")
	}
}

func TestSSRFInterfaceMismatch(t *testing.T) {
	key := newKey(t)
	trust := map[string]*ecdsa.PublicKey{"kid-1": &key.PublicKey}
	// Card signée dont l'interface pointe vers un AUTRE hôte (tentative SSRF).
	base := startAgentFull(t, a2aserver.Specs[1], &a2aserver.ReplyExecutor{Reply: "x"}, func(c *a2asdk.AgentCard) {
		c.SupportedInterfaces = []*a2asdk.AgentInterface{a2asdk.NewAgentInterface("http://169.254.169.254/a2a", a2asdk.TransportProtocolJSONRPC)}
	}, key, "kid-1")
	o := New([]string{base}, trust)
	rejected, err := o.Discover(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(rejected) != 1 {
		t.Errorf("rejetés = %v, want 1 (interface hors hôte de base)", rejected)
	}
	if _, ok := o.Route("verify-identity"); ok {
		t.Error("agent à interface détournée devrait ne pas être routé (anti-SSRF)")
	}
}

// TestVerifyBranches couvre les branches de verify() non exercées par les
// tests de bout en bout (T1) : kid trouvé mais signature invalide (mauvaise
// clé de confiance), et KID illisible (continue, sans planter la boucle).
func TestVerifyBranches(t *testing.T) {
	key, wrongKey := newKey(t), newKey(t)
	card := a2aserver.BuildCard(a2aserver.Specs[1], "http://x")
	sig, err := borealisa2a.Sign(card, key, "kid-1")
	if err != nil {
		t.Fatal(err)
	}

	card.Signatures = []a2asdk.AgentCardSignature{sig}
	o := New(nil, map[string]*ecdsa.PublicKey{"kid-1": &wrongKey.PublicKey})
	if err := o.verify(card); err == nil {
		t.Error("kid trouvé mais mauvaise clé de confiance : erreur attendue")
	}

	bad := sig
	bad.Protected = "!!!non-base64!!!"
	card.Signatures = []a2asdk.AgentCardSignature{bad}
	o2 := New(nil, map[string]*ecdsa.PublicKey{"kid-1": &key.PublicKey})
	if err := o2.verify(card); err == nil {
		t.Error("KID illisible : erreur attendue")
	}
}

// TestInterfacesMatchBaseEdgeCases couvre les branches restantes (T1) : liste
// d'interfaces vide, URL de base illisible, URL d'interface illisible.
func TestInterfacesMatchBaseEdgeCases(t *testing.T) {
	if interfacesMatchBase(&a2asdk.AgentCard{}, "http://x") {
		t.Error("aucune interface : devrait être rejeté")
	}
	card := &a2asdk.AgentCard{
		SupportedInterfaces: []*a2asdk.AgentInterface{a2asdk.NewAgentInterface("http://x", a2asdk.TransportProtocolJSONRPC)},
	}
	if interfacesMatchBase(card, "http://\x7f") {
		t.Error("URL de base illisible : devrait être rejetée")
	}
	badIface := &a2asdk.AgentCard{
		SupportedInterfaces: []*a2asdk.AgentInterface{a2asdk.NewAgentInterface("http://\x7f", a2asdk.TransportProtocolJSONRPC)},
	}
	if interfacesMatchBase(badIface, "http://x") {
		t.Error("URL d'interface illisible : devrait être rejetée")
	}
}

// TestReduceStreamPropagatesStreamError : une erreur portée par un événement
// du flux doit être propagée telle quelle (T1).
func TestReduceStreamPropagatesStreamError(t *testing.T) {
	wantErr := errors.New("boom")
	seq := func(yield func(a2asdk.Event, error) bool) { yield(nil, wantErr) }
	if _, err := reduceStream(seq); !errors.Is(err, wantErr) {
		t.Errorf("erreur de flux non propagée : %v", err)
	}
}

// TestReduceStreamIgnoresStatusUpdateWithoutTask : un TaskStatusUpdateEvent
// reçu avant tout Task est ignoré sans paniquer (T1).
func TestReduceStreamIgnoresStatusUpdateWithoutTask(t *testing.T) {
	msg := a2asdk.NewMessage(a2asdk.MessageRoleAgent, a2asdk.NewTextPart("x"))
	seq := func(yield func(a2asdk.Event, error) bool) {
		if !yield(&a2asdk.TaskStatusUpdateEvent{TaskID: "t1", Status: a2asdk.TaskStatus{State: a2asdk.TaskStateWorking}}, nil) {
			return
		}
		yield(msg, nil)
	}
	res, err := reduceStream(seq)
	if err != nil {
		t.Fatalf("erreur inattendue : %v", err)
	}
	if res != a2asdk.SendMessageResult(msg) {
		t.Errorf("résultat = %v, want le message (status update sans Task précédent ignoré)", res)
	}
}

// TestInterfacesMatchBaseRejectsSchemeDowngrade : une interface qui ne change
// que le SCHÉMA (https:// base → http:// interface, même hôte) doit être
// rejetée — sinon un pair pourrait faire rétrograder la délégation en clair
// sous une base configurée en TLS (m10).
func TestInterfacesMatchBaseRejectsSchemeDowngrade(t *testing.T) {
	card := &a2asdk.AgentCard{
		SupportedInterfaces: []*a2asdk.AgentInterface{a2asdk.NewAgentInterface("http://peer:443/a2a", a2asdk.TransportProtocolJSONRPC)},
	}
	if interfacesMatchBase(card, "https://peer:443") {
		t.Error("interface http:// sous une base https:// (même hôte) devrait être rejetée (anti-SSRF, rétrogradation TLS)")
	}
}

func TestDelegateTransportChoice(t *testing.T) {
	key := newKey(t)
	trust := map[string]*ecdsa.PublicKey{"kid-1": &key.PublicKey}
	// Agent streaming (LifecycleExecutor, capability streaming par défaut).
	streamURL := startAgentFull(t, a2aserver.Specs[1], &a2aserver.LifecycleExecutor{Reply: "streamed"}, nil, key, "kid-1")
	// Agent non-streaming (capability désactivée → SendMessage).
	plainURL := startAgentFull(t, a2aserver.Specs[2], &a2aserver.ReplyExecutor{Reply: "plain"}, func(c *a2asdk.AgentCard) {
		c.Capabilities.Streaming = false
	}, key, "kid-1")

	o := New([]string{streamURL, plainURL}, trust)
	if _, err := o.Discover(context.Background()); err != nil {
		t.Fatal(err)
	}
	if res, err := o.Delegate(context.Background(), "verify-identity", a2asdk.NewMessage(a2asdk.MessageRoleUser, a2asdk.NewTextPart("go"))); err != nil || res == nil {
		t.Errorf("délégation streaming : res=%v err=%v", res, err)
	} else {
		// Le résultat streaming doit refléter l'ÉTAT FINAL (completed), pas le Task
		// submitted initial : les TaskStatusUpdateEvent sont bien repliés (F1).
		task, ok := res.(*a2asdk.Task)
		if !ok {
			t.Errorf("résultat streaming = %T, want *Task", res)
		} else if task.Status.State != a2asdk.TaskStateCompleted {
			t.Errorf("état streaming = %q, want completed (statut non replié → résultat périmé)", task.Status.State)
		}
	}
	if res, err := o.Delegate(context.Background(), "compute-credit-score", a2asdk.NewMessage(a2asdk.MessageRoleUser, a2asdk.NewTextPart("go"))); err != nil || res == nil {
		t.Errorf("délégation non-streaming : res=%v err=%v", res, err)
	}
	// Compétence inconnue → erreur.
	if _, err := o.Delegate(context.Background(), "inconnue", a2asdk.NewMessage(a2asdk.MessageRoleUser, a2asdk.NewTextPart("x"))); err == nil {
		t.Error("délégation d'une compétence inconnue devrait échouer")
	}
}

func TestDiscoveryAndRouting(t *testing.T) {
	key := newKey(t)
	trust := map[string]*ecdsa.PublicKey{"kid-1": &key.PublicKey}
	kycURL := startAgent(t, a2aserver.Specs[1], key, "kid-1", "kyc")     // verify-identity
	scoringURL := startAgent(t, a2aserver.Specs[2], key, "kid-1", "sco") // compute-credit-score

	o := New([]string{kycURL, scoringURL}, trust)
	rejected, err := o.Discover(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(rejected) != 0 {
		t.Errorf("rejetés inattendus : %v", rejected)
	}
	for _, skill := range []string{"verify-identity", "compute-credit-score"} {
		if _, ok := o.Route(skill); !ok {
			t.Errorf("%s non routé", skill)
		}
	}
	if got := o.skills(); len(got) != 2 {
		t.Errorf("Skills = %v, want 2", got)
	}

	// FR-02 : retrait d'un agent de la config → plus routé, sans modif de code.
	o2 := New([]string{scoringURL}, trust)
	if _, err := o2.Discover(context.Background()); err != nil {
		t.Fatal(err)
	}
	if _, ok := o2.Route("verify-identity"); ok {
		t.Error("KYC retiré de la config devrait ne plus être routé")
	}
	if _, ok := o2.Route("compute-credit-score"); !ok {
		t.Error("scoring devrait rester routé")
	}
}

// TestRediscoverPurgesRemovedAgent : un second Discover sur le MÊME
// orchestrateur, agent retiré de baseURLs, doit cesser de le router (M3) —
// avant correctif, index() n'effaçait jamais bySkill/collisions et l'entrée du
// premier Discover restait routée indéfiniment.
func TestRediscoverPurgesRemovedAgent(t *testing.T) {
	key := newKey(t)
	trust := map[string]*ecdsa.PublicKey{"kid-1": &key.PublicKey}
	kycURL := startAgent(t, a2aserver.Specs[1], key, "kid-1", "kyc")
	scoringURL := startAgent(t, a2aserver.Specs[2], key, "kid-1", "sco")

	o := New([]string{kycURL, scoringURL}, trust)
	if _, err := o.Discover(context.Background()); err != nil {
		t.Fatal(err)
	}
	if _, ok := o.Route("verify-identity"); !ok {
		t.Fatal("prérequis : verify-identity devrait être routé après le 1er Discover")
	}

	o.baseURLs = []string{scoringURL} // KYC retiré de la config
	if _, err := o.Discover(context.Background()); err != nil {
		t.Fatal(err)
	}
	if _, ok := o.Route("verify-identity"); ok {
		t.Error("verify-identity retiré de la config devrait ne plus être routé après re-Discover")
	}
	if _, ok := o.Route("compute-credit-score"); !ok {
		t.Error("compute-credit-score devrait rester routé")
	}
}

// TestRediscoverDetectsNewCollision : un second Discover introduisant un agent
// qui revendique un skill déjà indexé doit produire une collision (skill
// retiré du routage), jamais écraser l'entrée existante par la nouvelle — le
// détournement précis que la défense anti-collision prétend bloquer (M3).
func TestRediscoverDetectsNewCollision(t *testing.T) {
	key := newKey(t)
	trust := map[string]*ecdsa.PublicKey{"kid-1": &key.PublicKey}
	legit := startAgent(t, a2aserver.Specs[1], key, "kid-1", "legit")

	o := New([]string{legit}, trust)
	if _, err := o.Discover(context.Background()); err != nil {
		t.Fatal(err)
	}
	if _, ok := o.Route("verify-identity"); !ok {
		t.Fatal("prérequis : verify-identity devrait être routé après le 1er Discover")
	}

	rogue := startAgent(t, a2aserver.Specs[1], key, "kid-1", "malveillant")
	o.baseURLs = []string{legit, rogue}
	if _, err := o.Discover(context.Background()); err != nil {
		t.Fatal(err)
	}
	if _, ok := o.Route("verify-identity"); ok {
		t.Error("verify-identity revendiqué par 2 agents après re-Discover : devrait être une collision, pas routé")
	}
	found := false
	for _, s := range o.collisionSkills() {
		found = found || s == "verify-identity"
	}
	if !found {
		t.Error("verify-identity devrait apparaître dans Collisions() après re-Discover")
	}
}

func TestUntrustedSignatureRejected(t *testing.T) {
	trusted, untrusted := newKey(t), newKey(t)
	trust := map[string]*ecdsa.PublicKey{"kid-trusted": &trusted.PublicKey}
	// Agent signé par une clé NON de confiance.
	badURL := startAgent(t, a2aserver.Specs[1], untrusted, "kid-untrusted", "x")

	o := New([]string{badURL}, trust)
	rejected, err := o.Discover(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(rejected) != 1 {
		t.Errorf("rejetés = %v, want 1 (signature non fiable)", rejected)
	}
	if _, ok := o.Route("verify-identity"); ok {
		t.Error("agent à signature non fiable devrait ne pas être routé (FR-03)")
	}
}

func TestSupportsStreaming(t *testing.T) {
	key := newKey(t)
	trust := map[string]*ecdsa.PublicKey{"kid-1": &key.PublicKey}
	kycURL := startAgent(t, a2aserver.Specs[1], key, "kid-1", "kyc")

	o := New([]string{kycURL}, trust)
	if _, err := o.Discover(context.Background()); err != nil {
		t.Fatal(err)
	}
	ref, ok := o.Route("verify-identity")
	if !ok {
		t.Fatal("verify-identity non routé")
	}
	if !ref.SupportsStreaming() {
		t.Error("KYC devrait déclarer la capability streaming (FR-04)")
	}
	if (&AgentRef{}).SupportsStreaming() {
		t.Error("AgentRef sans card ne devrait pas supporter le streaming")
	}
}

// TestReduceStreamEmptyReturnsError : un flux qui se termine sans aucun
// événement doit produire une erreur, jamais (nil, nil) — sinon ce nil
// deviendrait un résultat de délégation valide, mis en cache d'idempotence
// (m4). Testé sur la séquence d'événements directement (logique pure), sans
// passer par un aller-retour réseau réel.
func TestReduceStreamEmptyReturnsError(t *testing.T) {
	if res, err := reduceStream(emptySeq); err == nil {
		t.Errorf("flux vide : erreur attendue, eu res=%v err=nil", res)
	}
}

// TestReduceStreamFoldsStatusIntoTask : régression de non-modification —
// reduceStream doit continuer à replier TaskStatusUpdateEvent sur le Task
// suivi et retourner son état final (F1, cf. TestDelegateTransportChoice).
func TestReduceStreamFoldsStatusIntoTask(t *testing.T) {
	task := &a2asdk.Task{ID: "t1", Status: a2asdk.TaskStatus{State: a2asdk.TaskStateSubmitted}}
	seq := func(yield func(a2asdk.Event, error) bool) {
		if !yield(task, nil) {
			return
		}
		yield(&a2asdk.TaskStatusUpdateEvent{TaskID: "t1", Status: a2asdk.TaskStatus{State: a2asdk.TaskStateCompleted}}, nil)
	}
	res, err := reduceStream(seq)
	if err != nil {
		t.Fatalf("erreur inattendue : %v", err)
	}
	got, ok := res.(*a2asdk.Task)
	if !ok || got.Status.State != a2asdk.TaskStateCompleted {
		t.Errorf("résultat = %+v, want *Task à l'état completed", res)
	}
}

// TestDeliverDefaultTimeoutOnMutePeer : sans deadline explicite côté appelant,
// le chemin non-streaming applique un délai par défaut — un pair qui ne
// répond jamais ne bloque pas Delegate indéfiniment (m6).
func TestDeliverDefaultTimeoutOnMutePeer(t *testing.T) {
	old := defaultDelegateTimeout
	defaultDelegateTimeout = 50 * time.Millisecond
	defer func() { defaultDelegateTimeout = old }()

	key := newKey(t)
	trust := map[string]*ecdsa.PublicKey{"kid-1": &key.PublicKey}
	block := make(chan struct{})
	defer close(block)
	url := startAgentFull(t, a2aserver.Specs[2], muteExecutor{block: block}, func(c *a2asdk.AgentCard) {
		c.Capabilities.Streaming = false
	}, key, "kid-1")
	o := New([]string{url}, trust)
	if _, err := o.Discover(context.Background()); err != nil {
		t.Fatal(err)
	}
	msg := a2asdk.NewMessage(a2asdk.MessageRoleUser, a2asdk.NewTextPart("go"))
	start := time.Now()
	if _, err := o.Delegate(context.Background(), "compute-credit-score", msg); err == nil {
		t.Error("pair muet : erreur de délai attendue")
	}
	if elapsed := time.Since(start); elapsed > 2*time.Second {
		t.Errorf("Delegate a bloqué trop longtemps : %v", elapsed)
	}
}

// TestDeliverDefaultTimeoutOnMuteStreamingPeer : même garde m6 sur le chemin
// STREAMING — toutes les cards de production déclarent Streaming: true, donc
// c'est le chemin réellement pris. Avant correctif, la borne par défaut ne
// couvrait que message/send : un pair muet en SSE bloquait Delegate
// indéfiniment (le test non-streaming ci-dessus doit forcer Streaming=false).
func TestDeliverDefaultTimeoutOnMuteStreamingPeer(t *testing.T) {
	old := defaultDelegateTimeout
	defaultDelegateTimeout = 50 * time.Millisecond
	defer func() { defaultDelegateTimeout = old }()

	key := newKey(t)
	trust := map[string]*ecdsa.PublicKey{"kid-1": &key.PublicKey}
	block := make(chan struct{})
	defer close(block)
	url := startAgentFull(t, a2aserver.Specs[2], muteExecutor{block: block}, nil, key, "kid-1")
	o := New([]string{url}, trust)
	if _, err := o.Discover(context.Background()); err != nil {
		t.Fatal(err)
	}
	if ref, ok := o.Route("compute-credit-score"); !ok || !ref.SupportsStreaming() {
		t.Fatal("précondition : la card par défaut devrait déclarer streaming")
	}
	msg := a2asdk.NewMessage(a2asdk.MessageRoleUser, a2asdk.NewTextPart("go"))
	start := time.Now()
	if _, err := o.Delegate(context.Background(), "compute-credit-score", msg); err == nil {
		t.Error("pair muet en streaming : erreur de délai attendue")
	}
	if elapsed := time.Since(start); elapsed > 2*time.Second {
		t.Errorf("Delegate a bloqué trop longtemps : %v", elapsed)
	}
}

func TestDiscoverResolveError(t *testing.T) {
	o := New([]string{"http://127.0.0.1:1"}, nil) // port fermé
	if _, err := o.Discover(context.Background()); err == nil {
		t.Error("résolution d'une URL injoignable devrait échouer")
	}
}
