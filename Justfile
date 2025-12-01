set dotenv-load

default:
  @just --list

pc: pre-commit

[group('qa')]
pre-commit: tidy lint lint-integration generate

[group('qa')]
lint:
  golangci-lint run --fix --build-tags it --timeout 5m

[group('qa')]
lint-integration:
  cd ./tests/e2e && golangci-lint run --fix --build-tags it --timeout 5m

[group('qa')]
tidy:
  @go mod tidy

[group('qa')]
build:
  @go build -o ./build/terraform-provider-stack ./main.go

[group('test')]
tests: tests-unit tests-integration tests-e2e coverage

[group('test')]
coverage:
  @rm -rf coverage/coverage_merged.txt
  @head -n 1 coverage/coverage_unit.txt > coverage/coverage_merged.txt
  @tail -n +2 coverage/coverage_unit.txt | grep -Ev "generated|/sdk|tests/" >> coverage/coverage_merged.txt
  @tail -n +2 coverage/coverage_e2e.txt | grep -Ev "generated|/sdk|tests/" >> coverage/coverage_merged.txt
  @go tool cover -func=coverage/coverage_merged.txt

[group('test')]
generate:
  @go generate ./...

[group('test')]
tests-unit: 
  @mkdir -p ./coverage
  @go test -v -tags it ./internal/... -covermode=atomic -coverprofile=coverage/coverage_unit.txt -race -coverpkg=./internal/...

[group('test')]
tests-e2e tags="ci":
  @mkdir -p ./coverage
  @TF_ACC=1 go test -v -tags {{tags}} ./tests/e2e/... -covermode=atomic -coverprofile=coverage/coverage_e2e.txt -race -coverpkg=./internal/...,./cmd/...

tests-integration tags="ci":
  @mkdir -p ./coverage
  @TF_ACC=1 go test -v -tags {{tags}} ./tests/integration/... -covermode=atomic -coverprofile=coverage/coverage_integration.txt -race -coverpkg=./internal/...,./cmd/...

[group('terraform')]
init examples="install-verif": build
  @cd examples/{{examples}} && terraform init -upgrade

[group('terraform')]
plan examples="install-verif": build
  @cd examples/{{examples}} && terraform plan -generate-config-out=generated.tf

[group('terraform')]
apply examples="install-verif": build
  @cd examples/{{examples}} && terraform apply -auto-approve

[group('terraform')]
destroy examples="install-verif": build
  @cd examples/{{examples}} && terraform destroy -auto-approve 

[group('releases')]
release-local: pc
  @goreleaser release --nightly --skip=publish --clean

[group('releases')]
release-ci: pc
  @goreleaser release --nightly --clean

[group('releases')]
release: pc
  @echo "$GPG_PRIVATE_KEY" | gpg --batch --import
  @echo "$GPG_FULL_FP:6:" | gpg --import-ownertrust -
  @goreleaser release --clean

[group('deployment')]
connect-dev:
  vcluster connect $USER --server=https://kube.$USER.formance.dev


delete-all-stack force="true":
  #!/bin/bash
  set -euo pipefail

  # Récupération des IDs des stacks
  stacks=$(fctl stacks list -o json | jq -r '.data.stacks[].id')

  # Suppression de chaque stack
  for stack in $stacks; do
    if [[ "${force:-false}" == "true" ]]; then
      fctl stacks delete "$stack" --force
    else
      fctl stacks delete "$stack"
    fi
  done


generate-stack-client:
  #/usr/bin/env bash
  set -euo pipefail
  stack_version=$(yq e ".versions.base.version" ./openapi/versions.yaml)
  curl -o ./openapi/base.yaml https://raw.githubusercontent.com/formancehq/stack/refs/heads/main/releases/base.yaml
  for module in ledger payments reconciliation wallets flows webhooks auth gateway; do \
    mkdir -p ./openapi/${module}; \
    VERSION=$(yq e ".versions.${module}.version" ./openapi/versions.yaml); \
    curl -o ./openapi/${module}/openapi.yaml https://raw.githubusercontent.com/formancehq/${module}/refs/heads/${VERSION}/openapi.yaml; \
  done
  yq -i '.paths."/versions".get.security = [{"Authorization": []}]' openapi/gateway/openapi.yaml
  curl -o ./openapi/openapi-overlay.yaml https://raw.githubusercontent.com/formancehq/formance-sdk-typescript/refs/heads/main/overlay.yaml
  npx -y openapi-merge-cli -c openapi/openapi-merge.json
  speakeasy overlay apply -s openapi/generate.json -o openapi/openapi-overlay.yaml --out openapi/build.json
  speakeasy generate sdk -s openapi/build.json -o ./pkg/stack -l go