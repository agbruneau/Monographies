// Package security regroupe la vérification STRIDE consolidée du démonstrateur
// (PRD §11.3). Il ne contient pas de code de production : les défenses vivent
// dans les paquets concernés (pkg/a2a pour la signature de card, internal/audit
// pour la chaîne de hachage, internal/idpmock pour l'audience du jeton). La suite
// stride_test.go les exerce ensemble, en boîte noire, comme un unique artefact
// de sécurité traçable (M3.8) :
//
//	Menace STRIDE     Vecteur simulé              Défense vérifiée
//	----------------- --------------------------- ------------------------------
//	Spoofing          AgentCard falsifiée         pkg/a2a Verify rejette (ES256/JCS)
//	Tampering         Réécriture / troncature      audit.VerifyEntries (réécriture)
//	                  d'un log exporté             + VerifyExport ancré (troncature/effacement)
//	Elevation of priv. Jeton d'audience étrangère  idpmock.Verify rejette (RFC 8707)
//
// Les menaces Repudiation, Information disclosure et Denial of service sont
// traitées ailleurs (journal WORM non répudiable, pseudonymisation KYC, PEP
// fail-closed + disjoncteur) et hors du périmètre de cette suite ciblée.
package security
