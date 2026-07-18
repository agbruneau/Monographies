package a2a

import "testing"

// TestReservedAlgRejectedBySignatureCheck : une card dont l'en-tête annonce le
// placeholder PQ réservé est rejetée (algorithme non câblé).
func TestReservedAlgRejectedBySignatureCheck(t *testing.T) {
	key := newKey(t)
	card := testCard()
	sig, err := Sign(card, key, "kid-1")
	if err != nil {
		t.Fatal(err)
	}
	forged := sig
	forged.Protected = b64.EncodeToString([]byte(`{"alg":"` + SigAlgMLDSA65 + `","kid":"kid-1"}`))
	if err := Verify(card, forged, &key.PublicKey); err == nil {
		t.Error("une card à alg PQ réservé devrait être rejetée (placeholder non implémenté)")
	}
}
