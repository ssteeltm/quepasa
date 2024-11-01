#!/bin/sh
set -e

echo "Container's IP address: `awk 'END{print $1}' /etc/hosts`"
echo "Working dir: `pwd`"

cp -rf /builder/assets ./
cp -rf /builder/views ./
cp -rf /builder/service ./

# cleanup builder folder for avoid cache on updates
rm -fr /builder

./service
exec "$@"
