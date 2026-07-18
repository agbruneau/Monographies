#!/usr/bin/env bash
# check.sh — gate local pré-commit de Borealis (règle §2 du PRDPlan : « pas de CI
# distante, la rigueur est locale »). S'arrête à la première défaillance dure.
#
# Étapes : build → vet → test (+race si CGO) → plancher unitaire par paquet →
#   e2e (x2) → govulncheck → golangci-lint → plancher de couverture global.
#
# Plancher de couverture : variable COVERAGE_FLOOR (défaut 90 — directive utilisateur
#   « couverture > 90 % » du 2026-07-07, qui relève NFR-11 de 80 → 90 % global) ;
#   variable UNIT_FLOOR (défaut 95 — même directive : « unitaire >= 95 % sur
#   internal/* »), avec exemptions temporaires documentées (voir plus bas).
# Race : requiert CGO (gcc). Sur hôte sans C-toolchain (Windows nu), les tests
#   tournent SANS -race, avec avertissement — c'est une limitation d'hôte, pas
#   un choix de conception (DoD §16.9 : race requis sur hôte CGO/Linux ou conteneur).
# Lint & govulncheck : durs s'ils sont présents ; avertissement s'ils sont absents.
set -euo pipefail

COVERAGE_FLOOR="${COVERAGE_FLOOR:-90.0}"

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
cd "${REPO_ROOT}"

step() { printf '\n==> %s\n' "$1"; }

step "go build ./..."
go build ./...
echo "OK: build"

step "go vet ./..."
go vet ./...
echo "OK: vet"

# Tests (+race si un C-toolchain est disponible).
if [ "$(go env CGO_ENABLED)" = "1" ] && command -v gcc >/dev/null 2>&1; then
	RACE="-race"
else
	RACE=""
	echo "WARN: CGO/gcc absent — tests SANS -race (limitation d'hôte ; requis sur hôte CGO/Linux, DoD §16.9)."
fi
UNIT_FLOOR="${UNIT_FLOOR:-95.0}"
TEST_LOG="$(mktemp)"
trap 'rm -f "${TEST_LOG}"' EXIT

step "go test ${RACE} -coverprofile=coverage.out ./..."
# shellcheck disable=SC2086
go test ${RACE} -coverprofile=coverage.out ./... | tee "${TEST_LOG}"
echo "OK: test"

# Plancher unitaire par paquet (directive utilisateur 2026-07-07 : « unitaire
# >= 95 % sur internal/* »), distinct du plancher global ci-dessous.
# Paquets exemptés : auditexport et testfakes restent sous le plancher pour
# des branches défensives de connexion SDK non déclenchables sans
# instrumentation de test ad hoc (résidu documenté, non un oubli). Les
# exemptions Phase 4 (observability, conformance, fixtures) sont retirées :
# code mort supprimé, ces paquets sont à 100 %.
step "plancher unitaire par paquet (>= ${UNIT_FLOOR}%) sur internal/*"
EXEMPT_UNIT_FLOOR="internal/auditexport internal/testfakes"
fail_unit=0
while IFS= read -r line; do
	pkg="$(printf '%s' "${line}" | awk '{print $2}' | sed 's#github.com/agbruneau/borealis/##')"
	pct="$(printf '%s' "${line}" | grep -oE '[0-9]+\.[0-9]+%' | tr -d '%')" || true
	[ -z "${pct}" ] && continue # « [no statements] » : rien à couvrir
	skip=0
	for e in ${EXEMPT_UNIT_FLOOR}; do [ "${pkg}" = "${e}" ] && skip=1; done
	if [ "${skip}" = 1 ]; then
		printf 'SKIP (exempté, cf. commentaire ci-dessus) : %s = %s%%\n' "${pkg}" "${pct}"
		continue
	fi
	if awk -v t="${pct}" -v f="${UNIT_FLOOR}" 'BEGIN { exit (t+0 < f+0) }'; then
		printf 'OK: %s = %s%%\n' "${pkg}" "${pct}"
	else
		printf 'FAIL: %s = %s%% < %s%%\n' "${pkg}" "${pct}" "${UNIT_FLOOR}"
		fail_unit=1
	fi
done < <(grep -E '\sgithub\.com/agbruneau/borealis/internal/' "${TEST_LOG}")
if [ "${fail_unit}" = 1 ]; then
	exit 1
fi
echo "OK: plancher unitaire par paquet >= ${UNIT_FLOOR}%"

# e2e (déterminisme des golden, PRDPlan §2 : « make e2e x3 → sorties identiques »).
# Réduit à x2 dans le gate local pour un coût raisonnable (écart documenté,
# PRDPlan §8 T2) ; chaque exécution compare déjà sa propre sortie au golden.
step "go test -tags e2e ./test/e2e/... (x2)"
go test -tags e2e ./test/e2e/...
go test -tags e2e -count=1 ./test/e2e/...
echo "OK: e2e (x2)"

step "govulncheck ./..."
if command -v govulncheck >/dev/null 2>&1; then
	govulncheck ./...
	echo "OK: govulncheck"
else
	echo "WARN: govulncheck absent (go install golang.org/x/vuln/cmd/govulncheck@latest)."
fi

step "golangci-lint run ./..."
if command -v golangci-lint >/dev/null 2>&1; then
	golangci-lint run ./...
	echo "OK: golangci-lint"
else
	echo "WARN: golangci-lint absent."
fi

step "coverage floor (>= ${COVERAGE_FLOOR}%)"
total="$(go tool cover -func=coverage.out 2>/dev/null | awk '/^total:/ {gsub(/%/,"",$3); print $3}')" || true
if [ -z "${total}" ]; then
	echo "SKIP: aucun statement couvrable (squelette)."
else
	printf 'Total coverage: %s%%\n' "${total}"
	if awk -v t="${total}" -v f="${COVERAGE_FLOOR}" 'BEGIN { exit (t+0 < f+0) }'; then
		printf 'OK: coverage %s%% >= %s%%\n' "${total}" "${COVERAGE_FLOOR}"
	else
		printf 'FAIL: coverage %s%% < %s%%\n' "${total}" "${COVERAGE_FLOOR}"
		exit 1
	fi
fi

echo ""
echo "================ gate PASS ================"
