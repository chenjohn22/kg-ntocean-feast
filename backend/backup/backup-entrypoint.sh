#!/bin/sh
set -eu

backup_hour="${BACKUP_HOUR:-3}"
case "$backup_hour" in
  ''|*[!0-9]*) echo "BACKUP_HOUR must be an integer from 0 to 23" >&2; exit 1 ;;
esac
if [ "$backup_hour" -lt 0 ] || [ "$backup_hour" -gt 23 ]; then
  echo "BACKUP_HOUR must be an integer from 0 to 23" >&2
  exit 1
fi

echo "Backup scheduler started; daily backup hour: $(printf '%02d' "$backup_hour"):00 ${TZ:-local time}"

while true; do
  today="$(date '+%Y%m%d')"
  current_hour="$(date '+%H')"
  current_hour_number="${current_hour#0}"
  [ -n "$current_hour_number" ] || current_hour_number=0

  # If the container starts after the scheduled time, create today's missing backup immediately.
  existing="$(find /backups -maxdepth 1 -type f -name "ntocean-feast-${today}-*.sql.gz" -print -quit 2>/dev/null || true)"
  if [ "$current_hour_number" -ge "$backup_hour" ] && [ -z "$existing" ]; then
    if ! /bin/sh /opt/backup/backup-db.sh; then
      echo "Database backup failed; retrying in 5 minutes" >&2
      sleep 300
      continue
    fi
  fi

  sleep 60
done
