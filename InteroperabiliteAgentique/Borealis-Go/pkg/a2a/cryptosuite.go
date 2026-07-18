package a2a

// Suite cryptographique des signatures de card (crypto-agilité, ADR-0006). Le
// vérifieur lit l'algorithme dans l'en-tête JWS `alg`, il ne le présume pas :
// migrer d'algorithme n'exige donc pas de modifier les appelants.
const (
	// SigAlgES256 : algorithme classique ACTIF (ECDSA P-256, stdlib).
	SigAlgES256 = "ES256"

	// SigAlgMLDSA65 : emplacement post-quantique **réservé** (ML-DSA-65 / FIPS 204).
	// Placeholder de migration — AUCUNE crypto PQ réelle (NG-9) : Sign/Verify le
	// refusent tant qu'une implémentation n'est pas câblée (candidat #2, cf.
	// [[0006-crypto-agilite]] et [[0008-report-module-identite]]). Encodé ici pour
	// fixer la couture de migration (M5.6).
	SigAlgMLDSA65 = "ML-DSA-65"
)
