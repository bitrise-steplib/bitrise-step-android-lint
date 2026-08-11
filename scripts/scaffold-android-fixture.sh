#!/bin/bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
FIXTURE_DIR="$SCRIPT_DIR/android-fixture"
GRADLE_VERSION="${GRADLE_VERSION:-7.6.4}"

cp -R "$FIXTURE_DIR"/. .

gradle wrapper --gradle-version "$GRADLE_VERSION" --distribution-type all
