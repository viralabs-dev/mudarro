#!/usr/bin/env bash
set -euo pipefail
bin="${1:?binário}"
context="${2:?contexto do cluster de teste}"
mode="${MUDARRO_K8S_MODE:-manifests}"
case "$mode" in manifests|kustomize) ;; *) echo "Unsupported fixture mode: $mode" >&2; exit 1;; esac
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
          volumeMounts:
          - name: data
            mountPath: /data
      volumes:
      - name: data
        persistentVolumeClaim:
          claimName: keep-data
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
if [[ "$mode" = kustomize ]]; then
  printf 'apiVersion: kustomize.config.k8s.io/v1beta1\nkind: Kustomization\nresources: [app.yaml]\n' > "$root/k8s/kustomization.yaml"
fi
cat > "$root/mudarro.yaml" <<YAML
version: 1
name: Kubernetes smoke
services:
  - id: web
    dir: .
    infrastructure:
      kind: kubernetes
      mode: $mode
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
bash "$(dirname "${BASH_SOURCE[0]}")/verify-pvc-persistence.sh" "$bin" "$root" "$context" "$namespace"
