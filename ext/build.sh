#!/usr/bin/env bash
#
# Build Joblantern WebExtension zips for Chrome and Firefox.
#
# Output:
#   dist/joblantern-chrome-<version>.zip
#   dist/joblantern-firefox-<version>.zip
#
# Each zip contains the per-browser manifest at root + all shared assets.
#
set -euo pipefail

VERSION="${VERSION:-0.1.0}"
ROOT="$(cd "$(dirname "$0")" && pwd)"
DIST="$ROOT/dist"
mkdir -p "$DIST"

build_one() {
  local browser="$1"   # chrome | firefox
  local stage="$DIST/$browser-build"
  rm -rf "$stage"
  mkdir -p "$stage"

  cp "$ROOT/$browser/manifest.json" "$stage/manifest.json"
  cp -R "$ROOT/shared/"* "$stage/"

  ( cd "$stage" && zip -qr "$DIST/joblantern-$browser-$VERSION.zip" . )
  rm -rf "$stage"
  echo "  built dist/joblantern-$browser-$VERSION.zip"
}

echo "==> Joblantern WebExtension v$VERSION"
if [ $# -gt 0 ]; then
  for browser in "$@"; do build_one "$browser"; done
else
  build_one chrome
  build_one firefox
fi
echo "OK"
