#!/bin/sh
set -eu

. "$(dirname "$0")/dev-common.sh"

profile=${1:-all}
services=$(service_list "$profile")
child_pids=""
stopping=0

stop_children() {
  if [ "$stopping" -eq 1 ]; then
    return
  fi
  stopping=1

  for pid in $child_pids; do
    kill "$pid" >/dev/null 2>&1 || true
  done

  wait || true
}

trap 'stop_children; exit 0' INT TERM
trap 'stop_children' EXIT

for service in $services; do
  if is_port_busy "$service"; then
    echo "$service not started: port $(service_port "$service") already in use by pid $(port_listener_pid "$service")" >&2
    exit 1
  fi

  workdir=$(service_workdir "$service")
  command=$(service_command "$service")
  logfile=$(logfile_of "$service")
  : > "$logfile"

  (
    cd "$workdir" &&
    sh -lc "$command" 2>&1 |
      awk -v prefix="[$service] " '{ print prefix $0; fflush(); }' |
      tee "$logfile"
  ) &

  child_pids="$child_pids $!"
done

for service in $services; do
  if ! wait_for_http "$(service_url "$service")" 40; then
    echo "$service failed to become ready. Check $(logfile_of "$service")" >&2
    exit 1
  fi
done

echo "All services ready for profile: $profile"
echo "Press Ctrl-C to stop."

wait
