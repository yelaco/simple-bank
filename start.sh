#!/bin/sh

set -e

echo "run db migration"
if [ -f /app/app.env ]; then
	. /app/app.env
else
	echo "/app/app.env does not exist or is not readable"
	exit 1
fi
/app/migrate -path /app/migration -database "$DB_SOURCE" -verbose up

echo "start the app"
exec "$@"
