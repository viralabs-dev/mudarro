#!/usr/bin/env bash
set -euo pipefail
bin="${1:?binário}"
context="${2:?contexto do cluster de teste}"
namespace="mudarro-test-$(date +%s)-$$"
root="$(mktemp -d)"
trap 'kubectl --context "$context" delete namespace "$namespace" --wait=false || true; rm -rf -- "$root"' EXIT
kubectl --context "$context" create namespace "$namespace"
mkdir "$root/k8s"
cat > "$root/k8s/app.yaml" <<YAML
apiVersion: apps/v1
kind: Deployment
metadata:
  name: web
spec:
  replicas: 1
  selector:
    matchLabels: {app: mudarro-smoke}
  template:
    metadata:
      labels: {app: mudarro-smoke}
    spec:
      containers:
        - name: web
          image: alpine:3.21
          command: [sh, -c, 'echo ready; sleep 300']
---
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: keep-data
spec:
  accessModes: [ReadWriteOnce]
  resources:
    requests: {storage: 1Mi}
YAML
cat > "$root/mudarro.yaml" <<YAML
version: 1
name: Kubernetes smoke
services:
  - id: web
    dir: .
    infrastructure:
      kind: kubernetes
      mode: manifests
      file: k8s
      context: $context
      namespace: $namespace
YAML
"$bin" run web:up --root "$root"
kubectl --context "$context" -n "$namespace" rollout status deployment/web --timeout=120s
"$bin" run web:status --root "$root"
"$bin" run web:logs --root "$root"
"$bin" run web:restart --root "$root"
"$bin" run web:down --root "$root"
test "$(kubectl --context "$context" -n "$namespace" get deployment/web -o jsonpath='{.spec.replicas}')" = 0
kubectl --context "$context" -n "$namespace" get pvc keep-data
echo 'Kubernetes: OK; PVC preservado'
