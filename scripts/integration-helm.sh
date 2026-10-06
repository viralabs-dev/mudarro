#!/usr/bin/env bash
set -euo pipefail
bin="$(realpath "${1:?binary}")"
root="$(mktemp -d /tmp/mudarro-helm.XXXXXX)"
cluster="mudarro-helm-test-$$"
export KUBECONFIG="$root/kubeconfig"
trap 'kind delete cluster --name "$cluster"; rm -rf -- "$root"' EXIT
kind create cluster --name "$cluster" --wait 120s
kubectl create namespace mudarro-helm
mkdir -p "$root/chart/templates"
printf 'apiVersion: v2\nname: mudarro-test\nversion: 0.1.0\n' > "$root/chart/Chart.yaml"
cat > "$root/chart/templates/app.yaml" <<'YAML'
apiVersion: apps/v1
kind: Deployment
metadata:
  name: web
  labels:
    app.kubernetes.io/instance: web
spec:
  replicas: 1
  selector:
    matchLabels: {app: web}
  template:
    metadata:
      labels: {app: web}
    spec:
      containers:
      - name: web
        image: alpine:3.21
        command: [sh, -c, 'echo helm-ready; sleep 300']
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
cat > "$root/mudarro.yaml" <<YAML
version: 1
name: Helm test
services:
- id: web
  dir: .
  infrastructure:
    kind: kubernetes
    mode: helm
    file: chart
    context: kind-$cluster
    namespace: mudarro-helm
YAML
"$bin" doctor --root "$root"
"$bin" generate --root "$root"
"$bin" run web:up --root "$root"
kubectl -n mudarro-helm rollout status deployment/web --timeout=60s
"$bin" run web:status --root "$root"
"$bin" run web:logs --root "$root" | grep helm-ready
"$bin" run web:restart --root "$root"
kubectl -n mudarro-helm rollout status deployment/web --timeout=60s
"$bin" run web:down --root "$root"
test "$(kubectl -n mudarro-helm get deployment web -o jsonpath='{.spec.replicas}')" = 0
kubectl -n mudarro-helm get pvc keep-data
helm status web -n mudarro-helm >/dev/null
bash "$(dirname "${BASH_SOURCE[0]}")/verify-pvc-persistence.sh" "$bin" "$root" "kind-$cluster" mudarro-helm
