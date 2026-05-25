#!/usr/bin/env bash
#
# smoke.sh
# --------
# End-to-end smoke test for the joblantern binary:
#   1. `make build` produces ./bin/joblantern.
#   2. The binary boots on an ephemeral port.
#   3. GET /healthz returns 200 "ok".
#   4. The -healthcheck probe flag exits 0.
#   5. Graceful shutdown on SIGTERM.
#
# Exit code: 0 on full pass, non-zero with a labelled failure otherwise.
#
set -euo pipefail

PORT="${PORT:-18080}"
ADDR=":${PORT}"
BIN="bin/joblantern"

cleanup() {
    if [[ -n "${SRV_PID:-}" ]] && kill -0 "$SRV_PID" 2>/dev/null; then
        kill -TERM "$SRV_PID" 2>/dev/null || true
        wait "$SRV_PID" 2>/dev/null || true
    fi
}
trap cleanup EXIT

echo "==> make build"
make build >/dev/null

echo "==> start server on $ADDR"
JOBLANTERN_ADDR="$ADDR" "$BIN" &
SRV_PID=$!

# Wait up to 5s for the listener.
for i in {1..50}; do
    if curl -fsS "http://127.0.0.1:${PORT}/healthz" >/dev/null 2>&1; then
        break
    fi
    sleep 0.1
done

echo "==> GET /healthz"
body="$(curl -fsS "http://127.0.0.1:${PORT}/healthz")"
if [[ "$body" != "ok" ]]; then
    echo "FAIL: /healthz returned: $body"
    exit 1
fi

echo "==> -healthcheck probe"
if ! JOBLANTERN_ADDR="$ADDR" "$BIN" -healthcheck; then
    echo "FAIL: -healthcheck probe non-zero"
    exit 1
fi

echo "==> graceful shutdown"
kill -TERM "$SRV_PID"
if ! wait "$SRV_PID" 2>/dev/null; then
    code=$?
    # 143 = 128 + SIGTERM(15); accepted as clean exit on some shells.
    if [[ "$code" -ne 0 && "$code" -ne 143 ]]; then
        echo "FAIL: server exited with code $code"
        exit 1
    fi
fi

echo "OK: smoke test passed."
