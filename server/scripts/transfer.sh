#!/usr/bin/env bash
CURRENT_DIR="$(cd "$(dirname "$0")" && pwd -P)"
PARENT_DIR="$(cd "$CURRENT_DIR"/.. && pwd)"
USER="root"
SERVER_IP="159.223.91.60"
rsync -chav --progress "$PARENT_DIR"/main.go "$USER"@"$SERVER_IP":/home/repos/cloud-morph-host/server/
rsync -chav --progress "$PARENT_DIR"/pkg "$USER"@"$SERVER_IP":/home/repos/cloud-morph-host/server/
rsync -chav --progress "$PARENT_DIR"/scripts "$USER"@"$SERVER_IP":/home/repos/cloud-morph-host/server/
