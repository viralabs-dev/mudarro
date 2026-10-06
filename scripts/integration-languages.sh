#!/usr/bin/env bash
# Real Python and Go actions; no packages installed, disposable projects only.
set -euo pipefail
bin="$(realpath "${1:?binary}")"
root="$(mktemp -d '/tmp/mudarro languages.XXXXXX')"
mkdir "$root/tools"
ln -s "$bin" "$root/tools/mudarro"
export PATH="$root/tools:$PATH"
cleanup() { for lang in python go; do "$bin" run "app-$lang:down" --root "$root/$lang" >/dev/null 2>&1 || true; done; rm -rf -- "$root"; }
trap cleanup EXIT
mkdir "$root/python" "$root/go"
printf '' > "$root/python/requirements.txt"
cat > "$root/python/app.py" <<'PY'
import sys, time
if "--test" in sys.argv:
    assert sum([1, 2]) == 3
    print("python quality OK")
else:
    print("python ready", flush=True)
    while True: time.sleep(1)
PY
cat > "$root/go/go.mod" <<'GO'
module sample

go 1.23
GO
cat > "$root/go/main.go" <<'GO'
package main
import("fmt";"time")
func main(){fmt.Println("go ready");for{time.Sleep(time.Second)}}
GO
cat > "$root/go/main_test.go" <<'GO'
package main
import "testing"
func TestQuality(t *testing.T){if 1+2!=3{t.Fatal("unexpected result")};t.Log("go quality OK")}
GO
for lang in python go; do
  fixture="$root/$lang"
  "$bin" scan --root "$fixture" --json > "$fixture/scan.json"
  cat "$fixture/scan.json"
  "$bin" init --root "$fixture"
  # Python has no inferred entrypoint: provide explicit commands, as the CLI requires.
  python3 - "$fixture/mudarro.yaml" "$fixture/scan.json" "$lang" <<'PY'
import json,sys
p,scan,lang=sys.argv[1:]
c=json.load(open(scan))["config"]
s=c["services"][0]
if lang=="python":
    s["commands"].update(start={"args":["python3","app.py"],"group":"aplicacao"}, test={"args":["python3","app.py","--test"],"group":"qualidade"})
else:
    s["commands"]["test"]["args"]=["go","test","-v","./..."]
open(p,"w").write(json.dumps(c))
PY
  "$bin" generate --root "$fixture"
  before="$(sha256sum "$fixture/menu.sh")"
  "$bin" generate --root "$fixture"
  test "$before" = "$(sha256sum "$fixture/menu.sh")"
  bash "$fixture/.mudarro/scripts/app-$lang/test.sh" | grep "$lang quality OK"
  "$bin" run "app-$lang:up" --root "$fixture"
  "$bin" run "app-$lang:up" --root "$fixture"
  "$bin" run "app-$lang:status" --root "$fixture"
  for attempt in $(seq 1 30); do
    if "$bin" run "app-$lang:logs" --root "$fixture" | grep "$lang ready"; then break; fi
    if [[ "$attempt" = 30 ]]; then exit 1; fi
    sleep 1
  done
  "$bin" run "app-$lang:restart" --root "$fixture"
  "$bin" run "app-$lang:down" --root "$fixture"
  echo "$lang real generation/actions/supervisor: OK"
done
