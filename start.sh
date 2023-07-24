#!/bin/sh

set -e

echo "run db migrations"
/app/migrate -path /app/migration -database "$DB_SOURCE" -verbose up

echo "Starting the application..."
./main