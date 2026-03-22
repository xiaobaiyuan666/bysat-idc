#!/bin/sh
set -eu

. "$(dirname "$0")/dev-common.sh"

profile=${1:-all}
failed=0

for service in $(service_list "$profile"); do
  if ! stop_service "$service"; then
    failed=1
  fi
done

exit "$failed"
