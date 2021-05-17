#!/bin/sh
: ${WS_PORT="3000"} # Port on which to listen for control commands over websocket
: ${NG_PORT="3001"} # Port on which to listen for control commands over UDP
: ${MIN_PORT="20000"} # Minimum port for RTP allocation
: ${MAX_PORT="30000"} # Maximum port for RTP allocation
: ${CONTROL_IP="127.0.0.1"} # IP address to bind to for the control services.  This MUST NOT be exposed to insecure networks, as it is unauthenticated!

: ${CLOUD=""} # One of aws, azure, do, gcp, or empty
if [ "$CLOUD" != "" ]; then
   PROVIDER="-provider ${CLOUD}"
fi

: ${PRIVATE_IPV4="$(netdiscover -field privatev4 ${PROVIDER})"}
: ${PUBLIC_IPV4="$(netdiscover -field publicv4 ${PROVIDER})"}

: ${RTPENGINE_ARGS:="-f -i ${PRIVATE_IPV4}!${PUBLIC_IPV4} --listen-ng=${CONTROL_IP}:${NG_PORT} --listen-http=${CONTROL_IP}:${WS_PORT} --port-min 20000 --port-max 30000 --log-stderr"}

# If we were given arguments, run them instead
if [ $# -gt 0 ]; then
   exec "$@"
fi

# Run rtpengine
exec /usr/bin/rtpengine ${RTPENGINE_ARGS}
