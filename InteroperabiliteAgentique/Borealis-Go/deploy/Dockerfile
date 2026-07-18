# syntax=docker/dockerfile:1.7
# Dockerfile multi-étages **reproductible** (NFR-09, M4.4) : un binaire par
# service, sélectionné par l'argument BIN. Deux builds identiques → SHA256
# identiques (voir deploy/docker-build-verify.sh).
#
# ⚠ Bases épinglées par digest à FIXER avant tout build de production (les
# @sha256 ci-dessous sont des repères). Le démonstrateur n'a pas pu exécuter le
# build (Docker absent de l'hôte) — artefact livré, exécution « non vérifiée
# localement » (cf. docs/REPRODUCIBILITY.md).

FROM golang:1.26 AS build
WORKDIR /src

# Reproductibilité du BINAIRE comparé : elle provient de -trimpath + -ldflags
# "-s -w -buildid=" + CGO_ENABLED=0 (pas d'horodatage embarqué). SOURCE_DATE_EPOCH
# n'affecte PAS le binaire (go build ne le lit pas) ; il sert à l'horodatage de
# l'IMAGE côté BuildKit — conservé pour la reproductibilité au niveau image.
ARG SOURCE_DATE_EPOCH=0
ARG BIN
ENV CGO_ENABLED=0 GOFLAGS=-trimpath

# Dépendances d'abord (cache reproductible).
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -trimpath -ldflags "-s -w -buildid=" -o /out/app ./cmd/${BIN}

# Runtime minimal, sans shell (surface réduite). ⚠ épingler par digest.
FROM gcr.io/distroless/static:nonroot AS runtime
COPY --from=build /out/app /app
USER nonroot:nonroot
ENTRYPOINT ["/app"]
