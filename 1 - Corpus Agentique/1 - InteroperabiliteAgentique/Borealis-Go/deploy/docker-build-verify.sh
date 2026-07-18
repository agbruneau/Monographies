#!/usr/bin/env bash
# docker-build-verify.sh (M4.4, NFR-09) : construit DEUX fois la même image et
# compare les empreintes SHA256 du binaire produit. Reproductible → identiques.
#
# ⚠ NON VÉRIFIÉ LOCALEMENT : Docker absent de l'hôte du démonstrateur. Ce script
# est l'artefact d'exécution ; le lancer sur un hôte doté de Docker (voir
# docs/REPRODUCIBILITY.md). Données 100 % synthétiques.
set -euo pipefail

BIN="${1:-mcp-identity}"
EPOCH="${SOURCE_DATE_EPOCH:-1720310400}" # instant figé partagé par les deux builds
ROOT="$(cd "$(dirname "$0")/.." && pwd)"

sha_of_build() {
  local tag="$1"
  docker build \
    --file "$ROOT/deploy/Dockerfile" \
    --build-arg "BIN=$BIN" \
    --build-arg "SOURCE_DATE_EPOCH=$EPOCH" \
    --output "type=local,dest=$tag" \
    "$ROOT"
  sha256sum "$tag/app" | cut -d' ' -f1
}

A="$(mktemp -d)"; B="$(mktemp -d)"
trap 'rm -rf "$A" "$B"' EXIT

echo "Build 1/2 ($BIN)…"; H1="$(sha_of_build "$A")"
echo "Build 2/2 ($BIN)…"; H2="$(sha_of_build "$B")"

echo "SHA256 #1 : $H1"
echo "SHA256 #2 : $H2"
if [ "$H1" = "$H2" ]; then
  echo "OK : build reproductible (empreintes identiques)"
else
  echo "ÉCHEC : empreintes divergentes — build non reproductible" >&2
  exit 1
fi
