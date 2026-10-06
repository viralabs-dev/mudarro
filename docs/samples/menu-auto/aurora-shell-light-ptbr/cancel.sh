#!/usr/bin/env bash
printf 'cancel-ready\n'
trap '' TERM
sleep 30 &
printf '%s\n' $! > owned-child.pid
wait
