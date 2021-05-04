#!/bin/sh
: ${CLOUD=""} # One of aws, azure, do, gcp, or empty
if [ "$CLOUD" != "" ]; then
   PROVIDER="-provider ${CLOUD}"
fi

: ${PRIVATE_IPV4="$(netdiscover -field privatev4 ${PROVIDER})"}
: ${PUBLIC_IPV4="$(netdiscover -field publicv4 ${PROVIDER})"}

: ${RTPENGINE_ARGS:="-f -i ${PRIVATE_IPV4}!${PUBLIC_IPV4} --listen-ng=${PRIVATE_IPV4}:22222 --port-min 20000 --port-max 30000 --log-stderr"}

# If we were given arguments, run them instead
if [ $# -gt 0 ]; then
   exec "$@"
fi

# Run rtpengine
exec /usr/bin/rtpengine ${RTPENGINE_ARGS}
