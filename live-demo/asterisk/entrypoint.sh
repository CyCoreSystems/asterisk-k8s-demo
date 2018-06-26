#!/bin/sh
set -e

# Default values
: ${ASTERISK_ARGS:="-fp"}

# If we were given arguments, run them instead
if [ $# -gt 0 ]; then
   exec "$@"
fi

# Run Asterisk
exec /usr/sbin/asterisk ${ASTERISK_ARGS}
