#!/usr/bin/env bash
# Note: This is a hack to force goreleaser using Garble as a compiler.
$GARBLE_PATH -literals build -trimpath -ldflags="-s -w -X main.version=$RELEASE_VERSION -X main.commit=$RELEASE_COMMIT" -o ${@: -2} .
