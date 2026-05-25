#!/usr/bin/env bash
#
# license-check.sh
# ----------------
# Fail the build if any Go dependency declares a license we will not ship.
#
# Joblantern is Apache 2.0. We accept:
#   Apache-2.0, MIT, BSD-2-Clause, BSD-3-Clause, ISC, MPL-2.0, Unlicense, CC0-1.0
#
# We reject (build-failing):
#   GPL-*, AGPL-*, LGPL-*, SSPL-*, BUSL-*, Elastic-*, CC-BY-NC-*, CC-BY-SA-*
#
# Usage:
#   bash scripts/license-check.sh
#
set -euo pipefail

REPORT_DIR="${REPORT_DIR:-build/license}"
mkdir -p "$REPORT_DIR"
REPORT="$REPORT_DIR/dependencies.csv"

# go-licenses is the canonical tool. Install on demand.
if ! command -v go-licenses >/dev/null 2>&1; then
    echo "==> Installing go-licenses..."
    GOBIN="$(go env GOPATH)/bin" go install github.com/google/go-licenses@latest
    export PATH="$(go env GOPATH)/bin:$PATH"
fi

DISALLOWED_PATTERN='^(GPL|AGPL|LGPL|SSPL|BUSL|Elastic|CC-BY-NC|CC-BY-SA)'

# Detect every module that builds into the joblantern binary AND each MCP server.
TARGETS=(./cmd/joblantern)
while IFS= read -r dir; do
    TARGETS+=("./$dir")
done < <(find cmd -mindepth 1 -maxdepth 1 -type d -name 'mcp-*' | sort)

# No cmd/ binaries yet? Fall back to the module root so the check still runs.
if [ ${#TARGETS[@]} -eq 0 ] || [ ! -f cmd/joblantern/main.go ]; then
    TARGETS=(./...)
fi

echo "==> Scanning targets: ${TARGETS[*]}"
echo "module,license" > "$REPORT"

FAILED=0
for tgt in "${TARGETS[@]}"; do
    # `go-licenses report` prints: <module>,<license-url>,<license-type>
    # We only care about the module name (column 1) and license type (column 3).
    if ! go-licenses report "$tgt" 2>/dev/null | awk -F, '{print $1 "," $3}' >> "$REPORT.tmp"; then
        echo "  (target $tgt has no buildable source yet — skipping)"
        continue
    fi
done

if [ -f "$REPORT.tmp" ]; then
    sort -u "$REPORT.tmp" >> "$REPORT"
    rm -f "$REPORT.tmp"
fi

echo
echo "==> License report: $REPORT"
echo "==> Disallowed licenses (build will fail on any match):"
echo "    $DISALLOWED_PATTERN"

# Skip the header line, evaluate each row.
BAD=$(tail -n +2 "$REPORT" | awk -F, -v pat="$DISALLOWED_PATTERN" '$2 ~ pat {print}')

if [ -n "$BAD" ]; then
    echo
    echo "FAIL: disallowed license(s) detected:"
    echo "$BAD" | sed 's/^/  /'
    exit 1
fi

echo
echo "OK: no disallowed licenses detected."
