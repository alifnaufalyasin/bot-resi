#!/bin/sh

# https://github.com/docker-library/postgres/issues/151

set -x
POSTGRES="psql --username ${POSTGRES_USER}"

echo "Before"
echo "======"
$POSTGRES <<-SQL
\du
SQL

echo -n "[*] Creating database role: ${DB_USER}... "
$POSTGRES <<-SQL
IF NOT EXISTS (SELECT * FROM pg_user WHERE username = '${DB_USER}')
BEGIN
  CREATE ROLE ${DB_USER} LOGIN PASSWORD '${DB_PASSWORD}';
END;
SQL

echo -n "[*] Creating database ${DB_NAME}... "
$POSTGRES <<-SQL
SELECT ‘CREATE DATABASE ${DB_NAME} OWNER ${DB_USER} WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = '${DB_NAME}')\gexec
SQL

echo
echo "After"
echo "====="
$POSTGRES <<-SQL
\du
SQL