#!/bin/sh
set -e

# Default values
: ${ASTERISK_ARGS:="-fp"}

# If we were given arguments, run them instead
if [ $# -gt 0 ]; then
   exec "$@"
fi

if [ ! -e /etc/asterisk/manager.conf ]; then
   echo "configuration not available"
   exit 1
fi

# Run Asterisk
exec /usr/sbin/asterisk ${ASTERISK_ARGS}
