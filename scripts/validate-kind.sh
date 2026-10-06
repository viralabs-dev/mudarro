#!/usr/bin/env bash
set -euo pipefail
bin="${1:?binário}"
name="mudarro-validation-$$"
tmp="$(mktemp -d)"
trap 'kind delete cluster --name "$name"; rm -rf -- "$tmp"' EXIT
export KUBECONFIG="$tmp/kubeconfig"
kind create cluster --name "$name" --wait 120s
bash scripts/integration-k8s.sh "$bin" "kind-$name"
