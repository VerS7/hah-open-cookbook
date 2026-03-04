#!/bin/sh

set -eu

SOURCE_DIR="${BACKUP_SOURCE_DIR:-/app/data}"
BACKUP_BASE_DIR="${DEPLOY_DATA_DIR:-/app}"
BACKUP_DIR="${BACKUP_BASE_DIR}/backups"

mkdir -p "$BACKUP_DIR"

DATE="$(date +%Y%m%d_%H%M%S)"

echo "Starting backup..."
tar -czf "$BACKUP_DIR/backup_$DATE.tar.gz" "$SOURCE_DIR"
find "$BACKUP_DIR" -name "*.tar.gz" -mtime +7 -delete
echo "Backup completed"
