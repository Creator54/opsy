# Backup PostgreSQL Database

Production database backup procedure. Run during maintenance window.

## Check Database Connection

Verify PostgreSQL is accessible:
```bash
pg_isready -h localhost -p 5432
```

## Create Backup Directory
```bash
mkdir -p /tmp/postgres/$(date +%Y-%m)
```

## Dump Database

Create compressed backup with timestamp:
```bash
pg_dump -h localhost -U postgres -d production \
  | gzip > /tmp/postgres/$(date +%Y-%m)/backup_$(date +%Y%m%d_%H%M%S).sql.gz
```

## Verify Backup

Check backup file exists and has content:
```bash
ls -lh /tmp/postgres/$(date +%Y-%m)/ | tail -1
```

## Cleanup Old Backups

Remove backups older than 30 days:
```bash
find /tmp/postgres -name "*.sql.gz" -mtime +30 -delete
```

## Test Backup (Optional)

Verify backup integrity:
```bash
gunzip -t /tmp/postgres/$(date +%Y-%m)/backup_*.sql.gz | tail -1
```
