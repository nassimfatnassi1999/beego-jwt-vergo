#!/bin/bash

DB_USER="core"
DB_PASS="Core1234!"
DB_NAME="core"

echo "=============================="
echo " 🔄 RESET DATABASE: $DB_NAME"
echo "=============================="

sudo mysql <<EOF
DROP DATABASE IF EXISTS $DB_NAME;
CREATE DATABASE $DB_NAME;

-- Create user if not exists
CREATE USER IF NOT EXISTS '$DB_USER'@'localhost' IDENTIFIED BY '$DB_PASS';

-- Update password (just in case)
ALTER USER '$DB_USER'@'localhost' IDENTIFIED BY '$DB_PASS';

-- Grant privileges
GRANT ALL PRIVILEGES ON $DB_NAME.* TO '$DB_USER'@'localhost';
FLUSH PRIVILEGES;
EOF

echo "✔ DONE: database reset successful."
