#!/bin/sh
: ${CLOUD=""} # One of aws, azure, do, gcp, or empty
if [ "$CLOUD" != "" ]; then
   PROVIDER="-provider ${CLOUD}"
fi

: ${PRIVATE_IPV4="$(netdiscover -field privatev4 ${PROVIDER})"}
: ${PUBLIC_IPV4="$(netdiscover -field publicv4 ${PROVIDER})"}

: ${RTPPROXY_ARGS:="-f -A ${PUBLIC_IPV4} -F -l ${PRIVATE_IPV4} -m 20000 -M 30000 -s udp:127.0.0.1:7722 -d INFO"}

# If we were given arguments, run them instead
if [ $# -gt 0 ]; then
   exec "$@"
fi

# Run rtpproxy
exec /usr/bin/rtpproxy ${RTPPROXY_ARGS}
