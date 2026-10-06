#!/usr/bin/env bash
# Called only within a disposable integration namespace, after web:down.
set -euo pipefail
bin="${1:?binary}" root="${2:?fixture}" context="${3:?context}" namespace="${4:?namespace}"
k() { kubectl --context "$context" -n "$namespace" "$@"; }
k wait pvc/keep-data --for=jsonpath='{.status.phase}'=Bound --timeout=120s
uid="$(k get pvc keep-data -o jsonpath='{.metadata.uid}')"
test -n "$uid"
k wait --for=delete pod -l app --timeout=120s
"$bin" run web:up --root "$root"
k rollout status deployment/web --timeout=120s
marker="mudarro-owned-fixture-$$"
k exec deployment/web -- sh -c 'printf "%s" "$1" > /data/mudarro-probe' sh "$marker"
k exec deployment/web -- cat /data/mudarro-probe | grep -Fx "$marker"
"$bin" run web:down --root "$root"
test "$(k get deployment web -o jsonpath='{.spec.replicas}')" = 0
k wait --for=delete pod -l app --timeout=120s
test "$(k get pvc keep-data -o jsonpath='{.metadata.uid}')" = "$uid"
test "$(k get pvc keep-data -o jsonpath='{.status.phase}')" = Bound
"$bin" run web:up --root "$root"
k rollout status deployment/web --timeout=120s
test "$(k exec deployment/web -- cat /data/mudarro-probe)" = "$marker"
"$bin" run web:restart --root "$root"
k rollout status deployment/web --timeout=120s
test "$(k exec deployment/web -- cat /data/mudarro-probe)" = "$marker"
"$bin" run web:down --root "$root"
test "$(k get pvc keep-data -o jsonpath='{.metadata.uid}')" = "$uid"
echo "PVC persistence: OK; Bound, same UID, content survived down/up/restart ($namespace)"
