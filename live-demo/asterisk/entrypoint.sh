#!/bin/sh
set -e

# Default values
: ${ASTERISK_ARGS:="-fp"}
: ${CHECK_FILE:=".asterisk-config"}

# If we were given arguments, run them instead
if [ $# -gt 0 ]; then
   exec "$@"
fi

while [ ! -e "/etc/asterisk/$CHECK_FILE" ]; do
    echo "configuration not yet available, sleep 5 seconds..."
    sleep 5
done

# Run Asterisk
exec /usr/sbin/asterisk ${ASTERISK_ARGS}
