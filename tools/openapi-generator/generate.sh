#!/bin/bash

set -Eeuo pipefail

DEST="pkg/client"
SPEC_FILE="./api/spec/openapi.yml"

GH_ORG="polyxia-org"
GH_HOST="github.com"
GH_REPO="morty-gateway/pkg/client"

rm -rf "${DEST}" || true
mkdir -p "${DEST}"

echo "Generating Morty client into ${DEST}"
openapi-generator generate -i "${SPEC_FILE}" \
    -g go \
    -o "${DEST}" \
    --git-user-id "${GH_ORG}" \
    --git-repo-id "${GH_REPO}" \
    --git-host "${GH_HOST}" \
    -c ./tools/openapi-generator/config.yml

rm "${DEST}/git_push.sh" || true
rm "${DEST}/.travis.yml" || true
rm -rf "${DEST}/test" || true
