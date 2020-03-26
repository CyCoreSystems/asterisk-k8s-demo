#!/bin/sh
set -e

# Default values
: ${ASTERISK_ARGS:="-fp"}
: ${CHECK_FILE:=".asterisk-config"}

# If we were given arguments, run them instead
if [ $# -gt 0 ]; then
   exec "$@"
fi

if [ ! -e "/etc/asterisk/$CHECK_FILE" ]; then
   echo "configuration not yet available"
   exit 1
fi


# Run Asterisk
exec /usr/sbin/asterisk ${ASTERISK_ARGS}
