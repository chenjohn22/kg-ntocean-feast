#!/bin/sh
set -eu

: "${DB_HOST:?DB_HOST is required}"
: "${DB_PORT:?DB_PORT is required}"
: "${DB_USER:?DB_USER is required}"
: "${DB_PASSWORD:?DB_PASSWORD is required}"
: "${DB_NAME:?DB_NAME is required}"

backup_dir=/backups
retention_hours="${BACKUP_RETENTION_HOURS:-72}"
case "$retention_hours" in
  ''|*[!0-9]*) echo "BACKUP_RETENTION_HOURS must be a positive integer" >&2; exit 1 ;;
esac
if [ "$retention_hours" -le 0 ]; then
  echo "BACKUP_RETENTION_HOURS must be a positive integer" >&2
  exit 1
fi
timestamp="$(date '+%Y%m%d-%H%M%S')"
final_file="${backup_dir}/ntocean-feast-${timestamp}.sql.gz"
temporary_file="${final_file}.tmp"

mkdir -p "$backup_dir"
umask 027
export MYSQL_PWD="$DB_PASSWORD"

cleanup() {
  rm -f "$temporary_file"
}
trap cleanup EXIT HUP INT TERM

echo "Starting database backup: ${final_file}"
mysqldump \
  --host="$DB_HOST" \
  --port="$DB_PORT" \
  --user="$DB_USER" \
  --single-transaction \
  --quick \
  --routines \
  --triggers \
  --no-tablespaces \
  --set-gtid-purged=OFF \
  "$DB_NAME" lottery_codes registrations | gzip -c > "$temporary_file"

test -s "$temporary_file"
mv "$temporary_file" "$final_file"
trap - EXIT HUP INT TERM

# Only prune after a new backup has completed successfully.
find "$backup_dir" -type f -name 'ntocean-feast-*.sql.gz' -mmin "+$((retention_hours * 60))" -delete
echo "Database backup completed: ${final_file}"
