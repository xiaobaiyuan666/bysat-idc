#!/bin/sh

project_root=$(CDPATH= cd -- "$(dirname "$0")/.." && pwd)
state_dir="$project_root/.codex-dev"
pid_dir="$state_dir/pids"
log_dir="$state_dir/logs"

mkdir -p "$pid_dir" "$log_dir"

service_list() {
  case "${1:-all}" in
    root)
      echo "root-backend"
      ;;
    rebuild)
      echo "rebuild-api rebuild-admin rebuild-portal"
      ;;
    all)
      echo "root-backend rebuild-api rebuild-admin rebuild-portal"
      ;;
    *)
      echo "Unknown profile: $1" >&2
      return 1
      ;;
  esac
}

service_port() {
  case "$1" in
    root-backend) echo "3000" ;;
    rebuild-api) echo "18080" ;;
    rebuild-admin) echo "5177" ;;
    rebuild-portal) echo "5178" ;;
    *)
      echo "Unknown service: $1" >&2
      return 1
      ;;
  esac
}

service_url() {
  case "$1" in
    root-backend) echo "http://127.0.0.1:3000/api/health" ;;
    rebuild-api) echo "http://127.0.0.1:18080/api/v1/health" ;;
    rebuild-admin) echo "http://127.0.0.1:5177/" ;;
    rebuild-portal) echo "http://127.0.0.1:5178/" ;;
    *)
      echo "Unknown service: $1" >&2
      return 1
      ;;
  esac
}

service_workdir() {
  case "$1" in
    root-backend) echo "$project_root" ;;
    rebuild-api|rebuild-admin|rebuild-portal) echo "$project_root/idc-finance" ;;
    *)
      echo "Unknown service: $1" >&2
      return 1
      ;;
  esac
}

service_command() {
  case "$1" in
    root-backend) echo "npm run dev:backend" ;;
    rebuild-api) echo "npm run dev:api:unix" ;;
    rebuild-admin) echo "npm run dev:admin" ;;
    rebuild-portal) echo "npm run dev:portal" ;;
    *)
      echo "Unknown service: $1" >&2
      return 1
      ;;
  esac
}

pidfile_of() {
  echo "$pid_dir/$1.pid"
}

logfile_of() {
  echo "$log_dir/$1.log"
}

read_pid() {
  pidfile=$(pidfile_of "$1")
  if [ -f "$pidfile" ]; then
    cat "$pidfile"
  fi
}

is_pid_running() {
  pid=$(read_pid "$1")
  if [ -n "${pid:-}" ] && kill -0 "$pid" >/dev/null 2>&1; then
    return 0
  fi
  return 1
}

cleanup_stale_pid() {
  pidfile=$(pidfile_of "$1")
  if [ -f "$pidfile" ] && ! is_pid_running "$1"; then
    rm -f "$pidfile"
  fi
}

port_listener_pid() {
  lsof -ti "tcp:$(service_port "$1")" -sTCP:LISTEN 2>/dev/null | head -n 1
}

is_port_busy() {
  [ -n "$(port_listener_pid "$1")" ]
}

wait_for_http() {
  url=$1
  timeout=${2:-30}
  i=0
  while [ "$i" -lt "$timeout" ]; do
    if curl -fsS "$url" >/dev/null 2>&1; then
      return 0
    fi
    sleep 1
    i=$((i + 1))
  done
  return 1
}

start_service() {
  service=$1
  cleanup_stale_pid "$service"

  if is_pid_running "$service"; then
    echo "$service already running (pid $(read_pid "$service"))"
    return 0
  fi

  if is_port_busy "$service"; then
    echo "$service not started: port $(service_port "$service") already in use by pid $(port_listener_pid "$service")" >&2
    return 1
  fi

  workdir=$(service_workdir "$service")
  command=$(service_command "$service")
  logfile=$(logfile_of "$service")
  pidfile=$(pidfile_of "$service")
  url=$(service_url "$service")

  : > "$logfile"
  (
    cd "$workdir" &&
    nohup sh -lc "$command" >> "$logfile" 2>&1 &
    echo $! > "$pidfile"
  )

  if wait_for_http "$url" 40; then
    echo "$service started at $url"
    return 0
  fi

  echo "$service failed to become ready. Check $logfile" >&2
  tail -n 40 "$logfile" >&2 || true
  return 1
}

stop_service() {
  service=$1
  cleanup_stale_pid "$service"

  if ! is_pid_running "$service"; then
    if is_port_busy "$service"; then
      echo "$service is listening on port $(service_port "$service") but is not managed by pidfile" >&2
      return 1
    fi
    echo "$service already stopped"
    return 0
  fi

  pid=$(read_pid "$service")
  kill "$pid" >/dev/null 2>&1 || true

  i=0
  while [ "$i" -lt 10 ]; do
    if ! kill -0 "$pid" >/dev/null 2>&1; then
      rm -f "$(pidfile_of "$service")"
      echo "$service stopped"
      return 0
    fi
    sleep 1
    i=$((i + 1))
  done

  echo "$service did not exit after SIGTERM (pid $pid)" >&2
  return 1
}

status_service() {
  service=$1
  cleanup_stale_pid "$service"
  port=$(service_port "$service")
  url=$(service_url "$service")

  if is_pid_running "$service"; then
    echo "$service running pid=$(read_pid "$service") port=$port url=$url"
    return 0
  fi

  if is_port_busy "$service"; then
    echo "$service unmanaged port-in-use pid=$(port_listener_pid "$service") port=$port url=$url"
    return 0
  fi

  echo "$service stopped port=$port url=$url"
}
