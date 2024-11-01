#!/bin/sh
set -e

echo "Container's IP address: `awk 'END{print $1}' /etc/hosts`"
echo "Working dir: `pwd`"

cp -rf /builder/assets ./
cp -rf /builder/views ./

# moving builded service for avoid cache on updates
mv -f /builder/service ./ 2> /dev/null

./service
exec "$@"
