#!/bin/sh

set -e

echo "run db migration..."
/app/migrate -path /app/migrations -database "$POSTGRES_DRIVER://$POSTGRES_USER:$POSTGRES_PASSWORD@$POSTGRES_HOST:$POSTGRES_PORT/$POSTGRES_DB?sslmode=$SSL_MODE" -verbose up

echo "start the app..."
exec "$@"