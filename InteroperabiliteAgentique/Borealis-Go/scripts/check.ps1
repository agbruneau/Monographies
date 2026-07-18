# check.ps1 — gate local pré-commit de Borealis pour hôtes Windows (PowerShell 7+).
# Équivalent de check.sh. Race exécuté seulement si CGO (gcc) est disponible.
# Plancher de couverture globale : $env:COVERAGE_FLOOR (défaut 90 — directive
#   utilisateur > 90 %). Plancher unitaire par paquet internal/* :
#   $env:UNIT_FLOOR (défaut 95, même directive), avec exemptions temporaires
#   documentées plus bas.
$ErrorActionPreference = 'Stop'
$floor = if ($env:COVERAGE_FLOOR) { [double]$env:COVERAGE_FLOOR } else { 90.0 }
$unitFloor = if ($env:UNIT_FLOOR) { [double]$env:UNIT_FLOOR } else { 95.0 }

$repoRoot = Split-Path -Parent $PSScriptRoot
Set-Location $repoRoot

function Step($m) { Write-Host "`n==> $m" }

Step "go build ./..."
go build ./...; if ($LASTEXITCODE) { exit 1 }
Write-Host "OK: build"

Step "go vet ./..."
go vet ./...; if ($LASTEXITCODE) { exit 1 }
Write-Host "OK: vet"

$race = ''
$hasGcc = [bool](Get-Command gcc -ErrorAction SilentlyContinue)
if ((go env CGO_ENABLED) -eq '1' -and $hasGcc) { $race = '-race' }
else { Write-Host "WARN: CGO/gcc absent — tests SANS -race (limitation d'hote ; requis sur hote CGO/Linux, DoD 16.9)." }

Step "go test $race -coverprofile=coverage.out ./..."
$testLog = Join-Path ([System.IO.Path]::GetTempPath()) "borealis-check-test.log"
if ($race) { go test -race "-coverprofile=coverage.out" "./..." 2>&1 | Tee-Object -FilePath $testLog }
else { go test "-coverprofile=coverage.out" "./..." 2>&1 | Tee-Object -FilePath $testLog }
if ($LASTEXITCODE) { exit 1 }
Write-Host "OK: test"

# Plancher unitaire par paquet (directive utilisateur 2026-07-07 : « unitaire
# >= 95 % sur internal/* »), distinct du plancher global ci-dessous.
# Paquets exemptés : auditexport et testfakes restent sous le plancher pour
# des branches défensives de connexion SDK non déclenchables sans
# instrumentation de test ad hoc (résidu documenté, non un oubli). Les
# exemptions Phase 4 (observability, conformance, fixtures) sont retirées :
# code mort supprimé, ces paquets sont à 100 %.
Step "plancher unitaire par paquet (>= $unitFloor%) sur internal/*"
$exemptUnitFloor = @('internal/auditexport', 'internal/testfakes')
$failUnit = $false
Get-Content $testLog | Where-Object { $_ -match 'github\.com/agbruneau/borealis/internal/' } | ForEach-Object {
    if ($_ -notmatch 'github\.com/agbruneau/borealis/(internal/\S+)') { return }
    $pkg = $Matches[1]
    if ($_ -notmatch '([0-9]+\.[0-9]+)%') { return } # « [no statements] » : rien à couvrir
    $pct = [double]$Matches[1]
    if ($exemptUnitFloor -contains $pkg) {
        Write-Host ("SKIP (exempte, cf. commentaire ci-dessus) : {0} = {1}%" -f $pkg, $pct)
    } elseif ($pct -lt $unitFloor) {
        Write-Host ("FAIL: {0} = {1}% < {2}%" -f $pkg, $pct, $unitFloor)
        $failUnit = $true
    } else {
        Write-Host ("OK: {0} = {1}%" -f $pkg, $pct)
    }
}
Remove-Item $testLog -ErrorAction SilentlyContinue
if ($failUnit) { exit 1 }
Write-Host "OK: plancher unitaire par paquet >= $unitFloor%"

# e2e (déterminisme des golden, PRDPlan §2 : « make e2e x3 → sorties identiques »).
# Réduit à x2 dans le gate local pour un coût raisonnable (écart documenté,
# PRDPlan §8 T2) ; chaque exécution compare déjà sa propre sortie au golden.
# -count=1 sur la 2e exécution : sans cela, go test servirait le résultat mis
# en cache de la 1re au lieu de ré-exécuter.
Step "go test -tags e2e ./test/e2e/... (x2)"
go test -tags e2e ./test/e2e/...; if ($LASTEXITCODE) { exit 1 }
go test -tags e2e -count=1 ./test/e2e/...; if ($LASTEXITCODE) { exit 1 }
Write-Host "OK: e2e (x2)"

Step "govulncheck ./..."
if (Get-Command govulncheck -ErrorAction SilentlyContinue) { govulncheck ./...; if ($LASTEXITCODE) { exit 1 }; Write-Host "OK: govulncheck" }
else { Write-Host "WARN: govulncheck absent." }

Step "golangci-lint run ./..."
if (Get-Command golangci-lint -ErrorAction SilentlyContinue) { golangci-lint run ./...; if ($LASTEXITCODE) { exit 1 }; Write-Host "OK: golangci-lint" }
else { Write-Host "WARN: golangci-lint absent." }

Step "coverage floor (>= $floor%)"
$line = (go tool cover "-func=coverage.out" 2>$null | Select-String '^total:')
if (-not $line) { Write-Host "SKIP: aucun statement couvrable (squelette)." }
else {
    $total = [double](($line -split '\s+')[-1].TrimEnd('%'))
    Write-Host ("Total coverage: {0}%" -f $total)
    if ($total -lt $floor) { Write-Host ("FAIL: coverage {0}% < {1}%" -f $total, $floor); exit 1 }
    Write-Host ("OK: coverage {0}% >= {1}%" -f $total, $floor)
}
Write-Host "`n================ gate PASS ================"
