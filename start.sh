#!/bin/sh
set -e

echo "Run db migrations"
# if environment variable DB_SOURCE is not empty
if [ -z "$DB_SOURCE" ]; then
  echo "DB_SOURCE is empty"
  exit 1
fi
/app/migrate -path /app/migration -database "$DB_SOURCE" -verbose up

echo "Run app"
exec "$@"