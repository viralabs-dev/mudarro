#!/usr/bin/env bash
set -euo pipefail
for i in $(seq 1 90); do printf 'Real central output line %03d\n' "$i"; sleep .01; done
printf '\033[2J\033[H\033]0;fixture-title\007CONTROL_SAFE\n'
printf 'owned stderr\n' >&2
printf 'output-finished\n'
