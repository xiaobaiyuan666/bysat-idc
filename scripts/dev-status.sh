#!/bin/sh
set -eu

. "$(dirname "$0")/dev-common.sh"

profile=${1:-all}

for service in $(service_list "$profile"); do
  status_service "$service"
done
