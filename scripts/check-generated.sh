#!/bin/sh
set -eu
root=$(CDPATH='' cd -- "$(dirname "$0")/.." && pwd)
tmp=$(mktemp -d "${TMPDIR:-/tmp}/calendar-generated.XXXXXX")
cd "$root"
./node_modules/.bin/tsp compile api --output-dir "$tmp/openapi"
cmp api/generated/openapi.yaml "$tmp/openapi/openapi.yaml"
(cd backend && go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v2.5.1 -generate types,std-http-server -package generated ../api/generated/openapi.yaml > "$tmp/api.gen.go")
cmp backend/internal/api/generated/api.gen.go "$tmp/api.gen.go"
./node_modules/.bin/openapi-generator-cli generate -i api/generated/openapi.yaml -g typescript-fetch -o "$tmp/sdk" --additional-properties=supportsES6=true,typescriptThreePlus=true --global-property apis,models,supportingFiles=runtime.ts:index.ts
node scripts/patch-sdk.mjs "$tmp/sdk/runtime.ts"
diff -r frontend/src/api/generated "$tmp/sdk"
