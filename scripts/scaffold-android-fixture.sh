#!/bin/bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
FIXTURE_DIR="$SCRIPT_DIR/android-fixture"

cp -R "$FIXTURE_DIR"/. .
chmod +x gradlew
